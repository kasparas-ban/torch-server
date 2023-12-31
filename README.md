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

## Additional notes on TLS

If you're using AWS RDS database, to enable SSL/TSL encryption between the database instance and the server you'll need to install a certificate bundle. Donwload the bundle from [here](https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.SSL.html#UsingWithRDS.SSL.CertificatesAllRegions) and install it by following the instructions from [here](https://superuser.com/questions/437330/how-do-you-add-a-certificate-authority-ca-to-ubuntu).

## Fly.io database connection

To connect to the database forward the MySQL server port to your local machine using fly proxy:

`flyctl proxy 3306 -a torch-database`

## TODOS

- [x] Use nanoid's for userID's
- [x] Use nanoid's for itemID's
- [x] Check if code is SQL injection-prone
- [x] Refactor items
- [x] Stored Procedures for items
- [ ] Try to implement optional for Reccuring struct
- [ ] Setup CI/CD pipeline
- [ ] Do additional validation on the inputs
