# Feature Flags Management Service

[![Go](https://img.shields.io/badge/Go-1.24-blue)](https://golang.org)
[![MySQL](https://img.shields.io/badge/MySQL-8-orange)](https://www.mysql.com/)
[![Docker](https://img.shields.io/badge/Docker-Container-blue)](https://www.docker.com/)
[![Chi Router](https://img.shields.io/badge/Router-Chi-green)](https://github.com/go-chi/chi)

---

## Overview

This project is a backend service for managing **feature flags** with advanced support for dependencies and audit logging. Feature flags are a powerful technique to enable or disable functionality dynamically without deploying new code.

The service enables teams to:

- Create feature flags with optional dependencies on other flags.
- Toggle flags on/off, enforcing that all dependencies are active before enabling.
- Automatically disable dependent flags if a parent flag is disabled.
- Track all changes with a detailed audit log including actor, reason, and timestamp.

---

## Architecture & Tech Stack

- **Language:** Go (Golang 1.24.2)
- **HTTP Router:** [Chi](https://github.com/go-chi/chi) for lightweight, idiomatic routing and middleware.
- **Database:** MySQL 8 for durable relational storage.
- **Containerization:** Docker & Docker Compose for easy local setup and deployment.
- **Dependency management:** Go modules.
- **Database driver:** `github.com/go-sql-driver/mysql`

---

## Features

### Feature Flags
- Unique flags identified by name.
- Enabled or disabled states.
- Timestamps for creation and updates.

### Dependency Management
- Flags can depend on multiple other flags.
- Circular dependencies are detected and prevented.
- Flags cannot be enabled unless all dependencies are enabled.
- Disabling a flag automatically disables dependent flags recursively.

### Audit Logging
- Logs all actions (create, toggle, auto-disable) with timestamps.
- Stores actor information and reason for change.

### RESTful API Endpoints

| Method | Endpoint               | Description                          |
|--------|------------------------|------------------------------------|
| GET    | `/health`              | Health check endpoint               |
| POST   | `/flags`               | Create a new feature flag           |
| POST   | `/flags/{name}/toggle` | Enable or disable a feature flag    |
| GET    | `/flags/{name}`        | Retrieve details of a flag          |
| GET    | `/flags/{name}/logs`   | Retrieve audit logs for a flag      |

---

## Getting Started

### Prerequisites

- [Go 1.24+](https://golang.org/dl/)
- [Docker & Docker Compose](https://docs.docker.com/compose/install/)
- MySQL (via Docker or native)

---

### Environment Variables

Create a `.env` file with the following variables:

```env
DB_DSN=root:password@tcp(db:3306)/feature_flags?parseTime=true
PORT=8080