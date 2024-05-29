package main

import (
	"encoding/json"
	"log"
	"sync"
	"strings"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// Your node will receive a request message body that looks like this:
// {
//   "type": "broadcast",
//   "message": 1000
// }

// It should store the "message" value locally so it can be read later. In response, it should send an acknowledge with a broadcast_ok message:
// {
//   "type": "broadcast_ok"
// }

func main() {
	n := maelstrom.NewNode()
	store := map[int]bool{}
	mu := sync.Mutex{}
	topology := make(map[string][]string)
	n.Handle("broadcast", func(msg maelstrom.Message) error {

		var bodley struct {
			Type    string
			Message int
		}
		if err := json.Unmarshal(msg.Body, &bodley); err != nil {
			return err
		}
		mu.Lock()
		store[bodley.Message] = true
		mu.Unlock()

		// TODO:
		// [x] Fix concurrency issue
		// [x] Blow up github repo
		// [ ] Reduce unnecessary network chatter (gossip)


		whoami := msg.Dest
		from := msg.Src
		for _, node := range n.NodeIDs() {
			if whoami == node || from == node || strings.HasPrefix(from, "n") { continue }

			n.Send(node, msg.Body)
		}
		// for neighbor := range topology[whoami] {
		// 	n.Send(neighbor, msg)
		// }

		resp := map[string]any{}
		resp["type"] = "broadcast_ok"

		return n.Reply(msg, resp)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var nakamata struct {
			Type string
		}
		if err := json.Unmarshal(msg.Body, &nakamata); err != nil {
			return err
		}

		messages := []int{}

		for key := range store {
			messages = append(messages, key)
		}

		resp := map[string]any{}
		resp["type"] = "read_ok"
		resp["messages"] = messages

		return n.Reply(msg, resp)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {

		var peter struct {
			Type string
			Topology map[string][]string
		}
		if err := json.Unmarshal(msg.Body, &peter); err != nil {
			return err
		}

		for key, value := range peter.Topology {
			topology[key] = value
		}

		resp := map[string]any{}
		resp["type"] = "topology_ok"

		return n.Reply(msg, resp)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
