
# Ogigo API

## Overview

The Ogigo API is designed to handle wallet updates by processing events and placing them into a RedPanda queue. It also provides an endpoint to retrieve the current wallet state from a PostgreSQL database.

## Architecture

The system architecture includes the following components:

-   **Gin**: A web framework for building the API.
-   **Kafka**: Used for queuing wallet update events.
-   **PostgreSQL**: Database for storing and retrieving wallet data.
-   **Swagger**: API documentation generator.
-   **Docker**: Containerization for easy deployment.

  **Clone the repository**:
    
    `git clone https://github.com/yourusername/ogigo.git
    cd ogigo` 
    
   **Build and run Docker containers**:
    
    `docker-compose up --build` 
    
   **Access Swagger Documentation**: Open your browser and navigate to `http://localhost:80/swagger/index.html` to view the API swagger.

   ## Usage
   
### Endpoints

-   **POST /**: Process incoming wallet update events.
-   **GET /**: Retrieve the current wallet state.

#### Run Tests

`go test ./...`