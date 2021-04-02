<template>
    <v-container class="app-coming-container" fluid>
        <v-row align="center" justify="center" dense>
            <v-col align="center" justify="center" cols=6>
                This Application is not yet ready and available, but we are working on it!
            </v-col>
        </v-row>
        <v-row align="center" justify="center" dense>
            <v-col align="center" justify="center" cols=6>
                <v-img :height="500" :width="700" src="../../../assets/gopher.png" />
            </v-col>
        </v-row>
        <v-row align="center" justify="center" dense>
            <v-col align="center" justify="center" cols=6>
                Please come back soon for updates
            </v-col>
        </v-row>
    </v-container>
</template>

<script>
export default {
    name: "AppComing"
}
</script>

<style scoped>


.app-coming-container {
    margin-top: 50px;
    font-size: 18px;
}

</style>