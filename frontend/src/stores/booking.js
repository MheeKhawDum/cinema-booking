import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export const useBookingStore = defineStore('booking', () => {
  const seats          = ref([])
  const selectedSeat   = ref(null)
  const currentBooking = ref(null)

  async function loadSeats(showtimeId) {
    const { data } = await axios.get(`/api/showtimes/${showtimeId}/seats`)
    seats.value = data.seats
  }

  function updateSeatStatus(seatNumber, status) {
    const seat = seats.value.find(s => s.number === seatNumber)
    if (seat) seat.status = status
  }

  async function lockSeat(showtimeId, seatNumber) {
    const { data } = await axios.post('/api/bookings', {
      showtime_id: showtimeId,
      seat: seatNumber
    })
    currentBooking.value = data
    return data
  }

  async function confirmPayment(bookingId) {
    const { data } = await axios.post(`/api/bookings/${bookingId}/payment`, {
      booking_id: bookingId
    })
    currentBooking.value = null
    return data
  }

  return {
    seats, selectedSeat, currentBooking,
    loadSeats, updateSeatStatus, lockSeat, confirmPayment
  }
})