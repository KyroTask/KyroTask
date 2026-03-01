<template>
  <div class="space-y-8 pb-24 md:pb-12">
    <div v-if="loading" class="py-20 text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-4 border-slate-100 border-t-indigo-600 mx-auto"></div>
    </div>

    <div v-else-if="goal" class="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
      <!-- Page Header -->
      <div class="flex flex-col md:flex-row md:items-start justify-between gap-6 mb-6">
        <div class="flex-1 space-y-3">
          <div class="flex flex-wrap items-center gap-2">
            <span class="px-3 py-1 rounded-full bg-indigo-50 dark:bg-indigo-900/20 text-[10px] font-black uppercase tracking-widest text-indigo-600 dark:text-indigo-400 border border-indigo-100 dark:border-indigo-800">Objective</span>
            <span v-if="goal.target_date" class="px-3 py-1 rounded-full bg-amber-50 dark:bg-amber-900/20 text-[10px] font-black uppercase tracking-widest text-amber-600 dark:text-amber-400 border border-amber-100 dark:border-amber-800 flex items-center gap-1">
              <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
              {{ daysRemaining }} Days Left
            </span>
          </div>
          <h1 class="text-2xl md:text-4xl font-black tracking-tight text-slate-900 dark:text-white">{{ goal.title }}</h1>
          <p class="text-slate-500 dark:text-slate-400 text-sm md:text-base leading-relaxed max-w-2xl">{{ goal.description }}</p>
          
          <div v-if="goal.motivation" class="flex items-start gap-3 bg-indigo-50 dark:bg-indigo-900/10 rounded-2xl p-4 border border-indigo-100 dark:border-indigo-800 max-w-xl">
            <div class="w-8 h-8 rounded-xl bg-indigo-100 dark:bg-indigo-900/30 flex items-center justify-center shrink-0">
              <svg class="w-4 h-4 text-indigo-600 dark:text-indigo-400" fill="currentColor" viewBox="0 0 20 20"><path d="M10.394 2.08a1 1 0 00-.788 0l-7 3a1 1 0 000 1.84L5.25 8.051a.999.999 0 01.356-.257l4-1.714a1 1 0 11.788 1.838L7.667 9.088l1.94.831a1 1 0 00.787 0l7-3a1 1 0 000-1.838l-7-3zM3.31 9.397L5 10.12v4.102a8.969 8.969 0 00-1.05-.174 1 1 0 01-.89-.89 11.115 11.115 0 01.25-3.762zM9.3 16.573A9.026 9.026 0 007 14.935v-3.957l1.818.78a3 3 0 002.364 0l5.508-2.361a11.026 11.026 0 01.25 3.762 1 1 0 01-.89.89 8.968 8.968 0 00-5.35 2.524 1 1 0 01-1.4 0z"></path></svg>
            </div>
            <div>
              <p class="text-[10px] font-black text-indigo-600 dark:text-indigo-400 uppercase tracking-widest mb-1">Core Motivation</p>
              <p class="text-slate-700 dark:text-slate-300 italic font-medium">"{{ goal.motivation }}"</p>
            </div>
          </div>
        </div>

        <!-- Progress -->
        <div class="bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 rounded-2xl p-6 text-center shadow-sm shrink-0 w-36">
          <span class="text-3xl md:text-4xl font-black text-indigo-600">{{ calculateProgress }}%</span>
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mt-1">Complete</p>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
        <!-- Left Column: Milestones (8/12) -->
        <div class="lg:col-span-8 space-y-8">
          <div class="flex items-center justify-between px-4">
            <div class="flex items-center gap-4">
               <div class="w-12 h-12 rounded-2xl bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 flex items-center justify-center shadow-sm">
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4"></path></svg>
               </div>
               <div>
                 <h2 class="text-2xl font-black text-slate-900 dark:text-white tracking-tight">Strategic Milestones</h2>
                 <p class="text-sm text-slate-500 font-medium">Break down your vision into steps.</p>
               </div>
            </div>
            <button @click="showNewMilestoneModal = true" class="bg-slate-900 dark:bg-white text-white dark:text-slate-900 px-6 py-3 rounded-2xl font-black text-xs uppercase tracking-widest hover:scale-105 active:scale-95 transition-all shadow-xl shadow-slate-200 dark:shadow-none flex items-center gap-2">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4"></path></svg>
              <span>Add Step</span>
            </button>
          </div>

          <!-- Timeline Layout -->
          <div class="relative space-y-12 before:absolute before:inset-0 before:ml-6 before:-translate-x-px lg:before:ml-6 before:h-full before:w-0.5 before:bg-gradient-to-b before:from-indigo-500/50 before:via-slate-200 dark:before:via-slate-800 before:to-transparent">
            
            <div v-if="!goal.milestones || goal.milestones.length === 0" class="relative z-10 bg-white dark:bg-slate-900 p-12 rounded-[2.5rem] border border-slate-100 dark:border-slate-800 text-center shadow-sm">
               <div class="w-20 h-20 bg-slate-50 dark:bg-slate-800 rounded-full flex items-center justify-center mx-auto mb-6">
                 <svg class="w-10 h-10 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path></svg>
               </div>
               <h3 class="text-xl font-bold text-slate-900 dark:text-white mb-2">No milestones defined</h3>
               <p class="text-slate-500 dark:text-slate-400 max-w-xs mx-auto">Start by adding your first major milestone to track progress.</p>
            </div>

            <div v-for="(milestone, index) in goal.milestones" :key="milestone.id" class="relative flex items-start group">
               <!-- Timeline Indicator -->
               <div class="absolute left-0 top-0 mt-8 ml-6 -translate-x-1/2 w-4 h-4 rounded-full border-4 border-white dark:border-slate-950 z-10 transition-all duration-500 group-hover:scale-150 shadow-sm" :class="milestone.status === 'completed' ? 'bg-emerald-500' : 'bg-indigo-500'"></div>

               <div class="ml-14 w-full bg-white dark:bg-slate-900 rounded-[2.5rem] border border-slate-100 dark:border-slate-800 overflow-hidden shadow-sm hover:shadow-2xl hover:shadow-indigo-100/50 dark:hover:shadow-none transition-all duration-500 group-hover:-translate-y-1">
                  <!-- Milestone Header -->
                  <div class="p-8 border-b border-slate-50 dark:border-slate-800 flex items-start justify-between bg-slate-50/30 dark:bg-slate-800/20">
                     <div class="space-y-1">
                        <div class="flex items-center gap-3">
                          <span class="text-[10px] font-black text-indigo-500 uppercase tracking-[0.2em]">Milestone {{ index + 1 }}</span>
                          <span v-if="milestone.status === 'completed'" class="px-2 py-0.5 rounded-md bg-emerald-500/10 text-emerald-500 text-[8px] font-black uppercase tracking-widest">Completed</span>
                        </div>
                        <h3 class="text-2xl font-black text-slate-900 dark:text-white tracking-tight">{{ milestone.title }}</h3>
                     </div>
                     <button @click="deleteMilestone(milestone.id)" class="text-slate-300 hover:text-rose-500 transition-all p-2 hover:bg-rose-50 dark:hover:bg-rose-900/20 rounded-xl">
                        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                     </button>
                  </div>

                  <!-- Tasks List -->
                  <div class="p-6 space-y-3">
                     <div v-for="task in milestone.tasks" :key="task.id" 
                          class="flex items-center p-4 rounded-2xl transition-all duration-300 group/task cursor-pointer border border-transparent"
                          :class="task.status === 'completed' ? 'bg-emerald-50/30 dark:bg-emerald-900/5' : 'bg-slate-50/50 dark:bg-slate-800/30 hover:border-indigo-100 dark:hover:border-indigo-900/50 hover:bg-white dark:hover:bg-slate-800'"
                          @click="toggleTask(task)">
                        
                        <div class="w-7 h-7 rounded-xl border-2 mr-4 flex items-center justify-center transition-all duration-500"
                             :class="task.status === 'completed' ? 'bg-emerald-500 border-emerald-500 text-white scale-110 shadow-lg shadow-emerald-200' : 'border-slate-300 dark:border-slate-700 group-hover/task:border-indigo-500 group-hover/task:scale-110'">
                           <svg v-if="task.status === 'completed'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"></path></svg>
                        </div>
                        
                        <div class="flex-1">
                          <span class="text-base font-bold transition-all duration-500 block" 
                                :class="task.status === 'completed' ? 'text-slate-400 line-through' : 'text-slate-700 dark:text-slate-200'">
                             {{ task.title }}
                          </span>
                          <span v-if="task.priority" class="text-[8px] font-black uppercase tracking-widest" :class="task.priority === 'high' ? 'text-rose-500' : 'text-slate-400'">{{ task.priority }} Priority</span>
                        </div>
                     </div>

                     <button @click="openNewTaskModal(milestone.id)" class="w-full py-4 mt-2 text-[10px] font-black text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 dark:hover:bg-indigo-900/20 border-2 border-dashed border-slate-100 dark:border-slate-800 hover:border-indigo-200 dark:hover:border-indigo-900 rounded-2xl flex items-center justify-center gap-2 transition-all duration-300 uppercase tracking-[0.2em]">
                        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M12 4v16m8-8H4"></path></svg>
                        <span>Add Action Item</span>
                     </button>
                  </div>
               </div>
            </div>
          </div>
        </div>

        <!-- Right Column: Habits & Notes (4/12) -->
        <div class="lg:col-span-4 space-y-8">
          <!-- Habits Section -->
          <div class="space-y-6">
            <div class="flex items-center justify-between px-2">
              <div class="flex items-center gap-3">
                 <div class="w-10 h-10 rounded-xl bg-orange-100 dark:bg-orange-900/30 text-orange-600 dark:text-orange-400 flex items-center justify-center shadow-sm">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
                 </div>
                 <h2 class="text-xl font-black text-slate-900 dark:text-white tracking-tight">Daily Systems</h2>
              </div>
               <button @click="showNewHabitModal = true" class="text-[10px] font-black text-indigo-600 hover:text-indigo-700 bg-indigo-50 dark:bg-indigo-900/20 px-4 py-2 rounded-xl transition-all uppercase tracking-widest">
                <span>Add Habit</span>
              </button>
            </div>

            <div class="space-y-4">
              <div v-if="!goal.habits || goal.habits.length === 0" class="bg-white dark:bg-slate-900 p-8 rounded-[2rem] border border-slate-100 dark:border-slate-800 text-center">
                 <p class="text-slate-400 text-sm font-medium italic">No habits linked yet.</p>
              </div>

              <div v-for="habit in goal.habits" :key="habit.id" class="bg-white dark:bg-slate-900 p-6 rounded-[2rem] border border-slate-100 dark:border-slate-800 shadow-sm hover:shadow-xl transition-all duration-500 group">
                <div class="flex justify-between items-start mb-6">
                  <div>
                     <h3 class="font-black text-slate-900 dark:text-white text-lg tracking-tight">{{ habit.name }}</h3>
                     <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mt-1">{{ habit.frequency }}</p>
                  </div>
                  <div class="flex items-center gap-1.5 text-[10px] font-black text-orange-500 bg-orange-50 dark:bg-orange-900/30 px-3 py-1.5 rounded-xl border border-orange-100 dark:border-orange-900/50">
                    <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M12.395 2.553a1 1 0 00-1.45-.385c-.345.23-.614.558-.822.88-.214.33-.403.713-.57 1.116-.334.804-.614 1.768-.84 2.734a31.365 31.365 0 00-.613 3.58 2.64 2.64 0 01-.945-1.067c-.328-.68-.398-1.534-.398-2.654A1 1 0 005.05 6.05 6.981 6.981 0 003 11a7 7 0 1011.95-4.95c-.592-.591-.98-.985-1.348-1.467-.363-.476-.724-1.063-1.207-2.03zM12.12 15.12A3 3 0 017 13s.879.5 2.5.5c0-1 .5-4 1.25-4.5.5 1 .786 1.293 1.371 1.879A2.99 2.99 0 0113 13a2.99 2.99 0 01-.879 2.121z" clip-rule="evenodd"></path></svg>
                    <span>{{ habit.current_streak }} DAY STREAK</span>
                  </div>
                </div>
                
                <button class="w-full py-4 bg-slate-50 dark:bg-slate-800 hover:bg-indigo-600 hover:text-white text-slate-600 dark:text-slate-400 rounded-2xl font-black text-xs uppercase tracking-widest transition-all duration-500 flex items-center justify-center gap-2 group-hover:scale-[1.02] shadow-sm">
                   <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7"></path></svg>
                   <span>Check In</span>
                </button>
              </div>
            </div>
          </div>

          <!-- Notes Section -->
          <div class="bg-white dark:bg-slate-900 rounded-[2.5rem] p-8 border border-slate-100 dark:border-slate-800 shadow-sm space-y-6">
            <div class="flex items-center gap-4">
               <div class="w-10 h-10 rounded-xl bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 flex items-center justify-center shadow-sm">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-5M16.5 3.5a2.121 2.121 0 013 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
               </div>
               <h2 class="text-xl font-black text-slate-900 dark:text-white tracking-tight">Journal & Logs</h2>
            </div>
            <div class="relative group">
               <textarea v-model="goal.notes" @change="updateNotes" class="w-full bg-slate-50 dark:bg-slate-800 border-2 border-transparent focus:border-indigo-500/20 rounded-[2rem] p-6 text-slate-700 dark:text-slate-300 focus:ring-0 resize-none leading-relaxed transition-all duration-500 min-h-[200px]" placeholder="Document your journey, insights, and breakthroughs..."></textarea>
               <div class="absolute bottom-6 right-6 flex items-center gap-2 text-[8px] font-black text-slate-400 uppercase tracking-widest opacity-0 group-focus-within:opacity-100 transition-opacity">
                 <div class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></div>
                 <span>Syncing to Cloud</span>
               </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modals (Keep existing structure but improve styling if needed) -->
    <!-- New Milestone Modal -->
    <div v-if="showNewMilestoneModal" class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-white dark:bg-slate-900 rounded-3xl w-full max-w-md p-8 space-y-6 shadow-2xl border border-slate-100 dark:border-slate-800 transform transition-all scale-100">
        <h3 class="text-2xl font-bold text-slate-900 dark:text-white">Add Milestone</h3>
        <input v-model="newMilestone.title" type="text" placeholder="What's the next big step?" class="w-full bg-slate-50 dark:bg-slate-800 border-0 rounded-2xl px-6 py-4 dark:text-white focus:ring-2 focus:ring-indigo-500 text-lg">
        <div class="flex justify-end space-x-3">
          <button @click="showNewMilestoneModal = false" class="px-6 py-3 text-slate-500 font-bold hover:bg-slate-50 rounded-xl transition-colors">Cancel</button>
          <button @click="createMilestone" class="px-8 py-3 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-200 dark:shadow-none">Add Milestone</button>
        </div>
      </div>
    </div>

    <!-- New Task Modal -->
    <div v-if="showNewTaskModal" class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-white dark:bg-slate-900 rounded-3xl w-full max-w-md p-8 space-y-6 shadow-2xl border border-slate-100 dark:border-slate-800">
        <h3 class="text-2xl font-bold text-slate-900 dark:text-white">Add Task</h3>
        <input v-model="newTask.title" type="text" placeholder="What needs to be done?" class="w-full bg-slate-50 dark:bg-slate-800 border-0 rounded-2xl px-6 py-4 dark:text-white focus:ring-2 focus:ring-indigo-500 text-lg">
        <div class="flex justify-end space-x-3">
          <button @click="showNewTaskModal = false" class="px-6 py-3 text-slate-500 font-bold hover:bg-slate-50 rounded-xl transition-colors">Cancel</button>
          <button @click="createTask" class="px-8 py-3 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-200 dark:shadow-none">Add Task</button>
        </div>
      </div>
    </div>

    <!-- New Habit Modal -->
    <div v-if="showNewHabitModal" class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-50 flex items-center justify-center p-4">
      <div class="bg-white dark:bg-slate-900 rounded-3xl w-full max-w-md p-8 space-y-6 shadow-2xl border border-slate-100 dark:border-slate-800">
        <h3 class="text-2xl font-bold text-slate-900 dark:text-white">Add Habit</h3>
        <input v-model="newHabit.name" type="text" placeholder="Habit Name" class="w-full bg-slate-50 dark:bg-slate-800 border-0 rounded-2xl px-6 py-4 dark:text-white focus:ring-2 focus:ring-indigo-500 text-lg">
         <div class="space-y-2">
            <label class="text-xs font-bold text-slate-400 uppercase tracking-widest ml-2">Frequency</label>
            <select v-model="newHabit.frequency" class="w-full bg-slate-50 dark:bg-slate-800 border-0 rounded-2xl px-6 py-4 dark:text-white focus:ring-2 focus:ring-indigo-500 appearance-none">
               <option value="daily">Daily</option>
               <option value="weekly">Weekly</option>
            </select>
         </div>
        <div class="flex justify-end space-x-3">
          <button @click="showNewHabitModal = false" class="px-6 py-3 text-slate-500 font-bold hover:bg-slate-50 rounded-xl transition-colors">Cancel</button>
          <button @click="createHabit" class="px-8 py-3 bg-indigo-600 text-white rounded-xl font-bold hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-200 dark:shadow-none">Add Habit</button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted, computed, inject } from 'vue';
