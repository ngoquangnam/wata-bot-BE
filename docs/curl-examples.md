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
    "wata_reward": 50,
    "role": "user"
  }
}
```

## Get Bots API

### Basic Request
```bash
curl -X GET http://localhost:8888/api/bots
```

### Pretty Print Response
```bash
curl -X GET http://localhost:8888/api/bots | jq
```

### With Verbose Output
```bash
curl -v -X GET http://localhost:8888/api/bots
```

### Windows PowerShell
```powershell
Invoke-RestMethod -Uri "http://localhost:8888/api/bots" -Method GET
```

### Production URL Example
```bash
curl -X GET https://be.wataros.io/api/bots
```

## Expected Response for Get Bots
```json
{
  "message": "success",
  "data": [
    {
      "id": "1",
      "name": "BOT STAR",
      "iconLetter": "S",
      "riskLevel": "Very High",
      "durationDays": 5,
      "expectedReturnPercent": 15,
      "aprDisplay": "15% (total over 5 days)",
      "minInvestment": 10,
      "maxInvestment": 10000,
      "investmentRange": "$10 - $10,000",
      "subscribers": 10422,
      "author": "IYI Velocity Pro",
      "description": "Ultra-short 5-day investment package with 15% total return. Perfect for investors seeking quick capital turnover with very high risk.",
      "isActive": true,
      "metrics": {
        "lockupPeriod": "5 days",
        "expectedReturn": "15%",
        "minInvestment": "$10",
        "maxInvestment": "$10,000",
        "roi30d": "16.49%",
        "winRate": "79.45%",
        "tradingPair": "BTCUSDT",
        "totalTrades": 949,
        "pnl30d": 122840.71
      }
    },
    {
      "id": "2",
      "name": "BOT MINER",
      "iconLetter": "M",
      "riskLevel": "Very High",
      "durationDays": 15,
      "expectedReturnPercent": 25,
      "aprDisplay": "25% (total over 15 days)",
      "minInvestment": 10,
      "maxInvestment": 10000,
      "investmentRange": "$10 - $10,000",
      "subscribers": 8921,
      "author": "IYI Velocity Pro",
      "description": "15-day mid-term package offering 25% total profit. Balanced duration and attractive returns.",
      "isActive": true,
      "metrics": {
        "lockupPeriod": "15 days",
        "expectedReturn": "25%",
        "minInvestment": "$10",
        "maxInvestment": "$10,000",
        "roi30d": "26.12%",
        "winRate": "81.23%",
        "tradingPair": "BTCUSDT, ETHUSDT",
        "totalTrades": 1823,
        "pnl30d": 289451.30
      }
    }
  ]
}
```

### Empty Response (No Bots Found)
```json
{
  "message": "success",
  "data": []
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

