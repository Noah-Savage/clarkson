import { defineStore } from "pinia"

const API_URL = "http://localhost:3000/api"

export const useReminderStore = defineStore("reminders", () => {
  async function listReminders(token, vehicleId) {
    const response = await fetch(`${API_URL}/vehicles/${vehicleId}/reminders`, {
      headers: { Authorization: token },
    })
    return response.json()
  }

  async function createReminder(token, vehicleId, data) {
    const response = await fetch(`${API_URL}/vehicles/${vehicleId}/reminders`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: token,
      },
      body: JSON.stringify(data),
    })
    return response.json()
  }

  async function completeReminder(token, reminderId, data) {
    const response = await fetch(`${API_URL}/reminders/${reminderId}/complete`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: token,
      },
      body: JSON.stringify(data),
    })
    return response.json()
  }

  async function deleteReminder(token, reminderId) {
    const response = await fetch(`${API_URL}/reminders/${reminderId}`, {
      method: "DELETE",
      headers: { Authorization: token },
    })
    return response.json()
  }

  async function getOverdueReminders(token) {
    const response = await fetch(`${API_URL}/reminders/overdue`, {
      headers: { Authorization: token },
    })
    return response.json()
  }

  return {
    listReminders,
    createReminder,
    completeReminder,
    deleteReminder,
    getOverdueReminders,
  }
})
