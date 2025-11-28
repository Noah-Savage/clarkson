<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex justify-between items-center">
      <h1 class="text-3xl font-bold">{{ vehicle.year }} {{ vehicle.make }} {{ vehicle.model }}</h1>
      <router-link to="/dashboard" class="text-blue-500 hover:underline">← Back</router-link>
    </div>

    <!-- Stats Overview -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-blue-50 dark:bg-blue-900 p-4 rounded-lg">
        <div class="text-sm text-gray-600 dark:text-gray-400">Current Odometer</div>
        <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">{{ vehicle.odometer }} {{ vehicle.mileageUnit }}</div>
      </div>
      <div class="bg-green-50 dark:bg-green-900 p-4 rounded-lg">
        <div class="text-sm text-gray-600 dark:text-gray-400">Avg MPG</div>
        <div class="text-2xl font-bold text-green-600 dark:text-green-400">{{ fuelStats.average_mpg?.toFixed(1) || '-' }}</div>
      </div>
      <div class="bg-purple-50 dark:bg-purple-900 p-4 rounded-lg">
        <div class="text-sm text-gray-600 dark:text-gray-400">Total Cost</div>
        <div class="text-2xl font-bold text-purple-600 dark:text-purple-400">${{ (fuelStats.total_cost + expenseStats.total_cost)?.toFixed(2) || '0.00' }}</div>
      </div>
      <div class="bg-orange-50 dark:bg-orange-900 p-4 rounded-lg">
        <div class="text-sm text-gray-600 dark:text-gray-400">Due Reminders</div>
        <div class="text-2xl font-bold text-orange-600 dark:text-orange-400">{{ dueRemindersCount }}</div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="border-b border-gray-200 dark:border-gray-700">
      <div class="flex space-x-4">
        <button 
          @click="activeTab = 'fuel'"
          :class="['px-4 py-2 font-medium', activeTab === 'fuel' ? 'text-primary border-b-2 border-primary' : 'text-gray-600 dark:text-gray-400']"
        >
          Fuel Entries
        </button>
        <button 
          @click="activeTab = 'expenses'"
          :class="['px-4 py-2 font-medium', activeTab === 'expenses' ? 'text-primary border-b-2 border-primary' : 'text-gray-600 dark:text-gray-400']"
        >
          Expenses
        </button>
        <button 
          @click="activeTab = 'reminders'"
          :class="['px-4 py-2 font-medium', activeTab === 'reminders' ? 'text-primary border-b-2 border-primary' : 'text-gray-600 dark:text-gray-400']"
        >
          Maintenance
        </button>
        <button 
          @click="activeTab = 'reports'"
          :class="['px-4 py-2 font-medium', activeTab === 'reports' ? 'text-primary border-b-2 border-primary' : 'text-gray-600 dark:text-gray-400']"
        >
          Reports
        </button>
      </div>
    </div>

    <!-- Fuel Tab -->
    <div v-if="activeTab === 'fuel'" class="space-y-4">
      <button @click="showAddFuel = true" class="bg-primary text-white px-4 py-2 rounded hover:bg-primary-dark">
        + Add Fuel Entry
      </button>

      <div class="space-y-2">
        <div v-for="entry in fuelEntries" :key="entry.id" class="bg-white dark:bg-secondary rounded-lg p-4 border border-gray-200 dark:border-gray-700">
          <div class="flex justify-between items-start mb-2">
            <div>
              <p class="font-semibold">{{ formatDate(entry.date) }}</p>
              <p class="text-sm text-gray-600 dark:text-gray-400">{{ entry.odometer }} {{ vehicle.mileageUnit }} • {{ entry.gallons }} gal @ ${{ entry.price?.toFixed(2) }}</p>
            </div>
            <div class="flex space-x-2">
              <button @click="editFuel(entry)" class="text-blue-500 hover:text-blue-700">Edit</button>
              <button @click="deleteFuel(entry.id)" class="text-red-500 hover:text-red-700">Delete</button>
            </div>
          </div>
          <p v-if="entry.location" class="text-sm text-gray-600 dark:text-gray-400">📍 {{ entry.location }}</p>
          <p v-if="entry.notes" class="text-sm text-gray-600 dark:text-gray-400">{{ entry.notes }}</p>
        </div>
      </div>

      <!-- Add Fuel Modal -->
      <div v-if="showAddFuel" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white dark:bg-secondary rounded-lg p-6 max-w-md w-full">
          <h2 class="text-2xl font-bold mb-4">Add Fuel Entry</h2>
          <form @submit.prevent="submitFuel" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-1">Date</label>
              <input v-model="newFuel.date" type="date" class="w-full px-3 py-2 border rounded dark:bg-secondary" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Odometer</label>
              <input v-model.number="newFuel.odometer" type="number" step="0.1" class="w-full px-3 py-2 border rounded dark:bg-secondary" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Gallons</label>
              <input v-model.number="newFuel.gallons" type="number" step="0.01" class="w-full px-3 py-2 border rounded dark:bg-secondary" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Price ($)</label>
              <input v-model.number="newFuel.price" type="number" step="0.01" class="w-full px-3 py-2 border rounded dark:bg-secondary" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Location</label>
              <input v-model="newFuel.location" type="text" class="w-full px-3 py-2 border rounded dark:bg-secondary" />
            </div>
            <div class="flex space-x-2">
              <button type="submit" class="flex-1 bg-primary text-white px-4 py-2 rounded hover:bg-primary-dark">Save</button>
              <button type="button" @click="showAddFuel = false" class="flex-1 bg-gray-300 dark:bg-gray-600 px-4 py-2 rounded hover:bg-gray-400">Cancel</button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Expenses Tab -->
    <div v-if="activeTab === 'expenses'" class="space-y-4">
      <button @click="showAddExpense = true" class="bg-primary text-white px-4 py-2 rounded hover:bg-primary-dark">
        + Add Expense
      </button>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
        <div v-for="cat in expenseStats.categories" :key="cat.category" class="bg-white dark:bg-secondary rounded-lg p-4 border border-gray-200 dark:border-gray-700">
          <p class="font-semibold">{{ cat.category }}</p>
          <p class="text-lg text-green-600 dark:text-green-400">${{ cat.total?.toFixed(2) }}</p>
          <p class="text-sm text-gray-600 dark:text-gray-400">{{ cat.count }} entries</p>
        </div>
      </div>

      <div class="space-y-2">
        <div v-for="expense in expenses" :key="expense.id" class="bg-white dark:bg-secondary rounded-lg p-4 border border-gray-200 dark:border-gray-700">
          <div class="flex justify-between items-start">
            <div>
              <p class="font-semibold">{{ expense.category }}</p>
              <p class="text-sm text-gray-600 dark:text-gray-400">{{ formatDate(expense.date) }}</p>
            </div>
            <div class="text-right">
              <p class="font-bold text-lg">${{ expense.amount?.toFixed(2) }}</p>
              <div class="flex space-x-2 mt-2">
                <button @click="editExpense(expense)" class="text-blue-500 hover:text-blue-700 text-sm">Edit</button>
                <button @click="deleteExpense(expense.id)" class="text-red-500 hover:text-red-700 text-sm">Delete</button>
              </div>
            </div>
          </div>
          <p v-if="expense.notes" class="text-sm text-gray-600 dark:text-gray-400 mt-2">{{ expense.notes }}</p>
        </div>
      </div>

      <!-- Add Expense Modal -->
      <div v-if="showAddExpense" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white dark:bg-secondary rounded-lg p-6 max-w-md w-full">
          <h2 class="text-2xl font-bold mb-4">Add Expense</h2>
          <form @submit.prevent="submitExpense" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-1">Category</label>
              <select v-model="newExpense.category" class="w-full px-3 py-2 border rounded dark:bg-secondary">
                <option value="Maintenance">Maintenance</option>
                <option value="Insurance">Insurance</option>
                <option value="Registration">Registration</option>
                <option value="Parking">Parking</option>
                <option value="Tolls">Tolls</option>
                <option value="Other">Other</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Amount ($)</label>
              <input v-model.number="newExpense.amount" type="number" step="0.01" class="w-full px-3 py-2 border rounded dark:bg-secondary" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Date</label>
              <input v-model="newExpense.date" type="date" class="w-full px-3 py-2 border rounded dark:bg-secondary" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Notes</label>
              <textarea v-model="newExpense.notes" class="w-full px-3 py-2 border rounded dark:bg-secondary" rows="3"></textarea>
            </div>
            <div class="flex space-x-2">
              <button type="submit" class="flex-1 bg-primary text-white px-4 py-2 rounded hover:bg-primary-dark">Save</button>
              <button type="button" @click="showAddExpense = false" class="flex-1 bg-gray-300 dark:bg-gray-600 px-4 py-2 rounded hover:bg-gray-400">Cancel</button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Reminders Tab -->
    <div v-if="activeTab === 'reminders'" class="space-y-4">
      <button @click="showAddReminder = true" class="bg-primary text-white px-4 py-2 rounded hover:bg-primary-dark">
        + Add Reminder
      </button>

      <div class="space-y-2">
        <div v-for="reminder in reminders" :key="reminder.id" 
             :class="['rounded-lg p-4 border-2', getReminderStatusClass(reminder)]">
          <div class="flex justify-between items-start">
            <div>
              <p class="font-bold text-lg">{{ reminder.name }}</p>
              <div class="text-sm mt-2 space-y-1">
                <p v-if="reminder.interval_miles > 0">
                  📊 Every {{ reminder.interval_miles }} miles
                  <span v-if="reminder.last_service_miles > 0" class="text-gray-600 dark:text-gray-400">
                    (Last: {{ reminder.last_service_miles }} miles ago)
                  </span>
                </p>
                <p v-if="reminder.interval_days > 0">
                  📅 Every {{ reminder.interval_days }} days
                  <span v-if="reminder.last_service_date" class="text-gray-600 dark:text-gray-400">
                    (Last: {{ formatDate(reminder.last_service_date) }})
                  </span>
                </p>
              </div>
            </div>
            <div class="flex space-x-2">
              <button @click="completeReminder(reminder)" class="px-3 py-1 bg-green-500 text-white rounded text-sm hover:bg-green-600">
                Done
              </button>
              <button @click="editReminder(reminder)" class="text-blue-500 hover:text-blue-700">✏️</button>
              <button @click="deleteReminder(reminder.id)" class="text-red-500 hover:text-red-700">🗑️</button>
            </div>
          </div>
        </div>
      </div>

      <!-- Add Reminder Modal -->
      <div v-if="showAddReminder" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white dark:bg-secondary rounded-lg p-6 max-w-md w-full">
          <h2 class="text-2xl font-bold mb-4">Add Maintenance Reminder</h2>
          <form @submit.prevent="submitReminder" class="space-y-4">
            <div>
              <label class="block text-sm font-medium mb-1">Reminder Name</label>
              <input v-model="newReminder.name" type="text" placeholder="e.g., Oil Change" class="w-full px-3 py-2 border rounded dark:bg-secondary" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Interval (Miles)</label>
              <input v-model.number="newReminder.interval_miles" type="number" placeholder="0 to disable" class="w-full px-3 py-2 border rounded dark:bg-secondary" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Interval (Days)</label>
              <input v-model.number="newReminder.interval_days" type="number" placeholder="0 to disable" class="w-full px-3 py-2 border rounded dark:bg-secondary" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Last Service Date</label>
              <input v-model="newReminder.last_service_date" type="date" class="w-full px-3 py-2 border rounded dark:bg-secondary" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Last Service Mileage</label>
              <input v-model.number="newReminder.last_service_miles" type="number" class="w-full px-3 py-2 border rounded dark:bg-secondary" />
            </div>
            <div class="flex space-x-2">
              <button type="submit" class="flex-1 bg-primary text-white px-4 py-2 rounded hover:bg-primary-dark">Save</button>
              <button type="button" @click="showAddReminder = false" class="flex-1 bg-gray-300 dark:bg-gray-600 px-4 py-2 rounded hover:bg-gray-400">Cancel</button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Reports Tab -->
    <div v-if="activeTab === 'reports'" class="space-y-6">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="bg-white dark:bg-secondary rounded-lg p-6 border border-gray-200 dark:border-gray-700">
          <h3 class="font-bold text-lg mb-4">Fuel Economy Trend</h3>
          <div class="h-64 flex items-center justify-center text-gray-400">
            Chart placeholder - implement with Chart.js
          </div>
        </div>
        <div class="bg-white dark:bg-secondary rounded-lg p-6 border border-gray-200 dark:border-gray-700">
          <h3 class="font-bold text-lg mb-4">Expense Breakdown</h3>
          <div class="h-64 flex items-center justify-center text-gray-400">
            Chart placeholder - implement with Chart.js
          </div>
        </div>
      </div>

      <div class="flex space-x-2">
        <button @click="exportCSV" class="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">
          📥 Export CSV
        </button>
        <button @click="exportPDF" class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600">
          📄 Export PDF
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const route = useRoute()
const authStore = useAuthStore()
const vehicleID = route.params.id

