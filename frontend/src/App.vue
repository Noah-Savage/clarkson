<template>
  <div id="app" :class="[isDarkMode ? 'dark' : 'light']" class="min-h-screen transition-colors">
    <!-- Navigation -->
    <nav class="bg-primary text-white shadow-lg">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center space-x-2">
            <span class="text-2xl font-bold">Clarkson</span>
          </div>
          <div class="flex items-center space-x-4">
            <NotificationBell v-if="authStore.token" />
            <button @click="toggleDarkMode" class="p-2 hover:bg-primary-dark rounded transition">
              {{ isDarkMode ? '☀️' : '🌙' }}
            </button>
            <router-link v-if="authStore.token" to="/settings" class="p-2 hover:bg-primary-dark rounded transition">⚙️</router-link>
            <button v-if="authStore.token" @click="logout" class="p-2 hover:bg-primary-dark rounded transition">Logout</button>
          </div>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeMount } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'
import NotificationBell from './components/NotificationBell.vue'

const router = useRouter()
const authStore = useAuthStore()
const isDarkMode = ref(false)

onBeforeMount(() => {
  if (!authStore.token && router.currentRoute.value.path !== '/login' && router.currentRoute.value.path !== '/register') {
    router.push('/login')
  }
})

onMounted(() => {
  const saved = localStorage.getItem('darkMode')
  if (saved !== null) {
    isDarkMode.value = JSON.parse(saved)
  } else {
    isDarkMode.value = window.matchMedia('(prefers-color-scheme: dark)').matches
  }
  applyDarkMode()
})

function toggleDarkMode() {
  isDarkMode.value = !isDarkMode.value
  localStorage.setItem('darkMode', JSON.stringify(isDarkMode.value))
  applyDarkMode()
}

function applyDarkMode() {
  const html = document.documentElement
  if (isDarkMode.value) {
    html.classList.add('dark')
  } else {
    html.classList.remove('dark')
  }
}

function logout() {
  authStore.logout()
  router.push('/login')
}
</script>

<style>
:root {
  --color-primary: #2563eb;
  --color-primary-dark: #1e40af;
  --color-secondary: #64748b;
  --color-accent: #06b6d4;
  --color-success: #10b981;
  --color-warning: #f59e0b;
  --color-danger: #ef4444;

  --bg-light: #ffffff;
  --bg-secondary: #f8fafc;
  --text-light: #1e293b;
  --text-secondary: #64748b;
  --border-light: #e2e8f0;
}

html.dark {
  --bg-light: #0f172a;
  --bg-secondary: #1e293b;
  --text-light: #f1f5f9;
  --text-secondary: #cbd5e1;
  --border-light: #334155;
}

body {
  background-color: var(--bg-light);
  color: var(--text-light);
  transition: background-color 0.3s, color 0.3s;
  font-family: system-ui, -apple-system, sans-serif;
}

.bg-primary {
  background-color: var(--color-primary);
}

.text-primary {
  color: var(--color-primary);
}
</style>
