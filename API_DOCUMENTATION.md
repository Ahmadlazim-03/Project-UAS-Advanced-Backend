# Student Achievement System API Documentation

## Base URL
```
http://localhost:3000/api/v1
```

## Table of Contents
1. [Authentication](#authentication)
2. [User Management](#user-management)
3. [Student Operations](#student-operations)
4. [Lecturer Operations](#lecturer-operations)
5. [Achievement Operations](#achievement-operations)
6. [Verification Workflow](#verification-workflow)
7. [Reports & Analytics](#reports--analytics)
8. [Response Format](#response-format)
9. [Error Codes](#error-codes)

---

## Authentication

### 1. Login
**Endpoint:** `POST /auth/login`  
**Description:** Authenticate user and receive access token  
**Authorization:** None (Public)

**Request Body:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "b6a94218-bd99-4fcf-92f9-1dd33f355564",
      "username": "admin",
      "email": "admin@university.ac.id",
      "full_name": "System Administrator",
      "is_active": true,
      "role_id": "524c9756-5c56-491b-bb7d-c51b00748367",
      "role": {
        "id": "524c9756-5c56-491b-bb7d-c51b00748367",
        "name": "Admin",
        "description": "System administrator"
      }
    }
  }
}
```

**Error Response (401):**
```json
{
  "status": "error",
  "message": "Invalid credentials",
  "data": null
}
```

---

### 2. Refresh Token
**Endpoint:** `POST /auth/refresh`  
**Description:** Refresh access token using refresh token  
**Authorization:** None (Public)

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Token refreshed successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

---

### 3. Get Profile
**Endpoint:** `GET /auth/profile`  
**Description:** Get current user profile  
**Authorization:** Bearer Token (Required)

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Profile retrieved successfully",
  "data": {
    "id": "b6a94218-bd99-4fcf-92f9-1dd33f355564",
    "username": "admin",
    "email": "admin@university.ac.id",
    "full_name": "System Administrator",
    "is_active": true,
    "role": {
      "id": "524c9756-5c56-491b-bb7d-c51b00748367",
      "name": "Admin",
      "description": "System administrator"
    }
  }
}
```

---

### 4. Logout
**Endpoint:** `POST /auth/logout`  
**Description:** Logout current user  
**Authorization:** Bearer Token (Required)

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Logout successful",
  "data": null
}
```

---

## User Management

### 1. List Users
**Endpoint:** `GET /users?page=1&limit=10`  
**Description:** Get paginated list of all users  
**Authorization:** Bearer Token + Permission: `user:manage`  
**Role:** Admin only

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Users retrieved successfully",
  "data": {
    "users": [
      {
        "id": "b6a94218-bd99-4fcf-92f9-1dd33f355564",
        "username": "admin",
        "email": "admin@university.ac.id",
        "full_name": "System Administrator",
        "is_active": true,
        "role_id": "524c9756-5c56-491b-bb7d-c51b00748367",
        "role": {
          "id": "524c9756-5c56-491b-bb7d-c51b00748367",
          "name": "Admin",
          "description": "System administrator"
        },
        "created_at": "2025-12-03T23:28:28.690522+07:00",
        "updated_at": "2025-12-03T23:28:28.690522+07:00"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 3,
      "total_pages": 1
    }
  }
}
```

---

### 2. Get User Detail
**Endpoint:** `GET /users/:id`  
**Description:** Get specific user details  
**Authorization:** Bearer Token + Permission: `user:read`  
**Role:** Admin

**Success Response (200):**
```json
{
  "status": "success",
  "message": "User retrieved successfully",
  "data": {
    "id": "b6a94218-bd99-4fcf-92f9-1dd33f355564",
    "username": "admin",
    "email": "admin@university.ac.id",
    "full_name": "System Administrator",
    "is_active": true,
    "role_id": "524c9756-5c56-491b-bb7d-c51b00748367",
    "role": {
      "id": "524c9756-5c56-491b-bb7d-c51b00748367",
      "name": "Admin",
      "description": "System administrator"
    },
    "created_at": "2025-12-03T23:28:28.690522+07:00",
    "updated_at": "2025-12-03T23:28:28.690522+07:00"
  }
}
```

---

### 3. Create User
**Endpoint:** `POST /users`  
**Description:** Create new user (automatically creates Student or Lecturer record if applicable)  
**Authorization:** Bearer Token + Permission: `user:create`  
**Role:** Admin only

**Request Body (Create Student User):**
```json
{
  "username": "student001",
  "email": "student001@university.ac.id",
  "password": "student123",
  "full_name": "John Doe",
  "student_id": "STU001",
  "program_study": "Teknik Informatika",
  "academic_year": "2024"
}
```

**Request Body (Create Lecturer User):**
```json
{
  "username": "lecturer001",
  "email": "lecturer001@university.ac.id",
  "password": "lecturer123",
  "full_name": "Dr. Jane Smith",
  "lecturer_id": "LEC001",
  "department": "Computer Science"
}
```

**Request Body (Create Admin User):**
```json
{
  "username": "newadmin",
  "email": "newadmin@university.ac.id",
  "password": "admin123",
  "full_name": "New Administrator"
}
```

**Success Response (201):**
```json
{
  "status": "success",
  "message": "User created successfully",
  "data": {
    "id": "16677ddc-6a81-4341-a04e-1643e9d881eb",
    "username": "student001",
    "email": "student001@university.ac.id",
    "full_name": "John Doe",
    "is_active": true,
    "role_id": "394ff7a7-2653-45e7-9f61-0deff113c345",
    "role": {
      "id": "394ff7a7-2653-45e7-9f61-0deff113c345",
      "name": "Mahasiswa",
      "description": "Student"
    },
    "created_at": "2025-12-03T23:36:06.008006+07:00",
    "updated_at": "2025-12-03T23:36:06.008006+07:00"
  }
}
```

---

### 4. Update User
**Endpoint:** `PUT /users/:id`  
**Description:** Update user information  
**Authorization:** Bearer Token + Permission: `user:update`  
**Role:** Admin only

**Request Body:**
```json
{
  "full_name": "John Doe Updated",
  "email": "newemail@university.ac.id",
  "is_active": true
}
```

**Success Response (200):**
```json
{
  "status": "success",
  "message": "User updated successfully",
  "data": {
    "id": "16677ddc-6a81-4341-a04e-1643e9d881eb",
    "username": "student001",
    "email": "newemail@university.ac.id",
    "full_name": "John Doe Updated",
    "is_active": true,
    "role_id": "394ff7a7-2653-45e7-9f61-0deff113c345",
    "created_at": "2025-12-03T23:36:06.008006+07:00",
    "updated_at": "2025-12-03T23:40:15.123456+07:00"
  }
}
```

---

### 5. Delete User
**Endpoint:** `DELETE /users/:id`  
**Description:** Soft delete user (sets DeletedAt timestamp)  
**Authorization:** Bearer Token + Permission: `user:delete`  
**Role:** Admin only

**Success Response (200):**
```json
{
  "status": "success",
  "message": "User deleted successfully",
  "data": null
}
```

---

### 6. Assign Role
**Endpoint:** `PUT /users/:id/role`  
**Description:** Assign role to user  
**Authorization:** Bearer Token + Permission: `user:manage`  
**Role:** Admin only

**Request Body:**
```json
{
  "role_id": "394ff7a7-2653-45e7-9f61-0deff113c345"
}
```

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Role assigned successfully",
  "data": null
}
```

---

## Student Operations

### 1. List Students
**Endpoint:** `GET /students?page=1&limit=10`  
**Description:** Get paginated list of all students  
**Authorization:** Bearer Token + Permission: `user:read` or `user:manage`  
**Role:** Admin, Dosen Wali

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Students retrieved successfully",
  "data": {
    "students": [
      {
        "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "user_id": "16677ddc-6a81-4341-a04e-1643e9d881eb",
        "student_id": "STU001",
        "program_study": "Teknik Informatika",
        "academic_year": "2024",
        "advisor_id": "f1e2d3c4-b5a6-7890-1234-567890abcdef",
        "user": {
          "id": "16677ddc-6a81-4341-a04e-1643e9d881eb",
          "username": "student001",
          "email": "student001@university.ac.id",
          "full_name": "John Doe"
        },
        "advisor": {
          "id": "f1e2d3c4-b5a6-7890-1234-567890abcdef",
          "lecturer_id": "LEC001",
          "department": "Computer Science",
          "user": {
            "full_name": "Dr. Jane Smith"
          }
        },
        "created_at": "2025-12-03T23:28:28.500000+07:00"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

---

### 2. Get Student Detail
**Endpoint:** `GET /students/:id`  
**Description:** Get specific student details (use user_id)  
**Authorization:** Bearer Token + Permission: `user:read`  
**Role:** Admin, Dosen Wali, Own Student

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Student retrieved successfully",
  "data": {
    "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "user_id": "16677ddc-6a81-4341-a04e-1643e9d881eb",
    "student_id": "STU001",
    "program_study": "Teknik Informatika",
    "academic_year": "2024",
    "advisor_id": "f1e2d3c4-b5a6-7890-1234-567890abcdef",
    "user": {
      "id": "16677ddc-6a81-4341-a04e-1643e9d881eb",
      "username": "student001",
      "email": "student001@university.ac.id",
      "full_name": "John Doe",
      "is_active": true
    },
    "advisor": {
      "id": "f1e2d3c4-b5a6-7890-1234-567890abcdef",
      "lecturer_id": "LEC001",
      "department": "Computer Science",
      "user": {
        "full_name": "Dr. Jane Smith",
        "email": "lecturer001@university.ac.id"
      }
    },
    "created_at": "2025-12-03T23:28:28.500000+07:00"
  }
}
```

---

### 3. Get Student Achievements
**Endpoint:** `GET /students/:id/achievements`  
**Description:** Get all achievements for specific student  
**Authorization:** Bearer Token + Permission: `achievement:read`  
**Role:** Admin, Dosen Wali, Own Student

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Achievements retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "student_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "title": "National Programming Competition Winner",
      "description": "First place in National Coding Championship 2024",
      "achieved_date": "2024-06-15T00:00:00Z",
      "status": "verified",
      "data": {
        "competition_name": "National Coding Championship",
        "competition_level": "national",
        "rank": 1,
        "medal_type": "gold"
      },
      "created_at": "2024-06-20T10:30:00+07:00",
      "updated_at": "2024-06-25T15:45:00+07:00"
    }
  ]
}
```

---

### 4. Assign Advisor
**Endpoint:** `PUT /students/:id/advisor`  
**Description:** Assign lecturer as student's advisor  
**Authorization:** Bearer Token + Permission: `user:manage`  
**Role:** Admin only

**Request Body:**
```json
{
  "advisor_id": "f1e2d3c4-b5a6-7890-1234-567890abcdef"
}
```

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Advisor assigned successfully",
  "data": null
}
```

---

## Lecturer Operations

### 1. List Lecturers
**Endpoint:** `GET /lecturers?page=1&limit=10`  
**Description:** Get paginated list of all lecturers  
**Authorization:** Bearer Token + Permission: `user:read` or `user:manage`  
**Role:** Admin, Dosen Wali

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Lecturers retrieved successfully",
  "data": {
    "lecturers": [
      {
        "id": "f1e2d3c4-b5a6-7890-1234-567890abcdef",
        "user_id": "c1d2e3f4-a5b6-7890-1234-567890abcdef",
        "lecturer_id": "LEC001",
        "department": "Computer Science",
        "user": {
          "id": "c1d2e3f4-a5b6-7890-1234-567890abcdef",
          "username": "lecturer001",
          "email": "lecturer001@university.ac.id",
          "full_name": "Dr. Jane Smith"
        },
        "created_at": "2025-12-03T23:28:28.400000+07:00"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

---

### 2. Get Lecturer Advisees
**Endpoint:** `GET /lecturers/:id/advisees`  
**Description:** Get all students under this lecturer's supervision  
**Authorization:** Bearer Token + Permission: `achievement:verify`  
**Role:** Dosen Wali (own advisees), Admin

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Advisees retrieved successfully",
  "data": [
    {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "user_id": "16677ddc-6a81-4341-a04e-1643e9d881eb",
      "student_id": "STU001",
      "program_study": "Teknik Informatika",
      "academic_year": "2024",
      "advisor_id": "f1e2d3c4-b5a6-7890-1234-567890abcdef",
      "user": {
        "id": "16677ddc-6a81-4341-a04e-1643e9d881eb",
        "username": "student001",
        "email": "student001@university.ac.id",
        "full_name": "John Doe",
        "is_active": true
      },
      "created_at": "2025-12-03T23:28:28.500000+07:00"
    }
  ]
}
```

---

### 3. Get Advisee Achievements
**Endpoint:** `GET /lecturers/advisees/achievements`  
**Description:** Get all achievements from lecturer's advisees  
**Authorization:** Bearer Token + Permission: `achievement:verify`  
**Role:** Dosen Wali only

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Advisee achievements retrieved successfully",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "student_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "title": "National Programming Competition Winner",
      "description": "First place in National Coding Championship 2024",
      "achieved_date": "2024-06-15T00:00:00Z",
      "status": "pending_verification",
      "student": {
        "student_id": "STU001",
        "program_study": "Teknik Informatika",
        "user": {
          "full_name": "John Doe",
          "email": "student001@university.ac.id"
        }
      },
      "data": {
        "competition_name": "National Coding Championship",
        "competition_level": "national",
        "rank": 1,
        "medal_type": "gold"
      },
      "created_at": "2024-06-20T10:30:00+07:00"
    }
  ]
}
```

---

## Achievement Operations

### 1. List Achievements
**Endpoint:** `GET /achievements?page=1&limit=10&status=verified`  
**Description:** Get paginated list of achievements  
**Authorization:** Bearer Token + Permission: `achievement:read`  
**Role:** All authenticated users (see own achievements or all if admin/lecturer)

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)
- `status` (optional): Filter by status (draft, pending_verification, verified, rejected)

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Achievements retrieved successfully",
  "data": {
    "achievements": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "student_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "title": "National Programming Competition Winner",
        "description": "First place in National Coding Championship 2024",
        "achieved_date": "2024-06-15T00:00:00Z",
        "status": "verified",
        "data": {
          "competition_name": "National Coding Championship",
          "competition_level": "national",
          "rank": 1,
          "medal_type": "gold"
        },
        "created_at": "2024-06-20T10:30:00+07:00",
        "updated_at": "2024-06-25T15:45:00+07:00"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 1,
      "total_pages": 1
    }
  }
}
```

---

### 2. Get Achievement Detail
**Endpoint:** `GET /achievements/:id`  
**Description:** Get specific achievement details  
**Authorization:** Bearer Token + Permission: `achievement:read`  
**Role:** Owner, Advisor, Admin

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Achievement retrieved successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "student_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "title": "National Programming Competition Winner",
    "description": "First place in National Coding Championship 2024",
    "achieved_date": "2024-06-15T00:00:00Z",
    "status": "verified",
    "data": {
      "competition_name": "National Coding Championship",
      "competition_level": "national",
      "rank": 1,
      "medal_type": "gold",
      "certificate_url": "https://example.com/certificate.pdf",
      "organizer": "Indonesian Computer Society"
    },
    "verification": {
      "verified_by": "c1d2e3f4-a5b6-7890-1234-567890abcdef",
      "verified_at": "2024-06-25T15:45:00+07:00",
      "comments": "Excellent achievement! Certificate verified."
    },
    "student": {
      "student_id": "STU001",
      "program_study": "Teknik Informatika",
      "user": {
        "full_name": "John Doe",
        "email": "student001@university.ac.id"
      }
    },
    "created_at": "2024-06-20T10:30:00+07:00",
    "updated_at": "2024-06-25T15:45:00+07:00"
  }
}
```

