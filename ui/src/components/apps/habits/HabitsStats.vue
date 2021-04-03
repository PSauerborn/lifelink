<template>
    <v-container fluid>
        <v-row v-if="!onMobile" align="center" justify="center" dense>
            <v-col align="center" justify="center" xl=2 lg=3 md=3 sm=3>
                <span>{{ mostCompletedHabit.habit_name }}</span><br>
                <span style="font-size: 45px">{{ mostCompletedHabit.completions }}</span><br>
                <span style="font-size: 30px; font-weight: bold">Most Completed</span>
            </v-col>
            <v-divider class="stat-divider" :vertical="!onMobile"></v-divider>
            <v-col align="center" justify="center" xl=2 lg=3 md=3 sm=3>
                <span>{{ mostStreakedHabit.habit_name }}</span><br>
                <span style="font-size: 45px">{{ mostStreakedHabit.streak }}</span><br>
                <span style="font-size: 30px; font-weight: bold">Highest Streak</span>
            </v-col>
            <v-divider class="stat-divider" :vertical="!onMobile"></v-divider>
            <v-col align="center" justify="center" xl=2 lg=3 md=3 sm=3>
                <span style="font-size: 45px">{{ this.habits.length }}</span><br>
                <span style="font-size: 30px; font-weight: bold">Total Habits</span>
            </v-col>
        </v-row>
        <v-row v-if="onMobile" align="center" justify="center" dense>
            <v-col align="center" justify="center" cols=12>
                <span>{{ mostCompletedHabit.habit_name }}</span><br>
                <span style="font-size: 45px">{{ mostCompletedHabit.completions }}</span><br>
                <span style="font-size: 30px; font-weight: bold">Most Completed</span>
                <v-divider class="stat-divider-mobile"></v-divider>
            </v-col>
            <v-col align="center" justify="center" cols=12>
                <span>{{ mostStreakedHabit.habit_name }}</span><br>
                <span style="font-size: 45px">{{ mostStreakedHabit.streak }}</span><br>
                <span style="font-size: 30px; font-weight: bold">Highest Streak</span>
                <v-divider class="stat-divider-mobile"></v-divider>
            </v-col>
            <v-col align="center" justify="center" cols=12>
                <span style="font-size: 45px">{{ this.habits.length }}</span><br>
                <span style="font-size: 30px; font-weight: bold">Total Habits</span>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
export default {
    name: "HabitsStats",
    props: {
        habits: Array
    },
    computed: {
        onMobile() {
            console.log(this.$vuetify.breakpoint.name)
            switch (this.$vuetify.breakpoint.name) {
                case 'xs':
                    return true
                case 'sm':
                    return false
                default:
                    return false
            }
        },
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
    }
}
</script>

<style scoped>

.stat-divider {
    margin-left: 10px;
    margin-right: 10px;
}

.stat-divider-mobile {
    margin-top: 15px;
    margin-bottom: 15px;
}

</style>