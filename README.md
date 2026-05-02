# Task Manager API

A RESTful API for managing tasks built with Go and PostgreSQL.

## Features
- User registration and login with JWT authentication
- Task management with full CRUD operations
- Password hashing with bcrypt for security
- PostgreSQL database with Docker containerization

## Project Structure
```task-manager/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── auth/
│   │   └── jwt.go            # JWT generation and validation
│   ├── database/
│   │   ├── db.go             # Database connection
│   │   ├── tasks.go          # Task CRUD operations
│   │   └── users.go          # User CRUD operations
│   ├── handlers/
│   │   ├── tasks.go          # Task HTTP handlers
│   │   └── auth.go           # Auth HTTP handlers
│   └── models/
│       ├── task.go           # Task data models
│       └── user.go           # User data models
├── sql/
│   └── init.sql              # Database schema
├── Dockerfile                # Container build instructions
├── docker-compose.yml        # Multi-container setup
├── .dockerignore             # Files excluded from Docker image
├── go.mod
├── go.sum
└── README.md
```