# Receipt Processor

This application provides an API to process receipts and calculate points based on specific rules.

The API has two endpoints: one for submitting a receipt and another for retrieving points for a processed receipt.

## Installation

1. Clone this repository:

   ```
   git clone https://github.com/MattLui/receipt-processor.git

   cd receipt-processor
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

## Running the Application

To start the server, run: `go run main.go`

This will start the server at `http://localhost:8080`

## API Endpoints

### 1. Process a Receipt

- Endpoint: `/receipts/process`
- Method: `POST`
- Description: Submits a receipt for processing and returns a unique receipt ID.
- Request Body:
  ```json
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

  ```json
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
  ```json
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

## Examples

```json
{
  "retailer": "Target",
  "purchaseDate": "2022-01-01",
  "purchaseTime": "13:01",
  "items": [
    {
      "shortDescription": "Mountain Dew 12PK",
      "price": "6.49"
    },
    {
      "shortDescription": "Emils Cheese Pizza",
      "price": "12.25"
    },
    {
      "shortDescription": "Knorr Creamy Chicken",
      "price": "1.26"
    },
    {
      "shortDescription": "Doritos Nacho Cheese",
      "price": "3.35"
    },
    {
      "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
      "price": "12.00"
    }
  ],
  "total": "35.35"
}
```

```text
Total Points: 28
Breakdown:
     6 points - retailer name has 6 characters
    10 points - 5 items (2 pairs @ 5 points each)
     3 Points - "Emils Cheese Pizza" is 18 characters (a multiple of 3)
                item price of 12.25 * 0.2 = 2.45, rounded up is 3 points
     3 Points - "Klarbrunn 12-PK 12 FL OZ" is 24 characters (a multiple of 3)
                item price of 12.00 * 0.2 = 2.4, rounded up is 3 points
     6 points - purchase day is odd
  + ---------
  = 28 points
```

---

```json
{
  "retailer": "M&M Corner Market",
  "purchaseDate": "2022-03-20",
  "purchaseTime": "14:33",
  "items": [
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    },
    {
      "shortDescription": "Gatorade",
      "price": "2.25"
    }
  ],
  "total": "9.00"
}
```

```text
Total Points: 109
Breakdown:
    50 points - total is a round dollar amount
    25 points - total is a multiple of 0.25
    14 points - retailer name (M&M Corner Market) has 14 alphanumeric characters
                note: '&' is not alphanumeric
    10 points - 2:33pm is between 2:00pm and 4:00pm
    10 points - 4 items (2 pairs @ 5 points each)
  + ---------
  = 109 points
```
