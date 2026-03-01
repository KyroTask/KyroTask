<template>
  <div class="space-y-6 animate-fade-in pb-20 md:pb-6">
    <!-- Header -->
    <header class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-2">
      <div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white transition-colors">Settings</h1>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-1 transition-colors">
          Manage your profile, preferences, and connected accounts.
        </p>
      </div>
    </header>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      
      <!-- Left Column: Profile & Integrations -->
      <div class="lg:col-span-5 space-y-6">
        
        <!-- Profile Card -->
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 overflow-hidden transition-all hover:shadow-md relative group">
          <!-- Decorative Background -->
          <div class="h-24 bg-gradient-to-r from-primary-500 to-indigo-500 relative">
             <div class="absolute inset-0 bg-white/20 dark:bg-black/20 mix-blend-overlay"></div>
          </div>
          
          <div class="px-6 pb-6 pt-0 relative">
            <div class="flex justify-between items-end -mt-10 mb-4">
              <div class="w-20 h-20 rounded-2xl bg-white dark:bg-gray-800 p-1 shadow-md">
                <div class="w-full h-full rounded-xl bg-gray-100 dark:bg-gray-700 overflow-hidden border border-gray-100 dark:border-gray-600">
                  <img v-if="user?.photo_url" :src="user.photo_url" alt="Profile" class="w-full h-full object-cover" />
                  <div v-else class="w-full h-full flex items-center justify-center text-3xl font-black text-gray-400 dark:text-gray-500 bg-gray-50 dark:bg-gray-800">
                    {{ user?.first_name?.charAt(0) || "U" }}
                  </div>
                </div>
              </div>
              <div class="flex gap-2">
                <span v-if="user?.is_premium" class="px-2.5 py-1 bg-gradient-to-r from-amber-200 to-yellow-400 text-yellow-900 text-[10px] font-black rounded-lg uppercase tracking-widest shadow-sm">Pro</span>
                <span v-if="user?.is_verified" class="px-2.5 py-1 bg-blue-100 dark:bg-blue-900/40 text-blue-600 dark:text-blue-400 border border-blue-200 dark:border-blue-800/50 text-[10px] font-black rounded-lg uppercase tracking-widest">Verified</span>
              </div>
            </div>

            <div>
              <h2 class="text-xl font-bold text-gray-900 dark:text-white">{{ user?.first_name }} {{ user?.last_name }}</h2>
              <p class="text-sm text-gray-500 dark:text-gray-400 font-medium">@{{ user?.username || "explorer" }}</p>
            </div>
          </div>
        </div>

        <!-- Integrations Section -->
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6 transition-all hover:shadow-md">
           <div class="flex items-center gap-3 mb-6">
            <div class="w-10 h-10 rounded-full bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center text-blue-600 dark:text-blue-400">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path></svg>
            </div>
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Connected Apps</h2>
          </div>

          <div class="space-y-3">
            <!-- Telegram -->
            <button class="w-full flex items-center justify-between p-3.5 rounded-xl border border-gray-100 dark:border-gray-700 hover:border-blue-200 dark:hover:border-blue-800/50 hover:bg-blue-50/50 dark:hover:bg-blue-900/10 transition-all duration-300 group">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-[#0088cc] text-white rounded-xl flex items-center justify-center shadow-md shadow-blue-200 dark:shadow-none group-hover:scale-110 transition-transform overflow-hidden">
                  <img v-if="(user?.telegram_id !== 0 && user?.telegram_id != null) && user?.photo_url" :src="user.photo_url" alt="Telegram Profile" class="w-full h-full object-cover" />
                  <svg v-else class="w-5 h-5 ml-[-2px]" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm4.64 6.8c-.15 1.58-.8 5.42-1.13 7.19-.14.75-.42 1-.68 1.03-.58.05-1.02-.38-1.58-.75-.88-.58-1.38-.94-2.23-1.5-.99-.65-.35-1.01.22-1.59.15-.15 2.71-2.48 2.76-2.69.01-.03.01-.14-.07-.2-.08-.06-.19-.04-.27-.02-.12.02-1.96 1.25-5.54 3.67-.52.36-.99.53-1.42.52-.47-.01-1.37-.26-2.03-.48-.82-.27-1.47-.42-1.42-.88.03-.24.37-.48 1.02-.73 3.99-1.73 6.64-2.87 7.95-3.42 3.78-1.58 4.57-1.85 5.08-1.86.11 0 .37.03.53.17.14.11.18.26.2.37.01.06.03.24.01.37z"/></svg>
                </div>
                <div class="text-left">
                  <h4 class="text-sm font-bold text-gray-900 dark:text-white">Telegram Bot</h4>
                  <p class="text-xs text-gray-500 dark:text-gray-400 font-medium">Manage tasks remotely</p>
                </div>
              </div>
              <div class="px-3 py-1 rounded-lg text-xs font-bold transition-colors" :class="user?.telegram_id !== 0 && user?.telegram_id != null ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400' : 'bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-400'">
                {{ user?.telegram_id !== 0 && user?.telegram_id != null  ? 'Connected' : 'Connect' }}
              </div>
            </button>

            <!-- Google -->
            <button @click="connectGoogleAccount" class="w-full flex items-center justify-between p-3.5 rounded-xl border border-gray-100 dark:border-gray-700 hover:border-gray-200 dark:hover:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-all duration-300 group">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-white text-white rounded-xl flex items-center justify-center shadow-md shadow-gray-200 dark:shadow-none group-hover:scale-110 transition-transform overflow-hidden">
                  <img v-if="user?.firebase_uid && user?.photo_url" :src="user.photo_url" alt="Google Profile" class="w-full h-full object-cover" />
                  <svg v-else viewBox="0 0 24 24" class="w-5 h-5"><path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/><path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/><path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/><path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/></svg>
                </div>
                <div class="text-left">
                  <h4 class="text-sm font-bold text-gray-900 dark:text-white">Google Account</h4>
                  <p class="text-xs text-gray-500 dark:text-gray-400 font-medium">Backup and sync</p>
                </div>
              </div>
              <div class="px-3 py-1 rounded-lg text-xs font-bold transition-colors" :class="user?.firebase_uid ? 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-400' : 'bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-400'">
                {{ user?.firebase_uid ? 'Connected' : 'Connect' }}
              </div>
            </button>
          </div>
        </div>
      </div>

      <!-- Right Column: Preferences & Danger Zone -->
      <div class="lg:col-span-7 space-y-6">
        
        <!-- Preferences Card -->
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-700 p-6 transition-all hover:shadow-md">
          <div class="flex items-center gap-3 mb-6">
            <div class="w-10 h-10 rounded-full bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center text-primary-600 dark:text-primary-400">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"></path></svg>
            </div>
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">Preferences</h2>
          </div>

          <div class="space-y-2">
            <!-- Iterate Settings -->
            <div v-for="setting in settings" :key="setting.id" class="flex items-center justify-between p-3.5 rounded-xl hover:bg-gray-50 dark:hover:bg-gray-700/30 transition-colors group">
              <div class="flex items-center gap-4">
                <div :class="setting.bg" class="w-10 h-10 rounded-xl flex items-center justify-center">
                  <component :is="setting.icon" class="w-5 h-5" :class="setting.color" />
                </div>
                <div>
                  <h4 class="text-sm font-bold text-gray-900 dark:text-white">{{ setting.name }}</h4>
                  <p class="text-xs text-gray-500 dark:text-gray-400 font-medium">{{ setting.description }}</p>
                </div>
              </div>
              
              <!-- Dark Mode Switcher -->
              <div v-if="setting.name === 'Dark Mode'" class="flex items-center bg-gray-100 dark:bg-gray-900 p-1 rounded-lg border border-gray-200 dark:border-gray-700">
                <button v-for="theme in ['light', 'dark', 'system']" :key="theme" @click="authStore.setTheme(theme)" class="px-2.5 py-1 rounded-md text-[10px] font-bold uppercase tracking-wider transition-all" :class="authStore.theme === theme ? 'bg-white dark:bg-gray-700 text-primary-600 dark:text-primary-400 shadow-sm' : 'text-gray-500 hover:text-gray-700 dark:hover:text-gray-300'">
                  {{ theme }}
                </button>
              </div>
              
              <!-- Standard Toggle -->
              <button v-else @click="setting.enabled = !setting.enabled" class="w-12 h-6 rounded-full transition-colors relative" :class="setting.enabled ? 'bg-primary-500' : 'bg-gray-200 dark:bg-gray-700 border border-gray-300 dark:border-gray-600'">
                <div class="absolute top-[2px] left-[2px] w-5 h-5 bg-white rounded-full shadow-sm transition-transform duration-300" :class="setting.enabled ? 'translate-x-6' : 'translate-x-0'"></div>
              </button>
            </div>
          </div>
        </div>

        <!-- Danger Zone -->
        <div class="bg-red-50 dark:bg-red-900/10 rounded-2xl shadow-sm border border-red-100 dark:border-red-900/30 p-6 transition-all hover:shadow-md">
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div>
               <h3 class="text-base font-bold text-red-800 dark:text-red-400">Sign Out</h3>
               <p class="text-xs text-red-600/70 dark:text-red-400/70 mt-1 font-medium">Clear session data and return to login.</p>
            </div>
            <button @click="handleLogout" class="px-5 py-2.5 bg-white dark:bg-red-900/40 text-red-600 dark:text-red-300 border border-red-200 dark:border-red-800/50 rounded-xl font-bold text-sm tracking-wide hover:shadow-sm hover:bg-red-50 dark:hover:bg-red-900/60 active:scale-95 transition-all flex items-center justify-center gap-2 whitespace-nowrap">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"></path></svg>
              <span>Sign Out</span>
            </button>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, defineComponent, h } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";
