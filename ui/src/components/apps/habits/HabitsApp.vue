<template>
    <v-container class="habits-app-container" fluid>
        <v-row align="center" justify="center" dense>
            <v-col align="center" justify="center" cols=12>
                <v-row align="center" justify="center" dense>
                    <v-col align="center" justify="center" cols=2>
                            <span>{{ mostCompletedHabit.habit_name }}</span><br>
                            <span style="font-size: 45px">{{ mostCompletedHabit.completions }}</span><br>
                            <span style="font-size: 30px; font-weight: bold">Most Completed</span>
                    </v-col>
                    <v-divider vertical></v-divider>
                    <v-col align="center" justify="center" xl=2 lg=3 md=3>
                            <span>{{ mostStreakedHabit.habit_name }}</span><br>
                            <span style="font-size: 45px">{{ mostStreakedHabit.streak }}</span><br>
                            <span style="font-size: 30px; font-weight: bold">Highest Streak</span>
                    </v-col>
                    <v-divider vertical></v-divider>
                    <v-col align="center" justify="center" xl=2 lg=3 md=3>
                            <span style="font-size: 45px">{{ this.habits.length }}</span><br>
                            <span style="font-size: 30px; font-weight: bold">Total Habits</span>
                    </v-col>
                </v-row>
                <v-row style="margin-right: 50px" align="end" justify="end" dense>
                    <v-col align="end" justify="end">
                        <v-dialog v-model="dialog" width="600">
                            <template v-slot:activator="{ on, attrs }">
                                <v-btn color="primary" v-bind="attrs" v-on="on"><v-icon>mdi-plus</v-icon>New Habit</v-btn>
                            </template>
                            <new-habit-modal @reloadHabits="getHabits"
                                             @habitCreated="dialog = false"
                                             @logout="$emit('logout')" />
                        </v-dialog>
                    </v-col>
                </v-row>
                <v-row align="start" justify="start" dense>
                    <v-col align="start" justify="start"
                           v-for="habit in habits"
                           :key="habit.habit_id"
                           lg="4"
                           xl="3"
                           md="6">
                        <habit-card @reloadHabits="getHabits"
                                    @logout="$emit('logout')"
                                    :habit="habit" />
                    </v-col>
                </v-row>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
import axios from 'axios';

import HabitCard from './HabitCard.vue'
import NewHabitModal from './NewHabitModal.vue'
export default {
  components: { HabitCard, NewHabitModal },
    name: "HabitsApp",
    computed: {
        // computed property used to get most completed habit
        mostCompletedHabit() {
            if (this.habits.length == 0) {
                return {habit_name: 'None', completions: 0}
            }
            let topHabit = this.habits[0]
            this.habits.forEach(habit => {
                if (habit.completions > topHabit.completions) {
                    topHabit = habit
                }
            })
            return topHabit
        },
        // computed property used to get most streaked habit
        mostStreakedHabit() {
            if (this.habits.length == 0) {
                return {habit_name: 'None', streak: 0}
            }
            let topHabit = this.habits[0]
            this.habits.forEach(habit => {
                if (habit.streak > topHabit.streak) {
                    topHabit = habit
                }
            })
            return topHabit
        }
    },
    methods: {
        getHabits() {
            console.log(this.$vuetify.breakpoint.name)
            let vm = this
            axios({
                method: 'get',
                url: process.env.VUE_APP_API_URL + 'api/habits/all',
                headers: {'Authorization': 'Bearer ' + localStorage.getItem("lifelink_access_token")}
            }).then(function (response) {
                console.log(response)
                vm.habits = response.data.habits
            }).catch(function (error) {
                console.log(error)
                if (error.response.status == 401) {
                    vm.$emit('logout')
                }
            })
        }
    },
    data: () => ({
        habits: [],
        dialog: false
    }),
    mounted() {
        this.getHabits()
    }
}
</script>

<style scoped>

@import url('https://fonts.googleapis.com/css2?family=Source+Sans+Pro:wght@300;400&display=swap');

.habits-app-container {
    font-family: "Source Sans Pro";
    margin-top: 50px;
}
</style>