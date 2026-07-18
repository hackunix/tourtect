# WebSocket Realtime Protocol

> Tách từ `system-design.md` — mục 7.4.

### 7.4 WebSocket messages

Client control:

- <code>ptt.started</code>: role, sequence, audio format.
- Binary <code>MediaChunk</code> loại AUDIO_PCM trong thời gian giữ nút.
- <code>ptt.ended</code>: utterance_id, final media sequence, client timestamp.
- Binary <code>MediaChunk</code> loại CAMERA_JPEG cho context frame, không phải capture bằng chứng.
- <code>session.hold</code>, <code>session.resume</code>, <code>session.end</code>.
- <code>feedback.translation</code> và <code>feedback.price_insight</code>.

WebSocket dùng hai kiểu frame:

1. Text frame JSON cho control và server event.
2. Binary frame Protobuf cho media; không gửi PCM/JPEG trần và không base64 media trong JSON.

~~~protobuf
enum MediaKind {
  MEDIA_KIND_UNSPECIFIED = 0;
  AUDIO_PCM = 1;
  CAMERA_JPEG = 2;
}

message MediaChunk {
  uint32 protocol_version = 1;
  MediaKind kind = 2;
  string media_id = 3;
  optional string utterance_id = 4;
  uint64 sequence = 5;
  int64 captured_at_ms = 6;
  bool is_last = 7;
  bytes payload = 8;
}
~~~

Server trả <code>media.ack</code> chứa <code>media_id</code>, highest contiguous <code>sequence</code> và trạng thái accepted/duplicate/rejected. WebSocket giữ ordering trong một kết nối; sau reconnect, client chỉ gửi lại chunk chưa ACK trong buffer bộ nhớ. <code>kind</code> là discriminator bắt buộc để route PCM và JPEG.

Server event:

- <code>media.ack</code>
- <code>transcript.partial</code>
- <code>transcript.final</code>
- <code>translation.ready</code>
- <code>vision.observation</code>
- <code>price.candidate</code>
- <code>price.insight</code>
- <code>safety.escalation</code>
- <code>network.degraded</code>
- <code>session.expired</code>

Envelope:

~~~typescript
interface RealtimeEvent<T> {
  event_id: string;
  event_type: string;
  session_id: string;
  utterance_id?: string;
  sequence: number;
  occurred_at: string;
  payload: T;
  trace_id: string;
}
~~~

Yêu cầu protocol:

- Event idempotent theo <code>event_id</code>.
- Sequence của event và sequence của media là hai namespace riêng; media sequence tăng đơn điệu trong từng <code>media_id</code>.
- Client bỏ qua event có sequence cũ nhưng vẫn ACK.
- Backpressure làm giảm camera frame trước, không drop <code>ptt.ended</code>.
- Audio format được negotiate khi tạo session; pilot chuẩn hóa PCM mono 16-bit, 16 kHz.
- Reconnect dùng resume token; không resume camera/micro nếu app đang background.
