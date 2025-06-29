<template>
    <MapSvg style="position: absolute;"></MapSvg>

    <v-data-table-virtual
        :headers="(headers as any)"
        :items="API.getAgents().data.value"
        density="compact"
        item-key="name"
        class="table-el liquid-glass"
        fixed-header
    >
        <template v-slot:item.actions="{ item }">
            <ControlAgentDialog
                :agentId="item.id"
            />
        </template>

    </v-data-table-virtual>
</template>

<script setup lang="ts">

const headers = [
    // { title: 'ID', align: 'start', key: 'id' },
    { title: 'IP', align: 'start', key: 'ip' },
    { title: 'Hostname', align: 'start', key: 'hostname' },
    { title: 'OS', align: 'start', key: 'os' },
    { title: 'CPU Name', align: 'start', key: 'cpuName' },
    { title: 'CPU Arch', align: 'start', key: 'arch' },
    { title: 'RAM', align: 'start', key: 'ram' },
    { title: 'Last Seen', align: 'start', key: 'lastSeen' },
    { title: 'Actions', align: 'end', key: 'actions' },
]

</script>

<style>

.table-el {
    padding-top: 48px;
    height: 100vh;
    position: absolute;
    top: 0px;
}

</style>