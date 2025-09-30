# Needpam Wallet API Documentation

Needpam Wallet is a simple digital wallet API that allows users to manage their balance, deposits, and withdrawals.

## Table of Contents
- [Base URL](#base-url)
- [Authentication](#authentication)
- [API Endpoints](#api-endpoints)
  - [Authentication Endpoints](#authentication-endpoints)
    - [Register User](#1-register-user)
    - [Login User](#2-login-user)
  - [Wallet Endpoints](#wallet-endpoints)
    - [Get Wallet Balance](#3-get-wallet-balance)
    - [Deposit Funds](#4-deposit-funds)
    - [Withdraw Funds](#5-withdraw-funds)
    - [Get Transaction History](#6-get-transaction-history)
- [Error Handling](#error-handling)

## Base URL
All API endpoints are relative to the following base URL:
```
http://148.230.83.69:8080/api
```

## Authentication
Endpoints that require authentication are marked as "protected". To access these endpoints, you must include a JSON Web Token (JWT) in the `Authorization` header of your request.

The token is obtained from the [Login User](#2-login-user) endpoint.

**Header Format:**
```
Authorization: Bearer <your_jwt_token>
```

## API Endpoints

### Authentication Endpoints

#### 1. Register User
Creates a new user account and an associated wallet with a zero balance.

- **Method:** `POST`
- **Endpoint:** `/auth/register`

**Request Body:**
```json
{
    "email": "user@example.com",
    "password": "yoursecurepassword"
}
```

**Successful Response (201 Created):**
```json
{
    "message": "User registered successfully"
}
```

---

#### 2. Login User
Authenticates a user and returns a JWT for accessing protected endpoints.

- **Method:** `POST`
- **Endpoint:** `/auth/login`

**Request Body:**
```json
{
    "email": "user@example.com",
    "password": "yoursecurepassword"
}
```

**Successful Response (200 OK):**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3Mjc0NjYyMzQsInVzZXJfaWQiOjF9.abcdef123456"
}
```
- **`token`**: The JWT to be used for authenticating subsequent requests.

---

### Wallet Endpoints
*(Authentication required for all wallet endpoints)*

#### 3. Get Wallet Balance
Retrieves the current balance and details for the authenticated user's wallet.

- **Method:** `GET`
- **Endpoint:** `/wallet/balance`

**Request Headers:**
```
Authorization: Bearer <your_jwt_token>
```

**Successful Response (200 OK):**
```json
{
    "id": 1,
    "user_id": 1,
    "balance": "100.00",
    "currency": "USD",
    "created_at": "2025-10-01T12:00:00Z",
    "updated_at": "2025-10-01T12:30:00Z"
}
```
- **`id`**: The unique identifier for the wallet.
- **`user_id`**: The ID of the user who owns the wallet.
- **`balance`**: The current amount of money in the wallet.
- **`currency`**: The currency of the wallet (e.g., 'USD').
- **`created_at`**: The timestamp when the wallet was created.
- **`updated_at`**: The timestamp of the last update to the wallet.

---

#### 4. Deposit Funds
Adds a specified amount to the user's wallet balance.

- **Method:** `POST`
- **Endpoint:** `/wallet/deposit`

**Request Headers:**
```
Authorization: Bearer <your_jwt_token>
```

**Request Body:**
```json
{
    "amount": 50.75
}
```

**Successful Response (200 OK):**
```json
{
    "message": "Deposit successful"
}
```

---

#### 5. Withdraw Funds
Subtracts a specified amount from the user's wallet balance. The user must have sufficient funds.

- **Method:** `POST`
- **Endpoint:** `/wallet/withdraw`

**Request Headers:**
```
Authorization: Bearer <your_jwt_token>
```

**Request Body:**
```json
{
    "amount": 25.00
}
```

**Successful Response (200 OK):**
```json
{
    "message": "Withdrawal successful"
}
```

---

#### 6. Get Transaction History
Retrieves a list of all transactions (deposits and withdrawals) for the user's wallet.

- **Method:** `GET`
- **Endpoint:** `/wallet/transactions`

**Request Headers:**
```
Authorization: Bearer <your_jwt_token>
```

**Successful Response (200 OK):**
```json
[
    {
        "id": 2,
        "type": "withdrawal",
        "amount": "25.00",
        "balance_before": "150.75",
        "balance_after": "125.75",
        "description": "ATM withdrawal",
        "created_at": "2025-10-01T12:45:00Z"
    },
    {
        "id": 1,
        "type": "deposit",
        "amount": "50.75",
        "balance_before": "100.00",
        "balance_after": "150.75",
        "description": "Bank transfer",
        "created_at": "2025-10-01T12:30:00Z"
    }
]
```
- **`id`**: The unique identifier for the transaction.
- **`type`**: The type of transaction ('deposit' or 'withdrawal').
- **`amount`**: The amount of the transaction.
- **`balance_before`**: The wallet balance before the transaction occurred.
- **`balance_after`**: The wallet balance after the transaction occurred.
- **`description`**: An optional description for the transaction.
- **`created_at`**: The timestamp when the transaction was recorded.

## Error Handling
The API uses standard HTTP status codes to indicate the success or failure of a request.
- **`4xx` status codes** (e.g., `400 Bad Request`, `401 Unauthorized`, `404 Not Found`) indicate a client-side error.
- **`5xx` status codes** (e.g., `500 Internal Server Error`) indicate a server-side error.

Error responses will typically include a JSON body with an `error` key describing the issue.
```json
{
    "error": "A descriptive error message"
}
```