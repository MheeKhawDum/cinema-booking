package websocket

import "sync"

type SeatEvent struct {
    ShowtimeID string `json:"showtime_id"`
    Seat       string `json:"seat"`
    Status     string `json:"status"` 
    UserID     string `json:"user_id,omitempty"`
}

type Hub struct {
    clients    map[*Client]bool   
    broadcast  chan SeatEvent     
    register   chan *Client       
    unregister chan *Client       
    mu         sync.RWMutex     
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan SeatEvent, 256),  
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mu.Unlock()

        case event := <-h.broadcast:
            h.mu.RLock()
            for client := range h.clients {
                if client.showtimeID == event.ShowtimeID {
                    select {
                    case client.send <- event:
                    default:
                        close(client.send)
                        delete(h.clients, client)
                    }
                }
            }
            h.mu.RUnlock()
        }
    }
}

func (h *Hub) Broadcast(event SeatEvent) {
    h.broadcast <- event
}