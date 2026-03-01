<template>
  <div class="pomodoro-wrap flex flex-col">

    <div class="relative z-10 h-full flex flex-col flex-1 min-h-0">

      <!-- Loading -->
      <div v-if="pomodoroStore.loading" class="flex-1 flex items-center justify-center">
        <div class="w-10 h-10 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin"></div>
      </div>

      <template v-else-if="pomodoroStore.status">

        <!-- ═══ PHASE TRACK (compact) ═══ -->
        <div class="hidden sm:block bg-white dark:bg-slate-900 rounded-2xl p-4 border border-slate-200 dark:border-slate-800 shadow-md shrink-0 mb-3 w-full max-w-4xl mx-auto">
          <!-- Nodes row -->
          <div class="relative flex items-center justify-between px-4 sm:px-8">
            <!-- Bg line -->
            <div class="absolute left-8 right-8 top-1/2 -translate-y-1/2 h-0.5 bg-slate-200 dark:bg-slate-800 rounded-full hidden sm:block"></div>
            <!-- Filled line -->
            <div class="absolute left-8 top-1/2 -translate-y-1/2 h-0.5 rounded-full hidden sm:block transition-all duration-1000"
              :style="{ width: progressLineWidth }" :class="phaseStyle.progressLine"></div>

            <div v-for="phase in phases" :key="phase.id" class="relative z-10 flex flex-col items-center gap-1.5 w-16">
              <div class="w-10 h-10 rounded-full flex shrink-0 items-center justify-center text-base border-[3px] shadow transition-all duration-500 bg-white dark:bg-slate-900"
                :class="[
                  currentPhaseId === phase.id ? `${phase.glowClass} bg-white dark:bg-slate-900 scale-110` :
                  isPastPhase(phase.id) ? 'bg-slate-100 dark:bg-slate-800 border-slate-300 dark:border-slate-700 text-slate-400' :
                  'bg-slate-50 dark:bg-slate-900 border-slate-100 dark:border-slate-800 text-slate-300 dark:text-slate-600'
                ]">
                {{ phase.emoji }}
                <div v-if="currentPhaseId === phase.id" class="absolute inset-0 rounded-full animate-ping opacity-20" :class="phase.pulseRing"></div>
              </div>
              <span class="text-[9px] font-bold uppercase tracking-wider"
                :class="currentPhaseId === phase.id ? 'text-slate-900 dark:text-white' : 'text-slate-400 dark:text-slate-600'">
                {{ phase.name }}
              </span>
              <div v-if="currentPhaseId === phase.id"
                class="px-1.5 py-0.5 bg-indigo-50 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 text-[9px] font-black rounded border border-indigo-100 dark:border-indigo-800 whitespace-nowrap">
                LVL {{ currentLevel }}
              </div>
            </div>
          </div>
          <!-- Mini level bar -->
          <div class="mt-3 flex items-center gap-2 px-1 max-w-xl mx-auto">
            <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest w-6">Lvl</span>
            <div class="flex-1 h-1 bg-slate-200 dark:bg-slate-800 rounded-full overflow-hidden flex">
              <div v-for="i in phaseRange" :key="i" class="flex-1 border-r border-white/20 last:border-0 transition-all duration-500"
                :class="i <= currentLevel ? phaseStyle.barActive : ''"></div>
            </div>
            <span class="text-[9px] font-mono text-slate-500 w-7 text-right">{{ currentLevel }}/{{ phaseMax }}</span>
          </div>
        </div>

        <!-- ═══ NORMAL PHASE: Timer only, full-width centered ═══ -->
        <div v-if="status.can_start_today && isNormalPhase" class="flex-1 flex flex-col items-center justify-center min-h-0 py-4">
          <div class="bg-white dark:bg-slate-900 rounded-3xl border border-slate-200 dark:border-slate-800 shadow-xl p-6 sm:p-8 flex flex-col items-center relative overflow-hidden w-full max-w-lg md:max-w-xl lg:max-w-2xl mx-auto my-auto">

            <!-- Glow -->
            <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[30rem] h-[30rem] rounded-full blur-[100px] opacity-10 pointer-events-none"
              :class="isWorkCycle ? 'bg-indigo-500' : 'bg-emerald-500'"></div>

            <!-- Status badge -->
            <div class="inline-flex items-center gap-2 px-4 py-1.5 rounded-full bg-indigo-50 dark:bg-indigo-900/20 border border-indigo-100 dark:border-indigo-800 mb-5">
              <span class="relative flex h-2.5 w-2.5">
                <span class="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75" :class="isWorkCycle ? 'bg-indigo-400' : 'bg-emerald-400'"></span>
                <span class="relative inline-flex rounded-full h-2.5 w-2.5" :class="isWorkCycle ? 'bg-indigo-500' : 'bg-emerald-500'"></span>
              </span>
              <span class="text-xs font-bold uppercase tracking-widest" :class="isWorkCycle ? 'text-indigo-600 dark:text-indigo-400' : 'text-emerald-600 dark:text-emerald-400'">
                {{ isWorkCycle ? 'Focus Mode' : 'Break Time' }}
              </span>
            </div>

            <!-- Large SVG Timer -->
            <div class="relative cursor-pointer my-2" @click="togglePause">
              <div class="absolute inset-0 rounded-full blur-2xl opacity-20" :class="isWorkCycle ? 'bg-indigo-500' : 'bg-emerald-500'"></div>
              <svg class="timer-svg-large transform -rotate-90" viewBox="0 0 200 200">
                <circle cx="100" cy="100" r="88" fill="none" stroke="currentColor" stroke-width="6" class="text-slate-100 dark:text-slate-800/50"/>
                <circle cx="100" cy="100" r="88" fill="none" stroke="url(#timerGrad)" stroke-width="8"
                  stroke-linecap="round" :stroke-dasharray="553" :stroke-dashoffset="dashOffset"
                  class="transition-all duration-1000 ease-linear"/>
                <defs>
                  <linearGradient id="timerGrad" x1="0%" y1="0%" x2="100%" y2="0%">
                    <stop offset="0%" :stop-color="isWorkCycle ? '#6366f1' : '#10b981'"/>
                    <stop offset="100%" :stop-color="isWorkCycle ? '#a855f7' : '#3b82f6'"/>
                  </linearGradient>
                </defs>
              </svg>
              <div class="absolute inset-0 flex flex-col items-center justify-center">
                <div class="timer-time-large text-slate-800 dark:text-white">{{ formattedTime }}</div>
                <div class="text-xs font-bold text-slate-400 uppercase tracking-widest mt-2">{{ isWorkCycle ? 'Remaining' : 'Resting' }}</div>
              </div>
            </div>

            <!-- Info row -->
            <div class="flex items-center justify-center gap-3 sm:gap-6 mt-3 sm:mt-4 text-[10px] sm:text-xs text-slate-400 font-semibold w-full">
              <span>{{ status.required_work_duration }}m work</span>
              <span class="w-1 h-1 rounded-full bg-slate-300 shrink-0"></span>
              <span>{{ status.required_break_duration }}m break</span>
              <span class="w-1 h-1 rounded-full bg-slate-300 shrink-0"></span>
              <span>{{ status.target_cycles }} cycle{{ status.target_cycles > 1 ? 's' : '' }}</span>
            </div>

            <!-- Controls (pushes to bottom) -->
            <div class="mt-auto pt-3 sm:pt-6 pb-2 w-full flex justify-center">
              <div class="flex items-center gap-3 p-1.5 sm:p-2 bg-slate-50 dark:bg-slate-800/50 rounded-2xl border border-slate-200 dark:border-slate-700/50 backdrop-blur-sm">
                <button v-if="!hasActiveSession" @click="handleStart"
                  class="px-6 sm:px-8 py-2 sm:py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white rounded-xl font-bold text-xs sm:text-sm tracking-wide hover:scale-105 active:scale-95 transition-all flex items-center gap-2 shadow-lg shadow-indigo-500/30">
                  <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
                  Start Session
                </button>
                <template v-else>
                  <button @click="togglePause"
                    class="w-10 sm:w-12 h-10 sm:h-12 flex items-center justify-center rounded-xl bg-white dark:bg-slate-700 text-slate-900 dark:text-white shadow-md hover:bg-slate-50 dark:hover:bg-slate-600 transition-all active:scale-95">
                    <svg v-if="isPaused" class="w-5 h-5 ml-0.5" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
                    <svg v-else class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/></svg>
                  </button>
                  <div class="w-px h-6 sm:h-8 bg-slate-200 dark:bg-slate-600"></div>
                  <button @click="handleAbandon" class="px-3 sm:px-4 py-2 sm:py-2.5 rounded-xl font-bold text-[10px] sm:text-xs text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-900/20 transition-colors">
                    Abandon
                  </button>
                </template>
              </div>
            </div>
          </div>
        </div>

        <!-- ═══ OTHER PHASES: Timer + Sidebar grid ═══ -->
        <div v-else-if="status.can_start_today" class="flex-1 w-full max-w-5xl mx-auto grid grid-cols-1 lg:grid-cols-12 gap-4 xl:gap-6 min-h-0 py-2 sm:py-4">

          <!-- TIMER CARD -->
          <div class="lg:col-span-8 bg-white dark:bg-slate-900 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-lg p-4 flex flex-col items-center justify-center relative overflow-hidden min-h-0">

            <!-- Glow -->
            <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-48 h-48 rounded-full blur-[60px] opacity-25 pointer-events-none"
              :class="isWorkCycle ? 'bg-indigo-500' : 'bg-emerald-500'"></div>

            <!-- Status badge -->
            <div class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 mb-3">
              <span class="relative flex h-2 w-2">
                <span class="animate-ping absolute inline-flex h-full w-full rounded-full opacity-75" :class="isWorkCycle ? 'bg-indigo-400' : 'bg-emerald-400'"></span>
                <span class="relative inline-flex rounded-full h-2 w-2" :class="isWorkCycle ? 'bg-indigo-500' : 'bg-emerald-500'"></span>
              </span>
              <span class="text-[10px] font-bold uppercase tracking-[0.15em]" :class="isWorkCycle ? 'text-indigo-600 dark:text-indigo-400' : 'text-emerald-600 dark:text-emerald-400'">
                {{ isWorkCycle ? 'Deep Focus' : 'Recovery Break' }}
              </span>
            </div>

            <!-- SVG Timer (compact) -->
            <div class="relative cursor-pointer" @click="togglePause">
              <div class="absolute inset-0 rounded-full blur-lg opacity-15" :class="isWorkCycle ? 'bg-indigo-500' : 'bg-emerald-500'"></div>
              <svg class="timer-svg transform -rotate-90" viewBox="0 0 200 200">
                <circle cx="100" cy="100" r="88" fill="none" stroke="currentColor" stroke-width="8" class="text-slate-100 dark:text-slate-800/50"/>
                <circle cx="100" cy="100" r="88" fill="none" stroke="url(#timerGrad)" stroke-width="8"
                  stroke-linecap="round" :stroke-dasharray="553" :stroke-dashoffset="dashOffset"
                  class="transition-all duration-1000 ease-linear"/>
                <defs>
                  <linearGradient id="timerGrad" x1="0%" y1="0%" x2="100%" y2="0%">
                    <stop offset="0%" :stop-color="isWorkCycle ? '#6366f1' : '#10b981'"/>
                    <stop offset="100%" :stop-color="isWorkCycle ? '#a855f7' : '#3b82f6'"/>
                  </linearGradient>
                </defs>
              </svg>
              <div class="absolute inset-0 flex flex-col items-center justify-center">
                <div class="timer-time text-slate-800 dark:text-white">{{ formattedTime }}</div>
                <div class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">{{ isWorkCycle ? 'Remaining' : 'Resting' }}</div>
              </div>
            </div>

            <!-- Cycle dots -->
            <div class="flex items-center gap-2 mt-3">
              <div v-for="i in status.target_cycles" :key="i"
                class="w-2 h-2 rounded-full transition-all duration-300"
                :class="i <= completedCycles ? 'bg-indigo-500 scale-125' : 'bg-slate-200 dark:bg-slate-700'"></div>
              <span class="ml-1 text-[10px] font-bold text-slate-400">{{ completedCycles }}/{{ status.target_cycles }}</span>
            </div>

            <!-- Controls -->
            <div class="flex items-center gap-3 mt-4 p-1.5 bg-slate-100 dark:bg-slate-800/50 rounded-xl border border-slate-200 dark:border-slate-700/50">
              <button v-if="!hasActiveSession" @click="handleStart"
                :disabled="status.requires_project && !selectedProjectId"
                class="px-6 py-2 bg-slate-900 dark:bg-white text-white dark:text-slate-900 rounded-lg font-bold text-xs tracking-wide hover:scale-105 active:scale-95 disabled:opacity-40 disabled:cursor-not-allowed transition-all flex items-center gap-1.5">
                <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
                Start
              </button>
              <template v-else>
                <button @click="togglePause"
                  class="w-10 h-10 flex items-center justify-center rounded-lg bg-white dark:bg-slate-700 text-slate-900 dark:text-white shadow hover:bg-slate-50 dark:hover:bg-slate-600 transition-all active:scale-95">
                  <svg v-if="isPaused" class="w-4 h-4 ml-0.5" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z"/></svg>
                  <svg v-else class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M6 19h4V5H6v14zm8-14v14h4V5h-4z"/></svg>
                </button>
                <div class="w-px h-6 bg-slate-300 dark:bg-slate-600"></div>
                <button @click="handleAbandon" class="px-3 py-2 rounded-lg font-bold text-[10px] text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-900/20 transition-colors">
                  Abandon
                </button>
              </template>
            </div>
          </div>

          <!-- SIDEBAR -->
          <div class="lg:col-span-4 flex flex-col gap-3 min-h-0">
            <!-- Stats -->
            <div class="grid grid-cols-2 gap-3 shrink-0">
              <div class="bg-white dark:bg-slate-900 p-3 rounded-xl border border-slate-200 dark:border-slate-800 shadow-sm">
                <p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-0.5">Work</p>
                <p class="text-xl font-black text-slate-800 dark:text-white">{{ status.required_work_duration }}<span class="text-xs text-slate-400 ml-0.5">m</span></p>
              </div>
              <div class="bg-white dark:bg-slate-900 p-3 rounded-xl border border-slate-200 dark:border-slate-800 shadow-sm">
                <p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mb-0.5">Break</p>
                <p class="text-xl font-black text-emerald-600 dark:text-emerald-400">{{ status.required_break_duration }}<span class="text-xs text-slate-400 ml-0.5">m</span></p>
              </div>
            </div>

            <!-- Project -->
            <div class="bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 p-4 shadow-md flex-1 flex flex-col min-h-0">
              <h3 class="text-xs font-bold text-slate-900 dark:text-white mb-3 flex items-center gap-1.5">
                <svg class="w-3.5 h-3.5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>
                Project
              </h3>

              <div v-if="status.requires_project && !hasActiveSession">
                <div class="bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg p-2.5 mb-2">
                  <p class="text-[10px] font-bold text-amber-700 dark:text-amber-400">⚠️ Required for Hard Work</p>
                </div>
                <select v-model="selectedProjectId" class="w-full bg-white dark:bg-slate-800 border-2 border-amber-300 dark:border-amber-700 rounded-lg px-3 py-2 text-xs font-bold outline-none focus:ring-2 focus:ring-amber-400 text-slate-900 dark:text-white">
                  <option :value="null" disabled>Select project...</option>
                  <option v-for="p in projectStore.projects" :key="p.id" :value="p.id">{{ p.name }}</option>
                </select>
              </div>
              <div v-else>
                <select v-model="selectedProjectId" class="w-full bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg px-3 py-2 text-xs font-medium outline-none focus:ring-2 focus:ring-indigo-400 text-slate-900 dark:text-white">
                  <option :value="null">None</option>
                  <option v-for="p in projectStore.projects" :key="p.id" :value="p.id">{{ p.name }}</option>
                </select>
              </div>

              <!-- Phase info -->
              <div class="mt-auto pt-3 border-t border-slate-100 dark:border-slate-800 space-y-1.5">
                <div class="flex justify-between text-xs"><span class="text-slate-500">Phase</span><span class="font-bold" :class="phaseStyle.text">{{ phaseName }}</span></div>
                <div class="flex justify-between text-xs"><span class="text-slate-500">Level</span><span class="font-bold text-slate-900 dark:text-white">{{ currentLevel }}</span></div>
                <div class="flex justify-between text-xs"><span class="text-slate-500">Cycles</span><span class="font-bold text-rose-500">{{ status.target_cycles }}×</span></div>
              </div>
            </div>
          </div>
        </div>

        <!-- ═══ REST SCREEN ═══ -->
        <div v-else class="flex-1 flex items-center justify-center">
          <div class="text-center max-w-md">
            <div class="w-20 h-20 bg-gradient-to-br from-emerald-400 to-teal-500 rounded-2xl mx-auto flex items-center justify-center text-4xl shadow-xl shadow-emerald-500/30 mb-4 animate-bounce">🎉</div>
            <h2 class="text-2xl font-black text-slate-900 dark:text-white mb-1">Day Complete!</h2>
            <p class="text-slate-500 dark:text-slate-400 text-sm mb-6">Level {{ Math.max(1, currentLevel - 1) }} done. Come back tomorrow.</p>
            <div class="grid grid-cols-3 gap-3">
              <div class="p-3 rounded-xl bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700">
                <div class="text-indigo-600 font-black text-lg">{{ status.required_work_duration }}m</div>
                <div class="text-[9px] uppercase font-bold text-slate-400 mt-0.5">Next Work</div>
              </div>
              <div class="p-3 rounded-xl bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700">
                <div class="text-emerald-600 font-black text-lg">{{ status.required_break_duration }}m</div>
                <div class="text-[9px] uppercase font-bold text-slate-400 mt-0.5">Next Break</div>
              </div>
              <div class="p-3 rounded-xl bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700">
                <div class="text-rose-500 font-black text-lg">{{ status.target_cycles }}×</div>
                <div class="text-[9px] uppercase font-bold text-slate-400 mt-0.5">Cycles</div>
              </div>
            </div>
          </div>
        </div>

      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, inject } from 'vue'
