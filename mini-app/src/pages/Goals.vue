<template>
  <div class="space-y-8">
    <!-- Page Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-4">
      <div>
        <h1 class="text-2xl md:text-3xl font-black tracking-tight text-slate-900 dark:text-white">Goals</h1>
        <p class="text-slate-500 dark:text-slate-400 text-sm font-medium mt-1">Focus on what matters most.</p>
      </div>
      <button
        @click="showNewGoalModal = true"
        class="bg-indigo-600 text-white px-4 py-2 rounded-xl font-bold hover:bg-indigo-700 transition-all flex items-center space-x-2 shadow-lg shadow-indigo-200 dark:shadow-none group text-sm"
      >
        <svg class="w-4 h-4 group-hover:rotate-90 transition-transform duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4" />
        </svg>
        <span>New Goal</span>
      </button>
    </div>

    <!-- Stats Bar -->
    <div class="flex items-center space-x-3 overflow-x-auto pb-2 no-scrollbar mb-4">
      <div class="bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 px-4 py-2 rounded-xl flex items-center space-x-2 whitespace-nowrap shadow-sm">
        <span class="text-lg font-black text-slate-900 dark:text-white">{{ totalGoals }}</span>
        <span class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Total</span>
      </div>
      <div class="bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 px-4 py-2 rounded-xl flex items-center space-x-2 whitespace-nowrap shadow-sm">
        <span class="text-lg font-black text-emerald-600">{{ completedGoals }}</span>
        <span class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Done</span>
      </div>
      <div class="bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 px-4 py-2 rounded-xl flex items-center space-x-2 whitespace-nowrap shadow-sm">
        <span class="text-lg font-black text-amber-600">{{ inProgressGoals }}</span>
        <span class="text-[10px] font-bold uppercase tracking-widest text-slate-400">Active</span>
      </div>
    </div>

    <!-- Goals Bento Grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 md:gap-6 auto-rows-[minmax(160px,auto)]">
      <div v-if="goalStore.loading" class="col-span-full py-20 text-center">
        <div
          class="animate-spin rounded-full h-12 w-12 border-4 border-slate-100 border-t-indigo-600 mx-auto"
        ></div>
      </div>
      <div
        v-else
        v-for="(goal, index) in goalStore.goals"
        :key="goal.id"
        @click="router.push('/goals/' + (goal.slug || goal.id))"
        :class="[
          'bg-white dark:bg-slate-800 rounded-[2rem] md:rounded-[2rem] border border-slate-100 dark:border-slate-700 shadow-sm hover:shadow-2xl hover:shadow-indigo-100 dark:hover:shadow-none transition-all duration-500 p-6 md:p-6 group relative overflow-hidden flex flex-col cursor-pointer',
          index === 0 && goalStore.goals.length > 1 ? 'md:col-span-2 md:row-span-2' : '',
          goal.progress > 80 && index !== 0 && goalStore.goals.length > 2 ? 'md:row-span-2' : ''
        ]"
      >
        <!-- Hover Gradient Background -->
        <div class="absolute inset-0 bg-gradient-to-br from-transparent to-indigo-50/50 dark:to-indigo-900/10 opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none"></div>

        <div class="flex items-start justify-between mb-4 md:mb-5 relative z-10">
          <div
            :class="[
              'bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 dark:text-indigo-400 rounded-xl md:rounded-2xl flex items-center justify-center shadow-sm group-hover:scale-110 transition-transform duration-300',
              index === 0 && goalStore.goals.length > 1 ? 'w-14 h-14 md:w-16 md:h-16' : 'w-12 h-12 md:w-12 md:h-12'
            ]"
          >
            <svg
              :class="index === 0 && goalStore.goals.length > 1 ? 'w-8 h-8 md:w-9 md:h-9' : 'w-6 h-6 md:w-6 md:h-6'"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M13 10V3L4 14h7v7l9-11h-7z"
              />
            </svg>
          </div>
          <div class="flex items-center space-x-1">
            <button
              @click.stop="openEditModal(goal)"
              class="p-2 text-slate-300 dark:text-slate-600 hover:text-indigo-500 hover:bg-indigo-50 dark:hover:bg-indigo-900/20 rounded-xl transition-colors relative z-20"
            >
              <svg
                class="w-4 h-4 md:w-5 md:h-5"
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
              @click.stop="confirmDelete(goal)"
              class="p-2 text-slate-300 dark:text-slate-600 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-900/20 rounded-xl transition-colors relative z-20"
            >
              <svg
                class="w-4 h-4 md:w-5 md:h-5"
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

        <div class="flex-1">
           <h3 :class="[
             'font-bold text-slate-900 dark:text-white mb-1 md:mb-2 relative z-10 pointer-events-none group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors',
             index === 0 ? 'text-2xl md:text-4xl' : 'text-xl md:text-2xl'
           ]">
             {{ goal.title }}
           </h3>
           <p
             :class="[
               'text-slate-500 dark:text-slate-400 mb-6 md:mb-8 leading-relaxed relative z-10 pointer-events-none',
               index === 0 ? 'text-sm md:text-lg line-clamp-3' : 'text-xs md:text-sm line-clamp-2'
             ]"
           >
             {{ goal.description }}
           </p>
        </div>

        <div class="space-y-4 relative z-10 mt-auto">
          <div
            class="flex items-center justify-between text-xs font-black uppercase tracking-widest"
          >
            <span class="text-slate-400 dark:text-slate-500">Progress</span>
            <span class="text-indigo-600">{{ goal.progress || 0 }}%</span>
          </div>
          <div
            class="w-full bg-slate-100 dark:bg-slate-700 h-3 rounded-full overflow-hidden relative group/progress"
          >
            <div
              class="bg-indigo-600 h-full transition-all duration-1000 shadow-[0_0_10px_rgba(79,70,229,0.4)]"
              :style="{ width: (goal.progress || 0) + '%' }"
            ></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div
      v-if="!goalStore.loading && goalStore.goals.length === 0"
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
            d="M13 10V3L4 14h7v7l9-11h-7z"
          />
        </svg>
      </div>
      <h3 class="text-xl font-bold text-slate-900 dark:text-white">
        No goals set
      </h3>
      <p class="text-slate-500 dark:text-slate-400 mt-2 mb-8">
        Define your long-term objectives and track progress.
      </p>
      <button
        @click="showNewGoalModal = true"
        class="bg-indigo-600 text-white px-8 py-3 rounded-xl font-bold hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-200"
      >
        Create First Goal
      </button>
    </div>

    <!-- New Goal Modal -->
    <div
      v-if="showNewGoalModal"
      class="fixed inset-0 bg-slate-900/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
    >
      <div
        class="bg-white dark:bg-slate-900 rounded-3xl w-full max-w-lg shadow-2xl overflow-hidden border border-slate-100 dark:border-slate-800 flex flex-col max-h-[90vh]"
      >
        <div
          class="p-4 md:p-6 border-b border-slate-50 dark:border-slate-800 flex justify-between items-center shrink-0"
        >
          <h3 class="text-lg md:text-xl font-bold text-slate-900 dark:text-white">
            Create New Goal
          </h3>
          <button
            @click="showNewGoalModal = false"
            class="text-slate-400 hover:text-slate-600"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="p-4 md:p-6 space-y-4 overflow-y-auto">
          <div>
            <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Goal Title</label>
            <input
              v-model="newGoal.title"
              type="text"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
              placeholder="e.g. Launch New Product"
            />
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Description</label>
            <textarea
              v-model="newGoal.description"
              rows="2"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
              placeholder="What is the ultimate outcome?"
            ></textarea>
          </div>
          <div>
            <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Motivation (The Why)</label>
            <textarea
              v-model="newGoal.motivation"
              rows="2"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
              placeholder="Why does this goal matter?"
            ></textarea>
          </div>
          <div class="grid grid-cols-2 gap-4">
             <div>
              <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Deadline</label>
              <input
                v-model="newGoal.target_date"
                type="date"
                class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-lg px-2 py-2 md:rounded-xl md:px-4 md:py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white text-xs md:text-base"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2">Project</label>
              <select
                v-model="newGoal.project_id"
                class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 appearance-none text-slate-900 dark:text-white"
              >
                <option :value="null">No Project</option>
                <option v-for="project in projectStore.projects" :key="project.id" :value="project.id">{{ project.name }}</option>
              </select>
            </div>
          </div>
        </div>
        <div class="p-4 md:p-6 bg-slate-50 dark:bg-slate-800/50 flex justify-end space-x-3 shrink-0">
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

    <!-- Edit Goal Modal -->
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
            Edit Goal
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
              >Goal Title</label
            >
            <input
              v-model="editGoalData.title"
              type="text"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
            />
          </div>
          <div>
            <label
              class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2"
              >Description</label
            >
            <textarea
              v-model="editGoalData.description"
              rows="3"
              class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 text-slate-900 dark:text-white"
            ></textarea>
            <div>
              <label
                class="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-2 mt-4"
                >Project</label
              >
              <select
                v-model="editGoalData.project_id"
                class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-xl px-4 py-3 outline-none focus:ring-2 focus:ring-indigo-500 appearance-none text-slate-900 dark:text-white"
              >
                <option :value="null">No Project</option>
                <option
                  v-for="project in projectStore.projects"
                  :key="project.id"
                  :value="project.id"
                >
                  {{ project.name }}
                </option>
              </select>
            </div>
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
            @click="handleUpdateGoal"
            class="px-6 py-2 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100"
          >
            Save Changes
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
         <h3 class="text-xl font-bold text-slate-900 dark:text-white mb-2">Delete Goal?</h3>
         <p class="text-slate-500 dark:text-slate-400 mb-6">Are you sure you want to delete this goal? This action cannot be undone.</p>
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
import { useRouter } from "vue-router";
import { useGoalStore } from "../stores/goals";
import { useProjectStore } from "../stores/projects";

