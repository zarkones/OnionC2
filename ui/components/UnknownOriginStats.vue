<template>
    <v-btn @click="gotoUnknownOriginAgents" variant="outlined" density="compact">
        <v-icon icon="mdi-crosshairs-question" />
        {{ mockedAgentAmount }}
        <v-tooltip activator="parent" location="top" class="unknown-origin-tooltip" open-delay="600">
            <h4>Unknown Origin</h4>
            There {{ mockedAgentAmount === 1 ? 'is' : 'are' }} {{ mockedAgentAmount }}
            {{ mockedAgentAmount === 1 ? 'agent who\'s' : 'agents whose' }} IP {{ mockedAgentAmount === 1 ? 'address is' : 'addresses are' }} not yet known.
            {{ mockedAgentAmount === 1 ? 'This agent' : 'These agents' }} won't be shown in the map. 
            Clicking on this button would reveal which {{ mockedAgentAmount === 1 ? 'agent is' : 'agents are' }} in question.
        </v-tooltip>
    </v-btn>
</template>

<script setup lang="ts">

const mockedAgentAmount = ref(164)

const gotoUnknownOriginAgents = () => {
    // TODO: Redirect to the agents page with unknown origin filter applied.
}

</script>

<style>

.unknown-origin-tooltip {
    max-width: 340px;
}

</style>