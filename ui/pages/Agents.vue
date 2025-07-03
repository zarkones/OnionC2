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
        <template v-slot:top>
            <div class="pl-4 pt-4 pr-4 liquid-glass">
                <v-select
                    v-model="API.store.origins.selected"
                    :items="API.store.origins.data"
                    item-title="n"
                    item-value="i"
                    label="Select Origins"
                    multiple
                    density="compact"
                    variant="outlined"
                    v-on:update:model-value="originFilter"
                >
                    <template v-slot:selection="{ item, index }">
                        <v-chip :text="item.title"></v-chip>
                        <!-- <span
                            v-if="index === 2"
                            class="text-grey text-caption align-self-center"
                        >
                            (+{{ selectedOrigins.length - 2 }} others)
                        </span> -->
                    </template>
                </v-select>
            </div>
        </template>

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

const route = useRoute()

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

const originFilter = async () => {
    await API.value.generalUpdate()
}

onMounted(async () => {
    try {
        const origins = (route.query['origins'] as string).split(',')
        if (origins.length > 0) {
            API.value.store.origins.selected = origins
        }
        await API.value.generalUpdate()
    } catch(e) {}
})

</script>

<style>

.table-el {
    padding-top: 48px;
    height: 100vh;
    position: absolute;
    top: 0px;
}

</style>