import { useRoute } from 'vue-router';
import { useGoalStore } from '../stores/goals';
import { useTaskStore } from '../stores/tasks';
// import { useHabitStore } from '../stores/habits'; // Assuming habit store exists or we use goal store
import api from '../services/api';

const route = useRoute();
const goalStore = useGoalStore();
const taskStore = useTaskStore();
const showToast = inject('showToast');

const goal = ref(null);
const loading = ref(true);
const showNewMilestoneModal = ref(false);
const showNewTaskModal = ref(false);
const showNewHabitModal = ref(false);

const newMilestone = ref({ title: '' });
const newTask = ref({ title: '', milestone_id: null });
const newHabit = ref({ name: '', frequency: 'daily' });

const daysRemaining = computed(() => {
  if (!goal.value?.target_date) return 0;
  const target = new Date(goal.value.target_date);
  const now = new Date();
  const diff = target - now;
  return Math.ceil(diff / (1000 * 60 * 60 * 24));
});

const calculateProgress = computed(() => {
  if (!goal.value) return 0;
  // Formula: (Completed Tasks + Completed Habits / Total Items) * 100?
  // Or just based on milestones?
  // Let's use milestones for now if they exist, otherwise tasks.
  // Actually, let's count all tasks inside milestones.
  let totalTasks = 0;
  let completedTasks = 0;

  if (goal.value.milestones) {
    goal.value.milestones.forEach(m => {
      if (m.tasks) {
        m.tasks.forEach(t => {
          totalTasks++;
          if (t.status === 'completed') completedTasks++;
        });
      }
    });
  }
  
  // Also include direct tasks if any (though we are pushing for milestones)
  // And habits? Habits are continuous. Maybe they don't contribute to % completion of the goal directly, 
  // or they contribute to "daily success".
  // Let's stick to tasks for the progress bar for now as per "Goal Completion".
  
  if (totalTasks === 0) return goal.value.progress || 0; // Fallback to manual or 0
  return Math.round((completedTasks / totalTasks) * 100);
});

