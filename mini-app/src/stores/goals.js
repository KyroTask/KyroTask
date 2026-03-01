import { defineStore } from 'pinia'
import api from '../services/api'

export const useGoalStore = defineStore('goals', {
  state: () => ({
    goals: [],
    loading: false,
    error: null,
    lastFetched: null
  }),
  actions: {
    async fetchGoals(filters = {}, force = false) {
      const now = Date.now()
      const CACHE_TIMEOUT = 5 * 60 * 1000 // 5 minutes
      const REVALIDATE_TIMEOUT = 30000 // 30 seconds

      if (!force && this.goals.length > 0 && this.lastFetched && (now - this.lastFetched < REVALIDATE_TIMEOUT)) {
        return
      }

      if (!force && this.goals.length > 0 && this.lastFetched && (now - this.lastFetched < CACHE_TIMEOUT)) {
        api.get('/goals', { params: filters }).then(response => {
          this.goals = response.data
          this.lastFetched = Date.now()
        }).catch(err => console.error('Background fetchGoals failed:', err))
        return
      }

      this.loading = true
      try {
        const response = await api.get('/goals', { params: filters })
        this.goals = response.data
        this.lastFetched = now
      } catch (err) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    },
    async createGoal(goal) {
      try {
        const response = await api.post('/goals', goal)
        this.goals.push(response.data)
        return response.data
      } catch (err) {
        this.error = err.message
        throw err
      }
    },
    async deleteGoal(id) {
      try {
        await api.delete(`/goals/${id}`)
        this.goals = this.goals.filter(g => g.id !== id)
      } catch (err) {
        this.error = err.message
        throw err
      }
    },
    async updateGoal(id, goal) {
      try {
        const response = await api.put(`/goals/${id}`, goal)
        const index = this.goals.findIndex(g => g.id === id)
        if (index !== -1) this.goals[index] = response.data
        return response.data
      } catch (err) {
        this.error = err.message
        throw err
      }
    },
    async fetchGoal(idOrSlug) {
      this.loading = true
      try {
        // Check local state first
        let goal = this.goals.find(g => g.id === Number(idOrSlug) || g.slug === idOrSlug)
        
        // If not found or navigating directly, fetch from API
        const response = await api.get(`/goals/${idOrSlug}`)
        const fetchedGoal = response.data
        
        // Update local state
        const index = this.goals.findIndex(g => g.id === fetchedGoal.id)
        if (index !== -1) {
          this.goals[index] = fetchedGoal
        } else {
          this.goals.push(fetchedGoal)
        }
        
        return fetchedGoal
      } catch (err) {
        this.error = err.message
        throw err
      } finally {
        this.loading = false
      }
    },
    async createMilestone(milestone) {
      try {
        const response = await api.post('/milestones', milestone)
        // Find the goal and add the milestone
        const goal = this.goals.find(g => g.id === milestone.goal_id)
        if (goal) {
          if (!goal.milestones) goal.milestones = []
          goal.milestones.push(response.data)
        }
        return response.data
      } catch (err) {
        throw err
      }
    },
    async updateMilestone(id, updates) {
      try {
        const response = await api.put(`/milestones/${id}`, updates)
        // Update local state
        for (const goal of this.goals) {
          if (goal.milestones) {
            const index = goal.milestones.findIndex(m => m.id === id)
            if (index !== -1) {
              goal.milestones[index] = response.data
              break
            }
          }
        }
        return response.data
      } catch (err) {
        throw err
      }
    },
    async deleteMilestone(id) {
      try {
        await api.delete(`/milestones/${id}`)
        // Remove from local state
        for (const goal of this.goals) {
          if (goal.milestones) {
            goal.milestones = goal.milestones.filter(m => m.id !== id)
          }
        }
      } catch (err) {
        throw err
      }
    }
  }
})
