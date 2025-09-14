1. get : go get github.com/mattn/go-sqlite3
2. get : go get github.com/go-chi/chi/v5

BON Rewards Service
A simple backend service, written in Go, that simulates a user rewards program. A user receives a mock gift card reward only if they’ve paid their last 3 consecutive bills on time.

Tech Stack
Go: Language
SQLite: Database
Chi: HTTP Router

Running the Project
Clone the repository.
The project uses a pre-populated rewards.db file, so no database setup is needed.

Run the application from the root directory:
go run .
The server will start on http://localhost:8080.

API Endpoints

1. Create a User
Creates a new user.
Endpoint: POST http://localhost:8080/api/v1/users
Request Body:
JSON:
{
    "name": "John Doe"
}

2. Create a Bill
Creates a new, unpaid bill for a user. To test the reward logic, create at least three bills with future due dates.
Endpoint: POST http://localhost:8080/api/v1/bills
Request Body:
JSON : 
{
    "user_id": 1,
    "amount": 5000,
    "due_date": "2025-10-15T23:59:59Z"
}

3. Pay a Bill
Marks a bill as paid and triggers the reward check logic.
Endpoint: POST http://localhost:8080/api/v1/bills/{billID}/pay
Example URL: http://localhost:8080/api/v1/bills/1/pay
Request Body: (No body required)
Success Response (with reward):
JSON : 
{
    "status": "success",
    "bill": {
        "id": 3,
        "user_id": 1,
        "amount": 2200,
        "due_date": "2025-10-31T23:59:59Z",
        "payment_date": "2025-09-14T19:20:00Z",
        "status": "PAID_ON_TIME"
    },
    "reward_message": "Congratulations! You've earned a $10 Amazon Gift Card."
}
