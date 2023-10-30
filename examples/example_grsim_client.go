package examples

import "github.com/LiU-SeeGoals/controller/internal/client"

func GrsimClientExample() {
	c := client.NewSSLGrsimClient("127.0.0.1:20011")
	c.Connect()

	c.AddRobotCommand(false, 1, 1.0, 1.0, 0.0, 0.0, 0.0, false, false)
	c.Send()
}
