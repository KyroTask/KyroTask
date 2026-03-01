<template>
  <div class="space-y-6 md:space-y-8 pb-10">
    <!-- Header & Filters -->
    <div class="flex flex-col gap-4">
      <div>
        <h1 class="text-3xl font-black text-slate-900 dark:text-white tracking-tight">
          Activity Log
        </h1>
        <p class="text-slate-500 dark:text-slate-400 mt-1 font-medium">
          Your complete history of actions and focus sessions.
        </p>
      </div>

      <!-- Filter Chips -->
      <div class="flex items-center gap-2 overflow-x-auto pb-2 scrollbar-hide -mx-4 px-4 md:mx-0 md:px-0">
        <button
          v-for="filter in filters"
          :key="filter.value"
          @click="activeFilter = filter.value"
          class="px-4 py-2 rounded-xl text-sm font-bold whitespace-nowrap transition-all duration-300"
          :class="activeFilter === filter.value 
            ? 'bg-indigo-600 text-white shadow-md shadow-indigo-500/20 scale-105' 
            : 'bg-white dark:bg-slate-900 text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-800 hover:bg-slate-50 dark:hover:bg-slate-800'"
        >
          {{ filter.label }}
        </button>
      </div>
    </div>

    <!-- Activity Timeline -->
    <div class="relative before:absolute before:inset-0 before:left-6 md:before:left-8 before:w-1 before:bg-slate-100 dark:before:bg-slate-800/50 before:rounded-full">
      
      <!-- Loading State -->
      <div v-if="loading" class="py-20 flex flex-col items-center justify-center">
        <div class="w-12 h-12 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
        <p class="text-slate-400 font-bold uppercase tracking-widest text-xs animate-pulse">Loading History</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredGroups.length === 0" class="py-16 text-center relative z-10 pl-16">
        <div class="w-20 h-20 bg-slate-100 dark:bg-slate-800/50 rounded-3xl flex items-center justify-center mx-auto mb-6 rotate-3">
          <svg class="w-10 h-10 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h3 class="text-xl font-bold text-slate-900 dark:text-white mb-2">
          No Activity Found
        </h3>
        <p class="text-slate-500 dark:text-slate-400 text-sm max-w-xs mx-auto">
          {{ activeFilter === 'all' ? "You haven't logged any actions yet. Start working on your goals!" : "No activities match the current filter." }}
        </p>
      </div>

      <!-- Timeline Content -->
      <transition-group v-else name="list" tag="div" class="space-y-8 relative z-10 w-full overflow-hidden">
        <div v-for="(group, _index) in filteredGroups" :key="group.date" class="relative">
          
          <!-- Date Sticky Header -->
          <div class="sticky top-20 z-20 flex items-center mb-5 -ml-2">
            <div class="bg-slate-50 dark:bg-slate-900 p-1 md:bg-transparent md:dark:bg-transparent md:p-0">
              <span class="bg-white dark:bg-slate-800 border-2 border-slate-100 dark:border-slate-700 text-slate-800 dark:text-slate-200 text-[10px] md:text-xs font-black px-4 py-1.5 rounded-full uppercase tracking-wider shadow-sm">
                {{ group.date }}
              </span>
            </div>
          </div>

          <!-- Items vertically stacked -->
          <div class="flex flex-col gap-4">
            <div v-for="activity in group.items" :key="activity.id"
              class="group relative flex items-start gap-4 md:gap-6 pl-14 md:pl-24 transition-all duration-300 hover:translate-x-1"
            >
              <!-- Connection Line Connector (Desktop logic) -->
              <div class="absolute left-6 md:left-8 top-6 h-full w-1 bg-transparent group-last:hidden"></div>
              
              <!-- Timeline Dot -->
              <div class="absolute left-4 md:left-6 top-5 w-5 h-5 rounded-full border-4 border-slate-50 dark:border-slate-950 shadow-sm z-10 transition-transform duration-300 group-hover:scale-125" :class="getMeta(activity).dotClass"></div>

              <!-- Activity Card -->
              <div class="flex-1 bg-white dark:bg-slate-900 p-4 md:p-5 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm hover:shadow-md transition-all">
                <div class="flex sm:items-center flex-col sm:flex-row gap-3 sm:gap-4">
                  
                  <!-- Icon Box -->
                  <div class="w-10 h-10 md:w-12 md:h-12 rounded-xl flex items-center justify-center shrink-0 border mt-1 sm:mt-0" :class="getMeta(activity).iconBoxClass">
                    <component :is="getMeta(activity).icon" class="w-5 h-5 md:w-6 md:h-6" />
                  </div>

                  <!-- Details -->
                  <div class="flex-1 min-w-0">
                    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-1 mb-1">
                      <h3 class="text-sm md:text-base font-bold text-slate-900 dark:text-white truncate">
                        {{ formatActionTitle(activity) }}
                      </h3>
                      <span class="text-[10px] md:text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-widest shrink-0">
                        {{ formatTime(activity.created_at) }}
                      </span>
                    </div>
                    <p class="text-xs md:text-sm text-slate-500 dark:text-slate-400 font-medium">
                      {{ activity.description }}
                    </p>
                  </div>

                </div>
              </div>
            </div>
          </div>
        </div>
      </transition-group>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, defineComponent, h } from "vue";
