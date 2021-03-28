package main

import (
	"fmt"
	"strconv"
	
	"github.com/PSauerborn/lifelink/pkg/users"
	"github.com/PSauerborn/lifelink/pkg/utils"
)

var cfg = utils.NewConfigMapWithValues(map[string]string{
	"listen_port": "10866",
	"log_level": "DEBUG",
	"neo4j_host": "localhost",
	"neo4j_port": "7687",
	"neo4j_username": "neo4j",
	"neo4j_password": "development",
})

func main() {
	// configure log level
	cfg.ConfigureLogging()

	// retrieve API listen port and parse
	listenPort, err := strconv.Atoi(cfg.Get("listen_port"))
	if err != nil {
		panic(fmt.Errorf("invalid port %s", cfg.Get("listen_port")))
	}

	// retrieve port for neo4j and parse to 
	neo4jPort, err := strconv.Atoi(cfg.Get("neo4j_port"))
	if err != nil {
		panic(fmt.Errorf("invalid port %s", cfg.Get("neo4j_port")))
	}

	// set new graph peristence layer and defer closing
	persistence := users.SetGraphPersistence(cfg.Get("neo4j_host"), 
		neo4jPort, cfg.Get("neo4j_username"), cfg.Get("neo4j_password"))
	defer persistence.Driver.Close()
	// generate new instance of API and run
	users.NewUsersAPI().Run(fmt.Sprintf(":%d", listenPort))
}