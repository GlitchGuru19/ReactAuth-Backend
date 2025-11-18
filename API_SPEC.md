This document explains:

What each endpoint does

What body to send

What response to expect

Exactly how to test each one in Postman

Base URL
http://localhost:8000/api

------------------------------------
🔵 1. REGISTER USER
POST /api/register

Creates a new user.

✅ Request Body (Postman → Body → raw → JSON)
{
  "name": "John Doe",
  "email": "johndoe@example.com",
  "password": "SheLovesIce456"
}

🟢 Successful Response
{
  "id": 1,
  "name": "Nathan",
  "email": "nathan@example.com"
}

🧪 How to Test in Postman

Create a new request

Method: POST

URL: http://localhost:8000/api/register

Body → raw → JSON

Paste the JSON

Send

You should get user info back (no password returned).

------------------------------------
🟠 2. LOGIN
POST /api/login

Used to authenticate and receive access + refresh tokens.

✅ Request Body
{
  "email": "nathan@example.com",
  "password": "123456"
}

🟢 Successful Response
{
  "accessToken": "xxxxx",
  "refreshToken": "xxxxx",
  "user": {
    "id": 1,
    "name": "Nathan",
    "email": "nathan@example.com"
  }
}

🧪 How to Test in Postman

Method → POST

URL: http://localhost:8000/api/login

Body → JSON → paste the credentials

Send

Copy the accessToken for protected routes

Copy the refreshToken for refreshing token

------------------------------------
🟣 3. REFRESH TOKEN
POST /api/refresh

Takes a refreshToken and returns a new accessToken.

✅ Request Body
{
  "refreshToken": "YOUR_REFRESH_TOKEN_HERE"
}

🟢 Successful Response
{
  "accessToken": "new_access_token_here"
}

🧪 How to Test in Postman

Method: POST

URL: http://localhost:8000/api/refresh

Body → JSON

Paste the refresh token

Send

You will receive a fresh access token.

------------------------------------
🔴 4. LOGOUT
POST /api/logout

JWTs are stateless. Logout is done on the frontend by deleting the tokens.

🟢 Response
{
  "message": "Logout successful. Delete tokens on client side."
}

🧪 How to Test in Postman

Just send an empty POST request:

Method: POST

URL: http://localhost:8000/api/logout

Send

You will get the message above.

------------------------------------
🟩 5. GET AUTHENTICATED USER
GET /api/user

This route requires an Authorization header.

🔐 Required Header
Authorization: Bearer ACCESS_TOKEN_HERE

🟢 Successful Response Example
{
  "id": 1,
  "name": "Nathan",
  "email": "nathan@example.com"
}

🧪 How to Test in Postman

Login first and copy your accessToken

Create a new GET request

http://localhost:8000/api/user


Go to Headers

Add:

Key: Authorization

Value: Bearer YOUR_ACCESS_TOKEN

Send

If your access token is valid, you will receive the user.

------------------------------------
🧪 FULL WORKFLOW IN POSTMAN
✔ Step 1 — Register

→ POST /api/register

✔ Step 2 — Login

→ POST /api/login
Save both tokens.

✔ Step 3 — Access Protected Route

→ GET /api/user
Add header:
Authorization: Bearer <accessToken>

✔ Step 4 — Refresh Token

→ POST /api/refresh
Body: { "refreshToken": "..." }

✔ Step 5 — Logout

→ POST /api/logout