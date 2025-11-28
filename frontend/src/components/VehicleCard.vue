<template>
  <div class="bg-white dark:bg-secondary rounded-lg shadow-md p-6 hover:shadow-lg transition border border-gray-200 dark:border-gray-700">
    <div class="flex justify-between items-start mb-4">
      <div class="flex-1">
        <h3 class="text-xl font-bold">{{ vehicle.year }} {{ vehicle.make }} {{ vehicle.model }}</h3>
        <p class="text-sm text-gray-600 dark:text-gray-400">{{ vehicle.odometer }} {{ vehicle.mileageUnit }}</p>
      </div>
      <div class="flex space-x-2">
        <button @click="$emit('edit')" class="text-blue-500 hover:text-blue-700 text-lg">✏️</button>
        <button @click="$emit('delete')" class="text-red-500 hover:text-red-700 text-lg">🗑️</button>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-3 mb-4">
      <div class="bg-blue-50 dark:bg-blue-900 p-3 rounded">
        <div class="text-xs text-gray-600 dark:text-gray-400">Avg MPG</div>
        <div class="text-lg font-bold text-blue-600 dark:text-blue-400">{{ avgMPG }}</div>
      </div>
      <div class="bg-green-50 dark:bg-green-900 p-3 rounded">
        <div class="text-xs text-gray-600 dark:text-gray-400">Total Cost</div>
        <div class="text-lg font-bold text-green-600 dark:text-green-400">${{ totalCost }}</div>
      </div>
      <div class="bg-purple-50 dark:bg-purple-900 p-3 rounded">
        <div class="text-xs text-gray-600 dark:text-gray-400">Entries</div>
        <div class="text-lg font-bold text-purple-600 dark:text-purple-400">{{ entryCount }}</div>
      </div>
      <div v-if="dueReminders > 0" class="bg-orange-50 dark:bg-orange-900 p-3 rounded">
        <div class="text-xs text-gray-600 dark:text-gray-400">Due</div>
        <div class="text-lg font-bold text-orange-600 dark:text-orange-400">{{ dueReminders }}</div>
      </div>
    </div>

    <router-link :to="`/vehicle/${vehicle.id}`" class="w-full bg-primary text-white text-center px-3 py-2 rounded hover:bg-primary-dark block transition">
      View Details
    </router-link>
  </div>
</template>

<script setup>
defineProps({
  vehicle: {
    type: Object,
    required: true,
  },
  avgMPG: {
    type: String,
    default: '-',
  },
  totalCost: {
    type: String,
    default: '0.00',
  },
  entryCount: {
    type: Number,
    default: 0,
  },
  dueReminders: {
    type: Number,
    default: 0,
  },
})

defineEmits(['edit', 'delete'])
</script>
