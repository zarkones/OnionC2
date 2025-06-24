<template>
    <Map></Map>
    <div class="info-charts liquid-glass">
        <div class="actions">
            <v-btn @click="toggleShowStats" density="compact" variant="outlined">
                <!-- {{ hideStats ? 'Show' : 'Hide' }} -->
                <v-icon v-if="hideStats" icon="mdi-chevron-up" />
                <v-icon v-else icon="mdi-chevron-down" />
            </v-btn>
        </div>
        <v-divider v-show="!hideStats" class="mb-1 mt-3"/>
        <InfoCharts v-show="!hideStats" ></InfoCharts>
    </div>
</template>

<script setup lang="ts">

const hideStats = ref(false)

const toggleShowStats = () => {
    hideStats.value = !hideStats.value
}

</script>

<style scoped>
.info-charts {
    position: absolute;
    left: 0px;
    width: 100%;
    height: fit-content;
    bottom: 0px;
    padding: 8px;
}

.actions {
}
</style>