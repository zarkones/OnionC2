<template>
    <div id="map" ref="mapContainer" @wheel="handleWheel" @mousedown="handleMouseDown">
        <div class="map-container" :style="transformStyle">
            <MapSvg />
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from 'vue'

// Reactive state
const zoom = ref(1)
const panX = ref(0)
const panY = ref(0)
const isDragging = ref(false)
const startX = ref(0)
const startY = ref(0)
const startPanX = ref(0)
const startPanY = ref(0)
const mapContainer = ref(null)
const markers = ref([])

// Map dimensions (adjust these to match your map image)
const mapWidth = 2050  // Original width of the map image in pixels
const mapHeight = 2317/2 // Original height of the map image in pixels

// Computed style for the map container
const transformStyle = computed(() => ({
    transform: `translate(${panX.value}px, ${panY.value}px) scale(${zoom.value})`,
    transformOrigin: '0 0',
}))

// Clamp pan values to keep the map within bounds
const clampPan = () => {
    if (!mapContainer.value) return
    const containerWidth = mapContainer.value.clientWidth
    const containerHeight = mapContainer.value.clientHeight
    const scaledWidth = mapWidth * zoom.value
    const scaledHeight = mapHeight * zoom.value

    const minPanX = scaledWidth > containerWidth ? containerWidth - scaledWidth : 0
    const maxPanX = scaledWidth > containerWidth ? 0 : containerWidth - scaledWidth
    const minPanY = scaledHeight > containerHeight ? containerHeight - scaledHeight : 0
    const maxPanY = scaledHeight > containerHeight ? 0 : containerHeight - scaledHeight

    panX.value = Math.max(minPanX, Math.min(maxPanX, panX.value))
    panY.value = Math.max(minPanY, Math.min(maxPanY, panY.value))
}

// Handle zooming with the mouse wheel
const handleWheel = (event) => {
    event.preventDefault()

    // Get container dimensions
    const containerWidth = mapContainer.value.clientWidth
    const containerHeight = mapContainer.value.clientHeight

    // Calculate minimum zoom to ensure the map fills the container
    const minZoom = Math.max(containerWidth / mapWidth, containerHeight / mapHeight)
    const maxZoom = 200 // Maximum zoom limit (adjust as needed)
    const factor = event.deltaY < 0 ? 1.1 : 0.9 // Zoom in or out factor
    const newZoom = zoom.value * factor

    // Clamp the zoom between minZoom and maxZoom
    const clampedZoom = Math.max(minZoom, Math.min(maxZoom, newZoom))

    // Calculate the point under the mouse to keep it fixed during zoom
    const rect = mapContainer.value.getBoundingClientRect()
    const vx = event.clientX - rect.left
    const vy = event.clientY - rect.top
    const px = (vx - panX.value) / zoom.value
    const py = (vy - panY.value) / zoom.value

    // Update zoom and pan values
    zoom.value = clampedZoom
    panX.value = vx - px * clampedZoom
    panY.value = vy - py * clampedZoom

    // Ensure panning stays within bounds
    clampPan()
}

// Start panning
const handleMouseDown = (event) => {
    isDragging.value = true
    startX.value = event.clientX
    startY.value = event.clientY
    startPanX.value = panX.value
    startPanY.value = panY.value
    window.addEventListener('mousemove', handleMouseMove)
    window.addEventListener('mouseup', handleMouseUp)
}

// Update pan position during drag
const handleMouseMove = (event) => {
    if (!isDragging.value) return
    const dx = event.clientX - startX.value
    const dy = event.clientY - startY.value
    panX.value = startPanX.value + dx
    panY.value = startPanY.value + dy
    clampPan()
}

// Stop panning
const handleMouseUp = () => {
    isDragging.value = false
    window.removeEventListener('mousemove', handleMouseMove)
    window.removeEventListener('mouseup', handleMouseUp)
}

const setView = (lat, lon, zoomLevel = 1) => {
    const x = (lon + 180) * (mapWidth / 360)
    const y = (90 - lat) * (mapHeight / 180)
    const containerWidth = mapContainer.value.clientWidth
    const containerHeight = mapContainer.value.clientHeight
    panX.value = containerWidth / 2 - x * zoomLevel
    panY.value = containerHeight / 2 - y * zoomLevel
    zoom.value = zoomLevel
    clampPan()
}

const router = useRouter()

const gotoAgentsWithFilter = async (origins) => {
    await router.push({
        path: '/agents',
        query: { origins }
    })
}

// Set initial view to Europe on mount
onMounted(async () => {
    setView(48, 0, 2.5) // Center on Europe (approx. 50°N, 10°E) with zoom 1

    await API.value.generalUpdate()
    await nextTick()

    API.value.store.stats.countryCodes.forEach(cc => {
        const countries = document.querySelectorAll(`.cc-${cc}`)
        if (!countries || !countries.length) return
        countries.forEach(country => {
            country.style.fill = '#E53935'
            country.style.cursor = 'pointer'
            country.addEventListener('click', async () => {
                await gotoAgentsWithFilter(cc)
            })
        })
    })
})
</script>


<style scoped>
#map {
    width: 100%;
    height: 100vh;
    /* Adjust height as needed */
    overflow: hidden;
    user-select: none;
}

.map-container {
    position: relative;
    display: inline-block;
}

img {
    display: block;
    user-select: none;
}

.marker {
    position: absolute;
    width: 6px;
    height: 6px;
    background-color: red;
    border-radius: 50%;
    transform: translate(-50%, -50%);
}
</style>