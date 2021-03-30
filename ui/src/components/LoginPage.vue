<template>
    <v-container class="login-page-container" fluid>
        <v-row align="center" justify="center" dense>
            <v-col align="center" justify="center" cols=4>
                <v-card min-width="500">
                    <v-row class="component-container" align="center" justify="center" dense>
                        <v-col align="center" justify="center" cols=12>
                            <v-card-title class="justify-center">Login</v-card-title>
                            <v-avatar size="70"><img src="../assets/logo.jpg"></v-avatar>
                            <v-card-subtitle>Welcome to LifeLink</v-card-subtitle>
                        </v-col>
                    </v-row>
                    <v-row class="component-container" align="center" justify="center" dense>
                        <v-col align="center" justify="center" cols=6>
                            <v-divider></v-divider>
                        </v-col>
                    </v-row>
                    <v-card-text>
                        <v-row class="component-container" align="center" justify="center" dense>
                            <v-col align="center" justify="center" cols=7>
                                <v-text-field v-model="username" autocomplete="off"
                                              prepend-icon="mdi-account"
                                              label="Username" dense />
                                <v-text-field v-model="password" autocomplete="off"
                                              prepend-icon="mdi-lock"
                                              label="Password" :type="'password'" dense />
                                <v-btn color="primary" width="200" @click="login" depressed>Login</v-btn>
                            </v-col>
                        </v-row>
                        <v-row class="component-container" align="center" justify="center" dense>
                            <v-col align="center" justify="center" cols=10>
                                Not a Member? Signup here
                            </v-col>
                        </v-row>
                    </v-card-text>
                </v-card>
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
import axios from 'axios';

export default {
    components: {

    },
    name: 'LoginPage',
    methods: {
        login() {
            let vm = this
            axios({
                method: 'post',
                url: process.env.VUE_APP_IDP_URL + 'authenticate/token',
                data: {uid: vm.username, password: vm.password}
            }).then(function (response) {
                console.log(response)
                localStorage.setItem('lifelink_access_token', response.data.token)
                vm.$emit('loggedIn')
            }).catch(function (error) {
                console.log(error)
            })
        }
    },
    data: () => ({
        username: '',
        password: ''
    })
}
</script>

<style scoped>

.component-container {
    margin: 10px;
}

.login-page-container {
    margin-top: 120px;
}
</style>