# sre-bootcamp

This is the repository for the SRE bootcamp

### Running app locally 
- Make sure you have below tools installed 
    - [go](https://go.dev/doc/install)
    - [docker](https://docs.docker.com/engine/install/)
    - [docker-compose](https://docs.docker.com/compose/install/)
    - [make](https://www.gnu.org/software/make/)

- Clone the repository using the below command 

```bash
git clone https://github.com/JOSHUAJEBARAJ/sre-bootcamp.git
```

- Before setting up the application we have to setup few environment variables . Create the `.env` file in the current directory and store the values 

```bash
DB_USERNAME="root"
DB_PASSWORD="password"
DB_HOST="0.0.0.0"
DB_PORT="5432"
DB_DBNAME="dummy"
```

- Bring the DB UP 

```
make start-db
```


Do migration for the first time 

```
make db-migrate 
```

> Note execute this command only for the first time 

- Next run the application using the below command 

```bash
make run 
```

### Running using the containers


- Before running the application first we have to build the application 

```bash
IMAGE_TAG=0.0.1 make docker-app-build
```

- Next run the application using the below command 

```bash
IMAGE_TAG=0.0.1 make docker-app-run
```

### Cleaning 

- To stop the both containers

```bash
make docker-stop
```

- To remove the data 

```bash
make remove-db-data
```