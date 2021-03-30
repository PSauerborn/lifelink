<template>
    <v-container fluid>
        <v-row align="center" justify="center" dense>
            <v-col align="center" justify="center" cols=12>
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
                <v-row>
                    <v-col align="center" justify="center" cols=4
                           v-for="habit in habits" :key="habit.habit_id">
                        <habit-card @reloadHabits="getHabits" @logout="$emit('logout')" :habit="habit"/>
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
    methods: {
        getHabits() {
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

</style>