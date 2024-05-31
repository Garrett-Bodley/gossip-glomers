package main

import (
	"encoding/json"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main(){
	n := maelstrom.NewNode()
	s := &server{ n: n, gossip: make(map[int]bool) }

	n.Handle("broadcast", s.broadcastHandler)
}

// Pseudo class declaration

type server struct {
	n *maelstrom.Node
	gossip map[int]bool
	neighbors []int
	// maybe we need a mutex here
}

// Class method
func (s *server) broadcastHandler(msg maelstrom.Message) error {
	// Receive a juicy new message
	// Gossip it to our neighbors
	// Respond to the client
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	// we want to add the new message to our gossip map
	new_gossip := int(body["message"].(float64))
	if s.gossip[new_gossip]{
		return nil
	}

	type resp_t struct {
		Type string
		Id int
		Foobar string
	}

	type gossip_t {
		Type string
		Message string
	}

	// Gossip to our neighbors

	// Reply to the client
	resp := resp_t{ Type: "broadcast_ok" }

	s.n.Send(dest, resp)

	return nil
}