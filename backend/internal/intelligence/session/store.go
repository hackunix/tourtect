package session

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tourtect/backend/internal/intelligence/model"
)

var ErrNotFound = errors.New("assistant session not found")

type Store interface {
	Get(context.Context, string) (*model.Session, error)
	Put(context.Context, *model.Session, time.Duration) error
	Delete(context.Context, string) error
}

type memoryItem struct {
	session   *model.Session
	expiresAt time.Time
}

type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]memoryItem
}

func NewMemoryStore() *MemoryStore { return &MemoryStore{items: make(map[string]memoryItem)} }

func (s *MemoryStore) Get(_ context.Context, id string) (*model.Session, error) {
	s.mu.RLock()
	item, ok := s.items[id]
	s.mu.RUnlock()
	if !ok || time.Now().After(item.expiresAt) {
		if ok {
			s.mu.Lock()
			delete(s.items, id)
			s.mu.Unlock()
		}
		return nil, ErrNotFound
	}
	b, _ := json.Marshal(item.session)
	var cloned model.Session
	_ = json.Unmarshal(b, &cloned)
	return &cloned, nil
}

func (s *MemoryStore) Put(_ context.Context, value *model.Session, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	var cloned model.Session
	if err := json.Unmarshal(b, &cloned); err != nil {
		return err
	}
	s.mu.Lock()
	s.items[value.ID] = memoryItem{session: &cloned, expiresAt: time.Now().Add(ttl)}
	s.mu.Unlock()
	return nil
}

func (s *MemoryStore) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}

// RedisStore intentionally uses the small RESP subset needed for session TTLs,
// avoiding a second Redis abstraction in the modular monolith.
type RedisStore struct {
	addr, password, prefix string
	dialTimeout            time.Duration
}

func NewRedisStore(addr, password string) *RedisStore {
	return &RedisStore{addr: addr, password: password, prefix: "tourtect:assistant:session:", dialTimeout: 2 * time.Second}
}

func (s *RedisStore) Get(ctx context.Context, id string) (*model.Session, error) {
	resp, err := s.command(ctx, "GET", s.prefix+id)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, ErrNotFound
	}
	b, ok := resp.([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected redis response")
	}
	var value model.Session
	if err := json.Unmarshal(b, &value); err != nil {
		return nil, fmt.Errorf("decode assistant session: %w", err)
	}
	return &value, nil
}

func (s *RedisStore) Put(ctx context.Context, value *model.Session, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = s.command(ctx, "SETEX", s.prefix+value.ID, strconv.FormatInt(max(1, int64(ttl.Seconds())), 10), string(b))
	return err
}

func (s *RedisStore) Delete(ctx context.Context, id string) error {
	resp, err := s.command(ctx, "DEL", s.prefix+id)
	if err != nil {
		return err
	}
	if n, ok := resp.(int64); !ok || n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *RedisStore) command(ctx context.Context, args ...string) (any, error) {
	d := net.Dialer{Timeout: s.dialTimeout}
	conn, err := d.DialContext(ctx, "tcp", s.addr)
	if err != nil {
		return nil, fmt.Errorf("redis unavailable: %w", err)
	}
	defer conn.Close()
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}
	r := bufio.NewReader(conn)
	if s.password != "" {
		if err := writeRESP(conn, "AUTH", s.password); err != nil {
			return nil, err
		}
		if _, err := readRESP(r); err != nil {
			return nil, err
		}
	}
	if err := writeRESP(conn, args...); err != nil {
		return nil, err
	}
	return readRESP(r)
}

func writeRESP(w io.Writer, args ...string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "*%d\r\n", len(args))
	for _, arg := range args {
		fmt.Fprintf(&b, "$%d\r\n%s\r\n", len(arg), arg)
	}
	_, err := io.WriteString(w, b.String())
	return err
}

func readRESP(r *bufio.Reader) (any, error) {
	prefix, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimSuffix(strings.TrimSuffix(line, "\n"), "\r")
	switch prefix {
	case '+':
		return line, nil
	case '-':
		return nil, errors.New(line)
	case ':':
		return strconv.ParseInt(line, 10, 64)
	case '$':
		n, err := strconv.Atoi(line)
		if err != nil {
			return nil, err
		}
		if n == -1 {
			return nil, nil
		}
		buf := make([]byte, n+2)
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		return buf[:n], nil
	default:
		return nil, fmt.Errorf("unsupported redis response prefix %q", prefix)
	}
}
