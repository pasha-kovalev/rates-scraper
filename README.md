# Exchange Rate Service

### Usage

To start the service, run:
```sh
docker-compose up
```

This will start the server on http://localhost:8080.

####  Get All Rates  [GET /rates]
Retrieves all available exchange rates.

#### Get Rates by Date [GET /rates/:date]
Retrieves exchange rates for a specific date.
- Params:
  - date (string): The date in YYYY-MM-DD format.

