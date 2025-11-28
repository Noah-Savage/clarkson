import { defineStore } from "pinia"

const API_URL = "http://localhost:3000/api"

export const useFuelStore = defineStore("fuel", () => {
  async function listFuelEntries(token, vehicleId) {
    const response = await fetch(`${API_URL}/vehicles/${vehicleId}/fuel`, {
      headers: { Authorization: token },
    })
    return response.json()
  }

  async function createFuelEntry(token, vehicleId, data) {
    const response = await fetch(`${API_URL}/vehicles/${vehicleId}/fuel`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: token,
      },
      body: JSON.stringify(data),
    })
    return response.json()
  }

  async function getFuelStats(token, vehicleId) {
    const response = await fetch(`${API_URL}/vehicles/${vehicleId}/fuel-stats`, {
      headers: { Authorization: token },
    })
    return response.json()
  }

  async function deleteFuelEntry(token, fuelId) {
    const response = await fetch(`${API_URL}/fuel/${fuelId}`, {
      method: "DELETE",
      headers: { Authorization: token },
    })
    return response.json()
  }

  return {
    listFuelEntries,
    createFuelEntry,
    getFuelStats,
    deleteFuelEntry,
  }
})
