<template>
  <div class="space-y-6">
    <h2 class="text-2xl font-bold">Import Data</h2>

    <!-- Import Options -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <!-- Hammond Import -->
      <div class="bg-white dark:bg-secondary rounded-lg p-6 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition">
        <h3 class="text-lg font-bold mb-2">Hammond</h3>
        <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">Import from Hammond vehicle expense tracker</p>
        <label class="block">
          <input type="file" @change="importHammond" accept=".json" class="hidden" />
          <span class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 cursor-pointer inline-block">
            Choose File
          </span>
        </label>
      </div>

      <!-- Fuelly Import -->
      <div class="bg-white dark:bg-secondary rounded-lg p-6 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition">
        <h3 class="text-lg font-bold mb-2">Fuelly</h3>
        <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">Import from Fuelly CSV export</p>
        <label class="block">
          <input type="file" @change="importFuelly" accept=".csv" class="hidden" />
          <span class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600 cursor-pointer inline-block">
            Choose File
          </span>
        </label>
      </div>

      <!-- Clarkson Backup Import -->
      <div class="bg-white dark:bg-secondary rounded-lg p-6 border border-gray-200 dark:border-gray-700 hover:shadow-lg transition">
        <h3 class="text-lg font-bold mb-2">Clarkson Backup</h3>
        <p class="text-sm text-gray-600 dark:text-gray-400 mb-4">Import from Clarkson backup JSON</p>
        <label class="block">
          <input type="file" @change="importClarkson" accept=".json" class="hidden" />
          <span class="bg-purple-500 text-white px-4 py-2 rounded hover:bg-purple-600 cursor-pointer inline-block">
            Choose File
          </span>
        </label>
      </div>
    </div>

    <!-- Import Progress -->
    <div v-if="importing" class="bg-blue-50 dark:bg-blue-900 p-4 rounded-lg">
      <p class="text-sm font-semibold">Importing...</p>
      <div class="mt-2 bg-blue-200 dark:bg-blue-800 rounded-full h-2 overflow-hidden">
        <div class="bg-blue-600 h-full w-1/2 animate-pulse"></div>
      </div>
    </div>

    <!-- Import Results -->
    <div v-if="importResult" :class="['p-4 rounded-lg', importResult.error ? 'bg-red-50 dark:bg-red-900' : 'bg-green-50 dark:bg-green-900']">
      <h3 class="font-bold mb-2">Import Complete</h3>
      <ul class="text-sm space-y-1">
        <li v-if="importResult.vehicles">✓ {{ importResult.vehicles }} vehicles imported</li>
        <li v-if="importResult.fuel">✓ {{ importResult.fuel }} fuel entries imported</li>
        <li v-if="importResult.expenses">✓ {{ importResult.expenses }} expenses imported</li>
        <li v-if="importResult.reminders">✓ {{ importResult.reminders }} reminders imported</li>
      </ul>
      <div v-if="importResult.errors?.length" class="mt-3 text-red-700 dark:text-red-300">
        <p class="font-semibold">Errors:</p>
        <ul class="text-xs">
          <li v-for="(err, i) in importResult.errors" :key="i">- {{ err }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const importing = ref(false)
const importResult = ref(null)

const importHammond = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  await performImport(file, 'hammond')
}

const importFuelly = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  await performImport(file, 'fuelly')
}

const importClarkson = async (event) => {
  const file = event.target.files[0]
  if (!file) return

  await performImport(file, 'clarkson')
}

const performImport = async (file, type) => {
  const formData = new FormData()
  formData.append('file', file)

  importing.value = true
  importResult.value = null

  try {
    const response = await fetch(`http://localhost:3000/api/import/${type}`, {
      method: 'POST',
      headers: { 'Authorization': authStore.token },
      body: formData,
    })

    const data = await response.json()
    importResult.value = data
  } catch (error) {
    importResult.value = { error: true, message: 'Import failed' }
  } finally {
    importing.value = false
  }
}
</script>
