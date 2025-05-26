# Chirpy: Learning HTTP Servers in Go

This project is a simple environment for learning how to build HTTP servers using the Go programming language. It serves static files and provides a basic health check endpoint.

## Features

- Serves static files (including HTML and assets) from the project directory.
- Provides a `/healthz` endpoint for readiness checks.
- Uses Go's standard `net/http` library.

## Project Structure

- `main.go`: Main application entry point, sets up the HTTP server.
- `index.html`: Example static HTML file.
- `assets/`: Directory for static assets (e.g., images).
- `go.mod`: Go module definition.

## Running the Server

1. Make sure you have Go installed (version 1.24.2 or later).
2. In the `app` directory, run:

   ```sh
   go run main.go
   ```

3. Open your browser and visit http://localhost:8080/app/index.html to see the welcome page.
4. Check the health endpoint at http://localhost:8080/healthz

## License

This project is for educational purposes.
