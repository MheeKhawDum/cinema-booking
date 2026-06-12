package websocket

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

type Client struct {
    hub        *Hub
    conn       *websocket.Conn
    send       chan SeatEvent
    showtimeID string
}

func (c *Client) WritePump() {
    defer c.conn.Close()
    for event := range c.send {
        c.conn.WriteJSON(event)
    }
}

func Handler(hub *Hub) gin.HandlerFunc {
    upgrader := websocket.Upgrader{
        CheckOrigin: func(r *http.Request) bool { return true },
    }
    return func(c *gin.Context) {
        showtimeID := c.Query("showtime_id")
        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            return
        }
        client := &Client{
            hub:        hub,
            conn:       conn,
            send:       make(chan SeatEvent, 64),
            showtimeID: showtimeID,
        }
        hub.register <- client
        go client.WritePump()
    }
}