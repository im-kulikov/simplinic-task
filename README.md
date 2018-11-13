# Simplinic Test Task

You need to design the «Service of configurations» as a set of endpoints to create, edit, and search of tagged
templates and configurations to interact with the frontend / UI.
The design consists from the description of API, interaction Protocol, DB schemes and technologies used.

## Non-functional requirements
Templates and configurations should be immutable: append-only approach in the database, changing the entity
generates a new version of the entity. The number of schemes in the system is up to several hundred, one scheme
can correspond to several thousand configurations. Endpoints response time with less than 100 milliseconds

## Authorization / authentication not needed: DMZ, intranet or out-of-scope
Json themselves on the backend only stored, syntactic or semantic parsing, in particular, search for the content
of documents on the backend within this task is not expected, but should be possible in the future with
the development of the system. Delete entities («not present»): the records remain in the database but is not available
for UI. Objectives: internal audit and the ability to recover from human error.

## Example of scheme (Person, https://tools.ietf.org/html/draft-handrews-json-schema-hyperschema-01)

```json
{
  "title": "Person",
  "type": "object",
  "required": ["firstName", "lastName"],
  "properties": {
    "firstName": {
      "type": "string"
    },
    "lastName": {
      "type": "string"
    },
    "age": {
        "description": "Age in years",
        "type": "integer",
        "minimum": 0
    },
    "friends": {
        "type" : "array",
        "items" : { "title" : "REFERENCE", "$ref" : "#" }
    }
  }
}
```

## Example of documents (for multiple environments):

```json
{
  "environment": "stage",
  "firstName": "Evgeniy",
  "lastName": "Kulikov",
  "age": 28,
  "friends": []
}
```

and multiple document for another environment:

```json
{
  "total": 2,
  "offset": 0,
  "items": [
       {
         "environment": "stage",
         "firstName": "Evgeniy",
         "lastName": "Kulikov",
         "age": 28,
         "friends": []
       },
       {
         "environment": "dev",
         "first_name": "Evgeniy",
         "last_name": "Kulikov",
         "age": 28,
         "friends": []
       }
   ]
}
```