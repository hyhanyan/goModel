#!/bin/bash
go clean 
GOOS=linux
GOARCH=amd64
go build main.go