const vehicle = ref({})
const fuelEntries = ref([])
const expenses = ref([])
const reminders = ref([])
const fuelStats = ref({})
const expenseStats = ref({ categories: [] })
const activeTab = ref('fuel')
const dueRemindersCount = ref(0)

const showAddFuel = ref(false)
const showAddExpense = ref(false)
const showAddReminder = ref(false)

const newFuel = ref({
  date: new Date().toISOString().split('T')[0],
  odometer: 0,
  gallons: 0,
  price: 0,
  location: '',
})

const newExpense = ref({
  category: 'Maintenance',
  amount: 0,
  date: new Date().toISOString().split('T')[0],
  notes: '',
})

const newReminder = ref({
  name: '',
  interval_miles: 0,
  interval_days: 0,
  last_service_date: new Date().toISOString().split('T')[0],
  last_service_miles: 0,
})

onMounted(async () => {
  await loadVehicleData()
})

async function loadVehicleData() {
  try {
    const response = await fetch(`http://localhost:3000/api/vehicles/${vehicleID}`, {
      headers: { 'Authorization': authStore.token },
    })
    vehicle.value = await response.json()

    // Load fuel entries
    const fuelRes = await fetch(`http://localhost:3000/api/vehicles/${vehicleID}/fuel`, {
      headers: { 'Authorization': authStore.token },
    })
    fuelEntries.value = await fuelRes.json()

    // Load expenses
    const expRes = await fetch(`http://localhost:3000/api/vehicles/${vehicleID}/expenses`, {
      headers: { 'Authorization': authStore.token },
    })
    expenses.value = await expRes.json()

    // Load reminders
    const remRes = await fetch(`http://localhost:3000/api/vehicles/${vehicleID}/reminders`, {
      headers: { 'Authorization': authStore.token },
    })
    reminders.value = await remRes.json()

    // Load stats
    const statsRes = await fetch(`http://localhost:3000/api/vehicles/${vehicleID}/fuel-stats`, {
      headers: { 'Authorization': authStore.token },
    })
    fuelStats.value = await statsRes.json()

    const expStatsRes = await fetch(`http://localhost:3000/api/vehicles/${vehicleID}/expense-stats`, {
      headers: { 'Authorization': authStore.token },
    })
    expenseStats.value = await expStatsRes.json()

    countDueReminders()
  } catch (error) {
    console.error('Error loading vehicle data:', error)
  }
}