---

### 3. Create Achievement
**Endpoint:** `POST /achievements`  
**Description:** Create new achievement (Students only)  
**Authorization:** Bearer Token + Permission: `achievement:create`  
**Role:** Mahasiswa only

**Request Body (Competition Achievement):**
```json
{
  "title": "National Programming Competition Winner",
  "description": "First place in National Coding Championship 2024",
  "achieved_date": "2024-06-15",
  "data": {
    "competition_name": "National Coding Championship",
    "competition_level": "national",
    "rank": 1,
    "medal_type": "gold",
    "certificate_url": "https://example.com/certificate.pdf",
    "organizer": "Indonesian Computer Society",
    "location": "Jakarta, Indonesia",
    "participants_count": 150
  }
}
```

**Request Body (Publication Achievement):**
```json
{
  "title": "Published Research Paper on AI",
  "description": "Published paper about AI applications in education",
  "achieved_date": "2024-05-10",
  "data": {
    "publication_type": "journal",
    "publication_title": "AI Applications in Modern Education",
    "authors": ["John Doe", "Dr. Jane Smith", "Dr. Robert Brown"],
    "publisher": "IEEE",
    "issn": "1234-5678",
    "doi": "10.1109/example.2024.123456",
    "publication_date": "2024-05-10",
    "journal_name": "IEEE Transactions on Education",
    "volume": "67",
    "issue": "2",
    "pages": "123-145"
  }
}
```

