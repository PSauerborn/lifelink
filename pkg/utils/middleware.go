package utils

import (
    "net/http"

    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
)

// gin-gonic middleware used to check incoming request
// for user ID headers
func UserInjectionMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        uid := ctx.Request.Header.Get("X-Authenticated-Userid")
        if len(uid) > 0 {
            ctx.Set("uid", uid)
            ctx.Next()
        } else {
            log.Warn("unable to extract user ID from request header")
            ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "http_code": http.StatusUnauthorized, "success": false, 
                "message": "Unauthorized"})
        }
    }
}