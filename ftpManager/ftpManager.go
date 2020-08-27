package ftpManager

import (
	"fmt"
)

const host = "hostName"
const username = "userName"
const password = "password"
const port = 443 // You can use any port except the default port

func UploadSftp() {
	//Create new connection
	ftpClient, err := NewConn(host, username, password, port)
	if err != nil {
		fmt.Println(err)
	}
	err = ftpClient.Put("test.xml", "deneme.xml")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Upload Finished")
	}
}
