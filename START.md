# Hướng dẫn khởi động Wata Bot Backend

## Bước 1: Kiểm tra Prerequisites

Đảm bảo bạn đã cài đặt:
- ✅ Go 1.18+ (kiểm tra: `go version`)
- ✅ MySQL 5.7+ hoặc MariaDB 10.3+ (kiểm tra: `mysql --version`)

## Bước 2: Cài đặt Dependencies

```bash
# Cài đặt Go dependencies
go mod download
```

## Bước 3: Setup Database

### 3.1. Tạo database và tables

```bash
# Chạy SQL script để tạo database và bảng
mysql -u root -p < sql/schema.sql
```

Hoặc nếu bạn đã có user `wata_bot_app`:
```bash
mysql -u wata_bot_app -p < sql/schema.sql
```

### 3.2. Kiểm tra database đã được tạo

```bash
mysql -u root -p -e "USE wata_bot; SHOW TABLES;"
```

Bạn sẽ thấy bảng `user` được tạo.

## Bước 4: Cấu hình Environment Variables

### Option 1: Sử dụng .env file (Khuyến nghị)

1. Tạo file `.env` từ template:
```bash
# Windows PowerShell
Copy-Item .env.example .env

# Linux/Mac
cp .env.example .env
```

2. Chỉnh sửa file `.env` với thông tin của bạn:
```env
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8888

# JWT Secret Key (thay đổi thành secret key mạnh)
JWT_SECRET=your-secret-key-change-in-production

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=wata_bot_app
DB_PASSWORD=RcR6gdqnSZj8E
DB_NAME=wata_bot
DB_CHARSET=utf8mb4
DB_TIMEZONE=Asia/Ho_Chi_Minh

# Log Configuration
LOG_SERVICE_NAME=wata-bot-api
LOG_MODE=file
LOG_PATH=logs
LOG_LEVEL=info
LOG_COMPRESS=true
LOG_KEEP_DAYS=7
```

### Option 2: Sử dụng YAML config

Chỉnh sửa file `etc/wata-bot-api.yaml`:
```yaml
Database:
  DataSource: wata_bot_app:RcR6gdqnSZj8E@tcp(localhost:3306)/wata_bot?charset=utf8mb4&parseTime=true&loc=Asia%2FHo_Chi_Minh
```

**Lưu ý:** Nếu cả `.env` và YAML đều có config, `.env` sẽ override YAML.

## Bước 5: Chạy ứng dụng

### Development mode (khuyến nghị cho development)

```bash
go run wata-bot.go -f etc/wata-bot-api.yaml
```

### Production mode

```bash
# Build binary
go build -o wata-bot.exe wata-bot.go

# Chạy binary
./wata-bot.exe -f etc/wata-bot-api.yaml
```

Nếu thành công, bạn sẽ thấy:
```
Starting server at 0.0.0.0:8888...
```

## Bước 6: Kiểm tra API

### Test Hello endpoint

```bash
curl http://localhost:8888/api/hello?name=World
```

Response:
```json
{
  "message": "Hello, World!"
}
```

### Test Wallet Auth endpoint

```bash
curl -X POST http://localhost:8888/auth/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "signature": "0xe571ae6cb9072663d2be468b065871bf7217aed4e8f36c6aab7f3c99dde1656c2005df1232385a797557ae1f4daa6cf1f2dd5425b519f24eb5ae003c2558e21b1b",
    "message": "Please sign this message to confirm your account",
    "invite_code": "0x0000"
  }'
```

Response:
```json
{
  "message": "success",
  "data": {
    "access_token": "eyJhbGci...",
    "refresh_token": "...",
    "expires_in": 31536000,
    "aib_reward": 50,
    "role": "user"
  }
}
```

## Troubleshooting

### Lỗi kết nối database

```
Error: dial tcp 127.0.0.1:3306: connect: connection refused
```

**Giải pháp:**
1. Kiểm tra MySQL đang chạy: `mysql -u root -p`
2. Kiểm tra user và password trong `.env` hoặc `wata-bot-api.yaml`
3. Kiểm tra database `wata_bot` đã được tạo

### Lỗi port đã được sử dụng

```
Error: listen tcp :8888: bind: address already in use
```

**Giải pháp:**
1. Thay đổi port trong `.env`: `SERVER_PORT=8889`
2. Hoặc kill process đang sử dụng port 8888

### Lỗi không tìm thấy .env file

```
Warning: .env file not found, using default values or environment variables
```

**Giải pháp:**
- Không cần lo lắng, ứng dụng vẫn chạy được với config từ YAML file
- Hoặc tạo file `.env` từ `.env.example`

## API Endpoints

- `GET /api/hello` - Hello endpoint
- `POST /auth/wallet` - Wallet authentication

## Logs

Logs được lưu trong thư mục `logs/` (theo config trong `.env` hoặc YAML).

