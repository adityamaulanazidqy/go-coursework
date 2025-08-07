# go-coursework

A RESTful API built with Go, designed for university environments. It enables lecturers to assign coursework and allows students to submit their assignments efficiently.

![Go](https://img.shields.io/badge/Go-1.21-blue?logo=go)
![Build](https://img.shields.io/badge/build-passing-brightgreen?style=flat-square)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![Repo Size](https://img.shields.io/github/repo-size/adityamaulanazidqy/go-coursework)
![Last Commit](https://img.shields.io/github/last-commit/adityamaulanazidqy/go-coursework)
![Issues](https://img.shields.io/github/issues/adityamaulanazidqy/go-coursework)
![Stars](https://img.shields.io/github/stars/adityamaulanazidqy/go-coursework?style=social)

## Features

- 👨‍🎓 User roles: Instructor and Student with JWT-based authentication
- 📝 Task and subject CRUD endpoints
- 📤 File upload and download support
- 🚨 Real-time push notifications using Firebase Cloud Messaging (FCM) and FONTE
- 📩 Deadline reminder emails with Gomail
- 📄 Swagger UI for interactive API documentation
- 🧰 Redis-based caching and session handling
- 🧑‍💻 Postgres integration with database migration (`golang-migrate`)

## 🌐 API Documentation

Swagger UI is available when the server is running.  
Visit: http://localhost:8080/swagger/index.html

## Tech Stack

| Component        | Technology             |
|------------------|------------------------|
| Language         | Go 1.20                |
| Framework        | Fiber                  |
| Database         | PostgreSQL             |
| Cache            | Redis 7                |
| Email Sending    | Gomail                 |
| Real-time Notify | FCM & FONTE            |
| API Documentation| Swagger                |
| Logging          | Logrus                 |
| Authentication   | JWT Token              |

## 📦 Installation Guide
✅ Prerequisites
- Go 1.20+
- PostgreSQL
- Redis 7+
- Firebase account (FCM) for push notification
- SMTP access (for email reminders)

## 📥 Clone the Repository
```bash
git clone https://github.com/adityamaulanazidqy/go-coursework.git
cd go-coursework
```

## ⚙️ Configure `.env`

```env
JWT_KEY="your_jwt_secret"
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="username"
DB_PASS="password"
DB_NAME="your_db"
REDIS_ADDR="localhost:6379"
FCM_SERVER_KEY="your_fcm_key"
SMTP_HOST="smtp.gmail.com"
SMTP_PORT="587"
SMTP_USER="your_email@example.com"
SMTP_PASS="your_email_password"
```

## 📦 Install Dependencies
```bash
go mod tidy
```

## 🧾 Run Migrations
```bash
migrate -path db/migrations -database "postgres://username:password@localhost:5432/your_db?sslmode=disable" up
```

## 🧪 Generate Swagger Docs
```bash
swag init
```

## 🚀 Start Server
```bash
go run main.go 
```

## 📁 Project Structure

```bash
go-coursework/
├── main.go
├── config/           # Contains configuration files and logic (e.g., database setup, redis, set logrus).         
├── constants/        # Stores application-wide constant values such as messages, enums, and status codes.     
├── internal/         # Houses core business logic and implementation details, structured into subfolders:
│   └── dto/
│   └── handlers/
│   └── helpers/
│   └── logger/
│   └── mapper/
│   └── models/
│   └── repositories/
│   └── routes/
├── json/              # Stores sample or static JSON files, often used for testing or mocking data.
├── migrations/        # Contains database migration files (SQL or code-based) to manage schema changes over time.        
├── pkg/               # Reusable packages or modules that can be shared across the project or even other projects.    
├── test/unit_test/    # Includes unit test files for testing individual components/functions.       
├── ...
```

## License
This project is licensed under the MIT License.

## Contact
Aditya Maulana Zidqy
📧 Email: adityamaullana234@gmail.com
🐙 GitHub: @adityamaulanazidqy

Project Repo: go-coursework
