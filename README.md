# CRUD API by using Fiber framework

## Start application
### first of all you must boot up yor DB and metrics instances.
```
> docker-compose up
```
### then run server
```
> make run
```
### Add user
```
http://localhost:3000/api/v1/user

JSON body:
{
    "firstName": "Req6",
    "lastName": "Test",
    "email": "basr@foo.com",
    "password": "hunter123"
}
```
### Get user by ID
```
http://localhost:3000/api/v1/user/:id
```
### List all users
```
http://localhost:3000/api/v1/users
```
### Update user
```
http://localhost:3000/api/v1/user/:id

JSON body:
{
    "firstName": "f666",
    "lastName": "f6666",
    "email": "exampl@mail.com"
}
```
### Delete user
```
http://localhost:3000/api/v1/user/:id
```
## Prometheus metrics available on address  
```
http://localhost:3000/metrics
```
## Grafana dashboards
```
http://localhost:3001/
```

### Grafana Dashboard Template(JSON)
```
Paste all from custom_metrics.json into Grafana → Dashboards → + Import → Paste JSON
```

### Default dashboards
```
http://localhost:3001/d/go-app-handler-metrics
```
## Exanple of .env file
```
PG_HOST="localhost"
PG_PORT=5444
PG_USER="postgres"
PG_PASS="postgres"
PG_DB_NAME="Fiber_CRUD"
```

## Useful tools
### go-callvis
This tool build interactive diagram of calling directly from source code
Install:
```
go install github.com/ofabry/go-callvis@latest
```
Run:
```
go-callvis ./...
```