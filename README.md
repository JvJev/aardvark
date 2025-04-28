# Aardvark task

Project contains backend part (main part of the task) and frontend part. FE was not required but I built it for my own and evaluator's convenience. 

backend - GO

frontend - react typescript

## Features

- **Topic-based messaging**: Send and receive messages on specific topics
- **Auto-incrementing message IDs**: Each message gets a unique global ID
- **Automatic topic management**: Topics (URL endpoints) are created on-demand and cleaned up when no longer in use
- **Connection timeout**: Clients are disconnected after 30 seconds of inactivity with a timeout event. (after last message sent)
- **Concurrent handling**: Service handles multiple clients simultaneously. You can open few browsers and test


### Prerequisites

- Install GO
- Install Node.js (for frontend)

### Backend Setup

1. Clone the repository
2. Navigate to service directory
3. Run the backend service:
   ```
   go run main.go
   ```
   The service will start on port 8080 by default.

### Frontend Setup

1. Navigate to the frontend directory
2. Install dependencies:
   ```
   npm install
   ```
3. Start the development server:
   ```
   npm run start
   ```
   The frontend will be accessible at http://localhost:3000

## Testing

Run tests:
```
go test -v 
```
# Author
Jevgenij Voronov