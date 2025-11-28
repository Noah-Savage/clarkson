import { defineStore } from "pinia"

const API_URL = "http://localhost:3000/api"

export const useVehicleStore = defineStore("vehicles", () => {
  async function fetchVehicles(token) {
    const response = await fetch(`${API_URL}/vehicles`, {
      headers: { Authorization: token },
    })
    return response.json()
  }

  async function createVehicle(token, vehicleData) {
    const response = await fetch(`${API_URL}/vehicles`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: token,
      },
      body: JSON.stringify(vehicleData),
    })
    return response.json()
  }

  async function getVehicle(token, id) {
    const response = await fetch(`${API_URL}/vehicles/${id}`, {
      headers: { Authorization: token },
    })
    return response.json()
  }

  async function updateVehicle(token, id, vehicleData) {
    const response = await fetch(`${API_URL}/vehicles/${id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: token,
      },
      body: JSON.stringify(vehicleData),
    })
    return response.json()
  }

  async function deleteVehicle(token, id) {
    const response = await fetch(`${API_URL}/vehicles/${id}`, {
      method: "DELETE",
      headers: { Authorization: token },
    })
    return response.json()
  }

  async function checkReminders(token) {
    const response = await fetch(`${API_URL}/reminders/check`, {
      headers: { Authorization: token },
    })
    const data = await response.json()
    return data.alerts || []
  }

  return {
    fetchVehicles,
    createVehicle,
    getVehicle,
    updateVehicle,
    deleteVehicle,
    checkReminders,
  }
})
