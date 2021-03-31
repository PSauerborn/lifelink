package habits

import (
    "fmt"
    "time"
    "errors"

    "github.com/google/uuid"
    "github.com/neo4j/neo4j-go-driver/v4/neo4j"
    log "github.com/sirupsen/logrus"

    "github.com/PSauerborn/lifelink/pkg/utils"
)

var (
    // define custom errors
    ErrInvalidModule     = errors.New("Module does not exist")
    ErrModuleExists      = errors.New("Module already exists")
    ErrHabitDoesNotExist = errors.New("Habit does not exist")
)

type GraphPersistence struct {
    *utils.BaseGraphAccessor
}

// function used to generate a new graph persistence layer
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


type Habit struct {
    HabitName        string     `json:"habit_name" binding:"required"`
    HabitId          uuid.UUID  `json:"habit_id"`
    Created          time.Time  `json:"created"`
    HabitDescription string     `json:"habit_description" binding:"required"`
    HabitCycle       string     `json:"habit_cycle" binding:"required"`
    LastCompleted    *time.Time `json:"last_completed"`
}

type HabitCompletion struct {
    OnTarget       bool      `json:"on_target"`
    Streak         bool      `json:"streak"`
    EventTimestamp time.Time `json:"event_timestamp"`
}

// function used to retrieve all habits for given user
func(db *GraphPersistence) GetUserHabits(user string) ([]Habit, error) {
    log.Debug(fmt.Sprintf("retrieving all habits for user %s...", user))
    habits := []Habit{}

    // create new persitence session for graph and defer closing
    session := db.NewSession()
    defer session.Close()
    // generate config metadata for query
    cfg := map[string]interface{}{
        "uid": user,
    }
    handler := func(tx neo4j.Transaction) (interface{}, error) {
        query := `MATCH (:User {uid: $uid})-[:OWNS]->(h:Habit)
        RETURN h.habit_name, h.habit_id, h.habit_description,
        h.habit_cycle, h.created, h.last_completed`
        return neo4j.Collect(tx.Run(query, cfg))
    }
    // get all habits from graph using persistence session
    nodes, err := neo4j.AsRecords(session.ReadTransaction(handler))
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve nodes from graph: %+v", err))
        return habits, err
    }

    for _, node := range(nodes) {
        // parse habit ID to UUID type
        habitId, _ := uuid.Parse(node.Values[1].(string))
        lastCompleted := node.Values[5]
        habit := Habit{
            HabitName: node.Values[0].(string),
            HabitId: habitId,
            HabitDescription: node.Values[2].(string),
            HabitCycle: node.Values[3].(string),
            Created: node.Values[4].(time.Time),
        }
        // add last completed date if set else leave as null
        if lastCompleted != nil {
            completed := lastCompleted.(time.Time)
            habit.LastCompleted = &completed
        } else {
            habit.LastCompleted = nil
        }
        habits = append(habits, habit)
    }
    return habits, nil
}

// function used to retrieve habit by habit ID
func(db *GraphPersistence) GetHabitByHabitId(user string, habitId uuid.UUID) (Habit, error) {
    log.Debug(fmt.Sprintf("retrieving habit %s for user %s", user, habitId))
    // create new persistence session for graph and defer closing
    session := db.NewSession()
    defer session.Close()
    // generate config metadata for query
    cfg := map[string]interface{}{
        "uid": user,
        "habit_id": habitId.String(),
    }
    // define handler function used to retrieve node data
    handler := func(tx neo4j.Transaction) (interface{}, error) {
        query := `MATCH (:User {uid: $uid})-[:OWNS]->(h:Habit {habit_id: $habit_id})
        RETURN h.habit_name, h.habit_id, h.habit_description,
        h.habit_cycle, h.created, h.last_completed`
        results, err := tx.Run(query, cfg)
        if err != nil {
            log.Error(fmt.Errorf("unable to retrieve data from graph: %+v", err))
            return nil, err
        }
        node, err := neo4j.Single(results, err)
        if err != nil {
            log.Error(fmt.Errorf("unable to parse single node: %+v", err))
            return nil, ErrHabitDoesNotExist
        }
        return node, nil
    }
    // get all habits from graph using persistence session
    node, err := neo4j.AsRecord(session.ReadTransaction(handler))
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve habit: %+v", err))
        return Habit{}, err
    }

    // construct new habbit from node
    habit := Habit{
        HabitId: habitId,
        HabitName: node.Values[0].(string),
        HabitDescription: node.Values[2].(string),
        HabitCycle: node.Values[3].(string),
        Created: node.Values[4].(time.Time),
    }
    // parse last completed pointer if not nill
    lastCompleted := node.Values[5]
    // add last completed date if set else leave as null
    if lastCompleted != nil {
        completed := lastCompleted.(time.Time)
        habit.LastCompleted = &completed
    } else {
        habit.LastCompleted = nil
    }
    return habit, nil
}

// function used to generate a new habbit for a given user
func(db *GraphPersistence) CreateUserHabit(user string, habit Habit) error {
    log.Debug(fmt.Sprintf("creating new habit %+v for user %s...", habit, user))
    // create new persitence session for graph and defer closing
    session := db.NewSession()
    defer session.Close()
    // generate config metadata for query
    cfg := map[string]interface{}{
        "habit_name": habit.HabitName,
        "habit_id": uuid.New().String(),
        "habit_description": habit.HabitDescription,
        "habit_cycle": orderCyclesString(habit.HabitCycle),
        "created": time.Now().UTC(),
        "uid": user,
    }
    handler := func(tx neo4j.Transaction) (interface{}, error) {
        query := `CREATE (h:Habit {
            habit_name: $habit_name,
            habit_id: $habit_id,
            habit_description: $habit_description,
            habit_cycle: $habit_cycle,
            last_completed: null,
            created: $created
        })
        WITH h
        MATCH
            (u:User {uid: $uid})
        CREATE (u)-[:OWNS]->(h)`
        return tx.Run(query, cfg)
    }
    // get all habits from graph using persistence session
    _, err := session.WriteTransaction(handler)
    if err != nil {
        log.Error(fmt.Errorf("unable to create hew habit: %+v", err))
        return err
    }
    return nil
}

