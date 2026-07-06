# Customer Management System

> A scalable, production-ready Customer Management System built with Go, utilizing clean architecture and containerized for seamless deployment.

## 🚀 Overview

This repository houses a robust Customer Management backend service. Designed with modularity in mind, it leverages Go's performance, `sqlc` for type-safe database interactions, and Docker to provide a seamless developer experience from local testing to production deployments.

## 🛠 Tech Stack

* **Language:** Go
* **Database Tooling:** [sqlc](https://sqlc.dev/) (Type-safe SQL)
* **Containerization:** Docker & Docker Compose
* **API Client/Testing:** [Bruno](https://www.usebruno.com/) (Collections included)

## 📂 Project Structure

```text
.
├── cmd/server/                    # Main application entrypoint
├── internal/                      # Private core business logic and domain models
├── 0_document/                    # Architecture documentation and specs
├── Customer System API Endpoint/  # Bruno API request collections
├── sqlc.yaml                      # SQLC configuration for database code generation
├── Dockerfile                     # Instructions to build the Go binary container
└── docker-compose.yml             # Local multi-container orchestration

```

## ⚡ Getting Started

### Prerequisites

Ensure you have the following installed on your local development environment:

* [Go](https://golang.org/dl/) (1.20+)
* [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
* [sqlc](https://docs.sqlc.dev/en/latest/overview/install.html) *(Optional, for regenerating DB code)*
* [Bruno](https://www.usebruno.com/) *(Optional, for interacting with the API)*

### Installation & Run

1. **Clone the repository:**
```bash
# Clone the project repository to your local machine
git clone https://github.com/Amir-Golmoradi/Customer-Management-System.git
cd Customer-Management-System

```


2. **Run with Docker Compose (Recommended):**
The most straightforward way to run the service along with its database dependencies.
```bash
# Build the images and spin up the containers in detached mode
docker-compose up --build -d

```


3. **Run Locally (Standalone):**
```bash
# Download the required Go modules
go mod download

# Start the HTTP server
go run cmd/server/main.go

```



## 📖 API Documentation

API request collections are provided natively in the repository via Bruno (`.bru` files). To explore and test the available endpoints:

1. Open the Bruno app.
2. Load the `Customer System API Endpoint` directory as a collection.
3. Configure your environment variables to map to your local instance (e.g., `localhost:8080`).

## 🤝 Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

Distributed under the MIT License.
