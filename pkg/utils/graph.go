package utils

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	log "github.com/sirupsen/logrus"
)

// struct used to store base graph connection settings
type BaseGraphAccessor struct {
	// define connection settings for graph
	GraphHost string 
	GraphPort int

	// define properties for driver and session
	Driver neo4j.Driver
}

// function used to generate new graph accessor from connection
// settings and credentials
func NewGraphAccessor(host string, port int,
	username, password string) (*BaseGraphAccessor, error) {
	creds := neo4j.BasicAuth(username, password, "")
	// generate new driver instance
	driver, err := neo4j.NewDriver(fmt.Sprintf("%s:%d", host, port), creds)
	if err != nil {
		log.Error(fmt.Errorf("unable to generate neo4j driver: %+v", err))
		return nil, err
	}
	return &BaseGraphAccessor{
		GraphHost: host, 
		GraphPort: port,
		Driver: driver,
	}, nil
}

// function used to generate new driver and session 
// for graph queries
func(accessor *BaseGraphAccessor) NewSession() neo4j.Session {
	// generate new session and set attributes
	session := accessor.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite})
	return session
}