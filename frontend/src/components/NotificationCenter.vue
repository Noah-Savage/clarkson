<template>
  <div class="fixed top-16 right-4 w-96 max-h-96 bg-white dark:bg-secondary rounded-lg shadow-xl border border-gray-200 dark:border-gray-700 overflow-hidden z-40">
    <!-- Header -->
    <div class="bg-primary text-white px-4 py-3 flex justify-between items-center">
      <h3 class="font-bold">Notifications</h3>
      <button @click="$emit('close')" class="text-xl hover:opacity-75">×</button>
    </div>

    <!-- Notifications List -->
    <div class="overflow-y-auto max-h-80">
      <div v-if="notifications.length === 0" class="p-4 text-center text-gray-500">
        No notifications
      </div>

      <div v-for="notif in notifications" :key="notif.id"
           :class="['p-3 border-b border-gray-100 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800 transition cursor-pointer', notif.status === 'unread' ? 'bg-blue-50 dark:bg-blue-900' : '']">
        <div class="flex items-start space-x-3">
          <div :class="['w-2 h-2 rounded-full mt-1', notif.type === 'reminder_overdue' ? 'bg-red-600' : 'bg-yellow-600']"></div>
          <div class="flex-1">
            <p class="font-semibold text-sm">{{ notif.title }}</p>
            <p class="text-xs text-gray-600 dark:text-gray-400">{{ notif.message }}</p>
            <p class="text-xs text-gray-500 dark:text-gray-500 mt-1">{{ formatTime(notif.created_at) }}</p>
          </div>
          <button @click.stop="dismiss(notif.id)" class="text-gray-400 hover:text-gray-600 text-lg">×</button>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <div class="p-3 bg-gray-50 dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 flex space-x-2">
      <button @click="markAllRead" class="flex-1 text-sm bg-primary text-white px-3 py-1 rounded hover:bg-primary-dark">
        Mark all read
      </button>
      <button @click="clearAll" class="flex-1 text-sm bg-gray-300 dark:bg-gray-600 px-3 py-1 rounded hover:bg-gray-400">
        Clear all
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const notifications = ref([])

defineEmits(['close'])

onMounted(async () => {
  await loadNotifications()
  // Poll for new notifications every 30 seconds
  setInterval(loadNotifications, 30000)
})

async function loadNotifications() {
  try {
    const response = await fetch('http://localhost:3000/api/notifications', {
      headers: { 'Authorization': authStore.token },
    })
    notifications.value = await response.json()
  } catch (error) {
    console.error('Error loading notifications:', error)
  }
}

async function dismiss(id) {
  await fetch(`http://localhost:3000/api/notifications/${id}/dismiss`, {
    method: 'POST',
    headers: { 'Authorization': authStore.token },
  })
  await loadNotifications()
}

async function markAllRead() {
  for (const notif of notifications.value.filter(n => n.status === 'unread')) {
    await fetch(`http://localhost:3000/api/notifications/${notif.id}/read`, {
      method: 'POST',
      headers: { 'Authorization': authStore.token },
    })
  }
  await loadNotifications()
}

async function clearAll() {
  for (const notif of notifications.value) {
    await fetch(`http://localhost:3000/api/notifications/${notif.id}/dismiss`, {
      method: 'POST',
      headers: { 'Authorization': authStore.token },
    })
  }
  await loadNotifications()
}

function formatTime(date) {
  const now = new Date()
  const then = new Date(date)
  const diff = Math.floor((now - then) / 1000)

  if (diff < 60) return 'just now'
  if (diff < 3600) return `${Math.floor(diff / 60)}m ago`
  if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`
  return `${Math.floor(diff / 86400)}d ago`
}
</script>
