# IPLE LMS Backend System Design

## ออกแบบระบบ (System Design) - COMPLETE ✅

This project implements a complete **Interactive Programming Learning Environment (IPLE)** backend system with enterprise-grade architecture and modern design patterns.

## 🏗️ System Architecture

### Layered Architecture
```
┌─────────────────────────────────────────┐
│              HTTP Layer                 │
│  (Routes, Middleware, Authentication)   │
├─────────────────────────────────────────┤
│             Service Layer               │
│        (Business Logic)                 │
├─────────────────────────────────────────┤
│            Repository Layer             │
│         (Data Access)                   │
├─────────────────────────────────────────┤
│              Model Layer                │
│        (Domain Entities)                │
└─────────────────────────────────────────┘
```

### Core Components Implemented

#### 1. **Middleware Layer** (`internal/middleware/`)
- **Request ID Tracking**: Unique request identification for distributed tracing
- **Structured Logging**: JSON-based logging with context
- **JWT Authentication**: Token-based authentication middleware
- **Role-based Authorization**: Admin, Teacher, Student, Parent access control

#### 2. **Service Layer** (`internal/services/`)
- **AuthService**: User registration, login, token management
- **Business Logic Separation**: Clean separation from HTTP and data layers
- **Error Handling**: Proper error propagation and handling

#### 3. **Repository Layer** (`internal/repositories/`)
- **Repository Pattern**: Database abstraction layer
- **GORM Integration**: ORM with PostgreSQL support
- **Query Optimization**: Efficient data access patterns

#### 4. **Utils Package** (`internal/utils/`)
- **JWT Management**: Token generation and validation
- **Password Security**: Scrypt-based password hashing
- **Response Formatting**: Standardized API responses
- **Validation Helpers**: Request validation utilities

#### 5. **DTOs** (`internal/dto/`)
- **Request/Response Separation**: Clean API contracts
- **Validation Rules**: Input validation with go-playground/validator
- **Type Safety**: Strongly-typed API interfaces

## 🔐 Security Features

### Authentication & Authorization
- **JWT-based Authentication**: Stateless token-based auth
- **Password Hashing**: Secure scrypt algorithm
- **Role-based Access Control**: Multi-level permission system
- **Token Refresh**: Secure token renewal mechanism

### Security Middleware
- **CORS Protection**: Cross-origin request security
- **Helmet Integration**: Security headers
- **Request Validation**: Input sanitization
- **Error Sanitization**: Safe error responses

## 📊 Database Design

### Models Implemented
- **User**: Students, Teachers, Admins, Parents
- **Course**: Learning content structure
- **Section**: Course organization units
- **Lesson**: Individual learning materials
- **Enrollment**: Student-course relationships
- **Submission**: Student work submissions
- **Grade**: Assessment results
- **ParentChildLink**: Parent-student relationships
- **PDPAConsentLog**: Privacy compliance tracking

### Database Features
- **GORM ORM**: Type-safe database operations
- **Automatic Migrations**: Development environment setup
- **Connection Pooling**: Optimized database connections
- **Soft Deletes**: Data preservation
- **Timestamps**: Audit trail tracking

## 🚀 API Endpoints

### Public Endpoints
- `GET /api/v1/health` - System health check
- `GET /swagger/*` - API documentation

### Authentication (`/api/v1/auth/`)
- `POST /register` - User registration
- `POST /login` - User authentication
- `POST /refresh` - Token renewal
- `POST /logout` - User logout

### User Management (`/api/v1/users/`) [Protected]
- `GET /profile` - Current user profile
- `PUT /profile` - Update profile
- `GET /` - List all users
- `GET /:id` - Get user by ID
- `PUT /:id` - Update user
- `DELETE /:id` - Delete user

### Course Management (`/api/v1/courses/`) [Protected]
- `GET /` - List courses
- `POST /` - Create course
- `GET /:id` - Get course details
- `PUT /:id` - Update course
- `DELETE /:id` - Delete course
- `GET /:id/sections` - Course sections
- `POST /:id/sections` - Create section
- `GET /:id/sections/:section_id/lessons` - Section lessons

