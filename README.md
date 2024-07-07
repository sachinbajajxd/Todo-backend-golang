# TODO Application

This project is a simple TODO application built with Go and MongoDB. It includes user registration, login, JWT token generation, and CRUD APIs for managing TODO items.

## Setup

### Prerequisites

- Go
- MongoDB

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/sachinbajajxd/Todo-backend-golang.git

2. Install Dependencies
    ```sh
    go mod tidy

3. Run the application
   ```sh
   go run main.go

## Project Structure
* main.go: Entry point of the application
* controllers/: Contains the handler functions for the API endpoints and middleware
* models/: Contains the data models for the application

### API Endpoints
1 User Registration
* Endpoint: POST /register

```sh
  {
    "username": "your-username",
    "password": "your-password"
  }
```

* Response
![image](https://github.com/sachinbajajxd/Todo-backend-golang/assets/112001510/4ad71c3f-4687-4834-bea5-ee9d13f11e07)

2 User Login
* Endpoint: POST /login

```sh
  {
    "username": "your-username",
    "password": "your-password"
  }
```

* Response
![image](https://github.com/sachinbajajxd/Todo-backend-golang/assets/112001510/d82b1026-08fe-4507-b260-2d3c9cd959d1)


3. Get all posts (paginated)
* Endpoint: GET /posts?userId=_&sortBy=_&sortOrder=_&page=_&limit=_

```sh
  Headers -> Authorization: Bearer your-jwt-token
  Query params ->
  {
    "userId": "user id",
    "sortBy": "createdAt",
    "sortOrder": "asc/desc",
    "page" : number,
    "limit": number
  }
 ```

* Response
![image](https://github.com/sachinbajajxd/Todo-backend-golang/assets/112001510/0d8ac526-d194-451f-b9f6-d96c4bddff54)


4. Create a post
* Endpoint POST/posts

```sh
Headers-> Authorization: Bearer your-jwt-token
Request Body ->
{
  "title": "Your TODO title",
  "description": "Your TODO description",
  "status": "Your TODO status"
  "user_id": "User id"
}


```

* Response
![image](https://github.com/sachinbajajxd/Todo-backend-golang/assets/112001510/8cd78e5f-ea1d-46d0-9464-8257abed5144)


5. Update a post
* Endpoint PUT/posts/{id}

```sh
Headers-> Authorization: Bearer your-jwt-token
Request Body ->
{
  "title": "Your TODO title",
  "description": "Your TODO description",
  "status": "Your TODO status"
  "user_id": "User id"
}

```

* Response
![image](https://github.com/sachinbajajxd/Todo-backend-golang/assets/112001510/67a90051-be0c-427d-8a5a-13df4cc48d1b)


6. Delete a post
* Endpoint DELETE/posts/{id}

```sh
Headers-> Authorization: Bearer your-jwt-token
Request Body ->
{
  "user_id": "User id"
}

```

* Response
![image](https://github.com/sachinbajajxd/Todo-backend-golang/assets/112001510/eb3fc026-025e-404f-8457-19a8a2885c73)


## Authentication Flow
* Register: Create a new user account using the /register endpoint.
* Login: Obtain a JWT token using the /login endpoint.


## Conclusion
This project provides a basic structure for a TODO application with user authentication and CRUD operations. 
