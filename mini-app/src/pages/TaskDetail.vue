<template>
  <div class="space-y-8">
    <div v-if="loading" class="py-20 text-center">
      <div
        class="animate-spin rounded-full h-12 w-12 border-4 border-slate-100 border-t-indigo-600 mx-auto"
      ></div>
    </div>

    <div v-else-if="task" class="max-w-3xl mx-auto space-y-8">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <button
          @click="$router.back()"
          class="flex items-center space-x-2 text-slate-500 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200 transition-colors font-bold"
        >
          <svg
            class="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M10 19l-7-7m0 0l7-7m-7 7h18"
            />
          </svg>
          <span>Back</span>
        </button>
        <div class="flex items-center space-x-3">
          <button
            @click="openEditModal"
            class="p-2 text-slate-400 dark:text-slate-500 hover:text-indigo-600 hover:bg-indigo-50 dark:hover:bg-indigo-900/20 rounded-xl transition-all"
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
                d="M11 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-5M16.5 3.5a2.121 2.121 0 013 3L7 19l-4 1 1-4L16.5 3.5z"
              />
            </svg>
          </button>
          <button
            @click="handleDelete"
            class="p-2 text-slate-400 dark:text-slate-500 hover:text-rose-600 hover:bg-rose-50 dark:hover:bg-rose-900/20 rounded-xl transition-all"
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
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
              />
            </svg>
          </button>
        </div>
      </div>

      <!-- Task Card -->
      <div
        class="bg-white dark:bg-slate-900 rounded-[2.5rem] border border-slate-100 dark:border-slate-800 shadow-xl shadow-slate-200/50 dark:shadow-none overflow-hidden"
      >
        <div class="p-8 md:p-12 space-y-8">
          <div class="flex items-start justify-between gap-6">
            <div class="flex-1">
              <div class="flex items-center space-x-3 mb-4">
                <span
                  :class="priorityClass(task.priority)"
                  class="text-[10px] font-black px-3 py-1 rounded-full uppercase tracking-[0.2em]"
                >
                  {{ task.priority }}
                </span>
                <span class="text-slate-300 dark:text-slate-700">•</span>
                <span
                  class="text-xs font-bold text-slate-400 dark:text-slate-500 uppercase tracking-widest"
                  >{{ task.project?.name || "No Project" }}</span
                >
              </div>
              <h1
                class="text-3xl md:text-4xl font-black text-slate-900 dark:text-white leading-tight"
              >
                {{ task.title }}
              </h1>
            </div>
            <button
              @click="toggleStatus"
              class="w-12 h-12 rounded-2xl border-2 flex items-center justify-center transition-all flex-shrink-0"
              :class="[
                task.status === 'completed'
                  ? 'bg-indigo-600 border-indigo-600 text-white shadow-lg shadow-indigo-100'
                  : 'border-slate-100 dark:border-slate-800 hover:border-indigo-500',
              ]"
            >
              <svg
                v-if="task.status === 'completed'"
                class="w-6 h-6"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="3"
                  d="M5 13l4 4L19 7"
                />
              </svg>
            </button>
          </div>

          <div class="space-y-4">
            <label
              class="block text-xs font-black text-slate-400 uppercase tracking-[0.2em]"
              >Description</label
            >
            <p
              class="text-slate-600 dark:text-slate-300 leading-relaxed text-lg"
            >
              {{ task.description || "No description provided." }}
            </p>
          </div>

          <!-- Metadata Grid -->
          <div
            class="grid grid-cols-2 md:grid-cols-3 gap-8 pt-8 border-t border-slate-50 dark:border-slate-800"
          >
            <div>
              <label
                class="block text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] mb-2"
                >Created</label
              >
              <p class="text-sm font-bold text-slate-800 dark:text-slate-200">
                {{ formatDate(task.created_at) }}
              </p>
            </div>
            <div>
              <label
                class="block text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] mb-2"
                >Status</label
              >
              <p
                class="text-sm font-bold text-slate-800 dark:text-slate-200 capitalize"
              >
                {{ task.status }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Edit Task Modal -->
    <div
      v-if="showEditModal"
      class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    >
      <div
        class="bg-white dark:bg-slate-900 rounded-3xl w-full max-w-lg shadow-2xl overflow-hidden border border-slate-100 dark:border-slate-800"
      >
        <div
          class="p-6 border-b border-slate-50 dark:border-slate-800 flex justify-between items-center"
        >
          <h3 class="text-xl font-bold text-slate-900 dark:text-white">
            Edit Task
          </h3>
          <button
            @click="showEditModal = false"
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
              v-model="editTaskData.title"
              type="text"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
            />
          </div>
          <div>
            <label
              class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2"
              >Project</label
            >
            <select
              v-model="editTaskData.project_id"
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
              v-model="editTaskData.goal_id"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
            >
              <option :value="null">No Goal</option>
              <option v-for="g in goalStore.goals" :key="g.id" :value="g.id">
                {{ g.title }}
              </option>
            </select>
          </div>
          <div>
            <label
              class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2"
              >Priority</label
            >
            <select
              v-model="editTaskData.priority"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
            >
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
            </select>
          </div>
          <div>
            <label
              class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2"
              >Description</label
            >
            <textarea
              v-model="editTaskData.description"
              rows="3"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
            ></textarea>
          </div>
        </div>
        <div
          class="p-6 bg-slate-50 dark:bg-slate-800/50 flex justify-end space-x-3"
        >
          <button
            @click="showEditModal = false"
            class="px-6 py-2 text-slate-600 dark:text-slate-400 font-bold"
          >
            Cancel
          </button>
          <button
            @click="handleUpdateTask"
            class="px-6 py-2 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100"
          >
            Save Changes
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, inject } from "vue";
import { useRoute, useRouter } from "vue-router";
import api from "../services/api";
import { useTaskStore } from "../stores/tasks";
import { useProjectStore } from "../stores/projects";
import { useGoalStore } from "../stores/goals";

const route = useRoute();
const router = useRouter();
const taskStore = useTaskStore();
const projectStore = useProjectStore();
const goalStore = useGoalStore();
const showToast = inject("showToast");

const task = ref(null);
const loading = ref(true);
const showEditModal = ref(false);
const editTaskData = ref({});

async function fetchTask() {
  loading.value = true;
  try {
    const response = await api.get(`/tasks/${route.params.slug}`);
    task.value = response.data;
    // Fetch related data for edit modal
    projectStore.fetchProjects();
    goalStore.fetchGoals();
  } catch (err) {
    console.error(err);
  } finally {
    loading.value = false;
  }
}

function openEditModal() {
  editTaskData.value = { ...task.value };
  showEditModal.value = true;
}

async function handleUpdateTask() {
  try {
    const updated = await taskStore.updateTask(
      task.value.id,
      editTaskData.value
    );
    task.value = updated;
    showEditModal.value = false;
    showToast("Task updated successfully", "success");
  } catch (err) {
    showToast("Failed to update task", "error");
    console.error(err);
  }
}

async function toggleStatus() {
  try {
    const updated = await taskStore.toggleTaskStatus(task.value);
    task.value = updated;
  } catch (err) {
    console.error(err);
  }
}

async function handleDelete() {
  if (confirm("Are you sure you want to delete this task?")) {
    try {
      await taskStore.deleteTask(task.value.id);
      router.push("/tasks");
    } catch (err) {
      console.error(err);
    }
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

const formatDate = (dateStr) => {
  if (!dateStr) return "N/A";
  return new Date(dateStr).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
};

onMounted(fetchTask);
</script>
