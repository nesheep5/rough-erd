package rough_er

const (
	DataSourceMySQL      = "mysql"
	DataSourcePostgreSQL = "postgres"
)

type Option struct {
	SchemaNames []string
	Format      string
	Path        string
	DataSource  string
}

func MakeDoc(databaseConfig DatabaseConfig, option Option) error {
	var err error
	database, err := databaseConfig.Database()
	if err != nil {
		return err
	}
	defer database.Close()

	var schemas []*Schema
	schemas, err = database.Schemas()
	if err != nil {
		return err
	}

	return Render(option.DataSource, option.Format, option.Path, schemas)
}
