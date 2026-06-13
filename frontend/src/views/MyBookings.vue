<template>
  <div class="my-bookings-page">
    <div class="page-header">
      <router-link to="/" class="back-btn">← กลับ</router-link>
      <h2>ประวัติการจอง</h2>
    </div>

    <div v-if="loading" class="loading">กำลังโหลด...</div>

    <div v-else-if="bookings.length === 0" class="empty-state">
      <p>🎟 ยังไม่มีประวัติการจอง</p>
      <router-link to="/" class="btn-browse">ดูหนังทั้งหมด</router-link>
    </div>

    <div v-else class="booking-list">
      <div
        v-for="b in bookings"
        :key="b._id"
        class="booking-item"
      >
        <div class="booking-poster">
          <img :src="b.movie?.poster_url" :alt="b.movie?.title" />
        </div>
        <div class="booking-detail">
          <h3>{{ b.movie?.title ?? 'Unknown Movie' }}</h3>
          <p class="detail-row">🪑 ที่นั่ง <strong>{{ b.seat_number }}</strong></p>
          <p class="detail-row">🏛 {{ b.showtime?.hall }}</p>
          <p class="detail-row">📅 {{ formatDate(b.showtime?.start_time) }}</p>
          <p class="detail-row">🕐 จองเมื่อ {{ formatDate(b.created_at) }}</p>
        </div>
        <div class="booking-status">
          <span :class="['status-pill', b.status.toLowerCase()]">
            {{ statusLabel(b.status) }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const bookings = ref([])
const loading  = ref(true)

onMounted(async () => {
  try {
    const { data } = await axios.get('/api/bookings/my')
    bookings.value = data.bookings
  } finally {
    loading.value = false
  }
})

function statusLabel(status) {
  const map = {
    BOOKED:    '✅ จองแล้ว',
    PENDING:   '⏳ รอชำระ',
    TIMEOUT:   '⌛ หมดเวลา',
    CANCELLED: '❌ ยกเลิก',
  }
  return map[status] ?? status
}

function formatDate(iso) {
  if (!iso) return '-'
  return new Date(iso).toLocaleString('th-TH', {
    dateStyle: 'medium', timeStyle: 'short'
  })
}
</script>

<style scoped>
.my-bookings-page { max-width: 640px; margin: 0 auto; padding: 20px; }
.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}
.back-btn { color: #3b82f6; text-decoration: none; font-weight: 600; }
.page-header h2 { font-size: 20px; font-weight: 800; }

.booking-item {
  display: flex;
  gap: 14px;
  padding: 16px;
  border-radius: 14px;
  border: 1.5px solid #f0f0f0;
  margin-bottom: 12px;
  background: white;
}
.booking-poster img {
  width: 56px; height: 80px;
  border-radius: 8px;
  object-fit: cover;
  background: #f3f4f6;
}
.booking-detail { flex: 1; }
.booking-detail h3 { font-size: 15px; font-weight: 700; margin-bottom: 6px; }
.detail-row { font-size: 13px; color: #6b7280; margin-bottom: 3px; }
.detail-row strong { color: #111; }

.status-pill {
  padding: 5px 12px;
  border-radius: 99px;
  font-size: 12px;
  font-weight: 700;
  white-space: nowrap;
}
.status-pill.booked    { background: #dcfce7; color: #15803d; }
.status-pill.pending   { background: #fef9c3; color: #854d0e; }
.status-pill.timeout   { background: #fee2e2; color: #991b1b; }
.status-pill.cancelled { background: #f3f4f6; color: #6b7280; }

.empty-state { text-align: center; padding: 60px 0; color: #9ca3af; }
.btn-browse {
  display: inline-block;
  margin-top: 16px;
  background: #3b82f6;
  color: white;
  padding: 10px 24px;
  border-radius: 10px;
  text-decoration: none;
  font-weight: 600;
}
.loading { text-align: center; padding: 60px; color: #9ca3af; }
</style>