# Community Library Management System (CLMS) API

A robust, production-ready RESTful API developed for the **Belize National Library Service and Information System (BNLSIS)**. This system manages book collections, member registrations, and circulation logic with integrated security and performance optimizations.

## 🚀 Core Features
* **Full CRUD Operations**: Comprehensive management for Books, Members, Loans, and Fines.
* **Belizean Business Logic**: 
    * Standardized 14-day loan periods.
    * Fixed $3.00 BZD local membership registration fee.
* **Advanced Catalog Search**: Filter books dynamically by **Title** or **ISBN**.
* **Database Integrity**: Normalized PostgreSQL schema with strict constraints and relationship mapping.

---

## 🛠️ Technical Stack
* **Language**: Go (v1.25+)
* **Database**: PostgreSQL
* **Features**: Middleware-based Rate Limiting, CORS, Gzip Compression, and Real-time Metrics.

---

## 🚦 Getting Started

### 1. Prerequisites
Ensure you have **Go** and **PostgreSQL** installed and running on your system.

### 2. Database Setup
Initialize the schema and inject the Belizean library seed data:
```bash
make db/seed
```

### 3. Run the Application
Start the API server (defaults to port 4000):
```bash
make run/api
```

---

## 📡 API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/v1/healthcheck` | Check system status and environment. |
| `GET` | `/v1/books` | List books (Supports `?title=` and `?isbn=`). |
| `POST` | `/v1/books` | Add a new book to the collection. |
| `POST` | `/v1/members` | Register a new library member. |
| `POST` | `/v1/loans` | Process a book checkout (Circulation). |
| `POST` | `/v1/fines` | Record a membership fee or overdue fine. |
| `GET` | `/debug/vars` | View live system & business metrics. |

---

## 🛡️ Infrastructure & Security Proofs
This project implements several advanced requirements for modern web APIs:

* **Rate Limiting**: Protects against brute-force/DOS attacks (Configured for 2 RPS with a burst of 5).
* **CORS**: Configured to allow secure communication with trusted frontend origins.
* **Gzip Compression**: Custom middleware reduces response size for large JSON payloads (e.g., Metrics).
* **Dynamic Metrics**: Real-time tracking of:
    * `total_books_loaned`
    * `total_members_registered`
    * `total_fines_created`

---

## 🧪 Automated Auditing
To verify the system's compliance, two automated test suites are included:

**1. Functional API Audit** (Verifies database logic and validators):
```bash
make test/api
```

**2. Infrastructure Audit** (Verifies CORS, Rate Limiting, and Gzip):
```bash
make test/network
```