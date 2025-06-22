# Technical Specification

## Loyalty Accrual System "Gophermart"

---

### General Requirements

The system is an HTTP API with the following business logic requirements:

* User registration, authentication, and authorization;
* Accepting order numbers from registered users;
* Tracking and managing the list of submitted order numbers by each user;
* Tracking and managing the loyalty account of each registered user;
* Verifying submitted order numbers via the loyalty points calculation system;
* Crediting the appropriate rewards for valid order numbers to the user’s loyalty account.

![image](https://pictures.s3.yandex.net:443/resources/gophermart2x_1634502166.png)

### Abstract Interaction Flow

Below is an abstract business logic flow for user interaction with the system:

1. The user registers in the Gophermart loyalty system.
2. The user makes a purchase in the Gophermart online store.
3. The order is sent to the loyalty point calculation system.
4. The user submits the order number to the loyalty system.
5. The system links the order number to the user and checks it with the loyalty point calculation system.
6. If loyalty points are accrued, they are credited to the user's loyalty account.
7. The user spends available loyalty points to partially or fully pay for future orders in the Gophermart online store.

> **Note:**
> - Step 2 is hypothetical and not required to be implemented in this project.
> - Step 3 is implemented in the loyalty point calculation system and is not required in this project.

### Loyalty Points Calculation System

The loyalty point calculation system is an external service in a trusted environment. It operates as a black box and is not inspectable by external clients. It calculates the points awarded for a completed order using complex algorithms, which may change at any time.

The only data available to consumers is the number of points accrued for a specific order. The reasons for awarding or not awarding points are not available to consumers.

The communication protocol for this service is provided at the end.

### Summary of HTTP API

The Gophermart loyalty accrual system must expose the following HTTP handlers:

* `POST /api/user/register` — user registration
* `POST /api/user/login` — user authentication
* `POST /api/user/orders` — user submits an order number for processing
* `GET /api/user/orders` — retrieve submitted order numbers, their processing statuses, and accrual information
* `GET /api/user/balance` — get current loyalty point balance
* `POST /api/user/balance/withdraw` — request to withdraw points to pay for a new order
* `GET /api/user/withdrawals` — retrieve withdrawal history

### General Limitations and Requirements

* Data storage — PostgreSQL
* Table structure is at the student's discretion
* Types and format of stored data (including passwords and other sensitive data) are at the student's discretion
* The client may support compressed HTTP requests/responses
* The client is not required to follow this exact API specification; request validation is up to the student
* The authentication/authorization mechanism is at the student's discretion
* Order numbers are unique and never repeat
* An order number can be submitted only once by a user
* An order may result in no reward
* Rewards are credited and spent in virtual points: 1 point = 1 ruble

#### User Registration

**Endpoint:** `POST /api/user/register`

Users register with a login/password pair. Logins must be unique.
Upon successful registration, the user is automatically authenticated.

**Request:**

```json
POST /api/user/register HTTP/1.1
Content-Type: application/json

{
  "login": "<login>",
  "password": "<password>"
}
```

**Responses:**

- `200` — user successfully registered and authenticated
- `400` — invalid request format
- `409` — login already taken
- `500` — internal server error

#### User Authentication

**Endpoint:** `POST /api/user/login`

**Request:**

```json
POST /api/user/login HTTP/1.1
Content-Type: application/json

{
  "login": "<login>",
  "password": "<password>"
}
```

**Responses:**

- `200` — user successfully authenticated
- `400` — invalid request format
- `401` — invalid login/password pair
- `500` — internal server error

#### Submitting an Order Number

**Endpoint:** `POST /api/user/orders`

Authenticated users only. Order number is a string of digits.

You may validate the number using the [Luhn algorithm](https://en.wikipedia.org/wiki/Luhn_algorithm).

**Request:**

```http
POST /api/user/orders HTTP/1.1
Content-Type: text/plain

12345678903
```

**Responses:**

- `200` — order already submitted by this user
- `202` — new order accepted for processing
- `400` — invalid request format
- `401` — user not authenticated
- `409` — order submitted by another user
- `422` — invalid order number format
- `500` — internal server error

#### Retrieving Submitted Orders

**Endpoint:** `GET /api/user/orders`

Authenticated users only. Sort orders by upload time (newest first).

**Statuses:**
- `NEW`, `PROCESSING`, `INVALID`, `PROCESSED`

**Response Example:**

```json
[
  {
    "number": "9278923470",
    "status": "PROCESSED",
    "accrual": 500,
    "uploaded_at": "2020-12-10T15:15:45+03:00"
  },
  {
    "number": "12345678903",
    "status": "PROCESSING",
    "uploaded_at": "2020-12-10T15:12:01+03:00"
  },
  {
    "number": "346436439",
    "status": "INVALID",
    "uploaded_at": "2020-12-09T16:09:53+03:00"
  }
]
```

**Responses:**
- `200` — success
- `204` — no data
- `401` — not authenticated
- `500` — server error

#### Getting User Balance

**Endpoint:** `GET /api/user/balance`

**Response Example:**

```json
{
  "current": 500.5,
  "withdrawn": 42
}
```

**Responses:**
- `200` — success
- `401` — not authenticated
- `500` — server error

#### Withdrawing Points

**Endpoint:** `POST /api/user/balance/withdraw`

**Request:**

```json
{
  "order": "2377225624",
  "sum": 751
}
```

**Responses:**
- `200` — success
- `401` — not authenticated
- `402` — insufficient funds
- `422` — invalid order number
- `500` — server error

#### Retrieving Withdrawals

**Endpoint:** `GET /api/user/withdrawals`

**Response Example:**

```json
[
  {
    "order": "2377225624",
    "sum": 500,
    "processed_at": "2020-12-09T16:09:57+03:00"
  }
]
```

**Responses:**
- `200` — success
- `204` — no withdrawals
- `401` — not authenticated
- `500` — server error

### Accrual System Interaction

**Endpoint:** `GET /api/orders/{number}`

**Response Example:**

```json
{
  "order": "<number>",
  "status": "PROCESSED",
  "accrual": 500
}
```

**Statuses:** `REGISTERED`, `INVALID`, `PROCESSING`, `PROCESSED`

**Other Responses:**
- `204` — order not found
- `429` — rate limit exceeded (use `Retry-After`)
- `500` — server error

### Service Configuration

Supported via environment variables or flags:

* `RUN_ADDRESS` or `-a` — service address and port
* `DATABASE_URI` or `-d` — database connection URI
* `ACCRUAL_SYSTEM_ADDRESS` or `-r` — accrual system address
