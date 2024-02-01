# Anakonda

Anakonda is a container-based task runner to run user tasks in Docker.

## System Design Diagram

<p align="center"><img src='/docs/files/anakonda-system-design.png' alt='Anakonda System Design Diagram' /></p>

## Give a Star! :star:

If you like this repo or found it helpful, please give it a star. Thanks!

## Used Tools

1. [Echo as web framework](https://echo.labstack.com)
2. [Redis for queue](https://github.com/redis/redis)
3. [Asynq for task queue](https://github.com/hibiken/asynq)
4. [Postgresql as main database engine](https://github.com/postgres/postgres)
5. [PgAdmin as database management tool](https://github.com/pgadmin-org/pgadmin4)
6. [Koanf for configurations](https://github.com/knadh/koanf)
7. [Zap for logging](https://github.com/uber-go/zap)
8. [Gorm as ORM](https://github.com/go-gorm/gorm)
9. [Swagger for documentation](https://github.com/swaggo/swag)
10. [Validator for endpoint input Validation](https://github.com/go-playground/validator)
11. Docker compose for run project with all dependencies in docker

## How to run

### Docker start

```
docker compose -f "docker/docker-compose.yml" up -d --build
```

#### Web API

##### Run local manually [http://localhost:8080](http://localhost:8080)

#### PgAdmin

##### [http://localhost:8090](http://localhost:8090)

```
Username: h.darani@gmail.com
Password: 123456
```

Postgres Server info:

```
Host: postgres_container
Port: 5432
Username: postgres
Password: admin
```

## Project preview

## Swagger

##### [http://localhost:8090](http://localhost:8080/swagger/index.html)