import { usePomodoroStore } from '../stores/pomodoro'
import { useProjectStore } from '../stores/projects'

const pomodoroStore = usePomodoroStore()
const projectStore = useProjectStore()
const showToast = inject('showToast')

// ─── Backend-driven data ───
const status = computed(() => pomodoroStore.status)
const currentLevel = computed(() => status.value?.current_level || 1)

// ─── Phase definitions ───
const phases = [
  { id: 'normal',   name: 'Normal',    emoji: '🌱', levels: 'Lv 1–9',   min: 1,  max: 9,  glowClass: 'border-indigo-500 text-indigo-500',  pulseRing: 'bg-indigo-500' },
  { id: 'focus',    name: 'Focus',     emoji: '🔥', levels: 'Lv 20–29', min: 20, max: 29, glowClass: 'border-rose-500 text-rose-500',    pulseRing: 'bg-rose-500' },
  { id: 'deepwork', name: 'Deep Work', emoji: '🧠', levels: 'Lv 30–39', min: 30, max: 39, glowClass: 'border-violet-500 text-violet-500',  pulseRing: 'bg-violet-500' },
  { id: 'hardwork', name: 'Hard Work', emoji: '⚡', levels: 'Lv 40–49', min: 40, max: 49, glowClass: 'border-amber-500 text-amber-500',   pulseRing: 'bg-amber-500' },
]

