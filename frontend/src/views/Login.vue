<template>
  <div class="login-page">
    <div class="login-card">
      <h1>🎬 Cinema Booking</h1>
      <p>เข้าสู่ระบบเพื่อจองตั๋วหนัง</p>
      <div id="google-signin-btn"></div>
      <p v-if="error" class="error">{{ error }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router    = useRouter()
const authStore = useAuthStore()
const error     = ref('')

onMounted(() => {
  window.google?.accounts.id.initialize({
    client_id: import.meta.env.VITE_GOOGLE_CLIENT_ID,
    callback: handleGoogleResponse
  })
  window.google?.accounts.id.renderButton(
    document.getElementById('google-signin-btn'),
    { theme: 'outline', size: 'large' }
  )
})

async function handleGoogleResponse(response) {
  try {
    await authStore.loginWithGoogle(response.credential)
    router.push('/')
  } catch {
    error.value = 'เข้าสู่ระบบไม่สำเร็จ กรุณาลองใหม่'
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f9fafb;
}
.login-card {
  background: white;
  padding: 40px;
  border-radius: 16px;
  text-align: center;
  box-shadow: 0 4px 20px rgba(0,0,0,0.08);
}
.login-card h1 { font-size: 24px; margin-bottom: 8px; }
.login-card p  { color: #6b7280; margin-bottom: 24px; }
.error { color: #ef4444; margin-top: 12px; }
</style>