<template>
    <v-container class="loading-page-container" fluid>
        <v-row class="component-container" align="center" justify="center" dense>
            <v-col align="center" justify="center" cols=6>
                Welcome to LifeLink
            </v-col>
        </v-row>
        <v-row class="component-container" align="center" justify="center" dense>
            <v-col align="center" justify="center" cols=4>
                <v-divider></v-divider>
            </v-col>
        </v-row>
        <v-row class="component-container" align="center" justify="center" dense>
            <v-col align="center" justify="center" cols=6>
                <v-progress-circular :size="70" :width="7" indeterminate color="primary" />
            </v-col>
        </v-row>
        <v-row class="component-container" align="center" justify="center" dense>
            <v-col align="center" justify="center" cols=6>
                Checking Authorization... Please Wait
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
import axios from 'axios';

export default {
    name: "LoadingPage",
    methods: {
        checkAuth() {
            // get token from local storage. if not present, back to login
            let token = localStorage.getItem("lifelink_access_token")
            if (!token) {
                console.log('token not found in local storage. redirecting to login...')
                this.$emit('authChecked', false)
                return
            }
            // if access token is found, make request to health check route to check if
            // token is valid
            console.log('found access token in local storage. checking for validity...')
            let vm = this
            axios({
                method: 'get',
                url: process.env.VUE_APP_API_URL + 'api/habits/health_check',
                headers: {'Authorization': 'Bearer ' + token}
            }).then(function (response) {
                console.log(response)
                vm.$emit('authChecked', true)
            }).catch(function (error) {
                console.log(error)
                vm.$emit('authChecked', false)
            })
        }
    },
    mounted() {
        // check authorization when component is mounted
        this.checkAuth()
    }
}
</script>

<style scoped>


.loading-page-container {
    margin-top: 150px;
}

.component-container {
    margin: 30px;
}

</style>