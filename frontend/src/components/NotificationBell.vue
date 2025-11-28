<template>
  <div class="relative">
    <button @click="toggleNotifications" class="relative p-2 hover:bg-primary-dark rounded transition">
      <span class="text-2xl">🔔</span>
      <span v-if="unreadCount > 0" class="absolute top-0 right-0 bg-red-600 text-white text-xs font-bold rounded-full w-5 h-5 flex items-center justify-center">
        {{ unreadCount > 9 ? '9+' : unreadCount }}
      </span>
    </button>

    <NotificationCenter v-if="showNotifications" @close="showNotifications = false" />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import NotificationCenter from './NotificationCenter.vue'

const authStore = useAuthStore()
const showNotifications = ref(false)
const unreadCount = ref(0)
let pollInterval = null

onMounted(() => {
  loadNotificationSummary()
  // Poll for notification updates every 30 seconds
  pollInterval = setInterval(loadNotificationSummary, 30000)
})

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval)
})

async function loadNotificationSummary() {
  try {
    const response = await fetch('http://localhost:3000/api/notifications/summary', {
      headers: { 'Authorization': authStore.token },
    })
    const data = await response.json()
    unreadCount.value = data.unread_count
  } catch (error) {
    console.error('Error loading notification summary:', error)
  }
}

function toggleNotifications() {
  showNotifications.value = !showNotifications.value
}
</script>
