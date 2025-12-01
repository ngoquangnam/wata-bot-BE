# API Error Codes

Tất cả các lỗi từ API sẽ được trả về với format:

```json
{
  "error_code": "0001",
  "message": "invalid address format"
}
```

## Error Codes

### Validation Errors (0001-0099)

| Code | Message | Description |
|------|---------|-------------|
| 0001 | invalid address format | Địa chỉ wallet không đúng format (phải là 40 ký tự hex sau 0x) |
| 0002 | invalid signature | Signature không hợp lệ hoặc không khớp với message |
| 0003 | invalid message | Message không hợp lệ |

### Authentication Errors (0100-0199)

| Code | Message | Description |
|------|---------|-------------|
| 0100 | failed to generate tokens | Không thể tạo JWT tokens |

### Database Errors (0200-0299)

| Code | Message | Description |
|------|---------|-------------|
| 0200 | database error | Lỗi kết nối hoặc truy vấn database |
| 0201 | failed to create user | Không thể tạo user mới trong database |
| 0202 | failed to find user | Không thể tìm thấy user (sau khi tạo) |

### Server Errors (0500-0599)

| Code | Message | Description |
|------|---------|-------------|
| 0500 | internal server error | Lỗi server không xác định |

## HTTP Status Codes

- **400 Bad Request**: Validation errors, authentication errors, database errors (client-side issues)
- **500 Internal Server Error**: Server errors, unknown errors

## Examples

### Invalid Address Format
```bash
curl -X POST http://localhost:8888/auth/wallet-not-sign \
  -H "Content-Type: application/json" \
  -d '{"address": "invalid"}'
```

Response:
```json
{
  "error_code": "0001",
  "message": "invalid address format"
}
```

### Invalid Signature
```bash
curl -X POST http://localhost:8888/auth/wallet \
  -H "Content-Type: application/json" \
  -d '{
    "signature": "invalid",
    "message": "test"
  }'
```

Response:
```json
{
  "error_code": "0002",
  "message": "invalid signature"
}
```

### Database Error
```json
{
  "error_code": "0200",
  "message": "database error"
}
```