import api from "../services/api";

const activities = ref([]);
const loading = ref(true);
const activeFilter = ref("all");

const filters = [
  { label: "All Activity", value: "all" },
  { label: "Focus Timer", value: "pomodoro" },
  { label: "Tasks", value: "task" },
  { label: "Habits", value: "habit" },
  { label: "Projects", value: "project" },
];

const filteredGroups = computed(() => {
  let filtered = activities.value;
  
  if (activeFilter.value !== "all") {
    filtered = activities.value.filter(a => a.resource_type.toLowerCase() === activeFilter.value);
  }

  const groupsObj = {};
  filtered.forEach((activity) => {
    // Determine friendly date string
    const d = new Date(activity.created_at);
    const today = new Date();
    const yesterday = new Date(today);
    yesterday.setDate(yesterday.getDate() - 1);

    let dateStr = d.toLocaleDateString("en-US", { month: "short", day: "numeric", year: "numeric" });
    if (d.toDateString() === today.toDateString()) {
      dateStr = "Today";
    } else if (d.toDateString() === yesterday.toDateString()) {
      dateStr = "Yesterday";
    }

    if (!groupsObj[dateStr]) groupsObj[dateStr] = [];
    groupsObj[dateStr].push(activity);
  });

  // Convert to array so transition-group handles keys well
  return Object.keys(groupsObj).map(date => ({
    date,
    items: groupsObj[date]
  }));
});

async function fetchActivities() {
  loading.value = true;
  try {
    const response = await api.get("/activities");
    activities.value = response.data;
  } catch (err) {
    console.error(err);
  } finally {
    loading.value = false;
  }
}

function formatTime(dateStr) {
  return new Date(dateStr).toLocaleTimeString("en-US", {
    hour: "2-digit",
    minute: "2-digit",
  });
}

function formatActionTitle(activity) {
  const t = activity.resource_type;
  if (t === 'Pomodoro') {
    if (activity.action === 'Started') return 'Started Focus Session';
    if (activity.action === 'Completed') return 'Completed Focus Level';
    if (activity.action === 'Abandoned') return 'Abandoned Session';
  }
  return `${activity.action} ${t}`;
}

// ---- ICONS ----
const IconPomodoro = defineComponent({ render: () => h("svg", { fill: "none", stroke: "currentColor", viewBox: "0 0 24 24" }, [h("path", { "stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", d: "M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" })])});
const IconTask = defineComponent({ render: () => h("svg", { fill: "none", stroke: "currentColor", viewBox: "0 0 24 24" }, [h("path", { "stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", d: "M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" })])});
const IconHabit = defineComponent({ render: () => h("svg", { fill: "none", stroke: "currentColor", viewBox: "0 0 24 24" }, [h("path", { "stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", d: "M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" })])});
const IconProject = defineComponent({ render: () => h("svg", { fill: "none", stroke: "currentColor", viewBox: "0 0 24 24" }, [h("path", { "stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", d: "M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" })])});
const IconGeneric = defineComponent({ render: () => h("svg", { fill: "none", stroke: "currentColor", viewBox: "0 0 24 24" }, [h("path", { "stroke-linecap": "round", "stroke-linejoin": "round", "stroke-width": "2", d: "M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" })])});

// ---- METADATA (Colors & Icons) ----
function getMeta(activity) {
  const type = activity.resource_type.toLowerCase();
  
  if (type === 'pomodoro') {
    if (activity.action === 'Abandoned') return { icon: IconPomodoro, dotClass: 'bg-rose-500', iconBoxClass: 'bg-rose-50 dark:bg-rose-900/20 text-rose-500 border-rose-100 dark:border-rose-800' };
    return { icon: IconPomodoro, dotClass: 'bg-emerald-500', iconBoxClass: 'bg-emerald-50 dark:bg-emerald-900/20 text-emerald-500 border-emerald-100 dark:border-emerald-800' };
  }
  if (type === 'task') return { icon: IconTask, dotClass: 'bg-blue-500', iconBoxClass: 'bg-blue-50 dark:bg-blue-900/20 text-blue-500 border-blue-100 dark:border-blue-800' };
  if (type === 'habit') return { icon: IconHabit, dotClass: 'bg-indigo-500', iconBoxClass: 'bg-indigo-50 dark:bg-indigo-900/20 text-indigo-500 border-indigo-100 dark:border-indigo-800' };
  if (type === 'project') return { icon: IconProject, dotClass: 'bg-amber-500', iconBoxClass: 'bg-amber-50 dark:bg-amber-900/20 text-amber-500 border-amber-100 dark:border-amber-800' };
  
  return { icon: IconGeneric, dotClass: 'bg-slate-400', iconBoxClass: 'bg-slate-50 dark:bg-slate-800 text-slate-500 border-slate-200 dark:border-slate-700' };
}

onMounted(fetchActivities);
</script>

<style scoped>
.list-move,
.list-enter-active,
.list-leave-active {
  transition: all 0.5s ease;
}
.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}
.list-leave-active {
  position: absolute;
}
</style>
