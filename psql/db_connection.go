package psql

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

var sslModes []string = []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}

// DBOptions wrap options database
type DBOptions struct {
	Host           string
	Port           int
	Username       string
	Password       string
	DBName         string
	ConnectTimeout int
	SSLCert        string
	SSLKey         string
	SSLRootCert    string
	SSLMode        string

	DataSourceName string
	DB             *sql.DB
}

// DSN format connection
func (dbOption *DBOptions) DSN() string {
	sslMode := "sslmode=disable"

	if options.SSLMode != "" && options.SSLMode != "disable" {
		if !isValidSSLMode(options.SSLMode) {
			return nil, errors.New("lib-go: invalid ssl mode")
		}

		sslMode = fmt.Sprintf("sslmode=%s&sslrootcert=%s&sslcert=%s&sslkey=%s",
			options.SSLMode,
			options.SSLRootCert,
			options.SSLCert,
			options.SSLKey)
	}

	dbConfig := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?%s",
		options.Username,
		options.Password,
		options.Host,
		options.Port,
		options.DBName,
		sslMode,
	)

	dbOption.DataSourceName = dbConfig
}

func isValidSSLMode(sslMode string) bool {
	for _, v := range sslModes {
		if sslMode == v {
			return true
		}
	}

	return false
}

// Connect is connection to database
func Connect(options *DBOptions) (*sql.DB, error) {
	options.DSN()

	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// ConnectConfig is connection with return config
func ConnectConfig(options DBOptions) (*DBOptions, error) {
	options.DSN()

	db, err := sql.Open("postgres", dbConfig)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	options.DB = db
	return &options, nil
}
