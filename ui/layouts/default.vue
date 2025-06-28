<template>
    <v-app>
        <v-app-bar app density="compact"  class="liquid-glass">
            <v-tabs>
                <v-tab to="/">Dashboard</v-tab>
                <v-tab to="/agents">Agents</v-tab>
                <!-- <v-tab to="/files">Files</v-tab> -->
                <v-tab to="/settings">Settings</v-tab>
            </v-tabs>
        </v-app-bar>

        <NuxtPage />
    </v-app>
</template>

<script setup lang="ts">

</script>