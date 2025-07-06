<template>
    <v-dialog
        v-model="dialog"
        transition="dialog-bottom-transition"
        :scrim="false"
    >

        <template v-slot:activator="{ props: activatorProps }">
            <v-btn
                variant="plain"
                density="compact"
                v-bind="activatorProps"
                @click.stop="clicked"
            >
                <v-icon icon="mdi-file-search" />

                <v-tooltip activator="parent" location="top" open-delay="600">
                    <div class="tooltip-el">
                        <h4>Search File System</h4>
                        Search the file system for files and directories matching on of specified strings.
                    </div>
                </v-tooltip>
            </v-btn>
        </template>

        <v-card variant="outlined" class="liquid-glass" density="compact">
            <v-toolbar density="compact" class="liquid-glass" title="File Search">
                <v-toolbar-items>
                </v-toolbar-items>

                <v-spacer />

                <v-btn
                    icon="mdi-close"
                    @click="dialog = false"
                ></v-btn>
            </v-toolbar>

            <div class="pl-4 pr-4" style="height: 480px;">
                <v-text-field
                    v-model="startingDir"
                    variant="outlined"
                    density="compact"
                    label="Starting Directory's Path"
                />

                <div style="max-height: 450px; overflow-y: auto;">
                    <div
                        v-for="(_term, termIndex) in searchTerms"
                        :key="termIndex"
                    >
                        <v-text-field
                            :ref="(el) => (searchTermFields[termIndex] = el as HTMLElement)"
                            v-model="searchTerms[termIndex]"
                            variant="outlined"
                            density="compact"
                            placeholder="Search Term"
                            @keyup.enter="nextTerm"
                        />
                    </div>
                </div>

                <v-card-actions>
                    <v-spacer />

                    <v-btn
                        @click="startSearching"
                        variant="plain"
                        density="compact"
                        :disabled="!searchTerms.length || !searchTerms[0]?.length || !startingDir.length"
                    >
                        Start Searching

                        <v-tooltip activator="parent" location="top" open-delay="600">
                            <div class="tooltip-el">
                                <h4>Start File Search</h4>
                                Send a command to the agent to search files and directories which
                                match the specified conditions.
                            </div>
                        </v-tooltip>
                    </v-btn>
                </v-card-actions>
            </div>
        </v-card>
    </v-dialog>
</template>

<script setup lang=ts>

const dialog = ref(false)

const props = defineProps<{
    agentId: string
}>()

const searchTerms = ref([] as string[])
const startingDir = ref('')

const searchTermFields = ref({} as Record<number, HTMLElement>)

const nextTerm = async () => {
    if (searchTerms.value[searchTerms.value.length-1]?.length === 0) {
        return
    }
    searchTerms.value.push('')
    await nextTick()
    searchTermFields.value[Object.keys(searchTermFields.value).length-1]?.focus()
}

const startSearching = async () => {
    await API.value.sendMessage(props.agentId, `/find-files|${startingDir.value}|${searchTerms.value.filter(t => t.length).join(',')}`)
    dialog.value = false
    await API.value.generalUpdate()
}

const clicked = () => {
    searchTerms.value = ['']
}

</script>