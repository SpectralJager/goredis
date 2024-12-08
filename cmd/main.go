package main

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/SpectralJager/goredis"
)

var store = map[string][]byte{}

func main() {
	server := goredis.NewServer()
	server.Command("ping", func(ctx context.Context, args []goredis.Value) goredis.Value {
		if len(args) == 0 {
			return goredis.StringValue("PONG")
		}
		return goredis.StringValue(args[0].Bulk())
	})
	server.Command("set", func(ctx context.Context, args []goredis.Value) goredis.Value {
		if len(args) != 2 {
			return goredis.ErrorValue(fmt.Errorf("ERR wrong number of arguments: expect 2, got %d", len(args)))
		}
		key := args[0].Bulk()
		value := args[1].Marshall()
		store[key] = value
		return goredis.StringValue("OK")
	})
	server.Command("get", func(ctx context.Context, args []goredis.Value) goredis.Value {
		if len(args) != 1 {
			return goredis.ErrorValue(fmt.Errorf("ERR wrong number of arguments: expect 1, got %d", len(args)))
		}
		key := args[0].Bulk()
		content, ok := store[key]
		if !ok {
			return goredis.ErrorValue(fmt.Errorf("ERR haven't data for key: %s", key))
		}
		val, err := goredis.NewResp(bytes.NewReader(content)).Read()
		if err != nil {
			return goredis.ErrorValue(err)
		}
		return val
	})
	log.Println(server.Start(":6379"))
}