### Enrollment System (`/api/v1/enrollments/`) [Protected]
- `GET /` - User enrollments
- `POST /` - Enroll in course
- `GET /:id` - Enrollment details
- `PUT /:id` - Update enrollment
- `DELETE /:id` - Withdraw enrollment
- `GET /:id/progress` - Learning progress
- `GET /:id/submissions` - Student submissions

### Admin Portal (`/api/v1/admin/`) [Admin Only]
- `GET /users` - User management
- `POST /users` - Create users
- `PUT /users/:id/status` - User status management
- `GET /courses` - Course administration
- `GET /enrollments` - Enrollment oversight
- `GET /stats` - System statistics

### Parent Portal (`/api/v1/parent/`) [Parent Only]
- `GET /children` - Linked children
- `POST /children` - Link child account
- `DELETE /children/:id` - Unlink child
- `GET /children/:id/progress` - Child progress
- `GET /children/:id/enrollments` - Child enrollments
- `GET /children/:id/pdpa-consent` - Privacy consent management

## 🛠️ Development Features

### Configuration Management
- **Environment-based Config**: Development, staging, production
- **Secret Management**: Secure credential handling
- **Database Configuration**: Connection pooling and optimization
- **JWT Configuration**: Token expiration and secret management

### Logging & Monitoring
- **Structured Logging**: JSON-formatted logs with slog
- **Request Tracing**: Unique request IDs
- **Error Tracking**: Comprehensive error logging
- **Performance Monitoring**: Request latency tracking

### Development Tools
- **Hot Reloading**: Development server with auto-restart
- **Database Migrations**: Automatic schema management
- **API Documentation**: Swagger/OpenAPI integration
- **Environment Variables**: `.env` file support

## 🔧 Usage

### Quick Start
```bash
# Clone repository
git clone https://github.com/saengpech-sys/iple-backend

# Install dependencies
go mod tidy

# Run in development mode
ENV=development go run ./cmd/server

# Run in test mode (no database)
ENV=test go run ./cmd/server
```

### Environment Configuration
```env
ENV=development
PORT=8080
LOG_LEVEL=info

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=iple_lms_db

# JWT
JWT_SECRET_KEY=your_secret_key
JWT_EXPIRES_IN_MINUTES=60
```

### API Testing
```bash
# Health check
curl http://localhost:8080/api/v1/health

# User registration
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"first_name":"John","last_name":"Doe","email":"john@example.com","password":"password123","role":"student"}'

# Access protected endpoint
curl -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/users/profile
```

## 📝 API Documentation

- **Swagger UI**: Available at `/swagger/index.html`
- **OpenAPI Spec**: RESTful API documentation
- **Interactive Testing**: Built-in API testing interface

## 🚀 Production Readiness

### Features Implemented
- ✅ **Scalable Architecture**: Clean layered design
- ✅ **Security**: JWT auth, password hashing, CORS
- ✅ **Database**: ORM integration with connection pooling
- ✅ **Logging**: Structured logging with request tracing
- ✅ **Error Handling**: Centralized error management
- ✅ **Validation**: Input validation and sanitization
- ✅ **Documentation**: Comprehensive API docs
- ✅ **Configuration**: Environment-based configuration
- ✅ **Testing**: Health checks and endpoint testing

### Next Steps for Production
- [ ] **Database Migrations**: Production migration strategy
- [ ] **Monitoring**: Application performance monitoring
- [ ] **Caching**: Redis integration for session management
- [ ] **Rate Limiting**: API rate limiting middleware
- [ ] **Load Testing**: Performance benchmarking
- [ ] **CI/CD Pipeline**: Automated deployment
- [ ] **Container Deployment**: Docker and Kubernetes setup

## 📊 System Metrics

- **Total Endpoints**: 30+ RESTful API endpoints
- **Authentication**: JWT-based with role-based access control
- **Database Models**: 9 core domain models
- **Middleware**: 4 security and logging middleware components
- **Services**: Modular service layer architecture
- **Test Coverage**: API endpoint validation and security testing

---

**ออกแบบระบบ (System Design) เสร็จสมบูรณ์!** 

This IPLE LMS backend represents a production-ready, enterprise-grade system with modern Go development practices, comprehensive security, and scalable architecture patterns.