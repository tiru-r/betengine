# Betting API Documentation

## Overview
This is a RESTful API built in Go for managing a betting system. It allows users to place bets, settle bets, and check their balance. The service uses in-memory storage and ensures concurrency safety with mutexes. The API does not rely on external routing libraries, using the standard `net/http` package.

## Features
- **Place a Bet**: Users can place bets on events with specified odds and amounts.
- **Settle a Bet**: Admins can settle events, marking bets as won or lost and updating user balances.
- **Check Balance**: Users can query their current balance.
- **Concurrency Safety**: Uses mutexes to ensure thread-safe operations on in-memory data.
- **Error Handling**: Comprehensive error responses with appropriate HTTP status codes.
- **Logging**: Logs key actions (bet placement, settlement, errors) to stdout.

## Project Structure
```
betting-api/
├── main.go         # Entry point, sets up routes and server
├── api/
│   └── api.go      # HTTP handlers for API endpoints
├── models/
│   └── models.go   # Data structures (Bet, User)
├── storage/
│   └── storage.go  # In-memory storage with concurrency safety
├── service/
│   └── service.go  # Core logic  
└── README.md       # This documentation
```

## Setup and Running
### Prerequisites
- Go 1.24.1 or later

### Installation
1. Clone the repository or copy the code into a directory.
2. Navigate to the project directory:
   ```bash
   cd betengine
   ```
3. Run the application:
   ```bash
   go run .
   ```
4. The server will start on `http://localhost:8080`.

### Test Data
The API initializes with:
- Users: `user1` (balance: 1000.0), `user2` (balance: 500.0)
- Events: `event1` (open), `event2` (open)

## API Endpoints
### 1. Place a Bet
- **Endpoint**: `POST /api/bets/place`
- **Description**: Places a bet for a user on an event.
- **Request Body**:
  ```json
  {
    "user_id": "string",
    "event_id": "string",
    "odds": float64,
    "amount": float64
  }
  ```
- **Response**:
  - **201 Created**:
    ```json
    {
      "id": "string",
      "user_id": "string",
      "event_id": "string",
      "odds": float64,
      "amount": float64,
      "placed_at": "timestamp",
      "status": "pending"
    }
    ```
  - **400 Bad Request**: Invalid input or insufficient balance.
  - **404 Not Found**: User not found.
- **Example**:
  ```bash
  curl -X POST http://localhost:8080/api/bets/place -H "Content-Type: application/json" -d '{"user_id":"user1","event_id":"event1","odds":2.5,"amount":100}'
  ```

### 2. Settle a Bet
- **Endpoint**: `POST api/bets/settle`
- **Description**: Settles an event, updating all related bets and user balances.
- **Request Body**:
  ```json
  {
    "event_id": "string",
    "result": "win" | "lose"
  }
  ```
- **Response**:
  - **200 OK**:
    ```json
    {
      "message": "Event settled"
    }
    ```
  - **400 Bad Request**: Invalid event ID, result, or event already settled.
- **Example**:
  ```bash
  curl -X POST http://localhost:8080/api/bets/settle -H "Content-Type: application/json" -d '{"event_id":"event1","result":"win"}'
  ```

### 3. Check User Balance 
- **Endpoint**: `GET /api/bets/balance?user_id=user1`
- **Description**: Retrieves the balance for a specified user.
- **Response**:
  - **200 OK**:
    ```json
    {
      "balance": float64
    }
    ```
  - **404 Not Found**: User not found.
- **Example**:
  ```bash
  curl http://localhost:8080/api/balance?user_id=user1
  ```

## Error Handling
Errors are returned in the following format:
```json
{
  "error": "error message"
}
```
Common status codes:
- `400 Bad Request`: Invalid input or request body.
- `404 Not Found`: Resource (user, event) not found.
- `405 Method Not Allowed`: Incorrect HTTP method.
- `500 Internal Server Error`: Unexpected server error (rare due to in-memory storage).

## Concurrency
The API uses a `sync.RWMutex` to ensure thread-safe access to in-memory data:
- Write operations (placing bets, settling bets) acquire a write lock.
- Read operations (checking balance) acquire a read lock.

## Limitations
- **In-Memory Storage**: Data is not persisted; restarting the server resets all data.
- **No Authentication**: The API does not implement user authentication or authorization.
- **Basic Routing**: Uses standard `net/http` routing, which is less feature-rich than libraries like chi.

## Testing
You can test the API using `curl` or tools like Postman. Example workflow:
1. Check user balance: `GET /api/bets/balance?user_id=user1`
2. Place a bet: `POST /api/bets/place`
3. Settle the event: `POST api/bets/settle`
4. Check updated balance: `GET /api/bets/balance?user_id=user1`

For automated testing, consider using Go's `testing` package or tools like Postman collections.
