<template>
  <div class="space-y-8">
    <!-- Page Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-4">
      <div>
        <h1 class="text-2xl md:text-3xl font-black tracking-tight text-slate-900 dark:text-white">My Tasks</h1>
        <p class="text-slate-500 dark:text-slate-400 text-sm font-medium mt-1">Manage and track your daily activities.</p>
      </div>
      <div class="flex items-center gap-3">
        <div class="relative">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search tasks..."
            class="pl-10 pr-4 py-2 bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl focus:ring-2 focus:ring-indigo-500 outline-none transition-all w-full sm:w-56 text-slate-900 dark:text-white placeholder-slate-400 text-sm"
          />
          <svg class="w-4 h-4 text-slate-400 absolute left-3.5 top-2.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <button
          @click="showNewTaskModal = true"
          class="bg-indigo-600 text-white px-4 py-2 rounded-xl font-bold hover:bg-indigo-700 transition-all flex items-center space-x-2 shadow-lg shadow-indigo-200 dark:shadow-none group text-sm"
        >
          <svg class="w-4 h-4 group-hover:rotate-90 transition-transform duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4" />
          </svg>
          <span>Add Task</span>
        </button>
      </div>
    </div>

    <!-- Stats Bar -->
    <div class="flex items-center space-x-3 overflow-x-auto pb-2 no-scrollbar mb-4">
      <div class="bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 px-4 py-2 rounded-xl flex items-center space-x-2 whitespace-nowrap shadow-sm">
        <span class="text-lg font-black text-slate-900 dark:text-white">{{ taskStore.tasks.length }}</span>
        <span class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Total</span>
      </div>
      <div class="bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 px-4 py-2 rounded-xl flex items-center space-x-2 whitespace-nowrap shadow-sm">
        <span class="text-lg font-black text-emerald-600">{{ taskStore.tasks.filter(t => t.status === 'completed').length }}</span>
        <span class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Done</span>
      </div>
      <div class="bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 px-4 py-2 rounded-xl flex items-center space-x-2 whitespace-nowrap shadow-sm">
        <span class="text-lg font-black text-rose-600">{{ taskStore.tasks.filter(t => t.status === 'pending').length }}</span>
        <span class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Pending</span>
      </div>
    </div>

    <!-- Filters -->
    <div
      class="flex items-center space-x-2 overflow-x-auto pb-2 scrollbar-hide"
    >
      <button
        v-for="filter in filters"
        :key="filter"
        @click="activeFilter = filter"
        class="px-4 py-2 rounded-full text-sm font-semibold whitespace-nowrap transition-all"
        :class="[
          activeFilter === filter
            ? 'bg-slate-900 dark:bg-indigo-600 text-white shadow-md'
            : 'bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-400 border border-slate-200 dark:border-slate-700 hover:border-slate-300 dark:hover:border-slate-600',
        ]"
      >
        {{ filter }}
      </button>
    </div>

    <!-- Task List (Focus Mode) -->
    <div class="space-y-10">
      <div v-if="taskStore.loading" class="py-20 text-center">
        <div
          class="animate-spin rounded-full h-12 w-12 border-4 border-slate-100 border-t-indigo-600 mx-auto"
        ></div>
      </div>
      
      <div v-else v-for="(tasks, priority) in groupedTasks" :key="priority" class="space-y-4">
        <div v-if="tasks.length > 0" class="flex items-center space-x-4 px-2">
          <div :class="priorityIndicatorClass(priority)" class="w-1.5 h-6 rounded-full"></div>
          <h2 class="text-sm font-black uppercase tracking-[0.2em] text-slate-400 dark:text-slate-500">
            {{ priority }} Priority
          </h2>
          <div class="flex-1 h-px bg-slate-100 dark:bg-slate-800"></div>
          <span class="text-xs font-bold text-slate-400">{{ tasks.length }}</span>
        </div>

        <div class="grid grid-cols-1 gap-3">
          <div
            v-for="task in tasks"
            :key="task.isHabit ? 'habit-' + task.id : task.id"
            class="bg-white dark:bg-slate-900 p-4 md:p-4 rounded-2xl border border-slate-100 dark:border-slate-800 shadow-sm hover:shadow-md transition-all duration-300 flex items-center space-x-4 group cursor-pointer"
            @click="task.isHabit ? navigateToHabit(task) : $router.push(`/tasks/${task.slug}`)"
          >
            <button
              @click.stop="task.isHabit ? handleToggleHabit(task) : handleToggleTask(task)"
              class="w-6 h-6 rounded-lg border-2 flex items-center justify-center transition-all"
              :class="[
                task.status === 'completed'
                  ? (task.isHabit ? 'bg-rose-500 border-rose-500' : 'bg-indigo-600 border-indigo-600')
                  : 'border-slate-200 dark:border-slate-700 hover:border-indigo-500',
              ]"
            >
              <svg
                v-if="task.status === 'completed'"
                class="w-4 h-4 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
              </svg>
            </button>

            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <h3
                  class="text-base md:text-base font-bold transition-all truncate"
                  :class="[
                    task.status === 'completed'
                      ? 'text-slate-400 dark:text-slate-500 line-through'
                      : 'text-slate-800 dark:text-slate-200',
                  ]"
                >
                  {{ task.title }}
                </h3>
                <span v-if="task.isHabit" class="px-2 py-0.5 rounded-md bg-rose-50 dark:bg-rose-900/20 text-[8px] font-black text-rose-600 dark:text-rose-400 uppercase tracking-widest">Habit</span>
              </div>
            </div>

            <div v-if="!task.isHabit" class="flex items-center space-x-3 opacity-0 group-hover:opacity-100 transition-opacity">
               <button
                @click.stop="handleDeleteTask(task.id)"
                class="p-2 text-slate-300 hover:text-rose-500 transition-colors"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div
      v-if="!taskStore.loading && filteredTasks.length === 0"
      class="py-20 text-center"
    >
      <div
        class="w-24 h-24 bg-slate-100 rounded-full flex items-center justify-center mx-auto mb-6"
      >
        <svg
          class="w-12 h-12 text-slate-300"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"
          />
        </svg>
      </div>
      <h3 class="text-xl font-bold text-slate-900 dark:text-white">
        No tasks found
      </h3>
      <p class="text-slate-500 mt-2">
        Try adjusting your filters or create a new task.
      </p>
    </div>

    <!-- New Task Modal -->
    <div
      v-if="showNewTaskModal"
      class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    >
      <div
        class="bg-white dark:bg-slate-900 rounded-3xl w-full max-w-lg shadow-2xl overflow-hidden border border-slate-100 dark:border-slate-800 flex flex-col max-h-[90vh]"
      >
        <div
          class="p-4 md:p-6 border-b border-slate-50 dark:border-slate-800 flex justify-between items-center shrink-0"
        >
          <h3 class="text-lg md:text-xl font-bold text-slate-900 dark:text-white">
            Create New Task
          </h3>
          <button
            @click="showNewTaskModal = false"
            class="text-slate-400 hover:text-slate-600"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="p-4 md:p-6 space-y-4 overflow-y-auto">
          <div>
            <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Task Title</label>
            <input
              v-model="newTask.title"
              type="text"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
              placeholder="What needs to be done?"
            />
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Project</label>
            <select
               v-model="newTask.project_id"
               class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
             >
               <option :value="null">No Project</option>
               <option v-for="p in projectStore.projects" :key="p.id" :value="p.id">{{ p.name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Goal</label>
            <select
               v-model="newTask.goal_id"
               class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
             >
               <option :value="null">No Goal</option>
               <option v-for="g in goalStore.goals" :key="g.id" :value="g.id">{{ g.title }}</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Priority</label>
            <select
              v-model="newTask.priority"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
            >
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
            </select>
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Description</label>
            <textarea
              v-model="newTask.description"
              rows="2"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
              placeholder="Add some details..."
            ></textarea>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Start Time</label>
              <input
                v-model="newTask.start_time"
                type="datetime-local"
                class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-lg px-2 py-2 md:rounded-xl md:px-4 md:py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white text-xs md:text-base"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">End Time</label>
              <input
                v-model="newTask.end_time"
                type="datetime-local"
                class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-lg px-2 py-2 md:rounded-xl md:px-4 md:py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white text-xs md:text-base"
              />
            </div>
          </div>
          <div class="flex items-center space-x-3">
             <input
               v-model="newTask.is_daily"
               type="checkbox"
               id="is_daily"
               class="w-5 h-5 text-indigo-600 border-slate-200 rounded focus:ring-indigo-500"
             />
             <label for="is_daily" class="text-sm font-bold text-slate-700 dark:text-slate-300">Daily Task</label>
          </div>
        </div>
        <div class="p-4 md:p-6 bg-slate-50 dark:bg-slate-800/50 flex justify-end space-x-3 shrink-0">
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
import { ref, computed, onMounted, inject } from "vue";
import { useRouter } from "vue-router";
import { useTaskStore } from "../stores/tasks";
import { useProjectStore } from "../stores/projects";
import { useGoalStore } from "../stores/goals";
import { useHabitStore } from "../stores/habits";

const taskStore = useTaskStore();
const projectStore = useProjectStore();
const goalStore = useGoalStore();
const habitStore = useHabitStore();
const router = useRouter();
const showToast = inject("showToast");

const activeFilter = ref("All");
const filters = ["All", "Pending", "Completed", "High Priority"];
const searchQuery = ref("");
const showNewTaskModal = ref(false);
const newTask = ref({
  title: "",
  project_id: null,
  goal_id: null,
  priority: "medium",
  description: "",
  start_time: null,
  end_time: null,
  is_daily: false,
});

const filteredTasks = computed(() => {
  let result = taskStore.tasks;

  if (activeFilter.value === "Pending")
    result = result.filter((t) => t.status === "pending");
  if (activeFilter.value === "Completed")
    result = result.filter((t) => t.status === "completed");
  if (activeFilter.value === "High Priority")
    result = result.filter((t) => t.priority === "high");

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    result = result.filter((t) => t.title.toLowerCase().includes(query));
  }

  return result;
});

const scheduledHabitsForToday = computed(() => {
  const today = new Date();
  const dayOfWeek = today.getDay().toString(); // 0 = Sunday
  const todayStr = today.toISOString().split('T')[0];

  return habitStore.habits.filter(habit => {
    if (!habit.is_active) return false;
    
    if (habit.frequency === 'daily') return true;
    
    if (habit.frequency === 'weekly' && habit.scheduled_days) {
      return habit.scheduled_days.split(',').includes(dayOfWeek);
    }
    
    return false;
  }).map(habit => ({
    id: habit.id,
    title: habit.name,
    status: habit.logs?.some(log => log.log_date.startsWith(todayStr)) ? 'completed' : 'pending',
    priority: 'medium', // Habits are medium priority by default
    isHabit: true,
    habitData: habit
  }));
});

const groupedTasks = computed(() => {
  const tasks = [...filteredTasks.value, ...scheduledHabitsForToday.value];
  return {
    High: tasks.filter((t) => t.priority === "high"),
    Medium: tasks.filter((t) => t.priority === "medium"),
    Low: tasks.filter((t) => t.priority === "low"),
  };
});

const priorityIndicatorClass = (priority) => {
  switch (priority.toLowerCase()) {
    case "high":
      return "bg-rose-500";
    case "medium":
      return "bg-amber-500";
    case "low":
      return "bg-emerald-500";
    default:
      return "bg-slate-500";
  }
};

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
      description: "",
      start_time: null,
      end_time: null,
      is_daily: false,
    };
    showToast("Task created successfully", "success");
  } catch (err) {
    showToast("Failed to create task", "error");
    console.error(err);
  }
}

async function handleDeleteTask(id) {
  if (confirm("Are you sure you want to delete this task?")) {
    try {
      await taskStore.deleteTask(id);
      showToast("Task deleted", "success");
    } catch (err) {
      showToast("Failed to delete task", "error");
    }
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

async function handleToggleHabit(task) {
  const today = new Date().toISOString().split('T')[0];
  try {
    await habitStore.logHabit(task.id, today);
    showToast("Habit logged!", "success");
  } catch (err) {
    showToast("Failed to log habit", "error");
  }
}

function navigateToHabit(task) {
  router.push({
    path: '/habits',
    query: {
      id: task.id,
      expand: task.habitData.frequency === 'weekly' ? 'true' : 'false'
    }
  });
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
  habitStore.fetchHabits();
});
</script>

<style scoped>
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
