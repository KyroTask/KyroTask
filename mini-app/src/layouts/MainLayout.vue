<template>
  <div
    class="min-h-screen bg-slate-50 dark:bg-slate-950 flex flex-col md:flex-row"
  >
    <!-- Desktop Sidebar -->
    <aside
      class="hidden md:flex w-64 bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 flex-col sticky top-0 h-screen"
    >
      <div class="p-6">
        <div class="flex items-center space-x-3 mb-10">
          <div
            class="w-10 h-10 flex items-center justify-center transform -rotate-2"
          >
            <img src="/logo.png" alt="KyroTask Logo" class="w-full h-full object-contain drop-shadow-md" />
          </div>
          <span
            class="text-xl font-bold text-slate-800 dark:text-white tracking-tight"
            >KyroTask</span
          >
        </div>

        <nav class="space-y-1">
          <RouterLink
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            class="flex items-center space-x-3 px-4 py-3 rounded-xl transition-all duration-200 group"
            :class="[
              $route.path === item.path
                ? 'bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 dark:text-indigo-400 font-semibold'
                : 'text-slate-500 dark:text-slate-400 hover:bg-slate-50 dark:hover:bg-slate-800 hover:text-slate-700 dark:hover:text-slate-200',
            ]"
          >
            <component :is="item.icon" class="w-5 h-5" />
            <span>{{ item.name }}</span>
          </RouterLink>
        </nav>
      </div>

      <div class="mt-auto p-6 border-t border-slate-100 dark:border-slate-800">
        <div class="flex items-center space-x-3">
          <div
            class="w-10 h-10 rounded-full bg-slate-200 dark:bg-slate-800 overflow-hidden"
          >
            <img
              v-if="user?.photo_url"
              :src="user.photo_url"
              alt="Profile"
              class="w-full h-full object-cover"
            />
            <div
              v-else
              class="w-full h-full flex items-center justify-center text-slate-500 dark:text-slate-400 font-bold"
            >
              {{ user?.first_name?.charAt(0) || "U" }}
            </div>
          </div>
          <div class="flex-1 min-w-0">
            <p
              class="text-sm font-semibold text-slate-800 dark:text-slate-200 truncate"
            >
              {{ user?.first_name || "User" }}
            </p>
            <p class="text-xs text-slate-500 dark:text-slate-400 truncate">
              @{{ user?.username || "telegram" }}
            </p>
          </div>
        </div>
      </div>
    </aside>


    <!-- Main Content -->
    <main class="flex-1 overflow-y-auto" :class="[$route.name === 'Pomodoro' ? 'pb-16 md:pb-0' : 'pb-24 md:pb-0']">
      <div class="w-full h-full" :class="[$route.name === 'Pomodoro' ? 'p-0' : 'p-4 md:p-8']">
        <RouterView />
      </div>
    </main>

    <!-- Mobile Bottom Navigation -->
    <nav
      class="md:hidden fixed bottom-0 left-0 right-0 bg-white/80 dark:bg-slate-900/80 backdrop-blur-lg border-t border-slate-200 dark:border-slate-800 px-6 py-3 flex justify-between items-center z-20"
    >
      <RouterLink
        v-for="item in mobileNavItems"
        :key="item.path"
        :to="item.path"
        class="flex flex-col items-center space-y-1 transition-colors duration-200"
        :class="[
          $route.path === item.path
            ? 'text-indigo-600 dark:text-indigo-400'
            : 'text-slate-400 dark:text-slate-500',
        ]"
      >
        <component :is="item.icon" class="w-6 h-6" />
        <span class="text-[10px] font-medium">{{ item.name }}</span>
      </RouterLink>
    </nav>
  </div>
</template>

<script setup>
import { computed, defineComponent, h } from "vue";
import { useAuthStore } from "../stores/auth";

const authStore = useAuthStore();
const user = computed(() => authStore.user);

const IconDashboard = defineComponent({
  render: () =>
    h(
      "svg",
      {
        fill: "none",
        stroke: "currentColor",
        viewBox: "0 0 24 24",
        class: "w-inherit h-inherit",
      },
      [
        h("path", {
          "stroke-linecap": "round",
          "stroke-linejoin": "round",
          "stroke-width": "2",
          d: "M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z",
        }),
      ]
    ),
});

const IconTasks = defineComponent({
  render: () =>
    h(
      "svg",
      {
        fill: "none",
        stroke: "currentColor",
        viewBox: "0 0 24 24",
        class: "w-inherit h-inherit",
      },
      [
        h("path", {
          "stroke-linecap": "round",
          "stroke-linejoin": "round",
          "stroke-width": "2",
          d: "M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4",
        }),
      ]
    ),
});

