<template>
  <div class="space-y-8">
    <!-- Page Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-4">
      <div>
        <h1 class="text-2xl md:text-3xl font-black tracking-tight text-slate-900 dark:text-white">Projects</h1>
        <p class="text-slate-500 dark:text-slate-400 text-sm font-medium mt-1">Organize your work into projects.</p>
      </div>
      <button
        @click="showNewProjectModal = true"
        class="bg-indigo-600 text-white px-4 py-2 rounded-xl font-bold hover:bg-indigo-700 transition-all flex items-center space-x-2 shadow-lg shadow-indigo-200 dark:shadow-none group text-sm"
      >
        <svg class="w-4 h-4 group-hover:rotate-90 transition-transform duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4" />
        </svg>
        <span>New Project</span>
      </button>
    </div>

    <!-- Stats Bar -->
    <div class="flex items-center space-x-3 overflow-x-auto pb-2 no-scrollbar mb-4">
      <div class="bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 px-4 py-2 rounded-xl flex items-center space-x-2 whitespace-nowrap shadow-sm">
        <span class="text-lg font-black text-slate-900 dark:text-white">{{ totalProjects }}</span>
        <span class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Total</span>
      </div>
      <div class="bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 px-4 py-2 rounded-xl flex items-center space-x-2 whitespace-nowrap shadow-sm">
        <span class="text-lg font-black text-emerald-600">{{ activeProjects }}</span>
        <span class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Active</span>
      </div>
      <div class="bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 px-4 py-2 rounded-xl flex items-center space-x-2 whitespace-nowrap shadow-sm">
        <span class="text-lg font-black text-amber-600">{{ completedProjects }}</span>
        <span class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Done</span>
      </div>
    </div>

    <!-- Projects Roadmap List -->
    <div class="space-y-4 md:space-y-6">
      <div v-if="projectStore.loading" class="py-20 text-center">
        <div
          class="animate-spin rounded-full h-12 w-12 border-4 border-slate-100 border-t-indigo-600 mx-auto"
        ></div>
      </div>
      <div
        v-else
        v-for="(project, index) in projectStore.projects"
        :key="project.id"
        class="bg-white dark:bg-slate-800 rounded-[2rem] md:rounded-[2.5rem] border border-slate-100 dark:border-slate-700 shadow-sm hover:shadow-2xl hover:shadow-indigo-100/50 dark:hover:shadow-none transition-all duration-500 p-6 md:p-8 group cursor-pointer relative overflow-hidden flex flex-col md:flex-row md:items-center gap-6 md:gap-10"
        @click="$router.push(`/projects/${project.slug}`)"
      >
        <!-- Roadmap Number -->
        <div class="hidden md:flex absolute left-0 top-0 bottom-0 w-20 items-center justify-center bg-slate-50 dark:bg-slate-900/50 border-r border-slate-100 dark:border-slate-800">
          <span class="text-5xl font-black text-slate-200 dark:text-slate-800 group-hover:text-indigo-500/20 transition-colors">{{ index + 1 }}</span>
        </div>

        <div class="flex-1 md:ml-20 relative z-10">
          <div class="flex items-center space-x-4 mb-3 md:mb-4">
            <div
              :class="project.color || 'bg-indigo-500'"
              class="w-10 h-10 md:w-12 md:h-12 rounded-xl md:rounded-2xl flex items-center justify-center text-white shadow-lg shadow-indigo-200 dark:shadow-none group-hover:scale-110 transition-transform duration-300"
            >
              <svg class="w-5 h-5 md:w-6 md:h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
              </svg>
            </div>
            <div>
              <h3 class="text-lg md:text-2xl font-black text-slate-900 dark:text-white group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">
                {{ project.name }}
              </h3>
              <p class="text-slate-500 dark:text-slate-400 text-[10px] md:text-xs font-bold uppercase tracking-widest">
                Project Roadmap
              </p>
            </div>
          </div>
          <p class="text-slate-500 dark:text-slate-400 text-sm md:text-base mb-4 md:mb-6 line-clamp-2 leading-relaxed max-w-3xl">
            {{ project.description }}
          </p>
        </div>

        <!-- Roadmap Progress Section -->
        <div class="w-full md:w-64 space-y-3 md:space-y-4 relative z-10">
          <div class="flex items-center justify-between">
            <span class="text-[10px] md:text-xs font-black uppercase tracking-widest text-slate-400 dark:text-slate-500">Completion</span>
            <span class="text-indigo-600 dark:text-indigo-400 font-black">{{ calculateProjectProgress(project) }}%</span>
          </div>
          <div class="w-full bg-slate-100 dark:bg-slate-700 h-2.5 md:h-3 rounded-full overflow-hidden relative">
            <div
              class="bg-indigo-600 h-full transition-all duration-1000 shadow-[0_0_15px_rgba(79,70,229,0.5)]"
              :style="{ width: `${calculateProjectProgress(project)}%` }"
            ></div>
          </div>
          <div class="flex justify-between items-center pt-2">
            <button
              @click.stop="confirmDelete(project)"
              class="text-[10px] font-bold uppercase tracking-widest text-slate-400 hover:text-rose-500 transition-colors flex items-center space-x-1"
            >
              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              <span>Delete Project</span>
            </button>
            <div class="flex -space-x-2">
              <div v-for="i in 3" :key="i" class="w-6 h-6 rounded-full border-2 border-white dark:border-slate-800 bg-slate-100 dark:bg-slate-700 flex items-center justify-center text-[8px] font-bold text-slate-400">
                {{ i }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div
      v-if="!projectStore.loading && projectStore.projects.length === 0"
      class="py-20 text-center bg-white dark:bg-slate-800 rounded-[3rem] border border-slate-100 dark:border-slate-700 border-dashed"
    >
      <div
        class="w-24 h-24 bg-slate-50 dark:bg-slate-900 rounded-full flex items-center justify-center mx-auto mb-6"
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
            d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
          />
        </svg>
      </div>
      <h3 class="text-xl font-bold text-slate-900 dark:text-white">
        No projects yet
      </h3>
      <p class="text-slate-500 dark:text-slate-400 mt-2 mb-8">
        Create your first project to start organizing tasks.
      </p>
      <button
        @click="showNewProjectModal = true"
        class="bg-indigo-600 text-white px-8 py-3 rounded-xl font-bold hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-200"
      >
        Create First Project
      </button>
    </div>

    <!-- New Project Modal -->
    <div
      v-if="showNewProjectModal"
      class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    >
      <div
        class="bg-white dark:bg-slate-900 rounded-3xl w-full max-w-lg shadow-2xl overflow-hidden border border-slate-100 dark:border-slate-800 flex flex-col max-h-[90vh]"
      >
        <div
          class="p-4 md:p-6 border-b border-slate-50 dark:border-slate-800 flex justify-between items-center shrink-0"
        >
          <h3 class="text-lg md:text-xl font-bold text-slate-900 dark:text-white">
            Create New Project
          </h3>
          <button
            @click="showNewProjectModal = false"
            class="text-slate-400 hover:text-slate-600"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="p-4 md:p-6 space-y-4 overflow-y-auto">
          <div>
            <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Project Name</label>
            <input
              v-model="newProject.name"
              type="text"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
              placeholder="e.g. Marketing Campaign"
            />
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Description</label>
            <textarea
              v-model="newProject.description"
              rows="3"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
              placeholder="What is this project about?"
            ></textarea>
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Color</label>
            <div class="flex space-x-2">
              <button
                v-for="color in colors"
                :key="color"
                @click="newProject.color = color"
                class="w-8 h-8 rounded-full transition-all transform hover:scale-110"
                :class="[
                  color,
                  newProject.color === color
                    ? 'ring-2 ring-offset-2 ring-slate-400 scale-110'
                    : '',
                ]"
              ></button>
            </div>
          </div>
        </div>
        <div class="p-4 md:p-6 bg-slate-50 dark:bg-slate-800/50 flex justify-end space-x-3 shrink-0">
          <button
            @click="showNewProjectModal = false"
            class="px-6 py-2 text-slate-600 dark:text-slate-400 font-bold"
          >
            Cancel
          </button>
          <button
            @click="handleCreateProject"
            class="px-6 py-2 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100"
          >
            Create Project
          </button>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div
      v-if="showDeleteModal"
      class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    >
      <div
        class="bg-white dark:bg-slate-900 rounded-3xl w-full max-w-sm shadow-2xl overflow-hidden border border-slate-100 dark:border-slate-800 p-6 text-center"
      >
         <div class="w-16 h-16 bg-rose-100 dark:bg-rose-900/20 text-rose-500 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
         </div>
         <h3 class="text-xl font-bold text-slate-900 dark:text-white mb-2">Delete Project?</h3>
         <p class="text-slate-500 dark:text-slate-400 mb-6">Are you sure you want to delete this project? All associated tasks will be orphaned.</p>
         <div class="flex space-x-3 justify-center">
            <button @click="showDeleteModal = false" class="px-6 py-2 text-slate-600 dark:text-slate-400 font-bold bg-slate-100 dark:bg-slate-800 rounded-xl hover:bg-slate-200 dark:hover:bg-slate-700 transition-colors">Cancel</button>
            <button @click="executeDelete" class="px-6 py-2 bg-rose-500 text-white rounded-xl font-bold hover:bg-rose-600 transition-colors shadow-lg shadow-rose-200 dark:shadow-none">Delete</button>
         </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted, inject, computed } from "vue";
