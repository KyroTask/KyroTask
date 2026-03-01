import { defineStore } from 'pinia'
import api from '../services/api'

export const useProjectStore = defineStore('projects', {
  state: () => ({
    projects: [],
    loading: false,
    error: null,
    lastFetched: null
  }),
  actions: {
    async fetchProjects(force = false) {
      const now = Date.now()
      const CACHE_TIMEOUT = 5 * 60 * 1000 // 5 minutes
      const REVALIDATE_TIMEOUT = 30000 // 30 seconds

      if (!force && this.projects.length > 0 && this.lastFetched && (now - this.lastFetched < REVALIDATE_TIMEOUT)) {
        return
      }

      if (!force && this.projects.length > 0 && this.lastFetched && (now - this.lastFetched < CACHE_TIMEOUT)) {
        api.get('/projects').then(response => {
          this.projects = response.data
          this.lastFetched = Date.now()
        }).catch(err => console.error('Background fetchProjects failed:', err))
        return
      }

      this.loading = true
      try {
        const response = await api.get('/projects')
        this.projects = response.data
        this.lastFetched = now
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },
    async createProject(project) {
      try {
        const response = await api.post('/projects', project)
        this.projects.push(response.data)
        return response.data
      } catch (err) {
        this.error = err.message
        throw err
      }
    },
    async updateProject(id, project) {
      try {
        const response = await api.put(`/projects/${id}`, project)
        const index = this.projects.findIndex(p => p.id === id)
        if (index !== -1) this.projects[index] = response.data
        return response.data
      } catch (err) {
        this.error = err.message
        throw err
      }
    },
    async deleteProject(id) {
      try {
        await api.delete(`/projects/${id}`)
        this.projects = this.projects.filter(p => p.id !== id)
      } catch (err) {
        this.error = err.message
        throw err
      }
    }
  }
})
