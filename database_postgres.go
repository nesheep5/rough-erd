package rough_er

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type postgres struct {
	db     *sql.DB
	config *PostgresConfig
}

func (p *postgres) Close() {
	p.db.Close()
}

type PostgresConfig struct {
	DatabaseSchemaConfig
	URL *url.URL
}

func (config *PostgresConfig) String() string {
	return config.URL.String()
}

func (config *PostgresConfig) Database() (Database, error) {
	db, err := sql.Open("postgres", config.String())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("PostgresConfig#Database: fail to open postgres connection. %s", err.Error()))
	}
	return &postgres{db: db, config: config}, nil
}

func (p *postgres) splitDBSchema(dBSchema string) (string, string) {
	var database string
	var schema string
	ds := strings.Split(dBSchema, ".")
	if len(ds) == 1 {
		if ds[0] == "postgres" {
			database = ds[0]
			schema = "public"
		} else if ds[0] == "public" {
			database = "postgres"
			schema = ds[0]
		} else {
			panic(fmt.Sprintf("parse error database and schema. %s", dBSchema))
		}
	} else {
		database = ds[0]
		schema = ds[1]
	}
	return database, schema
}

func (p *postgres) pgDump(schemaName string, tableName string) (string, error) {
	databaseName, schemaName := p.splitDBSchema(schemaName)
	args := []string{
		"--schema-only",
		fmt.Sprintf("--dbname=%s", databaseName),
		fmt.Sprintf("--schema=%s", schemaName),
		fmt.Sprintf("--table=%s", tableName),
	}
	u := p.config.URL
	if u.Hostname() != "" {
		args = append(args, fmt.Sprintf("--host=%s", u.Hostname()))
	}
	user := u.User
	if user.Username() != "" {
		args = append(args, fmt.Sprintf("--username=%s", user.Username()))
	}
	if u.Port() != "" {
		args = append(args, fmt.Sprintf("--port=%s", u.Port()))
	}

	cmd := exec.Command("pg_dump", args...)
	cmd.Env = os.Environ()

	password, set := user.Password()
	if set {
		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", password))
	}
	values := u.Query()
	supportEnv := map[string]string{"sslmode": "PGSSLMODE"}
	for key, _ := range values {
		if env, ok := supportEnv[key]; ok {
			args = append(args, fmt.Sprintf("%s=%s", env, values.Get(key)))
		}
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(fmt.Sprintf("%s\n%s", err.Error(), string(out)))
	}

	return string(out), nil
}

func (p *postgres) tableDefinition(schemaName string, tableName string) (string, error) {
	ddl, err := p.pgDump(schemaName, tableName)
	if err != nil {
		return "", err
	}
	patterns := []string{
		"^--.*$", // comments
		"^\\\\\\.$",
		"^SET .*;$",
		"^CREATE EXTENSION .*;$",
		"^SELECT .*;$",
		"^COPY .*;$",
		"^ALTER TABLE [^ ;]+ OWNER TO .+;$",
		"^REVOKE .*;$",
		"^GRANT .*;$",
		"^$\n", // empty line
	}

	for _, pattern := range patterns {
		re := regexp.MustCompilePOSIX(pattern)
		ddl = re.ReplaceAllLiteralString(ddl, "")
	}

	return ddl, nil
}

func (p *postgres) viewDefinition(schemaName string, tableName string) (string, error) {
	databaseName, schemaName := p.splitDBSchema(schemaName)
	SQL := fmt.Sprintf("SELECT view_definition FROM %s.information_schema.views WHERE table_schema = $1 AND table_name = $2", databaseName)
	var definition string
	err := p.db.QueryRow(SQL, schemaName, tableName).Scan(&definition)
	if err != nil {
		return "", err
	}
	definition = fmt.Sprintf("CREATE OR REPLACE VIEW %s AS %s", tableName, definition)
	return definition, nil
}

