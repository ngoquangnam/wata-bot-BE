# CORS Debug Guide

## Các bước để kiểm tra và fix CORS

### 1. Restart Server
**QUAN TRỌNG**: Sau khi thay đổi code, bạn PHẢI restart server:

```bash
# Dừng server hiện tại (Ctrl+C)
# Sau đó chạy lại:
go run wata-bot.go -f etc/wata-bot-api.dev.yaml
```

### 2. Kiểm tra Logs
Sau khi restart, khi có request từ browser, bạn sẽ thấy log CORS trong console:
```
CORS Request - Origin: http://localhost:3000, Method: POST, Path: /auth/wallet-not-sign
CORS: Allowing origin http://localhost:3000
```

Nếu KHÔNG thấy log này, có nghĩa là:
- Server chưa được restart với code mới
- Hoặc middleware không được gọi

### 3. Test với curl để kiểm tra headers

#### Test preflight request (OPTIONS):
```bash
curl -X OPTIONS http://localhost:8888/auth/wallet-not-sign \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -v
```

**Kết quả mong đợi:**
```
< HTTP/1.1 204 No Content
< Access-Control-Allow-Origin: http://localhost:3000
< Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, PATCH
< Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With, Accept, Origin
< Access-Control-Allow-Credentials: true
```

#### Test actual request (POST):
```bash
curl -X POST http://localhost:8888/auth/wallet-not-sign \
  -H "Origin: http://localhost:3000" \
  -H "Content-Type: application/json" \
  -d '{"address":"0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"}' \
  -v
```

**Kết quả mong đợi:**
```
< HTTP/1.1 200 OK
< Access-Control-Allow-Origin: http://localhost:3000
< Access-Control-Allow-Credentials: true
```

### 4. Kiểm tra trong Browser

1. Mở **DevTools** (F12)
2. Vào tab **Network**
3. Gọi API từ frontend
4. Click vào request trong Network tab
5. Xem **Response Headers**:
   - Phải có `Access-Control-Allow-Origin: http://localhost:3000`
   - Phải có `Access-Control-Allow-Credentials: true`

### 5. Nếu vẫn gặp lỗi CORS

#### Kiểm tra Origin thực tế:
- Xem log trong console khi có request
- Origin có thể khác `http://localhost:3000` (ví dụ: `http://127.0.0.1:3000`)
- Nếu origin khác, thêm vào `wata-bot.go`:

```go
corsMiddleware := middleware.NewCorsMiddlewareWithOrigins([]string{
    "http://localhost:3000",
    "http://127.0.0.1:3000",  // Thêm nếu cần
})
```

#### Tạm thời cho phép tất cả origins (chỉ để test):
```go
// Trong wata-bot.go, thay đổi:
corsMiddleware := middleware.NewCorsMiddleware()  // Cho phép tất cả
```

#### Kiểm tra browser cache:
- Hard refresh: `Ctrl+Shift+R` (Windows) hoặc `Cmd+Shift+R` (Mac)
- Hoặc mở Incognito/Private window

### 6. Common Issues

#### Issue: "No 'Access-Control-Allow-Origin' header"
- **Nguyên nhân**: Origin không được phép hoặc middleware không được gọi
- **Giải pháp**: 
  - Kiểm tra log để xem origin thực tế
  - Đảm bảo server đã được restart
  - Thêm origin vào allowed list

#### Issue: "Credentials flag is true, but 'Access-Control-Allow-Origin' is '*"
- **Nguyên nhân**: Không thể dùng `*` với credentials
- **Giải pháp**: Đã được xử lý trong code - luôn set origin cụ thể khi có credentials

#### Issue: Preflight request fails
- **Nguyên nhân**: OPTIONS request không được xử lý
- **Giải pháp**: Đã được xử lý trong middleware - tự động handle OPTIONS requests

### 7. Test Script

Tạo file `test-cors.sh`:
```bash
#!/bin/bash

echo "Testing CORS preflight..."
curl -X OPTIONS http://localhost:8888/auth/wallet-not-sign \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -v

echo -e "\n\nTesting CORS actual request..."
curl -X POST http://localhost:8888/auth/wallet-not-sign \
  -H "Origin: http://localhost:3000" \
  -H "Content-Type: application/json" \
  -d '{"address":"0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"}' \
  -v
```

Chạy: `bash test-cors.sh`

