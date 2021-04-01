package habits

import (
    "fmt"
    "time"
    "strings"
    "errors"

    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"
)

var (
    // define custom errors
    ErrInvalidHabitStatus = errors.New("Habit has invalid status")
)

// define mappings for cycles to convert
// from cycle strings to integer days
var cycleMappings = map[string]int{
    "sun": 0,
    "mon": 1,
    "tue": 2,
    "wed": 3,
    "thu": 4,
    "fri": 5,
    "sat": 6,
}

// define reverse mappings for cycles to
// convert from integer days to strings
var reverseCycleMappings = map[int]string{
    0: "sun",
    1: "mon",
    2: "tue",
    3: "wed",
    4: "thu",
    5: "fri",
    6: "sat",
}

// function used to evaluate the last due date for a given habit
// based on it last completion date
func getHabitDueDate(habit Habit) time.Time {
    var ts time.Time
    if habit.LastCompleted != nil {
        // get ts for midnight following last completion date
        year, month, day := habit.LastCompleted.Date()
        ts = time.Date(year, month, day + 1, 0, 0, 0, 0, time.UTC)
    } else {
        // get ts for midnight of current day
        year, month, day := habit.Created.Date()
        ts = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
    }
    // convert habit cycle to array and find last due date
    cycle := strings.Split(habit.HabitCycle, ",")
    for {
        day := reverseCycleMappings[int(ts.Weekday())]
        // if next day is present in required cycle, break out
        // out of loop
        if stringSliceContains(cycle, day) {
            break
        }
        ts = ts.Add(time.Hour * 24)
    }
    return ts.Add(time.Hour * 24)
}

// function used to determine if a habit is due
// on any given day given the
func habitDueToday(habit Habit) bool {
    // get current due date for habit
    dueDate := getHabitDueDate(habit)
    log.Debug(fmt.Sprintf("checking if habit is due with reference due date %s", dueDate))
    // construct theoretical due date if due today and
    // compare to actual due date
    year, month, day := time.Now().UTC().Date()
    ts := time.Date(year, month, day + 1, 0, 0, 0, 0, time.UTC)
    log.Debug(fmt.Sprintf("%s: %s", dueDate, ts))
    return ts == dueDate
}

// function used to determine if a habit is overdue
// based on its current due date
func habitOverdue(habit Habit) bool {
    // get current due date for habit
    dueDate := getHabitDueDate(habit)
    log.Debug(fmt.Sprintf("checking if habit is overdue with reference due date %s", dueDate))
    return time.Now().UTC().After(dueDate)
}

// function used to retrieve habit status. habit status
// is returned as either due, overdue or on target
func getHabitStatus(habit Habit) string {
    if habitOverdue(habit) {
        return "overdue"
    } else if habitDueToday(habit) {
        return "due"
    } else {
        return "on-target"
    }
}

// function used to complete user habits. habits are evaluated
// as on-target and as streak contributers on each completion
func completeHabit(uid string, habitId uuid.UUID) error {
    // get current habit from graph
    habit, err := persistence.GetHabitByHabitId(uid, habitId)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve habit from graph: %+v", err))
        return err
    }

    log.Debug(fmt.Sprintf("habit due date evaluated as %s", getHabitDueDate(habit)))
    switch getHabitStatus(habit) {
    case "due":
        log.Debug("habit on target. adding with streak")
        return persistence.CompleteUserHabit(uid, habitId, true, habit.Streak + 1)
    case "overdue":
        log.Debug("habit overdue. adding without streak")
        return persistence.CompleteUserHabit(uid, habitId, false, 0)
    case "on-target":
        log.Debug("habit already on target. adding without streak")
        return persistence.CompleteUserHabit(uid, habitId, true, habit.Streak)
    default:
        return ErrInvalidHabitStatus
    }
}