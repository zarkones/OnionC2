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
                <p>{{ mockedCurrentDirPath }}</p>
                
                <v-toolbar-items>
                </v-toolbar-items>

                <v-spacer />

                <v-btn
                    icon="mdi-close"
                    @click="dialog = false"
                ></v-btn>
            </v-toolbar>

            <div class="pl-4 pr-4">
                <v-row style="min-height: 480px;">
                    <v-col cols="6">
                        <div style="height: 50%;">
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

                            <div
                                v-for="(record, recordIndex) in mockedDirContent"
                            >
                                <v-icon :icon="record.isDir ? 'mdi-folder' : 'mdi-file-outline'" class="mr-2"></v-icon>
                                <span>{{ record.name }}</span>
                                <span>{{ record.timestamp }}</span>
                                <v-btn
                                    variant="plain"
                                    density="compact"
                                    @click="agentToServerUpload(record)"
                                    style="float: right;"
                                >
                                    <v-icon icon="mdi-arrow-right-bold"></v-icon>

                                    <v-tooltip activator="parent" location="top" open-delay="600">
                                        <div class="tooltip-el">
                                            <h4>Upload File To Remote System</h4>
                                            Upload a file to the command and control server. Once a file
                                            is uploaded then you can order an agent to download it.
                                        </div>
                                    </v-tooltip>
                                </v-btn>
                                <v-divider v-if="recordIndex !== mockedDirContent.length-1" />
                            </div>
                        </div>

                        <h4>Files Uploaded By Agent:</h4>
                        <div style="height: 50%;">
                            <div
                            v-for="(record, recordIndex) in mockedDirContent"
                        >
                            <v-icon :icon="record.isDir ? 'mdi-folder' : 'mdi-file-outline'" class="mr-2"></v-icon>
                                <span>{{ record.name }}</span>
                                <span>{{ record.timestamp }}</span>
                                <v-btn
                                    variant="plain"
                                    density="compact"
                                    @click="agentToServerUpload(record)"
                                    style="float: right;"
                                >
                                    <v-icon icon="mdi-download"></v-icon>

                                    <v-tooltip activator="parent" location="top" open-delay="600">
                                        <div class="tooltip-el">
                                            <h4>Download File</h4>
                                            Download a file from the C2 server to your machine through the browser.
                                            This file was previously uploaded by the agent onto the C2 server.
                                        </div>
                                    </v-tooltip>
                                </v-btn>
                                <v-divider v-if="recordIndex !== mockedDirContent.length-1" />
                            </div>
                        </div>
                    </v-col>

                    <v-col>
                        <div class="d-flex">
                            <h4>Remote File System:</h4>
                            <v-spacer />
                            <v-btn
                                class="ml-4"
                                variant="plain"
                                density="compact"
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

                        <div
                            v-for="(record, recordIndex) in mockedDirContent"
                        >
                            <v-btn
                                variant="plain"
                                density="compact"
                                @click="agentToServerUpload(record)"
                            >
                                <v-icon icon="mdi-arrow-left-bold"></v-icon>

                                <v-tooltip activator="parent" location="top" open-delay="600">
                                    <div class="tooltip-el">
                                        <h4>Upload File To Server</h4>
                                        Upload a file from agent's system onto the C2 server.
                                        After it's uploaded to the C2 server, then you can download
                                        it from your browser.
                                    </div>
                                </v-tooltip>
                            </v-btn>
                            <v-icon :icon="record.isDir ? 'mdi-folder' : 'mdi-file-outline'" class="mr-2"></v-icon>
                            <span>{{ record.name }}</span>
                            <v-spacer />
                            <span>{{ record.timestamp }}</span>
                            <v-divider v-if="recordIndex !== mockedDirContent.length-1" />
                        </div>
                    </v-col>
                </v-row>
            </div>

            <v-card-actions>
                <v-btn
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
                </v-btn>
            </v-card-actions>

        </v-card>

    </v-dialog>
</template>

<script lang="ts" setup>

const dialog = ref(false)

const mockedCurrentDirPath = 'C:\\Users\\user\\Desktop'

const mockedDirContent = [
    { timestamp: '3h 13min ago', name: 'documents', isDir: true },
    { timestamp: '3h 13min ago', name: 'pics_yr2016', isDir: true },
    { timestamp: '3h 13min ago', name: 'Minecraft.exe', isDir: false },
    { timestamp: '3h 13min ago', name: 'PhotoShop.exe', isDir: false },
    { timestamp: '3h 13min ago', name: 'profilepic.png', isDir: false },
    { timestamp: '3h 13min ago', name: 'veryimportant.csv', isDir: false },
]

const selected = async () => {

}

const agentToServerUpload = async (record: FileRecord) => {

}

const apply = async () => {

}

</script>