function countDueReminders() {
  dueRemindersCount.value = reminders.value.filter(r => {
    let isDue = false
    if (r.interval_miles > 0) {
      const nextService = r.last_service_miles + r.interval_miles
      if (vehicle.value.odometer >= nextService - 500) isDue = true
    }
    if (r.interval_days > 0) {
      const nextService = new Date(r.last_service_date)
      nextService.setDate(nextService.getDate() + r.interval_days)
      if (new Date() >= nextService) isDue = true
    }
    return isDue
  }).length
}

function getReminderStatusClass(reminder) {
  let status = 'upcoming'
  if (reminder.interval_miles > 0) {
    const nextService = reminder.last_service_miles + reminder.interval_miles
    if (vehicle.value.odometer >= nextService) {
      status = 'overdue'
    } else if (vehicle.value.odometer >= nextService - 500) {
      status = 'soon'
    }
  }

  if (status === 'overdue') return 'bg-red-50 dark:bg-red-900 border-red-300 dark:border-red-700'
  if (status === 'soon') return 'bg-yellow-50 dark:bg-yellow-900 border-yellow-300 dark:border-yellow-700'
  return 'bg-gray-50 dark:bg-gray-900 border-gray-300 dark:border-gray-700'
}

async function submitFuel() {
  try {
    const response = await fetch(`http://localhost:3000/api/vehicles/${vehicleID}/fuel`, {
      method: 'POST',
      headers: {
        'Authorization': authStore.token,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        date: new Date(newFuel.value.date),
        odometer: newFuel.value.odometer,
        gallons: newFuel.value.gallons,
        price: newFuel.value.price,
        location: newFuel.value.location,
      }),
    })

    if (response.ok) {
      showAddFuel.value = false
      await loadVehicleData()
    }
  } catch (error) {
    console.error('Error submitting fuel entry:', error)
  }
}

