# Đăng ký, đăng nhập và Google OAuth

> Tách từ `system-design.md` — mục 3.12.

### 3.12 Đăng ký, đăng nhập và Google OAuth

Tourtect hỗ trợ ba trạng thái danh tính:

1. **Anonymous:** đọc nội dung, tìm kiếm, dùng SOS và tạo dữ liệu nháp cục bộ.
2. **Registered:** đăng bài, bình luận, vote, follow, đồng bộ saved list và nhận notification.
3. **Verified role:** moderator/editor hoặc đại diện doanh nghiệp đã qua quy trình xác minh riêng; đăng nhập Google không tự cấp vai trò này.

#### Đăng ký bằng email

- Người dùng nhập email, mật khẩu, display name, locale và chấp nhận Terms/Privacy/Community Guidelines thành các consent record có phiên bản.
- Gửi email verification bằng token dùng một lần, hash ở server và hết hạn sau 15–30 phút.
- Chỉ cho xuất bản nội dung sau khi xác minh email; vẫn cho chuẩn bị draft trong thời gian chờ.
- Mật khẩu tối thiểu 12 ký tự, cho phép password manager/paste, kiểm tra mật khẩu phổ biến/rò rỉ và hash bằng Argon2id với tham số có version.
- Quên mật khẩu dùng token một lần, vô hiệu hóa sau sử dụng; reset mật khẩu thu hồi các refresh session khác và gửi security notification.

#### Đăng nhập với Google

- Web dùng **Google Identity Services (GIS)**; Android dùng Credential Manager hoặc system browser theo OpenID Connect Authorization Code Flow với PKCE.
- Chỉ xin scope <code>openid email profile</code>. Đăng nhập không xin quyền Drive, Contacts, YouTube hoặc quyền đăng thay người dùng.
- Client tạo <code>state</code>, <code>nonce</code> và PKCE <code>code_verifier/code_challenge</code>; backend thực hiện code exchange và xác minh ID token.
- Backend kiểm tra chữ ký/JWKS, <code>iss</code>, <code>aud</code>, <code>exp</code>, <code>iat</code>, <code>nonce</code> và <code>email_verified</code>. Khóa liên kết bằng cặp <code>(issuer, sub)</code>, không dùng email làm định danh Google bất biến.
- Google access/refresh token không được lưu nếu chỉ dùng để đăng nhập. Authentication và authorization tới Google API là hai consent flow riêng.
- Redirect URI là allowlist chính xác theo môi trường; không nhận redirect do client tùy ý truyền.

#### Liên kết tài khoản và phiên

- Anonymous session được merge vào account sau đăng nhập theo preview rõ: saved item, draft và preference nào sẽ được chuyển; incident/private media không tự merge.
- Nếu email Google trùng account email/password đã xác minh, Tourtect yêu cầu người dùng chứng minh quyền kiểm soát account hiện tại trước khi link; không tự gộp chỉ dựa trên email.
- Một account có thể liên kết nhiều identity provider; phải giữ ít nhất một phương thức đăng nhập trước khi unlink.
- Web dùng session cookie <code>HttpOnly</code>, <code>Secure</code>, <code>SameSite=Lax/Strict</code> phù hợp; Android lưu refresh token bằng Android Keystore-backed encrypted storage. Access token sống ngắn, refresh token rotation và reuse detection.
- Trang Security hiển thị thiết bị/phiên gần đây, cho revoke từng phiên hoặc “đăng xuất tất cả”. Logout cục bộ phải xóa session; logout Google không đồng nghĩa xóa account Tourtect.
- Không dùng Google profile photo làm public avatar mặc định nếu người dùng chưa xác nhận phạm vi công khai.

~~~mermaid
sequenceDiagram
    actor U as Người dùng
    participant C as Web/Android
    participant I as Tourtect Identity
    participant G as Google Identity
    participant D as Account DB

    U->>C: Chọn Đăng nhập với Google
    C->>I: Tạo auth attempt + state/nonce/PKCE
    I-->>C: authorization URL đã allowlist
    C->>G: Authorization request
    G-->>C: code + state
    C->>I: callback(code, state, code_verifier)
    I->>I: Kiểm tra state và exchange code
    I->>G: Token/JWKS verification data
    G-->>I: ID token / keys
    I->>I: Verify iss/aud/exp/nonce/email_verified
    I->>D: Find/link bằng issuer + sub
    D-->>I: account + roles
    I-->>C: Tourtect session + rotated refresh credential
    C-->>U: Đăng nhập; hỏi trước khi merge anonymous data
~~~

---