import { useProjectStore } from "../stores/projects";

const projectStore = useProjectStore();
const showToast = inject("showToast");
const showNewProjectModal = ref(false);
const showDeleteModal = ref(false);
const projectToDelete = ref(null);
const newProject = ref({ name: "", description: "", color: "bg-indigo-500" });

const colors = [
  "bg-indigo-500",
  "bg-emerald-500",
  "bg-rose-500",
  "bg-amber-500",
  "bg-sky-500",
  "bg-violet-500",
];

const totalProjects = computed(() => projectStore.projects.length);
const activeProjects = computed(() => projectStore.projects.filter(p => calculateProjectProgress(p) < 100).length);
const completedProjects = computed(() => projectStore.projects.filter(p => calculateProjectProgress(p) === 100).length);

const calculateProjectProgress = (project) => {
  if (project.progress !== undefined && project.progress !== null) return project.progress;
  if (!project.goals || project.goals.length === 0) return 0;
  const totalProgress = project.goals.reduce(
    (acc, goal) => acc + (goal.progress || 0),
    0
  );
  return Math.round(totalProgress / project.goals.length);
};

async function handleCreateProject() {
  if (!newProject.value.name) return;
  try {
    await projectStore.createProject(newProject.value);
    showNewProjectModal.value = false;
    newProject.value = { name: "", description: "", color: "bg-indigo-500" };
    showToast("Project created successfully", "success");
  } catch (err) {
    showToast("Failed to create project", "error");
    console.error(err);
  }
}

function confirmDelete(project) {
  projectToDelete.value = project;
  showDeleteModal.value = true;
}

async function executeDelete() {
  if (!projectToDelete.value) return;
  try {
    await projectStore.deleteProject(projectToDelete.value.id);
    showToast("Project deleted", "success");
    showDeleteModal.value = false;
    projectToDelete.value = null;
  } catch (err) {
    showToast("Failed to delete project", "error");
  }
}

onMounted(() => {
  projectStore.fetchProjects();
});
</script>