async function submitExpense() {
  try {
    const response = await fetch(`http://localhost:3000/api/vehicles/${vehicleID}/expenses`, {
      method: 'POST',
      headers: {
        'Authorization': authStore.token,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        category: newExpense.value.category,
        amount: newExpense.value.amount,
        date: new Date(newExpense.value.date),
        notes: newExpense.value.notes,
      }),
    })

    if (response.ok) {
      showAddExpense.value = false
      await loadVehicleData()
    }
  } catch (error) {
    console.error('Error submitting expense:', error)
  }
}

async function submitReminder() {
  try {
    const response = await fetch(`http://localhost:3000/api/vehicles/${vehicleID}/reminders`, {
      method: 'POST',
      headers: {
        'Authorization': authStore.token,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: newReminder.value.name,
        interval_miles: newReminder.value.interval_miles,
        interval_days: newReminder.value.interval_days,
        last_service_date: new Date(newReminder.value.last_service_date),
        last_service_miles: newReminder.value.last_service_miles,
      }),
    })

    if (response.ok) {
      showAddReminder.value = false
      await loadVehicleData()
    }
  } catch (error) {
    console.error('Error submitting reminder:', error)
  }
}

async function completeReminder(reminder) {
  try {
    await fetch(`http://localhost:3000/api/reminders/${reminder.id}/complete`, {
      method: 'POST',
      headers: {
        'Authorization': authStore.token,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        service_date: new Date(),
        service_miles: vehicle.value.odometer,
      }),
    })
    await loadVehicleData()
  } catch (error) {
    console.error('Error completing reminder:', error)
  }
}

