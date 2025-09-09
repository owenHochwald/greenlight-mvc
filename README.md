# Greenlight Movie API

A comprehensive, production-ready RESTful API for movie data built with Go, implementing modern backend development practices and architectures.

## Project Overview

Greenlight is a robust movie API that demonstrates professional backend development patterns and practices. The project showcases a well-structured Go application that follows SOLID principles, dependency injection, and clean architecture.

## Key Features

- **RESTful API**: Full CRUD operations for movie resources
- **Rate Limiting**: 
  - Global rate limiting to protect API resources
  - IP-based rate limiting for fair usage
- **Middleware Pattern**: Modular request processing with middleware chains
- **Comprehensive Validation**: Request validation for data integrity
- **Advanced Pagination**: Efficient data retrieval with pagination support
- **Flexible Sorting**: Dynamic sorting options for API responses
- **Structured Logging**: Detailed logging for monitoring and debugging
- **Error Handling**: Consistent error responses across the API
- **MVC Architecture**: Clean separation of concerns
- **Dependency Injection**: Loosely coupled components for better testability
- **Interface Composition**: Go-idiomatic design using interfaces and receivers

## API Endpoints

The API offers the following endpoints:

- `GET /v1/healthcheck` - Service health status
- `GET /v1/movies` - List all movies with pagination and sorting
- `POST /v1/movies` - Create a new movie entry
- `GET /v1/movies/:id` - Retrieve a specific movie
- `PATCH /v1/movies/:id` - Update a movie's information
- `DELETE /v1/movies/:id` - Remove a movie from the database

## Technical Implementation

- **Docker Integration**: Containerized PostgreSQL database
- **Database Migrations**: Up and down migrations for database versioning
- **Compiled Binary**: Optimized for deployment on EC2 instances
- **Go Modules**: Modern dependency management
- **Gin Web Framework**: High-performance HTTP routing

## Deployment

The application is designed for cloud deployment with:
- Compiled binary for efficient execution on EC2 instances
- Containerized database for consistent development and production environments
- Structured logging for operational monitoring

## Future Roadmap

- User authentication and authorization with JWT
- User registration with email verification
- Account activation workflows
- Stateless authentication
- Real-time metrics and monitoring
- Performance optimization

---

