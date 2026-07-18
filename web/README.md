# Tourtect Web

Next.js 16 App Router client cho community feed, place discovery, Price Check và Safety Assessment. Web chỉ gọi Go API; không truy cập Postgres hoặc tự thực thi logic Price/Safety.

## Development

Backend API và PostgreSQL phải chạy trước. Với cấu hình localhost mặc định:

```bash
npm install
npm run dev
```

Mở <http://localhost:3000>. Khi API ở host khác, cấu hình cả hai biến:

```bash
API_URL=http://backend-host:8080
NEXT_PUBLIC_API_URL=http://backend-host:8080
npm run dev
```

- `API_URL`: request từ Next.js Server Components.
- `NEXT_PUBLIC_API_URL`: request từ Client Components như composer, useful, save, comments, Price và Safety.
- Không đặt API key hoặc secret vào biến `NEXT_PUBLIC_*`.

## Quality gates

```bash
npm run lint
npm test
npm run build
npx playwright install chromium
npm run test:e2e
```

Runbook đầy đủ: [../docs/operations/frontend-backend-runbook.md](../docs/operations/frontend-backend-runbook.md).
