# db2 backup tool




### get exeutable file
For Mac
```
CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64  go build -o db2_backup main.go
```
For Linux
```
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -o db2_backup main.go
```
For Windows
```
CGO_ENABLED=0 GOOS=windows  GOARCH=amd64  go build -o db2_backup.exe main.go
```


## Usage

Example
```
./go_backup_db2 -host=127.0.0.1 -port=50000 -db=DB_NAME -user=db2inst1 -password=passw0rd -json=./test/schema.json -output=. -generate=export
```

### Params
|param|description|
| ----- | ----- | 
| host | the db2 server ip |
| port| the db2 server port|
| db| the db2 db name|
| user| the db2 user id|
| password | the db2 password|
| json| the schema/table you want to export/import, [example](test/schema.json) |
| output | the export data file location|
| generate |possible value: `both`, `export`, `import`|




# Maintainers
[@jerrywang1981](https://github.com/jerrywang1981)

# Contributors
[![](https://github.com/jerrywang1981.png?size=50)](https://github.com/jerrywang1981)

# License
MIT License

