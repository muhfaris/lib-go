package psqlx

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
	DB             *pgx.Conn
	Pool           *pgxpool.Pool
}

// DSN is data source connection
func (options *DBOptions) DSN() error {
	sslMode := "sslmode=disable"

	if options.SSLMode != "" && options.SSLMode != "disable" {
		if !isValidSSLMode(options.SSLMode) {
			return errors.New("lib-go: invalid ssl mode")
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

	options.DataSourceName = dbConfig
	return nil
}

// isValidSSLMode is validate ssl mode
func isValidSSLMode(sslMode string) bool {
	for _, v := range sslModes {
		if sslMode == v {
			return true
		}
	}

	return false
}

// Connect is connection to database
func Connect(options *DBOptions, logger pgx.Logger) (*pgx.Conn, error) {
	var ctx = context.Background()
	if err := options.DSN(); err != nil {
		return nil, err
	}

	dbConfig, err := pgx.ParseConfig(options.DataSourceName)
	if err != nil {
		return nil, err
	}

	dbConfig.Logger = logger

	conn, err := pgx.ConnectConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	options.DB = conn
	return conn, nil
}

// ConnectPool is connection to database
func ConnectPool(options *DBOptions, logger pgx.Logger) (*pgxpool.Pool, error) {
	var ctx = context.Background()
	if err := options.DSN(); err != nil {
		return nil, err
	}

	dbConfig, err := pgxpool.ParseConfig(options.DataSourceName)
	if err != nil {
		return nil, err
	}

	dbConfig.ConnConfig.Logger = logger

	conn, err := pgxpool.ConnectConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	options.Pool = conn
	return conn, nil
}
