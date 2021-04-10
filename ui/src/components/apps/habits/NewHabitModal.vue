<template>
    <v-card min-height="325px">
        <v-row align="center" justify="center" dense>
            <v-col cols=12>
                <v-card-title>
                    Create New Habit
                </v-card-title>
                <v-card-subtitle>
                    Create a new habit to track and link to other lifelink applications
                </v-card-subtitle>
                <v-divider></v-divider>
            </v-col>
        </v-row>
        <v-row align="center" justify="center" dense>
            <v-col cols=8>
                <v-card-text>
                    <v-form ref="form">
                        <v-text-field v-model="habitName"
                                    ref="habitNameTextField"
                                    label="Habit Name"
                                    autocomplete="off"
                                    :rules="rules.habit_name"
                                    maxlength="20"
                                    counter dense />
                        <v-text-field v-model="habitDescription"
                                    ref="habitDescriptionTextField"
                                    label="Habit Description"
                                    autocomplete="off"
                                    :rules="rules.habit_description" dense />
                    </v-form>
                 </v-card-text>
            </v-col>
        </v-row>
        <v-row align="center" justify="center" style="margin-top: 10px" dense>
            <v-col style="text-transform: capitalize" align="center" justify="center" v-for="(day, i) in habitCycle" :key="i" cols=1>
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
        <v-row align="center" justify="center" dense>
            <v-col cols=12>
                <v-divider></v-divider>
                <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn color="primary"
                           @click="createHabit"
                           :disabled="!habitName || !habitDescription">
                        Create Habit
                    </v-btn>
                </v-card-actions>
            </v-col>
        </v-row>
      </v-card>
</template>

<script>
import axios from 'axios';

export default {
    name: "NewHabitModal",
    computed: {
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
        }
    },
    methods: {
        closeModal() {
            this.$refs.form.reset()
            this.$emit('habitCreated')
        },
        // function used to add a new habit to the graph
        // by sending a POST request to the backend along
        // with the habit name, description and cycle
        createHabit() {
            // return if habit name or description have not been set
            if (!this.habitName || !this.habitDescription || this.habitName.length > 20) {
                return
            }
            let vm = this
            axios({
                method: 'post',
                url: process.env.VUE_APP_API_URL + 'api/habits/new',
                headers: {'Authorization': 'Bearer ' + localStorage.getItem("lifelink_access_token")},
                data: {
                    habit_name: vm.habitName,
                    habit_description: vm.habitDescription,
                    habit_cycle: vm.selectedDays.join(",")
                }
            }).then(function (response) {
                console.log(response)
                vm.$emit('reloadHabits')
                vm.closeModal()
            }).catch(function (error) {
                console.log(error)
                if (error.response.status == 401) {
                    vm.$emit('logout')
                }
            })
        },
        // function used to add new day to habit cycle
        addDayToHabit(day) {
            this.selectedDays.push(day)
        },
        // function used to remove day from habit cycle
        removeDayFromHabit(day) {
            const index = this.selectedDays.indexOf(day)
            this.selectedDays.splice(index, 1)
        }
    },
    data: () => ({
        validDays: ["mon", "tue", "wed", "thu", "fri", "sat", "sun"],
        selectedDays: [],
        habitName: '',
        habitDescription: '',
        rules: {
            habit_name: [val => (val || '').length > 0 || 'This field is required', v => v.length <= 20 || 'Max 20 characters'],
            habit_description: [val => (val || '').length > 0 || 'This field is required']
        }
    })
}
</script>

<style scoped>

</style>