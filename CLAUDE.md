# Project Overview
This project is a finances tracker that allows me to upload my bank statements as a CSV file and track my expenses and income. This repository is for the API.

# Tech Stack
- Golang
- PostgreSQL

# Running application
- Run `go run cmd/api/main.go` to start the application.
- Run `go test ./...` to run the tests.

# Architecture
- Hexagonal Architecture but with a focus on simplicity and readability.
- I don't like directly naming folders ports and adapters
- I should be able to map from a domain model to a database model separately

# Code Conventions
- Use interfaces for dependency injection
- Use Go base interfaces as much as possible
- Start simple and build up
- Avoid comments, the code should be self-explanatory
- New code should ideally have tests or be written using TDD, mainly business logic

# Testing
- Many Unit Tests
- Couple of integration Tests
- Use Testify for assertions

# Code reviews
- This is a personal project so code does not have to be industry standard, but it should be readable and maintainable.
