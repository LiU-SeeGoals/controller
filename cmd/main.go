package main

import (
	"fmt"

	"github.com/LiU-SeeGoals/controller/internal/client"
)

func main() {
	fmt.Println("Hello World")

	clnt := client.NewSSLGrsimClient("127.0.0.1:20011")
	clnt.Connect()
	clnt.Send()
}
