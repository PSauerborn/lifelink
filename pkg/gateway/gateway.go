package gateway

import (
    "fmt"
    "strings"
    "net/http"
    "net/url"
    "net/http/httputil"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var persistence *GraphPersistence

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

// function used to generate new API gateway service
func NewAPIGateway(jwtSecret string) *gin.Engine {
	router := gin.Default()
	// add JWT middleware to parse access tokens
	router.Use(JWTMiddleware(jwtSecret, false))
	router.Any("/api/:application/*proxyPath", proxyHandler)
	return router
}

// API handler used to proxy request to downstream microservices
func proxyHandler(ctx *gin.Context) {
	uid := ctx.MustGet("uid").(string)
	log.Debug(fmt.Sprintf("proxying request for user %s", uid))
	// inject user ID into downstream headers
    ctx.Request.Header.Set("X-Authenticated-Userid", uid)
	// get module from graph persistence and handle errors
	module, err := persistence.GetModuleDetails(ctx.Param("application"))
	if err != nil {
		switch err {
		case ErrInvalidModule:
			log.Error(fmt.Sprintf("unable to retrieve module details: invalid module %s", 
				ctx.Param("application")))
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"http_code": http.StatusBadGateway, "success": false,
				"message": "Bad Gateway"})
		default:
			log.Error(fmt.Errorf("unable to retrieve module details: %+v", err))
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"http_code": http.StatusInternalServerError, "success": false,
				"message": "Internal server error"})
		}
		return
	}
    // proxy request to relevant microservices
    proxyRequest(module, ctx.Writer, ctx.Request)
}

// function used to set proxy headers headers on request
func SetProxyHeaders(request *http.Request, url *url.URL) {
    request.URL.Host = url.Host
    request.URL.Scheme = url.Scheme
    request.Header.Set("X-Forwarded-Host", request.Header.Get("Host"))
    request.Host = url.Host
}

// define function used to proxy request
func proxyRequest(app Module, response http.ResponseWriter, request *http.Request) {
    var redirectUrl string
    // trim app name from URL if specified in app config
    if app.TrimAppName {
        log.Debug(fmt.Sprintf("trimming app name from redirect for application %s", 
			app.ModuleName))
        replace := fmt.Sprintf("/%s", app.ModuleName)
        redirectUrl = strings.Replace(app.ModuleRedirect, replace, "", -1)
    } else {
        redirectUrl = app.ModuleRedirect
    }

	log.Info(fmt.Sprintf("proxying request to %s", redirectUrl))
    // construct new URL, set proxy headers and proxy
    redirect, _ := url.Parse(redirectUrl)
    SetProxyHeaders(request, redirect)
    // create reverse proxy instance and serve request
    proxy := httputil.NewSingleHostReverseProxy(redirect)
    proxy.ServeHTTP(response, request)
}