# 🏌️‍♂️ Golf Course Management API (Go + MySQL)

This project is a simple RESTful API built with **Go (Golang)** and **MySQL**, designed to manage a list of golf courses. It supports full CRUD operations: create, read, update, and delete golf course records.

---

## ✅ Features

- Fetch all golf courses
- Get a single course by ID
- Create a new golf course
- Update an existing course
- Delete a course by ID
- CORS support for frontend integration
- JSON-based API responses

---

## 🛠 Tech Stack

- **Language**: Go (Golang)
- **Database**: MySQL
- **Driver**: [`github.com/go-sql-driver/mysql`](https://github.com/go-sql-driver/mysql)
- **API Format**: REST / JSON

---

## 📦 API Endpoints

| Method | Endpoint                    | Description              |
|--------|-----------------------------|--------------------------|
| GET    | `/api/golfcourses`          | List all golf courses    |
| GET    | `/api/golfcourses/{id}`     | Get a course by ID       |
| POST   | `/api/golfcourses`          | Create a new course      |
| PUT    | `/api/golfcourses/{id}`     | Update a course by ID    |
| DELETE | `/api/golfcourses/{id}`     | Delete a course by ID    |

---
