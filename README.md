# Receipt Processor

This application provides an API to process receipts and calculate points based on specific rules.

The API has two endpoints: one for submitting a receipt and another for retrieving points for a processed receipt.

## Installation

1. Clone this repository:

   ```
   git clone https://github.com/mattlui/receipt-processor.git

   cd receipt-processor
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

## Running the Application

To start the server, run: `go run main.go `

This will start the server at `http://localhost:8080`.

## API Endpoints

### 1. Process a Receipt

- Endpoint: `/receipts/process`
- Method: `POST`
- Description: Submits a receipt for processing and returns a unique receipt ID.
- Request Body:
  ```
  {
    "retailer": "string",
    "purchaseDate": "YYYY-MM-DD",
    "purchaseTime": "HH:MM",
    "items": [
      {
        "shortDescription": "string",
        "price": "string"
      }
    ],
    "total": "string"
  }
  ```
- Response:

  ```
  {
    "id": "unique-receipt-id"
  }
  ```

### 2. Get Points for a Receipt

- Endpoint: `/receipts/{id}/points`
- Method: `GET`
- Description: Retrieves the points for the specified receipt.
- Path Parameter:
  - `{id}`: The unique ID of the receipt.
- Response:
  ```
  {
    "points": integer
  }
  ```

## Points Calculation Rules

The following rules are applied to determine the points awarded to each receipt:

1. One point for every alphanumeric character in the retailer name.
2. 50 points if the total is a round dollar amount with no cents.
3. 25 points if the total is a multiple of `0.25`.
4. 5 points for every two items on the receipt.
5. If the trimmed length of the item description is a multiple of 3, multiply the price by `0.2` and round up to the nearest integer. The result is the number of points earned.
6. 6 points if the day in the purchase date is odd.
7. 10 points if the time of purchase is after 2:00pm and before 4:00pm.
