# Notes on this project and Go

## Struct Tags

Struct tags are a way to attach metadata to struct fields. They provide additional information about how a struct field should be treated or interpreted by various libraries and tools that work with Go structs.

Here's a breakdown of the tags in the ConvertedBookmark struct:

json:"hash": This is a JSON struct tag. It tells the Go JSON encoder and decoder how to serialize and deserialize the struct field when converting between Go structs and JSON data. In this case, it specifies that the field name in the JSON representation should be "hash." It's a way to control the mapping between Go field names and JSON keys.

gorm:"unique": This is a GORM struct tag. It provides instructions to the GORM library on how to interact with the database. In particular, the "unique" tag indicates that the Hash field should have a unique constraint in the database schema. This constraint ensures that no two rows in the database can have the same value in the Hash column. It's a way to specify constraints and behaviors related to database interactions for the field.

Struct tags are specific to Go and are not found in most other programming languages. They are a powerful feature in Go, allowing you to provide metadata and instructions to various packages and tools that work with struct types. Other common uses of struct tags include specifying field names for database column mapping, defining validation rules, and more. Different libraries and tools may interpret and use struct tags in their own ways based on their specific requirements.

## SQL Queries

using a tool in the terminal or a gui like [sqlite studio](https://sqlitestudio.pl/) the following simple sql queries can be used as a guide to query the database


```sql
SELECT time,tags, href,description, extended,hash
  FROM converted_bookmarks
ORDER BY time ASC;
```
or a specific record based upon its hash

```sql
SELECT time,tags, href,description, extended,hash
  FROM converted_bookmarks
  WHERE hash="e723b64c2df57f09e9fc38d8efe36d92"
ORDER BY time ASC;
```
or a matched tag

```sql
SELECT time,tags, href,description, extended,hash
  FROM converted_bookmarks
  WHERE tags LIKE "%linux%"
ORDER BY time ASC;
```