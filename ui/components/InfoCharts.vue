<template>
    <v-row>
        <v-col>
            <span>Operating Systems</span>
            <div> <Pie :data="dataOs" :options="options" /> </div>
        </v-col>
        <v-col>
            <span>CPU Architectures</span>
            <div> <Pie :data="dataArchs" :options="options" /> </div>
        </v-col>
        <v-col>
            <span>CPU Brands</span>
            <div> <Pie :data="dataCpus" :options="options" /> </div>
        </v-col>
        <v-col>
            <Line :data="dataStats" :options="options" />
        </v-col>
    </v-row>
</template>

<script setup lang="ts">

import { Chart as ChartJS, ArcElement, Tooltip, Legend, CategoryScale, LinearScale, PointElement, LineElement, Title } from 'chart.js'
import { Pie, Line } from 'vue-chartjs'

const mocks = {
    cpus: {
        labels: ['AMD Ryzen 5 7600', 'Intel Core i5-13400F', 'AMD Ryzen 7 5700X', 'AMD Ryzen 9 9950X', 'Intel Core Ultra 9 285K', 'AMD Ryzen 9 9950X3D', 'AMD Ryzen Threadripper 7970X'],
        data: [40, 20, 80, 10, 34, 86]
    },
    archs: {
        labels: ['arm64', 'amd64'],
        data: [12, 120]
    },
    os: {
        labels: ['Windows 10', 'Windows 11', 'Windows 7'],
        data: [40, 60, 7]
    },
    stats: {
        labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
        data: [1, 4, 6, 9, 12, 15, 17]
    },
}

const backgroundColors = [
    // '#0D47A1',
    // '#1565C0',
    // '#1976D2',
    // '#1E88E5',
    // '#42A5F5',
    // '#64B5F6',

    // '#0D47A1',
    // '#1565C0',
    // '#1976D2',
    // '#1E88E5',
    // '#42A5F5',
    // '#64B5F6',

    // '#0D47A1',
    // '#1565C0',
    // '#1976D2',
    // '#1E88E5',
    // '#42A5F5',
    // '#64B5F6',

    '#263238',
    '#37474F',
    '#455A64',
    '#546E7A',
    '#607D8B',
    '#78909C',

    '#263238',
    '#37474F',
    '#455A64',
    '#546E7A',
    '#607D8B',
    '#78909C',

    '#263238',
    '#37474F',
    '#455A64',
    '#546E7A',
    '#607D8B',
    '#78909C',
]

ChartJS.register(
    ArcElement,
    Tooltip,
    Legend,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
)

const dataCpus = {
    labels: mocks.cpus.labels,
    datasets: [
        {
            backgroundColor: backgroundColors,
            borderWidth: 0,
            data: mocks.cpus.data,
        }
    ]
}

const dataArchs = {
    labels: mocks.archs.labels,
    datasets: [
        {
            backgroundColor: backgroundColors,
            borderWidth: 0,
            data: mocks.archs.data,
        }
    ]
}

const dataOs = {
    labels: mocks.os.labels,
    datasets: [
        {
            backgroundColor: backgroundColors,
            borderWidth: 0,
            data: mocks.os.data,
        },
    ]
}

const dataStats = {
    labels: mocks.stats.labels,
    datasets: [
        {
            label: 'data one',
            backgroundColor: backgroundColors,
            borderColor: '#263238',
            data: mocks.stats.data,
        }
    ]
}

const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
        legend: {
            display: false
        }
    }
}


</script>