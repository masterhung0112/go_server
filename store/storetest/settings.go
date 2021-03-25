package storetest

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"

	"database/sql"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/masterhung0112/hk_server/v5/model"
)

const (
	defaultMysqlDSN        = "hkuser:mostest@tcp(localhost:7306)/hungknow_test?charset=utf8mb4,utf8\u0026readTimeout=30s\u0026writeTimeout=30s&multiStatements=true"
	defaultPostgresqlDSN   = "postgres://hkuser:mostest@localhost:7432/hungknow_test?sslmode=disable&connect_timeout=10"
	defaultMysqlRootPWD    = "mostest"
	defaultMysqlReplicaDSN = "root:mostest@tcp(localhost:3307)/mattermost_test?charset=utf8mb4,utf8\u0026readTimeout=30s"
)

func getEnv(name, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}

func log(message string) {
	verbose := false
	if verboseFlag := flag.Lookup("test.v"); verboseFlag != nil {
		verbose = verboseFlag.Value.String() != ""
	}
	if verboseFlag := flag.Lookup("v"); verboseFlag != nil {
		verbose = verboseFlag.Value.String() != ""
	}

	if verbose {
		fmt.Println(message)
	}
}

// MySQLSettings returns the database settings to connect to the MySQL unittesting database.
// The database name is generated randomly and must be created before use.
func MySQLSettings(withReplica bool) *model.SqlSettings {
	dsn := getEnv("TEST_DATABASE_MYSQL_DSN", defaultMysqlDSN)
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		panic("failed to parse dsn " + dsn + ": " + err.Error())
	}

	cfg.DBName = "db" + model.NewId()

	mySQLSettings := databaseSettings("mysql", cfg.FormatDSN())

	if withReplica {
		mySQLSettings.DataSourceReplicas = []string{getEnv("TEST_DATABASE_MYSQL_REPLICA_DSN", defaultMysqlReplicaDSN)}
	}

	return mySQLSettings
}

// PostgresSQLSettings returns the database settings to connect to the PostgreSQL unittesting database.
// The database name is generated randomly and must be created before use.
func PostgreSQLSettings() *model.SqlSettings {
	dsn := getEnv("TEST_DATABASE_POSTGRESQL_DSN", defaultPostgresqlDSN)
	dsnUrl, err := url.Parse(dsn)
	if err != nil {
		panic("failed to parse dsn " + dsn + ": " + err.Error())
	}

	// Generate a random database name
	dsnUrl.Path = "db" + model.NewId()

	return databaseSettings("postgres", dsnUrl.String())
}

func mySQLRootDSN(dsn string) string {
	rootPwd := getEnv("TEST_DATABASE_MYSQL_ROOT_PASSWD", defaultMysqlRootPWD)
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		panic("failed to parse dsn " + dsn + ": " + err.Error())
	}

	cfg.User = "root"
	cfg.Passwd = rootPwd
	cfg.DBName = "mysql"

	return cfg.FormatDSN()
}

func postgreSQLRootDSN(dsn string) string {
	dsnUrl, err := url.Parse(dsn)
	if err != nil {
		panic("failed to parse dsn " + dsn + ": " + err.Error())
	}

	// // Assume the unittesting database has the same password.
	// password := ""
	// if dsnUrl.User != nil {
	// 	password, _ = dsnUrl.User.Password()
	// }

	// dsnUrl.User = url.UserPassword("", password)
	dsnUrl.Path = "postgres"

	return dsnUrl.String()
}

func mySQLDSNDatabase(dsn string) string {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		panic("failed to parse dsn " + dsn + ": " + err.Error())
	}

	return cfg.DBName
}

func postgreSQLDSNDatabase(dsn string) string {
	dsnUrl, err := url.Parse(dsn)
	if err != nil {
		panic("failed to parse dsn " + dsn + ": " + err.Error())
	}

	return path.Base(dsnUrl.Path)
}

