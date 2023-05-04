
# TIX-ID Movie Ticketing
---

This project aims to develop an application similar to [TIX-ID](https://www.tix.id/), a well-established platform for booking tickets online.
There is three roles here
- Admin*
- Customer*
- Guest  
<sub>*Need Authentification</sub>

The program uses [go-gin](https://github.com/gin-gonic/gin) as the framework, MySQL as the database, and using [go-migrate](https://github.com/golang-migrate/migrate) for database migrations

It also employs various libraries:
- [Go-MySQL-Driver](https://github.com/go-sql-driver/mysql)
- [jwt-go](https://github.com/dgrijalva/jwt-go)
- [goCron](github.com/claudiu/gocron)
- [gin-swagger](https://github.com/swaggo/gin-swagger)
- [GoDotEnv](https://github.com/joho/godotenv)
- [Redis client for Go](github.com/go-redis/redis)
- [Gomail](https://github.com/go-gomail/gomail)

## Get Started
### API Documentation
To view the API documentation, you need to create `/docs` directory by running the following command in the terminal:
```
swag init
```
and open `<URL>/swagger/index.html`

### Enviroment variables (.env)
Before starting the program, you need to set the `.env` file first:
1. Create `.env` file in the root directory
2. Copy the enviroment variables from `.env.example`
3. Fill the variables

### Database Migration
Ensure that you have installed [go-migrate](https://github.com/golang-migrate/migrate). Before migrating the database, create a database in your MySQL.  
To run the database migrations:
- UP Migration
  ```
  migrate -path config/migrations -database "<database_address>" up
  ```
- DOWN Migration
  ```
  migrate -path config/migrations -database "<database_address>" down
  ```

> "Note: Replace `database_address` with `mysql://user:password@tcp(host:port)/dbname?query`"

### Docker
To start this project in docker:
1. Build the Docker Compose first
   ```
   docker compose build
   ```
2. Execute the Docker Compose useing 'up' command
   ```
   docker compose up
   ```
   this docker will run in port `80`

- Stopping Docker Compose
  ```
  docker compose down
  ```
- Migrate Database in Docker
  ```
  docker run -v {{ migration dir }}:/config/migrations --network host migrate/migrate -path=/config/migrations/ -database mysql://user:password@tcp(host:port)/dbname?query up
  ```
