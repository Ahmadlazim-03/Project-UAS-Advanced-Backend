# Student Achievement System - API Documentation

**Base URL:** `http://localhost:3000/api/v1`  
**Version:** 1.1 🆕  
**Last Updated:** December 13, 2025

## 🎉 What's New in v1.1
- ✅ **Notification System** - 5 new endpoints for real-time notifications
- ✅ **Advanced Analytics** - Top students leaderboard, period analysis, competition level distribution
- ✅ **Auto-Notifications** - Automatic notifications when achievements are submitted, verified, or rejected
- ✅ **Enhanced Reporting** - New analytics endpoints for better insights
- ✅ **100% Compliance** - All functional requirements fully implemented

---

## Table of Contents
1. [Authentication](#authentication)
2. [Users Management](#users-management)
3. [Roles](#roles)
4. [Achievements](#achievements)
5. [Students](#students)
6. [Lecturers](#lecturers)
7. [Reports & Analytics](#reports--analytics)
8. [Notifications](#notifications) 🆕
9. [File Management](#file-management)
10. [Test Credentials](#test-credentials)

---

## Authentication

### 1. Login
**Endpoint:** `POST /auth/login`  
**Access:** Public  
**Description:** Authenticate user and receive JWT token

**Request Body:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response (200 OK):**
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
      "role": "Admin",
      "permissions": [
        "user:create",
        "user:read",
        "user:update",
        "user:delete",
        "user:manage",
        "achievement:create",
        "achievement:read",
        "achievement:update",
        "achievement:delete",
        "achievement:verify",
        "report:read"
      ]
    }
  }
}
```

---

### 2. Refresh Token
**Endpoint:** `POST /auth/refresh`  
**Access:** Public  
**Description:** Refresh expired access token

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200 OK):**
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

### 3. Logout
**Endpoint:** `POST /auth/logout`  
**Access:** Protected (Requires Bearer Token)  
**Description:** Logout user and invalidate token

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Logged out successfully"
}
```

---

### 4. Get Profile
**Endpoint:** `GET /auth/profile`  
**Access:** Protected (Requires Bearer Token)  
**Description:** Get current user profile

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Profile retrieved",
  "data": {
    "id": "b6a94218-bd99-4fcf-92f9-1dd33f355564",
    "username": "admin",
    "email": "admin@university.ac.id",
    "full_name": "System Administrator",
    "role": "Admin",
    "permissions": [
      "user:create",
      "user:read",
      "user:update",
      "user:delete",
      "user:manage",
      "achievement:create",
      "achievement:read",
      "achievement:update",
      "achievement:delete",
      "achievement:verify",
      "report:read"
    ]
  }
}
```

---

## Users Management

### 5. List Users
**Endpoint:** `GET /users`  
**Access:** Protected (Admin Only - `user:manage` permission)  
**Description:** Get paginated list of all users

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)
- `search` (optional): Search by username, email, or full name

**Response (200 OK):**
```json
{
  "status": "success",
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 3,
    "total_pages": 1,
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
            "description": "System administrator",
            "created_at": "2025-12-03T23:28:28.273556+07:00"
          },
          "created_at": "2025-12-03T23:28:28.690522+07:00",
          "updated_at": "2025-12-03T23:35:55.241898+07:00"
        },
        {
          "id": "2ce4318a-8e5d-47a2-a4b6-c32e29b723d1",
          "username": "student001",
          "email": "student001@university.ac.id",
          "full_name": "John Doe",
          "is_active": true,
          "role_id": "394ff7a7-2653-45e7-9f61-0deff113c345",
          "role": {
            "id": "394ff7a7-2653-45e7-9f61-0deff113c345",
            "name": "Mahasiswa",
            "description": "Student",
            "created_at": "2025-12-03T23:28:28.27836+07:00"
          },
          "created_at": "2025-12-13T10:12:22.597686+07:00",
          "updated_at": "2025-12-13T10:12:22.597686+07:00"
        },
        {
          "id": "8210511d-a5cb-4134-8d3d-04aa70d64e21",
          "username": "lecturer001",
          "email": "lecturer001@university.ac.id",
          "full_name": "Dr. Jane Smith",
          "is_active": true,
          "role_id": "7be08b8a-030e-4d2a-bef4-b90d81479b0a",
          "role": {
            "id": "7be08b8a-030e-4d2a-bef4-b90d81479b0a",
            "name": "Dosen Wali",
            "description": "Academic advisor",
            "created_at": "2025-12-03T23:28:28.276125+07:00"
          },
          "created_at": "2025-12-13T10:12:22.615454+07:00",
          "updated_at": "2025-12-13T10:12:22.615454+07:00"
        }
      ]
    }
  }
}
```

---

### 6. Get User by ID
**Endpoint:** `GET /users/:id`  
**Access:** Protected (Requires `user:read` permission)  
**Description:** Get specific user details by ID

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "User retrieved successfully",
  "data": {
    "id": "2ce4318a-8e5d-47a2-a4b6-c32e29b723d1",
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
    "created_at": "2025-12-13T10:12:22.597686+07:00",
    "updated_at": "2025-12-13T10:12:22.597686+07:00"
  }
}
```

---

### 7. Create User
**Endpoint:** `POST /users`  
**Access:** Protected (Admin Only - `user:create` permission)  
**Description:** Create a new user

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Request Body (Student):**
```json
{
  "username": "student002",
  "email": "student002@university.ac.id",
  "password": "student123",
  "full_name": "Jane Smith",
  "role_id": "394ff7a7-2653-45e7-9f61-0deff113c345",
  "student_id": "STU002",
  "program_study": "Computer Science",
  "academic_year": "2024"
}
```

**Request Body (Lecturer):**
```json
{
  "username": "lecturer002",
  "email": "lecturer002@university.ac.id",
  "password": "lecturer123",
  "full_name": "Dr. John Doe",
  "role_id": "7be08b8a-030e-4d2a-bef4-b90d81479b0a",
  "lecturer_id": "LEC002",
  "department": "Computer Science"
}
```

**Response (201 Created):**
```json
{
  "status": "success",
  "message": "User created successfully",
  "data": {
    "id": "new-uuid-here",
    "username": "student002",
    "email": "student002@university.ac.id",
    "full_name": "Jane Smith",
    "is_active": true,
    "role_id": "394ff7a7-2653-45e7-9f61-0deff113c345",
    "created_at": "2025-12-13T12:00:00+07:00"
  }
}
```

---

### 8. Update User
**Endpoint:** `PUT /users/:id`  
**Access:** Protected (Admin Only - `user:update` permission)  
**Description:** Update user information

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Request Body:**
```json
{
  "email": "newemail@university.ac.id",
  "full_name": "Updated Name",
  "is_active": true
}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "User updated successfully",
  "data": {
    "id": "2ce4318a-8e5d-47a2-a4b6-c32e29b723d1",
    "username": "student001",
    "email": "newemail@university.ac.id",
    "full_name": "Updated Name",
    "is_active": true,
    "updated_at": "2025-12-13T12:30:00+07:00"
  }
}
```

---

### 9. Delete User
**Endpoint:** `DELETE /users/:id`  
**Access:** Protected (Admin Only - `user:delete` permission)  
**Description:** Soft delete user (sets deleted_at timestamp)

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "User deleted successfully"
}
```

---

### 10. Assign Role to User
**Endpoint:** `PUT /users/:id/role`  
**Access:** Protected (Admin Only - `user:manage` permission)  
**Description:** Change user's role

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Request Body:**
```json
{
  "role_id": "7be08b8a-030e-4d2a-bef4-b90d81479b0a"
}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Role assigned successfully",
  "data": {
    "id": "2ce4318a-8e5d-47a2-a4b6-c32e29b723d1",
    "username": "student001",
    "role_id": "7be08b8a-030e-4d2a-bef4-b90d81479b0a",
    "role": {
      "name": "Dosen Wali"
    }
  }
}
```

---

## Roles

### 11. List Roles
**Endpoint:** `GET /roles`  
**Access:** Protected (Requires Bearer Token)  
**Description:** Get all available roles for user creation

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Roles retrieved successfully",
  "data": {
    "roles": [
      {
        "id": "524c9756-5c56-491b-bb7d-c51b00748367",
        "name": "Admin",
        "description": "System administrator",
        "created_at": "2025-12-03T23:28:28.273556+07:00"
      },
      {
        "id": "394ff7a7-2653-45e7-9f61-0deff113c345",
        "name": "Mahasiswa",
        "description": "Student",
        "created_at": "2025-12-03T23:28:28.27836+07:00"
      },
      {
        "id": "7be08b8a-030e-4d2a-bef4-b90d81479b0a",
        "name": "Dosen Wali",
        "description": "Academic advisor",
        "created_at": "2025-12-03T23:28:28.276125+07:00"
      }
    ]
  }
}
```

