<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <h1 class="text-3xl font-bold">Settings</h1>

    <!-- User Settings -->
    <div class="bg-white dark:bg-secondary rounded-lg p-6 border border-gray-200 dark:border-gray-700">
      <h2 class="text-xl font-bold mb-4">Profile</h2>
      <form @submit.prevent="updateProfile" class="space-y-4">
        <div>
          <label class="block text-sm font-medium mb-1">Name</label>
          <input v-model="profile.name" type="text" class="w-full px-3 py-2 border rounded dark:bg-secondary" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Email</label>
          <input v-model="profile.email" type="email" disabled class="w-full px-3 py-2 border rounded dark:bg-secondary opacity-50" />
        </div>
        <button type="submit" class="bg-primary text-white px-4 py-2 rounded hover:bg-primary-dark">Save Changes</button>
      </form>
    </div>

    <!-- Preferences -->
    <div class="bg-white dark:bg-secondary rounded-lg p-6 border border-gray-200 dark:border-gray-700">
      <h2 class="text-xl font-bold mb-4">Preferences</h2>
      <form @submit.prevent="updatePreferences" class="space-y-4">
        <div>
          <label class="block text-sm font-medium mb-1">Mileage Unit</label>
          <select v-model="preferences.units" class="w-full px-3 py-2 border rounded dark:bg-secondary">
            <option value="mi">Miles (mi)</option>
            <option value="km">Kilometers (km)</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Currency</label>
          <select v-model="preferences.currency" class="w-full px-3 py-2 border rounded dark:bg-secondary">
            <option value="USD">USD ($)</option>
            <option value="EUR">EUR (€)</option>
            <option value="GBP">GBP (£)</option>
            <option value="CAD">CAD (C$)</option>
            <option value="AUD">AUD (A$)</option>
          </select>
        </div>
        <button type="submit" class="bg-primary text-white px-4 py-2 rounded hover:bg-primary-dark">Save Preferences</button>
      </form>
    </div>

    <!-- Import/Export -->
    <div class="bg-white dark:bg-secondary rounded-lg p-6 border border-gray-200 dark:border-gray-700">
      <h2 class="text-xl font-bold mb-4">Import Data</h2>
      <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">Import vehicle data from other applications</p>
      <div class="space-y-3">
        <button @click="promptHammondImport" class="w-full bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600">
          Import from Hammond
        </button>
        <button @click="promptFuellyImport" class="w-full bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">
          Import from Fuelly
        </button>
        <button @click="promptClarksonImport" class="w-full bg-purple-500 text-white px-4 py-2 rounded hover:bg-purple-600">
          Import from Clarkson Backup
        </button>
      </div>
    </div>

    <!-- Danger Zone -->
    <div class="bg-red-50 dark:bg-red-900 rounded-lg p-6 border border-red-200 dark:border-red-700">
      <h2 class="text-xl font-bold mb-4 text-red-700 dark:text-red-300">Danger Zone</h2>
      <button @click="promptDeleteAccount" class="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700">
        Delete Account
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const profile = ref({ name: '', email: '' })
const preferences = ref({ units: 'mi', currency: 'USD' })

onMounted(() => {
  if (authStore.user) {
    profile.value = { ...authStore.user }
    preferences.value = {
      units: localStorage.getItem('units') || 'mi',
      currency: localStorage.getItem('currency') || 'USD',
    }
  }
})

async function updateProfile() {
  const response = await fetch(`http://localhost:3000/api/users/${authStore.user.id}`, {
    method: 'PUT',
    headers: {
      'Authorization': authStore.token,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ name: profile.value.name }),
  })
  if (response.ok) {
    alert('Profile updated!')
  }
}

async function updatePreferences() {
  localStorage.setItem('units', preferences.value.units)
  localStorage.setItem('currency', preferences.value.currency)
  alert('Preferences saved!')
}

function promptHammondImport() {
  alert('Hammond import feature coming soon!')
}

function promptFuellyImport() {
  alert('Fuelly import feature coming soon!')
}

function promptClarksonImport() {
  alert('Clarkson import feature coming soon!')
}

function promptDeleteAccount() {
  if (confirm('Are you sure? This cannot be undone.')) {
    alert('Account deletion feature coming soon!')
  }
}
</script>
