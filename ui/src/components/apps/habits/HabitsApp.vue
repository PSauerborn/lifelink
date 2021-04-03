<template>
    <v-container class="habits-app-container" fluid>
        <v-row align="center" justify="center" dense>
            <v-col align="center" justify="center" cols=12>
                <habits-stats :habits="habits" />
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
                    <v-col align="start" justify="start" v-for="habit in habits"
                           :key="habit.habit_id" lg="4" xl="3" md="6">
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
import HabitsStats from './HabitsStats.vue';
export default {
  components: { HabitCard, NewHabitModal, HabitsStats },
    name: "HabitsApp",
    computed: {

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