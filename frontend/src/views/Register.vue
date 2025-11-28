<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary to-primary-dark">
    <div class="bg-white dark:bg-secondary rounded-lg shadow-xl p-8 max-w-md w-full">
      <h1 class="text-3xl font-bold text-center mb-8">Create Account</h1>

      <form @submit.prevent="handleRegister" class="space-y-4">
        <div>
          <label class="block text-sm font-medium mb-1">Full Name</label>
          <input v-model="name" type="text" class="w-full px-4 py-2 border rounded-lg dark:bg-secondary focus:ring-2 focus:ring-primary" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Email</label>
          <input v-model="email" type="email" class="w-full px-4 py-2 border rounded-lg dark:bg-secondary focus:ring-2 focus:ring-primary" />
        </div>
        <div>
          <label class="block text-sm font-medium mb-1">Password (min 8 chars)</label>
          <input v-model="password" type="password" class="w-full px-4 py-2 border rounded-lg dark:bg-secondary focus:ring-2 focus:ring-primary" />
        </div>

        <button type="submit" class="w-full bg-primary text-white py-2 rounded-lg font-semibold hover:bg-primary-dark transition">
          Register
        </button>
      </form>

      <p class="text-center text-gray-600 dark:text-gray-400 mt-4">
        Already have an account?
        <router-link to="/login" class="text-primary hover:underline font-semibold">Login</router-link>
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

const name = ref('')
const email = ref('')
const password = ref('')

async function handleRegister() {
  const success = await authStore.register(email.value, password.value, name.value)
  if (success) {
    router.push('/login')
  } else {
    alert('Registration failed')
  }
}
</script>