---

## Achievements

### 12. List Achievements
**Endpoint:** `GET /achievements`  
**Access:** Protected (Requires `achievement:read` permission)  
**Description:** Get paginated list of achievements (filtered by role)

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)
- `type` (optional): Filter by achievement type (competition, publication, organization, certification, academic, other)
- `status` (optional): Filter by status (draft, submitted, verified, rejected, deleted)
- `search` (optional): Search in title and description

**Response (200 OK):**
```json
{
  "status": "success",
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1,
    "data": {
      "achievements": [
        {
          "id": "00000000-0000-0000-0000-000000000000",
          "mongo_achievement_id": "693d6046d6115cb20e966b7c",
          "student_id": "00000000-0000-0000-0000-000000000000",
          "achievement_type": "competition",
          "title": "International Programming Contest",
          "description": "Won 1st place at ICPC",
          "status": "draft",
          "points": 400,
          "achieved_date": "2024-11-15T00:00:00Z",
          "submitted_at": null,
          "verified_at": null,
          "verified_by": null,
          "rejection_note": "",
          "tags": [
            "programming",
            "algorithm",
            "international"
          ],
          "attachments": [],
          "details": {
            "competition_name": "ICPC World Finals",
            "competition_level": "international",
            "rank": 1,
            "medal_type": "gold",
            "event_date": "2024-11-15T00:00:00Z",
            "location": "Jakarta",
            "organizer": "ACM"
          },
          "student": {
            "id": "00000000-0000-0000-0000-000000000000",
            "name": "Unknown Student"
          },
          "created_at": "2025-12-13T19:47:02.524817+07:00",
          "updated_at": "2025-12-13T19:47:02.524817+07:00"
        }
      ]
    }
  }
}
```

