# Bonds Calculator
Bonds Calculator is a web service for tracking bond yields.
Its functions include calculating current income, calculating current income and income statistics for your bond portfolio.

### How to run project
Latest release build with docker compose
```sh
docker compose -f docker-compose-release.yml up web-release -d
```

Dev build with docker compose
```sh
docker compose -f docker-compose.yml up web -d
```

Run tests with docker compose
```sh
docker compose -f docker-compose.yml up web-tests --exit-code-from web-tests
```

You can run project without docker using `go run ./...` or run tests using `go test -v ./...` but make sure you configured all dependencies by yourself:
- generate code using `go generate` commands
- run postgres server
- run redis server
- move [bonds_calculator_client](https://github.com/Wanket/bonds_calculator_client) build artefacts to ./public folder

#### Shell Enviroments

You can change project settings by using Shell Enviroments:
- Moex client settings
  - CONFIG_MOEX_CLIENT_QUEUE_SIZE - Number of concurrent connections to Moex API, default 10
- HTTP settings
  -	CONFIG_HTTP_PORT - default 8080
  -	CONFIG_REQUEST_TIMEOUT - default 100ms
- JWT settings
  - CONFIG_ACCESS_TOKEN_TTL - default 15m
  - CONFIG_REFRESH_TOKEN_TTL - default 168h (1 week)
- Database settings
  -	CONFIG_DB_USER - default bonds_calculator
  -	CONFIG_DB_PASSWORD - default bonds_calculator
  -	CONFIG_DB_HOST - default localhost
  -	CONFIG_DB_PORT - default 5432
  -	CONFIG_DB_NAME - default - default bonds_calculator
- Redis settings
  -	CONFIG_REDIS_HOST - default localhost 
  -	CONFIG_REDIS_PORT - default 6379
  -	CONFIG_REDIS_PASSWORD - default ""
  -	CONFIG_REDIS_DB - default 0

### TODO list
- [x] Add Calculator models
  - [x] Tests and Fuzz Tests
- [x] Add Moex API support
- [x] Add Static services/controllers
  - [x] Search
  - [x] BondInfo
  - [x] StaticCalsulator
  - [x] StaticStore
  - [x] TimerService
  - [x] Add Tests for they
- [x] Use Docker and Docker Compose
- [x] Add linters
- [x] Add Redis and Postrges support
  - [ ] Integration tests
- [x] Setup GitHub Actions
- [x] Add Auth/Register support
  - [x] JWT Token support
  - [ ] Logins count/Tokens count limit
  - [ ] Mail registration/change password notifications
  - [ ] Add Tests
- [x] Simplify wire config
- [ ] Add fuzz testing to ci
- [ ] Add coverage to ci
- [ ] Add User services/controllers
  - [ ] UserCalculator
  - [ ] UserStore
  - [ ] UserBondsInfo
  - [ ] UserStatistic
  - [ ] Add Tests for they
- [ ] Add autotests
