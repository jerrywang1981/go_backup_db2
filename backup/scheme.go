package backup

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
)

type tableSchemaRow struct {
	Schema          string         `db:"SCHEMA_NAME"`
	Table           string         `db:"TABLE_NAME"`
	Column          string         `db:"COLUMN_NAME"`
	Position        int            `db:"POSITION"`
	Type            string         `db:"DATA_TYPE"`
	Length          int            `db:"LENGTH"`
	Scale           int            `db:"SCALE"`
	Description     sql.NullString `db:"DESCRIPTION"`
	Nullable        bool           `db:"NULLABLE"`
	DefaultValue    sql.NullString `db:"DEFAULT_VALUE"`
	IsIdentity      bool           `db:"IS_IDENTITY"`
	IsComputed      bool           `db:"IS_COMPUTED"`
	ComputedFormula sql.NullString `db:"COMPUTED_FORMULA"`
}

type tableSchema struct {
	Schema  string
	Table   string
	Columns []tableSchemaRow
}

func CreateTableSql(folder, schemaName, tableName string) (string, string) {
	absName, err := filepath.Abs(folder)
	if err != nil {
		fmt.Println("Invalid folder: ", err)
		return "", ""
	}
	schemaName, tableName = strings.ToUpper(schemaName), strings.ToUpper(tableName)
	if schema, ok := dbSchema[schemaName]; ok {
		if table, ok := schema[tableName]; ok {
			return table.createTableSql(absName)
		}
		return "", ""
	}
	return "", ""
}

func (ts *tableSchema) createTableSql(folder string) (string, string) {
	absName, err := filepath.Abs(folder)
	if err != nil {
		fmt.Println("Invalid folder: ", err)
		return "", ""
	}
	strCol := getColumnStrings(ts.Columns)
	expSql := fmt.Sprintf("export to %s/%s.%s.del of del select %s from %s.%s;", absName, ts.Schema, ts.Table, strCol, ts.Schema, ts.Table)
	impSql := fmt.Sprintf("load client from %s/%s.%s.del of del insert into %s.%s(%s) nonrecoverable", absName, ts.Schema, ts.Table, ts.Schema, ts.Table, strCol)
	return expSql, impSql
}

func getColumnStrings(colRows []tableSchemaRow) string {
	cols := []string{}
	for _, col := range colRows {
		cols = append(cols, col.Column)
	}
	strCol := strings.Join(cols, ",")
	return strCol
}

func GenerateSql(folder string, config map[string][]string) ([]string, []string) {
	absName, err := filepath.Abs(folder)
	if err != nil {
		fmt.Println("Invalid folder: ", err)
		return nil, nil
	}
	expSql, impSql := []string{}, []string{}
	expSql = append(expSql, "update command options using c off;", "select current timestamp from sysibm.sysdummy1;")

	impSql = append(impSql, "select current timestamp from sysibm.sysdummy1;")
	for schemaName, data := range config {
		for _, tableName := range data {
			eSql, iSql := CreateTableSql(absName, schemaName, tableName)
			expSql = append(expSql, eSql)
			impSql = append(impSql, iSql, "commit;")
		}
	}

	expSql = append(expSql, "select current timestamp from sysibm.sysdummy1;")
	expSql = append(expSql, "update command options using c on;")

	impSql = append(impSql, "select current timestamp from sysibm.sysdummy1;")
	return expSql, impSql
}