---

### 13. Get Achievement by ID
**Endpoint:** `GET /achievements/:id`  
**Access:** Protected (Requires `achievement:read` permission)  
**Description:** Get specific achievement details

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Achievement retrieved successfully",
  "data": {
    "id": "693d6046d6115cb20e966b7c",
    "student_id": "00000000-0000-0000-0000-000000000000",
    "achievement_type": "competition",
    "title": "International Programming Contest",
    "description": "Won 1st place at ICPC",
    "details": {
      "competition_name": "ICPC World Finals",
      "competition_level": "international",
      "rank": 1,
      "medal_type": "gold",
      "event_date": "2024-11-15T00:00:00Z",
      "location": "Jakarta",
      "organizer": "ACM"
    },
    "attachments": [],
    "tags": [
      "programming",
      "algorithm",
      "international"
    ],
    "points": 400,
    "created_at": "2025-12-13T12:47:02.473Z",
    "updated_at": "2025-12-13T12:47:02.473Z"
  }
}
```

---

### 14. Create Achievement
**Endpoint:** `POST /achievements`  
**Access:** Protected (Mahasiswa - `achievement:create` permission)  
**Description:** Create new achievement (supports 6 types with auto point calculation)

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Request Body - Competition Achievement:**
```json
{
  "achievement_type": "competition",
  "title": "International Programming Contest",
  "description": "Won 1st place at ICPC",
  "achieved_date": "2024-11-15",
  "data": {
    "competition_name": "ICPC World Finals",
    "competition_level": "international",
    "rank": 1,
    "medal_type": "gold",
    "event_date": "2024-11-15",
    "location": "Jakarta",
    "organizer": "ACM"
  },
  "tags": ["programming", "algorithm", "international"],
  "attachments": [
    {
      "file_name": "certificate.pdf",
      "file_url": "/uploads/certificate.pdf",
      "file_type": "application/pdf"
    }
  ]
}
```

**Request Body - Publication Achievement:**
```json
{
  "achievement_type": "publication",
  "title": "Machine Learning Research Paper",
  "description": "Published in IEEE Transactions",
  "achieved_date": "2024-10-20",
  "data": {
    "publication_type": "journal",
    "publication_title": "Deep Learning for NLP",
    "authors": ["John Doe", "Jane Smith"],
    "publisher": "IEEE Transactions",
    "issn": "1234-5678",
    "event_date": "2024-10-20"
  },
  "tags": ["research", "machine-learning", "publication"],
  "attachments": [
    {
      "file_name": "paper.pdf",
      "file_url": "/uploads/paper.pdf",
      "file_type": "application/pdf"
    }
  ]
}
```

**Request Body - Organization Achievement:**
```json
{
  "achievement_type": "organization",
  "title": "Student Council President",
  "description": "Led student council for 2024",
  "achieved_date": "2024-01-01",
  "data": {
    "organization_name": "Student Council",
    "position": "President",
    "period_start": "2024-01-01",
    "period_end": "2024-12-31",
    "location": "University Campus"
  },
  "tags": ["leadership", "organization"]
}
```

**Request Body - Certification Achievement:**
```json
{
  "achievement_type": "certification",
  "title": "AWS Solutions Architect",
  "description": "Professional AWS certification",
  "achieved_date": "2024-03-15",
  "data": {
    "certification_name": "AWS Solutions Architect",
    "issued_by": "Amazon Web Services",
    "certification_number": "AWS-123456",
    "valid_until": "2027-03-15",
    "event_date": "2024-03-15"
  },
  "tags": ["cloud", "AWS", "certification"]
}
```

**Request Body - Academic Achievement:**
```json
{
  "achievement_type": "academic",
  "title": "Dean's List",
  "description": "Achieved GPA 3.85",
  "achieved_date": "2024-06-30",
  "data": {
    "score": 3.85,
    "custom_fields": {
      "semester": "6",
      "major": "Computer Science",
      "dean_list": true
    }
  },
  "tags": ["academic", "GPA", "dean-list"]
}
```

**Request Body - Other Achievement:**
```json
{
  "achievement_type": "other",
  "title": "Community Service",
  "description": "Volunteered teaching programming",
  "achieved_date": "2024-05-10",
  "data": {
    "custom_fields": {
      "type": "Community Service",
      "hours": 40,
      "activity": "Teaching Programming"
    }
  },
  "tags": ["community-service", "volunteer"]
}
```

**Response (201 Created):**
```json
{
  "status": "success",
  "message": "Achievement created successfully",
  "data": {
    "id": "693d6046d6115cb20e966b7c",
    "student_id": "00000000-0000-0000-0000-000000000000",
    "achievement_type": "competition",
    "title": "International Programming Contest",
    "description": "Won 1st place at ICPC",
    "details": {
      "competition_name": "ICPC World Finals",
      "competition_level": "international",
      "rank": 1,
      "medal_type": "gold",
      "event_date": "2024-11-15T00:00:00Z",
      "location": "Jakarta",
      "organizer": "ACM"
    },
    "attachments": [
      {
        "file_name": "certificate.pdf",
        "file_url": "/uploads/certificate.pdf",
        "file_type": "application/pdf",
        "uploaded_at": "2025-12-13T12:47:02.473Z"
      }
    ],
    "tags": [
      "programming",
      "algorithm",
      "international"
    ],
    "points": 400,
    "created_at": "2025-12-13T12:47:02.473Z",
    "updated_at": "2025-12-13T12:47:02.473Z"
  }
}
```

**Point Calculation:**
- Competition: 100 (base) + 200 (international) + 100 (1st rank) = **400 points**
- Publication: 150 (base) + 100 (journal) = **250 points**
- Organization: **50 points**
- Certification: **75 points**
- Academic: **25 points**
- Other: **10 points**

---

### 15. Update Achievement
**Endpoint:** `PUT /achievements/:id`  
**Access:** Protected (Mahasiswa - `achievement:update` permission)  
**Description:** Update achievement (only if status is draft or rejected)

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Request Body:**
```json
{
  "title": "Updated Title",
  "description": "Updated description",
  "achievement_type": "competition",
  "achieved_date": "2024-11-16",
  "data": {
    "competition_name": "Updated Contest Name",
    "competition_level": "national",
    "rank": 2,
    "medal_type": "silver"
  },
  "tags": ["programming", "national"],
  "attachments": [
    {
      "file_name": "new_certificate.pdf",
      "file_url": "/uploads/new_certificate.pdf",
      "file_type": "application/pdf"
    }
  ]
}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Achievement updated successfully",
  "data": {
    "id": "693d6046d6115cb20e966b7c",
    "achievement_type": "competition",
    "title": "Updated Title",
    "description": "Updated description",
    "points": 275,
    "updated_at": "2025-12-13T13:00:00.000Z"
  }
}
```

---

### 16. Delete Achievement
**Endpoint:** `DELETE /achievements/:id`  
**Access:** Protected (Mahasiswa - `achievement:delete` permission)  
**Description:** Soft delete achievement (sets status to 'deleted')

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Achievement deleted successfully"
}
```

