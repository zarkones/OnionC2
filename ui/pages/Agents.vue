<template>
    <v-data-table
        :headers="(headers as any)"
        :items="mockedTableData"
        density="compact"
        item-key="name"
        class="table-el"
    >
        <template v-slot:item.actions="{ item }">
            <v-btn
                variant="plain"
                density="compact"
            >
                <v-icon icon="mdi-console" />

                <v-tooltip activator="parent" location="top" open-delay="600">
                    <div class="tooltip-el">
                        <h4>Channel Creation</h4>
                        Create a new chat channel.
                    </div>
                </v-tooltip>
            </v-btn>
        </template>
    </v-data-table>
</template>

<script setup lang="ts">

const headers = [
    // { title: 'ID', align: 'start', key: 'id' },
    { title: 'IP', align: 'start', key: 'ip' },
    { title: 'Hostname', align: 'start', key: 'hostname' },
    { title: 'OS', align: 'start', key: 'os' },
    { title: 'CPU Name', align: 'start', key: 'cpuName' },
    { title: 'CPU Arch', align: 'start', key: 'cpuArch' },
    { title: 'RAM', align: 'start', key: 'ram' },
    { title: 'Last Seen', align: 'start', key: 'lastSeen' },
    { title: 'Actions', align: 'end', key: 'actions' },
]

const mockedTableData = [
    { 
        ip: '192.168.1.100', 
        hostname: 'DESKTOP-7X9P2Q1', 
        os: 'Windows 10', 
        cpuName: 'i7-12700H', 
        cpuArch: 'amd64', 
        ram: '16GB', 
        lastSeen: '2025-06-26T12:34:56Z' 
    },
    { 
        ip: '192.168.1.101', 
        hostname: 'LAPTOP-K4M8N3V', 
        os: 'Windows 11', 
        cpuName: 'Ryzen 7 5800X', 
        cpuArch: 'amd64', 
        ram: '32GB', 
        lastSeen: '2025-06-26T15:22:10Z' 
    },
    { 
        ip: '192.168.1.102', 
        hostname: 'WORKSTATION-ZT5R8Y', 
        os: 'Windows 10', 
        cpuName: 'i5-11400', 
        cpuArch: 'amd64', 
        ram: '8GB', 
        lastSeen: '2025-06-26T09:10:30Z' 
    },
    { 
        ip: '192.168.1.103', 
        hostname: 'PC-9H2J4W8', 
        os: 'Windows 7', 
        cpuName: 'i7-9700K', 
        cpuArch: 'amd64', 
        ram: '16GB', 
        lastSeen: '2025-06-25T23:45:12Z' 
    },
    { 
        ip: '192.168.1.104', 
        hostname: 'DESKTOP-P3Q9M7X', 
        os: 'Windows 11', 
        cpuName: 'Ryzen 5 5600G', 
        cpuArch: 'amd64', 
        ram: '16GB', 
        lastSeen: '2025-06-26T18:00:00Z' 
    },
    { 
        ip: '192.168.1.105', 
        hostname: 'LAPTOP-6T2R5Y9', 
        os: 'Windows 10', 
        cpuName: 'i9-12900K', 
        cpuArch: 'amd64', 
        ram: '64GB', 
        lastSeen: '2025-06-26T10:15:22Z' 
    },
    { 
        ip: '192.168.1.106', 
        hostname: 'PC-X7Y8Z9W', 
        os: 'Windows 11', 
        cpuName: 'Ryzen 9 5900X', 
        cpuArch: 'amd64', 
        ram: '32GB', 
        lastSeen: '2025-06-26T14:30:45Z' 
    },
    { 
        ip: '192.168.1.107', 
        hostname: 'DESKTOP-4K5N2M8', 
        os: 'Windows 10', 
        cpuName: 'i7-13700F', 
        cpuArch: 'amd64', 
        ram: '16GB', 
        lastSeen: '2025-06-26T08:20:17Z' 
    },
    { 
        ip: '192.168.1.108', 
        hostname: 'WORKSTATION-8N3M7P', 
        os: 'Windows 7', 
        cpuName: 'Ryzen 3 5300U', 
        cpuArch: 'amd64', 
        ram: '8GB', 
        lastSeen: '2025-06-25T20:50:33Z' 
    },
    { 
        ip: '192.168.1.109', 
        hostname: 'LAPTOP-2Q9W4R6', 
        os: 'Windows 11', 
        cpuName: 'i5-1240P', 
        cpuArch: 'arm64', 
        ram: '16GB', 
        lastSeen: '2025-06-26T21:00:00Z' 
    },
]

</script>

<style>

.table-el {
    padding-top: 48px;
    height: 100%;
}

</style>