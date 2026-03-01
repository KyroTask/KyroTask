<template>
  <div class="space-y-8 animate-fade-in pb-20 md:pb-6">
    <!-- Page Header (Premium Glow Effect) -->
    <div class="relative mb-8">
      <!-- Decorative Background Glow -->
      <div class="absolute -top-8 -left-8 w-64 h-64 bg-indigo-500/10 dark:bg-indigo-500/5 rounded-full blur-3xl -z-10"></div>
      
      <div class="flex flex-col md:flex-row md:items-end justify-between gap-6">
        <div class="max-w-xl">
          <h1 class="text-3xl md:text-4xl font-black tracking-tight text-gray-900 dark:text-white mb-2 leading-tight">
            Welcome back, <span class="text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-indigo-600 dark:from-blue-400 dark:to-indigo-400">{{ user?.first_name || "User" }}</span>! 👋
          </h1>
          <p class="text-gray-500 dark:text-gray-400 font-medium">Here's your productivity overview for today.</p>
        </div>
        
        <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3 w-full md:w-auto">
          <RouterLink
            to="/calendar"
            class="px-5 py-3 sm:py-2.5 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-300 border border-gray-200 dark:border-gray-700 rounded-xl font-bold hover:bg-gray-50 dark:hover:bg-gray-750 hover:border-gray-300 dark:hover:border-gray-600 transition-all flex items-center justify-center gap-2 text-sm shadow-sm"
          >
            <svg class="w-4 h-4 text-gray-500 dark:text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            <span class="hidden sm:inline">Calendar</span>
          </RouterLink>
          
          <button
            @click="showNewTaskModal = true"
            class="px-5 py-3 sm:py-2.5 bg-gradient-to-r from-blue-600 to-indigo-600 hover:from-blue-500 hover:to-indigo-500 text-white rounded-xl font-bold transition-all flex items-center justify-center gap-2 shadow-lg shadow-indigo-500/30 hover:shadow-indigo-500/40 active:scale-[0.98] text-sm group"
          >
            <svg class="w-5 h-5 group-hover:rotate-90 transition-transform duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4" />
            </svg>
            <span>New Task</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Stats Overview Cards -->
    <div class="grid grid-cols-2 md:grid-cols-4 gap-4 md:gap-6 mb-8">
      
      <!-- Pomodoro Card -->
      <div class="col-span-2 md:col-span-1 bg-gradient-to-br from-gray-900 to-gray-800 dark:from-gray-800 dark:to-gray-900 rounded-2xl p-6 shadow-xl relative overflow-hidden group">
        <div class="absolute top-0 right-0 w-32 h-32 bg-red-500/20 rounded-full blur-2xl -mr-10 -mt-10 transition-transform group-hover:scale-110"></div>
        <div class="relative z-10 flex flex-col h-full justify-between">
          <RouterLink to="/analytics" class="flex justify-between items-start group-hover:opacity-80 transition-opacity">
            <div>
              <p class="text-[10px] font-black text-gray-400 uppercase tracking-widest mb-1 flex items-center gap-1">
                 <svg class="w-3 h-3 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                 Pomodoro
              </p>
              <div class="text-4xl font-black text-white tracking-tight" v-if="analyticsStore.analytics">
                LVL {{ analyticsStore.analytics.pomodoros.current_level }}
              </div>
              <div v-else class="text-4xl font-black text-white tracking-tight animate-pulse">...</div>
            </div>
            <div v-if="analyticsStore.analytics?.pomodoros?.current_phase" class="px-2 py-1 bg-red-500/20 text-red-400 text-[10px] font-bold rounded-lg border border-red-500/30 uppercase tracking-wider backdrop-blur-sm">
              {{ analyticsStore.analytics.pomodoros.current_phase }}
            </div>
          </RouterLink>
          <div class="mt-4 flex items-center justify-between text-sm">
             <span class="text-gray-400 font-medium">Total Focus</span>
             <span class="text-white font-bold" v-if="analyticsStore.analytics">{{ analyticsStore.analytics.pomodoros.total_focus_minutes }} min</span>
             <span class="text-white font-bold" v-else>...</span>
          </div>
        </div>
      </div>

      <!-- Specific Stats Cards -->
      <div v-for="stat in dashboardStats" :key="stat.label" class="bg-white dark:bg-gray-800 border border-gray-100 dark:border-gray-700 rounded-2xl p-5 shadow-sm hover:shadow-md transition-all group">
        <div class="flex items-center justify-between mb-4">
          <div :class="[stat.bg]" class="w-10 h-10 rounded-xl flex items-center justify-center transition-transform group-hover:scale-110">
            <component :is="stat.icon" class="w-5 h-5" :class="stat.color" />
          </div>
          <span class="text-xs font-bold px-2 py-1 rounded-lg" :class="[stat.bg, stat.color]">{{ stat.percent }}%</span>
        </div>
        <div>
           <div class="text-2xl font-black text-gray-900 dark:text-white tracking-tight">{{ stat.value }}</div>
           <p class="text-[10px] font-bold text-gray-400 uppercase tracking-widest mt-1">{{ stat.label }}</p>
        </div>
      </div>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
      
      <!-- Left Column: Recent Tasks -->
      <div class="lg:col-span-8 space-y-6">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
             <div class="w-8 h-8 rounded-lg bg-indigo-100 dark:bg-indigo-900/30 flex items-center justify-center text-indigo-600 dark:text-indigo-400">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path></svg>
             </div>
             <h2 class="text-xl font-bold text-gray-900 dark:text-white">Recent Tasks</h2>
          </div>
          <RouterLink to="/tasks" class="text-indigo-600 dark:text-indigo-400 text-sm font-bold hover:text-indigo-700 dark:hover:text-indigo-300 flex items-center gap-1 group">
            View All
            <svg class="w-4 h-4 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
          </RouterLink>
        </div>

        <div class="bg-white dark:bg-gray-800 rounded-2xl border border-gray-100 dark:border-gray-700 shadow-sm overflow-hidden">
          
          <!-- Loading State -->
          <div v-if="taskStore.loading" class="p-12 text-center">
            <div class="animate-spin rounded-full h-10 w-10 border-4 border-gray-100 border-t-indigo-600 mx-auto"></div>
          </div>
          
          <!-- Empty State -->
          <div v-else-if="taskStore.tasks.length === 0" class="p-12 md:p-16 text-center">
            <div class="w-20 h-20 bg-gray-50 dark:bg-gray-900 rounded-full flex items-center justify-center mx-auto mb-6 shadow-inner">
              <svg class="w-10 h-10 text-gray-300 dark:text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">You're all caught up!</h3>
            <p class="text-gray-500 dark:text-gray-400 text-sm max-w-xs mx-auto mb-6">Create a new task to keep your momentum going.</p>
            <button @click="showNewTaskModal = true" class="px-6 py-2 bg-gray-900 dark:bg-white text-white dark:text-gray-900 font-bold rounded-xl text-sm shadow-md hover:scale-105 transition-transform active:scale-95">
              Add Task
            </button>
          </div>
          
          <!-- Task List -->
          <div v-else class="divide-y divide-gray-50 dark:divide-gray-700/50">
            <div
              v-for="task in recentTasks"
              :key="task.id"
              class="p-4 md:p-5 hover:bg-gray-50 dark:hover:bg-gray-700/30 transition-all duration-300 flex items-center gap-4 group relative"
            >
               <!-- Active Indicator (Left styling) -->
              <div class="absolute left-0 top-0 bottom-0 w-1 bg-indigo-500 opacity-0 group-hover:opacity-100 transition-opacity rounded-r-full"></div>

              <!-- Checkbox -->
              <button
                @click.stop="handleToggleTask(task)"
                class="flex-shrink-0 w-6 h-6 rounded-full border-2 flex items-center justify-center transition-all group/btn"
                :class="[
                  task.status === 'completed'
                    ? 'bg-indigo-500 border-indigo-500 shadow-md shadow-indigo-500/20'
                    : 'border-gray-300 dark:border-gray-600 hover:border-indigo-500 dark:hover:border-indigo-400',
                ]"
              >
                <svg v-if="task.status === 'completed'" class="w-3.5 h-3.5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"/></svg>
                <div v-else class="w-2.5 h-2.5 rounded-full bg-indigo-500 opacity-0 group-hover/btn:opacity-100 transition-opacity"></div>
              </button>

              <!-- Content -->
              <div class="flex-1 min-w-0 cursor-pointer" @click="$router.push(`/tasks/${task.slug}`)">
                <h4 class="text-sm font-bold truncate transition-colors duration-200" :class="task.status === 'completed' ? 'line-through text-gray-400 dark:text-gray-500' : 'text-gray-900 dark:text-white group-hover:text-indigo-600 dark:group-hover:text-indigo-400'">
                  {{ task.title }}
                </h4>
                <div class="flex items-center gap-2 mt-1">
                   <p class="text-xs text-gray-500 dark:text-gray-400 truncate max-w-[150px] sm:max-w-xs block">
                     {{ task.project?.name || "No Project" }}
                   </p>
                   <span class="text-gray-300 dark:text-gray-600">•</span>
                   <p class="text-xs text-gray-400 font-medium whitespace-nowrap">
                     {{ new Date(task.due_date).toLocaleDateString(undefined, { month: 'short', day: 'numeric' }) }}
                   </p>
                </div>
              </div>

              <!-- Badges/Actions -->
              <div class="flex items-center gap-2" @click="$router.push(`/tasks/${task.slug}`)">
                <span :class="priorityClass(task.priority)" class="text-[10px] font-black px-2.5 py-1 rounded-lg uppercase tracking-widest hidden sm:inline-block">
                  {{ task.priority }}
                </span>
                <button class="w-8 h-8 rounded-lg flex items-center justify-center text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors opacity-0 group-hover:opacity-100 max-sm:opacity-100">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right Column: Projects Overview -->
      <div class="lg:col-span-4 space-y-6">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
             <div class="w-8 h-8 rounded-lg bg-teal-100 dark:bg-teal-900/30 flex items-center justify-center text-teal-600 dark:text-teal-400">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 002-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path></svg>
             </div>
             <h2 class="text-xl font-bold text-gray-900 dark:text-white">Active Projects</h2>
          </div>
          <RouterLink to="/projects" class="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-400 hover:text-gray-900 dark:hover:text-white transition-colors">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"></path></svg>
          </RouterLink>
        </div>

        <div class="space-y-3">
          <!-- Empty Projects -->
          <div v-if="projectStore.projects.length === 0" class="bg-white dark:bg-gray-800 p-6 rounded-2xl border border-dashed border-gray-200 dark:border-gray-700 text-center">
             <p class="text-sm font-medium text-gray-500 dark:text-gray-400 mb-3">No active projects.</p>
             <RouterLink to="/projects" class="text-xs font-bold text-indigo-600 hover:text-indigo-700 uppercase tracking-widest">Create Project</RouterLink>
          </div>

          <!-- Project Cards -->
          <RouterLink
            v-for="project in projectStore.projects.slice(0, 4)"
            :key="project.id"
            :to="`/projects/${project.slug}`"
            class="block bg-white dark:bg-gray-800 p-4 md:p-5 rounded-2xl border border-gray-100 dark:border-gray-700 shadow-sm hover:shadow-md hover:border-gray-200 dark:hover:border-gray-600 transition-all group relative overflow-hidden"
          >
            <!-- Hover Gradient -->
            <div class="absolute inset-0 bg-gradient-to-r from-transparent to-gray-50 dark:to-gray-700/30 opacity-0 group-hover:opacity-100 transition-opacity"></div>
            
            <div class="relative z-10">
              <div class="flex justify-between items-start mb-2">
                <h4 class="text-sm font-bold text-gray-900 dark:text-white group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">
                  {{ project.name }}
                </h4>
                <!-- Tiny Progress visual placeholder -->
                <div class="w-6 h-6 rounded-full border-2 border-gray-100 dark:border-gray-700 flex items-center justify-center">
                   <svg class="w-3 h-3 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
                </div>
              </div>
              <p class="text-xs text-gray-500 dark:text-gray-400 line-clamp-2 leading-relaxed">
                {{ project.description || 'No description provided.' }}
              </p>
            </div>
          </RouterLink>
        </div>
      </div>
    </div>

    <!-- New Task Modal (Simplified) -->
    <div
      v-if="showNewTaskModal"
      class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    >
      <div
        class="bg-white dark:bg-slate-900 rounded-3xl w-full max-w-lg shadow-2xl overflow-hidden border border-slate-100 dark:border-slate-800"
      >
        <div
          class="p-6 border-b border-slate-50 dark:border-slate-800 flex justify-between items-center"
        >
          <h3 class="text-xl font-bold text-slate-900 dark:text-white">
            Create New Task
          </h3>
          <button
            @click="showNewTaskModal = false"
            class="text-slate-400 hover:text-slate-600"
          >
            <svg
              class="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>
        <div class="p-6 space-y-4">
          <div>
            <label
              class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2"
              >Task Title</label
            >
            <input
              v-model="newTask.title"
              type="text"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
              placeholder="What needs to be done?"
            />
          </div>
          <div>
            <label
              class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2"
              >Project</label
            >
            <select
              v-model="newTask.project_id"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
            >
              <option :value="null">No Project</option>
              <option
                v-for="p in projectStore.projects"
                :key="p.id"
                :value="p.id"
              >
                {{ p.name }}
              </option>
            </select>
          </div>
          <div>
            <label
              class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2"
              >Goal</label
            >
            <select
              v-model="newTask.goal_id"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
            >
              <option :value="null">No Goal</option>
              <option v-for="g in goalStore.goals" :key="g.id" :value="g.id">
                {{ g.title }}
              </option>
            </select>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label
                class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2"
                >Priority</label
              >
              <select
                v-model="newTask.priority"
                class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
              >
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
              </select>
            </div>
          </div>
        </div>
        <div
          class="p-6 bg-slate-50 dark:bg-slate-800/50 flex justify-end space-x-3"
        >
          <button
            @click="showNewTaskModal = false"
            class="px-6 py-2 text-slate-600 dark:text-slate-400 font-bold"
          >
            Cancel
          </button>
          <button
            @click="handleCreateTask"
            class="px-6 py-2 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100"
          >
            Create Task
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, defineComponent, h, inject } from "vue";
import { useAuthStore } from "../stores/auth";
import { useTaskStore } from "../stores/tasks";
import { useProjectStore } from "../stores/projects";
import { useGoalStore } from "../stores/goals";
import { useAnalyticsStore } from "../stores/analytics";