**Request Body (Research Grant Achievement):**
```json
{
  "title": "Research Grant Recipient",
  "description": "Received national research grant for AI project",
  "achieved_date": "2024-03-20",
  "data": {
    "grant_name": "National Research Innovation Grant 2024",
    "funding_agency": "Ministry of Research and Technology",
    "amount": 50000000,
    "currency": "IDR",
    "research_topic": "Artificial Intelligence in Healthcare",
    "duration_months": 12,
    "principal_investigator": "John Doe",
    "co_investigators": ["Dr. Jane Smith"]
  }
}
```

**Success Response (201):**
```json
{
  "status": "success",
  "message": "Achievement created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "student_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "title": "National Programming Competition Winner",
    "description": "First place in National Coding Championship 2024",
    "achieved_date": "2024-06-15T00:00:00Z",
    "status": "draft",
    "data": {
      "competition_name": "National Coding Championship",
      "competition_level": "national",
      "rank": 1,
      "medal_type": "gold"
    },
    "created_at": "2024-06-20T10:30:00+07:00",
    "updated_at": "2024-06-20T10:30:00+07:00"
  }
}
```

**Error Response (400):**
```json
{
  "status": "error",
  "message": "Only students can create achievements",
  "data": null
}
```

---

### 4. Update Achievement
**Endpoint:** `PUT /achievements/:id`  
**Description:** Update achievement details (only draft or rejected status)  
**Authorization:** Bearer Token + Permission: `achievement:update`  
**Role:** Mahasiswa (owner only)

