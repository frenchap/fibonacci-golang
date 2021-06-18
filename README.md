# fibonacci-golang
## Run in windows:

### Prep
- Make sure docker desktop is installed and running
- Stop any running postgres containers
- If you have postgres installed locally and running, stop it in the services settings

### Powershell Commands (using your user is fine; admin ok, but not required)

```
$env:FI_API_DB_USER = 'postgres'
$env:FI_API_DB_PASSWORD = '12345-luggage-combo'
$env:FI_API_DB_NAME = 'postgres'
$env:FI_API_DB_PORT = '5432'
$env:FI_API_UPPER_BOUND = '50000'
docker run --name postgres-db  -e POSTGRES_PASSWORD=12345-luggage-combo -e POSTGRES_DB=postgres -d -p 5432:5432 postgres
go run main.go
```

## Run in linux
### Prep
- Make sure docker is installed and running (docker desktop if WSL2)
- Stop any running postgres containers
- If postgres is installed and running locally, stop the service

###
- Note, the tests won't currently run in WSL2 linux environments

### Shell commands

```
export FI_API_DB_USER="postgres"
export FI_API_DB_PASSWORD="12345-luggage-combo"
export FI_API_DB_NAME="postgres"
export FI_API_DB_PORT="5432"
export FI_API_UPPER_BOUND="50000"
docker run --name postgres-db  -e POSTGRES_PASSWORD=12345-luggage-combo -e POSTGRES_DB=postgres -d -p 5432:5432 postgres
go run main.go
```

## Cleanup

Use `docker ps` to identify any created images, `docker stop <id>` to stop the container, and `docker rm <id>` to remove them