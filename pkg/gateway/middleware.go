package gateway

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
)


// middleware used to parse JWTokens from request
func JWTMiddleware(jwtSecret string, adminOnly bool) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // set relevant cors headers and return options calls
        SetCorsHeaders(ctx.Writer, ctx.Request)
        if ctx.Request.Method == http.MethodOptions {
            log.Debug("received options calls. returning...")
            ctx.AbortWithStatus(http.StatusOK)
            return
        }

        log.Debug(fmt.Sprintf("received request for URL %s", ctx.Request.URL.Path))
        // authenticate user using JWToken present in request
        claims, err := authenticateUser(ctx.Request, jwtSecret)
        if err != nil {
            log.Error(fmt.Errorf("unable to authenticate user: %v", err))
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"http_code": http.StatusUnauthorized, "success": false, 
				"message": "Unauthorized"})
            return
        }
        // enforce admin-only uses on admin restricted routes
        if adminOnly && !claims.Admin {
            log.Error(fmt.Errorf("unable to authenticate user: %v", err))
            ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"http_code": http.StatusForbidden, "success": false,
				"message": "Forbidden"})
            return
        }

        // inject uid into request context
        log.Info(fmt.Sprintf("received proxy request for user %s", claims.Uid))
        ctx.Set("uid", claims.Uid)
        ctx.Next()
    }
}