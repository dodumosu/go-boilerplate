# Go Boilerplate

This project is a comprehensive boilerplate for building production-ready web applications in Go. It comes with a variety of features out of the box, allowing you to focus on your application's core logic.

## Features

*   **Authentication:** Secure user authentication with JWTs, including sign-up, login, and verification.
*   **Background Jobs:** Asynchronous job processing using Asynq and Redis, with examples for sending verification emails.
*   **Configuration Management:** Centralized configuration management using environment variables, with a clear structure in `internal/config`.
*   **Database Management:** Database migrations with `dbmate` and type-safe queries with `sqlc`.
*   **Structured Logging:** Structured logging with `slog` for clear and actionable logs.
*   **OAuth2 Support:** (Coming Soon) Support for OAuth2 providers like Google and Facebook.
*   **RESTful API:** A well-structured RESTful API with clear separation of concerns. API documentation is using [Huma](https://huma.rocks)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

*   [Go](https://golang.org/)
*   [Docker](https://www.docker.com/)
*   [Make](https://www.gnu.org/software/make/)
*   [dbmate](https://github.com/amacneil/dbmate)
*   [sqlc](https://github.com/sqlc-dev/sqlc)

### Installation

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/go-boilerplate.git
    cd go-boilerplate
    ```

2.  **Install dependencies:**

    ```bash
    go mod download
    ```

### Configuration

1.  **Create a `.env` file:**

    Copy the `.env.template` file to a new file named `.env`:

    ```bash
    cp .env.template .env
    ```

2.  **Update the environment variables:**

    Open the `.env` file and update the variables to match your local environment. At a minimum, you'll need to configure your database and email server settings.

### Running the Application

1.  **Start the database and other services:**

    This project uses Docker Compose to manage the required services (PostgreSQL and Redis). To start them, run:

    ```bash
    docker-compose up -d
    ```

2.  **Run the database migrations:**

    ```bash
    make migrate-up
    ```

3.  **Run the application:**

    You can run the application using the `make dev` command, which will start both the web server and the front end development server:

    ```bash
    make dev
    ```

    Alternatively, you can run the server by itself:

    *   **Run the server:**

        ```bash
        make run-server
        ```

    **Note:** The server includes a worker processing background jobs.

## Usage

The API documentation is available at `/api/docs` when the server is running.

## Database Migrations

This project uses `dbmate` to manage database migrations. You can create new migrations and run them using the following commands:

*   **Create a new migration:**

    ```bash
    make migrate-create
    ```

*   **Run all pending migrations:**

    ```bash
    make migrate-up
    ```

*   **Rollback the last migration:**

    ```bash
    make migrate-down
    ```

## Testing

(Coming Soon) Instructions on how to run the test suite will be added soon.

## License

This project is licensed under the Unlicense - see the [LICENSE.md](LICENSE.md) file for details.
