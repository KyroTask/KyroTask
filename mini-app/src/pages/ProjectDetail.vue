<template>
  <div class="space-y-8">
    <div v-if="loading" class="py-20 text-center">
      <div
        class="animate-spin rounded-full h-12 w-12 border-4 border-slate-100 border-t-indigo-600 mx-auto"
      ></div>
    </div>

    <div v-else-if="project" class="space-y-8">
      <!-- Header -->
      <div
        class="flex flex-col md:flex-row md:items-center justify-between gap-6"
      >
        <div class="flex items-center space-x-4">
          <div
            :class="project.color || 'bg-indigo-500'"
            class="w-16 h-16 rounded-3xl flex items-center justify-center text-white shadow-xl shadow-indigo-100 dark:shadow-none"
          >
            <svg
              class="w-8 h-8"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
              />
            </svg>
          </div>
          <div>
            <h1 class="text-3xl font-black text-slate-900 dark:text-white">
              {{ project.name }}
            </h1>
            <p class="text-slate-500 dark:text-slate-400 font-medium">
              {{ project.description }}
            </p>
          </div>
        </div>
        <div class="flex items-center space-x-3">
          <button
            @click="openEditModal"
            class="bg-white dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700 px-4 py-3 rounded-2xl font-bold hover:bg-slate-50 dark:hover:bg-slate-700 transition-all flex items-center space-x-2"
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
                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
              />
            </svg>
            <span>Edit</span>
          </button>
          <button
            @click="showNewGoalModal = true"
            class="bg-indigo-600 text-white px-6 py-3 rounded-2xl font-bold hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100 dark:shadow-none flex items-center space-x-2"
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
                d="M12 4v16m8-8H4"
              />
            </svg>
            <span>Add Goal</span>
          </button>
        </div>
      </div>

      <!-- Project Goals -->
      <div class="space-y-4">
        <h2 class="text-xl font-bold text-slate-900 dark:text-white">
          Project Goals
        </h2>
        <div class="grid grid-cols-1 gap-4">
          <div
            v-for="goal in project.goals"
            :key="goal.id"
            class="bg-white dark:bg-slate-800 p-5 rounded-2xl border border-slate-100 dark:border-slate-700 shadow-sm hover:shadow-md transition-all flex items-center justify-between group cursor-pointer"
          >
            <div class="flex items-center space-x-4">
              <div
                class="w-10 h-10 rounded-full bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 dark:text-indigo-400 flex items-center justify-center"
              >
                <span class="font-bold text-sm">{{ goal.progress }}%</span>
              </div>
              <div>
                <h3 class="font-bold text-slate-800 dark:text-slate-200">
                  {{ goal.title }}
                </h3>
                <p class="text-xs text-slate-500 dark:text-slate-400">
                  {{ goal.description }}
                </p>
              </div>
            </div>
            <span
              class="text-[10px] font-black px-2 py-1 rounded-full uppercase tracking-widest"
              :class="
                goal.status === 'completed'
                  ? 'bg-emerald-50 text-emerald-600'
                  : 'bg-slate-50 text-slate-600'
              "
            >
              {{ goal.status }}
            </span>
          </div>

          <div
            v-if="!project.goals || project.goals.length === 0"
            class="py-12 text-center bg-slate-50 dark:bg-slate-800/50 rounded-3xl border-2 border-dashed border-slate-200 dark:border-slate-700"
          >
            <p class="text-slate-500 dark:text-slate-400 font-bold">
              No goals in this project yet.
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- New Goal Modal -->
    <div
      v-if="showNewGoalModal"
      class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    >
      <div
        class="bg-white rounded-3xl w-full max-w-lg shadow-2xl overflow-hidden"
      >
        <div
          class="p-6 border-b border-slate-50 flex justify-between items-center"
        >
          <h3 class="text-xl font-bold text-slate-900">
            Add Goal to {{ project?.name }}
          </h3>
          <button
            @click="showNewGoalModal = false"
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
              >Goal Title</label
            >
            <input
              v-model="newGoal.title"
              type="text"
              class="w-full bg-slate-50 dark:bg-slate-700 border border-slate-100 dark:border-slate-600 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 dark:text-white placeholder-slate-400 dark:placeholder-slate-300"
              placeholder="What do you want to achieve?"
            />
          </div>
          <div>
            <label
              class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2"
              >Description</label
            >
            <textarea
              v-model="newGoal.description"
              rows="3"
              class="w-full bg-slate-50 dark:bg-slate-700 border border-slate-100 dark:border-slate-600 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 dark:text-white placeholder-slate-400 dark:placeholder-slate-300"
              placeholder="Details about this goal..."
            ></textarea>
          </div>
        </div>
        <div
          class="p-6 bg-slate-50 dark:bg-slate-800/50 flex justify-end space-x-3"
        >
          <button
            @click="showNewGoalModal = false"
            class="px-6 py-2 text-slate-600 dark:text-slate-400 font-bold"
          >
            Cancel
          </button>
          <button
            @click="handleCreateGoal"
            class="px-6 py-2 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100"
          >
            Create Goal
          </button>
        </div>
      </div>
    </div>

    <!-- Edit Project Modal -->
    <div
      v-if="showEditModal"
      class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    >
      <div
        class="bg-white dark:bg-slate-800 rounded-3xl w-full max-w-lg shadow-2xl overflow-hidden"
      >
        <div
          class="p-6 border-b border-slate-50 dark:border-slate-700 flex justify-between items-center"
        >
          <h3 class="text-xl font-bold text-slate-900 dark:text-white">
            Edit Project
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
              >Project Name</label
            >
            <input
              v-model="editProjectData.name"
              type="text"
              class="w-full bg-slate-50 dark:bg-slate-700 border border-slate-100 dark:border-slate-600 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 dark:text-white"
            />
          </div>
          <div>
            <label
              class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2"
              >Description</label
            >
            <textarea
              v-model="editProjectData.description"
              class="w-full bg-slate-50 dark:bg-slate-700 border border-slate-100 dark:border-slate-600 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 dark:text-white"
              rows="3"
            ></textarea>
          </div>
          <div>
            <label
              class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2"
              >Color</label
            >
            <div class="flex flex-wrap gap-2">
              <button
                v-for="color in [
                  'bg-indigo-500',
                  'bg-rose-500',
                  'bg-emerald-500',
                  'bg-amber-500',
                  'bg-sky-500',
                  'bg-violet-500',
                  'bg-pink-500',
                  'bg-slate-500',
                ]"
                :key="color"
                @click="editProjectData.color = color"
                class="w-8 h-8 rounded-full transition-transform hover:scale-110"
                :class="[
                  color,
                  editProjectData.color === color
                    ? 'ring-2 ring-offset-2 ring-slate-900 dark:ring-white'
                    : '',
                ]"
              ></button>
            </div>
          </div>
        </div>
        <div
          class="p-6 bg-slate-50 dark:bg-slate-900/50 flex justify-end space-x-3"
        >
          <button
            @click="showEditModal = false"
            class="px-6 py-2 text-slate-600 dark:text-slate-400 font-bold"
          >
            Cancel
          </button>
          <button
            @click="handleUpdateProject"
            class="px-6 py-2 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100 dark:shadow-none"
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
import { useRoute } from "vue-router";
import api from "../services/api";
import { useGoalStore } from "../stores/goals";
import { useProjectStore } from "../stores/projects";

