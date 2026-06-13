<template>
  <div class="payment-page">
    <h2>ยืนยันการจอง</h2>

    <div :class="['timer', timeLeft < 60 ? 'urgent' : '']">
      ⏱ เหลือเวลา {{ formatTime(timeLeft) }}
    </div>

    <div class="booking-summary">
      <p>ที่นั่ง: <strong>{{ route.query.seat }}</strong></p>
      <p>ราคา: <strong>฿ 150</strong></p>
    </div>

    <div class="payment-form">
      <p class="mock-note">🧪 Mock Payment — กดปุ่มเพื่อจำลองการจ่ายเงิน</p>
      <button class="btn-pay" @click="handlePayment" :disabled="isPaying || timeLeft <= 0">
        {{ isPaying ? 'กำลังดำเนินการ...' : 'ชำระเงิน ฿150' }}
      </button>
    </div>

    <div v-if="timeLeft <= 0" class="timeout-msg">
      หมดเวลา ที่นั่งถูกปล่อยแล้ว
      <router-link to="/">กลับหน้าหลัก</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useBookingStore } from '../stores/booking'

const route        = useRoute()
const router       = useRouter()
const bookingStore = useBookingStore()

const isPaying = ref(false)
const timeLeft = ref(300)  
let timer      = null

onMounted(() => {
  if (route.query.expires) {
    const expires = new Date(route.query.expires).getTime()
    timeLeft.value = Math.max(0, Math.floor((expires - Date.now()) / 1000))
  }

  timer = setInterval(() => {
    if (timeLeft.value <= 0) {
      clearInterval(timer)
      return
    }
    timeLeft.value--
  }, 1000)
})

onUnmounted(() => clearInterval(timer))

async function handlePayment() {
  isPaying.value = true
  try {
    await bookingStore.confirmPayment(route.params.bookingId)
    router.push({ name: 'booking-success', params: { id: route.params.bookingId } })
  } catch (err) {
    if (err.response?.status === 410) {  
      alert('หมดเวลาจ่ายเงิน ที่นั่งถูกปล่อยแล้ว')
      router.push('/')
    }
  } finally {
    isPaying.value = false
  }
}

function formatTime(seconds) {
  const m = Math.floor(seconds / 60).toString().padStart(2, '0')
  const s = (seconds % 60).toString().padStart(2, '0')
  return `${m}:${s}`
}
</script>

<style scoped>
.timer {
  font-size: 2rem;
  font-weight: 700;
  text-align: center;
  color: #3b82f6;
  margin: 20px 0;
  transition: color 0.3s;
}
.timer.urgent { color: #ef4444; animation: pulse 1s infinite; }
@keyframes pulse {
  0%, 100% { opacity: 1; }
  50%       { opacity: 0.5; }
}
</style>