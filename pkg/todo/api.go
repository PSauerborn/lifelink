package todo

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/PSauerborn/lifelink/pkg/utils"
)

var persistence *GraphPersistence

// function used to generate new TODO api
func NewTodoAPI() *gin.Engine {
	router := gin.Default()
	// add middleware to inject user ID into request context
	router.Use(utils.UserInjectionMiddleware())

	router.GET("/TODO/health_check", healthCheckHandler)
	router.GET("/TODO/items", getTodoItemsHandler)
	router.PATCH("/TODO/item/:itemId", completeTodoItemHandler)
	router.POST("/TODO/new", newTodoItemHandler)
	router.DELETE("/TODO/item/:itemId", deleteTodoItemHandler)
	return router
}

// function used to set new instance of graph persistence
// for global variables to use
func SetGraphPersistence(host string, port int,
	username, password string) *GraphPersistence {
	// generate new graph persistence and set globally
	db, err := NewGraphPersistence(fmt.Sprintf("neo4j://%s", host),
		port, username, password)
	persistence = db
	if err != nil {
		panic(fmt.Errorf("unable to generate graph persistence: %+v", err))
	}
	return persistence
}

// API handler used to serve health check routes
func healthCheckHandler(ctx *gin.Context) {
	log.Info("received request for health check handler")
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"success": true, "message": "Running"})
}

// API handler to generate a new TODO item for a given user
func newTodoItemHandler(ctx *gin.Context) {
	log.Info("received request for new todo item")
	var request TODOItem
	if err := ctx.ShouldBind(&request); err != nil {
		log.Error(fmt.Errorf("received invalid request: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"success": false, "message": "Invalid request body"})
		return
	}

	// retrieve user ID from context
	uid := ctx.MustGet("uid").(string)
	if err := persistence.CreateTodoItem(uid, request); err != nil {
		log.Error(fmt.Errorf("unable to create TODO item: %+v", err))
		switch err {
		case ErrInvalidTODOMetadata:
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
				"success": false, "message": "Invalid item metadata"})
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"http_code": http.StatusInternalServerError,
				"success": false, "message": "Internal server error"})
		}
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"http_code": http.StatusCreated,
		"success": false, "message": "Successfully created TODO item"})
}

// API handler to retrieve TODO item(s) for a given user
func getTodoItemsHandler(ctx *gin.Context) {
	log.Info("received request to retrieve todo item(s)")

	// retrieve user ID from context
	uid := ctx.MustGet("uid").(string)
	items, err := persistence.GetTodoItems(uid)
	if err != nil {
		log.Error(fmt.Errorf("unable to create TODO item: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"http_code": http.StatusInternalServerError,
			"success": false, "message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"success": true, "items": items})
}

// API handler to complete a TODO item for a given user
func completeTodoItemHandler(ctx *gin.Context) {
	log.Info("received request to complete todo")

	itemId, err := uuid.Parse(ctx.Param("itemId"))
	if err != nil {
		log.Error(fmt.Errorf("unable to parse item ID: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"success": false, "message": "Invalid item ID"})
		return
	}
	// retrieve user ID from context
	uid := ctx.MustGet("uid").(string)
	if err := persistence.CompleteTodoItem(uid, itemId); err != nil {
		log.Error(fmt.Errorf("unable to complete todo item: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"http_code": http.StatusInternalServerError,
			"success": false, "message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"success": true, "message": "Successfully compelted TODO item"})
}

// API handler to delete a TODO item for a given user
func deleteTodoItemHandler(ctx *gin.Context) {
	log.Info("received request to delete todo item")

	itemId, err := uuid.Parse(ctx.Param("itemId"))
	if err != nil {
		log.Error(fmt.Errorf("unable to parse item ID: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"http_code": http.StatusBadRequest,
			"success": false, "message": "Invalid item ID"})
		return
	}
	// retrieve user ID from context
	uid := ctx.MustGet("uid").(string)
	if err := persistence.DeleteTodoItem(uid, itemId); err != nil {
		log.Error(fmt.Errorf("unable to delete todo item: %+v", err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"http_code": http.StatusInternalServerError,
			"success": false, "message": "Internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
		"success": true, "message": "Successfully deleted TODO item"})
}
