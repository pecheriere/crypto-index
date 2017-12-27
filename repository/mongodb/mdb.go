package mongodb

import (
	"gopkg.in/mgo.v2"
	"crypto/tls"
	"net"
	"io/ioutil"
	"fmt"
	"strings"
	"log"
)

func GetSession() (*mgo.Session, error) {

	mdbHost, err := ioutil.ReadFile("/run/secrets/mdb_host")
	if err != nil {
		return nil, err
	}

	if string(mdbHost) == "localhost:27017\n" {
		log.Println("Connecting to MDB on localhost")
		return mgo.Dial("mongodb://localhost:27017")
	}

	mdbUser, err := ioutil.ReadFile("/run/secrets/mdb_user")
	if err != nil {
		return nil, err
	}

	mdbPassword, err := ioutil.ReadFile("/run/secrets/mdb_password")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s",
		strings.TrimSuffix(string(mdbUser), "\n"),
		strings.TrimSuffix(string(mdbPassword), "\n"),
		strings.TrimSuffix(string(mdbHost), "\n"))

	dialInfo, err := mgo.ParseURL(uri)

	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{}
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		return nil, err
	}

	return session, nil
}
