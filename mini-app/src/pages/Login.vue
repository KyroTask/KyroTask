<template>
  <div class="min-h-screen flex items-center justify-center bg-slate-50 p-4">
    <div class="w-full max-w-md">
      <!-- Logo & Header -->
      <div class="text-center mb-10">
        <div
          class="inline-flex items-center justify-center w-28 h-28 mb-4 transform -rotate-2"
        >
          <img src="/logo.png" alt="KyroTask Logo" class="w-full h-full object-contain drop-shadow-xl" />
        </div>
        <h1 class="text-4xl font-black text-slate-900 tracking-tight">
          KyroTask
        </h1>
        <p class="text-slate-500 mt-2 font-medium">
          Your personal productivity hub
        </p>
      </div>

      <!-- Main Card -->
      <div
        class="bg-white rounded-[2.5rem] shadow-2xl shadow-slate-200/60 p-8 md:p-10 border border-slate-100"
      >
        <div v-if="loading" class="text-center py-10">
          <div
            class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-slate-100 border-t-indigo-600 mb-4"
          ></div>
          <p class="text-slate-600 font-bold">Verifying account...</p>
        </div>

        <div v-else-if="error" class="text-center py-6">
          <div
            class="w-16 h-16 bg-rose-50 text-rose-500 rounded-full flex items-center justify-center mx-auto mb-4"
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
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
          </div>
          <h2 class="text-xl font-bold text-slate-900 mb-2">
            Authentication Error
          </h2>
          <p class="text-slate-500 mb-8 text-sm leading-relaxed">{{ error }}</p>

          <div class="space-y-3">
            <button
              @click="handleLogin"
              class="w-full bg-indigo-600 text-white py-4 rounded-2xl font-bold hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100"
            >
              Try Again
            </button>

          </div>
        </div>

        <div v-else class="text-center">
          <div class="mb-8">
            <h2 class="text-2xl font-bold text-slate-900 mb-2">Welcome!</h2>
            <p class="text-slate-500 text-sm">
              Please sign in with Telegram to continue to your workspace.
            </p>
          </div>

          <div class="space-y-4">
            <!-- Firebase Login -->
            <input v-model="email" type="email" placeholder="Email" class="w-full px-4 py-3 rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500">
            <input v-model="password" type="password" placeholder="Password" class="w-full px-4 py-3 rounded-xl border border-slate-200 focus:outline-none focus:ring-2 focus:ring-indigo-500" @keyup.enter="handleEmailLogin">
            <button @click="handleEmailLogin" class="w-full bg-slate-900 text-white font-bold py-3 rounded-xl hover:bg-slate-800 transition shadow-lg shadow-slate-200">Sign in with Email</button>
            
            <div class="relative flex py-2 items-center">
              <div class="flex-grow border-t border-slate-200"></div>
              <span class="flex-shrink-0 mx-4 text-slate-400 text-xs font-bold tracking-widest uppercase">OR</span>
              <div class="flex-grow border-t border-slate-200"></div>
            </div>

            <button @click="handleGoogleLogin" class="w-full flex items-center justify-center gap-3 bg-white border-2 border-slate-100 text-slate-700 font-bold py-3 rounded-xl hover:bg-slate-50 hover:border-slate-200 transition">
              <svg viewBox="0 0 24 24" class="w-5 h-5"><path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/><path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/><path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/><path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/></svg>
              Sign in with Google
            </button>

            <div class="relative flex py-2 items-center">
              <div class="flex-grow border-t border-slate-200"></div>
              <span class="flex-shrink-0 mx-4 text-slate-400 text-xs font-bold tracking-widest uppercase">TELEGRAM</span>
              <div class="flex-grow border-t border-slate-200"></div>
            </div>

            <!-- Telegram Widget Container -->
            <div id="telegram-widget-container" class="flex justify-center my-4"></div>

            <div class="relative flex py-2 items-center">
              <div class="flex-grow border-t border-slate-200"></div>
              <span class="flex-shrink-0 mx-4 text-slate-400 text-xs font-bold tracking-widest uppercase">DEV</span>
              <div class="flex-grow border-t border-slate-200"></div>
            </div>

            <button @click="handleDevLogin" class="w-full flex items-center justify-center gap-3 bg-gradient-to-r from-indigo-600 to-violet-600 text-white font-bold py-3.5 rounded-xl hover:from-indigo-700 hover:to-violet-700 transition shadow-lg shadow-indigo-200">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" /></svg>
              Quick Dev Login
            </button>
          </div>

          <p
            class="mt-8 text-[10px] font-bold text-slate-400 uppercase tracking-widest"
          >
            🔒 Secure & Private
          </p>
        </div>
      </div>

      <!-- Footer Info -->
      <div class="mt-10 grid grid-cols-3 gap-4">
        <div class="text-center">
          <div class="text-xl mb-1">⚡</div>
          <p
            class="text-[10px] font-bold text-slate-500 uppercase tracking-tighter"
          >
            Fast
          </p>
        </div>
        <div class="text-center">
          <div class="text-xl mb-1">🛡️</div>
          <p
            class="text-[10px] font-bold text-slate-500 uppercase tracking-tighter"
          >
            Secure
          </p>
        </div>
        <div class="text-center">
          <div class="text-xl mb-1">☁️</div>
          <p
            class="text-[10px] font-bold text-slate-500 uppercase tracking-tighter"
          >
            Sync
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "../stores/auth";
import api from "../services/api";
import { loginWithGoogle, loginWithEmail, registerWithEmail } from "../utils/firebase";

