package backup

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/jerrywang1981/go_backup_db2/db2"
	"github.com/jerrywang1981/go_backup_db2/tool"
)

func Backup(host, port, dbname, user, password, cert, generate, schema, output string) {
	db, _ := db2.Connect(host, port, dbname, user, password, cert)
	defer db.Disconnect()

	db.ReadAllTableSchema()

	m := tool.LoadSchemaTableMap(schema)
	exp, imp := db.GenerateSql(output, m)
	strCommand := fmt.Sprintf(
		"----generated by command : db2_backup -host=%s -port=%s -db=%s -user=%s -password=%s -json=%s -output=%s -generate=%s",
		host, port, dbname, user, password, schema, output, generate)

	exp = append([]string{strCommand}, exp...)
	imp = append([]string{strCommand}, imp...)

	switch generate {
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