const router = useRouter();
const goalStore = useGoalStore();
const projectStore = useProjectStore();
const showToast = inject("showToast");
const showNewGoalModal = ref(false);
const showEditModal = ref(false);
const showDeleteModal = ref(false);
const goalToDelete = ref(null);

const newGoal = ref({ title: "", description: "", motivation: "", target_date: "", project_id: null });
const editGoalData = ref({});

const totalGoals = computed(() => goalStore.goals.length);
const completedGoals = computed(() => goalStore.goals.filter(g => (g.progress || 0) === 100).length);
const inProgressGoals = computed(() => goalStore.goals.filter(g => (g.progress || 0) < 100 && (g.progress || 0) > 0).length);

async function handleCreateGoal() {
  if (!newGoal.value.title) return;
  try {
    // Convert date string to ISO if needed, but backend usually handles standard date strings
    const goalPayload = { ...newGoal.value };
    if (goalPayload.target_date) {
       goalPayload.target_date = new Date(goalPayload.target_date).toISOString();
    } else {
       delete goalPayload.target_date;
    }

    await goalStore.createGoal(goalPayload);
    showNewGoalModal.value = false;
    newGoal.value = { title: "", description: "", motivation: "", target_date: "", project_id: null };
    showToast("Goal created successfully", "success");
  } catch (err) {
    showToast("Failed to create goal", "error");
    console.error(err);
  }
}

function confirmDelete(goal) {
  goalToDelete.value = goal;
  showDeleteModal.value = true;
}

async function executeDelete() {
  if (!goalToDelete.value) return;
  try {
    await goalStore.deleteGoal(goalToDelete.value.id);
    showToast("Goal deleted", "success");
    showDeleteModal.value = false;
    goalToDelete.value = null;
  } catch (err) {
    showToast("Failed to delete goal", "error");
  }
}

function openEditModal(goal) {
  editGoalData.value = { ...goal };
  showEditModal.value = true;
}

async function handleUpdateGoal() {
  try {
    await goalStore.updateGoal(editGoalData.value.id, editGoalData.value);
    showEditModal.value = false;
    showToast("Goal updated successfully", "success");
  } catch (err) {
    showToast("Failed to update goal", "error");
    console.error(err);
  }
}

onMounted(() => {
  goalStore.fetchGoals();
  projectStore.fetchProjects();
});
</script>
