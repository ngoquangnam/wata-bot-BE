# API Curl Examples

## Wallet Auth Not Sign API

### Basic Request (without invite code)
```bash
curl -X POST http://localhost:8888/auth/wallet-not-sign \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
  }'
```

### Request with Invite Code
```bash
curl -X POST http://localhost:8888/auth/wallet-not-sign \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "invite_code": "ABC12345"
  }'
```

### Pretty Print Response
```bash
curl -X POST http://localhost:8888/auth/wallet-not-sign \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "invite_code": "ABC12345"
  }' | jq
```

### Windows PowerShell
```powershell
Invoke-RestMethod -Uri "http://localhost:8888/auth/wallet-not-sign" `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"address":"0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb","invite_code":"ABC12345"}'
```

## Expected Response
```json
{
  "message": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "base64_encoded_token...",
    "expires_in": 31536000,
    "aib_reward": 50,
    "role": "user"
  }
}
```

## Error Responses

### Invalid Address Format
```json
{
  "error": "invalid address format"
}
```

### Database Error
```json
{
  "error": "database error"
}
```

