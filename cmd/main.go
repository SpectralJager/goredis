package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/SpectralJager/goredis"
)

func main() {
	store := map[string][]byte{}
	hstore := map[string]map[string][]byte{}
	server := goredis.NewServer()
	server.Command("ping", func(ctx goredis.Context) goredis.Value {
		args := ctx.Args()
		if len(args) == 0 {
			return goredis.StringValue("PONG")
		}
		return goredis.StringValue(args[0].Bulk())
	})
	server.Command("set", func(ctx goredis.Context) goredis.Value {
		args := ctx.Args()
		if len(args) != 2 {
			return goredis.ErrorValue(fmt.Errorf("ERR wrong number of arguments: expect 2, got %d", len(args)))
		}
		key := args[0].Bulk()
		value := args[1].Marshall()
		store[key] = value
		return goredis.StringValue("OK")
	})
	server.Command("get", func(ctx goredis.Context) goredis.Value {
		args := ctx.Args()
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
	server.Command("hset", func(ctx goredis.Context) goredis.Value {
		args := ctx.Args()
		if len(args) != 3 {
			return goredis.ErrorValue(fmt.Errorf("ERR wrong number of arguments: expect 3, got %d", len(args)))
		}
		hash := args[0].Bulk()
		key := args[1].Bulk()
		val := args[2].Marshall()
		if _, ok := hstore[hash]; !ok {
			hstore[hash] = map[string][]byte{}
		}
		hstore[hash][key] = val
		return goredis.StringValue("OK")
	})
	server.Command("hget", func(ctx goredis.Context) goredis.Value {
		args := ctx.Args()
		if len(args) != 2 {
			return goredis.ErrorValue(fmt.Errorf("ERR wrong number of arguments: expect 2, got %d", len(args)))
		}
		hash := args[0].Bulk()
		key := args[1].Bulk()

		store, ok := hstore[hash]
		if !ok {
			return goredis.ErrorValue(fmt.Errorf("ERR haven't store for hash: %s", hash))
		}
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
	server.Command("hgetall", func(ctx goredis.Context) goredis.Value {
		args := ctx.Args()
		if len(args) != 1 {
			return goredis.ErrorValue(fmt.Errorf("ERR wrong number of arguments: expect 1, got %d", len(args)))
		}
		hash := args[0].Bulk()

		store, ok := hstore[hash]
		if !ok {
			return goredis.ErrorValue(fmt.Errorf("ERR haven't store for hash: %s", hash))
		}
		values := []goredis.Value{}
		for key, content := range store {
			val, err := goredis.NewResp(bytes.NewReader(content)).Read()
			if err != nil {
				return goredis.ErrorValue(err)
			}
			values = append(values, goredis.BulkValue(key), val)
		}
		return goredis.ArrayValue(values...)
	})
	log.Println(server.Start(":6379"))
}
