<template>
  <div class="space-y-6 animate-fade-in pb-20 md:pb-6">
    <!-- Header -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white transition-colors">Analytics</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1 transition-colors">
          Track your progress across tasks, habits, pomodoros, and goals.
        </p>
      </div>
    </header>

    <!-- Error State -->
    <div v-if="analyticsStore.error" class="bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 p-4 rounded-xl border border-red-100 dark:border-red-900/30 text-sm">
      {{ analyticsStore.error }}
    </div>

    <!-- Loading State -->
    <div v-if="analyticsStore.loading" class="flex justify-center items-center py-20">
      <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-primary-500"></div>
    </div>

    <!-- Dashboard Content -->
    <div v-else-if="analyticsStore.analytics" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-2 gap-6">

      <!-- Pomodoros Card -->
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6 transition-all hover:shadow-md relative overflow-hidden">
        <!-- Decorative subtle background ring for Pomodoro Focus effect -->
        <div class="absolute -right-10 -top-10 w-32 h-32 bg-red-500/5 dark:bg-red-400/5 rounded-full blur-xl pointer-events-none"></div>

        <div class="flex items-start justify-between mb-4 relative z-10">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-full bg-red-100 dark:bg-red-900/30 flex items-center justify-center text-red-600 dark:text-red-400">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            </div>
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Focus Time</h2>
          </div>
          <!-- Phase Badge Display -->
          <div v-if="analyticsStore.analytics.pomodoros.current_phase" class="px-3 py-1 bg-red-50 dark:bg-red-900/40 text-red-700 dark:text-red-300 text-xs font-semibold rounded-full border border-red-100 dark:border-red-800/50">
            {{ analyticsStore.analytics.pomodoros.current_phase }} Phase
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4 relative z-10">
          <div>
            <p class="text-sm text-gray-500 dark:text-gray-400">Total Minutes</p>
            <p class="text-3xl font-bold text-gray-900 dark:text-white">{{ analyticsStore.analytics.pomodoros.total_focus_minutes }}</p>
          </div>
          <div class="flex flex-col items-end">
             <p class="text-sm text-gray-500 dark:text-gray-400 mr-2 text-right w-full">Current Level</p>
             <div class="flex items-baseline gap-1">
               <span class="text-3xl font-bold text-red-600 dark:text-red-400">LVL {{ analyticsStore.analytics.pomodoros.current_level }}</span>
             </div>
          </div>
          <div class="col-span-2 pt-2 border-t border-gray-100 dark:border-gray-700 mt-2">
            <p class="text-sm text-gray-500 dark:text-gray-400">Sessions Completed: <span class="font-medium text-gray-900 dark:text-white">{{ analyticsStore.analytics.pomodoros.total_sessions_completed }}</span></p>
          </div>
        </div>
      </div>

      <!-- Tasks Card -->
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6 transition-all hover:shadow-md">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 rounded-full bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center text-blue-600 dark:text-blue-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"></path></svg>
          </div>
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Tasks Overview</h2>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <p class="text-sm text-gray-500 dark:text-gray-400">Completed</p>
            <p class="text-3xl font-bold text-green-600 dark:text-green-400">{{ analyticsStore.analytics.tasks.total_completed }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500 dark:text-gray-400">Pending</p>
            <p class="text-3xl font-bold text-orange-500 dark:text-orange-400">{{ analyticsStore.analytics.tasks.total_pending }}</p>
          </div>
        </div>
      </div>

      <!-- Habits Card -->
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6 transition-all hover:shadow-md">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 rounded-full bg-emerald-100 dark:bg-emerald-900/30 flex items-center justify-center text-emerald-600 dark:text-emerald-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
          </div>
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Habit Consistency</h2>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <p class="text-sm text-gray-500 dark:text-gray-400">Max Streak</p>
            <p class="text-3xl font-bold text-gray-900 dark:text-white">{{ analyticsStore.analytics.habits.highest_streak }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500 dark:text-gray-400">Total Logs</p>
            <p class="text-3xl font-bold text-gray-900 dark:text-white">{{ analyticsStore.analytics.habits.total_logs }}</p>
          </div>
          <div class="col-span-2">
            <p class="text-sm text-gray-500 dark:text-gray-400">Active Habits Tracking: <span class="font-medium text-gray-900 dark:text-white">{{ analyticsStore.analytics.habits.active_habits }}</span></p>
          </div>
        </div>
      </div>

      <!-- Goals Card -->
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6 transition-all hover:shadow-md">
        <div class="flex items-center gap-3 mb-4">
          <div class="w-10 h-10 rounded-full bg-purple-100 dark:bg-purple-900/30 flex items-center justify-center text-purple-600 dark:text-purple-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
          </div>
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Goals Achievement</h2>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <p class="text-sm text-gray-500 dark:text-gray-400">Completed</p>
            <p class="text-3xl font-bold text-green-600 dark:text-green-400">{{ analyticsStore.analytics.goals.completed_goals }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500 dark:text-gray-400">Active</p>
            <p class="text-3xl font-bold text-purple-600 dark:text-purple-400">{{ analyticsStore.analytics.goals.active_goals }}</p>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useAnalyticsStore } from '../stores/analytics'

const analyticsStore = useAnalyticsStore()

onMounted(() => {
  analyticsStore.fetchAnalytics()
})
</script>
