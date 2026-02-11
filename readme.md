# Homework Test from Mampu.io

## Overview
A simple web api to check and withdraw balance, built with Go and use PostgreSQL to store the data.

## Requirements
- Go 1.24+
- PostgreSQL 18+

## Environment Variables
- `DB_URL`: PostgreSQL connection string

## Usage
- Run `go run .` to start the server
- Use `curl` or any HTTP client to send requests to the server

## API Endpoints
- `GET /users/{id}/balance`: Get the current balance
- `POST /users/{id}/withdraw`: Withdraw a certain amount of balance

## Examples
```
curl -X GET http://localhost:8080/users/1/balance
curl -X POST http://localhost:8080/users/1/withdraw -H "content-type: application/json" -d '{"amount": 10000}'
```
