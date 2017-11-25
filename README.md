# GO-PEX

[![CircleCI](https://circleci.com/gh/joaosilva2095/go-pex.svg?style=svg)](https://circleci.com/gh/joaosilva2095/go-pex)
[![Go Report Card](https://goreportcard.com/badge/github.com/joaosilva2095/go-pex)](https://goreportcard.com/report/github.com/joaosilva2095/go-pex)
[![GoDoc](https://godoc.org/github.com/joaosilva2095/go-pex?status.svg)](https://godoc.org/github.com/joaosilva2095/go-pex)

### A permissions system for Go structs

Developing APIs in Go is very common but so far there is no easy way to choose what to return accordingly
to the user that did the request and to the action.
To solve that, I created a library that allow developers to easily set permissions of the fields of a struct with Go tags.

## How it works

The system uses the _pex_ tag in each field to determine if a user has or not permission for that action in that field.

It is considered that a certain user type has permission if the permissions tag is not defined or if explicitly written
in the tag that for a certain action, it has permission.
Invalid actions or user types that doesn't exist in the tag are considered as **not** having permission.

## Tag structure

The permission tag is a set of pairs between user type and permission like `pex:"user:r,admin:rw"`.
In this case _user_ would have permission to _read_ while _admin_ would have permission to _read_ and _write_.

## Extract fields
Imagine you have this two structs

```go
type Person struct {
    ID     int  `pex:"user:r,admin:rw"`
    Name string `pex:"user:r,admin:rw" json:"full_name"`
}

type Employee struct {
    Parent
    Income float32 `pex:"user:,admin:rw"`
}
```

And you queried the database to get the employees. Now suppose you want to return (_ActionRead_) the result to a regular
user (userType = "user"). Of course you don't want to show the income of the employees to a regular user, you have to remove it.
For that you leave the permission of **user** empty in the **income** field as you can see above. Then you just have to call
extract fields function.

```go
fields := ExtractFields(employee, userType, ActionRead)
```

This will return an interface that contains all the fields in the struct that the user has permission.
The key in the result is the JSON key if the JSON tag exists otherwise its the field name.

```json
{
  "ID": 1,
  "full_name": "John Doe"
}
```

If the **userType = admin** the result has the income included

```json
{
  "ID": 1,
  "full_name": "John Doe",
  "Income": 1000.0
}
```

This can be applied to slices, pointers, maps and any kind of variables.

```json
[
  {
    "ID": 1,
    "full_name": "John Doe",
    "Income": 1000.0
  },
  {
    "ID": 2,
    "full_name": "Jack Sparrow",
    "Income": 9999.99
  }
]
```

## Clean struct
It is also possible to take advantage of the fields extraction to clean a struct, that is to set fields that user does
not have permission to their zero values and the others to the result of the field extraction.

```go
var cleanedObject *Employee
cleanedObject := CleanObject(employee, userType, ActionRead).(*Employee)
```

## Possible actions

`ActionRead`: 0 

`ActionWrite`: 1 

## Permission values

`PermissionNone`: Empty string

`PermissionRead`: "r"

`PermissionWrite`: "w"

`PermissionReadWrite`: "rw"
