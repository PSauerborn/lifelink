package idp

import (
    "fmt"
    "errors"

    "github.com/neo4j/neo4j-go-driver/v4/neo4j"
    log "github.com/sirupsen/logrus"

    "github.com/PSauerborn/lifelink/pkg/utils"
)

var (
    // define custom errors 
    ErrUserDoesNotExist = errors.New("User does not exist")
)

type GraphPersistence struct {
    *utils.BaseGraphAccessor
}

// function used to generate a new graph persistence layer
func NewGraphPersistence(host string, port int,
    username, password string) (*GraphPersistence, error) {
    // generate new base accessor and return persistence layers
    baseAccessor, err := utils.NewGraphAccessor(host, port, username, password)
    if err != nil {
        log.Error(fmt.Errorf("unable to generate base graph connection: %+v", err))
        return nil, err
    }
    return &GraphPersistence{
        BaseGraphAccessor: baseAccessor,
    }, nil
}

// function used to set credentials for a given user
func(db *GraphPersistence) SetUserCredentials(uid, password string) error {
    log.Debug(fmt.Sprintf("setting credentials for user %s", uid))
    session := db.NewSession()
    defer session.Close()
    
    cfg := map[string]interface{}{
        "uid": uid,
        "password": hashAndSalt(password),
    }
    // generate handler function to process graph query
    handler := func(tx neo4j.Transaction) (interface{}, error) {
        query := `MATCH (u:User {uid: $uid})-[:OWNS]->(c:Credentials) 
        SET c.password = $password`
        return tx.Run(query, cfg)
    }

    _, err := session.WriteTransaction(handler)
    if err != nil {
        log.Error(fmt.Errorf("unable to update user credentials: %+v", err))
        return err
    }
    return nil
}

// function used to set credentials for a given user
func(db *GraphPersistence) AddUserCredentials(uid, password string) error {
    log.Debug(fmt.Sprintf("adding credentials for user %s", uid))
    session := db.NewSession()
    defer session.Close()
    
    cfg := map[string]interface{}{
        "uid": uid,
        "password": hashAndSalt(password),
    }
    // generate handler function to process graph query
    handler := func(tx neo4j.Transaction) (interface{}, error) {
        query := `CREATE (c:Credentials {password: $password})
        WITH c
        MATCH 
            (u:User {uid: $uid})
        CREATE (u)-[:OWNS]->(h)` 
        return tx.Run(query, cfg)
    }
    _, err := session.WriteTransaction(handler)
    if err != nil {
        log.Error(fmt.Errorf("unable to update user credentials: %+v", err))
        return err
    }
    return nil
}

// function used to set credentials for a given user
func(db *GraphPersistence) GetUserCredentials(uid string) (string, error) {
    log.Debug(fmt.Sprintf("fetching credentials for user %s", uid))
    session := db.NewSession()
    defer session.Close()
    
    cfg := map[string]interface{}{
        "uid": uid,
    }
    // generate handler function to process graph query
    handler := func(tx neo4j.Transaction) (interface{}, error) {
        query := `MATCH (u:User {uid: $uid}-[:OWNS]->(c:Credentials)
        RETURN c.password` 
        result, err := tx.Run(query, cfg)
        if err != nil {
            log.Error(fmt.Errorf("unable to retrieve user credentials: %+v", err))
            return nil, err
        }
        // get first result from node query and return
        node, err := neo4j.Single(result, err)
        if err != nil {
            log.Error(fmt.Errorf("unable to retrieve user credentials: %+v", err))
            return nil, ErrUserDoesNotExist
        }
        return node, nil
    }
    // get credentials node from graph and return
    node, err := neo4j.AsRecord(session.ReadTransaction(handler))
    if err != nil {
        log.Error(fmt.Errorf("unable to get user credentials: %+v", err))
        return "", err
    }
    return node.Values[0].(string), nil
}