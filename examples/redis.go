package main

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	"github.com/SpectralJager/resp"
)

func main() {
	store := map[string][]byte{}
	hstore := map[string]map[string][]byte{}
	rwx := sync.Mutex{}
	hrwx := sync.Mutex{}

	server := resp.NewServer()

	server.Command("COMMAND", func(ctx resp.Context) resp.Value {
		return resp.StringValue("OK")
	})

	server.Command("ping", func(ctx resp.Context) resp.Value {
		args := ctx.Args()
		if len(args) == 0 {
			return resp.StringValue("PONG")
		}
		return resp.StringValue(args[0].Bulk())
	})

	server.Command("set", func(ctx resp.Context) resp.Value {
		args := ctx.Args()
		if len(args) != 2 {
			return resp.ErrorValue(fmt.Errorf("ERR wrong number of arguments: expect 2, got %d", len(args)))
		}
		rwx.Lock()
		defer rwx.Unlock()
		key := args[0].Bulk()
		value := args[1].Marshall()
		store[key] = value
		return resp.StringValue("OK")
	})

	server.Command("get", func(ctx resp.Context) resp.Value {
		args := ctx.Args()
		if len(args) != 1 {
			return resp.ErrorValue(fmt.Errorf("ERR wrong number of arguments: expect 1, got %d", len(args)))
		}
		rwx.Lock()
		defer rwx.Unlock()
		key := args[0].Bulk()
		content, ok := store[key]
		if !ok {
			return resp.ErrorValue(fmt.Errorf("ERR haven't data for key: %s", key))
		}
		val, err := resp.NewResp(bytes.NewReader(content)).Read()
		if err != nil {
			return resp.ErrorValue(err)
		}
		return val
	})

	server.Command("hset", func(ctx resp.Context) resp.Value {
		args := ctx.Args()
		if len(args) != 3 {
			return resp.ErrorValue(fmt.Errorf("ERR wrong number of arguments: expect 3, got %d", len(args)))
		}
		hrwx.Lock()
		defer hrwx.Unlock()
		hash := args[0].Bulk()
		key := args[1].Bulk()
		val := args[2].Marshall()
		if _, ok := hstore[hash]; !ok {
			hstore[hash] = map[string][]byte{}
		}
		hstore[hash][key] = val
		return resp.StringValue("OK")
	})

	server.Command("hget", func(ctx resp.Context) resp.Value {
		args := ctx.Args()
		if len(args) != 2 {
			return resp.ErrorValue(fmt.Errorf("ERR wrong number of arguments: expect 2, got %d", len(args)))
		}
		hrwx.Lock()
		defer hrwx.Unlock()
		hash := args[0].Bulk()
		key := args[1].Bulk()
		store, ok := hstore[hash]
		if !ok {
			return resp.ErrorValue(fmt.Errorf("ERR haven't store for hash: %s", hash))
		}
		content, ok := store[key]
		if !ok {
			return resp.ErrorValue(fmt.Errorf("ERR haven't data for key: %s", key))
		}
		val, err := resp.NewResp(bytes.NewReader(content)).Read()
		if err != nil {
			return resp.ErrorValue(err)
		}
		return val
	})

	server.Command("hgetall", func(ctx resp.Context) resp.Value {
		args := ctx.Args()
		if len(args) != 1 {
			return resp.ErrorValue(fmt.Errorf("ERR wrong number of arguments: expect 1, got %d", len(args)))
		}
		hrwx.Lock()
		defer hrwx.Unlock()
		hash := args[0].Bulk()
		store, ok := hstore[hash]
		if !ok {
			return resp.ErrorValue(fmt.Errorf("ERR haven't store for hash: %s", hash))
		}
		values := []resp.Value{}
		for key, content := range store {
			val, err := resp.NewResp(bytes.NewReader(content)).Read()
			if err != nil {
				return resp.ErrorValue(err)
			}
			values = append(values, resp.BulkValue(key), val)
		}
		return resp.ArrayValue(values...)
	})

	log.Println(server.Start(":6379"))
}
