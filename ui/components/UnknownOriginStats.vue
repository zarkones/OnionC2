<template>
    <v-btn @click="gotoUnknownOriginAgents" variant="plain" density="compact">
        <v-icon icon="mdi-crosshairs-question" class="mr-1" />
        {{ API.store.stats.unknownOriginCount }}
        <v-tooltip activator="parent" location="top" class="unknown-origin-tooltip" open-delay="600">
            <h4>Unknown Origin</h4>
            There {{ API.store.stats.unknownOriginCount === 1 ? 'is' : 'are' }} {{ API.store.stats.unknownOriginCount }}
            {{ API.store.stats.unknownOriginCount === 1 ? 'agent who\'s' : 'agents whose' }} IP {{ API.store.stats.unknownOriginCount === 1 ? 'address is' : 'addresses are' }} not yet known.
            {{ API.store.stats.unknownOriginCount === 1 ? 'This agent' : 'These agents' }} won't be shown in the map. 
            Clicking on this button would reveal which {{ API.store.stats.unknownOriginCount === 1 ? 'agent is' : 'agents are' }} in question.
        </v-tooltip>
    </v-btn>
</template>

<script setup lang="ts">

const router = useRouter()

const gotoUnknownOriginAgents = () => {
    router.push({
        path: '/agents',
        query: { origins: 'unknown' }
    })
}

</script>

<style>

.unknown-origin-tooltip {
    max-width: 340px;
}

</style>