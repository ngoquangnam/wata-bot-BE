#### API
1. auth/wallet


##### INPUT
curl 'https://api.aibotiyi.com/auth/wallet' \
  -H 'content-type: application/json' \
  --data-raw '{"signature":"0xe571ae6cb9072663d2be468b065871bf7217aed4e8f36c6aab7f3c99dde1656c2005df1232385a797557ae1f4daa6cf1f2dd5425b519f24eb5ae003c2558e21b1b","message":"Please sign this message to confirm your account","invite_code":"0x0000"}'
##### OUTPUT
{"message":"success","data":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIweEQ2MDMzMTY4MDVlMzFCOTQyYjg4RGRlNDc2NmZGMDY2REQ2NzM5N0IiLCJleHAiOjE3OTU5NjE1NjIsImp0aSI6IjU5NjkxZGQ0LTA3ZDgtNGYyNS04ZWNkLTA0YmM1YzVlYzBkNiIsImlhdCI6MTc2NDQyNTU2MiwiaXNzIjoicHJvZC1haWJvdC1iYWNrZW5kLWlzc3VlciIsInN1YiI6ImF1dGgiLCJ1c2VyX2lkIjoiMmNjNzUwMDItMDU0Mi00MDAxLThmYWItMzhlNjZlMmJkZTM3IiwiYWRkcmVzcyI6IjB4RDYwMzMxNjgwNWUzMUI5NDJiODhEZGU0NzY2ZkYwNjZERDY3Mzk3QiIsInJlZmVycmFsX2NvZGUiOiJERDY3Mzk3QiIsInJvbGUiOiJ1c2VyIn0.7I-H4iaOfRQBoUYGzrK21Qwiwl0RPipZGRB9XmTUQOk","refresh_token":"jr8S6K33AdUJBTZln9LT9EN1I6qdHq7jHwZ+rIyHklgN4UE+GFOpioj8DPeMDg8ZLkddAhDnjPdlGIwn2byW7A==","expires_in":31536000,"aib_reward":50,"role":"user"}}


##### OUTPUT ERROR
{"error_code":"0001","message":"invalid address format"}

2. auth/wallet-not-sign

#### INPUT
 curl -X POST http://localhost:8888/auth/wallet-not-sign   -H "Content-Type: application/json"   -d '{
    "address": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
    "invite_code": "ABC12345"
  }'

  ##### OUTPUT
{"message":"success","data":{"access_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIweEQ2MDMzMTY4MDVlMzFCOTQyYjg4RGRlNDc2NmZGMDY2REQ2NzM5N0IiLCJleHAiOjE3OTU5NjE1NjIsImp0aSI6IjU5NjkxZGQ0LTA3ZDgtNGYyNS04ZWNkLTA0YmM1YzVlYzBkNiIsImlhdCI6MTc2NDQyNTU2MiwiaXNzIjoicHJvZC1haWJvdC1iYWNrZW5kLWlzc3VlciIsInN1YiI6ImF1dGgiLCJ1c2VyX2lkIjoiMmNjNzUwMDItMDU0Mi00MDAxLThmYWItMzhlNjZlMmJkZTM3IiwiYWRkcmVzcyI6IjB4RDYwMzMxNjgwNWUzMUI5NDJiODhEZGU0NzY2ZkYwNjZERDY3Mzk3QiIsInJlZmVycmFsX2NvZGUiOiJERDY3Mzk3QiIsInJvbGUiOiJ1c2VyIn0.7I-H4iaOfRQBoUYGzrK21Qwiwl0RPipZGRB9XmTUQOk","refresh_token":"jr8S6K33AdUJBTZln9LT9EN1I6qdHq7jHwZ+rIyHklgN4UE+GFOpioj8DPeMDg8ZLkddAhDnjPdlGIwn2byW7A==","expires_in":31536000,"aib_reward":50,"role":"user"}}
