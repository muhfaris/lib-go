package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// DBOptions is wrap options database
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
}

// Connect is connection to mysql
func Connect(options DBOptions) (*sql.DB, error) {
	dbConfig := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		options.Username,
		options.Password,
		options.Host,
		options.Port,
		options.DBName,
	)

	db, err := sql.Open("mysql", dbConfig)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
