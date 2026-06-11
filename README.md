# Full-Stack Task Management Application

A secure, performance-focused task management system engineered with a Go REST API backend and a Next.js frontend, backed by an isolated PostgreSQL Docker instance.

## 🚀 Architectural Decisions & Assumptions
- **State Management Persistence:** Implemented clean functional hooks initialization for React Context states to extract values safely from `localStorage` on initial mount, preventing server/client hydration layout flickering and double-rendering sync cycles.
- **Data Isolation Rules:** Tightened security scopes across all single-resource API route handlers (`GET`, `PATCH`, `DELETE`). Queries natively intersect specific database task indices using both the `id` and the context-extracted authenticated `userID` to eliminate cross-tenant data modification vulnerabilities.
- **Route Trailing-Slash Consistency:** Standardized routing blocks in Go to use explicit full paths rather than nested prefix grouping configurations, eliminating HTTP path matching inconsistencies.

## 🛠️ Environment Prerequisites
Make sure your system has the following runtimes installed:
- Docker & Docker Compose
- Go (Version 1.20+)
- Node.js (Version 18+)

---

## ⚡ Setup & Local Run Instructions

### 1. Database Layer Engine Activation
From the root repository path, spin up your PostgreSQL engine background daemon:
```bash
docker compose up -d