**Request Body:**
```json
{
  "title": "National Programming Competition Winner - Updated",
  "description": "Updated description with more details about the competition",
  "data": {
    "competition_name": "National Coding Championship",
    "competition_level": "national",
    "rank": 1,
    "medal_type": "gold",
    "certificate_url": "https://example.com/new-certificate.pdf",
    "organizer": "Indonesian Computer Society",
    "location": "Jakarta Convention Center, Indonesia",
    "participants_count": 150,
    "team_members": ["John Doe", "Alice Johnson"]
  }
}
```

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Achievement updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "student_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "title": "National Programming Competition Winner - Updated",
    "description": "Updated description with more details about the competition",
    "achieved_date": "2024-06-15T00:00:00Z",
    "status": "draft",
    "data": {
      "competition_name": "National Coding Championship",
      "competition_level": "national",
      "rank": 1,
      "medal_type": "gold",
      "certificate_url": "https://example.com/new-certificate.pdf",
      "organizer": "Indonesian Computer Society",
      "location": "Jakarta Convention Center, Indonesia",
      "participants_count": 150,
      "team_members": ["John Doe", "Alice Johnson"]
    },
    "created_at": "2024-06-20T10:30:00+07:00",
    "updated_at": "2024-06-21T14:20:00+07:00"
  }
}
```

---

### 5. Delete Achievement
**Endpoint:** `DELETE /achievements/:id`  
**Description:** Delete achievement (only draft status)  
**Authorization:** Bearer Token + Permission: `achievement:delete`  
**Role:** Mahasiswa (owner only), Admin

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Achievement deleted successfully",
  "data": null
}
```

