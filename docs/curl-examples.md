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
      "durationDays": [5, 15, 30, 60, 90, 180],
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
      "durationDays": [5, 15, 30, 60, 90, 180],
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

## User Bot Subscription APIs

### Get User's Subscribed Bots
```bash
curl -X POST http://localhost:8888/api/user/bots \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
  }'
```

### Subscribe to a Bot
```bash
curl -X POST http://localhost:8888/api/user/bots/subscribe \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb",
    "bot_id": "1",
    "duration_days": 5
  }'
```

### Unsubscribe from a Bot
```bash
curl -X POST http://localhost:8888/api/user/bots/unsubscribe \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "bot_id": "1"
  }'
```

### Pretty Print Response
```bash
curl -X POST http://localhost:8888/api/user/bots \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
  }' | jq
```

## Expected Response for Get User Bots
```json
{
  "message": "success",
  "data": [
    {
      "id": "1",
      "name": "BOT STAR",
      "iconLetter": "S",
      "riskLevel": "Very High",
      "durationDays": [5, 15, 30, 60, 90, 180],
      "expectedReturnPercent": 15,
      "aprDisplay": "15% (total over 5 days)",
      "minInvestment": 10,
      "maxInvestment": 10000,
      "investmentRange": "$10 - $10,000",
      "subscribers": 10423,
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
    }
  ]
}
```

## Expected Response for Subscribe/Unsubscribe
```json
{
  "message": "Subscribed successfully",
  "data": {
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
    "subscribers": 10423,
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
  }
}
```

### Already Subscribed Response
```json
{
  "message": "Already subscribed",
  "data": {
    "id": "1",
    "name": "BOT STAR",
    ...
  }
}
```

### Not Subscribed Response (Unsubscribe)
```json
{
  "message": "Not subscribed to this bot"
}
```

## Get User Profile API

### Basic Request
```bash
curl -X POST http://localhost:8888/api/user/profile \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb"
  }'
```

### Pretty Print Response
```bash
curl -X POST http://localhost:8888/api/user/profile \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
  }' | jq
```

### Production URL Example
```bash
curl -X POST https://be.wataros.io/api/user/profile \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"
  }'
```

## Expected Response for Get Profile
```json
{
  "message": "success",
  "data": {
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "referral_code": "ABC12345",
    "invite_code": "",
    "wata_reward": 50,
    "wata_balance": "1000.5",
    "usdt_balance": "500.25",
    "role": "user",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### User Not Found Response
```json
{
  "error_code": "0202",
  "message": "User not found"
}
```

## Deposit API

### Deposit WATA
```bash
curl -X POST http://localhost:8888/api/user/deposit \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb",
    "currency": "wata",
    "amount": "100.5",
    "tx_hash": "0x1234567890abcdef..."
  }'
```

### Deposit USDT
```bash
curl -X POST http://localhost:8888/api/user/deposit \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb",
    "currency": "usdt",
    "amount": "500.25",
    "tx_hash": "0xabcdef1234567890..."
  }'
```

### Pretty Print Response
```bash
curl -X POST http://localhost:8888/api/user/deposit \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb",
    "currency": "wata",
    "amount": "100.5"
  }' | jq
```

## Withdraw API

### Withdraw WATA
```bash
curl -X POST http://localhost:8888/api/user/withdraw \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb",
    "currency": "wata",
    "amount": "50.25",
    "tx_hash": "0x9876543210fedcba..."
  }'
```

### Withdraw USDT
```bash
curl -X POST http://localhost:8888/api/user/withdraw \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb",
    "currency": "usdt",
    "amount": "200.75",
    "tx_hash": "0xfedcba0987654321..."
  }'
```

### Pretty Print Response
```bash
curl -X POST http://localhost:8888/api/user/withdraw \
  -H "Content-Type: application/json" \
  -d '{
    "address": "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb",
    "currency": "usdt",
    "amount": "200.75"
  }' | jq
```

## Expected Response for Deposit/Withdraw
```json
{
  "message": "Deposit successful",
  "data": {
    "type": "deposit",
    "currency": "wata",
    "amount": "100.5",
    "balance_before": "0",
    "balance_after": "100.5",
    "status": "completed",
    "tx_hash": "0x1234567890abcdef...",
    "created_at": "2025-12-01T16:30:00+07:00"
  }
}
```

### Insufficient Balance Response (Withdraw)
```json
{
  "error_code": "0302",
  "message": "insufficient balance"
}
```

### Invalid Currency Response
```json
{
  "error_code": "0300",
  "message": "invalid currency. Must be 'wata' or 'usdt'"
}
```

### Invalid Amount Response
```json
{
  "error_code": "0301",
  "message": "invalid amount"
}
```

