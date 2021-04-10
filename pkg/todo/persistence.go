package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	log "github.com/sirupsen/logrus"

	"github.com/PSauerborn/lifelink/pkg/utils"
)

var (
	// define custom errors
	ErrTODOItemNotFound    = errors.New("cannot find specified TODO item")
	ErrInvalidTODOMetadata = errors.New("invalid TODO metadata")
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

type TODOItem struct {
	ItemId      uuid.UUID              `json:"item_id"`
	ItemTitle   string                 `json:"item_title" binding:"required"`
	ItemContent string                 `json:"item_content" binding:"required"`
	Created     time.Time              `json:"created"`
	Completed   bool                   `json:"completed"`
	Metadata    map[string]interface{} `json:"metadata" binding:"required"`
}

// function used to create a new todo item for a given user
func (db *GraphPersistence) CreateTodoItem(uid string, item TODOItem) error {
	log.Debug(fmt.Sprintf("creating new TODO item for user %s", uid))
	session := db.NewSession()
	defer session.Close()

	// convert metadata to JSON string and store
	meta, err := json.Marshal(item.Metadata)
	if err != nil {
		log.Error(fmt.Errorf("unable to convert toto metadata: %+v", err))
		return ErrInvalidTODOMetadata
	}
	cfg := map[string]interface{}{
		"uid":          uid,
		"item_id":      uuid.New().String(),
		"item_title":   item.ItemTitle,
		"item_content": item.ItemContent,
		"metadata":     string(meta),
	}
	// define handler to execute node query
	handler := func(tx neo4j.Transaction) (interface{}, error) {
		query := `CREATE (t:TODO {
            item_id: $item_id,
            item_title: $item_title,
            item_content: $item_content,
            created: datetime(),
            completed: false,
            metadata: $metadata
        })
        WITH t
        MATCH (u:User {uid: $uid})
        CREATE (u)-[:OWNS]->(t)`
		return tx.Run(query, cfg)
	}
	// execute handler within write transaction
	_, err = session.WriteTransaction(handler)
	if err != nil {
		log.Error(fmt.Errorf("unable to execute graph query: %+v", err))
		return err
	}
	return nil
}

// function used to retrieve todo items for a given user
func (db *GraphPersistence) GetTodoItem(uid string, itemId uuid.UUID) (TODOItem, error) {
	log.Debug(fmt.Sprintf("retrieving TODO item %s for user %s", itemId, uid))
	session := db.NewSession()
	defer session.Close()

	cfg := map[string]interface{}{
		"uid":     uid,
		"item_id": itemId.String(),
	}
	handler := func(tx neo4j.Transaction) (interface{}, error) {
		query := `MATCH (u:User {uid: $uid})-[:OWNS]->(t:TODO {item_id: $item_id})
        RETURN t.item_id, t.item_title, t.item_content, t.completed, t.created, t.metadata
        `
		// run graph query to fetch TODO item
		results, err := tx.Run(query, cfg)
		if err != nil {
			return nil, err
		}
		// retrieve single entry; return not found it err is
		// not nil
		node, err := neo4j.Single(results, err)
		if err != nil {
			return nil, ErrTODOItemNotFound
		}
		return node, nil
	}
	// execute read transaction with handler to retreive item
	node, err := neo4j.AsRecord(session.ReadTransaction(handler))
	if err != nil {
		return TODOItem{}, err
	}

	var metadata map[string]interface{}
	// convert metadata from JSON string to struct
	json.Unmarshal([]byte(node.Values[5].(string)), &metadata)
	return TODOItem{
		ItemId:      node.Values[0].(uuid.UUID),
		ItemTitle:   node.Values[1].(string),
		ItemContent: node.Values[2].(string),
		Completed:   node.Values[3].(bool),
		Created:     node.Values[4].(time.Time),
		Metadata:    metadata,
	}, nil
}

// function used to retrieve a todo item with given todo ID
func (db *GraphPersistence) GetTodoItems(uid string) ([]TODOItem, error) {
	log.Debug(fmt.Sprintf("retreiving TODO item(s) for user %s", uid))
	session := db.NewSession()
	defer session.Close()

	items := []TODOItem{}
	cfg := map[string]interface{}{
		"uid": uid,
	}
	handler := func(tx neo4j.Transaction) (interface{}, error) {
		query := `MATCH (u:User {uid: $uid})-[:OWNS]->(t:TODO)
        RETURN t.item_id, t.item_title, t.item_content, t.completed,
        t.created, t.metadata`
		return neo4j.Collect(tx.Run(query, cfg))
	}
	nodes, err := neo4j.AsRecords(session.ReadTransaction(handler))
	if err != nil {
		log.Error(fmt.Errorf("unable to execute graph query: %+v", err))
		return items, err
	}
	for _, node := range nodes {
		itemId, _ := uuid.Parse(node.Values[0].(string))
		var metadata map[string]interface{}
		// convert metadata from JSON string to struct
		json.Unmarshal([]byte(node.Values[5].(string)), &metadata)
		items = append(items, TODOItem{
			ItemId:      itemId,
			ItemTitle:   node.Values[1].(string),
			ItemContent: node.Values[2].(string),
			Completed:   node.Values[3].(bool),
			Created:     node.Values[4].(time.Time),
			Metadata:    metadata,
		})
	}
	return items, nil
}

// function used to complete a given todo item
func (db *GraphPersistence) CompleteTodoItem(uid string, itemId uuid.UUID) error {
	log.Debug(fmt.Sprintf("completing TODO item %s for user %s", itemId, uid))
	session := db.NewSession()
	defer session.Close()

	cfg := map[string]interface{}{
		"uid":     uid,
		"item_id": itemId.String(),
	}
	handler := func(tx neo4j.Transaction) (interface{}, error) {
		query := `MATCH (u:User {uid: $uid})-[:OWNS]->(t:TODO {item_id: $item_id})
        SET t.completed = true`
		return tx.Run(query, cfg)
	}
	_, err := session.WriteTransaction(handler)
	if err != nil {
		log.Error(fmt.Errorf("unable to execute graph query: %+v", err))
		return err
	}
	return nil
}

// function used to create a new todo item for a given user
func (db *GraphPersistence) DeleteTodoItem(uid string, itemId uuid.UUID) error {
	log.Debug(fmt.Sprintf("deleting TODO item %s for user %s", itemId, uid))
	session := db.NewSession()
	defer session.Close()

	cfg := map[string]interface{}{
		"uid":     uid,
		"item_id": itemId.String(),
	}
	handler := func(tx neo4j.Transaction) (interface{}, error) {
		query := `MATCH (u:User {uid: $uid})-[:OWNS]->(t:TODO {item_id: $item_id})
        DETACH DELETE t`
		return tx.Run(query, cfg)
	}
	_, err := session.WriteTransaction(handler)
	if err != nil {
		log.Error(fmt.Errorf("unable to execute graph query: %+v", err))
		return err
	}
	return nil
}