---

## Verification Workflow

### 1. Submit for Verification
**Endpoint:** `POST /achievements/:id/submit`  
**Description:** Submit achievement for lecturer verification  
**Authorization:** Bearer Token + Permission: `achievement:update`  
**Role:** Mahasiswa (owner only)

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Achievement submitted for verification",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "student_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "title": "National Programming Competition Winner",
    "status": "pending_verification",
    "submitted_at": "2024-06-22T09:15:00+07:00"
  }
}
```

**Error Response (400):**
```json
{
  "status": "error",
  "message": "Achievement is already submitted",
  "data": null
}
```

---

### 2. Verify Achievement
**Endpoint:** `POST /achievements/:id/verify`  
**Description:** Verify and approve achievement  
**Authorization:** Bearer Token + Permission: `achievement:verify`  
**Role:** Dosen Wali (for own advisees), Admin

**Request Body:**
```json
{
  "comments": "Excellent achievement! Certificate and supporting documents verified. Keep up the great work!"
}
```

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Achievement verified successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "student_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "title": "National Programming Competition Winner",
    "status": "verified",
    "verification": {
      "verified_by": "c1d2e3f4-a5b6-7890-1234-567890abcdef",
      "verified_at": "2024-06-25T15:45:00+07:00",
      "comments": "Excellent achievement! Certificate and supporting documents verified. Keep up the great work!"
    },
    "verified_by": {
      "id": "c1d2e3f4-a5b6-7890-1234-567890abcdef",
      "full_name": "Dr. Jane Smith",
      "email": "lecturer001@university.ac.id"
    }
  }
}
```

