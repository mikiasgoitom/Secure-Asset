# Tools and Packages for SecureAsset

This document outlines the recommended Go packages and external tools for building the SecureAsset management system.

## Go Packages

Here is a list of Go packages that will help you implement the required security features. You can install them all with the following command:

```bash
go get github.com/gin-gonic/gin \
       github.com/spf13/viper \
       go.mongodb.org/mongo-driver/mongo \
       golang.org/x/crypto/bcrypt \
       github.com/golang-jwt/jwt/v4 \
       github.com/pquerna/otp/totp \
       github.com/dchest/captcha \
       github.com/casbin/casbin/v2 \
       github.com/casbin/mongo-adapter/v3 \
       github.com/uber-go/zap \
       github.com/go-playground/validator/v10
```

### Package Breakdown

| Category                       | Package                                  | Purpose                                                                                                                       |
| ------------------------------ | ---------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------- |
| **Web Framework**              | `github.com/gin-gonic/gin`               | The core web framework for building RESTful APIs.                                                                             |
| **Configuration**              | `github.com/spf13/viper`                 | For managing application configuration from files (e.g., `config.yaml`), environment variables, and command-line arguments.   |
| **Database**                   | `go.mongodb.org/mongo-driver/mongo`      | The official MongoDB driver for Go to interact with your database.                                                            |
| **Password Hashing**           | `golang.org/x/crypto/bcrypt`             | A secure, standard library for hashing and verifying passwords. It includes salting to protect against rainbow table attacks. |
| **Token Authentication (JWT)** | `github.com/golang-jwt/jwt/v4`           | For creating and validating JSON Web Tokens (JWTs) to manage user sessions.                                                   |
| **Multi-Factor Auth (MFA)**    | `github.com/pquerna/otp/totp`            | To generate and validate Time-based One-Time Passwords (TOTP) for MFA.                                                        |
| **CAPTCHA**                    | `github.com/dchest/captcha`              | A simple library to generate and verify image or audio CAPTCHAs to prevent bot-driven account creation.                       |
| **Access Control**             | `github.com/casbin/casbin/v2`            | A powerful authorization library that supports multiple access control models like MAC, DAC, RBAC, ABAC, and RuBAC.           |
| **Casbin DB Adapter**          | `github.com/casbin/mongo-adapter/v3`     | The MongoDB adapter for Casbin, allowing you to store your access control policies in the database.                           |
| **Logging**                    | `github.com/uber-go/zap`                 | A high-performance, structured logging library. It can be configured for centralized and encrypted logging.                   |
| **Request Validation**         | `github.com/go-playground/validator/v10` | For validating incoming request data (structs) based on tags, essential for ensuring data integrity.                          |

## Other Recommended Tools

| Tool                | Purpose                                                                                                                                                   |
| ------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Postman**         | Essential for testing your RESTful APIs. You can create collections for all your endpoints to test authentication, authorization, and business logic.     |
| **MongoDB Compass** | A GUI for MongoDB that allows you to visualize, query, and manage your database, which is very helpful during development.                                |
| **Docker**          | For containerizing your Go application and MongoDB instance. This simplifies setup, ensures consistency across environments, and makes deployment easier. |
| **Git**             | For version control. Essential for tracking changes, collaborating, and managing your codebase.                                                           |
| **Go Linter**       | Tools like `golangci-lint` help enforce code quality and catch common mistakes before they become bugs.                                                   |