async function loadGoal() {
  loading.value = true;
  try {
    // Use the updated store action that handles both ID and slug
    goal.value = await goalStore.fetchGoal(route.params.slug);
  } catch (err) {
    console.error(err);
    showToast('Failed to load goal', 'error');
  } finally {
    loading.value = false;
  }
}

async function createMilestone() {
  if (!newMilestone.value.title) return;
  try {
    await goalStore.createMilestone({
      goal_id: goal.value.id,
      title: newMilestone.value.title
    });
    await loadGoal(); // Refresh to show new milestone
    showNewMilestoneModal.value = false;
    newMilestone.value = { title: '' };
    showToast('Milestone added', 'success');
  } catch (err) {
    showToast('Failed to add milestone', 'error');
  }
}

async function deleteMilestone(id) {
  if (!confirm('Delete milestone?')) return;
  try {
    await goalStore.deleteMilestone(id);
    await loadGoal();
  } catch (err) {
    showToast('Failed to delete milestone', 'error');
  }
}

function openNewTaskModal(milestoneId) {
  newTask.value.milestone_id = milestoneId;
  showNewTaskModal.value = true;
}

async function createTask() {
  if (!newTask.value.title) return;
  try {
    // We need to create a task linked to this goal and milestone
    await taskStore.createTask({
      title: newTask.value.title,
      goal_id: goal.value.id,
      milestone_id: newTask.value.milestone_id,
      project_id: goal.value.project_id // Optional, inherit project
    });
    await loadGoal();
    showNewTaskModal.value = false;
    newTask.value = { title: '', milestone_id: null };
    showToast('Task added', 'success');
  } catch (err) {
    showToast('Failed to add task', 'error');
  }
}

