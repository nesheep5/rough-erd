package rough_erd

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type Database interface {
	Tables(schemaName string) ([]*Table, error)
	Close()
}

type DatabaseConfig interface {
	String() string
	Names() map[string]string
}

type DatabaseSchemaConfig struct {
	Schemas []string
}

type ConnectInfo struct {
	User     string
	Password string
	Port     int
	Protocol string
	DBName   string
}

func CreateDatabase(dbType string, conn *ConnectInfo) (Database, error) {
	switch dbType {
	case "mysql":
		return MysqlDatabase(conn)
	}
	return nil, fmt.Errorf("dbtype is invalid. type: %s", dbType)
}
