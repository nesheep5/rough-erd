package rough_erd

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"strings"
)

type Database interface {
	Schemas() ([]*Schema, error)
	Close()
}

type DatabaseConfig interface {
	String() string
	Database() (Database, error)
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
	Socket   string
}

func (config *DatabaseSchemaConfig) Names() map[string]string {
	var names = map[string]string{}
	for _, name := range config.Schemas {
		var alias string
		tmp := strings.Split(name, "@")
		if len(tmp) > 1 {
			name = tmp[0]
			alias = tmp[1]

		} else {
			alias = name
		}
		names[alias] = name
	}
	return names
}

func MakeDBConfig(typec ConnectInfo) DatabaseConfig {

	return nil
}
