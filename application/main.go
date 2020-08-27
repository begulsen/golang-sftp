package main

import (
	"../ftpManager"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Sftp Upload started ", time.Now().String())
	ftpManager.UploadSftp()
	fmt.Println("Sftp Upload finished ", time.Now().String())
}
