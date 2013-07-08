#!/bin/sh

go run drop_database.go
go run create_schema.go
go run seed.go
