package sockets

import (
	"log"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/googollee/go-socket.io"
	"github.com/thesocialapp/conversation_ai/backend/go/config"
)

const (
	MaxUnackedMessages = 1000
	MessageExpiration  = time.Minute * 5
	ChannelBufferSize  = 100
	MaxRetries         = 3
)

type Subscriber struct {
	sendChan chan Message
	mu       sync.Mutex
	rdb      *redis.Client
	server   *go_socket_io.Server
}

func NewSubscriber(server *go_socket_io.Server, rdb *redis.Client) *Subscriber {
	return &Subscriber{
		sendChan: make(chan Message, ChannelBufferSize),
		server:   server,
		rdb:      rdb,
	}
}

func (s *Subscriber) InitializeHandlers() {
	s.server.OnConnect("/", func(so go_socket_io.Conn) error {
		so.SetContext("")
		log.Println("connected:", so.ID())
		return nil
	})

	s.server.OnEvent("/", "transcript_received", func(so go_socket_io.Conn, transcript string) string {
		log.Println("transcript received:", transcript)

		cfg := config.LoadConfig()
		gptResponse, err := sendTextToGPT(transcript, cfg.GPTAPIURL)
		if err != nil {
			log.Printf("Error processing text with GPT-3: %v", err)
			return "Error processing input"
		}

		return gptResponse.Text
	})

	s.server.OnDisconnect("/", func(so go_socket_io.Conn, reason string) {
		log.Println("closed", reason)
	})
}

func CreateSocketServer(s *Subscriber) {
	http.Handle("/socket.io/", s.server)
	log.Println("Serving socket server at localhost:8000...")
}
