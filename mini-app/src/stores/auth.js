import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('jwt') || null)
  const user = ref(null)
  const loading = ref(false)
  const error = ref(null)
  const theme = ref(localStorage.getItem('theme') || 'light')

  const isAuthenticated = computed(() => !!token.value)

  function setTheme(newTheme) {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    applyTheme()
  }

  function applyTheme() {
    const isDark = theme.value === 'dark' ||
      (theme.value === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)

    if (isDark) {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  // Listen for system theme changes
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
    if (theme.value === 'system') {
      applyTheme()
    }
  })

  async function loginWithTelegram(initData) {
    loading.value = true
    error.value = null

    try {
      const response = await api.post('/auth/telegram/verify', { initData })
      token.value = response.data.token
      user.value = response.data.user
      localStorage.setItem('jwt', token.value)
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Login failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function loginWithWidget(user_data) {
    loading.value = true
    error.value = null

    try {
      const response = await api.post('/auth/telegram/widget', user_data)
      token.value = response.data.token
      user.value = response.data.user
      localStorage.setItem('jwt', token.value)
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Login failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function loginWithFirebase(idToken, linkAccount = false) {
    loading.value = true
    error.value = null

    try {
      const response = await api.post('/auth/firebase', {
        id_token: idToken,
        link_account: linkAccount
      })

      if (!linkAccount) {
        token.value = response.data.token
        user.value = response.data.user
        localStorage.setItem('jwt', token.value)
      } else {
        await fetchCurrentUser(true)
      }
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Firebase authentication failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  const lastFetched = ref(0)
  const CACHE_TIMEOUT = 5 * 60 * 1000 // 5 minutes

  async function fetchCurrentUser(force = false) {
    const now = Date.now()
    if (!force && user.value && (now - lastFetched.value < CACHE_TIMEOUT)) {
      return user.value
    }

    try {
      const response = await api.get('/user/me')
      user.value = response.data
      lastFetched.value = now
      return response.data
    } catch (err) {
      console.error('Failed to fetch user:', err)
      throw err
    }
  }

  function logout() {
    token.value = null
    user.value = null
    localStorage.removeItem('jwt')
  }

  async function loginDev() {
    loading.value = true
    error.value = null

    try {
      const response = await api.post('/auth/dev-login')
      token.value = response.data.token
      user.value = response.data.user
      localStorage.setItem('jwt', token.value)
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Dev login failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    token,
    user,
    loading,
    error,
    theme,
    isAuthenticated,
    loginWithTelegram,
    loginWithWidget,
    loginWithFirebase,
    loginDev,
    fetchCurrentUser,
    logout,
    setTheme,
    applyTheme
  }
})
