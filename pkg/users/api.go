package users

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"

    "github.com/PSauerborn/lifelink/pkg/utils"
)

var persistence *GraphPersistence

func NewUsersAPI() *gin.Engine {
    router := gin.Default()
    // add middleware to inject user ID into request context
    router.Use(utils.UserInjectionMiddleware())

    router.GET("/users/health_check", healthCheckHandler)
    router.GET("/users/user", getUserHandler)
    router.GET("/users/details/:uid", AdminProtected(), getUserDetailsHandler)
    router.POST("/users/new", AdminProtected(), createUserHandler)
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

// handler used to serve health check routes
func healthCheckHandler(ctx *gin.Context) {
    log.Info("received request for health check handler")
    ctx.JSON(http.StatusCreated, gin.H{"http_code": http.StatusCreated,
        "success": true, "message": "Running"})
}

// API handler used to retrieve user from graph
func getUserHandler(ctx *gin.Context) {
    log.Info("received request to get user")
    // retrieve user ID from context
    uid := ctx.MustGet("uid").(string)

    // get user details from graph
    user, err := persistence.GetUserDetails(uid)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve user details: %+v", err))
        switch err {
        case ErrUserDoesNotExist:
            ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
                "http_code": http.StatusNotFound, "success": false,
                "message": "Cannot find user"})
        default:
            ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                "http_code": http.StatusInternalServerError, "success": false,
                "message": "Internal server error"})
        }
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
        "success": true, "user": user})
}

// API handler used to retrieve user details from graph
func getUserDetailsHandler(ctx *gin.Context) {
    log.Info("received request to get user details")
    // retrieve user ID from context
    targetUser := ctx.Param("uid")

    // get user details from graph
    user, err := persistence.GetUserDetails(targetUser)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve user details: %+v", err))
        switch err {
        case ErrUserDoesNotExist:
            ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
                "http_code": http.StatusNotFound, "success": false,
                "message": "Cannot find user"})
        default:
            ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                "http_code": http.StatusInternalServerError, "success": false,
                "message": "Internal server error"})
        }
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
        "success": true, "user": user})
}

// API handler used to generate a new user
func createUserHandler(ctx *gin.Context) {
    log.Info("received request to generate new user")
    var request User
    if err := ctx.ShouldBind(&request); err != nil {
        log.Error(fmt.Errorf("received invalid request: %+v", err))
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "http_code": http.StatusBadRequest, "success": false,
            "message": "Invalid request body"})
        return
    }

    // try to get user details from database to check if user exists
    _, err := persistence.GetUserDetails(request.Uid)
    if err != ErrUserDoesNotExist {
        log.Error(fmt.Errorf("cannot create entry for user %s: already exists", request.Uid))
        ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
            "http_code": http.StatusConflict, "success": false,
            "message": "Username already exists"})
        return
    }

    // add new user to persistence graph
    if err := persistence.AddUser(request); err != nil {
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "http_code": http.StatusInternalServerError, "success": false,
            "message": "Internal server error"})
        return
    }
    ctx.JSON(http.StatusCreated, gin.H{"http_code": http.StatusCreated,
        "success": true, "message": "Successfully created user"})
}