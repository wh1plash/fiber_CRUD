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
## Prometheus metrics available on address  
```
http://localhost:3000/metrics
```
## Grafana dashboards
```
http://localhost:3001/
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

