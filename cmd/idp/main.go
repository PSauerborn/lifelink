package main

import (
    "fmt"
    "strconv"

    "github.com/PSauerborn/lifelink/pkg/idp"
    "github.com/PSauerborn/lifelink/pkg/utils"
)

var cfg = utils.NewConfigMapWithValues(map[string]string{
    "listen_port": "10867",
    "log_level": "DEBUG",
    "neo4j_host": "localhost",
    "neo4j_port": "7687",
    "neo4j_username": "neo4j",
    "neo4j_password": "development",
    "admin_api_host": "localhost",
    "admin_api_port": "8081",
    "users_api_host": "localhost",
    "users_api_port": "10866",
})

// funcion used to retrieve downstream microservice
// configuration/connection settings for admin API
func getAdminAPIConfig() utils.APIDependencyConfig {
    apiPort, err := strconv.Atoi(cfg.Get("admin_api_port"))
    if err != nil {
        panic("received invalid api port for admin API")
    }
    return utils.APIDependencyConfig{
        Host: cfg.Get("admin_api_host"),
        Port: &apiPort,
        Protocol: "http",
    }
}

// funcion used to retrieve downstream microservice
// configuration/connection settings for users API
func getUsersAPIConfig() utils.APIDependencyConfig {
    apiPort, err := strconv.Atoi(cfg.Get("users_api_port"))
    if err != nil {
        panic(("received invalid api port for users API"))
    }
    return utils.APIDependencyConfig{
        Host: cfg.Get("users_api_host"),
        Port: &apiPort,
        Protocol: "http",
    }
}

func main() {
    // configure log level
    cfg.ConfigureLogging()

    // retrieve API listen port and parse
    listenPort, err := strconv.Atoi(cfg.Get("listen_port"))
    if err != nil {
        panic(fmt.Errorf("invalid port %s", cfg.Get("listen_port")))
    }

    // retrieve port for neo4j and parse to integer
    neo4jPort, err := strconv.Atoi(cfg.Get("neo4j_port"))
    if err != nil {
        panic(fmt.Errorf("invalid port %s", cfg.Get("neo4j_port")))
    }

    // set new graph peristence layer and defer closing
    persistence := idp.SetGraphPersistence(cfg.Get("neo4j_host"),
        neo4jPort, cfg.Get("neo4j_username"), cfg.Get("neo4j_password"))
    defer persistence.Driver.Close()

    // generate new instance of API (with config for users and admin API's)
    service := idp.NewIdentityProvider(getUsersAPIConfig(), getAdminAPIConfig())
    service.Run(fmt.Sprintf(":%d", listenPort))
}