# pismo- transactions routines

This repository contains the backend APIs for the Transactions routines.
Components:
- Accounts
- Transactions
- Operation Types
- Users

## Table of Contents

- [Prerequisites](#prerequisites)
- [Build](#build)
- [Makefile](#makefile)
- [Local Development](#local-development)
- [Architecture Overview](#architecture-overview)
- [Current Limitations and Future Improvements](#current-limitations-and-future-improvements)
- [Production Considerations](#production-considerations)
- [API Endpoints](#api-endpoints)


### Prerequisites

- GoLang - 1.23+ - Follow the instructions here to [install Golang](https://go.dev/doc/install)

```shell
$ go version
go version go1.23.0 darwin/arm64
```

- Ensure that, `$GOPATH` is set and the `$GOPATH/bin` is added to `$PATH` in the shell config.

```shell
export GOPATH=$(go env GOPATH)
export PATH=$PATH:$(go env GOPATH)/bin
```

- Docker - To install visit [https://docs.docker.com/get-docker/](https://docs.docker.com/get-docker/)

### Build

- Clone the repository
- Copy `.env.default` to `.env` and modify as required.
- Build & Run

```shell
git clone https://github.com/imjenal/transaction-service.git
cd transaction-service
make dev
```

### Makefile

The Makefile contains a set of targets that can be used to build, test, and run the API. It also contains helper targets
to reduce the amount of typing required to run common commands.

It is recommended that you use the Makefile targets instead of running the commands directly. This will ensure that the
correct commands are issued with correct order and arguments.

To see the list of available targets, run ``make help``

### Local Development

Head over to [docs/LOCAL_DEV.md](docs/LOCAL_DEV.md) for detailed instructions to setup local development.

## Architecture Overview

Here's a brief overview:

- **Database (PostgreSQL)**: Serves as the primary data store for users, accounts, operation_types and transactions.
- **Go Application**:
    - Database Layer: Manages all interactions with the PostgreSQL database, including CRUD operations.
    - API Layer: Exposes RESTful endpoints to interact with the service.
    - Server: Handles HTTP requests and routes them to appropriate handlers.
  
## Current Limitations and Future Improvements
- Testing: Comprehensive unit and integration tests are essential for ensuring code reliability and ease of maintenance. Expanding the test suite would be a priority for future development.
- Logging and Monitoring: Implementing a more robust logging and monitoring system would be crucial for production deployment, allowing for better observability and troubleshooting.
- Security: Enhancements in security measures, such as securing API endpoints and database connections, would be necessary for a production environment.
- Deployment and Scalability: While the current setup is suitable for development and small-scale deployment, considerations for containerization (e.g., using Docker) and orchestration (e.g., Kubernetes) would be vital for larger-scale production deployment.

## Production Considerations
For production deployment, the following practices should be followed:
- Environment Configuration: Secure management of environment variables and configuration settings, potentially using a service like HashiCorp Vault.
- Database Scalability: Implementing database replication and sharding for improved performance and fault tolerance.
- Load Balancing and High Availability: Deploying the application across multiple servers or regions with load balancing to ensure high availability.
- CI/CD Pipelines: Establishing continuous integration and continuous deployment pipelines for streamlined development and deployment processes.

By addressing these areas, this project can be evolved into a robust, production-ready backend service capable of handling large-scale data and traffic.

## API Endpoints

This service exposes several RESTful endpoints for interacting with accounts and transactions. Below is a list of the available endpoints:

- **Create Account**:
    - `POST /api/v1/accounts`
    - creates an account

- **Fetch Account Details by AccountID**:
    - `GET /api/v1/accounts/{accountID}`
    - Retrieves details of a specific account.

- **Create Transactions**:
    - `POST /api/v1/transactions`
    - creates a transaction
  
- **Fetch Transaction Details by TransactionID**:
    - `GET /api/v1/transactions/{transactionID}`
    - Retrieves details of a specific transaction.

Our API implements versioning to ensure backward compatibility and a smooth transition for clients when introducing changes. The version of the API is specified in the URL, making it clear and easy to manage different versions of the API. The current version is v1.
