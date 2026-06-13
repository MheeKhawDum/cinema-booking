<template>
  <div class="home-page">
    <header class="app-header">
      <h1>🎬 Cinema Booking</h1>
      <div class="user-info">
        <img :src="authStore.user?.picture" class="avatar" />
        <span>{{ authStore.user?.name }}</span>
        <router-link v-if="authStore.isAdmin" to="/admin" class="btn-admin">Admin</router-link>
        <router-link to="/bookings/my" class="btn-outline">ประวัติการจอง</router-link>
        <button @click="handleLogout" class="btn-logout">ออกจากระบบ</button>
      </div>
    </header>

    <div v-if="loading" class="loading">กำลังโหลด...</div>

    <div v-else class="movie-grid">
      <div v-for="movie in movies" :key="movie.id"
           class="movie-card" @click="selectMovie(movie)">
        <img :src="movie.poster_url" :alt="movie.title" class="movie-poster" />
        <div class="movie-info">
          <h3>{{ movie.title }}</h3>
          <p class="duration">⏱ {{ movie.duration }} นาที</p>
          <p class="description">{{ movie.description }}</p>
        </div>
      </div>
    </div>

    <Transition name="fade">
      <div v-if="selectedMovie" class="modal-overlay" @click.self="selectedMovie = null">
        <div class="modal">
          <button class="modal-close" @click="selectedMovie = null">✕</button>
          <h2>{{ selectedMovie.title }}</h2>
          <p class="modal-subtitle">เลือกรอบฉาย</p>
          <div v-if="loadingShowtimes" class="loading">กำลังโหลด...</div>
          <div v-else class="showtime-list">
            <div v-for="st in showtimes" :key="st.id"
                 class="showtime-item" @click="goToSeatMap(st)">
              <div class="st-time">{{ formatTime(st.start_time) }}</div>
              <div class="st-info">
                <span class="st-hall">{{ st.hall }}</span>
                <span class="st-date">{{ formatDate(st.start_time) }}</span>
              </div>
              <span class="st-arrow">→</span>
            </div>
            <p v-if="showtimes.length === 0" class="empty">ไม่มีรอบฉาย</p>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import axios from 'axios'

const router    = useRouter()
const authStore = useAuthStore()

const movies           = ref([])
const showtimes        = ref([])
const selectedMovie    = ref(null)
const loading          = ref(true)
const loadingShowtimes = ref(false)

onMounted(async () => {
  try {
    const { data } = await axios.get('/api/movies')
    movies.value = data.movies
  } finally {
    loading.value = false
  }
})

async function selectMovie(movie) {
  selectedMovie.value    = movie
  loadingShowtimes.value = true
  showtimes.value        = []
  try {
    const { data } = await axios.get(`/api/movies/${movie.id}/showtimes`)
    showtimes.value = data.showtimes
  } finally {
    loadingShowtimes.value = false
  }
}

function goToSeatMap(st) {
  router.push({
    name:  'seat-map',
    params: { id: st.id },
    query:  { movieId: selectedMovie.value.id }
  })
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}

const formatTime = (iso) => new Date(iso).toLocaleTimeString('th-TH', { hour: '2-digit', minute: '2-digit' })
const formatDate = (iso) => new Date(iso).toLocaleDateString('th-TH', { weekday: 'short', month: 'short', day: 'numeric' })
</script>

<style scoped>
.home-page { max-width: 1100px; margin: 0 auto; padding: 20px; }
.app-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 32px; padding-bottom: 16px; border-bottom: 1px solid #f0f0f0; }
.app-header h1 { font-size: 24px; font-weight: 800; }
.user-info { display: flex; align-items: center; gap: 10px; }
.avatar { width: 36px; height: 36px; border-radius: 50%; object-fit: cover; }
.btn-admin { background: #7c3aed; color: white; padding: 6px 14px; border-radius: 8px; text-decoration: none; font-size: 13px; font-weight: 600; }
.btn-outline { border: 1px solid #e5e7eb; padding: 6px 14px; border-radius: 8px; text-decoration: none; font-size: 13px; color: #374151; }
.btn-logout { background: none; border: 1px solid #e5e7eb; padding: 6px 14px; border-radius: 8px; cursor: pointer; font-size: 13px; }
.movie-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(220px, 1fr)); gap: 20px; }
.movie-card { border-radius: 14px; overflow: hidden; border: 1.5px solid #f0f0f0; cursor: pointer; transition: transform 0.2s, box-shadow 0.2s; background: white; }
.movie-card:hover { transform: translateY(-4px); box-shadow: 0 8px 24px rgba(0,0,0,0.1); }
.movie-poster { width: 100%; aspect-ratio: 2/3; object-fit: cover; }
.movie-info { padding: 14px; }
.movie-info h3 { font-size: 15px; font-weight: 700; margin-bottom: 4px; }
.duration { font-size: 12px; color: #6b7280; margin-bottom: 6px; }
.description { font-size: 12px; color: #9ca3af; overflow: hidden; display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; }
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; z-index: 100; padding: 20px; }
.modal { background: white; border-radius: 20px; padding: 28px; width: 100%; max-width: 480px; max-height: 80vh; overflow-y: auto; position: relative; }
.modal-close { position: absolute; top: 16px; right: 16px; background: #f3f4f6; border: none; border-radius: 50%; width: 32px; height: 32px; cursor: pointer; }
.modal h2 { font-size: 20px; font-weight: 800; margin-bottom: 4px; }
.modal-subtitle { color: #6b7280; margin-bottom: 20px; font-size: 14px; }
.showtime-item { display: flex; align-items: center; gap: 14px; padding: 14px 16px; border-radius: 12px; border: 1.5px solid #f0f0f0; margin-bottom: 10px; cursor: pointer; transition: all 0.15s; }
.showtime-item:hover { border-color: #3b82f6; background: #eff6ff; }
.st-time { font-size: 22px; font-weight: 800; color: #3b82f6; width: 64px; }
.st-info { flex: 1; }
.st-hall { display: block; font-weight: 600; font-size: 14px; }
.st-date { font-size: 12px; color: #6b7280; }
.st-arrow { color: #9ca3af; }
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.loading, .empty { text-align: center; padding: 40px; color: #9ca3af; }
</style>