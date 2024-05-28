package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main(){
	var n *maelstrom.Node
	n = maelstrom.NewNode()

	n.Handle("echo", func(msg maelstrom.Message) error {
		// Unmarshal the message body as a loosely-typed map.
		var bodley map[string]any
		if err := json.Unmarshal(msg.Body, &bodley); err != nil {
			return err
		}
		// Update the message type to return back.
		bodley["type"] = "echo_ok"

		// Echo the original message back with the updated message type.
		return n.Reply(msg, bodley)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}