func (p *postgres) columns(schemaName string, tableName string) ([]*Column, error) {
	databaseName, schemaName := p.splitDBSchema(schemaName)
	SQL := fmt.Sprintf("SELECT ordinal_position, column_name, data_type, COALESCE(column_default, 'null'), is_nullable FROM %s.information_schema.COLUMNS WHERE table_schema = $1 AND table_name = $2 ORDER BY ordinal_position ASC", databaseName)
	rows, err := p.db.Query(SQL, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var columns []*Column
	for rows.Next() {
		var columnPosition uint
		var columnName string
		var columnType string
		var columnDefault string
		var columnNullable string
		if err := rows.Scan(&columnPosition, &columnName, &columnType, &columnDefault, &columnNullable); err != nil {
			return nil, err
		}
		column := Column{
			Position: columnPosition,
			Name:     columnName,
			Type:     columnType,
			Default:  columnDefault,
			Nullable: columnNullable,
		}
		columns = append(columns, &column)
	}

	comments, err := p.columnComments(schemaName, tableName)
	if err != nil {
		return nil, err
	}
	for pos, comment := range comments {
		for _, column := range columns {
			if column.Position == pos {
				column.Comment = comment
			}
		}
	}

	return columns, nil
}

func (p *postgres) tableComment(schemaName string, tableName string) (string, error) {
	databaseName, schemaName := p.splitDBSchema(schemaName)
	var comment string
	SQL := fmt.Sprintf("select COALESCE(obj_description(st.relid), '') from %s.pg_catalog.pg_statio_all_tables as st where st.schemaname = $1 and st.relname = $2", databaseName)
	err := p.db.QueryRow(SQL, schemaName, tableName).Scan(&comment)
	if err != nil {
		return "", err
	}
	return comment, nil
}

func (p *postgres) columnComments(schemaName string, tableName string) (map[uint]string, error) {
	databaseName, schemaName := p.splitDBSchema(schemaName)
	SQL := fmt.Sprintf("select col.ordinal_position, COALESCE(col_description(st.relid, col.ordinal_position), '') from %s.pg_catalog.pg_statio_all_tables as st inner join %s.information_schema.columns as col on col.table_schema = st.schemaname and col.table_name = st.relname where st.schemaname = $1 and st.relname = $2", databaseName, databaseName)
	rows, err := p.db.Query(SQL, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := map[uint]string{}

	for rows.Next() {
		var position uint
		var comment string
		if err := rows.Scan(&position, &comment); err != nil {
			return nil, err
		}
		comments[position] = comment
	}
	return comments, nil
}

func (p *postgres) tables(schemaName string) ([]*Table, error) {
	databaseName, schemaName := p.splitDBSchema(schemaName)
	SQL := fmt.Sprintf("SELECT table_name, table_type FROM %s.information_schema.tables WHERE table_schema = $1", databaseName)
	rows, err := p.db.Query(SQL, schemaName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*Table
	for rows.Next() {
		var tableName string
		var tableType string
		if err := rows.Scan(&tableName, &tableType); err != nil {
			return nil, err
		}
		table := Table{
			Name:    tableName,
			Type:    tableType,
			Comment: "",
		}
		table.Comment, err = p.tableComment(schemaName, table.Name)
		if err != nil {
			return nil, err
		}
		table.Columns, err = p.columns(schemaName, table.Name)
		if err != nil {
			return nil, err
		}
		if table.Type == "VIEW" {
			table.Definition, err = p.viewDefinition(schemaName, table.Name)
			if err != nil {
				return nil, err
			}
		} else {
			table.Definition, err = p.tableDefinition(schemaName, table.Name)
			if err != nil {
				return nil, err
			}
		}
		tables = append(tables, &table)
	}
	return tables, nil
}

func (p *postgres) Schemas() ([]*Schema, error) {
	var ss []*Schema
	for alias, name := range p.config.Names() {
		tables, err := p.tables(name)
		if err != nil {
			return nil, err
		}
		s := Schema{
			Name:   alias,
			Tables: tables,
		}
		ss = append(ss, &s)
	}
	return ss, nil
}
