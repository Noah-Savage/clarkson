<template>
  <div v-if="alerts.length" class="space-y-3">
    <div v-for="alert in alerts" :key="alert.reminderID"
         :class="['p-4 rounded-lg border-l-4', getAlertClass(alert)]">
      <div class="flex justify-between items-start">
        <div>
          <p class="font-bold text-lg">{{ alert.name }}</p>
          <p class="text-sm">{{ alert.vehicleName }}</p>
          <div class="mt-2 text-sm space-y-1">
            <p v-if="alert.milesToGo">📊 {{ Math.round(alert.milesToGo) }} miles remaining</p>
            <p v-if="alert.daysUntilDue">📅 {{ alert.daysUntilDue }} days remaining</p>
          </div>
        </div>
        <div>
          <span :class="['px-3 py-1 rounded-full text-white text-sm font-semibold', getStatusColor(alert)]">
            {{ alert.status.toUpperCase() }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  alerts: {
    type: Array,
    default: () => [],
  },
})

function getAlertClass(alert) {
  if (alert.status === 'overdue') return 'bg-red-50 dark:bg-red-900 border-red-500'
  if (alert.status === 'soon') return 'bg-yellow-50 dark:bg-yellow-900 border-yellow-500'
  return 'bg-gray-50 dark:bg-gray-900 border-gray-500'
}

function getStatusColor(alert) {
  if (alert.status === 'overdue') return 'bg-red-600'
  if (alert.status === 'soon') return 'bg-yellow-600'
  return 'bg-gray-600'
}
</script>
