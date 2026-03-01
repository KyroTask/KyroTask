import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import MainLayout from '../layouts/MainLayout.vue'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../pages/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: MainLayout,
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('../pages/Dashboard.vue')
      },
      {
        path: 'projects',
        name: 'Projects',
        component: () => import('../pages/Projects.vue')
      },
      {
        path: 'projects/:slug',
        name: 'ProjectDetail',
        component: () => import('../pages/ProjectDetail.vue')
      },
      {
        path: 'tasks',
        name: 'Tasks',
        component: () => import('../pages/Tasks.vue')
      },
      {
        path: 'tasks/:slug',
        name: 'TaskDetail',
        component: () => import('../pages/TaskDetail.vue')
      },
      {
        path: 'goals',
        name: 'Goals',
        component: () => import('../pages/Goals.vue')
      },
      {
        path: 'goals/:slug',
        name: 'GoalDetail',
        component: () => import('../pages/GoalDetail.vue')
      },
      {
        path: 'habits',
        name: 'Habits',
        component: () => import('../pages/Habits.vue')
      },
      {
        path: 'pomodoro',
        name: 'Pomodoro',
        component: () => import('../pages/Pomodoro.vue')
      },
      {
        path: 'calendar',
        name: 'Calendar',
        component: () => import('../pages/Calendar.vue')
      },
      {
        path: 'activity',
        name: 'Activity',
        component: () => import('../pages/Activity.vue')
      },
      {
        path: 'analytics',
        name: 'Analytics',
        component: () => import('../pages/Analytics.vue')
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('../pages/Settings.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard for authentication
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router
