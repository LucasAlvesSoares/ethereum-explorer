package websocket

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin in development
		// In production, you should check the origin properly
		return true
	},
}

// Message represents a WebSocket message
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// BlockUpdate represents a new block notification
type BlockUpdate struct {
	Number           int64  `json:"number"`
	Hash             string `json:"hash"`
	TransactionCount int    `json:"transaction_count"`
	GasUsed          uint64 `json:"gas_used"`
	GasLimit         uint64 `json:"gas_limit"`
	Timestamp        string `json:"timestamp"`
	Miner            string `json:"miner"`
}

// TransactionUpdate represents a new transaction notification
type TransactionUpdate struct {
	Hash        string `json:"hash"`
	BlockNumber int64  `json:"block_number"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address,omitempty"`
	Value       string `json:"value"`
	GasPrice    string `json:"gas_price"`
}

// Client represents a WebSocket client
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	topics map[string]bool // subscribed topics
	mu     sync.RWMutex
}

// Hub maintains the set of active clients and broadcasts messages to them
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			logrus.Infof("WebSocket client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			logrus.Infof("WebSocket client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					delete(h.clients, client)
					close(client.send)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastBlockUpdate sends a block update to all connected clients
func (h *Hub) BroadcastBlockUpdate(block BlockUpdate) {
	message := Message{
		Type: "block_update",
		Data: block,
	}

	data, err := json.Marshal(message)
	if err != nil {
		logrus.Errorf("Failed to marshal block update: %v", err)
		return
	}

	select {
	case h.broadcast <- data:
	default:
		logrus.Warn("Broadcast channel is full, dropping block update")
	}
}

// BroadcastTransactionUpdate sends a transaction update to all connected clients
func (h *Hub) BroadcastTransactionUpdate(tx TransactionUpdate) {
	message := Message{
		Type: "transaction_update",
		Data: tx,
	}

	data, err := json.Marshal(message)
	if err != nil {
		logrus.Errorf("Failed to marshal transaction update: %v", err)
		return
	}

	select {
	case h.broadcast <- data:
	default:
		logrus.Warn("Broadcast channel is full, dropping transaction update")
	}
}

// BroadcastNetworkStats sends network statistics to all connected clients
func (h *Hub) BroadcastNetworkStats(stats interface{}) {
	message := Message{
		Type: "network_stats",
		Data: stats,
	}

	data, err := json.Marshal(message)
	if err != nil {
		logrus.Errorf("Failed to marshal network stats: %v", err)
		return
	}

	select {
	case h.broadcast <- data:
	default:
		logrus.Warn("Broadcast channel is full, dropping network stats")
	}
}

// HandleWebSocket handles WebSocket connections
func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("Failed to upgrade WebSocket connection: %v", err)
		return
	}

	client := &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
	}

	client.hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		var msg map[string]interface{}
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("WebSocket error: %v", err)
			}
			break
		}

		// Handle subscription messages
		if msgType, ok := msg["type"].(string); ok {
			switch msgType {
			case "subscribe":
				if topic, ok := msg["topic"].(string); ok {
					c.mu.Lock()
					c.topics[topic] = true
					c.mu.Unlock()
					logrus.Debugf("Client subscribed to topic: %s", topic)
				}
			case "unsubscribe":
				if topic, ok := msg["topic"].(string); ok {
					c.mu.Lock()
					delete(c.topics, topic)
					c.mu.Unlock()
					logrus.Debugf("Client unsubscribed from topic: %s", topic)
				}
			}
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logrus.Errorf("WebSocket write error: %v", err)
				return
			}
		}
	}
}

// IsSubscribed checks if client is subscribed to a topic
func (c *Client) IsSubscribed(topic string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.topics[topic]
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
