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

### Params
|param|description|
| ----- | ----- | 
| host | the db2 server ip |
