<template>
    <v-container>
        <v-app-bar class="app-bar-main"
        absolute
        color="white"
        elevate-on-scroll
        >
        <v-app-bar-nav-icon></v-app-bar-nav-icon>

        <v-toolbar-title>Lifelink</v-toolbar-title>
        <v-spacer></v-spacer>

        <v-btn icon>
            <v-icon>mdi-magnify</v-icon>
        </v-btn>
        <v-menu left bottom>
            <template v-slot:activator="{ on, attrs }">
            <v-btn
                icon
                v-bind="attrs"
                v-on="on"
            >
                <v-icon>mdi-dots-vertical</v-icon>
            </v-btn>
            </template>

            <v-list>
            <v-list-item>
                <v-list-item-title>Profile</v-list-item-title>
            </v-list-item>
            <v-list-item @click="logout">
                <v-list-item-title>Logout</v-list-item-title>
            </v-list-item>
            </v-list>
        </v-menu>
        <template v-slot:extension>
        <v-tabs align-with-title>
            <v-tab>Habits</v-tab>
            <v-tab-item><habits-app @logout="logout" /></v-tab-item>
            <v-tab>Timesheets</v-tab>
            <v-tab-item><timesheets-app @logout="logout" /></v-tab-item>
            <v-tab>Guitar Pro</v-tab>
            <v-tab-item><guitar-app @logout="logout" /></v-tab-item>
            </v-tabs>
        </template>
        </v-app-bar>
    </v-container>
</template>

<script>
import GuitarApp from './apps/guitar/GuitarApp.vue'
import HabitsApp from './apps/habits/HabitsApp.vue'
import TimesheetsApp from './apps/timesheets/TimesheetsApp.vue'
export default {
  components: { HabitsApp, TimesheetsApp, GuitarApp },
    name: "ApplicationFrame",
    methods: {
        logout() {
            this.$emit('logout')
        }
    }
}
</script>

<style scoped>

.app-bar-main {
    margin-bottom: 20px;
}

</style>