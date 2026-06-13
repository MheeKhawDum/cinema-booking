<template>
  <div class="admin-page">
    <h1>Admin Dashboard</h1>

    <div class="stats-grid">
      <div class="stat-card" v-for="stat in stats" :key="stat.status">
        <p class="stat-value">{{ stat.count }}</p>
        <p class="stat-label">{{ stat.status }}</p>
      </div>
    </div>

    <div class="filters">
      <select v-model="filter.status" @change="loadBookings">
        <option value="">ทุก Status</option>
        <option value="PENDING">PENDING</option>
        <option value="BOOKED">BOOKED</option>
        <option value="TIMEOUT">TIMEOUT</option>
        <option value="CANCELLED">CANCELLED</option>
      </select>

      <input type="date" v-model="filter.date" @change="loadBookings" />

      <button @click="resetFilters">รีเซ็ต</button>
    </div>

    <div class="table-wrapper">
      <table>
        <thead>
          <tr>
            <th>Booking ID</th>
            <th>ผู้ใช้</th>
            <th>ที่นั่ง</th>
            <th>Status</th>
            <th>เวลาจอง</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="b in bookings" :key="b._id">
            <td class="mono">
              {{ b._id?.$oid?.slice(-8) ?? b._id?.toString().slice(-8) }}
            </td>
            <td>{{ b.user?.[0]?.name ?? "Unknown" }}</td>
            <td>{{ b.seat_number ?? "-" }}</td>
            <td>
              <span :class="['badge', b.status?.toLowerCase()]">
                {{ b.status }}
              </span>
            </td>
            <td>{{ formatDate(b.created_at) }}</td>
          </tr>
        </tbody>
      </table>
      <p v-if="bookings.length === 0" class="empty">ไม่มีรายการ</p>
    </div>

    <div class="audit-section">
      <h2>Audit Logs</h2>
      <select v-model="auditFilter" @change="loadAuditLogs">
        <option value="">ทุก Event</option>
        <option value="BOOKING_SUCCESS">BOOKING_SUCCESS</option>
        <option value="BOOKING_TIMEOUT">BOOKING_TIMEOUT</option>
        <option value="SEAT_RELEASED">SEAT_RELEASED</option>
        <option value="LOCK_FAIL">LOCK_FAIL</option>
      </select>

      <div class="log-list">
        <div
          v-for="log in auditLogs"
          :key="log._id"
          :class="['log-item', getLogClass(log.event)]"
        >
          <span class="log-time">{{ formatDate(log.created_at) }}</span>
          <span class="log-event">{{ log.event }}</span>
          <span class="log-user">{{ log.user_id }}</span>
          <pre class="log-data">{{ JSON.stringify(log.data, null, 2) }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import axios from "axios";

const bookings = ref([]);
const auditLogs = ref([]);
const stats = ref([]);
const auditFilter = ref("");

const filter = ref({
  status: "",
  date: "",
});

onMounted(async () => {
  await Promise.all([loadBookings(), loadAuditLogs(), loadStats()]);
});

async function loadBookings() {
  const params = {};
  if (filter.value.status) params.status = filter.value.status;
  if (filter.value.date) params.date = filter.value.date;

  const { data } = await axios.get("/api/admin/bookings", { params });
  bookings.value = data.bookings ?? [];
}

async function loadAuditLogs() {
  const params = auditFilter.value ? { event: auditFilter.value } : {};
  const { data } = await axios.get("/api/admin/audit-logs", { params });
  auditLogs.value = data.logs ?? [];
}

async function loadStats() {
  const { data } = await axios.get("/api/admin/stats");
  stats.value = data.booking_stats ?? [];
}

function resetFilters() {
  filter.value = { status: "", date: "" };
  loadBookings();
}

function getLogClass(event) {
  const map = {
    BOOKING_SUCCESS: "success",
    BOOKING_TIMEOUT: "warning",
    SEAT_RELEASED: "info",
    LOCK_FAIL: "error",
  };
  return map[event] ?? "info";
}

function formatDate(iso) {
  if (!iso) return "";
  return new Date(iso).toLocaleString("th-TH");
}
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 12px;
  margin-bottom: 24px;
}
.stat-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  text-align: center;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}
.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: #3b82f6;
}
table {
  width: 100%;
  border-collapse: collapse;
}
th,
td {
  padding: 10px 12px;
  text-align: left;
  border-bottom: 1px solid #f0f0f0;
}
th {
  background: #f9fafb;
  font-weight: 600;
}
.badge {
  padding: 3px 10px;
  border-radius: 99px;
  font-size: 12px;
  font-weight: 600;
}
.badge.booked {
  background: #dcfce7;
  color: #15803d;
}
.badge.pending {
  background: #fef9c3;
  color: #854d0e;
}
.badge.timeout {
  background: #fee2e2;
  color: #991b1b;
}
.log-item {
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 8px;
}
.log-item.success {
  background: #f0fdf4;
  border-left: 3px solid #22c55e;
}
.log-item.warning {
  background: #fffbeb;
  border-left: 3px solid #f59e0b;
}
.log-item.error {
  background: #fef2f2;
  border-left: 3px solid #ef4444;
}
.log-item.info {
  background: #eff6ff;
  border-left: 3px solid #3b82f6;
}
.log-data {
  font-size: 11px;
  color: #6b7280;
  margin: 4px 0 0;
}
</style>
