package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

// Transport handles the network layer for bridge communication
type Transport interface {
	Connect(ctx context.Context, cfg BridgeConfig) error
	Send(msg ProtocolMessage) error
	Receive() (ProtocolMessage, error)
	Close() error
	IsConnected() bool
}

// WebSocketTransport implements bridge over WebSocket
type WebSocketTransport struct {
	mu        sync.Mutex
	conn      net.Conn
	listener  net.Listener
	inbox     chan ProtocolMessage
	outbox    chan ProtocolMessage
	done      chan struct{}
	connected bool
}

func NewWebSocketTransport() *WebSocketTransport {
	return &WebSocketTransport{
		inbox:  make(chan ProtocolMessage, 100),
		outbox: make(chan ProtocolMessage, 100),
		done:   make(chan struct{}),
	}
}

func (t *WebSocketTransport) Connect(ctx context.Context, cfg BridgeConfig) error {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("bridge listen: %w", err)
	}
	t.listener = listener
	t.connected = true

	// Accept connections in background
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-t.done:
					return
				default:
					continue
				}
			}
			t.mu.Lock()
			t.conn = conn
			t.mu.Unlock()
			go t.readLoop(conn)
		}
	}()

	return nil
}

func (t *WebSocketTransport) readLoop(conn net.Conn) {
	decoder := json.NewDecoder(conn)
	for {
		var msg ProtocolMessage
		if err := decoder.Decode(&msg); err != nil {
			return
		}
		select {
		case t.inbox <- msg:
		case <-t.done:
			return
		}
	}
}

func (t *WebSocketTransport) Send(msg ProtocolMessage) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.conn == nil {
		return fmt.Errorf("not connected")
	}
	msg.Timestamp = time.Now().UnixMilli()
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	_, err = t.conn.Write(append(data, '\n'))
	return err
}

func (t *WebSocketTransport) Receive() (ProtocolMessage, error) {
	select {
	case msg := <-t.inbox:
		return msg, nil
	case <-t.done:
		return ProtocolMessage{}, fmt.Errorf("transport closed")
	}
}

func (t *WebSocketTransport) Close() error {
	close(t.done)
	t.connected = false
	if t.listener != nil {
		t.listener.Close()
	}
	t.mu.Lock()
	if t.conn != nil {
		t.conn.Close()
	}
	t.mu.Unlock()
	return nil
}

func (t *WebSocketTransport) IsConnected() bool { return t.connected }

// PollingTransport implements bridge over HTTP long-polling
type PollingTransport struct {
	mu        sync.Mutex
	server    *http.Server
	inbox     chan ProtocolMessage
	outbox    []ProtocolMessage
	done      chan struct{}
	connected bool
}

func NewPollingTransport() *PollingTransport {
	return &PollingTransport{
		inbox: make(chan ProtocolMessage, 100),
		done:  make(chan struct{}),
	}
}

func (t *PollingTransport) Connect(ctx context.Context, cfg BridgeConfig) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/poll", t.handlePoll)
	mux.HandleFunc("/send", t.handleSend)
	mux.HandleFunc("/status", t.handleStatus)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	t.server = &http.Server{Addr: addr, Handler: mux}
	t.connected = true

	go t.server.ListenAndServe() //nolint:errcheck
	return nil
}

func (t *PollingTransport) handlePoll(w http.ResponseWriter, r *http.Request) {
	t.mu.Lock()
	msgs := t.outbox
	t.outbox = nil
	t.mu.Unlock()
	json.NewEncoder(w).Encode(msgs) //nolint:errcheck
}

func (t *PollingTransport) handleSend(w http.ResponseWriter, r *http.Request) {
	var msg ProtocolMessage
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	t.inbox <- msg
	w.WriteHeader(http.StatusOK)
}

func (t *PollingTransport) handleStatus(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{"connected": t.connected}) //nolint:errcheck
}

func (t *PollingTransport) Send(msg ProtocolMessage) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	msg.Timestamp = time.Now().UnixMilli()
	t.outbox = append(t.outbox, msg)
	return nil
}

func (t *PollingTransport) Receive() (ProtocolMessage, error) {
	select {
	case msg := <-t.inbox:
		return msg, nil
	case <-t.done:
		return ProtocolMessage{}, fmt.Errorf("closed")
	}
}

func (t *PollingTransport) Close() error {
	close(t.done)
	t.connected = false
	if t.server != nil {
		t.server.Shutdown(context.Background()) //nolint:errcheck
	}
	return nil
}

func (t *PollingTransport) IsConnected() bool { return t.connected }
