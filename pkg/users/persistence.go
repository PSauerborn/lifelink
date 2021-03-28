package users

import (
    "fmt"
    "time"
    "errors"

    "github.com/neo4j/neo4j-go-driver/v4/neo4j"
    log "github.com/sirupsen/logrus"

    "github.com/PSauerborn/lifelink/pkg/utils"
)

var (
    // define custom errors 
    ErrUserDoesNotExist  = errors.New("User does not exist")
    ErrUserAlreadyExists = errors.New("User already exists")
)

type GraphPersistence struct {
    *utils.BaseGraphAccessor
}

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

type User struct {
    Uid     string    `json:"uid" binding:"required"`
    Email   string    `json:"email" binding:"required"`
    Created time.Time `json:"created"`
    Admin   bool      `json:"admin"`
}

// function used to retrieve all users from the graph
func(db *GraphPersistence) GetAllUsers() ([]User, error) {
    log.Debug("retreiving all users...")
    users := []User{}

    // create new persitence session for graph and defer closing
    session := db.NewSession()
    defer session.Close()
    // generate config metadata for query
    handler := func(tx neo4j.Transaction) (interface{}, error) {
        query := `MATCH (n:User) 
        RETURN n.uid, n.email, n.created, n.admin`
        return neo4j.Collect(tx.Run(query, nil))
    }
    // get all habits from graph using persistence session
    nodes, err := neo4j.AsRecords(session.ReadTransaction(handler))
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve user details: %+v", err))
        return users, err
    }
    for _, node := range(nodes) {
        users = append(users, User{
            Uid: node.Values[0].(string),
            Email: node.Values[1].(string),
            Created: node.Values[2].(time.Time),
            Admin: node.Values[3].(bool),
        })
    }
    return users, nil
}

// function used to retrieve all users from the graph
func(db *GraphPersistence) GetUserDetails(uid string) (User, error) {
    log.Debug(fmt.Sprintf("retreiving details for user %s...", uid))
    // create new persitence session for graph and defer closing
    session := db.NewSession()
    defer session.Close()

    cfg := map[string]interface{}{
        "uid": uid,
    }
    // generate config metadata for query
    handler := func(tx neo4j.Transaction) (interface{}, error) {
        query := `MATCH (n:User{uid: $uid})
        RETURN n.uid, n.email, n.created, n.admin`
        result, err := tx.Run(query, cfg)
        if err != nil {
            return nil, err
        }
        // retrieve first instance of user node
        node, err := neo4j.Single(result, err)
        if err != nil {
            return nil, ErrUserDoesNotExist
        }
        return node, nil
    }
    // get all habits from graph using persistence session
    user, err := neo4j.AsRecord(session.ReadTransaction(handler))
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve user details: %+v", err))
        return User{}, err
    }
    return User{
        Uid: user.Values[0].(string),
        Email: user.Values[1].(string),
        Created: user.Values[2].(time.Time),
        Admin: user.Values[3].(bool),
    }, nil
}

// function used to retrieve all users from the graph
func(db *GraphPersistence) AddUser(user User) error {
    log.Debug(fmt.Sprintf("creating new user %+v...", user))
    // create new persitence session for graph and defer closing
    session := db.NewSession()
    defer session.Close()

    cfg := map[string]interface{}{
        "uid": user.Uid,
        "email": user.Email,
        "created": time.Now().UTC(),
        "admin": user.Admin,
    }
    // generate config metadata for query
    handler := func(tx neo4j.Transaction) (interface{}, error) {
        query := `CREATE (n:User{
            uid: $uid,
            email: $email,
            created: $created,
            admin: $admin
        }) RETURN n.uid`
        return tx.Run(query, cfg)
    }
    // get all habits from graph using persistence session
    _, err := session.WriteTransaction(handler)
    if err != nil {
        log.Error(fmt.Errorf("unable to generate new user: %+v", err))
        switch err.(type) {
        case *neo4j.Neo4jError:
            return ErrUserAlreadyExists
        default:
            return err
        }
    }
    return nil
}
