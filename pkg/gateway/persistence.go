package gateway

import (
	"fmt"
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	log "github.com/sirupsen/logrus"

	"github.com/PSauerborn/lifelink/pkg/utils"
)

var (
	// define custom errors 
	ErrInvalidModule = errors.New("Module does not exist")
	ErrModuleExists  = errors.New("Module already exists")
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


type Module struct {
	ModuleName        string `json:"module_name" validate:"required"`
	ModuleRedirect    string `json:"module_redirect" validate:"required"`
	TrimAppName       bool   `json:"trim_app_name" validate:"required"`
	ModuleDescription string `json:"module_description" validate:"required"`
}

// function used to retrieve module details from graph
func(db *GraphPersistence) GetModuleDetails(name string) (Module, error) {
	log.Debug(fmt.Sprintf("fetching module details for %s...", name))
 
	// create new persitence session for graph and defer closing
	session := db.NewSession()
	defer session.Close()
	// generate config metadata for query
	cfg := map[string]interface{}{
		"module_name": name,
	}
	// define handle used to execute graph function
	handler := func(tx neo4j.Transaction) (interface{}, error) {
		query := `MATCH (n:Module{module_name: $module_name})
		RETURN n.module_name, n.module_redirect, n.module_description, 
		n.trim_app_name LIMIT 1`

		result, err := tx.Run(query, cfg)
		if err != nil {
			log.Error(fmt.Errorf("unable to execute graph query: %+v", err))
			return nil, err
		}
		// retrieve first instance of modules. note that invalid module
		// error is returned here if the module cannot be found
		node, err := neo4j.Single(result, err)
		if err != nil {
			return nil, ErrInvalidModule
		}
		return node, err
	}
	// get module details from graph using persistence session
	result, err := neo4j.AsRecord(session.ReadTransaction(handler))
	if err != nil {
		log.Error(fmt.Errorf("unable to get module details: %+v", err))
		return Module{}, err
	}
	return Module{
		ModuleName: result.Values[0].(string),
		ModuleRedirect: result.Values[1].(string),
		ModuleDescription: result.Values[2].(string),
		TrimAppName: result.Values[3].(bool),
	}, nil
}

// function used to add a new module to the graph. note that
// all module lables are constrained so that the module name
// must be unique to avoid duplicate modules
func(db *GraphPersistence) AddModule(module Module) error {
	log.Debug(fmt.Sprintf("generating new module %+v...", module))

	// create new persitence session for graph and defer closing
	session := db.NewSession()
	defer session.Close()
	// generate config metadata for query
	cfg := map[string]interface{}{
		"module_name": module.ModuleName,
		"module_redirect": module.ModuleRedirect,
		"module_description": module.ModuleDescription,
		"trim_app_name": module.TrimAppName,
	}
	// define handle used to execute graph function
	handler := func(tx neo4j.Transaction) (interface{}, error) {
		query := `CREATE (n:Module{
			module_name: $module_name,
			module_redirect: $module_redirect,
			module_description: $module_description,
			trim_app_name: $trim_app_name
		})`
		return tx.Run(query, cfg)
	}
	// get module details from graph using persistence session
	_, err := session.WriteTransaction(handler)
	if err != nil {
		log.Error(fmt.Errorf("unable to generate new module: %+v", err))
		return err
	}
	return nil
}
