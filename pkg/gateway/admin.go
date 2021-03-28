package gateway

import (
	"fmt"
	"time"
	"net/http"

    "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var (
	jwtSecret string
	tokenExpiryMinutes int
)

// function used to generate new API gateway admin service
func NewGatewayAdminAPI(secret string, tokenExpiry int) *gin.Engine {
	// set variables to be used globally
	jwtSecret, tokenExpiryMinutes = secret, tokenExpiry 

	router := gin.Default()
	router.GET("/admin/health_check", healthCheckHandler)
	router.POST("/admin/token", getTokenHandler)
	return router
}

// function used to generate JWToken with UID and expiry date
func generateJWToken(uid string, admin bool) (string, error) {
    // evaluate expiry time
    expiry := time.Now().UTC()
    expiry = expiry.Add(time.Duration(tokenExpiryMinutes) * time.Minute)
    // generate token and sign with secret key
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "uid": uid,
        "exp": expiry.Unix(),
        "admin": admin,
    })
    return token.SignedString([]byte(jwtSecret))
}

// handler used to serve health check routes
func healthCheckHandler(ctx *gin.Context) {
	log.Info("received request for health check handler")
	ctx.JSON(http.StatusCreated, gin.H{"http_code": http.StatusCreated,
		"success": true, "message": "Running"})
}


// API handler used to generate and retreive new
// JWT token
func getTokenHandler(ctx *gin.Context) {
	log.Info("received request for token")
    var request struct {
        Uid   string `json:"uid"   binding:"required"`
        Admin *bool  `json:"admin" binding:"required"`
    }
    // parse request body from context
    if err := ctx.ShouldBind(&request); err != nil {
        log.Error(fmt.Errorf("received invalid request body: %+v", err))
        ctx.JSON(http.StatusBadRequest, gin.H{
            "http_code": http.StatusBadRequest, "message": "Invalid request body"})
        return
    }
    // generate JWToken for user and return
    token, err := generateJWToken(request.Uid, *request.Admin)
    if err != nil {
        log.Error(fmt.Errorf("unable to generate JWToken: %+v", err))
        ctx.JSON(http.StatusInternalServerError, gin.H{"http_code": http.StatusInternalServerError, 
			"success": false, "message": "Internal server error"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK, 
		"success": true, "token": token})
}