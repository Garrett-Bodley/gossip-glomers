package main

import (
	"encoding/json"
	"log"
	"context"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main(){
	n := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(n)
	// counter := 0
	// kv.Write(context.Background(), "count", 0)

	n.Handle("add", func(msg maelstrom.Message) error {
		var nakamata struct {
			Type string
			Delta int
		}

		if err := json.Unmarshal(msg.Body, &nakamata); err != nil {
			return err
		}

		cur_count, err := kv.Read(context.Background(), "count")
		if err != nil { return err }

		kv.Write(context.Background(), "count", cur_count.(int) + nakamata.Delta)
		// counter += nakamata.Delta

		resp := map[string]any{}
		resp["type"] = "add_ok"

		return n.Reply(msg, resp)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		var kang struct {
			Type string
		}

		if err := json.Unmarshal(msg.Body, &kang); err != nil {
			return err
		}

		cur_count, err := kv.Read(context.Background(), "count")
		if err != nil { return err }

		resp := map[string]any{}
		resp["type"] = "read_ok"
		resp["value"] = cur_count

		return n.Reply(msg, resp)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}