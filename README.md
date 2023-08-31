# Avito test task
executed by Artem Shmakov
## Installation & Run
```bash
# Download this project
git clone ...
```
You can either use the docker to get the application up or run it on a local machine.
# Docker-compose
```bash
# Build and Run
cd app
docker-compose up --build app
```
# Local machine
```bash
#build and run
cd app
go run cmd/main/main.go
```
you can change the startup parameters in the files `configs/config.yml` or `configs/dockerConfig.yml`

## Structure
```
├── app
│   ├── cmd
|   |   └── main
|   |        └── main.go
│   ├── internal               // Our API core handlers
│   │   ├── config             // Common response functions
|   |        └── config.go
│   │   └── handlers           // Our API core handlers
|   |        └── handler.go
│   └── configs                // Configuration
│   |    ├── config.yml        // For local machine
|   |    └── dockerConfig.yml  // For docker-compose
|   └── pkg
|        ├── client
|        |    ├── models
|        |    |    └── model.go // Models for our application
|        |    └── postgres
|        |        └── db.go     // Сode of initialization and connection to our database
|        └── utils              // Helper code
|            └── helpers.go
└── Dockerfile
└── docker-compose.yml
```

## API

#### /segment
* `POST` : create a new segment
* `DELETE` : delete a segment

#### /segment/{id}
* `GET` : get all user segments by user_id
#### /user/{id}
* `POST` : creates a new user if a user with this id does not exist yet. Adds a list of segments to the user and deletes the list of segments.
