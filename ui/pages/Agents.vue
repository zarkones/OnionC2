<template>
    <MapSvg style="position: absolute;"></MapSvg>

    <v-data-table-virtual
        :headers="(headers as any)"
        :items="API.store.agents.data"
        density="compact"
        item-key="name"
        class="table-el liquid-glass"
        fixed-header
    >
        <template v-slot:item.lastSeen="{ item }">
            {{ formatUnixNanoTime(item.lastSeen) }}
        </template>

        <template v-slot:item.actions="{ item }">
            <ControlAgentDialog
                :agentId="item.id"
            />
        </template>

    </v-data-table-virtual>
</template>

<script setup lang="ts">

const formatUnixNanoTime = (nanoStr: number) => {
    const nano = BigInt(nanoStr);
    const milli = Number(nano / 1000000n);
    const date = new Date(milli);
    const now = new Date();
    // @ts-ignore
    const diff = now - date;  // difference in milliseconds
    const diffSeconds = Math.floor(diff / 1000);

    if (diffSeconds < 60) {
        return `${diffSeconds}s ago`;
    } else if (diffSeconds < 3600) {
        const minutes = Math.floor(diffSeconds / 60);
        return `${minutes}min ago`;
    } else if (diffSeconds < 86400) {
        const hours = Math.floor(diffSeconds / 3600);
        const minutes = Math.floor((diffSeconds % 3600) / 60);
        return `${hours}h ${minutes}min ago`;
    } else {
        const months = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
        const day = date.getDate();
        const month = months[date.getMonth()];
        const year = date.getFullYear();
        const hours = date.getHours().toString().padStart(2, '0');
        const minutes = date.getMinutes().toString().padStart(2, '0');
        return `${day} ${month} ${year} / ${hours}:${minutes}`;
    }
}

const headers = [
    // { title: 'ID', align: 'start', key: 'id' },
    { title: 'Country', align: 'start', key: 'country' },
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