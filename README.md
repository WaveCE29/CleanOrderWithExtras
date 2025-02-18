# Contact Information

- **Name**: Sarin Hongthong
- **Tel**: 093-9195903
- **Email**: <sarin8845@gmail.com>

## 🚀 Getting Started

## Prerequisites

Ensure you have the following installed:

- [Go](https://go.dev/dl/) (version 1.18+)
- `make` (optional, for running commands)

## Install Dependencies

```sh
go mod tidy
```

### 🏃 Running the Service

Start the server:

```sh
go run cmd/server/main.go
```

The API will be available at <http://localhost:8080>

## API Endpoints

### Clean Order

Endpoint: `POST /clean-order`

Description: Cleans and normalizes an order input.

- Request Body (JSON):

  ```json
    [
        {
        "no": 1,
        "platformProductId": "--FG0A-CLEAR-OPPOA3*2/FG0A-MATTE-OPPOA3",
        "qty": 1,
        "totalPrice": 120
        } 
    ]
    ```

- Response (JSON):

  ```json
    [
        {
            "no": 1,
            "productId": "FG0A-CLEAR-OPPOA3",
            "materialId": "FG0A-CLEAR",
            "modelId": "OPPOA3",
            "qty": 2,
            "unitPrice": 40,
            "totalPrice": 80
        },
        {
            "no": 2,
            "productId": "FG0A-MATTE-OPPOA3",
            "materialId": "FG0A-MATTE",
            "modelId": "OPPOA3",
            "qty": 1,
            "unitPrice": 40,
            "totalPrice": 40
        }
    ]
    ```

### Running Tests

To run unit tests:

```code
go test ./internal/services -v
```
