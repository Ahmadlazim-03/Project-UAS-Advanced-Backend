# Student Achievement System (Backend & Frontend)

This project is a comprehensive system for managing student achievements, featuring a robust backend API and a simple web interface.

## 🚀 Features

*   **Role-Based Access Control (RBAC)**: Admin, Student (Mahasiswa), Advisor (Dosen Wali).
*   **Hybrid Database Architecture**:
    *   **PostgreSQL**: Manages Users, Roles, Permissions, and Achievement References (Relational Data).
    *   **MongoDB**: Stores dynamic and flexible Achievement details (NoSQL).
*   **RESTful API**: Fully documented with Swagger.
*   **Web Interface**: A simple dashboard for students to manage achievements.

## 🛠 Tech Stack

*   **Language**: Go (Golang)
*   **Framework**: Fiber v2
*   **Databases**:
    *   PostgreSQL (via GORM)
    *   MongoDB (via official mongo-driver)
*   **Authentication**: JWT (JSON Web Tokens)
*   **Documentation**: Swagger (Swaggo)
*   **Frontend**: HTML5, Bootstrap 5, Vanilla JavaScript

## 📂 Project Structure

```
.
├── cmd/                # Main application entry point
├── database/           # Database connection logic
├── docs/               # Swagger documentation files
├── internal/           # Internal application code (if any)
├── middleware/         # HTTP Middleware (Auth, Logging)
├── models/             # Data models (Structs)
├── public/             # Static frontend files (HTML, CSS, JS)
├── repository/         # Data Access Layer
├── routes/             # API Route definitions
├── services/           # Business Logic Layer
├── tests/              # Unit tests
├── utils/              # Utility functions
├── .env                # Environment variables
├── go.mod              # Go module definition
└── main.go             # Application entry point
```

## ⚙️ Setup & Installation

1.  **Prerequisites**:
    *   Go 1.20+
    *   PostgreSQL
    *   MongoDB

2.  **Clone the Repository**:
    ```bash
    git clone https://github.com/Ahmadlazim-03/Project-UAS-Advanced-Backend.git
    cd Project-UAS-Advanced-Backend
    ```

3.  **Configure Environment**:
    Create a `.env` file in the root directory:
    ```env
    DB_HOST=localhost
    DB_USER=postgres
    DB_PASSWORD=postgres
    DB_NAME=achievement_db
    DB_PORT=5432
    MONGO_URI=mongodb://localhost:27017
    MONGO_DB_NAME=achievement_db
    JWT_SECRET=supersecretkey
    PORT=3000
    ```

4.  **Install Dependencies**:
    ```bash
    go mod tidy
    ```

5.  **Run the Application**:
    ```bash
    go run main.go
    ```
    The server will start on `http://localhost:3000`.

## 📖 API Documentation

Swagger documentation is available at:
**[http://localhost:3000/swagger/index.html](http://localhost:3000/swagger/index.html)**

## 🖥️ Web Interface

Access the web interface at:
**[http://localhost:3000/](http://localhost:3000/)**

### Default Roles (Seeded automatically)
*   **Admin**: Can manage users.
*   **Mahasiswa**: Can submit achievements.
*   **Dosen Wali**: Can verify achievements.

*Note: You need to register users via the API (`/api/v1/auth/register`) or seed them manually to log in.*

## 🧪 Testing

Run unit tests:
```bash
go test ./tests/...
```
