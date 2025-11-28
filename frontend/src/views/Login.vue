<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary to-primary-dark">
    <div class="bg-white dark:bg-secondary rounded-lg shadow-xl p-8 max-w-md w-full">
      <h1 class="text-3xl font-bold text-center mb-2">Clarkson</h1>
      <p class="text-center text-gray-600 dark:text-gray-400 mb-8">Vehicle Expense Tracker</p>

      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="block text-sm font-medium mb-1">Email</label>
          <input v-model="email" type="email" class="w-full px-4 py-2 border rounded-lg dark:bg-secondary focus:ring-2 focus:ring-primary" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Password</label>
          <input v-model="password" type="password" class="w-full px-4 py-2 border rounded-lg dark:bg-secondary focus:ring-2 focus:ring-primary" />
        </div>

        <button type="submit" class="w-full bg-primary text-white py-2 rounded-lg font-semibold hover:bg-primary-dark transition">
          Login
        </button>
      </form>

      <p class="text-center text-gray-600 dark:text-gray-400 mt-4">
        Don't have an account?
        <router-link to="/register" class="text-primary hover:underline font-semibold">Register</router-link>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')

async function handleLogin() {
  const success = await authStore.login(email.value, password.value)
  if (success) {
    router.push('/dashboard')
  } else {
    alert('Login failed')
  }
}
</script>
