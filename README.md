# Project decisions

- **Standard Library:** do without known web frameworks like Echo or Gin
in order to showcase knowledge over the language's standard library;
- **Provisioning:** use Docker/Docker Compose to facilitate the project's
portability through virtualization;

# Project dependencies

- This project was build using **go1.21.0 linux/amd64**;
- The laguage's dependencies are listed in the **go.mod** file;
- A **docker-compose** recipe is also provided in this repository.

# Implemetation decisions

- A simple three layered "Main/Entrypoint -> REST API -> Database"
architecture sufices for the solution here, i.e., main.go ->
api/server.go -> db/db.go;
- A way to uniquely identify password cards is required in order to
edit and delete a specifc card. To keep it simple, the decision was
to use a small and easy to handle UUID library available at:
https://pkg.go.dev/github.com/rs/xid;
- The project's files tree can be read as:
```
├── api/                       REST API package
│   ├── corsMiddleware.go
│   ├── loggerMiddleware.go
│   ├── server.go              API route handlers
│   └── server_test.go         Unit tests for API layer
├── db/                        Database layer package
│   ├── data.go                An initial load for the in memory data base
│   ├── db.go                  The database layer implementation
│   └── dbmock.go              A mocked implementation used for testing
└── main.go
```

# Running

To run the project, either build the binary locally and execute it,
i.e.:
```
> go build && ./velozient-backend
```

Or, use the docker-compose service, i.e., bring the virtualized
environment up and execute the desired commands in it:
```
> docker-compose up -d backend
> docker-compose exec backend sh -c "go build && ./velozient-backend"
```

# Testing strategy

Some unit tests are provided in the `package api`.

Improvements: the present unit tests could be incremented and some
other "integration test" strategy could be applied.
