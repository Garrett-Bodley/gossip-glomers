package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"github.com/google/uuid"
)

// Your node will receive a request message body that looks like this:
// {
//   "type": "generate"
// }

// It will need to return a "generate_ok" message with a unique ID:
// {
//   "type": "generate_ok",
//   "id": 123
// }


func main(){
	var n *maelstrom.Node
	n = maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		var bodley map[string]any
		if err := json.Unmarshal(msg.Body, &bodley); err != nil {
			return err
		}

		bodley["type"] = "generate_ok"
		bodley["id"] = uuid.New()

		return n.Reply(msg, bodley)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}