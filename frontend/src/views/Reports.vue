<template>
  <div class="space-y-8">
    <h1 class="text-3xl font-bold">Reports & Analytics</h1>

    <!-- Report Type Selection -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <button @click="setReportType('vehicle')" 
              :class="['p-4 rounded-lg border-2 transition', reportType === 'vehicle' ? 'border-primary bg-blue-50 dark:bg-blue-900' : 'border-gray-200 dark:border-gray-700']">
        <p class="font-bold">Vehicle Report</p>
        <p class="text-sm text-gray-600 dark:text-gray-400">Detailed vehicle statistics</p>
      </button>
      <button @click="setReportType('comparison')" 
              :class="['p-4 rounded-lg border-2 transition', reportType === 'comparison' ? 'border-primary bg-blue-50 dark:bg-blue-900' : 'border-gray-200 dark:border-gray-700']">
        <p class="font-bold">Vehicle Comparison</p>
        <p class="text-sm text-gray-600 dark:text-gray-400">Compare all vehicles</p>
      </button>
      <button @click="setReportType('search')" 
              :class="['p-4 rounded-lg border-2 transition', reportType === 'search' ? 'border-primary bg-blue-50 dark:bg-blue-900' : 'border-gray-200 dark:border-gray-700']">
        <p class="font-bold">Search</p>
        <p class="text-sm text-gray-600 dark:text-gray-400">Find entries</p>
      </button>
    </div>

    <!-- Vehicle Report -->
    <div v-if="reportType === 'vehicle'" class="space-y-6">
      <div>
        <label class="block text-sm font-medium mb-2">Select Vehicle</label>
        <select v-model="selectedVehicleId" @change="loadVehicleReport" class="px-4 py-2 border rounded dark:bg-secondary">
          <option value="">Choose a vehicle</option>
          <option v-for="v in vehicles" :key="v.id" :value="v.id">
            {{ v.year }} {{ v.make }} {{ v.model }}
          </option>
        </select>
      </div>

      <div v-if="vehicleReport" class="space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="bg-blue-50 dark:bg-blue-900 p-4 rounded-lg">
            <p class="text-sm text-gray-600 dark:text-gray-400">Total Distance</p>
            <p class="text-2xl font-bold">{{ vehicleReport.total_distance?.toFixed(1) }} {{ vehicleReport.vehicle.mileage_unit }}</p>
          </div>
          <div class="bg-green-50 dark:bg-green-900 p-4 rounded-lg">
            <p class="text-sm text-gray-600 dark:text-gray-400">Average MPG</p>
            <p class="text-2xl font-bold">{{ vehicleReport.average_mpg?.toFixed(1) }}</p>
          </div>
          <div class="bg-purple-50 dark:bg-purple-900 p-4 rounded-lg">
            <p class="text-sm text-gray-600 dark:text-gray-400">Total Cost</p>
            <p class="text-2xl font-bold">${{ vehicleReport.total_cost?.toFixed(2) }}</p>
          </div>
          <div class="bg-orange-50 dark:bg-orange-900 p-4 rounded-lg">
            <p class="text-sm text-gray-600 dark:text-gray-400">Cost/Mile</p>
            <p class="text-2xl font-bold">${{ (vehicleReport.total_cost / vehicleReport.total_distance)?.toFixed(3) }}</p>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="bg-white dark:bg-secondary rounded-lg p-4 border border-gray-200 dark:border-gray-700">
            <p class="font-semibold mb-2">Fuel Costs</p>
            <p class="text-2xl font-bold text-blue-600 dark:text-blue-400">${{ vehicleReport.fuel_costs?.toFixed(2) }}</p>
          </div>
          <div class="bg-white dark:bg-secondary rounded-lg p-4 border border-gray-200 dark:border-gray-700">
            <p class="font-semibold mb-2">Maintenance</p>
            <p class="text-2xl font-bold text-green-600 dark:text-green-400">${{ vehicleReport.maintenance_cost?.toFixed(2) }}</p>
          </div>
          <div class="bg-white dark:bg-secondary rounded-lg p-4 border border-gray-200 dark:border-gray-700">
            <p class="font-semibold mb-2">Other</p>
            <p class="text-2xl font-bold text-orange-600 dark:text-orange-400">${{ vehicleReport.other_costs?.toFixed(2) }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Comparison Report -->
    <div v-if="reportType === 'comparison'" class="space-y-4">
      <div v-for="comp in comparisonReport" :key="comp.vehicle.id" class="bg-white dark:bg-secondary rounded-lg p-6 border border-gray-200 dark:border-gray-700">
        <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
          <div>
            <p class="font-bold">{{ comp.vehicle.year }} {{ comp.vehicle.make }} {{ comp.vehicle.model }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-600 dark:text-gray-400">Total Cost</p>
            <p class="font-bold">${{ comp.total_cost?.toFixed(2) }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-600 dark:text-gray-400">Distance</p>
            <p class="font-bold">{{ comp.total_miles?.toFixed(1) }} mi</p>
          </div>
          <div>
            <p class="text-xs text-gray-600 dark:text-gray-400">Avg MPG</p>
            <p class="font-bold">{{ comp.average_mpg?.toFixed(1) }}</p>
          </div>
          <div>
            <p class="text-xs text-gray-600 dark:text-gray-400">Cost/Mile</p>
            <p class="font-bold">${{ comp.cost_per_mile?.toFixed(3) }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Search -->
    <div v-if="reportType === 'search'" class="space-y-4">
      <input v-model="searchQuery" @keyup.enter="performSearch" type="text" placeholder="Search entries..." 
             class="w-full px-4 py-2 border rounded dark:bg-secondary" />
      <button @click="performSearch" class="bg-primary text-white px-4 py-2 rounded hover:bg-primary-dark">
        Search
      </button>

      <div v-if="searchResults" class="space-y-4">
        <p class="text-sm text-gray-600">{{ searchResults.count }} results found</p>
        <div v-if="searchResults.fuel?.length" class="space-y-2">
          <h3 class="font-bold">Fuel Entries</h3>
          <div v-for="f in searchResults.fuel" :key="f.id" class="bg-blue-50 dark:bg-blue-900 p-3 rounded text-sm">
            {{ formatDate(f.date) }} - {{ f.location }} ({{ f.gallons }} gal @ ${{ f.price }})
          </div>
        </div>
        <div v-if="searchResults.expenses?.length" class="space-y-2">
          <h3 class="font-bold">Expenses</h3>
          <div v-for="e in searchResults.expenses" :key="e.id" class="bg-green-50 dark:bg-green-900 p-3 rounded text-sm">
            {{ formatDate(e.date) }} - {{ e.category }}: ${{ e.amount }}
          </div>
        </div>
      </div>
    </div>

    <!-- Export Options -->
    <div class="bg-white dark:bg-secondary rounded-lg p-6 border border-gray-200 dark:border-gray-700">
      <h2 class="text-xl font-bold mb-4">Export Data</h2>
      <div class="flex flex-wrap gap-2">
        <button @click="exportCSV" class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">
          Download CSV
        </button>
        <button @click="exportJSON" class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
          Download JSON Backup
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const reportType = ref('vehicle')
const vehicles = ref([])
const selectedVehicleId = ref('')
const vehicleReport = ref(null)
const comparisonReport = ref([])
const searchQuery = ref('')
const searchResults = ref(null)

const loadVehicles = async () => {
  try {
    const response = await fetch('http://localhost:3000/api/vehicles', {
      headers: { 'Authorization': authStore.token },
    })
    vehicles.value = await response.json()
  } catch (error) {
    console.error('Error loading vehicles:', error)
  }
}

const setReportType = (type) => {
  reportType.value = type
}

const loadVehicleReport = async () => {
  if (!selectedVehicleId.value) return

  try {
    const response = await fetch(`http://localhost:3000/api/vehicles/${selectedVehicleId.value}/report`, {
      headers: { 'Authorization': authStore.token },
    })
    vehicleReport.value = await response.json()
  } catch (error) {
    console.error('Error loading report:', error)
  }
}

const loadComparisonReport = async () => {
  try {
    const response = await fetch('http://localhost:3000/api/report/comparison', {
      headers: { 'Authorization': authStore.token },
    })
    const data = await response.json()
    comparisonReport.value = data.vehicles
  } catch (error) {
    console.error('Error loading comparison:', error)
  }
}

const performSearch = async () => {
  if (!searchQuery.value) return

  try {
    const response = await fetch(`http://localhost:3000/api/search?q=${encodeURIComponent(searchQuery.value)}`, {
      headers: { 'Authorization': authStore.token },
    })
    searchResults.value = await response.json()
  } catch (error) {
    console.error('Error searching:', error)
  }
}

const exportCSV = async () => {
  const response = await fetch('http://localhost:3000/api/export/csv', {
    headers: { 'Authorization': authStore.token },
  })
  const blob = await response.blob()
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'clarkson-export.csv'
  a.click()
}

const exportJSON = async () => {
  const response = await fetch('http://localhost:3000/api/export/json', {
    headers: { 'Authorization': authStore.token },
  })
  const blob = await response.blob()
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'clarkson-backup.json'
  a.click()
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString()
}

onMounted(async () => {
  await loadVehicles()
})
</script>
