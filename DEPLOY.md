# Hướng dẫn Deploy Wata Bot Backend lên Production Ubuntu

Hướng dẫn chi tiết để deploy ứng dụng Wata Bot Backend lên server Ubuntu production.

## Mục lục

1. [Prerequisites](#prerequisites)
2. [Chuẩn bị Server](#chuẩn-bị-server)
3. [Cài đặt Dependencies](#cài-đặt-dependencies)
4. [Setup Database](#setup-database)
5. [Cấu hình Ứng dụng](#cấu-hình-ứng-dụng)
6. [Build và Deploy](#build-và-deploy)
7. [Setup Systemd Service](#setup-systemd-service)
8. [Setup Nginx Reverse Proxy](#setup-nginx-reverse-proxy)
9. [SSL/TLS với Let's Encrypt](#ssltls-với-lets-encrypt)
10. [Firewall Configuration](#firewall-configuration)
11. [Monitoring và Logging](#monitoring-và-logging)
12. [Troubleshooting](#troubleshooting)

---

## Prerequisites

- Server Ubuntu 20.04 LTS hoặc cao hơn
- Quyền root hoặc sudo
- Domain name (nếu cần SSL)
- Tối thiểu 2GB RAM, 20GB disk space

---

## Chuẩn bị Server

### 1. Cập nhật hệ thống

```bash
sudo apt update
sudo apt upgrade -y
```

### 2. Tạo user cho ứng dụng (khuyến nghị)

```bash
# Tạo user mới
sudo adduser wata-bot

# Thêm user vào nhóm sudo (nếu cần)
sudo usermod -aG sudo wata-bot

# Chuyển sang user mới
su - wata-bot
```

### 3. Tạo thư mục cho ứng dụng

```bash
sudo mkdir -p /opt/wata-bot
sudo chown wata-bot:wata-bot /opt/wata-bot
```

---

## Cài đặt Dependencies

### 1. Cài đặt Go

```bash
# Tải Go (kiểm tra phiên bản mới nhất tại https://go.dev/dl/)
cd /tmp
wget https://go.dev/dl/go1.24.0.linux-amd64.tar.gz

# Xóa bản cũ nếu có
sudo rm -rf /usr/local/go

# Giải nén và cài đặt
sudo tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

# Thêm vào PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Kiểm tra cài đặt
go version
```

### 2. Cài đặt MySQL

```bash
# Cài đặt MySQL Server
sudo apt install mysql-server -y

# Bảo mật MySQL
sudo mysql_secure_installation

# Khởi động và enable MySQL
sudo systemctl start mysql
sudo systemctl enable mysql
```

### 3. Cài đặt Nginx (cho reverse proxy)

```bash
sudo apt install nginx -y
sudo systemctl start nginx
sudo systemctl enable nginx
```

---

## Setup Database

### 1. Tạo database và user

```bash
# Đăng nhập MySQL
sudo mysql -u root -p

# Trong MySQL prompt, chạy các lệnh sau:
```

```sql
-- Tạo database
CREATE DATABASE IF NOT EXISTS `wata_bot` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Tạo user cho ứng dụng (thay đổi password mạnh)
CREATE USER 'wata_bot_app'@'localhost' IDENTIFIED BY 'YOUR_STRONG_PASSWORD_HERE';

-- Cấp quyền
GRANT ALL PRIVILEGES ON wata_bot.* TO 'wata_bot_app'@'localhost';
FLUSH PRIVILEGES;

-- Kiểm tra
SHOW DATABASES;
EXIT;
```

### 2. Import schema

```bash
# Clone repository hoặc upload file schema.sql
cd /opt/wata-bot

# Nếu bạn đã clone repo:
mysql -u wata_bot_app -p wata_bot < sql/schema.sql

# Hoặc nếu bạn upload file:
mysql -u wata_bot_app -p wata_bot < /path/to/schema.sql

# Kiểm tra tables đã được tạo
mysql -u wata_bot_app -p -e "USE wata_bot; SHOW TABLES;"
```

---

## Cấu hình Ứng dụng

### 1. Clone hoặc upload code

```bash
cd /opt/wata-bot

# Option 1: Clone từ Git repository
git clone <your-repo-url> .

# Option 2: Upload code bằng SCP từ máy local
# scp -r /path/to/wata-bot-BE/* user@server:/opt/wata-bot/
```

### 2. Tạo file cấu hình production

```bash
cd /opt/wata-bot

# Tạo thư mục etc nếu chưa có
mkdir -p etc

# Tạo file cấu hình production
nano etc/wata-bot-api.prod.yaml
```

Nội dung file `etc/wata-bot-api.prod.yaml`:

```yaml
Name: wata-bot-api
Host: 127.0.0.1
Port: 8888

# JWT Secret Key - THAY ĐỔI THÀNH KEY MẠNH VÀ BẢO MẬT
JWTSecret: YOUR_STRONG_JWT_SECRET_KEY_HERE_CHANGE_THIS

# Database configuration
Database:
  DataSource: wata_bot_app:YOUR_DB_PASSWORD@tcp(localhost:3306)/wata_bot?charset=utf8mb4&parseTime=true&loc=Asia%2FHo_Chi_Minh

# Cache configuration (disabled)
Cache:
  - Host: localhost:6379
    Type: node
    Pass: ""
    DB: 0

# Log settings
Log:
  ServiceName: wata-bot-api
  Mode: file
  Path: /opt/wata-bot/logs
  Level: info
  Compress: true
  KeepDays: 30
  StackCooldownMillis: 100
```

**Lưu ý quan trọng:**
- Thay `YOUR_STRONG_JWT_SECRET_KEY_HERE_CHANGE_THIS` bằng JWT secret key mạnh (ít nhất 32 ký tự)
- Thay `YOUR_DB_PASSWORD` bằng password database đã tạo ở bước trên
- Đặt `Host: 127.0.0.1` để chỉ lắng nghe localhost (Nginx sẽ reverse proxy)

### 3. Tạo file .env (tùy chọn, nếu muốn override config)

```bash
nano /opt/wata-bot/.env
```

```env
# Server Configuration
SERVER_HOST=127.0.0.1
SERVER_PORT=8888

# JWT Secret Key
JWT_SECRET=YOUR_STRONG_JWT_SECRET_KEY_HERE_CHANGE_THIS

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=wata_bot_app
DB_PASSWORD=YOUR_DB_PASSWORD
DB_NAME=wata_bot
DB_CHARSET=utf8mb4
DB_TIMEZONE=Asia/Ho_Chi_Minh

# Log Configuration
LOG_SERVICE_NAME=wata-bot-api
LOG_MODE=file
LOG_PATH=/opt/wata-bot/logs
LOG_LEVEL=info
LOG_COMPRESS=true
LOG_KEEP_DAYS=30
```

### 4. Tạo thư mục logs

```bash
mkdir -p /opt/wata-bot/logs
chmod 755 /opt/wata-bot/logs
```

### 5. Set permissions

```bash
# Đảm bảo user wata-bot có quyền
sudo chown -R wata-bot:wata-bot /opt/wata-bot
```

---

## Build và Deploy

### 1. Build ứng dụng

```bash
cd /opt/wata-bot

# Download dependencies
go mod download

# Build binary cho Linux
# LƯU Ý: Phải build trên Linux server hoặc cross-compile
go build -o wata-bot wata-bot.go

# Nếu build từ máy khác (Windows/Mac), sử dụng cross-compile:
# GOOS=linux GOARCH=amd64 go build -o wata-bot wata-bot.go

# Kiểm tra binary đã được tạo
ls -lh wata-bot

# Kiểm tra file type (phải là Linux ELF binary)
file wata-bot
# Kết quả mong đợi: "ELF 64-bit LSB executable, x86-64" hoặc tương tự

# Set executable permission
chmod +x wata-bot

# Test chạy thử (Ctrl+C để dừng)
./wata-bot -f etc/wata-bot-api.prod.yaml
```

### 2. Tạo script khởi động

```bash
nano /opt/wata-bot/start.sh
```

Nội dung:

```bash
#!/bin/bash
cd /opt/wata-bot
./wata-bot -f etc/wata-bot-api.prod.yaml
```

```bash
chmod +x /opt/wata-bot/start.sh
```

---

## Setup Systemd Service

Tạo systemd service để ứng dụng tự động khởi động và quản lý như một service.

### 1. Tạo service file

```bash
sudo nano /etc/systemd/system/wata-bot.service
```

Nội dung:

```ini
[Unit]
Description=Wata Bot Backend API Service
After=network.target mysql.service
Requires=mysql.service

[Service]
Type=simple
User=wata-bot
Group=wata-bot
WorkingDirectory=/opt/wata-bot
ExecStart=/opt/wata-bot/wata-bot -f /opt/wata-bot/etc/wata-bot-api.prod.yaml
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=wata-bot

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/wata-bot/logs

# Resource limits
LimitNOFILE=65536
LimitNPROC=4096

[Install]
WantedBy=multi-user.target
```

### 2. Reload systemd và khởi động service

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service (tự động khởi động khi boot)
sudo systemctl enable wata-bot

# Khởi động service
sudo systemctl start wata-bot

# Kiểm tra status
sudo systemctl status wata-bot

# Xem logs
sudo journalctl -u wata-bot -f
```

### 3. Các lệnh quản lý service

```bash
# Khởi động
sudo systemctl start wata-bot

# Dừng
sudo systemctl stop wata-bot

# Khởi động lại
sudo systemctl restart wata-bot

# Xem status
sudo systemctl status wata-bot

# Xem logs
sudo journalctl -u wata-bot -n 100
sudo journalctl -u wata-bot -f

# Disable auto-start
sudo systemctl disable wata-bot
```

---

## Setup Nginx Reverse Proxy

### 1. Tạo Nginx config

```bash
sudo nano /etc/nginx/sites-available/wata-bot
```

Nội dung (thay `your-domain.com` bằng domain của bạn):

```nginx
server {
    listen 80;
    server_name be.wataros.io www.be.wataros.io;

    # Logging
    access_log /var/log/nginx/wata-bot-access.log;
    error_log /var/log/nginx/wata-bot-error.log;

    # Client body size limit
    client_max_body_size 10M;

    # Timeouts
    proxy_connect_timeout 60s;
    proxy_send_timeout 60s;
    proxy_read_timeout 60s;

    # Proxy settings
    location / {
        proxy_pass http://127.0.0.1:8888;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # CORS headers (nếu cần)
        add_header Access-Control-Allow-Origin *;
        add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS";
        add_header Access-Control-Allow-Headers "Authorization, Content-Type";
    }

    # Health check endpoint
    location /health {
        proxy_pass http://127.0.0.1:8888/api/hello;
        access_log off;
    }
}
```

### 2. Enable site

```bash
# Tạo symbolic link
sudo ln -s /etc/nginx/sites-available/wata-bot /etc/nginx/sites-enabled/

# Xóa default site (nếu không cần)
sudo rm /etc/nginx/sites-enabled/default

# Test Nginx config
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

### 3. Kiểm tra

```bash
# Test từ server
curl http://localhost:8888/api/hello?name=Test

# Test từ bên ngoài (nếu đã có domain)
curl http://your-domain.com/api/hello?name=Test
```

---

## SSL/TLS với Let's Encrypt

### 1. Cài đặt Certbot

```bash
sudo apt install certbot python3-certbot-nginx -y
```

### 2. Cấu hình SSL

```bash
# Chạy certbot (thay your-domain.com)
sudo certbot --nginx -d your-domain.com -d www.your-domain.com

# Certbot sẽ tự động:
# - Tạo SSL certificate
# - Cập nhật Nginx config
# - Setup auto-renewal
```

### 3. Kiểm tra auto-renewal

```bash
# Test renewal
sudo certbot renew --dry-run

# Certbot tự động setup cron job để renew
# Kiểm tra: sudo systemctl status certbot.timer
```

### 4. Cập nhật Nginx config sau SSL

Sau khi chạy certbot, file config sẽ được tự động cập nhật. Bạn có thể kiểm tra:

```bash
sudo nano /etc/nginx/sites-available/wata-bot
```

---

## Firewall Configuration

### 1. Cấu hình UFW

```bash
# Cho phép SSH (quan trọng!)
sudo ufw allow 22/tcp

# Cho phép HTTP và HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable

# Kiểm tra status
sudo ufw status
```

### 2. Hoặc sử dụng iptables (nếu không dùng UFW)

```bash
# Allow HTTP
sudo iptables -A INPUT -p tcp --dport 80 -j ACCEPT

# Allow HTTPS
sudo iptables -A INPUT -p tcp --dport 443 -j ACCEPT

# Save rules
sudo iptables-save > /etc/iptables/rules.v4
```

---

## Monitoring và Logging

### 1. Xem logs ứng dụng

```bash
# Logs từ systemd
sudo journalctl -u wata-bot -f

# Logs từ file
tail -f /opt/wata-bot/logs/access.log
tail -f /opt/wata-bot/logs/error.log

# Logs Nginx
sudo tail -f /var/log/nginx/wata-bot-access.log
sudo tail -f /var/log/nginx/wata-bot-error.log
```

### 2. Monitoring với systemctl

```bash
# Kiểm tra service đang chạy
sudo systemctl is-active wata-bot

# Kiểm tra service có lỗi không
sudo systemctl is-failed wata-bot
```

### 3. Monitoring disk space

```bash
# Kiểm tra disk usage
df -h

# Kiểm tra log size
du -sh /opt/wata-bot/logs/*
```

### 4. Setup log rotation (tùy chọn)

```bash
sudo nano /etc/logrotate.d/wata-bot
```

Nội dung:

```
/opt/wata-bot/logs/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 0644 wata-bot wata-bot
    sharedscripts
    postrotate
        systemctl reload wata-bot > /dev/null 2>&1 || true
    endscript
}
```

---

## Troubleshooting

### 1. Lỗi "cannot execute binary file: Exec format error"

Lỗi này xảy ra khi binary được build trên hệ điều hành hoặc kiến trúc khác (ví dụ: build trên Windows nhưng chạy trên Linux).

**Giải pháp:**

```bash
# 1. Kiểm tra file type của binary hiện tại
file /opt/wata-bot/wata-bot

# Nếu thấy "PE32" hoặc "Windows", đây là binary Windows - cần build lại

# 2. Xóa binary cũ (nếu có)
rm -f /opt/wata-bot/wata-bot
rm -f /opt/wata-bot/wata-bot.exe

# 3. Kiểm tra Go đã được cài đặt đúng
go version

# 4. Build lại binary cho Linux (QUAN TRỌNG: phải build trên Linux server)
cd /opt/wata-bot
go mod download
go build -o wata-bot wata-bot.go

# Hoặc nếu cần cross-compile từ máy khác, sử dụng:
# GOOS=linux GOARCH=amd64 go build -o wata-bot wata-bot.go

# 5. Kiểm tra binary mới
file /opt/wata-bot/wata-bot
# Phải thấy: "ELF 64-bit LSB executable, x86-64" hoặc tương tự

# 6. Set executable permission
chmod +x /opt/wata-bot/wata-bot

# 7. Test chạy
./wata-bot -f etc/wata-bot-api.prod.yaml
```

**Lưu ý quan trọng:**
- **Phải build trên Linux server** hoặc cross-compile với `GOOS=linux GOARCH=amd64`
- Không thể chạy file `.exe` (Windows) trên Linux
- Kiểm tra kiến trúc: `uname -m` (phải là x86_64 hoặc amd64)
- Nếu server là ARM, sử dụng: `GOOS=linux GOARCH=arm64 go build -o wata-bot wata-bot.go`

### 2. Service không khởi động

```bash
# Kiểm tra logs
sudo journalctl -u wata-bot -n 50

# Kiểm tra file binary có tồn tại
ls -la /opt/wata-bot/wata-bot

# Kiểm tra permissions
ls -la /opt/wata-bot/

# Kiểm tra file type (phải là Linux ELF binary)
file /opt/wata-bot/wata-bot

# Test chạy thủ công
cd /opt/wata-bot
./wata-bot -f etc/wata-bot-api.prod.yaml
```

### 3. Lỗi kết nối database

```bash
# Kiểm tra MySQL đang chạy
sudo systemctl status mysql

# Test kết nối database
mysql -u wata_bot_app -p -h localhost wata_bot

# Kiểm tra config trong file YAML
cat /opt/wata-bot/etc/wata-bot-api.prod.yaml
```

### 4. Port đã được sử dụng

```bash
# Kiểm tra port 8888
sudo netstat -tlnp | grep 8888
# hoặc
sudo ss -tlnp | grep 8888

# Kill process nếu cần
sudo kill -9 <PID>
```

### 5. Nginx không proxy được

```bash
# Test Nginx config
sudo nginx -t

# Kiểm tra Nginx đang chạy
sudo systemctl status nginx

# Xem Nginx error logs
sudo tail -f /var/log/nginx/error.log

# Test kết nối từ Nginx đến backend
curl http://127.0.0.1:8888/api/hello
```

### 6. Permission denied

```bash
# Kiểm tra ownership
ls -la /opt/wata-bot/

# Fix ownership
sudo chown -R wata-bot:wata-bot /opt/wata-bot

# Fix permissions
sudo chmod +x /opt/wata-bot/wata-bot
sudo chmod 755 /opt/wata-bot/logs
```

### 7. Out of memory

```bash
# Kiểm tra memory
free -h

# Kiểm tra processes
top
# hoặc
htop

# Nếu cần, tăng swap
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

---

## Backup và Recovery

### 1. Backup Database

```bash
# Tạo script backup
nano /opt/wata-bot/backup-db.sh
```

Nội dung:

```bash
#!/bin/bash
BACKUP_DIR="/opt/wata-bot/backups"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

mysqldump -u wata_bot_app -p'YOUR_DB_PASSWORD' wata_bot > $BACKUP_DIR/wata_bot_$DATE.sql

# Xóa backup cũ hơn 30 ngày
find $BACKUP_DIR -name "*.sql" -mtime +30 -delete

echo "Backup completed: $BACKUP_DIR/wata_bot_$DATE.sql"
```

```bash
chmod +x /opt/wata-bot/backup-db.sh

# Test backup
/opt/wata-bot/backup-db.sh
```

### 2. Setup cron job cho backup tự động

```bash
# Mở crontab
crontab -e

# Thêm dòng sau để backup mỗi ngày lúc 2 giờ sáng
0 2 * * * /opt/wata-bot/backup-db.sh >> /opt/wata-bot/logs/backup.log 2>&1
```

### 3. Restore Database

```bash
# Restore từ backup
mysql -u wata_bot_app -p wata_bot < /opt/wata-bot/backups/wata_bot_YYYYMMDD_HHMMSS.sql
```

---

## Update và Maintenance

### 1. Update ứng dụng

```bash
# Dừng service
sudo systemctl stop wata-bot

# Backup database trước khi update
/opt/wata-bot/backup-db.sh

# Pull code mới hoặc upload code mới
cd /opt/wata-bot
# git pull origin main  # nếu dùng Git
# hoặc upload code mới

# Rebuild
go mod download
GOOS=linux GOARCH=amd64 go build -o wata-bot wata-bot.go

# Khởi động lại
sudo systemctl start wata-bot

# Kiểm tra
sudo systemctl status wata-bot
```

### 2. Rollback

```bash
# Dừng service
sudo systemctl stop wata-bot

# Restore binary cũ (nếu đã backup)
cp /opt/wata-bot/wata-bot.backup /opt/wata-bot/wata-bot

# Hoặc rebuild từ commit cũ
cd /opt/wata-bot
git checkout <old-commit>
go build -o wata-bot wata-bot.go

# Khởi động lại
sudo systemctl start wata-bot
```

---

## Security Best Practices

1. **Đổi tất cả password mặc định**
   - Database password
   - JWT secret key
   - User passwords

2. **Giới hạn quyền truy cập**
   - Chỉ cho phép user `wata-bot` truy cập thư mục ứng dụng
   - Sử dụng firewall để giới hạn truy cập

3. **Bảo mật MySQL**
   - Không expose MySQL ra ngoài
   - Sử dụng strong password
   - Chỉ cấp quyền cần thiết

4. **SSL/TLS**
   - Luôn sử dụng HTTPS trong production
   - Setup auto-renewal cho SSL certificate

5. **Regular updates**
   - Cập nhật hệ thống thường xuyên
   - Cập nhật ứng dụng khi có bản mới

6. **Monitoring**
   - Theo dõi logs thường xuyên
   - Setup alerts cho errors

---

## Checklist Deploy

- [ ] Server Ubuntu đã được cập nhật
- [ ] Go đã được cài đặt và cấu hình
- [ ] MySQL đã được cài đặt và cấu hình
- [ ] Database và user đã được tạo
- [ ] Schema đã được import
- [ ] Ứng dụng đã được build thành công
- [ ] File config production đã được tạo và cấu hình đúng
- [ ] Systemd service đã được tạo và enable
- [ ] Service đang chạy và không có lỗi
- [ ] Nginx đã được cấu hình và reload
- [ ] SSL certificate đã được cài đặt (nếu có domain)
- [ ] Firewall đã được cấu hình
- [ ] Backup script đã được setup
- [ ] Logs đang được ghi đúng
- [ ] API đã được test và hoạt động

---

## Liên hệ và Hỗ trợ

Nếu gặp vấn đề trong quá trình deploy, vui lòng:
1. Kiểm tra logs: `sudo journalctl -u wata-bot -n 100`
2. Kiểm tra file config
3. Kiểm tra database connection
4. Xem phần Troubleshooting ở trên

---

**Chúc bạn deploy thành công! 🚀**

