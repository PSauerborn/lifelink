package idp

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    log "github.com/sirupsen/logrus"

    "github.com/PSauerborn/lifelink/pkg/utils"
    api "github.com/PSauerborn/lifelink/pkg/utils/accessors"
)

var (
    // define global persistence layer
    persistence *GraphPersistence

    // define global accessors for microservice mesh
    usersAPIAccessor *api.UsersAPIAccessor
    adminAPIAccessor *api.GatewayAdminAPIAccessor
)

// function used to generate new instance of idP
// service. note that accessors for both the users
//  and the API gateway admin console are generated
func NewIdentityProvider(usersCfg utils.APIDependencyConfig,
    adminCfg utils.APIDependencyConfig) *gin.Engine {

    // create API accessors for users API and for gateway
    usersAPIAccessor = api.NewUsersApiAccessorFromConfig(usersCfg)
    adminAPIAccessor = api.NewGatewayAdminApiAccessorFromConfig(adminCfg)

    router := gin.Default()
    router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

    router.GET("/authenticate/health_check", healthCheckHandler)
    router.POST("/authenticate/register", registerHandler)
    router.POST("/authenticate/token", authenticateHandler)
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
    log.Info(fmt.Sprintf("received request for health check handler"))
    ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
        "success": true, "message": "Running"})
}

// API handler used to register a new user
func registerHandler(ctx *gin.Context) {
    log.Info("received request to register new user")
    var request struct {
        Uid      string `json:"uid" binding:"required"`
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    if err := ctx.ShouldBind(&request); err != nil {
        log.Error(fmt.Errorf("received invalid authentication request: %+v", err))
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "http_code": http.StatusBadRequest, "success": false,
            "message": "Invalid request body"})
        return
    }

    // get user details from users API to get admin status
    success, err := usersAPIAccessor.CreateUser("lifelink_idp", request)
    if err != nil || !success {
        log.Error(fmt.Errorf("unable to fetch user details: %+v", err))
        switch err {
        case api.ErrUserAlreadyExists:
            ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
                "http_code": http.StatusBadRequest, "success": false,
                "message": "User already exists"})
        default:
            ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                "http_code": http.StatusInternalServerError, "success": false,
                "message": "Internal server error"})
        }
        return
    }

    // add node for user credentials to graph
    if err := persistence.AddUserCredentials(request.Uid, request.Password); err != nil {
        log.Error(fmt.Errorf("unable to add user credentials: %+v", err))
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "http_code": http.StatusInternalServerError, "success": false,
            "message": "Internal server error"})
        return
    }
    ctx.JSON(http.StatusCreated, gin.H{"http_code": http.StatusCreated,
        "success": true, "message": "successfully created user"})
}

// API handler used to authenticate users
func authenticateHandler(ctx *gin.Context) {
    log.Info("received token request")
    var request struct {
        Uid      string `json:"uid" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    if err := ctx.ShouldBind(&request); err != nil {
        log.Error(fmt.Errorf("received invalid authentication request: %+v", err))
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "http_code": http.StatusBadRequest, "success": false,
            "message": "Invalid request body"})
        return
    }
    // get hashed password from database
    creds, err := persistence.GetUserCredentials(request.Uid)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve user credentials: %+v", err))
        switch err {
        case ErrUserDoesNotExist:
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "http_code": http.StatusUnauthorized, "success": false,
                "message": "Unauthorized"})
        default:
            ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                "http_code": http.StatusInternalServerError, "success": false,
                "message": "Internal server error"})
        }
        return
    }
    // compare given password to hashed password
    if !comparePasswords(request.Password, creds) {
        ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
            "http_code": http.StatusUnauthorized, "success": false,
            "message": "Unauthorized"})
        return
    }

    // get user details from users API to get admin status
    details, err := usersAPIAccessor.GetUserDetails("lifelink_idp", request.Uid)
    if err != nil {
        log.Error(fmt.Errorf("unable to fetch user details: %+v", err))
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "http_code": http.StatusInternalServerError, "success": false,
            "message": "Internal server error"})
        return
    }
    // get token from API gateway
    token, err := adminAPIAccessor.GetAccessToken(request.Uid, details.User.Admin)
    if err != nil {
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "http_code": http.StatusInternalServerError, "success": false,
            "message": "Internal server error"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
        "success": true, "token": token.Token})
}