---

### 17. Submit Achievement for Verification
**Endpoint:** `POST /achievements/:id/submit`  
**Access:** Protected (Mahasiswa - `achievement:update` permission)  
**Description:** Submit achievement for lecturer verification (changes status from draft to submitted)

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Achievement submitted for verification",
  "data": {
    "id": "693d6046d6115cb20e966b7c",
    "status": "submitted",
    "submitted_at": "2025-12-13T19:47:34.0705938+07:00"
  }
}
```

---

### 18. Verify Achievement
**Endpoint:** `POST /achievements/:id/verify`  
**Access:** Protected (Dosen Wali - `achievement:verify` permission)  
**Description:** Verify achievement (only for advisees)

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Achievement verified successfully",
  "data": {
    "id": "693d6046d6115cb20e966b7c",
    "status": "verified",
    "verified_at": "2025-12-13T20:00:00+07:00",
    "verified_by": "8210511d-a5cb-4134-8d3d-04aa70d64e21"
  }
}
```

---

### 19. Reject Achievement
**Endpoint:** `POST /achievements/:id/reject`  
**Access:** Protected (Dosen Wali - `achievement:verify` permission)  
**Description:** Reject achievement with note (only for advisees)

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Request Body:**
```json
{
  "note": "Insufficient evidence provided. Please attach certificate."
}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Achievement rejected",
  "data": {
    "id": "693d6046d6115cb20e966b7c",
    "status": "rejected",
    "rejection_note": "Insufficient evidence provided. Please attach certificate.",
    "rejected_at": "2025-12-13T20:05:00+07:00"
  }
}
```

