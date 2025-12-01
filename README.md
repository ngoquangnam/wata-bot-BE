# Wata Bot Backend

Backend API service built with go-zero framework.

## Prerequisites

- Go 1.18 or higher
- go-zero framework
- MySQL 5.7+ or MariaDB 10.3+

## Installation

```bash
# Install dependencies
go mod download

# Or install go-zero CLI tool (optional)
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

## Project Structure

```
.
├── api/                    # API definition files
│   └── wata-bot.api
├── etc/                    # Configuration files
│   └── wata-bot-api.yaml
├── internal/              # Internal application code
│   ├── config/           # Configuration
│   ├── handler/          # HTTP handlers
│   ├── logic/            # Business logic
│   ├── model/            # Database models
│   ├── svc/              # Service context
│   └── types/            # Request/Response types
├── sql/                   # SQL migration files
│   └── schema.sql
├── wata-bot.go          # Main entry point
├── go.mod
└── README.md
```

## Quick Start

Xem file [START.md](START.md) để có hướng dẫn chi tiết từng bước.

### Development Environment

1. **Setup database cho dev:**
   ```bash
   mysql -u root -p < sql/schema.dev.sql
   ```

2. **Chạy trên dev (Windows PowerShell):**
   ```powershell
   .\scripts\run-dev.ps1
   ```

   **Hoặc (Linux/Mac):**
   ```bash
   chmod +x scripts/run-dev.sh
   ./scripts/run-dev.sh
   ```

   **Hoặc chạy trực tiếp:**
   ```bash
   go run wata-bot.go -f etc/wata-bot-api.dev.yaml
   ```

### Production Environment

1. **Setup database:**
   ```bash
   mysql -u root -p < sql/schema.sql
   ```

2. **Cấu hình (chọn một trong hai):**
   - Tạo file `.env` từ `.env.example` và chỉnh sửa
   - Hoặc chỉnh sửa `etc/wata-bot-api.yaml`

3. **Chạy ứng dụng:**
   ```bash
   go run wata-bot.go -f etc/wata-bot-api.yaml
   ```

## Running the Service

```bash
# Run the service
go run wata-bot.go -f etc/wata-bot-api.yaml

# Or build and run
go build -o wata-bot.exe wata-bot.go
./wata-bot.exe -f etc/wata-bot-api.yaml
```

## API Endpoints

- `GET /api/hello` - Hello endpoint

Example:
```bash
curl http://localhost:8888/api/hello?name=World
```

## Database Setup

1. Create database and tables:
```bash
mysql -u root -p < sql/schema.sql
```

## Configuration

The application supports configuration via environment variables (`.env` file) or YAML config file.

### Using .env file (Recommended)

1. Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

2. Edit `.env` file with your configuration:
```env
# Server Configuration
SERVER_HOST=0.0.0.0
SERVER_PORT=8888

# JWT Secret Key
JWT_SECRET=your-secret-key-change-in-production

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=wata_bot_app
DB_PASSWORD=your_password_here
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

3. The application will automatically load `.env` file on startup.

### Using YAML config file

Edit `etc/wata-bot-api.yaml` to configure the service:
- Database connection string
- JWT secret key
- Server host and port
- Log settings

**Note:** Environment variables will override YAML config values if both are set.

## Development

### Generate Code from API Definition

If you have goctl installed:

```bash
goctl api go -api api/wata-bot.api -dir . -style gozero
```

## License

MIT

