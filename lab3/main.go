package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const numPeers = 4

var peerPorts = []int{8081, 8082, 8083, 8084}

type PeerStatus struct {
	Port   int    `json:"port"`
	Status string `json:"status"` // "живой" или "мертв"
}

type Peer struct {
	Port       int
	Status     string
	killChan   chan bool
	statusLock sync.Mutex
}

var peers = make(map[int]*Peer)

var statusUpdates = make(chan []PeerStatus)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	for _, port := range peerPorts {
		peers[port] = &Peer{
			Port:     port,
			Status:   "живой",
			killChan: make(chan bool),
		}
		go startPeer(peers[port])
	}

	go monitorPeers()

	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/kill", handleKillPeer)

	go func() {
		for {
			var statuses []PeerStatus
			for _, peer := range peers {
				peer.statusLock.Lock()
				statuses = append(statuses, PeerStatus{
					Port:   peer.Port,
					Status: peer.Status,
				})
				peer.statusLock.Unlock()
			}
			statusUpdates <- statuses
			time.Sleep(2 * time.Second)
		}
	}()

	log.Println("Запуск основного сервера на порту 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func startPeer(peer *Peer) {
	peer.statusLock.Lock()
	peer.Status = "живой"
	peer.statusLock.Unlock()

	mux := http.NewServeMux()
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		peer.statusLock.Lock()
		defer peer.statusLock.Unlock()
		fmt.Fprintf(w, "Пир %d: %s", peer.Port, peer.Status)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", peer.Port),
		Handler: mux,
	}

	go func() {
		log.Printf("Пир %d запущен\n", peer.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Пир %d остановлен\n", peer.Port)
		}
	}()

	<-peer.killChan

	peer.statusLock.Lock()
	peer.Status = "мертв"
	peer.statusLock.Unlock()
	log.Printf("Пир %d убит\n", peer.Port)
	server.Close()
}

func monitorPeers() {
	for {
		for _, peer := range peers {
			go func(p *Peer) {
				if p.Status == "мертв" {

					log.Printf("Перезапуск пира %d\n", p.Port)
					go startPeer(p)
				} else {

					conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", p.Port), time.Second)
					if err != nil {
						log.Printf("Пир %d недоступен, помечен как мертвый\n", p.Port)
						p.statusLock.Lock()
						p.Status = "мертв"
						p.statusLock.Unlock()
					} else {
						conn.Close()
					}
				}
			}(peer)
		}
		time.Sleep(5 * time.Second)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при апгрейде WebSocket:", err)
		return
	}
	defer conn.Close()

	for statuses := range statusUpdates {
		if err := conn.WriteJSON(statuses); err != nil {
			log.Println("Ошибка при отправке JSON через WebSocket:", err)
			break
		}
	}
}

func handleKillPeer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	type KillRequest struct {
		Port int `json:"port"`
	}

	var killReq KillRequest
	err := json.NewDecoder(r.Body).Decode(&killReq)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	peer, exists := peers[killReq.Port]
	if !exists {
		http.Error(w, "Пир не найден", http.StatusNotFound)
		return
	}

	go func(p *Peer) {
		p.killChan <- true
	}(peer)

	w.WriteHeader(http.StatusOK)
}
