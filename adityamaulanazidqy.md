# go-coursework

A RESTful API built with Go, designed for university environments. It enables lecturers to assign coursework and allows students to submit their assignments efficiently.

## Features

- Role-based authentication with JWT tokens
- Lecturers can create and manage coursework
- Students can view and submit assignments
- Modular and scalable code structure
- Environment-based configuration
- Integrated logging and database migration

## Tech Stack

| Technology        | Description                                 |
|-------------------|---------------------------------------------|
| Go                | Core programming language                   |
| Fiber             | Web framework for building fast HTTP APIs   |
| PostgreSQL        | Relational database for data storage        |
| Redis             | Caching and session management              |
| JWT               | Role-based authentication                   |
| logrus (sirupsen) | Structured logging                          |
| database/migrate  | Database migration management               |
| RESTful Box       | API testing tool                            |

## Getting Started

1. Clone the repository.
2. Copy `.env.example` to `.env` and set the appropriate values.
3. Run database migrations.
4. Start the server with:

```bash
go run main.go
```