# Go Backend: Orders API

> [!TIP]
> **Build this project yourself, step by step, on [Codealong](https://codealong.dev/courses/backend-go).**
>
> This repository is the finished code. In the *Backend API Development* course you build it lesson by lesson: start with your own server on `net/http`, move to a layered architecture that scales, then add orders with stock reservation. Pause the video, change the code, run it, and get instant feedback on each challenge.
>
> **[Start building your own API →](https://codealong.dev/courses/backend-go)**

A small REST API for users, products and orders, written in Go as part of the [Backend API Development](https://codealong.dev/courses/backend-go) course on Codealong.

It uses only the standard library (`net/http` with method-and-path routing) plus [`google/uuid`](https://github.com/google/uuid). Data is kept in memory, so you can run it with no database or other setup.

## What it covers

- **A layered structure for each domain:** handler (HTTP) → service (business rules) → repository (storage).
- **Repository interfaces** with in-memory implementations, so a real database can be added later without touching the services.
- **Placing an order across domains:** the order service checks that the customer and products exist, then reserves stock before it saves the order.
- **All-or-nothing stock reservation:** every line item is checked before any stock is taken, under a single lock.
- **The same JSON shape for every response**, for both success and error.

## Requirements

- Go 1.26 or newer

## Getting started

```sh
git clone https://github.com/codealong-dev/go-backend-ecom.git
cd go-backend-ecom
go run ./api
```

The server listens on `http://localhost:8080`. Everything is stored in memory, so data resets each time you restart the server. Some demo data is loaded at startup (see [Seed data](#seed-data)).

## Project structure

```
.
├── api/
│   └── main.go              # wiring: repositories → services → handlers, seed, start server
└── internal/
    ├── platform/httpx/      # JSON response helpers shared by every handler
    ├── user/
    ├── product/
    └── order/
```

Each domain package follows the same layout:

| File         | Responsibility                                          |
| ------------ | ------------------------------------------------------- |
| `<domain>.go` | Entity types, request types and errors                 |
| `handler.go` | Route registration, request decoding, maps errors to HTTP status codes |
| `service.go` | Business rules and the `Repository` interface          |
| `memory.go`  | In-memory `Repository`, safe for concurrent use with `sync.RWMutex` |
| `seed.go`    | Demo data loaded at startup (users and products only)  |

## API

Every response is JSON with one of two shapes:

```json
{ "data": ... }
```

```json
{ "error": "message" }
```

Prices are whole numbers in cents (`2499` = 24.99).

| Method | Path            | Description                     |
| ------ | --------------- | ------------------------------- |
| GET    | `/users`        | List users                      |
| POST   | `/users`        | Create a user                   |
| GET    | `/users/{id}`   | Get a user by ID                |
| GET    | `/products`     | List products and current stock |
| POST   | `/orders`       | Place an order                  |
| GET    | `/orders/{id}`  | Get an order by ID              |

### Examples

List products:

```sh
curl localhost:8080/products
```

Create a user:

```sh
curl -X POST localhost:8080/users \
  -d '{"name": "Grace"}'
```

```json
{"data":{"id":"ce81b23d-fa43-4fd1-b879-64cf7dea1001","name":"Grace"}}
```

Place an order for the seeded customer:

```sh
curl -X POST localhost:8080/orders \
  -d '{
    "customer_id": "9a8b7c6d-5e4f-4a3b-8c2d-000000000001",
    "items": [
      { "product_id": "7c1e4d2a-3b5f-4e8a-9c6d-000000000001", "quantity": 2 }
    ]
  }'
```

```json
{
  "data": {
    "id": "92ffbde5-9d71-4cdc-8e5c-d032661a2536",
    "status": "pending",
    "items": [
      {
        "id": "ce1181e3-0c92-4391-9682-384374fbe494",
        "product_id": "7c1e4d2a-3b5f-4e8a-9c6d-000000000001",
        "order_id": "92ffbde5-9d71-4cdc-8e5c-d032661a2536",
        "price": 2499,
        "quantity": 2,
        "created_at": "2026-10-24T18:13:14.66855Z"
      }
    ],
    "customer_id": "9a8b7c6d-5e4f-4a3b-8c2d-000000000001",
    "created_at": "2026-10-24T18:13:14.668552Z"
  }
}
```

Fetch it back with the returned ID:

```sh
curl localhost:8080/orders/92ffbde5-9d71-4cdc-8e5c-d032661a2536
```

Each order item stores the product's price at the moment the order was placed, so changing a product's price later does not change past orders.

### Errors

| Status | When                                                                 |
| ------ | -------------------------------------------------------------------- |
| 400    | Invalid JSON, missing `name` / `customer_id` / `items`, or an order ID that is not a UUID |
| 404    | User or order not found                                              |
| 422    | The order refers to a customer or product that does not exist       |

## Seed data

**User**

| ID                                     | Name |
| -------------------------------------- | ---- |
| `9a8b7c6d-5e4f-4a3b-8c2d-000000000001` | Ada  |

**Products**

| ID                                     | Name                | Price | Stock |
| -------------------------------------- | ------------------- | ----- | ----- |
| `7c1e4d2a-3b5f-4e8a-9c6d-000000000001` | Wireless Mouse      | 2499  | 25    |
| `7c1e4d2a-3b5f-4e8a-9c6d-000000000002` | Mechanical Keyboard | 8999  | 3     |
| `7c1e4d2a-3b5f-4e8a-9c6d-000000000003` | 27-inch Monitor     | 22900 | 5     |
| `7c1e4d2a-3b5f-4e8a-9c6d-000000000004` | USB-C Hub           | 3999  | 0     |

The USB-C Hub starts with no stock, so you can try ordering something that is out of stock.

## License

[MIT](LICENSE)