const IconPomodoro = defineComponent({
  render: () =>
    h(
      "svg",
      {
        fill: "none",
        stroke: "currentColor",
        viewBox: "0 0 24 24",
        class: "w-inherit h-inherit",
      },
      [
        h("circle", { cx: "12", cy: "14", r: "7", "stroke-width": "2", "stroke-linecap": "round", "stroke-linejoin": "round" }),
        h("path", { "stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", d: "M12 10v4l2.5 2.5" }),
        h("path", { "stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", d: "M10 3h4" }),
        h("path", { "stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", d: "M12 3v4" }),
        h("path", { "stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", d: "M19.5 7.5l-1.5 1.5" }),
      ]
    ),
});

const IconProjects = defineComponent({
  render: () =>
    h(
      "svg",
      {
        fill: "none",
        stroke: "currentColor",
        viewBox: "0 0 24 24",
        class: "w-inherit h-inherit",
      },
      [
        h("path", {
          "stroke-linecap": "round",
          "stroke-linejoin": "round",
          "stroke-width": "2",
          d: "M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z",
        }),
      ]
    ),
});

const IconGoals = defineComponent({
  render: () =>
    h(
      "svg",
      {
        fill: "none",
        stroke: "currentColor",
        viewBox: "0 0 24 24",
        class: "w-inherit h-inherit",
      },
      [
        h("path", {
          "stroke-linecap": "round",
          "stroke-linejoin": "round",
          "stroke-width": "2",
          d: "M13 10V3L4 14h7v7l9-11h-7z",
        }),
      ]
    ),
});

const IconHabits = defineComponent({
  render: () =>
    h(
      "svg",
      {
        fill: "none",
        stroke: "currentColor",
        viewBox: "0 0 24 24",
        class: "w-inherit h-inherit",
      },
      [
        h("path", {
          "stroke-linecap": "round",
          "stroke-linejoin": "round",
          "stroke-width": "2",
          d: "M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z",
        }),
      ]
    ),
});

const IconSettings = defineComponent({
  render: () =>
    h(
      "svg",
      {
        fill: "none",
        stroke: "currentColor",
        viewBox: "0 0 24 24",
        class: "w-inherit h-inherit",
      },
      [
        h("path", {
          "stroke-linecap": "round",
          "stroke-linejoin": "round",
          "stroke-width": "2",
          d: "M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z",
        }),
        h("path", {
          "stroke-linecap": "round",
          "stroke-linejoin": "round",
          "stroke-width": "2",
          d: "M15 12a3 3 0 11-6 0 3 3 0 016 0z",
        }),
      ]
    ),
});

const IconCalendar = defineComponent({
  render: () =>
    h(
      "svg",
      {
        fill: "none",
        stroke: "currentColor",
        viewBox: "0 0 24 24",
        class: "w-inherit h-inherit",
      },
      [
        h("path", {
          "stroke-linecap": "round",
          "stroke-linejoin": "round",
          "stroke-width": "2",
          d: "M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z",
        }),
      ]
    ),
});

const IconAnalytics = defineComponent({
  render: () =>
    h(
      "svg",
      {
        fill: "none",
        stroke: "currentColor",
        viewBox: "0 0 24 24",
        class: "w-inherit h-inherit",
      },
      [
        h("path", {
          "stroke-linecap": "round",
          "stroke-linejoin": "round",
          "stroke-width": "2",
          d: "M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z",
        }),
      ]
    ),
});

const navItems = [
  { name: "Dashboard", path: "/", icon: IconDashboard },
  { name: "Tasks", path: "/tasks", icon: IconTasks },
  { name: "Pomodoro", path: "/pomodoro", icon: IconPomodoro },
  { name: "Projects", path: "/projects", icon: IconProjects },
  { name: "Goals", path: "/goals", icon: IconGoals },
  { name: "Habits", path: "/habits", icon: IconHabits },
  { name: "Calendar", path: "/calendar", icon: IconCalendar },
  { name: "Analytics", path: "/analytics", icon: IconAnalytics },
  { name: "Settings", path: "/settings", icon: IconSettings },
];

const mobileNavItems = [
  { name: "Home", path: "/", icon: IconDashboard },
  { name: "Tasks", path: "/tasks", icon: IconTasks },
  { name: "Pomodoro", path: "/pomodoro", icon: IconPomodoro },
  { name: "Goals", path: "/goals", icon: IconGoals },
  { name: "Calendar", path: "/calendar", icon: IconCalendar },
];
</script>

<style scoped>
.router-link-active {
  color: #4f46e5;
}
</style>
