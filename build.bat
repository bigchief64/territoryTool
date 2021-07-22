@echo off
go generate
go build -ldflags "-H windowsgui" -o territory_tool.exe 
