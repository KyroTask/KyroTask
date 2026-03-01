import { defineStore } from 'pinia'
import api from '../services/api'

export const useCalendarStore = defineStore('calendar', {
  state: () => ({
    events: [],
    loading: false,
    error: null,
    lastRange: null, // { start, end }
  }),
  actions: {
    async fetchEvents(start, end) {
      // Format dates as YYYY-MM-DD
      const startStr = start instanceof Date ? start.toISOString().split('T')[0] : start
      const endStr = end instanceof Date ? end.toISOString().split('T')[0] : end

      // Skip if we already fetched this exact range
      if (this.lastRange && this.lastRange.start === startStr && this.lastRange.end === endStr && this.events.length > 0) {
        return
      }

      this.loading = true
      try {
        const response = await api.get('/calendar', {
          params: { start: startStr, end: endStr }
        })
        this.events = response.data.events || []
        this.lastRange = { start: startStr, end: endStr }
      } catch (err) {
        this.error = err.message
        console.error('Failed to fetch calendar events:', err)
      } finally {
        this.loading = false
      }
    },

    getEventsForDay(date) {
      const dateStr = date instanceof Date
        ? date.toISOString().split('T')[0]
        : date

      return this.events.filter(event => {
        const eventDate = new Date(event.date)
        return eventDate.toISOString().split('T')[0] === dateStr
      })
    }
  }
})
