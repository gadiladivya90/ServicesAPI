# konnect-takehome

Konnect take home exercise. Service Rest API

## Getting Started

### Makefile

There is a `Makefile` included in the repository. To get started:

```shell
1. `go mod download` - download modules
2. `make up` - bring up database and application
```

More helpful commands can be found in Makefile:

```shell

 make run                  -> just run application
 make restart              -> runtime: restart environment

```

### Implemented Endpoints

$BASEURL : <http://localhost:9090>

GET $BASEURL/services
GET $BASEURL/services/{service_id}
POST $BASEURL/services
PUT $BASEURL//services/{service_id}
Delete $BASEURL//services/{service_id}
GET $BASEURL//services/{service_id}/versions
POST $BASEURL//services/{service_id}/versions
DELETE $BASEURL//services/{service_id}/versions

### Implementation details

This entire API is based on Hexagonal architecture(Ports and Adapters). All the wiring across db<->request is based in this architecture.
<https://www.qwan.eu/2020/08/20/hexagonal-architecture.html>

Database
-Using Postgres relational database, can be spinned up by `make db`. Persisted memory
-Due to time crunch utilized go packages that I am familiar with(sqlx)
-Using ENT would have made database interactions easier/cleaner
-Spinned up docker via docker-compose and using sql script to create services database and add tables

Multiplexer: Used gorilla Mux as URL router/dispatcher

### Filtering, validation and Pagination

Added validation of payload and pagination(getManyServices request). Few examples of usage

```shell

filtering:
GET $BASEURL/services?filter=name%3DserviceName

offset/limit
GET $BASEURL/services?filter=name%3DserviceName?offset=2&limit=3


example: 
 POST $BASEURL/services
 {
 "name":"newService", 
 "description":"demo service" 
}
```

### Error Handling/Logging

Handled and returning errors based HTTP norms. Added framework for logging and errorHandling
Existing code can definitely use more debugging logs and inline comments

### Unit Test

To run test `make tests`. This will take care of generating mocks and runnning tests

### Improvements/Pending implementation

ServiceVersions and servicepackage coupling queries can be improved
Add unit test coverage
Add inline comments in code for readability