const route = useRoute();
const goalStore = useGoalStore();
const projectStore = useProjectStore();
const showToast = inject("showToast");

const project = ref(null);
const loading = ref(true);
const showNewGoalModal = ref(false);
const showEditModal = ref(false);
const newGoal = ref({ title: "", description: "" });
const editProjectData = ref({});

async function fetchProject() {
  loading.value = true;
  try {
    // route.params.slug can be ID or Slug. Backend handles both.
    const response = await api.get(`/projects/${route.params.slug}`);
    project.value = response.data;
  } catch (err) {
    console.error(err);
    showToast("Failed to load project", "error");
  } finally {
    loading.value = false;
  }
}

async function handleCreateGoal() {
  if (!newGoal.value.title) return;
  try {
    await goalStore.createGoal({
      ...newGoal.value,
      project_id: project.value.id, // Use the ID from the loaded project
    });
    // Refresh project to show new goal (or push to local list if store doesn't auto-update)
    // Since goalStore.createGoal adds to its own state, we might need to manually add to project.goals
    // OR just re-fetch the project. Re-fetching is safer to get the updated list with ID.
    await fetchProject(); 
    
    showNewGoalModal.value = false;
    newGoal.value = { title: "", description: "" };
    showToast("Goal created successfully", "success");
  } catch (err) {
    console.error(err);
    showToast("Failed to create goal", "error");
  }
}

function openEditModal() {
  editProjectData.value = { ...project.value };
  showEditModal.value = true;
}

async function handleUpdateProject() {
  try {
    const updated = await projectStore.updateProject(
      project.value.id,
      editProjectData.value
    );
    project.value = updated;
    showEditModal.value = false;
    showToast("Project updated successfully", "success");
  } catch (err) {
    showToast("Failed to update project", "error");
    console.error(err);
  }
}

onMounted(fetchProject);
</script>
