<template>
  <RouterView />
  <Toast ref="toastRef" />
</template>

<script setup>
import { onMounted, ref, provide } from "vue";
import { useAuthStore } from "./stores/auth";
import { useTaskStore } from "./stores/tasks";
import { useProjectStore } from "./stores/projects";
import { useGoalStore } from "./stores/goals";
import { useHabitStore } from "./stores/habits";
import Toast from "./components/Toast.vue";

const authStore = useAuthStore();
const toastRef = ref(null);

provide("showToast", (message, type, duration) => {
  toastRef.value?.addToast(message, type, duration);
});

onMounted(async () => {
  authStore.applyTheme();

  if (authStore.isAuthenticated) {
    try {
      await authStore.fetchCurrentUser();
      
      // Pre-fetch all data in background to make navigation instant
      const taskStore = useTaskStore();
      const projectStore = useProjectStore();
      const goalStore = useGoalStore();
      const habitStore = useHabitStore();
      
      Promise.all([
        taskStore.fetchTasks(),
        projectStore.fetchProjects(),
        goalStore.fetchGoals(),
        habitStore.fetchHabits()
      ]).catch(e => console.error("Failed to pre-fetch data:", e));
      
    } catch (e) {
      console.error("Failed to fetch user session:", e);
      authStore.logout();
    }
  }

  if (window.Telegram?.WebApp) {
    window.Telegram.WebApp.ready();
    window.Telegram.WebApp.expand();
  }
});
</script>

<style>
/* Global styles are in style.css */
</style>
