Simple redis serialization protocol ([RESP](https://deploy-preview-1964--redis-doc.netlify.app/docs/reference/protocol-spec/)) and it's tcp server implementation.

# Install
```
go get github.com/SpectralJager/resp
```

# Usage
Small redis server
```go
package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/SpectralJager/resp"
)

func main() {
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

	log.Println(server.Start(":6379"))
}
```