**Error Response (400):**
```json
{
  "status": "error",
  "message": "Achievement is not pending verification",
  "data": null
}
```

---

### 3. Reject Achievement
**Endpoint:** `POST /achievements/:id/reject`  
**Description:** Reject achievement with reason  
**Authorization:** Bearer Token + Permission: `achievement:verify`  
**Role:** Dosen Wali (for own advisees), Admin

**Request Body:**
```json
{
  "reason": "Missing supporting documents. Please provide:\n1. Official certificate from organizer\n2. Competition announcement/flyer\n3. Photo during the event\n\nPlease resubmit after uploading the required documents."
}
```

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Achievement rejected",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "student_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
    "title": "National Programming Competition Winner",
    "status": "rejected",
    "rejection": {
      "rejected_by": "c1d2e3f4-a5b6-7890-1234-567890abcdef",
      "rejected_at": "2024-06-23T11:30:00+07:00",
      "reason": "Missing supporting documents. Please provide:\n1. Official certificate from organizer\n2. Competition announcement/flyer\n3. Photo during the event\n\nPlease resubmit after uploading the required documents."
    },
    "rejected_by": {
      "id": "c1d2e3f4-a5b6-7890-1234-567890abcdef",
      "full_name": "Dr. Jane Smith",
      "email": "lecturer001@university.ac.id"
    }
  }
}
```

---

## Reports & Analytics

### 1. Get Statistics
**Endpoint:** `GET /reports/statistics`  
**Description:** Get system-wide statistics  
**Authorization:** Bearer Token + Permission: `report:read`  
**Role:** Admin, Dosen Wali

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Statistics retrieved successfully",
  "data": {
    "total_students": 125,
    "total_lecturers": 15,
    "total_achievements": 342,
    "achievements_by_status": {
      "draft": 45,
      "pending_verification": 28,
      "verified": 245,
      "rejected": 24
    },
    "achievements_by_level": {
      "international": 35,
      "national": 128,
      "regional": 95,
      "university": 84
    },
    "achievements_by_type": {
      "competition": 215,
      "publication": 87,
      "research_grant": 25,
      "certification": 15
    },
    "recent_verifications": 12,
    "pending_verifications": 28,
    "top_performers": [
      {
        "student_id": "STU001",
        "full_name": "John Doe",
        "program_study": "Teknik Informatika",
        "total_achievements": 15,
        "verified_achievements": 12
      }
    ]
  }
}
```

