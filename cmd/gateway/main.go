package main

import (
    "fmt"
    "strconv"
    
    "github.com/PSauerborn/lifelink/pkg/gateway"
    "github.com/PSauerborn/lifelink/pkg/utils"
    log "github.com/sirupsen/logrus"
)

var cfg = utils.NewConfigMapWithValues(map[string]string{
    "listen_port": "8080",
    "listen_port_admin": "8081",
    "log_level": "DEBUG",
    "neo4j_host": "localhost",
    "neo4j_port": "7687",
    "neo4j_username": "neo4j",
    "neo4j_password": "development",
    "jwt_secret": "secret",
})

func runService() {
    
}

func main() {
    // configure log level
    cfg.ConfigureLogging()

    // retrieve API listen port and parse
    listenPort, err := strconv.Atoi(cfg.Get("listen_port"))
    if err != nil {
        panic(fmt.Errorf("invalid port %s", cfg.Get("listen_port")))
    }

    // retrieve admin listen port and parse
    listenPortAdmin, err := strconv.Atoi(cfg.Get("listen_port_admin"))
    if err != nil {
        panic(fmt.Errorf("invalid port %s", cfg.Get("listen_port_admin")))
    }

    // retrieve port for neo4j and parse to 
    neo4jPort, err := strconv.Atoi(cfg.Get("neo4j_port"))
    if err != nil {
        panic(fmt.Errorf("invalid port %s", cfg.Get("neo4j_port")))
    }

    // set new graph peristence layer and defer closing
    persistence := gateway.SetGraphPersistence(cfg.Get("neo4j_host"), 
        neo4jPort, cfg.Get("neo4j_username"), cfg.Get("neo4j_password"))
    defer persistence.Driver.Close()

    // generate new admin API and run on admin port
    admin := gateway.NewGatewayAdminAPI(cfg.Get("jwt_secret"), 180)
    go func() {
        defer func() {
            if r := recover(); r != nil {
                log.Warn(fmt.Sprintf("recovered paniced admin API: %+v", r))
                admin.Run(fmt.Sprintf(":%d", listenPortAdmin))
            }
        }()
        admin.Run(fmt.Sprintf(":%d", listenPortAdmin))
    }()
    // generate new instance of API and run
    gateway.NewAPIGateway(cfg.Get("jwt_secret")).Run(fmt.Sprintf(":%d", listenPort))
}