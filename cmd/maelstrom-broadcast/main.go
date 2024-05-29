package main

import (
	"encoding/json"
	"log"
	"sync"
	"strings"
	"fmt"
	"context"
	"time"

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

	// TODO:
	// [ ] Handler for broadcast_ok, read_ok ?

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

		whoami := msg.Dest
		from := msg.Src
		for _, node := range n.NodeIDs() {
			if whoami == node || from == node || strings.HasPrefix(from, "n") { continue }
			if strings.HasPrefix(node, "c") { log.Printf("wtf why is this not logging ever the problem is not here: %s\n", node); panic(1) }

			delay := time.Millisecond * 200
			for i := 0; i < 10; i++ {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second * 2)
				defer cancel()

				msg2, err := n.SyncRPC(ctx, node, msg.Body)
				if err != nil {
					log.Printf("Maybe here? %+v\n", err)
					time.Sleep(delay)
					delay *= 2
					log.Printf("Trying again %d\n", i)
					continue
				}

				// overkill but it's fine
				var kaplan struct {
					Type string
					Message int
					MsgId int
				}
				if err := json.Unmarshal(msg2.Body, &kaplan); err != nil {
					return err
				}
				if kaplan.Type != "broadcast_ok" {
					panic(fmt.Sprintf("broadcast is not okay: %s", kaplan.Type))
				}
				break
			}

		}


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
