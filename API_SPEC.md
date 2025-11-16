# Authentication API Specification

Base URL:
bash
Copy code
http://localhost:8000/api

### 1. Register User
POST /register
Request (JSON)
json
Copy code
{
  "name": "Peter Kabwe",
  "email": "peter@example.com",
  "password": "mypassword123"
}
Successful Response (JSON)
json
Copy code
{
  "ID": 1,
  "CreatedAt": "2025-01-01T12:00:00Z",
  "UpdatedAt": "2025-01-01T12:00:00Z",
  "DeletedAt": null,
  "name": "Peter Kabwe",
  "email": "peter@example.com",
  "password": "$2a$12$encryptedhash..."
}


### 2. Login User
POST /login
Request (JSON)
json
Copy code
{
  "email": "peter@example.com",
  "password": "mypassword123"
}
Successful Response (JSON)
json
Copy code
{
  "message": "success"
}
Notes
A JWT token is set in HTTP-only cookie named jwt

You do NOT need to store the token in frontend localStorage

### 3. Get Authenticated User
GET /user
Request
No body required.
The cookie jwt must be present.

Successful Response (JSON)
json
Copy code
{
  "ID": 1,
  "CreatedAt": "2025-01-01T12:00:00Z",
  "UpdatedAt": "2025-01-01T12:00:00Z",
  "DeletedAt": null,
  "name": "Peter Kabwe",
  "email": "peter@example.com"
}

Failure Response (JWT missing or invalid)
json
Copy code
{
  "message": "unauthenticated"
}


### 4. Logout User
POST /logout
Request
No body required.

Successful Response
json
Copy code
{
  "message": "success"
}

Notes
The backend deletes the JWT cookie by setting an expired cookie.

Summary of Endpoints
Method	Endpoint	Description	Auth          Required
POST	/register	Create a new user	      ❌ No
POST	/login	    Log in user, set cookie	  ❌ No
GET	/user	        Get current user	      ✔️ Yes
POST	/logout	    Logout user	              ✔️ Yes

Frontend Notes (React)
When sending requests from frontend:

Must include cookies:
js
Copy code
fetch("http://localhost:8000/api/user", {
  method: "GET",
  credentials: "include"
});
Otherwise the cookie will not be sent.