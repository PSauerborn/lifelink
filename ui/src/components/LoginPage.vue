<template>
    <v-container class="login-page-container" fluid>
        <v-row align="center" justify="center" dense>
            <v-col align="center" justify="center" cols=4>
                <v-card min-width="500" flat>
                    <v-row class="component-container" align="center" justify="center" dense>
                        <v-col align="center" justify="center" cols=12>
                            <v-card-title class="justify-center" v-if="!signup">Login</v-card-title>
                            <v-card-title class="justify-center" v-if="signup">Register</v-card-title>
                            <v-img src="../assets/logo.jpg" :height="100" :width="100"></v-img>
                            <v-card-subtitle>Welcome to Lifelink</v-card-subtitle>
                        </v-col>
                    </v-row>
                    <v-row class="component-container" align="center" justify="center" dense>
                        <v-col align="center" justify="center" cols=6>
                            <v-divider></v-divider>
                        </v-col>
                    </v-row>
                    <v-card-text v-if="!signup">
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
                                <span class="signup-span" @click="signup = true">Not a Member? Signup Here</span>
                            </v-col>
                        </v-row>
                    </v-card-text>
                    <v-card-text v-if="signup">
                        <v-row class="component-container" align="center" justify="center" dense>
                            <v-col align="center" justify="center" cols=7>
                                <v-form v-model="formValid">
                                    <v-text-field v-model="username" autocomplete="off"
                                                class="register-text-field"
                                                prepend-icon="mdi-account"
                                                :rules="[rules.required, () => !userAlreadyExists || 'Username taken']"
                                                ref="usernameInput"
                                                @input="userAlreadyExists = false"
                                                label="Username" dense />
                                    <v-text-field v-model="email" autocomplete="off"
                                                class="register-text-field"
                                                prepend-icon="mdi-at"
                                                :rules="[rules.required, rules.email]"
                                                label="Email" :type="'email'" dense />
                                    <v-text-field v-model="password" autocomplete="off"
                                                class="register-text-field"
                                                prepend-icon="mdi-lock"
                                                :rules="[rules.required]"
                                                label="Password"
                                                :type="show1 ? 'text' : 'password'"
                                                :append-icon="show1 ? 'mdi-eye' : 'mdi-eye-off'"
                                                @click:append="show1 = !show1" dense />
                                    <v-text-field v-model="confirmPassword" autocomplete="off"
                                                class="register-text-field"
                                                prepend-icon="mdi-lock"
                                                :rules="[rules.required, () => password == confirmPassword || 'Passwords must match']"
                                                label="Confirm Password"
                                                :type="show2 ? 'text' : 'password'"
                                                :append-icon="show2 ? 'mdi-eye' : 'mdi-eye-off'"
                                                @click:append="show2 = !show2" dense />
                                    <v-btn :disabled="!formValid" color="primary" width="200" @click="register" depressed>Register</v-btn>
                                </v-form>
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
        },
        toggleSignup() {
            this.signup = true
        },
        register() {
            let vm = this
            axios({
                method: 'post',
                url: process.env.VUE_APP_IDP_URL + 'authenticate/register',
                data: {uid: vm.username, password: vm.password, email: vm.email}
            }).then(function (response) {
                console.log(response)
                localStorage.setItem('lifelink_access_token', response.data.token)
                vm.login()
            }).catch(function (error) {
                console.log(error)
                if (error.response.status == 400) {
                    console.log('setting user exists status')
                    vm.userAlreadyExists = true
                    vm.$refs.usernameInput.validate()
                }
            })
        }
    },
    data: () => ({
        username: '',
        password: '',
        email: '',
        confirmPassword: '',
        signup: false,
        show1: false,
        show2: false,
        userAlreadyExists: false,
        formValid: false,
        rules: {
          required: value => !!value || 'Required.',
          min: v => v.length >= 8 || 'Min 8 characters',
          email: value => {
            const pattern = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
            return pattern.test(value) || 'Invalid e-mail.'
          },
        },
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

.signup-span:hover {
    cursor: pointer;
    color: #449FFF;
}

.register-text-field {
    margin-bottom: 10px;
}

</style>