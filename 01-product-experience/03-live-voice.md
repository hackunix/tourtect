# Tourtect Live — Live Voice

> Tách từ `system-design.md` — mục 3.4.

### 3.4 Tourtect Live (Live Voice)

#### Giao diện

- Header: ngôn ngữ du khách, trạng thái mạng, thời lượng, Hold và End.
- Khu transcript: câu gốc và bản dịch, tách theo vai.
- Hai nút lớn:
  - “Du khách nói”
  - “Người Việt nói”
- Nút loa/tai nghe, phát lại, báo dịch sai.
- Price insight drawer mặc định thu gọn.
- Nút SOS luôn thấy, không chỉ biểu đạt bằng màu đỏ.

#### Luồng một lượt nói

1. Người dùng giữ nút đúng vai.
2. App phát tín hiệu rung ngắn, bắt đầu lấy PCM.
3. PCM được xử lý VAD/noise reduction và stream dần.
4. Khi nhả nút, app gửi <code>ptt.ended</code>.
5. ASR trả transcript có confidence.
6. Translation lane tạo bản dịch và trả ngay cho client.
7. TTS hệ điều hành phát bản dịch theo audio route hiện tại; router có thể dùng server TTS khi cần.
8. Intelligence lane độc lập trích giá, món/dịch vụ, scam signal và critical safety signal.
9. Nếu có price candidate đủ rõ, Price Engine trả insight kín.
10. Nếu có red flag, Safety Engine ưu tiên escalation card nhưng không làm mất transcript/dịch.

#### Dual-lane sequence

~~~mermaid
sequenceDiagram
    actor Tourist as Du khách
    participant App as Android App
    participant RT as Realtime Gateway
    participant ASR as ASR Router
    participant MT as Translation Service
    participant Intel as Intelligence Extractor
    participant Price as Price Engine
    participant Safety as Safety Engine

    Tourist->>App: Giữ nút PTT và nói
    App->>RT: ptt.started + PCM chunks
    Tourist->>App: Nhả nút
    App->>RT: ptt.ended
    RT->>ASR: Finalize utterance
    ASR-->>RT: transcript.final + confidence

    par Lane dịch
        RT->>MT: source, target, transcript
        MT-->>RT: translation.ready
        RT-->>App: Text dịch
        App-->>Tourist: System TTS + caption
    and Lane intelligence
        RT->>Intel: Transcript đã chuẩn hóa
        Intel->>Safety: Critical/scam signals
        Safety-->>RT: safe action hoặc no-op
        Intel->>Price: PriceCandidate nếu có
        Price-->>RT: PriceInsight
        RT-->>App: Thẻ kín + haptic
    end
~~~

#### Quy tắc âm thanh

- PTT-only, không mở micro liên tục.
- Mỗi utterance chỉ có một <code>SpeakerRole</code>.
- Không tự phát price insight qua loa.
- Nếu TTS locale không tồn tại hoặc engine lỗi:
  1. Hiển thị caption chữ lớn.
  2. Cho phép người dùng đưa màn hình cho bên kia đọc.
  3. Dùng audio phrase đã đóng gói cho câu khẩn cấp đã duyệt.
- Khi dùng tai nghe, người dùng có thể bấm “Đọc riêng” để nghe giải thích giá.
