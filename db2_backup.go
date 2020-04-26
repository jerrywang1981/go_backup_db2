package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jerrywang1981/go_backup_db2/backup"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println(`

    the format would be : 
    ./db2_backup -host=127.0.0.1 -port=50000 -db=DB_NAME -user=db2inst1 -password=passw0rd -json=./test/schema.json -output=. -cert=xxx.arm -generate=export
    `)
		os.Exit(0)
	}

	host := flag.String("host", "127.0.0.1", "db2 server ip address")
	port := flag.String("port", "50000", "db2 server port")
	cert := flag.String("cert", "", "cert location")
	db := flag.String("db", "", "database name")
	user := flag.String("user", "", "db2 user")
	password := flag.String("password", "", "db2 password")
	generate := flag.String("generate", "both", "possible export, import, both")
	schema := flag.String("json", "", "json file")
	output := flag.String("output", ".", "the file where to export")
	flag.Parse()
	backup.Backup(*host, *port, *db, *user, *password, *cert, *generate, *schema, *output)
}

