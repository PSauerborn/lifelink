package habits

import (
    "fmt"
    "time"
    "strings"

    "github.com/google/uuid"
    log "github.com/sirupsen/logrus"
)

var cycleMappings = map[string]int{
    "sun": 0,
    "mon": 1,
    "tue": 2,
    "wed": 3,
    "thur": 4,
    "fri": 5,
    "sat": 6,
}

// function used to determine if a habit is currently on target
func isHabitOnTarget(habit Habit) bool {
    // return false if never completed
    if habit.LastCompleted == nil {
        return false
    }
    // de-reference pointer to get datetime instance and get weekday
    lastCompleted := *(habit.LastCompleted)
    lastCompletedWeekday := int(lastCompleted.Weekday())

    // get current week and year for last completed and current timestamp
    yearNow, weekNow := time.Now().ISOWeek()
    yearCompleted, weekCompleted := lastCompleted.ISOWeek()
    // if week and year of now and last completed dont match, return false
    if yearNow != yearCompleted || weekNow != weekCompleted {
        return false
    }

    cycle := strings.Split(habit.HabitCycle, ",")
    // get last due date
    lastDueDay := cycle[len(cycle) - 1]
    return lastCompletedWeekday >= cycleMappings[lastDueDay]
}

// function used to complete user habit
func completeUserHabit(uid string, habitId uuid.UUID) error {
    // get current habit from graph
    habit, err := persistence.GetHabitByHabitId(uid, habitId)
    if err != nil {
        log.Error(fmt.Errorf("unable to retrieve habit from graph: %+v", err))
        return err
    }
    // if habit is already on target, add as exta/non-streaked completion
    if isHabitOnTarget(habit) {
        log.Debug(fmt.Sprintf("completing habit %s for user %s as not streaked and on target", habitId, uid))
        return persistence.CompleteUserHabit(uid, habitId, true, false)
    }
    now := time.Now().UTC()
    // generate n + 1 habit by copying current habit and inserting current timestamp
    futureHabit := habit
    futureHabit.LastCompleted = &now
    // if habit is currently not on target, but inserting n + 1 will
    // cause it to be on target, add to graph as a streaked, on-target completion
    if !isHabitOnTarget(habit) && isHabitOnTarget(futureHabit) {
        log.Debug(fmt.Sprintf("completing habit %s for user %s as streaked and on target", habitId, uid))
        return persistence.CompleteUserHabit(uid, habitId, true, true)
    // else add to graph as streaked, overdue habit
    } else {
        log.Debug(fmt.Sprintf("completing habit %s for user %s as streaked and overdue", habitId, uid))
        return persistence.CompleteUserHabit(uid, habitId, false, true)
    }
}