---

### 20. Get Achievement History
**Endpoint:** `GET /achievements/:id/history`  
**Access:** Protected (Requires `achievement:read` permission)  
**Description:** Get status change history of achievement

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Achievement history retrieved successfully",
  "data": {
    "achievement_id": "693d6046d6115cb20e966b7c",
    "current_status": "verified",
    "history": [
      {
        "id": "history-uuid-1",
        "achievement_reference_id": "ref-uuid",
        "previous_status": "draft",
        "new_status": "submitted",
        "changed_by": "2ce4318a-8e5d-47a2-a4b6-c32e29b723d1",
        "changed_by_user": {
          "username": "student001",
          "full_name": "John Doe"
        },
        "note": "",
        "created_at": "2025-12-13T19:47:34+07:00"
      },
      {
        "id": "history-uuid-2",
        "achievement_reference_id": "ref-uuid",
        "previous_status": "submitted",
        "new_status": "verified",
        "changed_by": "8210511d-a5cb-4134-8d3d-04aa70d64e21",
        "changed_by_user": {
          "username": "lecturer001",
          "full_name": "Dr. Jane Smith"
        },
        "note": "Approved",
        "created_at": "2025-12-13T20:00:00+07:00"
      }
    ]
  }
}
```

---

### 21. Upload Achievement Attachments
**Endpoint:** `POST /achievements/:id/attachments`  
**Access:** Protected (Requires `achievement:create` permission)  
**Description:** Upload additional attachments to achievement

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Request Body:**
```json
{
  "attachments": [
    {
      "file_name": "additional_doc.pdf",
      "file_url": "/uploads/additional_doc.pdf",
      "file_type": "application/pdf"
    },
    {
      "file_name": "photo.jpg",
      "file_url": "/uploads/photo.jpg",
      "file_type": "image/jpeg"
    }
  ]
}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Attachments added successfully",
  "data": {
    "achievement_id": "693d6046d6115cb20e966b7c",
    "attachments": [
      {
        "file_name": "certificate.pdf",
        "file_url": "/uploads/certificate.pdf",
        "file_type": "application/pdf",
        "uploaded_at": "2025-12-13T12:47:02.473Z"
      },
      {
        "file_name": "additional_doc.pdf",
        "file_url": "/uploads/additional_doc.pdf",
        "file_type": "application/pdf",
        "uploaded_at": "2025-12-13T20:10:00.000Z"
      },
      {
        "file_name": "photo.jpg",
        "file_url": "/uploads/photo.jpg",
        "file_type": "image/jpeg",
        "uploaded_at": "2025-12-13T20:10:00.000Z"
      }
    ]
  }
}
```

---

## Students

### 22. List Students
**Endpoint:** `GET /students`  
**Access:** Protected (Requires `user:read` or `user:manage` permission)  
**Description:** Get paginated list of all students

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)
- `search` (optional): Search by name or student ID

**Response (200 OK):**
```json
{
  "status": "success",
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1,
    "data": {
      "students": [
        {
          "id": "student-uuid",
          "user_id": "2ce4318a-8e5d-47a2-a4b6-c32e29b723d1",
          "student_id": "STU001",
          "program_study": "Computer Science",
          "academic_year": "2024",
          "advisor_id": "d1359b6d-4b28-4319-b09c-66377fe34a4b",
          "user": {
            "id": "2ce4318a-8e5d-47a2-a4b6-c32e29b723d1",
            "username": "student001",
            "email": "student001@university.ac.id",
            "full_name": "John Doe",
            "is_active": true
          },
          "advisor": {
            "id": "d1359b6d-4b28-4319-b09c-66377fe34a4b",
            "lecturer_id": "LEC001",
            "user": {
              "full_name": "Dr. Jane Smith"
            }
          },
          "created_at": "2025-12-13T10:12:22.617674+07:00"
        }
      ]
    }
  }
}
```

---

### 23. Get Student by ID
**Endpoint:** `GET /students/:id`  
**Access:** Protected (Requires `user:read` permission)  
**Description:** Get specific student details

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Student retrieved successfully",
  "data": {
    "id": "student-uuid",
    "user_id": "2ce4318a-8e5d-47a2-a4b6-c32e29b723d1",
    "student_id": "STU001",
    "program_study": "Computer Science",
    "academic_year": "2024",
    "advisor_id": "d1359b6d-4b28-4319-b09c-66377fe34a4b",
    "user": {
      "username": "student001",
      "email": "student001@university.ac.id",
      "full_name": "John Doe",
      "is_active": true
    },
    "advisor": {
      "lecturer_id": "LEC001",
      "user": {
        "full_name": "Dr. Jane Smith"
      }
    }
  }
}
```

