<template>
    <v-card min-width="450" :style="{'border-left': borderColor}">
        <v-expansion-panels v-model="panel" :disabled=disabled>
            <v-expansion-panel>
                <v-expansion-panel-header>
                    <v-card-text>
                        <v-row style="margin-bottom: 10px" dense>
                            <v-col cols=7>
                                <v-card-subtitle>{{ habit.habit_name }}</v-card-subtitle>
                            </v-col>
                            <v-divider vertical />
                            <v-col align="end" justify="end" cols=5>
                                <v-btn @click.native.stop icon @click="completeHabit"><v-icon>mdi-check</v-icon></v-btn>
                                <v-btn @click.native.stop icon><v-icon>mdi-pencil</v-icon></v-btn>
                                <v-btn @click.native.stop icon @click="deleteHabit"><v-icon>mdi-delete</v-icon></v-btn>
                            </v-col>
                        </v-row>
                        <v-row style="margin-top: 10px" dense>
                            <v-col style="text-transform: capitalize" v-for="(day, i) in habitCycle" :key="i" cols=1.5>
                                <v-btn @click.native.stop v-if="day.included"
                                       color="primary"
                                       @click="removeDayFromHabit(day.day)"
                                       x-small rounded dense depressed>{{ day.day }}</v-btn>
                                <v-btn @click.native.stop v-if="!day.included"
                                       color="primary"
                                       @click="addDayToHabit(day.day)"
                                       outlined x-small rounded dense depressed>{{ day.day }}</v-btn>
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-expansion-panel-header>
                <v-expansion-panel-content>
                    <v-row align="center" justify="center" dense>
                        <v-col align="center" justify="center" cols=12>
                            <v-divider></v-divider>
                            <v-row align="center" justify="center" dense>
                                <v-col align="center" justify="center">
                                    <v-card-text>
                                        {{ createdDate }}<br>
                                        <span style="font-size: 20px; font-weight: bold">Created</span>
                                    </v-card-text>
                                </v-col>
                                <v-col>
                                    <v-card-text>
                                        {{ lastCompleted }}<br>
                                        <span style="font-size: 20px; font-weight: bold">Completed</span>
                                    </v-card-text>
                                </v-col>
                                <v-col>
                                    <v-card-text>
                                        {{ streak }}<br>
                                        <span style="font-size: 20px; font-weight: bold">Streak</span>
                                    </v-card-text>
                                </v-col>
                            </v-row>
                            <v-row justify="start" align="start" dense>
                                <v-col justify="start" align="start">
                                    <v-card-subtitle>{{ habit.habit_description }}</v-card-subtitle>
                                </v-col>
                            </v-row>
                        </v-col>
                    </v-row>
                </v-expansion-panel-content>
            </v-expansion-panel>
        </v-expansion-panels>
    </v-card>
</template>

<script>
import axios from 'axios';
import moment from 'moment';

export default {
    name: "HabitCard",
    props: {
        habit: {
            type: Object,
            default: () => ({
                habit_title: "Default Habit Title",
                habit_description: "Some random habit used for testing",
                habit_cycle: "mon,wed,fri",
                created: "2021-01-01"
            })
        },
        habit_id: Number
    },
    computed: {
        // computed property used to evaluate the current
        // habit cycle as an containing information on
        // whether or not the current day is actice
        habitCycle() {
            const cycle = this.selectedDays
            let formatted = []
            this.validDays.forEach(element => {
                if (cycle.includes(element)) {
                    formatted.push({day: element, included: true})
                } else {
                    formatted.push({day: element, included: false})
                }
            })
            return formatted
        },
        // computed property used to determine if a habit is due on the present day
        isHabitDue() {
            const d = new Date().getDay()
            return this.selectedDays.includes(this.dayMappings[d])
        },
        // computed property to return the border color
        // of the habit card based on whether or not the habit
        // is due on the present day
        borderColor() {
            if (this.isHabitDue) {
                return "5px solid red"
            } else {
                return ""
            }
        },
        createdDate() {
            return moment(this.habit.created).format("YYYY-MM-DD")
        },
        lastCompleted() {
            if (!this.habit.last_completed) {
                return 'Never'
            }
            return moment(this.habit.last_completed).format("YYYY-MM-DD")
        },
        streak() {
            return 0
        }
    },
    methods: {
        // function used to add a given day to the current
        // habits cycle
        addDayToHabit(day) {
            console.log("adding day to habit: " + day)
            this.selectedDays.push(day)
            this.updateHabit()
        },
        // function used to remove a given day from the current
        // habits cycle
        removeDayFromHabit(day) {
            console.log("removing day from habit: " + day)
            const index = this.selectedDays.indexOf(day)
            this.selectedDays.splice(index, 1)
            this.updateHabit()
        },
        // function used to complete a habit with given habit
        // ID for that specific day
        completeHabit() {
            let vm = this
            axios({
                method: 'patch',
                url: process.env.VUE_APP_API_URL + 'api/habits/complete/' + vm.habit.habit_id,
                headers: {'Authorization': 'Bearer ' + localStorage.getItem("lifelink_access_token")},
            }).then(function (response) {
                console.log(response)
                vm.$emit('reloadHabits')
            }).catch(function (error) {
                console.log(error)
                if (error.response.status == 401) {
                    vm.$emit('logout')
                }
            })
        },
        // function used to update a habit. habits are updated
        // by sending the new habit using a PUT request with
        // given habit ID to the backend
        updateHabit() {
            let vm = this
            axios({
                method: 'put',
                url: process.env.VUE_APP_API_URL + 'api/habits/update/' + vm.habit.habit_id,
                headers: {'Authorization': 'Bearer ' + localStorage.getItem("lifelink_access_token")},
                data: {
                    habit_name: vm.habit.habit_name,
                    habit_description: vm.habit.habit_description,
                    habit_cycle: vm.selectedDays.join(",")
                }
            }).then(function (response) {
                console.log(response)
                vm.$emit('reloadHabits')
            }).catch(function (error) {
                console.log(error)
                if (error.response.status == 401) {
                    vm.$emit('logout')
                }
            })
        },
        // function used to delete a habit from the graph.
        // habits are deleted by sending a DELETE request
        // with a supplied habit ID to the backend
        deleteHabit() {
            let vm = this
            axios({
                method: 'delete',
                url: process.env.VUE_APP_API_URL + 'api/habits/delete/' + vm.habit.habit_id,
                headers: {'Authorization': 'Bearer ' + localStorage.getItem("lifelink_access_token")}
            }).then(function (response) {
                console.log(response)
                vm.$emit('reloadHabits')
            }).catch(function (error) {
                console.log(error)
                if (error.response.status == 401) {
                    vm.$emit('logout')
                }
            })
        }
    },
    data: () => ({
        disabled: false,
        panel: [0, 1],
        validDays: ["mon", "tue", "wed", "thu", "fri", "sat", "sun"],
        selectedDays: [],
        dayMappings: {
            0: "sun",
            1: "mon",
            2: "tue",
            3: "wed",
            4: "thu",
            5: "fri",
            6: "sat"
        }
    }),
    mounted() {
        this.selectedDays = this.habit.habit_cycle.split(",")
    }
}
</script>

<style scoped>

.habit-due {
    border-left: 5px solid red;
}

</style>