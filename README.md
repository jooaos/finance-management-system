# Finance Management System

This project is a backend system for managing user accounts on exchanges

The main objective is to build a set of RESTful APIs to manage exchanges, users, and their respective accounts by applying specific business rules, such as age verification and transaction limits.

The system was developed following Golang  and interacting with a MariaDB database.


## Features
* **Exchange Management:** Allows the creation of new exchanges with a name, minimum required age for users, and a maximum transfer amount per transaction.
* **User Management:** Allows the creation of users with unique details like a username and document number.
* **Account Management:** Enables a user to open an account with a specific exchange, enforcing the exchange's minimum age requirement.
* **Financial Operations:**
    * **Deposit:** Add funds to an account , validating the amount against the exchange's maximum transfer limit.
    * **Withdrawal:** Withdraw funds from an account, validating the amount against the exchange's limit and the user's current balance.
* **Queries and Reports:**
    * Fetches a user's consolidated balance across all exchange accounts.
    * Generates a report of daily transaction volume per exchange within a given date range.

## Tech Stack
* **Language:** Golang (v1.24)
* **Database:** MariaDB (v10.8)
* **Infrastructure:**
	* Docker
	* Docker-Compose

## Prerequisites

All commands here are using make, but if you don't want use it, you can access Makefile, copy command and paste in your terminal

## Getting Started

Follow these steps to get the project running locally.

### 1. Copy env
```
cp .env.example .env
```

### 2. Up services
```
make up-dettached
```
Run the follow command and check if your app is up
```
docker container ps | grep "finance_management"
```

It's necessary to appear two containes: `finance_management_api_stable` and `finance_management_db`, if it's showing, you can go to next step

### 3. Run migrations
```
make migration-up
```
Now you can use the application using the endpoint `http://localhost:8080` 

### Developer Mode
If you want to test in developer mode, it's necessary to have [air](https://github.com/air-verse/air) in your local environment

#### 1. Up database
```
make up-service SERVICE=db
```

#### 2. Run migrations (if it's necessary)
```
make migration-up
```

#### 3. Setup depencies
```
make set-up
```
#### 4. Up API using air
```
make air-up SERVICE=api
```

Now you can use the application using the endpoint `http://localhost:8080` 

### Testing
To run unit test, use the following command
```
make test-unit
```
## API Documentation
Below are the available API endpoints as required 

#### Testing with Postman 
To make it easier to test the API endpoints, this project includes a Postman collection that you can import. The collection file, named `Finance Management.postman_collection`, is located in the `docs/postman/` directory at the root of the project. It contains all the available requests documented below. An environment file, `Finance Management.postman_environment`, is also included. It's recommended that you import it as well, as it contains the `api_url` variable pre-configured to `http://localhost:8080/api`.

-----

### Exchange

#### 1. Create Exchange
This endpoint is used to create exchanges in the system. Each exchange should have a name, a minimum age for users, and a maximum transfer amount. 
* **Endpoint:** `POST /api/exchange`
* **Body (Request):**
```json
{
	"name": "CryptoExchange_V3",
	"min_age": 5,
	"max_amount": 5000
}
```
* **Response (201 Created):**
```json
{
	"data": {
		"id": 6,
		"name": "CryptoExchange_V3"
	}
}
```
-----

### User

#### 2. Create User
This endpoint is used to create users in the system. Each user should have a unique username, a date of birth, and a unique document number.

* **Endpoint:** `POST /api/user`
* **Body (Request):**
```json
{
	"username": "user_test_02",
	"document_number": "12345678900",
	"date_of_birth": "1990-01-01"
}
```
* **Response (201 Created):**
```json
{
	"data": {
		"id": 2,
		"username": "user_test_03",
		"document_number": "12345678901",
		"date_of_birth": "1990-01-01T00:00:00Z"
	}
}
```
-----
### Account
#### 3. Create Account
This endpoint creates a user account for a specific exchange. An account can only be created if the user's age is greater than or equal to the exchange's minimum age requirement. The API must generate and return an account ID in the response.

* **Endpoint:** `POST /api/account`

* **Body (Request):**

```json
{
	"user_id": 1,
	"exchange_id": 5
}
```
* **Response (201 Created):**

```json
{
    "data": {
        "account_id": "0d96dfb4-7587-4a82-b92f-f349dfa661da"
    }
}
```
### Transaction
#### 4. Make a Deposit
This endpoint is used to add funds to a user's account.
* The deposit amount must be greater than 0. 
* The amount must not exceed the exchange's maximum transfer amount. 
* **Endpoint:** `POST /api/transaction/deposit`

* **Body (Request):**
```json
{
	"account_id": "67f13780-b962-4430-bd5d-52bef38fb231",
	"amount": 5000
}
```
* **Response (201 CREATED):**
```json

{
	"data": {
		"account_id": "67f13780-b962-4430-bd5d-52bef38fb231",
		"amount": 5000,
		"new_balance": 45000
	}
}
```

#### 5. Make a Withdrawal
This endpoint is used to withdraw a specific amount from a user's balance.
* The withdrawal amount must be greater than 0. 
* The amount must not exceed the exchange's maximum transfer limit.
* The amount must not exceed the account's current balance.
* **Endpoint:** `POST /api/transaction/withdrawal`
* **Body (Request):**
```json
{
	"account_id": "2c64242a-ce18-4193-891e-234937f80cfd",
	"amount": 500.25
}
```
* **Response (201 CREATED):**
```json
{
	"data": {
		"account_id": "2c64242a-ce18-4193-891e-234937f80cfd",
		"amount": 500.25,
		"new_balance": 44499.75
	}
}
```

#### 6. Get User Balance
This endpoint retrieves the balances an individual user holds across all of their exchange accounts.
* The response must include an individual balance for each exchange where the user has funds, excluding accounts with a zero balance. 
* The response must also include the total balance from all accounts.
* **Endpoint:** `GET /api/user/{userId}/balance`
* **Response (200 OK):**
```json
{
    "data": {
        "total_balance": 70000,
        "balances_by_exchange": [
            {
                "user_id": 1,
                "exchange_id": 1,
                "account_id": "2c64242a-ce18-4193-891e-234937f80cfd",
                "balance": 35000,
                "exchange": {
                    "name": "CryptoExchange"
                }
            },
            {
                "user_id": 1,
                "exchange_id": 5,
                "account_id": "67f13780-b962-4430-bd5d-52bef38fb231",
                "balance": 35000,
                "exchange": {
                    "name": "CryptoExchange_V2"
                }
            }
        ]
    }
}
```

#### 7. Daily Transaction Volume Report
This endpoint reports transaction volumes per exchange for each day within a date range.
* The API must accept a start date and an end date as query parameters.
* It returns the total transaction amount per exchange for each day in the specified range. 
* **Endpoint:** `GET /api/transactions/daily-volume?start_date=2025-10-01&end_date=2025-10-05`
* **Response (200 OK):**
```json
{
    "data": [
        {
            "date": "2025-10-09",
            "exchanges": [
                {
                    "exchange_name": "CryptoExchange",
                    "amount": 30000
                },
                {
                    "exchange_name": "CryptoExchange_V2",
                    "amount": 35000
                }
            ]
        },
        {
            "date": "2025-10-08",
            "exchanges": [
                {
                    "exchange_name": "CryptoExchange",
                    "amount": 5000
                }
            ]
        }
    ]
}
```