---

### 24. Get Student Achievements
**Endpoint:** `GET /students/:id/achievements`  
**Access:** Protected (Requires `achievement:read` permission)  
**Description:** Get all achievements of a specific student

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)
- `status` (optional): Filter by status

**Response (200 OK):**
```json
{
  "status": "success",
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 3,
    "total_pages": 1,
    "data": {
      "student": {
        "id": "student-uuid",
        "student_id": "STU001",
        "full_name": "John Doe"
      },
      "achievements": [
        {
          "id": "693d6046d6115cb20e966b7c",
          "achievement_type": "competition",
          "title": "International Programming Contest",
          "status": "verified",
          "points": 400,
          "verified_at": "2025-12-13T20:00:00+07:00"
        }
      ],
      "total_points": 400
    }
  }
}
```

---

### 25. Assign Advisor to Student
**Endpoint:** `PUT /students/:id/advisor`  
**Access:** Protected (Admin - `user:manage` permission)  
**Description:** Assign or change student's academic advisor

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: application/json
```

**Request Body:**
```json
{
  "advisor_id": "d1359b6d-4b28-4319-b09c-66377fe34a4b"
}
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Advisor assigned successfully",
  "data": {
    "student_id": "student-uuid",
    "advisor_id": "d1359b6d-4b28-4319-b09c-66377fe34a4b",
    "advisor": {
      "lecturer_id": "LEC001",
      "user": {
        "full_name": "Dr. Jane Smith"
      }
    }
  }
}
```

---

## Lecturers

### 26. List Lecturers
**Endpoint:** `GET /lecturers`  
**Access:** Protected (Requires `user:read` or `user:manage` permission)  
**Description:** Get paginated list of all lecturers

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)

**Response (200 OK):**
```json
{
  "status": "success",
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1,
    "data": {
      "lecturers": [
        {
          "id": "d1359b6d-4b28-4319-b09c-66377fe34a4b",
          "user_id": "8210511d-a5cb-4134-8d3d-04aa70d64e21",
          "lecturer_id": "LEC001",
          "department": "Computer Science",
          "user": {
            "id": "8210511d-a5cb-4134-8d3d-04aa70d64e21",
            "username": "lecturer001",
            "email": "lecturer001@university.ac.id",
            "full_name": "Dr. Jane Smith",
            "is_active": true
          },
          "created_at": "2025-12-13T10:12:22.617674+07:00"
        }
      ]
    }
  }
}
```

---

### 27. Get Lecturer's Advisees
**Endpoint:** `GET /lecturers/:id/advisees`  
**Access:** Protected (Requires `achievement:verify` permission)  
**Description:** Get all students advised by specific lecturer

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Advisees retrieved successfully",
  "data": {
    "lecturer": {
      "id": "d1359b6d-4b28-4319-b09c-66377fe34a4b",
      "lecturer_id": "LEC001",
      "full_name": "Dr. Jane Smith"
    },
    "advisees": [
      {
        "id": "student-uuid",
        "student_id": "STU001",
        "user": {
          "full_name": "John Doe",
          "email": "student001@university.ac.id"
        },
        "program_study": "Computer Science",
        "academic_year": "2024",
        "achievements_count": 3,
        "total_points": 400
      }
    ],
    "total_advisees": 1
  }
}
```

---

## Reports & Analytics

### 28. Get Statistics
**Endpoint:** `GET /reports/statistics`  
**Access:** Protected (Admin - `report:read` permission)  
**Description:** Get system-wide statistics

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Statistics retrieved successfully",
  "data": {
    "students": 10,
    "lecturers": 5,
    "achievements": {
      "draft": 5,
      "submitted": 8,
      "verified": 20,
      "rejected": 2,
      "total": 35
    },
    "achievement_types": {
      "competition": 12,
      "publication": 5,
      "organization": 8,
      "certification": 4,
      "academic": 4,
      "other": 2
    },
    "total_points": 5200
  }
}
```

---

### 29. Get Student Report
**Endpoint:** `GET /reports/student/:id`  
**Access:** Protected (Requires `report:read` permission)  
**Description:** Get detailed report for specific student

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Student report retrieved successfully",
  "data": {
    "student": {
      "id": "student-uuid",
      "student_id": "STU001",
      "full_name": "John Doe",
      "program_study": "Computer Science",
      "academic_year": "2024"
    },
    "advisor": {
      "full_name": "Dr. Jane Smith",
      "lecturer_id": "LEC001"
    },
    "achievements": {
      "total": 5,
      "verified": 3,
      "submitted": 1,
      "draft": 1
    },
    "achievements_by_type": {
      "competition": 2,
      "publication": 1,
      "organization": 1,
      "certification": 1
    },
    "total_points": 925,
    "recent_achievements": [
      {
        "id": "693d6046d6115cb20e966b7c",
        "achievement_type": "competition",
        "title": "International Programming Contest",
        "status": "verified",
        "points": 400,
        "created_at": "2025-12-13T12:47:02.473Z"
      }
    ]
  }
}
```

