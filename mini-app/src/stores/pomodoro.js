import { defineStore } from 'pinia'
import api from '../services/api'

export const usePomodoroStore = defineStore('pomodoro', {
  state: () => ({
    status: null,
    loading: false,
    error: null,
  }),

  actions: {
    async fetchStatus() {
      this.loading = true
      try {
        const response = await api.get('/pomodoro/status')
        this.status = response.data
        this.error = null
        return response.data
      } catch (error) {
        this.error = "Failed to fetch pomodoro status"
        console.error('Error fetching pomodoro status:', error)
        throw error
      } finally {
        this.loading = false
      }
    },

    async startSession(projectId = null) {
      try {
        const response = await api.post('/pomodoro/session/start', {
          project_id: projectId
        })
        await this.fetchStatus()
        return response.data
      } catch (error) {
        console.error('Error starting pomodoro session:', error)
        throw error
      }
    },

    async logCycle() {
      try {
        const response = await api.post('/pomodoro/session/cycle')
        await this.fetchStatus()
        return response.data
      } catch (error) {
        console.error('Error logging cycle:', error)
        throw error
      }
    },

    async resumeWork() {
      try {
        const response = await api.post('/pomodoro/session/resume-work')
        await this.fetchStatus()
        return response.data
      } catch (error) {
        console.error('Error resuming work:', error)
        throw error
      }
    },

    async pauseSession() {
      try {
        const response = await api.post('/pomodoro/session/pause')
        await this.fetchStatus()
        return response.data
      } catch (error) {
        console.error('Error pausing session:', error)
        throw error
      }
    },

    async resumeSession() {
      try {
        const response = await api.post('/pomodoro/session/resume')
        await this.fetchStatus()
        return response.data
      } catch (error) {
        console.error('Error resuming session:', error)
        throw error
      }
    },

    async completeLevel() {
      try {
        const response = await api.post('/pomodoro/session/complete')
        await this.fetchStatus()
        return response.data
      } catch (error) {
        console.error('Error completing level:', error)
        throw error
      }
    },

    async abandonSession() {
      try {
        await api.post('/pomodoro/session/abandon')
        await this.fetchStatus()
      } catch (error) {
        console.error('Error abandoning session:', error)
        throw error
      }
    }
  }
})
