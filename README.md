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
## Grafana homepage
```
http://localhost:3001/
```

### Grafana datasources prometheus
```
http://prometheus:9090
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

## Testing
### for now it need to be created database 'test' in postgres DB
Run test by
```
make test
```

## GitLab CI/CD
### Docker container for local runner on remote server
```
docker run -d --name gitlab-runner \
--restart always \
-v $HOME/gitlab-runner/config:/etc/gitlab-runner \
-v /var/run/docker.sock:/var/run/docker.sock \
gitlab/gitlab-runner:alpine

docker run --rm -it \
-v $HOME/gitlab-runner/config:/etc/gitlab-runner \
gitlab/gitlab-runner:alpine register

//glrt-VF8Ga5npyanSiViES_JXLm86MQpwOjE1bDd4dAp0OjMKdToybGJhdRg.01.1j1aoxkl9

> https://gitlab.com/
> tocken from your GitLab repository - CI/CD Settings - Register runner
> deploy
> docker
> docker:dind
```
### important to edit ~/gitlab-runner/config/config.toml
```
privileged = true
volumes = ["/var/run/docker.sock:/var/run/docker.sock", "/cache"]
```

## Kubernetes. K3d as a lightweight Kubernetes distribution
install:
```
curl -s https://raw.githubusercontent.com/k3d-io/k3d/main/install.sh | bash
```
Create cluster
```
k3d cluster create dev-cluster --agents 2 --port "3030:80@loadbalancer"
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