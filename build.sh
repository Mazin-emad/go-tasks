#!/bin/bash
go build -tags netgo -ldflags '-s -w' -o app ./cmd/main.go
