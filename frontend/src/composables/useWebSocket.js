import { ref, onUnmounted } from 'vue'
import { useBookingStore } from '../stores/booking'

export function useWebSocket(showtimeId) {
  const bookingStore = useBookingStore()
  const connected    = ref(false)
  let ws             = null
  let reconnectTimer = null

  function connect() {
    ws = new WebSocket(`ws://localhost:8080/ws?showtime_id=${showtimeId}`)

    ws.onopen = () => {
      connected.value = true
      if (reconnectTimer) clearTimeout(reconnectTimer)
    }

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      bookingStore.updateSeatStatus(data.seat, data.status)
    }

    ws.onclose = () => {
      connected.value = false
      reconnectTimer = setTimeout(connect, 3000)
    }

    ws.onerror = () => ws.close()
  }

  function disconnect() {
    if (reconnectTimer) clearTimeout(reconnectTimer)
    ws?.close()
  }

  onUnmounted(disconnect)

  return { connected, connect, disconnect }
}