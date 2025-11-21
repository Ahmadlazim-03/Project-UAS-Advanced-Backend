# 🎓 Student Achievement System

Fullstack application untuk manajemen prestasi mahasiswa dengan Go Fiber backend dan SvelteKit frontend.

## 🚀 Features

*   **Multi-role Authentication**: Admin, Mahasiswa, Dosen Wali
*   **Achievement Management**: Create, Read, Update, Delete achievements
*   **Verification System**: Dosen dapat verify/reject achievements
*   **User Management**: Admin dapat mengelola users
*   **Statistics Dashboard**: Real-time metrics dan reports
*   **Hybrid Database**: PostgreSQL + MongoDB
*   **RESTful API**: Fully documented dengan Swagger
*   **Modern Frontend**: SvelteKit 5 dengan TypeScript
*   **Responsive Design**: Mobile-friendly interface

## 🛠 Tech Stack

### Backend
*   **Language**: Go (Golang) 1.21+
*   **Framework**: Fiber v2
*   **Databases**: PostgreSQL (GORM) + MongoDB
*   **Authentication**: JWT (JSON Web Tokens)
*   **Documentation**: Swagger/OpenAPI
*   **Deployment**: Vercel Serverless Functions

### Frontend
*   **Framework**: SvelteKit 5 with Runes
*   **Language**: TypeScript
*   **Styling**: Tailwind CSS v4
*   **State Management**: Svelte Stores
*   **Build Tool**: Vite
*   **Deployment**: Vercel

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

### Features by Role

#### All Users
- **Register**: Create new account with username, email, full name, password, and role selection
- **Login**: Authenticate with username and password
- **Dashboard**: View achievement statistics and recent activities
- **My Achievements**: Manage personal achievements with full CRUD operations
- **Statistics**: View data visualizations (charts) of achievements by status and type

#### Mahasiswa (Student)
- ✅ Create new achievements with title, type, description, points, and tags
- ✅ View all personal achievements in table format
- ✅ Submit draft achievements for verification
- ✅ Delete draft achievements
- ✅ View detailed information for each achievement

#### Dosen Wali (Advisor)
- ✅ All Mahasiswa features
- ✅ **Verify Achievements**: View and verify/reject submitted achievements from students
- ✅ Provide rejection notes when rejecting achievements

#### Admin
- ✅ **User Management**: View all users, toggle user active/inactive status
- ✅ Access to all statistics and reports
- ✅ Full system oversight

### Technology Stack
- **Frontend**: Vanilla JavaScript (ES6+) with Bootstrap 5
- **Charts**: Chart.js for data visualization
- **Icons**: Bootstrap Icons
- **State Management**: localStorage for authentication
- **API Communication**: Fetch API with async/await

### Key Features
1. **Role-Based UI**: Navigation menu adapts based on user role
2. **Real-time Updates**: Data refreshes after each action
3. **Responsive Design**: Works on desktop, tablet, and mobile
4. **Toast Notifications**: User-friendly feedback for all actions
5. **Modal Dialogs**: Add/view achievements without page refresh
6. **Data Visualization**: Interactive charts for statistics
7. **Optimized Performance**: GZIP compression, connection pooling, database indexes

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
