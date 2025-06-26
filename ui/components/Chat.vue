<template>
    <div class="liquid-glass" style="width: 100%; height: 100%;">

        <div class="chat-actions-bar">
            <v-btn
                variant="plain"
                density="compact"
            >
                <v-icon icon="mdi-chat-plus" />

                <v-tooltip activator="parent" location="top" open-delay="600">
                    <div class="tooltip-el">
                        <h4>Channel Creation</h4>
                        Create a new chat channel.
                    </div>
                </v-tooltip>
            </v-btn>

            <v-btn
                variant="plain"
                density="compact"
            >
                <v-icon icon="mdi-chat-minus" />

                <v-tooltip activator="parent" location="top" open-delay="600">
                    <div class="tooltip-el">
                        <h4>Channel Removal</h4>
                        Removes a chat channel. This would also delete all messages of that channel form the C2 server.
                    </div>
                </v-tooltip>
            </v-btn>

            <v-btn
                variant="plain"
                density="compact"
            >
                <v-icon icon="mdi-account-plus" />

                <v-tooltip activator="parent" location="top" open-delay="600">
                    <div class="tooltip-el">
                        <h4>Invite To Channel</h4>
                        Invite an operator to a chat channel.
                    </div>
                </v-tooltip>
            </v-btn>

            <v-btn
                variant="plain"
                density="compact"
            >
                <v-icon icon="mdi-account-minus" />

                <v-tooltip activator="parent" location="top" open-delay="600">
                    <div class="tooltip-el">
                        <h4>Remove From Channel</h4>
                        Remove an operator from a chat channel. This would trigger
                        rotation of channel's encryption secret.
                    </div>
                </v-tooltip>
            </v-btn>

            <v-btn
                variant="plain"
                density="compact"
            >
                <v-icon icon="mdi-delete" />

                <v-tooltip activator="parent" location="top" open-delay="600">
                    <div class="tooltip-el">
                        <h4>Delete Messages</h4>
                        Delete messages of a specific chat.
                    </div>
                </v-tooltip>
            </v-btn>

            <v-btn
                variant="plain"
                density="compact"
            >
                <v-icon icon="mdi-timer-remove" />

                <v-tooltip activator="parent" location="top" open-delay="600">
                    <div class="tooltip-el">
                        <h4>Time-based Messages Deletion</h4>
                        Configure periodic deletion of messages.
                    </div>
                </v-tooltip>
            </v-btn>

            <v-spacer />
        </div>

        <v-divider />

        <v-row class="full-height">
            <v-col cols="4" class="channels-col">
                <div class="categories-wrapper">
                    <div v-for="category in channelCategories" :key="category" class="chat-expansion-panel">
                        <h3 class="chat-category-label">{{ category.toUpperCase() }}</h3>
                        
                        <div class="chat-channels">
                            <v-btn
                                v-for="channel in channels[category]" 
                                :key="channel" 
                                density="compact"
                                variant="text"
                                class="justify-start"
                            >
                                {{ channel }}
                            </v-btn>
                        </div>
                    </div>
                </div>
            </v-col>
            
            <v-col class="channel-messages-container">
                <div class="messages-wrapper">
                    <div v-for="(msg, msgIndex) in mockedMessages" :key="msgIndex" class="message">
                        <h4>{{ msg.name }}</h4>
                        <p>{{ msg.content }}</p>
                    </div>
                </div>
                <v-text-field
                    variant="outlined"
                    density="compact"
                    placeholder="Send Message"
                    class="mt-3"
                />
            </v-col>
        </v-row>

    </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

const mockedMessages = [
    { name: 'Mulder', content: 'Scully, I just found a case file about a small town where people vanish during full moons. Sounds like werewolves to me.' },
    { name: 'Scully', content: 'Mulder, there’s no scientific evidence for werewolves. It’s probably just a coincidence or a local prank.' },
    { name: 'Scully', content: 'Have you checked for environmental factors? Maybe there’s a pattern tied to weather or local wildlife.' },
    { name: 'Mulder', content: 'Wildlife? Scully, the last guy who disappeared left behind claw marks and fur. Explain that with science.' },
    { name: 'Scully', content: 'I’d need to see lab results on that fur before jumping to "werewolf." Could be a bear.' },
    { name: 'Scully', content: 'I’m ordering forensic tests on those samples. We’ll know more in 48 hours.' },
    { name: 'Scully', content: 'And Mulder, don’t go chasing shadows in the woods without me.' },
    { name: 'Mulder', content: 'Wouldn’t dream of it, Scully. But bring a silver bullet, just in case.' },
    { name: 'Mulder', content: 'Scully, I just found a case file about a small town where people vanish during full moons. Sounds like werewolves to me.' },
    { name: 'Scully', content: 'Mulder, there’s no scientific evidence for werewolves. It’s probably just a coincidence or a local prank.' },
    { name: 'Scully', content: 'Have you checked for environmental factors? Maybe there’s a pattern tied to weather or local wildlife.' },
    { name: 'Mulder', content: 'Wildlife? Scully, the last guy who disappeared left behind claw marks and fur. Explain that with science.' },
    { name: 'Scully', content: 'I’d need to see lab results on that fur before jumping to "werewolf." Could be a bear.' },
    { name: 'Scully', content: 'I’m ordering forensic tests on those samples. We’ll know more in 48 hours.' },
    { name: 'Scully', content: 'And Mulder, don’t go chasing shadows in the woods without me.' },
    { name: 'Mulder', content: 'Wouldn’t dream of it, Scully. But bring a silver bullet, just in case.' },
]

const channels = {
    general: ['scripts', 'evasion', 'data-analysis'],
    devops: ['pipeline', 'pentesting'],
} as const

type ChannelKey = keyof typeof channels

const channelCategories = computed(() => Object.keys(channels) as ChannelKey[])
</script>

<style scoped>
.chat-channels {
    display: flex;
    flex-direction: column;
}

.chat-expansion-panel {
    background: none;
}

.full-height {
    height: 100%;
}

.channel-messages-container {
    display: flex;
    flex-direction: column;
    height: calc(100% - 40px);
    padding-top: 32px;
    padding-right: 32px;
    padding-bottom: 0px;
}

.messages-wrapper {
    flex-grow: 1;
    overflow-y: auto;
}

.message {
    margin-bottom: 16px;
}

.chat-category-label {
    margin: 16px;
}

/* New styles for channel categories scrolling */
.channels-col {
    display: flex;
    flex-direction: column;
    height: 100%;
}

.categories-wrapper {
    flex-grow: 1;
    overflow-y: auto;
}

.chat-actions-bar {
    display: flex;
    flex-direction: row;
    gap: 8px;
    margin: 16px;
}
</style>