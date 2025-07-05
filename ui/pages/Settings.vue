<template>
    <v-container
        class="settings-container"
    >
        <v-card style="height: 100%;">

            <v-row>
                <v-col cols="3">
                    <v-tabs
                        v-model="tab"
                        color="primary"
                        direction="vertical"
                    >
                        <v-tab
                            v-for="(category) in Object.values(tabs)"
                            :key="category.key"
                            :prepend-icon="category.icon" 
                            :text="category.key"
                            :value="category.key"
                        />
                    </v-tabs>
                </v-col>

                <v-col>
                    <v-tabs-window v-model="tab">

                        <v-tabs-window-item class="settings-tab-container" :value="tabs.Authentication.key">
                            <v-text-field
                                v-model="API.c2HostURL"
                                variant="outlined"
                                density="compact"
                                label="C2 Host URL"
                            />

                            <v-text-field
                                v-model="API.username"
                                variant="outlined"
                                density="compact"
                                label="Username"
                            />

                            <v-textarea
                                v-model="privateKeyHexPem"
                                @update:model-value="updatePrivateKey"
                                variant="outlined"
                                density="compact"
                                :error="keyParsingErrored"
                                label="Operator's Private Key (Hex Encoded)"
                            />
                        </v-tabs-window-item>

                        <v-tabs-window-item class="settings-tab-container" :value="tabs.Operators.key">
                            <v-select
                                v-model="selectedOperatorUsername"
                                :items="API.store.operators.data"
                                item-title="username"
                                item-value="username"
                                label="Select Operator"
                                density="compact"
                                variant="outlined"
                                v-on:update:model-value="switchOperator"
                            >
                            </v-select>

                            <div v-if="selectedOperator && Object.keys(selectedOperator as any).length > 0">
                                <v-expansion-panels variant="accordion">
                                    <v-expansion-panel
                                        title="Show Public Key"
                                        :text="selectedOperator?.publicKeyHex"
                                    />
                                </v-expansion-panels>

                                <v-row class="mt-2">
                                    <v-col cols="2"><h4>Created At:</h4></v-col>
                                    <v-col>{{ formatUnixNanoTime(selectedOperator?.createdAt as bigint) }}</v-col>
                                </v-row>

                                <!-- {{ permissions }} -->

                                <!-- <v-row class="mt-4"
                                    v-for="p in permissions"
                                    :key="p.id"
                                >
                                    <v-col cols="2">{{ PERMISSIONS[p.key] }}</v-col>
                                </v-row> -->
                                <v-data-table-virtual
                                    :headers="(headers as any)"
                                    :items="Object.keys(permissions)"
                                    density="compact"
                                    item-key="id"
                                    class="liquid-glass mt-4"
                                    fixed-header
                                >
                                    <template v-slot:item.key="{ item }">
                                        {{ PERMISSIONS[permissions[item].key] }}
                                    </template>

                                    <template v-slot:item.acquired="{ item }">
                                        <v-switch class="mt-5" color="primary" v-model="permissions[item].acquired" label="Switch"></v-switch>
                                    </template>
                                    
                                    <template v-slot:item.createdAt="{ item }">
                                        {{ permissions[item].acquired === true ? formatUnixNanoTime(permissions[item].createdAt) : '' }}
                                    </template>

                                    <template v-slot:item.actions="{ item }">
                                        <v-btn density="compact" variant="plain">
                                            <!-- // TODO: -->
                                            <v-icon icon="mdi-delete"></v-icon>
                                        </v-btn>
                                    </template>
                                </v-data-table-virtual>
                            </div>
                        </v-tabs-window-item>

                    </v-tabs-window>
                </v-col>
            </v-row>

        </v-card>
    </v-container>
</template>

<script setup lang="ts">
const tabs = {
    Authentication: { key: 'Authentication', icon: 'mdi-cog' },
    Operators: { key: 'Operators', icon: 'mdi-account-group' },
} as const

const tab = ref('')
const privateKeyHexPem = ref('')
const keyParsingErrored = ref(false)
const selectedOperatorUsername = ref('')
const selectedOperator: Ref<Operator | null> = ref(null)
const permissions = ref({} as GetPermissionsRespCtx)

const headers = [
    { title: 'Key', align: 'start', key: 'key' },
    { title: 'Acquired', align: 'start', key: 'acquired' },
    { title: 'Created At', align: 'start', key: 'createdAt' },
    { title: 'Actions', align: 'end', key: 'actions' },
]

const switchOperator = async () => {
    try {
    selectedOperator.value = API.value.store.operators.data.filter(o => o.username === selectedOperatorUsername.value)[0] as Operator
    } catch(e) {
        console.error('could not find stored operator with username of:', selectedOperatorUsername.value)
        return
    }

    permissions.value = await API.value.fetchPermissions(selectedOperator.value.username)
    console.log(permissions.value)
}

const updatePrivateKey = async () => {
    try {
        // @ts-ignore
        const bytes = new Uint8Array(privateKeyHexPem.value.match(/.{1,2}/g).map(byte => parseInt(byte, 16)))
        const decoder = new TextDecoder('utf-8')
        const pemEncodedKey = decoder.decode(bytes)
        await API.value.setPrivateKey(pemEncodedKey)
        keyParsingErrored.value = false
    } catch(e) {
        console.error('failed to parse private key from hex:', e)
        keyParsingErrored.value = true
    }
}

</script>

<style scoped>

.settings-container {
    margin-top: auto;
    margin-bottom: auto;
}

.settings-tab-container {
    padding: 16px;
}

</style>