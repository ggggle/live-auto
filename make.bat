@echo off
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-extldflags '-static'"