const currentPhaseId = computed(() => {
  const l = currentLevel.value
  if (l <= 9) return 'normal'
  if (l <= 29) return 'focus'
  if (l <= 39) return 'deepwork'
  return 'hardwork'
})

const currentPhase = computed(() => phases.find(p => p.id === currentPhaseId.value))
const phaseName = computed(() => currentPhase.value?.name || 'Normal')
const phaseMax = computed(() => currentPhase.value?.max || 9)
const isNormalPhase = computed(() => currentPhaseId.value === 'normal')

const phaseStyle = computed(() => ({
  normal:   { icon: '🌱', headerBg: 'bg-indigo-100 dark:bg-indigo-900/30', text: 'text-indigo-600 dark:text-indigo-400', barActive: 'bg-indigo-500', progressLine: 'bg-indigo-500', badgeClass: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300' },
  focus:    { icon: '🔥', headerBg: 'bg-rose-100 dark:bg-rose-900/30',     text: 'text-rose-600 dark:text-rose-400',     barActive: 'bg-rose-500',   progressLine: 'bg-rose-500',   badgeClass: 'bg-rose-100 text-rose-700 dark:bg-rose-900/40 dark:text-rose-300' },
  deepwork: { icon: '🧠', headerBg: 'bg-violet-100 dark:bg-violet-900/30', text: 'text-violet-600 dark:text-violet-400', barActive: 'bg-violet-500', progressLine: 'bg-violet-500', badgeClass: 'bg-violet-100 text-violet-700 dark:bg-violet-900/40 dark:text-violet-300' },
  hardwork: { icon: '⚡', headerBg: 'bg-amber-100 dark:bg-amber-900/30',   text: 'text-amber-600 dark:text-amber-400',   barActive: 'bg-amber-500',  progressLine: 'bg-amber-500',  badgeClass: 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300' },
}[currentPhaseId.value]))

const phaseRange = computed(() => {
  const p = currentPhase.value
  if (!p) return []
  return Array.from({ length: p.max - p.min + 1 }, (_, i) => i + p.min)
})

function isPastPhase(phaseId) {
  const order = ['normal', 'focus', 'deepwork', 'hardwork']
  return order.indexOf(phaseId) < order.indexOf(currentPhaseId.value)
}

const progressLineWidth = computed(() => {
  const order = ['normal', 'focus', 'deepwork', 'hardwork']
  const idx = order.indexOf(currentPhaseId.value)
  return `${(idx / (phases.length - 1)) * 100}%`
})

// ─── Timer State ───
const secondsLeft = ref(0)
const isWorkCycle = ref(true)
const isPaused = ref(false)
const timerRunning = ref(false)   // TRUE = session active & timer ticking
const completedCycles = ref(0)
const timerInterval = ref(null)
const selectedProjectId = ref(null)

// Whether to show Start vs Pause/Abandon
const hasActiveSession = computed(() => timerRunning.value || !!status.value?.active_session)

const formattedTime = computed(() => {
  const m = Math.floor(secondsLeft.value / 60)
  const s = secondsLeft.value % 60
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
})

const totalSeconds = computed(() => {
  if (!status.value) return 1
  return (isWorkCycle.value ? status.value.required_work_duration : status.value.required_break_duration) * 60
})

const dashOffset = computed(() => 553 - (secondsLeft.value / totalSeconds.value) * 553)

// ─── Timer Engine ───
function startLocalTimer() {
  clearInterval(timerInterval.value)
  timerRunning.value = true
  timerInterval.value = setInterval(tick, 1000)
}

function stopLocalTimer() {
  clearInterval(timerInterval.value)
  timerRunning.value = false
  isPaused.value = false
}

function tick() {
  if (isPaused.value) return
  if (secondsLeft.value > 0) {
    secondsLeft.value--
  } else {
    handleTimerEnd()
  }
}

async function handleTimerEnd() {
  clearInterval(timerInterval.value)

  if (isWorkCycle.value) {
    // Work cycle finished → log to backend (sets is_on_break=true)
    try {
      await pomodoroStore.logCycle()
      completedCycles.value = pomodoroStore.status?.active_session?.completed_cycles || completedCycles.value + 1

      if (completedCycles.value >= status.value.target_cycles) {
        const completeRes = await pomodoroStore.completeLevel()
        showToast(completeRes?.message || 'Level complete! 🎉', 'success')
        stopLocalTimer()
        syncFromBackend()
        return
      }

      // Backend now has is_on_break=true and new seconds_remaining
      isWorkCycle.value = false
      secondsLeft.value = pomodoroStore.status?.seconds_remaining ?? (status.value.required_break_duration * 60)
      showToast('Work done! Take a break ☕', 'success')
      startLocalTimer()
    } catch (e) {
      showToast('Failed to sync cycle', 'error')
      stopLocalTimer()
    }
  } else {
    // Break finished → tell backend to start next work cycle
    try {
      await pomodoroStore.resumeWork()
      isWorkCycle.value = true
      secondsLeft.value = pomodoroStore.status?.seconds_remaining ?? (status.value.required_work_duration * 60)
      showToast('Break over — focus up! 🎯', 'info')
      startLocalTimer()
    } catch (e) {
      showToast('Failed to resume work', 'error')
      stopLocalTimer()
    }
  }
}

// ─── User Actions ───
async function handleStart() {
  try {
    await pomodoroStore.startSession(selectedProjectId.value)
    syncFromBackend()
    showToast('Session started! 🎯', 'success')
  } catch (err) {
    const msg = err?.response?.data?.error || ''
    if (msg.includes('already in progress')) {
      // Another tab/device already started — just sync to the existing session
      await pomodoroStore.fetchStatus()
      syncFromBackend()
      showToast('Resuming existing session', 'info')
    } else {
      showToast(msg || 'Failed to start session', 'error')
    }
  }
}

async function togglePause() {
  if (!timerRunning.value) return
  try {
    if (isPaused.value) {
      await pomodoroStore.resumeSession()
      isPaused.value = false
      secondsLeft.value = pomodoroStore.status?.seconds_remaining ?? secondsLeft.value
    } else {
      await pomodoroStore.pauseSession()
      isPaused.value = true
    }
  } catch (e) {
    showToast('Failed to toggle pause', 'error')
  }
}

async function handleAbandon() {
  if (!confirm('Abandon this session? You can retry today.')) return
  stopLocalTimer()
  try {
    await pomodoroStore.abandonSession()
    showToast('Session abandoned', 'info')
  } catch (e) {
    showToast('Failed to abandon', 'error')
  }
  syncFromBackend()
}

// ─── Sync from backend (page load, after actions) ───
function syncFromBackend() {
  const s = pomodoroStore.status
  if (!s) return

  const sess = s.active_session
  if (sess) {
    completedCycles.value = sess.completed_cycles || 0
    isWorkCycle.value = !s.is_on_break
    isPaused.value = !!s.is_paused
    secondsLeft.value = s.seconds_remaining ?? 0

    timerRunning.value = true
    clearInterval(timerInterval.value)
    timerInterval.value = setInterval(tick, 1000)
  } else {
    stopLocalTimer()
    completedCycles.value = 0
    isWorkCycle.value = true
    secondsLeft.value = s.required_work_duration * 60
  }
}

// ─── WebSocket for real-time sync across tabs/devices ───
const ws = ref(null)
const wsReconnectTimer = ref(null)

function connectWebSocket() {
  const token = localStorage.getItem('jwt')
  if (!token) return

  const proto = location.protocol === 'https:' ? 'wss:' : 'ws:'
  const url = `${proto}//${location.host}/api/v1/ws/pomodoro?token=${encodeURIComponent(token)}`

  const socket = new WebSocket(url)

  socket.onopen = () => {
    console.log('[WS] Connected to pomodoro sync')
  }

  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      if (data.type === 'pomodoro_status') {
        // Update the store with the incoming real-time data
        pomodoroStore.status = {
          ...pomodoroStore.status,
          current_level: data.current_level,
          required_work_duration: data.required_work_duration,
          required_break_duration: data.required_break_duration,
          target_cycles: data.target_cycles,
          requires_project: data.requires_project,
          active_session: data.active_session,
          seconds_remaining: data.seconds_remaining,
          is_on_break: data.is_on_break,
          is_paused: data.is_paused,
        }
        syncFromBackend()
      }
    } catch (e) {
      console.error('[WS] Parse error:', e)
    }
  }

  socket.onclose = () => {
    console.log('[WS] Disconnected, reconnecting in 3s...')
    ws.value = null
    wsReconnectTimer.value = setTimeout(connectWebSocket, 3000)
  }

  socket.onerror = (err) => {
    console.error('[WS] Error:', err)
    socket.close()
  }

  ws.value = socket
}

