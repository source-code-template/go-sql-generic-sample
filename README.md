# go-sql-generic-sample

In modern software development, efficient data management is a crucial aspect of building robust applications. Leveraging a generic CRUD (Create, Read, Update, Delete) repository can significantly simplify data operations and enhance code maintainability.
Inspired by [CrudRepository](https://docs.spring.io/spring-data/commons/docs/current/api/org/springframework/data/repository/CrudRepository.html) of Spring in Java, I developed a Generic CRUD Repository in Golang that supports "database/sql". This article introduces the structure and functionalities of this repository.

### Overview
The Generic CRUD Repository in Golang provides a set of standard methods to perform CRUD operations on a database. It implemented this interface:
```go
package repository

import "context"

type GenericRepository[T any, K any] interface {
  All(ctx context.Context) ([]T, error)
  Load(ctx context.Context, id K) (*T, error)
  Create(ctx context.Context, model T) (int64, error)
  Update(ctx context.Context, model T) (int64, error)
  Patch(context.Context, map[string]interface{}) (int64, error)
  Save(ctx context.Context, model T) (int64, error)
  Delete(ctx context.Context, id K) (int64, error)
}
```
- All: Retrieves all records from a database table.
- Load: Fetches a specific record by its ID.
- Create: Adds a new record to the database.
- Update: Modifies an existing record.
- Delete: Removes a record from the database.
- Save: Inserts a new record or updates an existing one.
- Patch: Perform a partial update of a resource

### Benefits
- <b>Simplicity</b>: provides a set of standard CRUD (Create, Read, Update, Delete) operations out of the box, reducing the amount of boilerplate code developers need to write.
  - Especially, it provides "Save" method, to build an insert or update statement, specified for Oracle, MySQL, MS SQL, Postgres, SQLite.
- <b>Consistency</b>: By using Repository, the code follows a consistent pattern for data access across the application, making it easier to understand and maintain.
- <b>Rapid Development</b>: reducing boilerplate code and ensuring transactional integrity.
- <b>Flexibility</b>: offers flexibility and control over complex queries, because it uses "database/sql" at GO SDK level.
- <b>Type Safety</b>: being a generic interface, it provides type-safe access to the entity objects, reducing the chances of runtime errors.
- <b>Learning Curve</b>: it supports utilities at GO SDK level. So, a developer who works with "database/sql" at GO SDK can quickly understand and use it.
- <b>Composite primary key</b>: it supports composite primary key.
  - You can look at the sample at [go-sql-composite-key](https://github.com/go-tutorials/go-sql-composite-key).
  - In this sample, the company_users has 2 primary keys: company_id and user_id
  - You can define a GO struct, which contains 2 fields: CompanyId and UserId
    ```go
    package model
    
    type UserId struct {
      CompanyId string `json:"companyId" gorm:"column:company_id;primary_key"`
      UserId    string `json:"userId" gorm:"column:user_id;primary_key"`
    }
    ```

### Use Cases for Generic CRUD Repository
- <b>Basic CRUD Operations</b>: ideal for applications that require standard create, read, update, and delete operations on entities.
- <b>Prototyping and Rapid Development</b>: useful in the early stages of development for quickly setting up data access layers.
- <b>Admin/Back Office Web Application</b>: In admin application where services often perform straightforward CRUD operations, CrudRepository can be very effective.
- <b>Microservices</b>: In microservices architectures where services often perform straightforward CRUD operations, CrudRepository can be very effective.

### Conclusion
- The Generic CRUD Repository in Golang provides a robust and flexible solution for managing database operations. By abstracting common CRUD operations into a generic repository, developers can write cleaner, more maintainable code. This repository structure is inspired by CrudRepository of Spring in Java and adapted for the Go programming language, leveraging the power of "database/sql".
- This implementation can be further extended to include additional functionalities and optimizations based on specific project requirements. By adopting this generic repository pattern, developers can streamline their data management tasks and focus on building feature-rich applications.

## Architecture
### Simple Layer Architecture
![Layer Architecture](https://cdn-images-1.medium.com/max/800/1*JDYTlK00yg0IlUjZ9-sp7Q.png)

### Layer Architecture with full features
![Layer Architecture with standard features: config, health check, logging, middleware log tracing](https://cdn-images-1.medium.com/max/800/1*8UjJSv_tW0xBKFXKZu86MA.png)

#### [core-go/search](https://github.com/core-go/search)
- Build the search model at http handler
- Build dynamic SQL for search
  - Build SQL for paging by page index (page) and page size (limit)
  - Build SQL to count total of records
### Search users: Support both GET and POST 
#### POST /users/search
##### *Request:* POST /users/search
In the below sample, search users with these criteria:
- get users of page "1", with page size "20"
- email="tony": get users with email starting with "tony"
- dateOfBirth between "min" and "max" (between 1953-11-16 and 1976-11-16)
- sort by phone ascending, id descending
```json
{
    "page": 1,
    "limit": 20,
    "sort": "phone,-id",
    "email": "tony",
    "dateOfBirth": {
        "min": "1953-11-16T00:00:00+07:00",
        "max": "1976-11-16T00:00:00+07:00"
    }
}
```
##### GET /users/search?page=1&limit=2&email=tony&dateOfBirth.min=1953-11-16T00:00:00+07:00&dateOfBirth.max=1976-11-16T00:00:00+07:00&sort=phone,-id
In this sample, search users with these criteria:
- get users of page "1", with page size "20"
- email="tony": get users with email starting with "tony"
- dateOfBirth between "min" and "max" (between 1953-11-16 and 1976-11-16)
- sort by phone ascending, id descending

#### *Response:*
- total: total of users, which is used to calculate numbers of pages at client 
- list: list of users
```json
{
    "list": [
        {
            "id": "ironman",
            "username": "tony.stark",
            "email": "tony.stark@gmail.com",
            "phone": "0987654321",
            "dateOfBirth": "1963-03-24T17:00:00Z"
        }
    ],
    "total": 1
}
```

## API Design
### Common HTTP methods
- GET: retrieve a representation of the resource
- POST: create a new resource
- PUT: update the resource
- PATCH: perform a partial update of a resource, refer to [core-go/core](https://github.com/core-go/core) and [core-go/sql](https://github.com/core-go/sql)
- DELETE: delete a resource

## API design for health check
To check if the service is available.
#### *Request:* GET /health
#### *Response:*
```json
{
    "status": "UP",
    "details": {
        "sql": {
            "status": "UP"
        }
    }
}
```

## API design for users
#### *Resource:* users

### Get all users
#### *Request:* GET /users
#### *Response:*
```json
[
    {
        "id": "spiderman",
        "username": "peter.parker",
        "email": "peter.parker@gmail.com",
        "phone": "0987654321",
        "dateOfBirth": "1962-08-25T16:59:59.999Z"
    },
    {
        "id": "wolverine",
        "username": "james.howlett",
        "email": "james.howlett@gmail.com",
        "phone": "0987654321",
        "dateOfBirth": "1974-11-16T16:59:59.999Z"
    }
]
```

### Get one user by id
#### *Request:* GET /users/:id
```shell
GET /users/wolverine
```
#### *Response:*
```json
{
    "id": "wolverine",
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```

### Create a new user
#### *Request:* POST /users 
```json
{
    "id": "wolverine",
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```
#### *Response:* 1: success, 0: duplicate key, -1: error
```json
1
```

### Update one user by id
#### *Request:* PUT /users/:id
```shell
PUT /users/wolverine
```
```json
{
    "username": "james.howlett",
    "email": "james.howlett@gmail.com",
    "phone": "0987654321",
    "dateOfBirth": "1974-11-16T16:59:59.999Z"
}
```
#### *Response:* 1: success, 0: not found, -1: error
```json
1
```

### Patch one user by id
Perform a partial update of user. For example, if you want to update 2 fields: email and phone, you can send the request body of below.
#### *Request:* PATCH /users/:id
```shell
PATCH /users/wolverine
```
```json
{
    "email": "james.howlett@gmail.com",
    "phone": "0987654321"
}
```
#### *Response:* 1: success, 0: not found, -1: error
```json
1
```

#### Problems for patch
If we pass a struct as a parameter, we cannot control what fields we need to update. So, we must pass a map as a parameter.
```go
type UserService interface {
    Update(ctx context.Context, user *User) (int64, error)
    Patch(ctx context.Context, user map[string]interface{}) (int64, error)
}
```
We must solve 2 problems:
1. At http handler layer, we must convert the user struct to map, with json format, and make sure the nested data types are passed correctly.
2. At repository layer, from json format, we must convert the json format to database format (in this case, we must convert to column)

#### Solutions for patch  
At http handler layer, we use [core-go/core](https://github.com/core-go/core), to convert the user struct to map, to make sure we just update the fields we need to update
```go
import "github.com/core-go/core"

func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
    var user User
    userType := reflect.TypeOf(user)
    _, jsonMap := core.BuildMapField(userType)
    body, _ := core.BuildMapAndStruct(r, &user)
    json, er1 := core.BodyToJson(r, user, body, ids, jsonMap, nil)

    result, er2 := h.service.Patch(r.Context(), json)
    if er2 != nil {
        http.Error(w, er2.Error(), http.StatusInternalServerError)
        return
    }
    respond(w, result)
}
```

### Delete a new user by id
#### *Request:* DELETE /users/:id
```shell
DELETE /users/wolverine
```
#### *Response:* 1: success, 0: not found, -1: error
```json
1
```

## Common libraries
- [core-go/health](https://github.com/core-go/health): include HealthHandler, HealthChecker, SqlHealthChecker
- [core-go/config](https://github.com/core-go/config): to load the config file, and merge with other environments (SIT, UAT, ENV)
- [core-go/log](https://github.com/core-go/log): logging
- [core-go/middleware](https://github.com/core-go/log): middleware log tracing

### core-go/health
To check if the service is available, refer to [core-go/health](https://github.com/core-go/health)
#### *Request:* GET /health
#### *Response:*
```json
{
    "status": "UP",
    "details": {
        "sql": {
            "status": "UP"
        }
    }
}
```
To create health checker, and health handler
```go
    db, err := sql.Open(conf.Driver, conf.DataSourceName)
    if err != nil {
        return nil, err
    }

    sqlChecker := s.NewSqlHealthChecker(db)
    healthHandler := health.NewHealthHandler(sqlChecker)
```

To handler routing
```go
    r := mux.NewRouter()
    r.HandleFunc("/health", healthHandler.Check).Methods("GET")
```

### core-go/config
To load the config from "config.yml", in "configs" folder
```go
package main

import "github.com/core-go/config"

type Root struct {
    DB DatabaseConfig `mapstructure:"db"`
}

type DatabaseConfig struct {
    Driver         string `mapstructure:"driver"`
    DataSourceName string `mapstructure:"data_source_name"`
}

func main() {
    var conf Root
    err := config.Load(&conf, "configs/config")
    if err != nil {
        panic(err)
    }
}
```

### core-go/log *&* core-go/middleware
```go
import (
    "github.com/core-go/config"
    "github.com/core-go/log"
    m "github.com/core-go/middleware"
    "github.com/gorilla/mux"
)

func main() {
    var conf app.Root
    config.Load(&conf, "configs/config")

    r := mux.NewRouter()

    log.Initialize(conf.Log)
    r.Use(m.BuildContext)
    logger := m.NewStructuredLogger()
    r.Use(m.Logger(conf.MiddleWare, log.InfoFields, logger))
    r.Use(m.Recover(log.ErrorMsg))
}
```
To configure to ignore the health check, use "skips":
```yaml
middleware:
  skips: /health
```

#### To run the application
```shell
go run main.go
```
