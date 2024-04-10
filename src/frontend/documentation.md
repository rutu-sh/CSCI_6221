## Authentication

- Authentication is required for all endpoints.
- Authentication is done via JWT token.
- Include the JWT token in the Authorization header of your request.

## Endpoints

### 1. User Authentication

#### POST /auth/signup

- Description: Creates a new user account.
- Request Body:
```json
{
    "email": "user@example.com",
    "password": "password123"
}
```
- Response:
```json
{
    "success": true,
    "message": "User account created successfully"
}
```

#### POST /auth/login

- Description: Logs in an existing user.
- Request Body:
```json
{
    "email": "user@example.com",
    "password": "password123"
}
```
- Response:
```json
{
    "success": true,
    "token": "JWT_TOKEN"
}
```

### 2. Subscription Management

#### POST /subscriptions

- Description: Creates a new subscription.
- Request Body:
```json
{
    "vendor_name": "Vendor ABC",
    "vendor_url": "https://vendorabc.com",
    "duration": "30"
}
```
- Response:
```json
{
    "success": true,
    "message": "Subscription created successfully"
}
```

#### GET /subscriptions

- Description: Retrieves all subscriptions for the authenticated user.
- Response:
```json
{
    "subscriptions": [
        {
            "id": 1,
            "vendor_name": "Vendor ABC",
            "vendor_url": "https://vendorabc.com",
            "duration": "30"
        },
        {
            "id": 2,
            "vendor_name": "Vendor XYZ",
            "vendor_url": "https://vendorxyz.com",
            "duration": "45"
        }
    ]
}
```

#### GET /subscriptions/{subscription_id}

- Description: Retrieves details of a specific subscription.
- Response:
```json
{
    "id": 1,
    "vendor_name": "Vendor ABC",
    "vendor_url": "https://vendorabc.com",
    "duration": "30"
}
```

#### DELETE /subscriptions/{subscription_id}

- Description: Deletes a specific subscription.
- Response:
```json
{
    "success": true,
    "message": "Subscription deleted successfully"
}
```

## Error Responses

- HTTP status codes in the 4xx or 5xx range indicate an error.
- Error responses will contain a JSON object with an error message explaining the issue.

### 400 Bad Request

- **Description:** The request was invalid or missing required parameters.
- **Response Body:**
```json
{
    "error": "Bad Request",
    "message": "Invalid request parameters"
}
```

### 401 Unauthorized

- **Description:** Authentication credentials were missing or invalid.
- **Response Body:**
```json
{
    "error": "Unauthorized",
    "message": "Authentication failed"
}
```

### 403 Forbidden

- **Description:** The user is not authorized to access the requested resource.
- **Response Body:**
```json
{
    "error": "Forbidden",
    "message": "Access denied"
}
```

### 404 Not Found

- **Description:** The requested resource does not exist.
- **Response Body:**
```json
{
    "error": "Not Found",
    "message": "Resource not found"
}
```

### 500 Internal Server Error

- **Description:** An unexpected error occurred on the server.
- **Response Body:**
```json
{
    "error": "Internal Server Error",
    "message": "An unexpected error occurred"
}
```