---

### 2. Get Student Report
**Endpoint:** `GET /reports/student/:id`  
**Description:** Get detailed report for specific student  
**Authorization:** Bearer Token + Permission: `report:read`  
**Role:** Admin, Dosen Wali (for own advisees), Own Student

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Student report retrieved successfully",
  "data": {
    "student": {
      "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "student_id": "STU001",
      "full_name": "John Doe",
      "email": "student001@university.ac.id",
      "program_study": "Teknik Informatika",
      "academic_year": "2024",
      "advisor": {
        "lecturer_id": "LEC001",
        "full_name": "Dr. Jane Smith",
        "department": "Computer Science"
      }
    },
    "summary": {
      "total_achievements": 15,
      "verified_achievements": 12,
      "pending_achievements": 2,
      "rejected_achievements": 1,
      "draft_achievements": 0
    },
    "achievements_by_level": {
      "international": 3,
      "national": 7,
      "regional": 3,
      "university": 2
    },
    "achievements_by_type": {
      "competition": 10,
      "publication": 3,
      "research_grant": 1,
      "certification": 1
    },
    "timeline": [
      {
        "year": 2024,
        "month": 6,
        "achievements_count": 3,
        "verified_count": 2
      },
      {
        "year": 2024,
        "month": 5,
        "achievements_count": 2,
        "verified_count": 2
      }
    ],
    "recent_achievements": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "National Programming Competition Winner",
        "achieved_date": "2024-06-15T00:00:00Z",
        "status": "verified",
        "verification_date": "2024-06-25T15:45:00+07:00"
      }
    ]
  }
}
```

---

## Response Format

### Success Response
```json
{
  "status": "success",
  "message": "Operation completed successfully",
  "data": { ... }
}
```

### Error Response
```json
{
  "status": "error",
  "message": "Error description",
  "data": null
}
```

### Paginated Response
```json
{
  "status": "success",
  "message": "Data retrieved successfully",
  "data": {
    "items": [ ... ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 100,
      "total_pages": 10
    }
  }
}
```

---

## Error Codes

| Status Code | Description |
|------------|-------------|
| 200 | OK - Request successful |
| 201 | Created - Resource created successfully |
| 400 | Bad Request - Invalid request data |
| 401 | Unauthorized - Missing or invalid token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource not found |
| 500 | Internal Server Error - Server error |

---

## Authentication Flow

1. **Login** → Receive `access_token` and `refresh_token`
2. **Use access token** → Include in `Authorization: Bearer <token>` header
3. **Token expires** → Use refresh token to get new access token
4. **Logout** → Invalidate tokens

---

## Permission Matrix

| Role | Permissions |
|------|------------|
| **Admin** | All permissions (full access) |
| **Mahasiswa** | `achievement:create`, `achievement:read`, `achievement:update`, `achievement:delete` (own only) |
| **Dosen Wali** | `achievement:read`, `achievement:verify`, `report:read`, `user:read` |

---

## Status Flow

### Achievement Status Lifecycle:
```
draft → pending_verification → verified
                              ↘ rejected → (can resubmit)
```

1. **draft** - Initial state when created
2. **pending_verification** - Submitted for verification
3. **verified** - Approved by lecturer
4. **rejected** - Rejected by lecturer (can be edited and resubmitted)

---

## Notes

- All dates are in ISO 8601 format
- All timestamps include timezone (+07:00 for WIB)
- UUIDs are used for all IDs
- Pagination defaults: page=1, limit=10
- Maximum limit per request: 100
- Soft delete is used (records are not permanently deleted)
- MongoDB is used for achievement data storage
- PostgreSQL is used for user, role, and relational data

---

## Testing Credentials

### Admin
```
Username: admin
Password: admin123
```

### Student
```
Username: student001
Password: student123
```

### Lecturer
```
Username: lecturer001
Password: lecturer123
```

---

**API Version:** v1  
**Last Updated:** December 3, 2025  
**Base URL:** http://localhost:3000/api/v1
