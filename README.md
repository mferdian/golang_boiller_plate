# Go Boilerplate

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)
![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)

A production-ready, scalable Golang boilerplate featuring clean architecture, JWT authentication, and modern development practices. Built to accelerate your backend development with industry best practices.

## Features

- **Clean Architecture** - Service & Repository pattern for maintainable code
- **JWT Authentication** - Secure access & refresh token implementation
- **Structured Logging** - Advanced logging with multiple levels and formats
- **Configuration Management** - TOML-based config with environment variables
- **Live Reload** - Hot reload during development with Air
- **Database Agnostic** - Easy database switching (PostgreSQL default)
- **Middleware Support** - CORS, rate limiting, request logging
- **Error Handling** - Centralized error handling and validation

## Architecture

This boilerplate follows the **Clean Architecture** principles:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│    Handler      │───▶│     Service     │───▶│   Repository    │
│   (HTTP Layer)  │    │ (Business Logic)│    │  (Data Layer)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```


## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL (or your preferred database)

### 1. Clone and Setup

```bash
git clone https://github.com/yourusername/go-boilerplate.git
cd go-boilerplate

# Copy environment variables
cp .env.example .env

# Install dependencies
go mod tidy
```

### 2. Configure Environment

Edit `.env` file:

```env
# Database configuration
DB_HOST=localhost
DB_USER=your_db_user
DB_PASS=your_db_password
DB_NAME=your_db_name
DB_PORT=your_db_port

# Application environment
APP_ENV=development

# App ports
PORT=8000
NGINX_PORT=8080
GOLANG_PORT=8888

# JWT secret key (use a long random string)
JWT_SECRET=your_jwt_secret_key

# SMTP configuration
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_SENDER_NAME="Your App Name <no-reply@example.com>"
SMTP_AUTH_EMAIL=your_email@example.com
SMTP_AUTH_PASSWORD=your_email_password


```

### 3. Database Setup

```bash
# Run migrations
./scripts/migrate.sh up

# Or manually create database and run migrations
createdb your_database
psql -d your_database -f migrations/001_create_users_table.sql
```

### 4. Run the Application

#### Development with Live Reload

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

#### Manual Run

```bash
go run main.go
```

## Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feat/amazing-feature`
3. Commit your changes: `git commit -m 'feat: add amazing feature'`
4. Push to the branch: `git push origin feat/amazing-feature`
5. Open a Pull Request

Please read our [Contributing Guidelines](CONTRIBUTING.md) for details on our code of conduct and development process.


## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library
- [Air](https://github.com/cosmtrek/air) - Live reloading

## Inspiration & Credits

### Primary References
This boilerplate was heavily inspired by:
- **[go-boilerplate](https://github.com/Amierza/go-boiler-plate)** - Project structure & Clean architecture patterns

### Mentors 
Special thanks to:
- [Ahmad Mirza Rafiq Azmi](https://github.com/Amierza) for guidance on Go best practices


**Note:** While this project draws inspiration from various sources, all code has been written from scratch or significantly modified to fit our specific use case.

## Support

- Email: eikhapoetra@gmail.com
- Issues: [GitHub Issues](https://github.com/mferdian/golang_boiller_plate/issues)
- Discussions: [GitHub Discussions](https://github.com/mferdian/go-golang_boiller_plate/discussions)

---

<div align="center">

**[Star this repository](https://github.com/mferdian/golang_boiller_plate)** if you find it helpful!

Made by Maulana (https://github.com/mferdian)

</div>
