# Currency API

[![Golang](https://golang.org/lib/godoc/images/go-logo-blue.svg)](https://golang.org/)

Currency API is a simple microservice capable of handling [historical rates](https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml) via REST API and authorize users using JWT Authentication.

### Tech

Currency API uses a number of open source projects to work properly:

* [Golang](https://golang.org/) - Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.
* [GORM](https://gorm.io/) - The fantastic ORM library for Golang
* [gorilla/mux](https://github.com/gorilla/mux) - Package mux implements a request router and dispatcher.
* [PostgreSQL](https://www.postgresql.org/) - The World's Most Advanced Open Source Relational Database
* [Docker](https://www.docker.com/) - Securely build, share and run modern applications anywhere
* [jwt-go](https://github.com/dgrijalva/jwt-go) - A go (or 'golang' for search engine friendliness) implementation of JSON Web Tokens

### Installation

Currency API requires [Docker](https://www.docker.com/) and [docker-compose](https://docs.docker.com/compose/) to run.

Install Docker and docker-compose to start the server
 - [Docker Desktop on Windows](https://docs.docker.com/docker-for-windows/install/)
 - [Docker on Linux](https://docs.docker.com/install/linux/docker-ce/centos/)
 - [Docker Desktop on MacOS](https://docs.docker.com/docker-for-mac/install/)
 - [Install docker-compose](https://docs.docker.com/compose/install/)

```sh
$ cd currency-api
$ docker-compose up
```

### Usage
    Currently uses inMemoryValidation please use: user123/pass123 or useruser/passpass
    - POST "https://{HOST}:9988/api/auth"
        {
            "username": user123,
            "password": pass123
        }
        returns: {JwtToken}
        
    Requires: Header {"Authorization": "Bearer {JwtToken}"}
    - GET "https://{HOST}:9988/rates/latest"
    - GET "https://{HOST}:9988/rates/{YYYY-MM-DD}"
    - GET "https://{HOST}:9988/rates/analyze"
### Todos
 - Validate credentials against DB

