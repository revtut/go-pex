# GO-PEX
### A permissions system for Go structs

Developing APIs in Go is very common but so far there is no easy way to choose what to return accordingly
to the user that did the request and to the action.
To solve that, I created a library that allow developers to easily set permissions of the fields of a struct with Go tags.

## How it works

The system uses the _pex_ tag in each field to determine if a user has or not permission for that action in that field.

If the tag is not found or the action is invalid, the system will consider that the user has permission for that field so
it will be added to the result.

## Tag structure

The permission tag is a set of numbers like `pex:"120123"`. Each index in the number string corresponds to a user type,
that is, imagine that a regular user has the **userType = 1**, then his permission would be **2**, which is the
corresponding index on the _120123_ string.

## Extract fields
Imagine you have this two structs

```go
type Person struct {
    ID     int  `pex:"11"`
    Name string `pex:"31" json:"full_name"`
}

type Employee struct {
    Parent
    Income float32 `pex:"30"`
}
```

And you queried the database to get the employees. Now suppose you want to return (_ActionRead_) the result to a regular
user (userType = 1). Of course you don't want to show the income of the employees to a regular user, you have to remove it.
For that you set the permission to **0** in the **income** field in the **index = 1** and then you just have to call
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

If the **userType = 0** the result has the income included

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

`PermissionNone`: 0

`PermissionRead`: 1

`PermissionWrite`: 2

`PermissionReadWrite`: 3
