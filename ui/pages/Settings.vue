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

                    </v-tabs-window>
                </v-col>
            </v-row>

        </v-card>
    </v-container>
</template>

<script setup lang="ts">
const tabs = {
    Authentication: { key: 'Authentication', icon: 'mdi-account' },
} as const

const tab = ref('')
const privateKeyHexPem = ref('')
const keyParsingErrored = ref(false)

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

<style>

.settings-container {
    margin-top: auto;
    margin-bottom: auto;
}

.settings-tab-container {
    padding: 16px;
}

</style>