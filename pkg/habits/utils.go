package habits

import (
    "fmt"
    "strings"

    log "github.com/sirupsen/logrus"
)

var (
    validCycles = []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
)

// helper function used to determine is a string
// slice contains a given element
func stringSliceContains(slice []string, element string) bool {
    for _, ele := range(slice) {
        if ele ==  element {
            return true
        }
    }
    return false
}

// function used to cast a slice of strings to
// lowercase
func stringSliceToLower(slice []string) []string {
    lowered := []string{}
    for _, ele := range(slice) {
        lowered = append(lowered, strings.ToLower(ele))
    }
    return lowered
}

// function used to check if a given cycle is valid
func isValidCycle(cycle string) bool {
    log.Debug(fmt.Sprintf("checking cycle %s for validity", cycle))
    // iterate over days given in comma-separarted cycle
    // and check that all days are valid
    for _, day := range(strings.Split(cycle, ",")) {
        if !stringSliceContains(validCycles, strings.ToLower(day)) {
            log.Warn(fmt.Sprintf("cannot process cycle: invalid cycle day %s", day))
            return false
        }
    }
    return len(cycle) > 0
}

// function used to order cycles and convert
// comma separated list of days into a slice
func orderCyclesSlice(cycle string) []string {
    ordered := []string{}
    // convert original cycle to string array
    cycleSlice := stringSliceToLower(strings.Split(cycle, ","))
    // iterate over valid days (already ordered) and append
    // present values into ordered array
    for _, day := range(validCycles) {
        if stringSliceContains(cycleSlice, day) {
            ordered = append(ordered, day)
        }
    }
    return ordered
}

// function used to order cycles and convert
// comma separated list of days into a slice
func orderCyclesString(cycle string) string {
    ordered := orderCyclesSlice(cycle)
    return strings.Join(ordered, ",")
}

