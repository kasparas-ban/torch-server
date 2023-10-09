# Torch-server

API server for the Torch frontend.

## Testing

Testing basic CRUD operations on the server requires a test database. There are two ways to setup the database and run the tests.

The default way requires no involvment from the user and works by automatically creating, launching, and removing a Docker container that contains the database. Run it with the usual command for go tests:

`go test ./tests`

Container setup, however, does take some time (>30s). To avoid waiting everytime a test is run, you can lauch a Docker container locally and reuse it across different test runs. To setup and lauch the container run the following command

`cd tests/ && docker compose up -d`

To stop and remove Docker containers run

`docker-compose down`

Note that all images built this way will be cached and reused with each later launch.

## TODOS

- [x] Use nanoid's for userID's
- [x] Use nanoid's for itemID's
- [ ] Check if code is SQL injection-prone
- [ ] Try to implement optional for Reccuring struct
- [x] Refactor items
- [ ] Stored Procedures for items