---

### 30. Get Top Students (Leaderboard) 🆕
**Endpoint:** `GET /reports/top-students`  
**Access:** Protected (Requires `report:read` permission)  
**Description:** Get leaderboard of top students based on verified achievements

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Query Parameters:**
- `limit` (optional): Number of top students to return (default: 10, max: 100)

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Top students retrieved successfully",
  "data": {
    "total": 5,
    "students": [
      {
        "student_id": "STU001",
        "full_name": "John Doe",
        "program_study": "Computer Science",
        "academic_year": "2024",
        "achievement_count": 15,
        "total_points": 3200,
        "verified_achievements": 15
      },
      {
        "student_id": "STU002",
        "full_name": "Jane Smith",
        "program_study": "Information Systems",
        "academic_year": "2024",
        "achievement_count": 12,
        "total_points": 2850,
        "verified_achievements": 12
      }
    ]
  }
}
```

---

### 31. Get Statistics by Period 🆕
**Endpoint:** `GET /reports/statistics/period`  
**Access:** Protected (Requires `report:read` permission)  
**Description:** Get achievement statistics grouped by month/year

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Query Parameters:**
- `start_date` (optional): Start date in YYYY-MM-DD format (default: 12 months ago)
- `end_date` (optional): End date in YYYY-MM-DD format (default: today)

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Period statistics retrieved successfully",
  "data": {
    "period": {
      "start": "2024-01-01",
      "end": "2025-12-31"
    },
    "statistics": {
      "2024-12": 15,
      "2025-01": 23,
      "2025-02": 18,
      "2025-03": 20
    },
    "total": 76
  }
}
```

---

### 32. Get Competition Level Distribution 🆕
**Endpoint:** `GET /reports/statistics/competition-levels`  
**Access:** Protected (Requires `report:read` permission)  
**Description:** Get distribution of achievements by competition level

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Competition level distribution retrieved successfully",
  "data": {
    "distribution": {
      "local": 45,
      "regional": 32,
      "national": 28,
      "international": 15
    },
    "total": 120,
    "by_type": {
      "competition": {
        "local": 25,
        "regional": 18,
        "national": 15,
        "international": 12
      },
      "publication": {
        "national": 8,
        "international": 3
      }
    }
  }
}
```

---

## Notifications

### 33. Get My Notifications 🆕
**Endpoint:** `GET /notifications`  
**Access:** Protected (Any authenticated user)  
**Description:** Get paginated list of user's notifications

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Notifications retrieved successfully",
  "data": {
    "total": 25,
    "unread_count": 5,
    "notifications": [
      {
        "id": "notification-uuid",
        "user_id": "user-uuid",
        "type": "achievement_verified",
        "title": "Achievement Verified",
        "message": "Congratulations! 'International Programming Contest' has been verified",
        "data": {
          "achievement_id": "693d6046d6115cb20e966b7c",
          "verified_by": "verifier-uuid"
        },
        "is_read": false,
        "read_at": null,
        "created_at": "2025-12-13T15:30:00Z"
      },
      {
        "id": "notification-uuid-2",
        "user_id": "user-uuid",
        "type": "achievement_submitted",
        "title": "New Achievement Submitted",
        "message": "John Doe submitted: Research Publication",
        "data": {
          "achievement_id": "693d6046d6115cb20e966b7d",
          "student_id": "student-uuid"
        },
        "is_read": true,
        "read_at": "2025-12-13T16:00:00Z",
        "created_at": "2025-12-13T14:00:00Z"
      }
    ]
  }
}
```

**Notification Types:**
- `achievement_submitted` - Student submitted achievement for verification
- `achievement_verified` - Achievement has been verified
- `achievement_rejected` - Achievement has been rejected
- `advisor_assigned` - Advisor has been assigned to student

---

