# Torch-server

API server for the Torch frontend.

## Testing

Testing basic CRUD operations on the server requires a test database. With the following comand Docker initializes this test database

`cd tests/ && docker compose up -d`

To stop and remove Docker containers run

`docker-compose down`

You still need to remove server image by hand whenever you update your server code, otherwise the old version of your code will be used.
