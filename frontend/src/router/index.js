import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/Login.vue'),
    beforeEnter: (to, from, next) => {
      const auth = useAuthStore()
      auth.isLoggedIn ? next('/') : next()
    }
  },
  {
    path: '/',
    name: 'home',
    component: () => import('../views/Home.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/showtimes/:id/seats',
    name: 'seat-map',
    component: () => import('../views/SeatMap.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/payment/:bookingId',
    name: 'payment',
    component: () => import('../views/Payment.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/booking/:id/success',
    name: 'booking-success',
    component: () => import('../views/BookingSuccess.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/bookings/my',
    name: 'my-bookings',
    component: () => import('../views/MyBookings.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/admin',
    name: 'admin',
    component: () => import('../views/Admin.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth && !auth.isLoggedIn) {
    next('/login')
    return
  }
  if (to.meta.requiresAdmin && auth.user?.role !== 'ADMIN') {
    next('/')
    return
  }
  next()
})

export default router