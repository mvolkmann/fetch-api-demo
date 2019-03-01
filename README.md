# fetch-api-demo

This is a web application in Vue demonstrates
use of the Fetch API to talk to a REST server.
There are two implementations of the REST server,
one in Node.js (using Express) and one in Go (using gin).
Both talk to a PostgreSQL database.

![screenshot](./fetch-api-demo.png)

## To setup and start database server

Install the Postgres database server.
On a Mac this can be done with `brew install postgresql`
assuming homebrew is installed.

```bash
# Initialize the database.
sudo mkdir /usr/local/pgsql
sudo chown {your-user-name} /usr/local/pgsql
initdb -D /usr/local/pgsql/data

# Start the Postgres server.
pg_ctl -D /usr/local/pgsql/data start

# Create the "postgres" user.
create user postgres (enter "postgres" for the password)
```

## To use Go server (instead of Node.js server)

```bash
cd go
go get github.com/gin-gonic/gin # REST server library
go get github.com/lib/pq # PostgreSQL driver
go run main.go
```

## To use Node.js server (instead of Go server)

```bash
cd server
npm install # installs all dependencies
npm run dbsetup # recreates database tables, losing data
npm start # starts Express server
```

## To build and run web client

cd to the top project directory.

```bash
npm install # installs all dependencies
npm run serve # starts local HTTP server
```

browse localhost:8080
