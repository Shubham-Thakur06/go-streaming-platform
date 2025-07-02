# Go Streaming Platform

A modern, scalable streaming platform built with Go, featuring user authentication, media upload/streaming, and multi-cloud storage support.

## 🚀 Features

- **User Authentication**: JWT-based authentication system with secure password hashing
- **Media Management**: Upload, stream, and manage audio/video content
- **Multi-Cloud Storage**: Support for AWS S3, Azure Blob Storage, and Google Cloud Storage
- **RESTful API**: Clean, well-structured REST API with proper error handling
- **Database Integration**: PostgreSQL with GORM for efficient data management
- **CORS Support**: Cross-origin resource sharing enabled for web applications
- **Activity Tracking**: User activity monitoring for analytics
- **Scalable Architecture**: Clean separation of concerns with modular design

## 🏗️ Architecture

The project follows a clean, layered architecture:

```
go-streaming-platform/
├── cmd/api/                 # Application entry point
├── internal/                # Private application code
│   ├── app/                 # Application initialization
│   ├── config/              # Configuration management
│   ├── database/            # Database connection and setup
│   ├── handlers/            # HTTP request handlers
│   ├── interfaces/          # Interface definitions
│   ├── middleware/          # HTTP middleware (auth, CORS)
│   ├── models/              # Data models and database schemas
│   ├── repository/          # Data access layer
│   ├── server/              # HTTP server setup and routing
│   ├── service/             # Business logic layer
│   └── storage/             # Cloud storage providers
├── pkg/                     # Public packages
│   ├── logger/              # Logging utilities
│   └── validator/           # Input validation
└── docs/                    # Documentation
```

## 🛠️ Tech Stack

- **Language**: Go 1.23.2
- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT
- **Storage**: AWS S3, Azure Blob Storage, Google Cloud Storage
- **Configuration**: Environment variables with godotenv
- **Validation**: Custom validators
- **Logging**: Structured logging

## 📋 Prerequisites

- Go 1.23.2 or higher
- PostgreSQL database
- Cloud storage account (AWS S3, Azure Blob Storage, or Google Cloud Storage)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/Shubham-Thakur06/go-streaming-platform.git
cd go-streaming-platform
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Environment Configuration

Create a `.env` file in the root directory:

```env
# Server Configuration
PORT=8080
HOST=localhost

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=streaming_platform
DB_SSLMODE=disable

# JWT Configuration
JWT_SECRET_KEY=your-secret-key-here
JWT_EXPIRY_HOURS=24

# Storage Configuration (choose one provider) aws | azure | gcp
STORAGE_PROVIDER=aws

# AWS Configuration
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_REGION=region-of-bucket
AWS_BUCKET_NAME=your-bucket-name

# Azure Configuration (if using Azure)
AZURE_ACCOUNT_NAME=your-account-name
AZURE_ACCOUNT_KEY=your-account-key
AZURE_CONTAINER_NAME=your-container-name

# Google Cloud Configuration (if using GCP)
GCP_PROJECT_ID=your-project-id
GCP_BUCKET_NAME=your-bucket-name
GCP_CREDENTIALS_FILE=path/to/credentials.json

# Host Configuration
HOST_USERNAME=host
HOST_PASSWORD=host123
HOST_EMAIL=host@streaming-platform.com
```

### 4. Database Setup

Create a PostgreSQL database and run the application. The database tables will be automatically created using GORM's auto-migration.

### 5. Run the Application

```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Authentication

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "host",
  "password": "host123"
}
```

### Media Management

#### List Media (Public)
```http
GET /api/v1/media
```

#### Get Media Details (Public)
```http
GET /api/v1/media/{id}
```

#### Get Stream URL (Public)
```http
GET /api/v1/media/{id}/stream
```

#### Upload Media (Protected - Host Only)
```http
POST /api/v1/media/upload
Authorization: Bearer {jwt_token}
Content-Type: multipart/form-data

{
  "file": "media_file",
  "title": "Media Title",
  "description": "Media Description",
  "genre": "Music",
  "tags": "tag1,tag2,tag3",
  "is_public": true
}
```

#### Delete Media (Protected - Host Only)
```http
DELETE /api/v1/media/{id}
Authorization: Bearer {jwt_token}
```

### User Management

#### Get Profile (Protected)
```http
GET /api/v1/profile
Authorization: Bearer {jwt_token}
```

#### Update Profile (Protected)
```http
PUT /api/v1/profile
Authorization: Bearer {jwt_token}
Content-Type: application/json

{
  "username": "new_username",
  "email": "new_email@example.com"
}
```

### Health Check

```http
GET /health
```

## 🗄️ Database Schema

### Users Table
- `id` (UUID, Primary Key)
- `username` (String, Unique)
- `email` (String, Unique)
- `password` (String, Hashed)
- `created_at` (Timestamp)
- `updated_at` (Timestamp)

### Media Table
- `id` (UUID, Primary Key)
- `title` (String)
- `description` (String)
- `filename` (String)
- `file_size` (Int64)
- `file_type` (String)
- `genre` (String)
- `tags` (String)
- `duration` (Int)
- `storage_url` (String)
- `thumbnail_url` (String)
- `user_id` (UUID, Foreign Key)
- `is_public` (Boolean)
- `view_count` (Int)
- `created_at` (Timestamp)
- `updated_at` (Timestamp)

### User Activity Table
- `id` (UUID, Primary Key)
- `user_id` (UUID, Foreign Key)
- `media_id` (UUID, Foreign Key)
- `action` (String)
- `duration` (Int)
- `created_at` (Timestamp)

## 🔧 Configuration

The application uses environment variables for configuration. Key configuration areas include:

- **Server**: Host and port settings
- **Database**: PostgreSQL connection parameters
- **Storage**: Cloud storage provider configuration
- **JWT**: Authentication token settings
- **Host**: Default host user credentials

## 🚀 Deployment

### Docker (Recommended)

1. Build the Docker image:
```bash
docker build -t go-streaming-platform .
```

2. Run the container:
```bash
docker run -p 8080:8080 --env-file .env go-streaming-platform
```

### Manual Deployment

1. Build the binary:
```bash
go build -o streaming-platform cmd/api/main.go
```

2. Run the binary:
```bash
./streaming-platform
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Shubham Thakur**
- GitHub: [@Shubham-Thakur06](https://github.com/Shubham-Thakur06)

## 🙏 Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library for Go
- [AWS SDK for Go](https://github.com/aws/aws-sdk-go) - AWS SDK
- [JWT Go](https://github.com/golang-jwt/jwt) - JWT implementation

## 📞 Support

If you have any questions or need help, please open an issue on GitHub or contact the maintainer.