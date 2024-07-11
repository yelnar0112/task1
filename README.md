# Task1

This project is a simple HTTP proxy server implemented in Go. It supports forwarding requests to specified URLs and returning the responses. Additionally, it provides Swagger documentation for the API.


## Getting Started

### Prerequisites

- Go 1.16 or later
- `swag` and `http-swagger` packages

### Installing

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yelnar0112/task1.git







API Endpoints

- Method: POST
- Description: Forwards a request to the specified URL and returns the response.

```yaml
Request Body:
{
  "method": "GET",
  "url": "http://example.com",
  "headers": {
    "Content-Type": "application/json"
  }
}

Response:
{
  "id": 1,
  "status": 200,
  "headers": {
    "Content-Type": "application/json"
  },
  "length": 1234
}
