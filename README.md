# API Raya OJT

## Requirement 

Installed Go 
```bash
$ go version
go version go1.18.4 linux/amd64
```

Create a postgresql database
```sql
CREATE DATABASE api_raya_ojt
```

## Setup Project

Clone Project 
```bash
git clone https://github.com/ajikamaludin/api-raya-ojt
cd api-raya-ojt
```

Run Install Depedency
```bash
go mod tidy
```

Create env file
```bash
cp .env.example .env
```

Edit `.env` file change your database connection, redis connection and google cloud credentials 
```
APP_NAME=note-app
APP_ENV=dev
APP_PORT=3000

DB_NAME=test
DB_HOST=localhost
DB_PORT=5432
DB_USER=aji
DB_PASS=eta
DB_TIMEZONE=Asia/Jakarta
DB_ISMIGRATE=true

JWT_SECRET=IyMjIEdvIEZpYmVyIChCYXNpYyBTZXR1cCk.
JWT_EXPIRED_SECOND=3600

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

GOOGLE_PROJECT_NAME=project-name
GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json
```

Run Rest Api Project (Keep Running) (Automated Migrate Database and Seed)
```bash
go run .
```

Run Service Transaction (Keep Running)
```bash
go run services/main.go
```
## REST Api docs
[Postman Collection](https://raw.githubusercontent.com/ajikamaludin/api-raya-ojt/dev/assets/postman/ApiRaya.Postman_collection.json)

[Postman Environment](https://raw.githubusercontent.com/ajikamaludin/api-raya-ojt/dev/assets/postman/ApiRaya.Postman_environment.json)

Public Api Doc : https://documenter.getpostman.com/view/1829038/UzXM1JEa

![run results](https://github.com/ajikamaludin/api-raya-ojt/raw/dev/assets/results.png)
## Code overview

### Folders

- `app/models` - Contains all the GORM models and models for request and response
- `app/controllers` - Contains all the controllers
- `app/configs` - Contains all the application configuration files
- `app/repository` - Contains all function to access database layer and redis 
- `app/services/services.go` - Contains all package/service use by app to interact
- `pkg/` - Contains all third party librari to access service like gorm to access database or go-redis/client to access redist
- `router/` - Contains all the app routes 
- `main.go` - Main app file to start the app
- `services` - Contains service app that works with app

## Project Design

### REST Api Design with Swagger Editor 
https://raw.githubusercontent.com/ajikamaludin/api-raya-ojt/dev/assets/schema.yaml
### Database Table Design
![table structure](https://github.com/ajikamaludin/api-raya-ojt/raw/dev/assets/case1_database.png)
### TechStack In Plan
![tech stack](https://github.com/ajikamaludin/api-raya-ojt/raw/dev/assets/case1_techstack.png)