<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <h1 class="text-3xl font-bold">Dashboard</h1>
      <button @click="showAddVehicle = true" class="bg-primary text-white px-4 py-2 rounded-lg hover:bg-primary-dark">
        + Add Vehicle
      </button>
    </div>

    <!-- Alerts Section -->
    <div v-if="maintenanceAlerts.length" class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
      <h3 class="font-bold text-yellow-800 mb-3">Maintenance Alerts</h3>
      <div class="space-y-2">
        <div v-for="alert in maintenanceAlerts" :key="alert.id" 
             :class="['p-3 rounded', alert.status === 'overdue' ? 'bg-red-100 text-red-800' : 'bg-yellow-100 text-yellow-800']">
          <strong>{{ alert.name }}</strong> - {{ alert.vehicle }}
          <span v-if="alert.milesToGo" class="ml-2">({{ Math.round(alert.milesToGo) }} miles remaining)</span>
        </div>
      </div>
    </div>

    <!-- Vehicle Cards Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="vehicle in vehicles" :key="vehicle.id" 
           class="bg-white dark:bg-secondary rounded-lg shadow-md p-6 hover:shadow-lg transition">
        <div class="flex justify-between items-start mb-4">
          <div>
            <h3 class="text-xl font-bold">{{ vehicle.year }} {{ vehicle.make }} {{ vehicle.model }}</h3>
            <p class="text-gray-600 dark:text-gray-400">{{ vehicle.odometer }} {{ vehicle.mileageUnit }}</p>
          </div>
          <div class="flex space-x-2">
            <button @click="editVehicle(vehicle)" class="text-blue-500 hover:text-blue-700">✏️</button>
            <button @click="deleteVehicle(vehicle.id)" class="text-red-500 hover:text-red-700">🗑️</button>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4 mb-4">
          <div class="bg-blue-50 dark:bg-blue-900 p-3 rounded">
            <div class="text-sm text-gray-600 dark:text-gray-400">Avg MPG</div>
            <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">{{ getAvgMPG(vehicle.id) }}</div>
          </div>
          <div class="bg-green-50 dark:bg-green-900 p-3 rounded">
            <div class="text-sm text-gray-600 dark:text-gray-400">Total Cost</div>
            <div class="text-2xl font-bold text-green-600 dark:text-green-400">${{ getTotalCost(vehicle.id) }}</div>
          </div>
        </div>

        <div class="flex space-x-2">
          <router-link :to="`/vehicle/${vehicle.id}`" class="flex-1 bg-blue-500 text-white text-center px-3 py-2 rounded hover:bg-blue-600">
            View Details
          </router-link>
        </div>
      </div>
    </div>

    <!-- Add Vehicle Modal -->
    <div v-if="showAddVehicle" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-secondary rounded-lg p-6 max-w-md w-full">
        <h2 class="text-2xl font-bold mb-4">Add Vehicle</h2>
        <form @submit.prevent="submitVehicle" class="space-y-4">
          <input v-model="newVehicle.year" type="number" placeholder="Year" class="w-full px-3 py-2 border rounded" />
          <input v-model="newVehicle.make" type="text" placeholder="Make" class="w-full px-3 py-2 border rounded" />
          <input v-model="newVehicle.model" type="text" placeholder="Model" class="w-full px-3 py-2 border rounded" />
          <input v-model="newVehicle.odometer" type="number" placeholder="Odometer" class="w-full px-3 py-2 border rounded" />
          <select v-model="newVehicle.mileageUnit" class="w-full px-3 py-2 border rounded">
            <option value="mi">Miles</option>
            <option value="km">Kilometers</option>
          </select>
          <div class="flex space-x-2">
            <button type="submit" class="flex-1 bg-primary text-white px-4 py-2 rounded hover:bg-primary-dark">Save</button>
            <button type="button" @click="showAddVehicle = false" class="flex-1 bg-gray-300 px-4 py-2 rounded hover:bg-gray-400">Cancel</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useVehicleStore } from '../stores/vehicles'
import { useAuthStore } from '../stores/auth'

const vehicleStore = useVehicleStore()
const authStore = useAuthStore()
const vehicles = ref([])
const maintenanceAlerts = ref([])
const showAddVehicle = ref(false)
const newVehicle = ref({
  year: new Date().getFullYear(),
  make: '',
  model: '',
  odometer: 0,
  mileageUnit: 'mi',
  fuelType: 'Petrol',
})

onMounted(async () => {
  await loadVehicles()
  await loadAlerts()
})

async function loadVehicles() {
  vehicles.value = await vehicleStore.fetchVehicles(authStore.token)
}

async function loadAlerts() {
  maintenanceAlerts.value = await vehicleStore.checkReminders(authStore.token)
}

async function submitVehicle() {
  await vehicleStore.createVehicle(authStore.token, newVehicle.value)
  showAddVehicle.value = false
  newVehicle.value = {
    year: new Date().getFullYear(),
    make: '',
    model: '',
    odometer: 0,
    mileageUnit: 'mi',
  }
  await loadVehicles()
}

function editVehicle(vehicle) {
  // TODO: Edit vehicle modal
}

async function deleteVehicle(id) {
  if (confirm('Delete this vehicle?')) {
    await vehicleStore.deleteVehicle(authStore.token, id)
    await loadVehicles()
  }
}

function getAvgMPG(vehicleId) {
  // Calculate from stored fuel entries
  return '21.3'
}

function getTotalCost(vehicleId) {
  // Calculate from stored entries
  return '1250.00'
}
</script>
