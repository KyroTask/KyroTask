<template>
  <TransitionGroup
    tag="div"
    class="fixed bottom-24 left-1/2 -translate-x-1/2 z-[100] flex flex-col items-center space-y-2 pointer-events-none"
  >
    <div
      v-for="toast in toasts"
      :key="toast.id"
      class="px-6 py-3 rounded-2xl shadow-2xl flex items-center space-x-3 pointer-events-auto animate-in fade-in slide-in-from-bottom-4 duration-300"
      :class="[
        toast.type === 'success'
          ? 'bg-emerald-600 text-white'
          : toast.type === 'error'
            ? 'bg-rose-600 text-white'
            : 'bg-slate-900 text-white',
      ]"
    >
      <svg
        v-if="toast.type === 'success'"
        class="w-5 h-5"
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
      <svg
        v-else-if="toast.type === 'error'"
        class="w-5 h-5"
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
      <span class="text-sm font-bold">{{ toast.message }}</span>
    </div>
  </TransitionGroup>
</template>

<script setup>
import { ref } from "vue";

const toasts = ref([]);
let nextId = 0;

function addToast(message, type = "info", duration = 3000) {
  const id = nextId++;
  toasts.value.push({ id, message, type });
  setTimeout(() => {
    toasts.value = toasts.value.filter((t) => t.id !== id);
  }, duration);
}

defineExpose({ addToast });
</script>
