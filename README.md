# Automatic Auction Closing with Goroutines

## Description

This project implements an automatic auction closing system in Go using Goroutines.

It started as a lab task provided by the Full Cycle team, and I added the automatic auction closing feature requested for this implementation.

- Auctions are created via a REST API and stored in MongoDB.
- Users can be created via the REST API and the generated user ID can be reused in bid requests.
- When an auction is created, a background Goroutine is launched to monitor its duration.
- Once the configured duration expires, the Goroutine updates the auction status to **Closed** in the database automatically, with no manual intervention.
- Bids submitted to a closed or expired auction are silently rejected.

---

## Running with Docker Compose

1. Create a copy of the example file and add your environment variables in `cmd/auction/.env` and adjust the environment variables as needed:

```env
AUCTION_DURATION=60s
```

2. Start the full environment:

```zsh
docker compose up -d --build
```

Available services:

- Auction API: `http://localhost:8080`
- MongoDB: `localhost:27017`

---

## Environment Variables

All variables are configured in `cmd/auction/.env`:

| Variable                     | Description                                                | Default         |
| ---------------------------- | ---------------------------------------------------------- | --------------- |
| `AUCTION_DURATION`           | How long an auction stays open before auto-closing         | `5m` (if unset) |
| `AUCTION_INTERVAL`           | Interval used by the bid service to check expiry in memory | —               |
| `BATCH_INSERT_INTERVAL`      | Interval between batch bid inserts                         | —               |
| `MAX_BATCH_SIZE`             | Max number of bids per batch insert                        | —               |
| `MONGODB_URL`                | MongoDB connection string                                  | —               |
| `MONGODB_DB`                 | MongoDB database name                                      | —               |
| `MONGO_INITDB_ROOT_USERNAME` | MongoDB root username                                      | —               |
| `MONGO_INITDB_ROOT_PASSWORD` | MongoDB root password                                      | —               |

> **Note:** `AUCTION_DURATION` and `AUCTION_INTERVAL` should be kept in sync so both the auto-close goroutine and the bid service's in-memory check agree on the auction lifetime.

---

## API Endpoints

### Create User

```http
POST http://localhost:8080/user
Content-Type: application/json

{
  "name": "Pablo"
}
```

Example response:

```json
{
  "id": "generated-uuid",
  "name": "Pablo"
}
```

### Create Auction

```http
POST http://localhost:8080/auction
Content-Type: application/json

{
  "product_name": "iPhone 17 Pro Max",
  "category": "Electronics",
  "description": "Brand new iPhone 17 Pro Max sealed in box",
  "condition": 1
}
```

### List Auctions

```http
GET http://localhost:8080/auction
```

### Find Auction by ID

```http
GET http://localhost:8080/auction/{auctionId}
```

### Submit a Bid

```http
POST http://localhost:8080/bid
Content-Type: application/json

{
  "user_id": "<userId>",
  "auction_id": "<auctionId>",
  "amount": 1500.00
}
```

### Find Winning Bid

```http
GET http://localhost:8080/auction/winner/{auctionId}
```

---

## Ready-made Requests

Use [api/api.http](/Users/pablo/Documents/dev/go-expert/labs/go-concurrency/api/api.http) to run the main scenarios.

Run `Create User` first. The request file captures the returned `id` and reuses it automatically in the bid requests.

---

To run all tests:

```zsh
go test ./...
```

---

Repository: <https://github.com/pablorodrigovieira/go-concurrency>
