<template>
  <div class="success-page">
    <div class="success-card">
      <div class="checkmark-circle">
        <svg viewBox="0 0 52 52" class="checkmark-svg">
          <circle class="checkmark-circle-bg" cx="26" cy="26" r="25" fill="none"/>
          <path class="checkmark-check" fill="none" d="M14 27l7 7 17-17"/>
        </svg>
      </div>

      <h1>จองสำเร็จ! 🎉</h1>
      <p class="subtitle">ตั๋วของคุณได้รับการยืนยันแล้ว</p>

      <div class="ticket">
        <div class="ticket-header">
          <span class="ticket-label">🎬 CINEMA TICKET</span>
          <span class="ticket-id">#{{ bookingId?.slice(-8).toUpperCase() }}</span>
        </div>

        <div class="ticket-divider">
          <span class="ticket-dot left"/>
          <div class="ticket-dashes"/>
          <span class="ticket-dot right"/>
        </div>

        <div class="ticket-body">
          <div class="ticket-row">
            <span class="t-label">ที่นั่ง</span>
            <span class="t-value seat-badge">{{ route.query.seat }}</span>
          </div>
          <div class="ticket-row">
            <span class="t-label">รอบฉาย</span>
            <span class="t-value">{{ formatDate(route.query.time) }}</span>
          </div>
          <div class="ticket-row">
            <span class="t-label">สถานะ</span>
            <span class="t-value status-badge">✅ CONFIRMED</span>
          </div>
        </div>
      </div>

      <div class="actions">
        <router-link to="/" class="btn-home">🏠 กลับหน้าหลัก</router-link>
        <router-link to="/bookings/my" class="btn-history">ดูประวัติการจอง</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'

const route     = useRoute()
const bookingId = route.params.id

function formatDate(iso) {
  if (!iso) return '-'
  return new Date(iso).toLocaleString('th-TH', {
    dateStyle: 'medium', timeStyle: 'short'
  })
}
</script>

<style scoped>
.success-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #eff6ff 0%, #f0fdf4 100%);
  padding: 20px;
}
.success-card {
  background: white;
  border-radius: 24px;
  padding: 40px 32px;
  text-align: center;
  max-width: 400px;
  width: 100%;
  box-shadow: 0 20px 60px rgba(0,0,0,0.08);
}

.checkmark-circle {
  width: 80px; height: 80px;
  margin: 0 auto 24px;
}
.checkmark-svg { width: 100%; }
.checkmark-circle-bg {
  stroke: #22c55e;
  stroke-width: 2;
  stroke-dasharray: 166;
  stroke-dashoffset: 166;
  animation: stroke 0.6s cubic-bezier(0.65, 0, 0.45, 1) forwards;
}
.checkmark-check {
  stroke: #22c55e;
  stroke-width: 3;
  stroke-linecap: round;
  stroke-dasharray: 48;
  stroke-dashoffset: 48;
  animation: stroke 0.4s cubic-bezier(0.65, 0, 0.45, 1) 0.5s forwards;
}
@keyframes stroke {
  100% { stroke-dashoffset: 0; }
}

h1 { font-size: 28px; font-weight: 800; margin-bottom: 8px; }
.subtitle { color: #6b7280; margin-bottom: 28px; }

.ticket {
  background: #f9fafb;
  border-radius: 16px;
  overflow: hidden;
  margin-bottom: 28px;
  border: 1.5px solid #e5e7eb;
}
.ticket-header {
  background: #1e40af;
  color: white;
  padding: 14px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  font-weight: 700;
}
.ticket-id { font-family: monospace; opacity: 0.8; }

.ticket-divider {
  position: relative;
  display: flex;
  align-items: center;
  padding: 0 -10px;
}
.ticket-dot {
  width: 20px; height: 20px;
  border-radius: 50%;
  background: white;
  border: 1.5px solid #e5e7eb;
  position: absolute;
  top: -10px;
}
.ticket-dot.left  { left: -10px; }
.ticket-dot.right { right: -10px; }
.ticket-dashes {
  flex: 1;
  border-top: 2px dashed #e5e7eb;
  margin: 0 10px;
}

.ticket-body { padding: 20px; }
.ticket-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f3f4f6;
}
.ticket-row:last-child { border-bottom: none; }
.t-label { font-size: 13px; color: #9ca3af; }
.t-value { font-size: 14px; font-weight: 600; }

.seat-badge {
  background: #3b82f6;
  color: white;
  padding: 4px 14px;
  border-radius: 99px;
  font-size: 16px;
  font-weight: 800;
}
.status-badge { color: #16a34a; }

.actions { display: flex; flex-direction: column; gap: 10px; }
.btn-home {
  background: #3b82f6;
  color: white;
  padding: 14px;
  border-radius: 12px;
  text-decoration: none;
  font-weight: 700;
  font-size: 15px;
}
.btn-history {
  background: none;
  border: 1.5px solid #e5e7eb;
  color: #374151;
  padding: 12px;
  border-radius: 12px;
  text-decoration: none;
  font-weight: 600;
  font-size: 14px;
}
</style>