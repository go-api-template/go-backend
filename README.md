# go-backend

This is a template for a RESTful API written in Go.


## Installation

**Clone the repository**

```bash
git clone https://github.com/go-api-template/go-backend
```

**Setup dependencies**

```bash
go mod download
```

**Get the list of all predefined tasks**

```bash
task --list
```

**Build all (docker images & go) then run the application**

```bash
task build run
```

## Introduction

## Folder Structure

``` 
├── controllers
├── models
├── modules
├── routes
├── services
``` 

### Controllers

Controllers are responsible for handling the HTTP requests coming into the router.
The controller layer should not implement service logic and data access.
The service and data access layers should be done separately.

Controllers must implement services through their interface.
Service interface implementations should NOT be done in the controller to maintain decoupled logic.
The implementation will be injected during compile time.

### Models

The models folder contains the structs that represent the data in the database.

### Modules

### Routes

### Services

The service layer is responsible for implementing the business logic of the application.



