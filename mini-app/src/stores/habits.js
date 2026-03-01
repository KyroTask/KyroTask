import { defineStore } from 'pinia'
import api from '../services/api'

export const useHabitStore = defineStore('habits', {
  state: () => ({
    habits: [],
    loading: false,
    error: null,
    lastFetched: null
  }),
  actions: {
    async fetchHabits(force = false) {
      const now = Date.now()
      const CACHE_TIMEOUT = 5 * 60 * 1000 // 5 minutes
      const REVALIDATE_TIMEOUT = 30000 // 30 seconds

      if (!force && this.habits.length > 0 && this.lastFetched && (now - this.lastFetched < REVALIDATE_TIMEOUT)) {
        return
      }

      if (!force && this.habits.length > 0 && this.lastFetched && (now - this.lastFetched < CACHE_TIMEOUT)) {
        api.get('/habits').then(response => {
          this.habits = response.data
          this.lastFetched = Date.now()
        }).catch(err => console.error('Background fetchHabits failed:', err))
        return
      }

      this.loading = true
      try {
        const response = await api.get('/habits')
        this.habits = response.data
        this.lastFetched = now
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },
    async createHabit(habit) {
      try {
        const response = await api.post('/habits', habit)
        this.habits.push(response.data)
        return response.data
      } catch (err) {
        this.error = err.message
        throw err
      }
    },
    async logHabit(id, date) {
      try {
        // Go backend expects RFC3339 formatted dates, e.g. "2026-02-28T00:00:00Z"
        // date parameter is typically "YYYY-MM-DD"
        const formattedDate = date.includes('T') ? date : `${date}T00:00:00Z`
        const response = await api.post(`/habits/${id}/log`, { log_date: formattedDate })
        const habit = this.habits.find(h => h.id === id)
        if (habit) {
          if (!habit.logs) habit.logs = []
          habit.logs.push(response.data)
        }
        return response.data
      } catch (err) {
        this.error = err.message
        throw err
      }
    },
    async deleteHabit(id) {
      try {
        await api.delete(`/habits/${id}`)
        this.habits = this.habits.filter(h => h.id !== id)
      } catch (err) {
        this.error = err.message
        throw err
      }
    },
    async updateHabit(id, habit) {
      try {
        const response = await api.put(`/habits/${id}`, habit)
        const index = this.habits.findIndex(h => h.id === id)
        if (index !== -1) this.habits[index] = response.data
        return response.data
      } catch (err) {
        this.error = err.message
        throw err
      }
    }
  }
})
