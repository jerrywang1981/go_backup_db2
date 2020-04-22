package backup

import (
	"flag"
	"fmt"

	_ "github.com/ibmdb/go_ibm_db"
	"github.com/jmoiron/sqlx"
)

//global db connection
var db *sqlx.DB
var dbSchema map[string]map[string]tableSchema = make(map[string]map[string]tableSchema, 1)

func Connect(hostname, port, database, user, password, cert string) {
	var connParam string
	if cert != "" {
		connParam = fmt.Sprintf("HOSTNAME=%s;DATABASE=%s;PORT=%s;UID=%s;PWD=%s;SECURITY=SSL;SSLSERVERCERTIFICATE=%s", hostname, database, port, user, password, cert)
	} else {
		connParam = fmt.Sprintf("HOSTNAME=%s;DATABASE=%s;PORT=%s;UID=%s;PWD=%s", hostname, database, port, user, password)
	}
	connStr := flag.String("conn", connParam, "connection string")
	db1, err := sqlx.Connect("go_ibm_db", *connStr)
	db = db1
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func Disconnect() {
	defer db.Close()
}

func ReadAllTableSchema() {
	schemaData := []tableSchemaRow{}
	err := db.Select(&schemaData, `
        Select c.tabschema as schema_name,
             c.tabname as table_name, 
             c.colname as column_name,
             c.colno as position,
             c.typename as data_type,
             c.length,
             c.scale,
             c.remarks as description,   
             case when  c.nulls = 'Y' then 1 else 0 end as nullable,
             default as default_value,
             case when c.identity ='Y' then 1 else 0 end as is_identity,
             case when c.generated ='' then 0 else 1 end as  is_computed,
             c.text as computed_formula
      from syscat.columns c
      inner join syscat.tables t on 
            t.tabschema = c.tabschema and t.tabname = c.tabname
      where t.type = 'T'
      order by schema_name,
               table_name
  `)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// for _, r := range schemaData {
	// fmt.Printf("%+v\n", r)
	// }
	for _, row := range schemaData {
		schemaName, tableName := row.Schema, row.Table
		// fmt.Printf("schema: %s, table:%s, row:%+v \n", schemaName, tableName, row)
		schema, ok := dbSchema[schemaName]
		if !ok {
			dbSchema[schemaName] = make(map[string]tableSchema, 1)
		}
		schema = dbSchema[schemaName]
		table, ok := schema[tableName]
		if !ok {
			schema[tableName] = tableSchema{Schema: schemaName, Table: tableName, Columns: []tableSchemaRow{}}
		}
		table = schema[tableName]
		table.Columns = append(table.Columns, row)
		dbSchema[schemaName][tableName] = table
	}
}

func PrintTableSchema() {
	for s, schema := range dbSchema {
		fmt.Println("Schema: \v\n", s)
		for t, tbl := range schema {
			fmt.Println("Table: \v\n", t)
			for _, col := range tbl.Columns {
				fmt.Printf("%+v\n", col)
			}
		}
	}
}

func PrintOneTableSchema(schema, table string) {
	if s, ok := dbSchema[schema]; ok {
		if t, ok := s[table]; ok {
			for _, r := range t.Columns {
				fmt.Printf("%+v\n", r)
			}
		}
	}
}
