package mongo

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/globalsign/mgo"
)

type DBOption struct {
	Host     []string
	Port     []string
	Username string
	Password string
	DBName   string
}

// Connect function need to be update
// can not support for multiple server
func Connect(option DBOption) (*mgo.Session, error) {
	//connect URL:
	// "mongodb://<username>:<password>@<hostname>:<port>,<hostname>:<port>/<db-name>
	var host string
	for _, v := range option.Host {
		tmpHost = fmt.Sprintf("%s:%s")
		host = fmt.Sprintf("%s,%s", host, tmpHost)
	}

	url := fmt.Sprintf("mongodb://%s:%s@%s/%s", option.Username, option.Password, host, option.DBName)
	dialInfo, err := mgo.ParseURL(url)
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	//Here it is the session. Up to you from here ;)
	session, err := mgo.DialWithInfo(dialInfo)
	return session, err
}