### 34. Get Unread Notifications 🆕
**Endpoint:** `GET /notifications/unread`  
**Access:** Protected (Any authenticated user)  
**Description:** Get list of unread notifications

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Unread notifications retrieved successfully",
  "data": {
    "count": 3,
    "notifications": [
      {
        "id": "notification-uuid",
        "user_id": "user-uuid",
        "type": "achievement_verified",
        "title": "Achievement Verified",
        "message": "Congratulations! 'International Programming Contest' has been verified",
        "data": {
          "achievement_id": "693d6046d6115cb20e966b7c",
          "verified_by": "verifier-uuid"
        },
        "is_read": false,
        "read_at": null,
        "created_at": "2025-12-13T15:30:00Z"
      }
    ]
  }
}
```

---

### 35. Get Unread Count 🆕
**Endpoint:** `GET /notifications/unread/count`  
**Access:** Protected (Any authenticated user)  
**Description:** Get count of unread notifications

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Unread count retrieved successfully",
  "data": {
    "count": 5
  }
}
```

---

### 36. Mark Notification as Read 🆕
**Endpoint:** `PUT /notifications/:id/read`  
**Access:** Protected (Any authenticated user)  
**Description:** Mark a specific notification as read

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "Notification marked as read"
}
```

**Error Response (404):**
```json
{
  "status": "error",
  "error": "Notification not found"
}
```

---

### 37. Mark All Notifications as Read 🆕
**Endpoint:** `PUT /notifications/read-all`  
**Access:** Protected (Any authenticated user)  
**Description:** Mark all user's notifications as read

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "All notifications marked as read",
  "data": {
    "marked_count": 5
  }
}
```

---

## File Management

### 38. Upload File
**Endpoint:** `POST /files/upload`  
**Access:** Protected (Requires authentication)  
**Description:** Upload achievement attachment file

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Content-Type: multipart/form-data
```

**Form Data:**
- `file`: File to upload (max 10MB)
- Allowed types: PDF, JPG, JPEG, PNG, DOC, DOCX

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "File uploaded successfully",
  "data": {
    "filename": "1702456789123_certificate.pdf",
    "original_name": "certificate.pdf",
    "size": 245760,
    "url": "/uploads/1702456789123_certificate.pdf"
  }
}
```

---

### 39. Delete File
**Endpoint:** `DELETE /files/:filename`  
**Access:** Protected (Requires authentication)  
**Description:** Delete uploaded file

**Headers:**
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Response (200 OK):**
```json
{
  "status": "success",
  "message": "File deleted successfully"
}
```

---

## Test Credentials

Use these credentials to test the API:

| Role | Username | Password | Permissions |
|------|----------|----------|-------------|
| **Admin** | admin | admin123 | Full system access |
| **Student** | student001 | student123 | Create, read, update, delete own achievements |
| **Lecturer** | lecturer001 | lecturer123 | Read achievements, verify advisee achievements |

---

## Error Responses

All error responses follow this format:

**401 Unauthorized:**
```json
{
  "status": "error",
  "error": "Missing authorization header"
}
```

**403 Forbidden:**
```json
{
  "status": "error",
  "error": "You don't have permission to access this resource"
}
```

**404 Not Found:**
```json
{
  "status": "error",
  "error": "Resource not found"
}
```

**400 Bad Request:**
```json
{
  "status": "error",
  "error": "Invalid request body",
  "details": {
    "field": "achievement_type",
    "message": "achievement_type is required"
  }
}
```

**500 Internal Server Error:**
```json
{
  "status": "error",
  "error": "Internal server error"
}
```

---

## Notes

1. **Authentication:** All protected endpoints require Bearer token in Authorization header
2. **Pagination:** List endpoints support pagination with `page` and `limit` query parameters
3. **Filtering:** Many endpoints support filtering via query parameters
4. **Point Calculation:** Achievement points are automatically calculated based on type and details
5. **Status Workflow:** Achievements follow: draft → submitted → verified/rejected
6. **Soft Delete:** Delete operations use soft delete (status='deleted'), data is preserved
7. **RBAC:** Role-based access control enforces permissions on all endpoints
8. **Notifications:** 🆕 Automatically created when achievements are submitted, verified, or rejected
9. **Real-time Updates:** 🆕 Use GET /notifications/unread/count for polling notification updates
10. **Analytics:** 🆕 Advanced analytics available for admins and advisors to track student progress

---

## API Summary

**Total Endpoints:** 39

| Category | Endpoints | New Features |
|----------|-----------|--------------|
| Authentication | 4 | - |
| Users | 11 | - |
| Achievements | 10 | - |
| Students | 4 | - |
| Lecturers | 3 | - |
| Reports & Analytics | 6 | 3 🆕 (Top Students, Period Stats, Competition Levels) |
| Notifications | 5 | 5 🆕 (All endpoints) |
| Files | 2 | - |

**New in v1.1 (December 2025):**
- ✅ Notification system with 5 endpoints
- ✅ Advanced analytics with leaderboard and period analysis
- ✅ Competition level distribution reporting
- ✅ Auto-notification on achievement workflow events
- ✅ 100% functional requirements compliance

---

**For Swagger UI documentation, visit:** `http://localhost:3000/swagger/index.html`