import { loginWithGoogle } from "../utils/firebase";

const router = useRouter();
const authStore = useAuthStore();
const user = computed(() => authStore.user);

const connectGoogleAccount = async () => {
  if (user.value?.firebase_uid) return;
  try {
    const result = await loginWithGoogle();
    const idToken = await result.user.getIdToken();
    await authStore.loginWithFirebase(idToken, true);
    alert("Google Account connected successfully!");
  } catch (err) {
    console.error(err);
    alert("Failed to connect Google Account: " + err.message);
  }
};

const IconBell = defineComponent({
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
          d: "M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9",
        }),
      ]
    ),
});

const IconMoon = defineComponent({
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
          d: "M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z",
        }),
      ]
    ),
});

const IconLock = defineComponent({
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
          d: "M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z",
        }),
      ]
    ),
});

const settings = ref([
  {
    id: 1,
    name: "Push Notifications",
    description: "Receive alerts for upcoming tasks and goal progress.",
    icon: IconBell,
    bg: "bg-indigo-50",
    color: "text-indigo-600",
    enabled: true,
  },
  {
    id: 2,
    name: "Dark Mode",
    description: "Switch to a darker theme for better night viewing.",
    icon: IconMoon,
    bg: "bg-slate-900",
    color: "text-white",
    enabled: false,
  },
  {
    id: 3,
    name: "Privacy Mode",
    description: "Hide sensitive task details in the dashboard view.",
    icon: IconLock,
    bg: "bg-rose-50",
    color: "text-rose-600",
    enabled: false,
  },
]);

const handleLogout = () => {
  authStore.logout();
  router.push("/login");
};
</script>