const authStore = useAuthStore();
const taskStore = useTaskStore();
const projectStore = useProjectStore();
const goalStore = useGoalStore();
const analyticsStore = useAnalyticsStore();
const showToast = inject("showToast");

const user = computed(() => authStore.user);
const showNewTaskModal = ref(false);
const newTask = ref({
  title: "",
  project_id: null,
  goal_id: null,
  priority: "medium",
});

const recentTasks = computed(() => taskStore.tasks.slice(0, 5));

const IconCheck = defineComponent({
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

const IconClock = defineComponent({
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
          d: "M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z",
        }),
      ]
    ),
});

const dashboardStats = computed(() => {
  const totalTasks = taskStore.tasks.length;
  const completedTasks = taskStore.tasks.filter((t) => t.status === "completed").length;
  const taskPercent = totalTasks > 0 ? Math.round((completedTasks / totalTasks) * 100) : 0;

  const totalProjects = projectStore.projects.length;
  const completedProjects = projectStore.projects.filter(p => {
    const projectTasks = taskStore.tasks.filter(t => t.project_id === p.id);
    return projectTasks.length > 0 && projectTasks.every(t => t.status === 'completed');
  }).length;
  const projectPercent = totalProjects > 0 ? Math.round((completedProjects / totalProjects) * 100) : 0;

  const totalGoals = goalStore.goals.length;
  const completedGoals = goalStore.goals.filter(g => g.progress === 100).length;
  const goalPercent = totalGoals > 0 ? Math.round((completedGoals / totalGoals) * 100) : 0;

  return [
    {
      label: "Tasks",
      value: `${completedTasks}/${totalTasks}`,
      percent: taskPercent,
      icon: IconCheck,
      bg: "bg-emerald-500/10",
      color: "text-emerald-400",
      barColor: "bg-emerald-500",
    },
    {
      label: "Projects",
      value: `${completedProjects}/${totalProjects}`,
      percent: projectPercent,
      icon: IconClock,
      bg: "bg-indigo-500/10",
      color: "text-indigo-400",
      barColor: "bg-indigo-500",
    },
    {
      label: "Goals",
      value: `${completedGoals}/${totalGoals}`,
      percent: goalPercent,
      icon: IconClock,
      bg: "bg-amber-500/10",
      color: "text-amber-400",
      barColor: "bg-amber-500",
    }
  ];
});

