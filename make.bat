@echo off
set GOOS=linux
set GOARCH=arm64
set GOPATH=I:\go
go build -ldflags "-extldflags '-static'"