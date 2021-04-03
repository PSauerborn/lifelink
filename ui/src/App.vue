<template>
  <v-app class="main-app">
    <v-main>
      <loading-page @authChecked="setAuth" v-if="!authChecked" />
      <login-page v-if="authChecked && !hasValidToken"
                  @loggedIn="hasValidToken = true" />
      <application-frame @logout="logout" v-if="authChecked && hasValidToken" />
    </v-main>
  </v-app>
</template>

<script>
import ApplicationFrame from './components/ApplicationFrame.vue';
import LoadingPage from './components/login/LoadingPage.vue';
import LoginPage from './components/LoginPage.vue';

export default {
  name: 'App',

  components: {
    LoginPage,
    LoadingPage,
    ApplicationFrame
  },
  methods: {
    setAuth(tokenValid) {
      // set variables to indicate if token if auth has been checked
      // and token is valid
      this.hasValidToken = tokenValid
      this.authChecked = true
    },
    logout() {
      localStorage.removeItem("lifelink_access_token")
      this.hasValidToken = false
    }
  },
  computed: {

  },
  data: () => ({
    authChecked: false,
    hasValidToken: false
  }),
  mounted() {

  }
};
</script>

<style scoped>
html {
  overflow-y: auto;
}
</style>