const router = useRouter();
const authStore = useAuthStore();
const loading = ref(false);
const error = ref(null);
const email = ref("");
const password = ref("");

async function handleGoogleLogin() {
  loading.value = true;
  error.value = null;
  try {
    const result = await loginWithGoogle();
    const idToken = await result.user.getIdToken();
    await authStore.loginWithFirebase(idToken, false);
    router.push("/");
  } catch (err) {
    error.value = err.message || "Google Login failed";
    console.error(err);
  } finally {
    loading.value = false;
  }
}

async function handleEmailLogin() {
  if (!email.value || !password.value) {
    error.value = "Email and password are required";
    return;
  }
  loading.value = true;
  error.value = null;
  try {
    let result;
    try {
      result = await loginWithEmail(email.value, password.value);
    } catch (e) {
      if (e.code === 'auth/user-not-found' || e.code === 'auth/invalid-credential') {
        // Option to sign up
        result = await registerWithEmail(email.value, password.value);
      } else {
        throw e;
      }
    }
    const idToken = await result.user.getIdToken();
    await authStore.loginWithFirebase(idToken, false);
    router.push("/");
  } catch (err) {
    error.value = err.message || "Email Login/Signup failed";
    console.error(err);
  } finally {
    loading.value = false;
  }
}

async function handleLogin() {
  loading.value = true;
  error.value = null;

  try {
    if (window.Telegram?.WebApp?.initData) {
      await authStore.loginWithTelegram(window.Telegram.WebApp.initData);
      router.push("/");
    } else {
      // If not in Telegram WebApp, show error but widget is also available
      throw new Error(
        "No Telegram data available. Please open this app from Telegram or use the login widget."
      );
    }
  } catch (err) {
    error.value =
      err.response?.data?.error || err.message || "Authentication failed";
    console.error("Login error:", err);
  } finally {
    loading.value = false;
  }
}

async function handleDevLogin() {
  loading.value = true;
  error.value = null;
  try {
    await authStore.loginDev();
    router.push("/");
  } catch (err) {
    error.value = err.message || "Dev Login failed";
    console.error(err);
  } finally {
    loading.value = false;
  }
}


// Global callback for Telegram Widget
window.onTelegramAuth = async (user) => {
  loading.value = true;
  error.value = null;
  try {
    await authStore.loginWithWidget(user);
    router.push("/");
  } catch (err) {
    error.value = "Widget login failed";
    console.error(err);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  if (window.Telegram?.WebApp) {
    window.Telegram.WebApp.ready();
    if (window.Telegram.WebApp.initData) {
      handleLogin();
    }
  } else if (!authStore.isAuthenticated) {
    // Automatically enter demo mode if not in Telegram and not authenticated
    // enterPreviewMode(); 
    // Commented out auto-demo to allow widget usage
  }

  // Inject Telegram Widget script dynamically
  const script = document.createElement("script");
  script.async = true;
  script.src = "https://telegram.org/js/telegram-widget.js?22";
  script.setAttribute("data-telegram-login", "KyroTaskBot");
  script.setAttribute("data-size", "large");
  script.setAttribute("data-onauth", "onTelegramAuth(user)");
  script.setAttribute("data-request-access", "write");
  document.getElementById("telegram-widget-container").appendChild(script);
});
</script>
