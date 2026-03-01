import { defineStore } from 'pinia'
import api from '../services/api'

export const useAnalyticsStore = defineStore('analytics', {
  state: () => ({
    analytics: null,
    loading: false,
    error: null
  }),

  actions: {
    async fetchAnalytics() {
      this.loading = true
      this.error = null
      try {
        const response = await api.get('/analytics')
        this.analytics = response.data
      } catch (error) {
        this.error = error.response?.data?.error || 'Failed to fetch analytics data'
        console.error('Error fetching analytics:', error)
      } finally {
        this.loading = false
      }
    }
  }
})