func databaseSettings(driver, dataSource string) *model.SqlSettings {
	settings := &model.SqlSettings{
		DriverName:                  &driver,
		DataSource:                  &dataSource,
		DataSourceReplicas:          []string{},
		DataSourceSearchReplicas:    []string{},
		MaxIdleConns:                new(int),
		ConnMaxLifetimeMilliseconds: new(int),
		ConnMaxIdleTimeMilliseconds: new(int),
		MaxOpenConns:                new(int),
		Trace:                       model.NewBool(false),
		AtRestEncryptKey:            model.NewString(model.NewRandomString(32)),
		QueryTimeout:                new(int),
	}
	*settings.MaxIdleConns = 10
	*settings.ConnMaxLifetimeMilliseconds = 3600000
	*settings.ConnMaxIdleTimeMilliseconds = 300000
	*settings.MaxOpenConns = 100
	*settings.QueryTimeout = 60

	return settings
}

// execAsRoot executes the given sql as root against the testing database
func execAsRoot(settings *model.SqlSettings, sqlCommand string) error {
	var dsn string
	var driver = *settings.DriverName

	switch driver {
	case model.DATABASE_DRIVER_MYSQL:
		dsn = mySQLRootDSN(*settings.DataSource)
	case model.DATABASE_DRIVER_POSTGRES:
		dsn = postgreSQLRootDSN(*settings.DataSource)
	default:
		return fmt.Errorf("unsupported driver %s", driver)
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to %s database as root", driver)
	}
	defer db.Close()
	if _, err = db.Exec(sqlCommand); err != nil {
		return errors.Wrapf(err, "failed to execute `%s` against %s database as root", sqlCommand, driver)
	}

	return nil
}

func replaceMySQLDatabaseName(dsn, newDBName string) string {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		panic("failed to parse dsn " + dsn + ": " + err.Error())
	}
	cfg.DBName = newDBName
	return cfg.FormatDSN()
}

// MakeSqlSettings creates a randomly named database and returns the corresponding sql settings
func MakeSqlSettings(driver string, withReplica bool) *model.SqlSettings {
	var settings *model.SqlSettings
	var dbName string

	switch driver {
	case model.DATABASE_DRIVER_MYSQL:
		settings = MySQLSettings(withReplica)
		dbName = mySQLDSNDatabase(*settings.DataSource)
		newDSRs := []string{}
		for _, dataSource := range settings.DataSourceReplicas {
			newDSRs = append(newDSRs, replaceMySQLDatabaseName(dataSource, dbName))
		}
		settings.DataSourceReplicas = newDSRs
	case model.DATABASE_DRIVER_POSTGRES:
		settings = PostgreSQLSettings()
		dbName = postgreSQLDSNDatabase(*settings.DataSource)
	default:
		panic("unsupported driver " + driver)
	}

	if err := execAsRoot(settings, "CREATE DATABASE "+dbName); err != nil {
		panic("failed to create temporary database " + dbName + ": " + err.Error())
	}

	switch driver {
	case model.DATABASE_DRIVER_MYSQL:
		if err := execAsRoot(settings, "GRANT ALL PRIVILEGES ON "+dbName+".* TO 'hkuser'"); err != nil {
			panic("failed to grant hkuser permission to " + dbName + ":" + err.Error())
		}
	case model.DATABASE_DRIVER_POSTGRES:
		if err := execAsRoot(settings, "GRANT ALL PRIVILEGES ON DATABASE \""+dbName+"\" TO hkuser"); err != nil {
			panic("failed to grant hkuser permission to " + dbName + ":" + err.Error())
		}
	default:
		panic("unsupported driver " + driver)
	}

	log("Created temporary " + driver + " database " + dbName)

	return settings
}

func CleanupSqlSettings(settings *model.SqlSettings) {
	var driver = *settings.DriverName
	var dbName string

	switch driver {
	case model.DATABASE_DRIVER_MYSQL:
		dbName = mySQLDSNDatabase(*settings.DataSource)
	case model.DATABASE_DRIVER_POSTGRES:
		dbName = postgreSQLDSNDatabase(*settings.DataSource)
	default:
		panic("unsupported driver " + driver)
	}

	if err := execAsRoot(settings, "DROP DATABASE "+dbName); err != nil {
		panic("failed to drop temporary database " + dbName + ": " + err.Error())
	}

	log("Dropped temporary database " + dbName)
}
