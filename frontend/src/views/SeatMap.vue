<template>
  <div class="seat-map">
    <h2>เลือกที่นั่ง — รอบ {{ showtime?.start_time }}</h2>

    <div class="legend">
      <span class="dot available"></span> ว่าง
      <span class="dot locked"></span> กำลังถูกเลือก
      <span class="dot booked"></span> จองแล้ว
      <span class="dot selected"></span> ที่คุณเลือก
    </div>

    <div class="seat-grid">
      <button
        v-for="seat in seats"
        :key="seat.number"
        :class="['seat', getSeatClass(seat)]"
        :disabled="seat.status !== 'AVAILABLE'"
        @click="selectSeat(seat)"
      >
        {{ seat.number }}
      </button>
    </div>

    <button
      v-if="selectedSeat"
      class="btn-book"
      @click="confirmBooking"
    >
      จองที่นั่ง {{ selectedSeat }} — 150 บาท
    </button>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const seats = ref([])
const selectedSeat = ref(null)
const ws = ref(null)

onMounted(async () => {
  const { data } = await axios.get(`/api/showtimes/${route.params.id}/seats`, {
    headers: { Authorization: `Bearer ${authStore.token}` }
  })
  seats.value = data.seats

  connectWebSocket()
})

function connectWebSocket() {
  const wsHost = window.location.hostname
  ws.value = new WebSocket(`ws://${wsHost}:8080/ws?showtime_id=${route.params.id}`)

  ws.value.onmessage = (event) => {
    const update = JSON.parse(event.data)
    const seat = seats.value.find(s => s.number === update.seat)
    if (seat) {
      seat.status = update.status
    }
  }

  ws.value.onclose = () => {
    setTimeout(connectWebSocket, 3000)
  }
}

function selectSeat(seat) {
  if (seat.status !== 'AVAILABLE') return
  selectedSeat.value = seat.number
}

async function confirmBooking() {
  try {
    const { data } = await axios.post('/api/bookings', {
      showtime_id: route.params.id,
      seat: selectedSeat.value
    }, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })

    router.push({
      name: 'payment',
      params: { bookingId: data.booking_id },
      query:  { expires: data.expires_at, seat: selectedSeat.value }
    })
  } catch (err) {
    if (err.response?.status === 409) {
      alert('ที่นั่งนี้ถูกเลือกโดยคนอื่นแล้ว กรุณาเลือกที่นั่งใหม่')
      selectedSeat.value = null
    }
  }
}

function getSeatClass(seat) {
  if (seat.number === selectedSeat.value) return 'selected'
  return seat.status.toLowerCase() 
}

onUnmounted(() => ws.value?.close())
</script>