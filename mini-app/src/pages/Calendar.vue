<template>
  <div class="space-y-8 pb-24 md:pb-12">
    <!-- Page Header -->
    <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 mb-4">
      <div>
        <h1 class="text-2xl md:text-3xl font-black tracking-tight text-slate-900 dark:text-white">Your Timeline</h1>
        <p class="text-slate-500 dark:text-slate-400 text-sm font-medium mt-1">Master your time and stay ahead of your commitments.</p>
      </div>
      <div class="flex items-center gap-3">
        <div class="flex items-center space-x-2 text-sm">
          <div class="bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 px-3 py-1.5 rounded-lg flex items-center gap-2 shadow-sm">
            <span class="font-black text-indigo-600">{{ eventsThisWeek }}</span>
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Week</span>
          </div>
          <div class="bg-white dark:bg-slate-900 border border-slate-100 dark:border-slate-800 px-3 py-1.5 rounded-lg flex items-center gap-2 shadow-sm">
            <span class="font-black text-emerald-600">{{ eventsToday }}</span>
            <span class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Today</span>
          </div>
        </div>
        <div class="flex items-center bg-slate-100 dark:bg-slate-800 rounded-xl p-1">
          <button
            v-for="view in ['Month', 'Week', 'Day']"
            :key="view"
            @click="activeView = view"
            class="px-3 py-1.5 rounded-lg text-xs font-bold transition-all"
            :class="[
              activeView === view
                ? 'bg-white dark:bg-slate-700 text-slate-900 dark:text-white shadow-sm'
                : 'text-slate-400 hover:text-slate-600 dark:hover:text-slate-200',
            ]"
          >
            {{ view }}
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="calendarStore.loading" class="py-20 text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-4 border-slate-100 border-t-indigo-600 mx-auto"></div>
    </div>

    <!-- Calendar Interface -->
    <div v-else class="bg-white dark:bg-slate-900 rounded-[2.5rem] md:rounded-[3.5rem] border border-slate-100 dark:border-slate-800 shadow-sm overflow-hidden animate-in fade-in slide-in-from-bottom-8 duration-1000">
      <!-- Navigation Header -->
      <div class="p-4 md:p-12 border-b border-slate-50 dark:border-slate-800 flex flex-col md:flex-row items-center justify-between gap-4">
        <div class="flex items-center gap-6">
          <h2 class="text-3xl md:text-4xl font-black text-slate-900 dark:text-white tracking-tighter">
            {{ currentMonthName }} <span class="text-indigo-600 dark:text-indigo-400">{{ currentYear }}</span>
          </h2>
          <div class="flex items-center bg-slate-100 dark:bg-slate-800 rounded-2xl p-1">
            <button @click="prevMonth" class="p-2 hover:bg-white dark:hover:bg-slate-700 hover:shadow-sm rounded-xl transition-all text-slate-600 dark:text-slate-400">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M15 19l-7-7 7-7" /></svg>
            </button>
            <button @click="nextMonth" class="p-2 hover:bg-white dark:hover:bg-slate-700 hover:shadow-sm rounded-xl transition-all text-slate-600 dark:text-slate-400">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 5l7 7-7 7" /></svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Month View Grid -->
      <div v-if="activeView === 'Month'" class="overflow-x-auto no-scrollbar">
        <div class="w-full md:min-w-[800px]">
          <div class="grid grid-cols-7 bg-slate-50/50 dark:bg-slate-800/50 border-b border-slate-100 dark:border-slate-800">
            <div v-for="day in ['S', 'M', 'T', 'W', 'T', 'F', 'S']" :key="day" class="py-2 md:py-6 text-center text-[10px] font-black text-slate-400 dark:text-slate-500 uppercase tracking-[0.3em]">
              {{ day }}
            </div>
          </div>

          <div class="grid grid-cols-7">
            <div
              v-for="(day, index) in calendarDays"
              :key="index"
              @dblclick="handleDayDoubleClick(day.date)"
              class="min-h-[60px] md:min-h-[180px] p-1 md:p-4 border-r border-b border-slate-50 dark:border-slate-800 last:border-r-0 group transition-all duration-500 hover:bg-indigo-50/30 dark:hover:bg-indigo-900/10 cursor-pointer"
              :class="{ 'bg-slate-50/20 dark:bg-slate-800/10': !day.isCurrentMonth }"
            >
              <div class="flex justify-between items-start mb-1 md:mb-4">
                <span
                  class="w-6 h-6 md:w-10 md:h-10 flex items-center justify-center rounded-lg md:rounded-2xl text-[10px] md:text-base font-black transition-all duration-500"
                  :class="[
                    day.isToday
                      ? 'bg-indigo-600 text-white shadow-2xl shadow-indigo-200 dark:shadow-none scale-110 ring-4 ring-indigo-50 dark:ring-indigo-900/50'
                      : day.isCurrentMonth
                        ? 'text-slate-900 dark:text-slate-100'
                        : 'text-slate-300 dark:text-slate-700',
                  ]"
                >
                  {{ day.date.getDate() }}
                </span>
              </div>

              <div class="space-y-1 md:space-y-1.5">
                <div
                  v-for="event in getEventsForDay(day.date).slice(0, 3)"
                  :key="event.id + '-' + event.type + '-' + event.date"
                  class="px-2 py-1 rounded-lg text-[8px] md:text-[10px] font-bold truncate transition-all duration-300 hover:scale-[1.02] border border-transparent"
                  :class="eventClass(event)"
                >
                  {{ event.title }}
                </div>
                <div v-if="getEventsForDay(day.date).length > 3" class="text-[9px] font-black text-slate-400 uppercase tracking-widest pl-2">
                  + {{ getEventsForDay(day.date).length - 3 }} More
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Week View -->
      <div v-else-if="activeView === 'Week'" class="divide-y divide-slate-50 dark:divide-slate-800">
        <div v-for="day in weekDays" :key="day.date.toISOString()" class="p-8 md:p-12 flex flex-col md:flex-row items-start gap-8 md:gap-16 hover:bg-slate-50/50 dark:hover:bg-slate-800/30 transition-all duration-500 group">
          <div class="text-center md:text-left min-w-[120px]">
            <p class="text-[10px] font-black text-slate-400 dark:text-slate-500 uppercase tracking-[0.2em] mb-2">{{ day.label }}</p>
            <div
              class="w-16 h-16 flex items-center justify-center rounded-[1.5rem] text-3xl font-black transition-all duration-500"
              :class="[
                day.isToday
                  ? 'bg-indigo-600 text-white shadow-2xl shadow-indigo-200 dark:shadow-none scale-110'
                  : 'text-slate-900 dark:text-slate-100 bg-slate-100 dark:bg-slate-800',
              ]"
            >
              {{ day.date.getDate() }}
            </div>
          </div>
          
          <div class="flex-1 w-full space-y-4">
            <div v-if="getEventsForDay(day.date).length === 0" class="py-8 text-slate-400 dark:text-slate-600 text-lg font-medium italic flex items-center gap-3">
              <svg class="w-6 h-6 opacity-20" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4"></path></svg>
              No commitments scheduled
            </div>
            <div v-for="event in getEventsForDay(day.date)" :key="event.id + '-' + event.type + '-' + event.date" class="bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700 p-6 rounded-[2rem] shadow-sm flex items-center justify-between group/task hover:shadow-xl hover:shadow-indigo-100/50 dark:hover:shadow-none transition-all duration-500">
              <div class="flex items-center gap-6">
                <div class="w-3 h-3 rounded-full shadow-sm" :class="eventDotClass(event)"></div>
                <div>
                  <h4 class="text-lg font-black text-slate-900 dark:text-white tracking-tight" :class="{ 'line-through text-slate-400 dark:text-slate-600': event.status === 'completed' }">{{ event.title }}</h4>
                  <p class="text-xs text-slate-500 font-medium">{{ event.project || eventTypeLabel(event.type) }}</p>
                </div>
              </div>
              <span class="px-4 py-1.5 rounded-xl text-[10px] font-black uppercase tracking-widest" :class="event.status === 'completed' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-slate-100 dark:bg-slate-800 text-slate-400'">{{ event.status || event.type }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Day View -->
      <div v-else-if="activeView === 'Day'" class="p-8 md:p-20">
        <div v-if="getEventsForDay(currentDate).length === 0" class="text-center space-y-8">
          <div class="w-32 h-32 bg-slate-50 dark:bg-slate-800 rounded-[2.5rem] flex items-center justify-center mx-auto shadow-inner">
            <svg class="w-16 h-16 text-slate-200 dark:text-slate-700" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
          </div>
          <div class="space-y-2">
            <h3 class="text-3xl font-black text-slate-900 dark:text-white tracking-tight">Clear Horizon</h3>
            <p class="text-slate-500 dark:text-slate-400 text-lg">No tasks scheduled for today. Time to recharge.</p>
          </div>
        </div>
        <div v-else class="space-y-6 max-w-4xl mx-auto">
          <div v-for="event in getEventsForDay(currentDate)" :key="event.id + '-' + event.type + '-' + event.date" class="bg-slate-50 dark:bg-slate-800 p-8 rounded-[2.5rem] flex items-center justify-between group hover:bg-white dark:hover:bg-slate-700 hover:shadow-2xl hover:shadow-indigo-100/50 dark:hover:shadow-none transition-all duration-500 border border-transparent hover:border-indigo-100 dark:hover:border-indigo-900/50">
            <div class="flex items-center gap-8">
              <div class="w-16 h-16 rounded-[1.5rem] flex items-center justify-center shadow-lg group-hover:scale-110 transition-transform duration-500" :class="eventIconClass(event)">
                <svg v-if="event.type === 'habit' || event.type === 'habit_log'" class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" /></svg>
                <svg v-else-if="event.type === 'milestone'" class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M3 21v-4m0 0V5a2 2 0 012-2h6.5l1 1H21l-3 6 3 6h-8.5l-1-1H5a2 2 0 00-2 2zm9-13.5V9" /></svg>
                <svg v-else class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" /></svg>
              </div>
              <div>
                <h4 class="text-2xl font-black text-slate-900 dark:text-white tracking-tight" :class="{ 'line-through text-slate-400 dark:text-slate-600': event.status === 'completed' }">{{ event.title }}</h4>
                <p class="text-sm text-slate-500 font-medium">{{ event.project || eventTypeLabel(event.type) }}</p>
              </div>
            </div>
            <div class="flex items-center gap-4">
              <span class="px-6 py-2 rounded-2xl text-[10px] font-black uppercase tracking-widest" :class="event.status === 'completed' ? 'bg-emerald-500/10 text-emerald-500' : 'bg-slate-200 dark:bg-slate-700 text-slate-600 dark:text-slate-400'">{{ event.status || event.type }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import { useCalendarStore } from "../stores/calendar";

const calendarStore = useCalendarStore();
const activeView = ref("Month");
const currentDate = ref(new Date());

// Compute the date range we need to fetch based on view & currentDate
const fetchRange = computed(() => {
  const d = currentDate.value;
  if (activeView.value === "Month") {
    // Fetch the entire month + padding days
    const firstDay = new Date(d.getFullYear(), d.getMonth(), 1);
    const lastDay = new Date(d.getFullYear(), d.getMonth() + 1, 0);
    // Include padding from prev/next month
    const start = new Date(firstDay);
    start.setDate(start.getDate() - firstDay.getDay());
    const end = new Date(lastDay);
    end.setDate(end.getDate() + (6 - lastDay.getDay()));
    return { start, end };
  } else if (activeView.value === "Week") {
    const startOfWeek = new Date(d);
    startOfWeek.setDate(d.getDate() - d.getDay());
    const endOfWeek = new Date(startOfWeek);
    endOfWeek.setDate(startOfWeek.getDate() + 6);
    return { start: startOfWeek, end: endOfWeek };
  } else {
    // Day view
    return { start: new Date(d), end: new Date(d) };
  }
});

// Watch for changes and refetch
watch(
  fetchRange,
  (range) => {
    calendarStore.fetchEvents(range.start, range.end);
  },
  { immediate: true }
);

const eventsToday = computed(() => {
  const todayStr = new Date().toISOString().split("T")[0];
  return calendarStore.events.filter((e) => {
    return new Date(e.date).toISOString().split("T")[0] === todayStr;
  }).length;
});

const eventsThisWeek = computed(() => {
  const today = new Date();
  const startOfWeek = new Date(today);
  startOfWeek.setDate(today.getDate() - today.getDay());
  startOfWeek.setHours(0, 0, 0, 0);

  const endOfWeek = new Date(startOfWeek);
  endOfWeek.setDate(startOfWeek.getDate() + 6);
  endOfWeek.setHours(23, 59, 59, 999);

  return calendarStore.events.filter((e) => {
    const d = new Date(e.date);
    return d >= startOfWeek && d <= endOfWeek;
  }).length;
});

const currentMonth = computed(() => currentDate.value.getMonth());
const currentYear = computed(() => currentDate.value.getFullYear());
const currentMonthName = computed(() => {
  if (activeView.value === "Day") {
    return new Intl.DateTimeFormat("en-US", {
      month: "long",
      day: "numeric",
    }).format(currentDate.value);
  }
  return new Intl.DateTimeFormat("en-US", { month: "long" }).format(
    currentDate.value
  );
});

const calendarDays = computed(() => {
  const days = [];
  const firstDayOfMonth = new Date(currentYear.value, currentMonth.value, 1);
  const lastDayOfMonth = new Date(currentYear.value, currentMonth.value + 1, 0);

  // Padding from previous month
  const prevMonthLastDay = new Date(currentYear.value, currentMonth.value, 0);
  for (let i = firstDayOfMonth.getDay(); i > 0; i--) {
    days.push({
      date: new Date(
        currentYear.value,
        currentMonth.value - 1,
        prevMonthLastDay.getDate() - i + 1
      ),
      isCurrentMonth: false,
      isToday: false,
    });
  }

  // Current month days
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  for (let i = 1; i <= lastDayOfMonth.getDate(); i++) {
    const date = new Date(currentYear.value, currentMonth.value, i);
    days.push({
      date,
      isCurrentMonth: true,
      isToday: date.getTime() === today.getTime(),
    });
  }

  // Padding for next month
  const remainingDays = 42 - days.length;
  for (let i = 1; i <= remainingDays; i++) {
    days.push({
      date: new Date(currentYear.value, currentMonth.value + 1, i),
      isCurrentMonth: false,
      isToday: false,
    });
  }

  return days;
});

const weekDays = computed(() => {
  const days = [];
  const startOfWeek = new Date(currentDate.value);
  startOfWeek.setDate(currentDate.value.getDate() - currentDate.value.getDay());

  const today = new Date();
  today.setHours(0, 0, 0, 0);

  for (let i = 0; i < 7; i++) {
    const date = new Date(startOfWeek);
    date.setDate(startOfWeek.getDate() + i);
    date.setHours(0, 0, 0, 0);
    days.push({
      date,
      label: date.toLocaleDateString("en-US", { weekday: "short" }),
      isToday: date.getTime() === today.getTime(),
    });
  }
  return days;
});

function getEventsForDay(date) {
  return calendarStore.getEventsForDay(date);
}

function eventClass(event) {
  if (event.status === "completed") {
    return "bg-emerald-50/50 dark:bg-emerald-900/10 text-emerald-600/50 dark:text-emerald-400/50 line-through";
  }
  if (event.type === "habit" || event.type === "habit_log") {
    return "bg-rose-50 dark:bg-rose-900/20 text-rose-600 dark:text-rose-400 border-rose-100 dark:border-rose-900/50";
  }
  if (event.type === "milestone") {
    return "bg-violet-50 dark:bg-violet-900/20 text-violet-600 dark:text-violet-400 border-violet-100 dark:border-violet-900/50";
  }
  if (event.priority === "high") {
    return "bg-rose-50 dark:bg-rose-900/20 text-rose-600 dark:text-rose-400 border-rose-100 dark:border-rose-900/50";
  }
  return "bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 dark:text-indigo-400 border-indigo-100 dark:border-indigo-900/50";
}

function eventDotClass(event) {
  if (event.type === "habit" || event.type === "habit_log") return "bg-rose-500 animate-pulse";
  if (event.type === "milestone") return "bg-violet-500";
  if (event.priority === "high") return "bg-rose-500 animate-pulse";
  return "bg-indigo-500";
}

function eventIconClass(event) {
  if (event.type === "habit" || event.type === "habit_log") return "bg-rose-500 text-white shadow-rose-200";
  if (event.type === "milestone") return "bg-violet-600 text-white shadow-violet-200";
  if (event.priority === "high") return "bg-rose-500 text-white shadow-rose-200";
  return "bg-indigo-600 text-white shadow-indigo-200";
}

function eventTypeLabel(type) {
  switch (type) {
    case "task": return "Task";
    case "habit": return "Habit";
    case "habit_log": return "Habit Log";
    case "milestone": return "Milestone";
    default: return "Personal";
  }
}

function prevMonth() {
  if (activeView.value === "Month") {
    currentDate.value = new Date(currentYear.value, currentMonth.value - 1, 1);
  } else if (activeView.value === "Week") {
    const d = new Date(currentDate.value);
    d.setDate(d.getDate() - 7);
    currentDate.value = d;
  } else {
    const d = new Date(currentDate.value);
    d.setDate(d.getDate() - 1);
    currentDate.value = d;
  }
}

function nextMonth() {
  if (activeView.value === "Month") {
    currentDate.value = new Date(currentYear.value, currentMonth.value + 1, 1);
  } else if (activeView.value === "Week") {
    const d = new Date(currentDate.value);
    d.setDate(d.getDate() + 7);
    currentDate.value = d;
  } else {
    const d = new Date(currentDate.value);
    d.setDate(d.getDate() + 1);
    currentDate.value = d;
  }
}

function goToToday() {
  currentDate.value = new Date();
}

function handleDayDoubleClick(date) {
  currentDate.value = new Date(date);
  activeView.value = "Day";
}
</script>
