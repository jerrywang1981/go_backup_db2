package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jerrywang1981/go_backup_db2/backup"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println(`

    the format would be : 
    ./db2_backup -host=127.0.0.1 -port=50000 -db=DB_NAME -user=db2inst1 -password=passw0rd -json=./test/schema.json -output=. -generate=export
    `)
		os.Exit(0)
	}

	host := flag.String("host", "127.0.0.1", "db2 server ip address")
	port := flag.String("port", "50000", "db2 server port")
	db := flag.String("db", "", "database name")
	user := flag.String("user", "", "db2 user")
	password := flag.String("password", "", "db2 password")
	generate := flag.String("generate", "both", "possible export, import, both")
	schema := flag.String("json", "", "json file")
	output := flag.String("output", ".", "the file where to export")
	flag.Parse()
	backup.Connect(*host, *port, *db, *user, *password)
	defer backup.Disconnect()
	backup.ReadAllTableSchema()

	m := backup.LoadSchemaTableMap(*schema)
	exp, imp := backup.GenerateSql(*output, m)
	switch *generate {
	case "both":
		writeSql(exp, "export.sql")
		writeSql(imp, "import.sql")
	case "export":
		writeSql(exp, "export.sql")
	case "import":
		writeSql(imp, "import.sql")
	default:
	}
}

func writeSql(sqls []string, fileName string) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, line := range sqls {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatalf("Got error while writing to a file. Err: %s", err.Error())
		}
	}
	writer.Flush()
}
