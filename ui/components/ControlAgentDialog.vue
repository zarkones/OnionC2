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
            >
                <v-icon icon="mdi-console" />

                <v-tooltip activator="parent" location="top" open-delay="600">
                    <div class="tooltip-el">
                        <h4>todo</h4>
                        // TODO
                    </div>
                </v-tooltip>
            </v-btn>
        </template>

        <v-card class="liquid-glass" density="compact" >
            <v-toolbar density="compact" class="liquid-glass">
                <v-toolbar-items>
                    <v-btn
                        variant="plain"
                        density="compact"
                    >
                        <v-icon icon="mdi-folder-arrow-up-down" />

                        <v-tooltip activator="parent" location="top" open-delay="600">
                            <div class="tooltip-el">
                                <h4>File Management</h4>
                                Browse files of the remote system with an ability to upload
                                and download files.
                            </div>
                        </v-tooltip>
                    </v-btn>

                    <v-btn
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
                    </v-btn>

                    <v-btn
                        variant="plain"
                        density="compact"
                    >
                        <v-icon icon="mdi-clipboard-arrow-down" />

                        <v-tooltip activator="parent" location="top" open-delay="600">
                            <div class="tooltip-el">
                                <h4>Grab Clipboard</h4>
                                Grab text copied into the clipboard.
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

            <div class="terminal-output ma-2">
                <p>
                    {{ mockedTerminalOutput }}
                </p>
            </div>

            <v-card-actions class="terminal-actions">
                <v-text-field
                    variant="outlined"
                    density="compact"
                    placeholder="Send Command"
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

const dialog = shallowRef(false)
const awaitExecution = ref(true)

const toggleExecutionAwaiting = () => {
    awaitExecution.value = !awaitExecution.value
}

const mockedTerminalOutput = `[INFO] 2025-06-28 13:19:23 Starting application bootstrap...
[DEBUG] Loading configuration from /etc/app/config.yaml
[DEBUG] Environment variables set: APP_ENV=production, LOG_LEVEL=debug
[INFO] Initializing database connection...
[DEBUG] Connecting to PostgreSQL at localhost:5432
[DEBUG] Connection pool size: 10, timeout: 30s
[SUCCESS] Database connection established
[INFO] Starting user authentication module
[DEBUG] JWT secret validated
[WARNING] Deprecated API call detected in auth module (v1.2.3). Upgrade to v2.0.0 recommended
[INFO] Starting HTTP server on port 8080
[DEBUG] Middleware stack: [cors, logger, auth, rate_limiter]
[INFO] Server ready, accepting connections
[DEBUG] Incoming request: GET /api/v1/health
[DEBUG] Response: 200 OK, body: {"status": "healthy", "uptime": "0h0m12s"}
[INFO] Processing batch job #12345
[DEBUG] Fetching data from external API: https://api.example.com/v1/data
[DEBUG] API response: 200 OK, 2.1 MB received
[DEBUG] Parsing JSON payload...
[SUCCESS] Batch job #12345 completed in 1.234s
[ERROR] Failed to sync data for user_id=789: Timeout after 5s
[DEBUG] Retrying sync for user_id=789 (attempt 1/3)
[INFO] Sync retry successful for user_id=789
[INFO] Processing batch job #12346
[DEBUG] Validating input data: 5000 records
[DEBUG] Loading user data for 1024 users
[INFO] User data cache warmed up in 0.342s
[WARNING] Deprecated API call detected in auth module (v1.2.3). Upgrade to v2.0.0 recommended
[INFO] Starting HTTP server on port 8080
[DEBUG] Middleware stack: [cors, logger, auth, rate_limiter]
[INFO] Server ready, accepting connections
[DEBUG] Incoming request: GET /api/v1/health
[WARNING] Deprecated API call detected in auth module (v1.2.3). Upgrade to v2.0.0 recommended
[INFO] Starting HTTP server on port 8080
[DEBUG] Middleware stack: [cors, logger, auth, rate_limiter]
[INFO] Server ready, accepting connections
[DEBUG] Incoming request: GET /api/v1/health
[DEBUG] Response: 200 OK, body: {"status": "healthy", "uptime": "0h0m12s"}
[INFO] Processing batch job #12345
[DEBUG] Fetching data from external API: https://api.example.com/v1/data
[DEBUG] API response: 200 OK, 2.1 MB received
[DEBUG] Parsing JSON payload...
[SUCCESS] Batch job #12345 completed in 1.234s
[ERROR] Failed to sync data for user_id=789: Timeout after 5s
[DEBUG] Retrying sync for user_id=789 (attempt 1/3)
[INFO] Sync retry successful for user_id=789
[INFO] Processing batch job #12346
[DEBUG] Validating input data: 5000 records
[DEBUG] Response: 200 OK, body: {"status": "healthy", "uptime": "0h0m12s"}
[INFO] Processing batch job #12345
[DEBUG] Fetching data from external API: https://api.example.com/v1/data
[DEBUG] API response: 200 OK, 2.1 MB received
[DEBUG] Parsing JSON payload...
[SUCCESS] Batch job #12345 completed in 1.234s
[ERROR] Failed to sync data for user_id=789: Timeout after 5s
[DEBUG] Retrying sync for user_id=789 (attempt 1/3)
[INFO] Sync retry successful for user_id=789
[INFO] Processing batch job #12346
[DEBUG] Validating input data: 5000 records
`
</script>

<style scoped>
.terminal-output {
    height: 100%;
    max-height: 100%;
    overflow-y: auto;
    white-space: pre-line;
}

.terminal-actions {
    padding-bottom: 0px;
}
</style>