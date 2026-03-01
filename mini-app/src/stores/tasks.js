import { defineStore } from 'pinia'
import api from '../services/api'

export const useTaskStore = defineStore('tasks', {
  state: () => ({
    tasks: [],
    loading: false,
    error: null,
    lastFetched: null
  }),
  actions: {
    async fetchTasks(filters = {}, force = false) {
      const now = Date.now()
      const CACHE_TIMEOUT = 5 * 60 * 1000 // 5 minutes
      const REVALIDATE_TIMEOUT = 30000 // 30 seconds

      // If we have data and it's fresh (within 30s), just return
      if (!force && this.tasks.length > 0 && this.lastFetched && (now - this.lastFetched < REVALIDATE_TIMEOUT)) {
        return
      }

      // If we have data but it's stale (between 30s and 5m), return immediately but fetch in background
      if (!force && this.tasks.length > 0 && this.lastFetched && (now - this.lastFetched < CACHE_TIMEOUT)) {
        // Trigger background fetch without setting loading = true
        api.get('/tasks', { params: filters }).then(response => {
          this.tasks = response.data
          this.lastFetched = Date.now()
        }).catch(err => console.error('Background fetchTasks failed:', err))
        return
      }

      // Otherwise, do a full fetch with loading state
      this.loading = true
      try {
        const response = await api.get('/tasks', { params: filters })
        this.tasks = response.data
        this.lastFetched = now
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },
    async createTask(task) {
      try {
        // Build a clean payload for the Go backend
        const payload = {
          title: task.title,
          priority: task.priority || 'medium',
          description: task.description || '',
          is_daily: !!task.is_daily,
        }

        // Only include optional foreign keys if they have a value
        if (task.project_id) payload.project_id = task.project_id
        if (task.goal_id) payload.goal_id = task.goal_id

        // Convert datetime-local strings to RFC 3339 format for Go's time.Time
        if (task.start_time) {
          payload.start_time = new Date(task.start_time).toISOString()
        }
        if (task.end_time) {
          payload.end_time = new Date(task.end_time).toISOString()
        }

        const response = await api.post('/tasks', payload)
        this.tasks.push(response.data)
        return response.data
      } catch (err) {
        this.error = err.message
        throw err
      }
    },
    async updateTask(id, task) {
      try {
        const response = await api.put(`/tasks/${id}`, task)
        const index = this.tasks.findIndex(t => t.id === id)
        if (index !== -1) this.tasks[index] = response.data
        return response.data
      } catch (err) {
        this.error = err.message
        throw err
      }
    },
    async deleteTask(id) {
      try {
        await api.delete(`/tasks/${id}`)
        this.tasks = this.tasks.filter(t => t.id !== id)
      } catch (err) {
        this.error = err.message
        throw err
      }
    },
    async toggleTaskStatus(task) {
      const newStatus = task.status === 'completed' ? 'pending' : 'completed'
      return this.updateTask(task.id, { ...task, status: newStatus })
    }
  }
})
