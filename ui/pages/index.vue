<template>
    <Map></Map>
    <div v-show="showChat" ref="chatDiv" class="chat-container">
        <Chat />
    </div>
    <div ref="infoCharts" class="info-charts liquid-glass">
        <div class="actions">
            <v-btn @click="toggleShowStats" density="compact" variant="plain">
                <v-icon v-if="hideStats" icon="mdi-chevron-up" />
                <v-icon v-else icon="mdi-chevron-down" />
            </v-btn>

            <UnknownOriginStats />

            <v-spacer />

            <v-btn @click="toggleShowChat" density="compact" variant="plain">
                <v-icon v-if="showChat" icon="mdi-chat" />
                <v-icon v-else icon="mdi-chat-outline" />

                <v-tooltip activator="parent" location="top" open-delay="600">
                    <div class="tooltip-el">
                        <h4>E2EE Chat</h4>
                        Opens/Closes end-to-end encrypted chat. This chat is used to securely communicate
                        with other operators. Not even the C2 server is able to read chat's content.
                    </div>
                </v-tooltip>
            </v-btn>
        </div>
        
        <v-divider v-show="!hideStats" class="mb-1 mt-3" />
        <InfoCharts v-show="!hideStats" />
    </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

const hideStats = ref(false)
const showChat = ref(false)

const toggleShowStats = () => {
    hideStats.value = !hideStats.value
}

const toggleShowChat = () => {
    showChat.value = !showChat.value
}

const chatDiv: Ref<HTMLElement | null> = ref(null);
const infoCharts: Ref<HTMLElement | null> = ref(null);

const adjustChatHeight = (chat: HTMLElement, infoCharts: HTMLElement) => {
    const topBarHeight = 48

    const observer = new ResizeObserver(() => {
        // @ts-ignore
        const height = infoCharts.offsetHeight
        chat.style.height = `calc(100vh - ${height + topBarHeight}px)`
    })

    // Start observing the .info-charts div
    observer.observe(infoCharts)

    // Set the initial height
    // @ts-ignore
    const height = infoCharts.offsetHeight
    chat.style.height = `calc(100vh - ${height}px)`
}

onMounted(() => {
    if (chatDiv.value && infoCharts.value) {
        adjustChatHeight(chatDiv.value, infoCharts.value)
    }
})

</script>

<style scoped>
.info-charts {
    position: absolute;
    left: 0px;
    width: 100%;
    height: fit-content;
    bottom: 0px;
    padding: 8px;
}

.actions {
    display: flex;
    gap: 16px;
}

.chat-container {
    position: absolute;
    width: 660px;
    right: 0px;
    top: 48px;
}
</style>