async function deleteFuel(id) {
  if (confirm('Delete this fuel entry?')) {
    await fetch(`http://localhost:3000/api/fuel/${id}`, {
      method: 'DELETE',
      headers: { 'Authorization': authStore.token },
    })
    await loadVehicleData()
  }
}

async function deleteExpense(id) {
  if (confirm('Delete this expense?')) {
    await fetch(`http://localhost:3000/api/expenses/${id}`, {
      method: 'DELETE',
      headers: { 'Authorization': authStore.token },
    })
    await loadVehicleData()
  }
}

async function deleteReminder(id) {
  if (confirm('Delete this reminder?')) {
    await fetch(`http://localhost:3000/api/reminders/${id}`, {
      method: 'DELETE',
      headers: { 'Authorization': authStore.token },
    })
    await loadVehicleData()
  }
}

function editFuel(entry) {
  // TODO: Edit modal
}

function editExpense(expense) {
  // TODO: Edit modal
}

function editReminder(reminder) {
  // TODO: Edit modal
}

async function exportCSV() {
  const response = await fetch(`http://localhost:3000/api/export/csv`, {
    headers: { 'Authorization': authStore.token },
  })
  const blob = await response.blob()
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'clarkson-export.csv'
  a.click()
}

async function exportPDF() {
  alert('PDF export coming soon!')
}

function formatDate(date) {
  return new Date(date).toLocaleDateString()
}
</script>
