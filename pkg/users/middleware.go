package users

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
)

// middleware used to protect routes with admin-only
// access	
func AdminProtected() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // get user ID from header and get admin access
        uid := ctx.Request.Header.Get("X-Authenticated-Userid")
        admin, err := isAdminUser(uid)
        if err != nil {
            log.Error(fmt.Errorf("unable to check admin status for user: %+v", err))
            ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                "http_code": http.StatusInternalServerError, "success": false, 
                "message": "Internal server error"})
        }
        
        // return 403 if user does not have admin rights
        if !admin {
            ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
                "http_code": http.StatusForbidden, "success": false,
                "message": "Forbidden"})
        } else {
            ctx.Next()
        }
    }
}