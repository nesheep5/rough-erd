package rough_erd

import (
	"database/sql"
	"errors"
	"fmt"

	mysql_driver "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	db     *sql.DB
	config *MySQLConfig
}

func (m *Mysql) Close() {
	m.db.Close()
}

type MySQLConfig struct {
	Config *mysql_driver.Config
}

func MysqlDatabase(c *ConnectInfo) (Database, error) {
	config := &MySQLConfig{
		Config: &mysql_driver.Config{
			User:                 c.User,
			Passwd:               c.Password,
			Net:                  c.Protocol,
			Addr:                 fmt.Sprintf("%s:%d", c.Host, c.Port),
			DBName:               c.DBName,
			AllowNativePasswords: true,
		},
	}
	con := config.Config.FormatDSN()
	// fmt.Printf("  connection param: %s\n", con)
	db, err := sql.Open("mysql", con)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("MySQLConfig#Database: fail to open pmysql connection. %s", err.Error()))
	}
	return &Mysql{db: db, config: config}, nil
}

func (m *Mysql) columns(schemaName string, tableName string) ([]*Column, error) {
	q := `SELECT ORDINAL_POSITION, COLUMN_NAME, COLUMN_TYPE, COALESCE(COLUMN_DEFAULT, 'null'), IS_NULLABLE, COLUMN_COMMENT
         FROM information_schema.COLUMNS
         WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ?
         ORDER BY ORDINAL_POSITION`
	rows, err := m.db.Query(q, schemaName, tableName)

	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var columns []*Column
	for rows.Next() {
		var columnPosition uint
		var columnName string
		var columnType string
		var columnDefault string
		var columnNullable string
		var columnComment string
		if err := rows.Scan(&columnPosition, &columnName, &columnType, &columnDefault, &columnNullable, &columnComment); err != nil {
			return nil, err
		}
		column := Column{
			Position: columnPosition,
			Name:     columnName,
			Type:     columnType,
			Default:  columnDefault,
			Nullable: columnNullable,
			Comment:  columnComment,
		}
		columns = append(columns, &column)
	}
	return columns, nil
}

func (m *Mysql) Tables(schemaName string) ([]*Table, error) {
	rows, err := m.db.Query("SELECT TABLE_NAME, TABLE_TYPE, TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA = ?", schemaName)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var tables []*Table
	for rows.Next() {
		var err error
		var tableName string
		var tableType string
		var tableComment string
		if err = rows.Scan(&tableName, &tableType, &tableComment); err != nil {
			return nil, err
		}
		table := Table{
			Name:    tableName,
			Type:    tableType,
			Comment: tableComment,
		}
		table.Columns, err = m.columns(schemaName, table.Name)
		if err != nil {
			return nil, err
		}
		tables = append(tables, &table)
	}
	return tables, nil
}

//
//func (m *mysql) Schemas() ([]*Schema, error) {
//	var ss []*Schema
//	for alias, name := range m.config.Names() {
//		tables, err := m.tables(name)
//		if err != nil {
//			return nil, err
//		}
//		s := Schema{
//			Name:   alias,
//			Tables: tables,
//		}
//		ss = append(ss, &s)
//	}
//	return ss, nil
//}
