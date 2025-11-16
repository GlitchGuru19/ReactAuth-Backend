# ReactAuthBackend API Specification

All endpoints are prefixed with `/api`.

---

## **1. Register User**

**Endpoint:** `/api/register`  
**Method:** `POST`  
**Description:** Registers a new user and stores hashed password in the database.

### Request Body (JSON)
```json
{
  "name": "John Doe",
  "email": "johndoe@example.com",
  "password": "password123"
}
Success Response (201)
json
{
  "id": 1,
  "name": "John Doe",
  "email": "johndoe@example.com"
}

Error Response (400)
json
{
  "message": "failed to parse request"
}
```

## **2. Login User**
**Endpoint:** `/api/login`
Method: POST
Description: Logs in a user, generates JWT, and sets a secure HTTP-only cookie.

### Request Body (JSON)
```json
{
  "email": "johndoe@example.com",
  "password": "password123"
}

Success Response (200)
```json
{
  "message": "success"
}

Cookie set: jwt=<token> (HTTP-only, expires in 24h)

Error Responses
User not found (404)

```json
{
  "message": "user not found"
}

Incorrect password (400)

```json
{
  "message": "incorrect password"
}

Server error (500)

```json
{
  "message": "could not log in"
}
```

## **3. Get Current User (User Info)**
**Endpoint**: `/api/user`
**Method**: `GET`
Description: Returns the currently authenticated user based on JWT in cookie.

Headers
Cookie: jwt=<token>
Success Response (200)
json
{
  "id": 1,
  "name": "John Doe",
  "email": "johndoe@example.com"
}
Error Response (401)
json
Copy code
{
  "message": "unauthenticated"
}


## 4. Logout User
Endpoint: /api/logout
Method: POST
Description: Logs out the user by clearing the JWT cookie.

Headers

Cookie: jwt=<token>

Success Response (200)
json
{
  "message": "success"
}
Notes
All requests and responses are in JSON.

Passwords are always hashed and never returned in any response.

JWT is stored in HTTP-only cookie; frontend should send credentials with fetch:

 fetch("/api/user", {
  method: "GET",
  credentials : "include"
})