const overallCompletion = computed(() => {
  const stats = dashboardStats.value;
  if (stats.length === 0) return 0;
  const totalPercent = stats.reduce((acc, curr) => acc + curr.percent, 0);
  return totalPercent / stats.length;
});

async function handleCreateTask() {
  if (!newTask.value.title) return;
  try {
    await taskStore.createTask(newTask.value);
    showNewTaskModal.value = false;
    newTask.value = {
      title: "",
      project_id: null,
      goal_id: null,
      priority: "medium",
    };
    showToast("Task created successfully", "success");
  } catch (err) {
    showToast("Failed to create task", "error");
    console.error(err);
  }
}

async function handleToggleTask(task) {
  try {
    await taskStore.toggleTaskStatus(task);
    showToast(
      task.status === "completed" ? "Task completed!" : "Task reopened",
      "success"
    );
  } catch (err) {
    showToast("Failed to update task", "error");
  }
}

const priorityClass = (priority) => {
  switch (priority) {
    case "high":
      return "bg-rose-50 text-rose-600";
    case "medium":
      return "bg-amber-50 text-amber-600";
    case "low":
      return "bg-emerald-50 text-emerald-600";
    default:
      return "bg-slate-50 text-slate-600";
  }
};

onMounted(() => {
  taskStore.fetchTasks();
  projectStore.fetchProjects();
  goalStore.fetchGoals();
  analyticsStore.fetchAnalytics();
});
</script>
