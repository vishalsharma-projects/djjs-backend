# DJJS Event Reporting Backend

This is a Go backend application for event reporting. It includes a **login/logout API** with JWT authentication and role-based access. The backend connects to a **Postgres database** (running via Docker) and uses **GORM** for ORM.

Other APIs will be added later. 

---

## **Prerequisites**

Before starting, make sure you have the following installed:

- [Go 1.21+](https://go.dev/doc/install)
- [Docker & Docker Compose](https://docs.docker.com/compose/install/)
- Git


## **Clone the Repository**

git clone https://github.com/yourusername/djjs-event-reporting-backend.git

cd djjs-event-reporting-backend

## **Start the Database**

The project uses Docker Compose to start a Postgres container with the database and initial tables.

docker-compose up -d

## **Install Go Dependencies**

go mod tidy

## **Run the Backend**

go run main.go

## **Access the APIs**

Once the server is running:

Base URL: http://localhost:8080
Swagger UI: http://localhost:8080/swagger/index.html
Swagger JSON (raw): http://localhost:8080/swagger/doc.json

## **Generating Swagger Docs (if needed)**

If you make changes to your routes or handlers, regenerate Swagger docs using the following command:

swag init -g main.go -o docs

## **Test the APIs**

Login Request (run in terminal)

Linux:

curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}'

Windows:

Invoke-RestMethod -Uri http://localhost:8080/login `
 -Method POST `
 -Headers @{ "Content-Type" = "application/json" } `
 -Body '{"email":"admin@example.com","password":"admin123"}'

Response:
{
  "token": "<JWT_TOKEN>"
}

Logout Request (run in terminal)

Use the token from login in the Authorization header:

Linux:

curl -X POST http://localhost:8080/logout \
  -H "Authorization: Bearer <JWT_TOKEN>"

Windows:

 Invoke-RestMethod -Uri http://localhost:8080/logout `
  -Method POST `
  -Headers @{ "Authorization" = "Bearer $token" }