// ─── Lifecycle ───
onMounted(async () => {
  await pomodoroStore.fetchStatus()
  await projectStore.fetchProjects()
  syncFromBackend()
  connectWebSocket()
})

onUnmounted(() => {
  clearInterval(timerInterval.value)
  clearTimeout(wsReconnectTimer.value)
  if (ws.value) ws.value.close()
})
</script>

<style scoped>
/* Fill the entire available viewport area from MainLayout */
.pomodoro-wrap {
  flex: 1;
  position: relative;
  height: 100%;
  min-height: calc(100vh - 4rem);
  overflow: hidden;
}

/* Timer ring - responsive to fill available space */
.timer-svg {
  width: min(52vh, 16rem);
  height: min(52vh, 16rem);
}
.timer-time {
  font-size: clamp(2.5rem, 8vh, 4.5rem);
  font-weight: 900;
  font-variant-numeric: tabular-nums;
  line-height: 1;
  letter-spacing: -0.03em;
}

/* Large timer for Normal phase */
.timer-svg-large {
  width: min(50vh, 18rem);
  height: min(50vh, 18rem);
}
.timer-time-large {
  font-size: clamp(2.5rem, 8vh, 4.5rem);
  font-weight: 900;
  font-variant-numeric: tabular-nums;
  line-height: 1;
  letter-spacing: -0.04em;
}

/* Entrance animation */
@keyframes fade-in-up {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>