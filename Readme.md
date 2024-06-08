# SagMi | API Health Checker

SagMi is a service that monitors the health of multiple API endpoints and provides visualization and notification capabilities for users. It is developed using the following tech stack:

- Backend: Golang
- Web Server: Gin
- Frontend: htmx
- Database: SQLite
- Containerization: Docker
- Configuration: TOML

## Features

- Background service to automatically check the health of API endpoints and save the logs in the database.
- Capability to print logs in the console.
- Sending Slack messages and email using Mailgun.
- Web UI for monitoring the health status of endpoints visually.
- Subscribing new endpoints, deleting/editing existing ones, and unsubscribing endpoints.
- CLI manual endpoint health check.

## Installation

To install and run the API Health Checker service, follow these steps:

1. Clone the repository: `git clone https://github.com/your-username/api-health-checker.git`
2. Build the Docker image: `docker build -t api-health-checker .`
3. Create a `config.toml` file with the necessary configuration options.
4. Run the Docker container: `docker run -d -v /path/to/config.toml:/app/config.toml --name api-health-checker api-health-checker`

## Configuration

The API Health Checker service uses a `config.toml` file for configuration. The available options are:

- `slack_webhook_url`: The URL of the Slack webhook for sending messages.
- `mailgun_api_key`: The API key for Mailgun.
- `mailgun_domain`: The Mailgun domain.
- `mailgun_sender`: The email address of the sender.
- `mailgun_recipient`: The email address of the recipient.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
