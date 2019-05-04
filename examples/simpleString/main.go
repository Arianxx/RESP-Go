package main

import (
	"fmt"
	"github.com/arianxx/RESP-Go"
	"github.com/arianxx/camellia-io"
)

func handler(cmd *resp.Command, conn *camellia.Conn, res *[]byte) {
	if cmd == resp.Nil {
		*res = resp.AppendError(*res, []byte("--err nil"))
		return
	}

	fmt.Println("Type:", cmd.Type, "Msg:", cmd.Raw)

	switch string(cmd.Args[0]) {
	case "echo":
		if len(cmd.Args) != 2 {
			*res = resp.AppendError(*res, []byte("--err error arg length"))
		} else {
			*res = resp.AppendSimpleString(*res, cmd.Args[1])
		}
	default:
		*res = resp.AppendError(*res, []byte("--err unknown command"))
	}
}

func main() {
	server, err := resp.NewServer("tcp4", "127.0.0.1:12131", handler)
	if err != nil {
		panic(err)
	}
	if err = server.StartServe(); err != nil {
		panic(err)
	}
}