// function used to complete a habit with given habit ID for user
func(db *GraphPersistence) CompleteUserHabit(user string, habitId uuid.UUID,
    onTarget bool, streak bool) error {
    log.Debug(fmt.Sprintf("completing habit %s for user %s...", habitId, user))
    // create new persitence session for graph and defer closing
    session := db.NewSession()
    defer session.Close()
    // generate config metadata for query
    cfg := map[string]interface{}{
        "habit_id": habitId.String(),
        "uid": user,
        "streak": streak,
        "on_target": onTarget,
    }
    handler := func(tx neo4j.Transaction) (interface{}, error) {
        var (query string; err error)
        // set completion date on habit as property
        query = `MATCH (u:User {uid: $uid})-[:OWNS]->(h:Habit {habit_id: $habit_id})
        SET h.last_completed = datetime()`
        _, err = tx.Run(query, cfg)
        if err != nil {
            log.Error(fmt.Errorf("unable to update habit with last completion: %+v", err))
            return nil, err
        }

        // create new completion node and set ownership to habbit
        query = `CREATE (c:HabitCompletion {
            event_timestamp: datetime(),
            on_target: $on_target,
            streak: $streak
        })
        WITH c
        MATCH (u:User {uid: $uid})-[:OWNS]->(h:Habit {habit_id: $habit_id})
        CREATE (h)-[:OWNS]->(c)
        `
        _, err = tx.Run(query, cfg)
        if err != nil {
            log.Error(fmt.Errorf("unable to update habit with last completion: %+v", err))
            return nil, err
        }
        return nil, nil
    }
    // get all habits from graph using persistence session
    _, err := session.WriteTransaction(handler)
    if err != nil {
        log.Error(fmt.Errorf("unable to complete habit: %+v", err))
        return err
    }
    return nil
}

// function used to complete a habit with given habit ID for user
func(db *GraphPersistence) DeleteUserHabit(user string, habitId uuid.UUID) error {
    log.Debug(fmt.Sprintf("deleting habit %s for user %s...", habitId, user))
    // create new persitence session for graph and defer closing
    session := db.NewSession()
    defer session.Close()
    // generate config metadata for query
    cfg := map[string]interface{}{
        "habit_id": habitId.String(),
        "uid": user,
    }
    handler := func(tx neo4j.Transaction) (interface{}, error) {
        query := `MATCH (u:User {uid: $uid})-[:OWNS]->(h:Habit {habit_id: $habit_id})
        DETACH DELETE h`
        return tx.Run(query, cfg)
    }
    // get all habits from graph using persistence session
    _, err := session.WriteTransaction(handler)
    if err != nil {
        log.Error(fmt.Errorf("unable to delete habit: %+v", err))
        return err
    }
    return nil
}

// function used to update a habit with given habit ID for user
func(db *GraphPersistence) UpdateUserHabit(user string, habitId uuid.UUID,
    habit Habit) error {
    log.Debug(fmt.Sprintf("updating habit %s for user %s...", habitId, user))
    // create new persitence session for graph and defer closing
    session := db.NewSession()
    defer session.Close()
    // generate config metadata for query
    cfg := map[string]interface{}{
        "habit_id": habitId.String(),
        "habit_name": habit.HabitName,
        "habit_description": habit.HabitDescription,
        "habit_cycle": orderCyclesString(habit.HabitCycle),
        "uid": user,
    }
    handler := func(tx neo4j.Transaction) (interface{}, error) {
        query := `MATCH (u:User {uid: $uid})-[:OWNS]->(h:Habit {habit_id: $habit_id})
        SET h.habit_name = $habit_name, h.habit_description = $habit_description,
        h.habit_cycle = $habit_cycle`
        return tx.Run(query, cfg)
    }
    // get all habits from graph using persistence session
    _, err := session.WriteTransaction(handler)
    if err != nil {
        log.Error(fmt.Errorf("unable to update habit: %+v", err))
        return err
    }
    return nil
}

func(db *GraphPersistence) GetHabitCompletions(user string, habitId uuid.UUID) ([]HabitCompletion, error) {
    log.Debug(fmt.Sprintf("fetching habit completions for habit %s...", habitId))
    completions := []HabitCompletion{}
    // create new persistence session for graph and defer closing
    session := db.NewSession()
    defer session.Close()
    // generate config metadata for query
    cfg := map[string]interface{}{
        "uid": user,
        "habit_id": habitId.String(),
    }
    // define handler function used to retrieve node data
    handler := func(tx neo4j.Transaction) (interface{}, error) {
        query := `MATCH (:User {uid: $uid})-[:OWNS]->(h:Habit {habit_id: $habit_id})-[:OWNS]->(c.HabitCompletion)
        RETURN c.event_timestamp, c.streak, c.on_target ORDER BY c.event_timestamp DESC`
        return neo4j.Collect(tx.Run(query, cfg))
    }
    // iterate over
    nodes, err := neo4j.AsRecords(session.ReadTransaction(handler))
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve habit completions: %+v", err))
        return completions, nil
    }
    for _, node := range(nodes) {
        completions = append(completions, HabitCompletion{
            EventTimestamp: node.Values[0].(time.Time),
            Streak: node.Values[1].(bool),
            OnTarget: node.Values[2].(bool),
        })
    }
    return completions, nil
}