async function toggleTask(task) {
  try {
    const newStatus = task.status === 'completed' ? 'pending' : 'completed';
    await taskStore.updateTask(task.id, { status: newStatus });
    // Update local state immediately for responsiveness
    task.status = newStatus;
    // Optionally reload goal to update progress bar if it depends on backend calculation
    // But our computed property handles it if we update the local object correctly.
  } catch (err) {
    showToast('Failed to update task', 'error');
  }
}

async function createHabit() {
  if (!newHabit.value.name) return;
  try {
    // Assuming we have a habit store or use api directly
    await api.post('/habits', {
      ...newHabit.value,
      goal_id: goal.value.id
    });
    await loadGoal();
    showNewHabitModal.value = false;
    newHabit.value = { name: '', frequency: 'daily' };
    showToast('Habit added', 'success');
  } catch (err) {
    showToast('Failed to add habit', 'error');
  }
}

async function updateNotes() {
  try {
    await goalStore.updateGoal(goal.value.id, { notes: goal.value.notes });
    showToast('Notes saved', 'success');
  } catch (err) {
    showToast('Failed to save notes', 'error');
  }
}

onMounted(loadGoal);
</script>

<style scoped>
.dashed-border {
  background-image: url("data:image/svg+xml,%3csvg width='100%25' height='100%25' xmlns='http://www.w3.org/2000/svg'%3e%3crect width='100%25' height='100%25' fill='none' rx='8' ry='8' stroke='%2394A3B8FF' stroke-width='2' stroke-dasharray='6%2c 14' stroke-dashoffset='0' stroke-linecap='square'/%3e%3c/svg%3e");
}
</style>
