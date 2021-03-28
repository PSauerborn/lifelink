package habits

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"

    "github.com/PSauerborn/lifelink/pkg/utils"
)

var persistence *GraphPersistence

// function used to generate new API 
func NewHabitsAPI() *gin.Engine {
    router := gin.Default()
    // add middleware to inject user ID into request context
    router.Use(utils.UserInjectionMiddleware())

    router.GET("/habits/health_check", healthCheckHandler)
    router.GET("/habits/all", getHabitsHandler)

    router.POST("/habits/new", createHabitHandler)

    router.PUT("/habits/update/:habitId", updateHabitHandler)
    router.PATCH("/habits/complete/:habitId", completeHabitHandler)
    router.DELETE("/habits/delete/:habitId", deleteHabitHandler)
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
    log.Info(fmt.Sprintf("received request for health check handler"))
    ctx.JSON(http.StatusCreated, gin.H{"http_code": http.StatusCreated,
        "success": true, "message": "Running"})
}

// handler used to retrieve user habits
func getHabitsHandler(ctx *gin.Context) {
    log.Info(fmt.Sprintf("received request to retrieve habits"))
    // retrieve user ID from context
    uid := ctx.MustGet("uid").(string)

    // get habits from graph for given user
    habits, err := persistence.GetUserHabits(uid)
    if err != nil {
        log.Error(fmt.Errorf("unable to retreive user habits: %+v", err))
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "http_code": http.StatusInternalServerError, "success": false, 
            "message": "Internal server error"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
        "success": true, "habits": habits})
}

// function used to generate a new habit for given user
func createHabitHandler(ctx *gin.Context) {
    log.Info("received request to create new habit")
    var request Habit 
    if err := ctx.ShouldBind(&request); err != nil {
        log.Error(fmt.Sprintf("received invalid request body"))
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "http_code": http.StatusBadRequest, "success": false,
            "message": "Invalid request body"})
        return
    }

    // check that given habit cycle is valid
    if !isValidCycle(request.HabitCycle) {
        log.Error(fmt.Sprintf("received invalid cycle %s", request.HabitCycle))
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "http_code": http.StatusBadRequest, "success": false,
            "message": "Invalid habit cycle"})
        return
    }

    // retrieve user ID from context
    uid := ctx.MustGet("uid").(string)
    if err := persistence.CreateUserHabit(uid, request); err != nil {
        log.Error(fmt.Errorf("unable to create new habit for user %s: %+v", uid, err))
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "http_code": http.StatusInternalServerError, "success": false, 
            "message": "Internal server error"})
        return
    }
    ctx.JSON(http.StatusCreated, gin.H{"http_code": http.StatusCreated,
        "success": true, "message": "Successfully created habit"})
}

// function used to generate a new habit for given user
func updateHabitHandler(ctx *gin.Context) {
    log.Info("received request to update habit")
    var request Habit 
    // parse habit ID from request path
    habitId, err := uuid.Parse(ctx.Param("habitId"))
    if err != nil {
        log.Error(fmt.Sprintf("received invalid habit ID %s", ctx.Param("habitId")))
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "http_code": http.StatusBadRequest, "success": false,
            "message": "Invalid habit ID"})
        return
    }
    // parse request body
    if err := ctx.ShouldBind(&request); err != nil {
        log.Error(fmt.Sprintf("received invalid request body"))
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "http_code": http.StatusBadRequest, "success": false,
            "message": "Invalid request body"})
        return
    }

    // check that given habit cycle is valid
    if !isValidCycle(request.HabitCycle) {
        log.Error(fmt.Sprintf("received invalid cycle %s", request.HabitCycle))
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "http_code": http.StatusBadRequest, "success": false,
            "message": "Invalid habit cycle"})
        return
    }

    // retrieve user ID from context
    uid := ctx.MustGet("uid").(string)
    if err := persistence.UpdateUserHabit(uid, habitId, request); err != nil {
        log.Error(fmt.Errorf("unable to update habit for user %s: %+v", uid, err))
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "http_code": http.StatusInternalServerError, "success": false, 
            "message": "Internal server error"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
        "success": true, "message": "Successfully updated habit"})
}

// function used to complete a habit with given habit ID
func completeHabitHandler(ctx *gin.Context) {
    log.Info("received request to complete habit")
    habitId, err := uuid.Parse(ctx.Param("habitId"))
    if err != nil {
        log.Error(fmt.Sprintf("received invalid habit ID %s", ctx.Param("habitId")))
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "http_code": http.StatusBadRequest, "success": false,
            "message": "Invalid habit ID"})
        return
    }

    // retrieve user ID from context
    uid := ctx.MustGet("uid").(string)
    if err := persistence.CompleteUserHabit(uid, habitId); err != nil {
        log.Error(fmt.Errorf("unable to complete user habit: %+v", err))
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "http_code": http.StatusInternalServerError, "success": false, 
            "message": "Internal server error"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
        "success": true, "message": "Successfully completed habit"})
}

// function used to delete a habit with given habit ID
func deleteHabitHandler(ctx *gin.Context) {
    log.Info("received request to delete habit")
    habitId, err := uuid.Parse(ctx.Param("habitId"))
    if err != nil {
        log.Error(fmt.Sprintf("received invalid habit ID %s", ctx.Param("habitId")))
        ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
            "http_code": http.StatusBadRequest, "success": false,
            "message": "Invalid habit ID"})
        return
    }

    // retrieve user ID from context
    uid := ctx.MustGet("uid").(string)
    if err := persistence.DeleteUserHabit(uid, habitId); err != nil {
        log.Error(fmt.Errorf("unable to delete user habit: %+v", err))
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
            "http_code": http.StatusInternalServerError, "success": false, 
            "message": "Internal server error"})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"http_code": http.StatusOK,
        "success": true, "message": "Successfully deleted habit"})
}