package ftpManager

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const sshTurstedKey = "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBHHpFhk/28mT/rtxi77vz/NY35GCHoZwi7cAImtQITsy0uQRWH2PmhNwXRzzaZTtDFChzjE9BZwT7CtSGDt549M="

type sftpClient struct {
	host, user, password string
	port                 int
	*sftp.Client
}

// Create a new SFTP connection by given parameters
func NewConn(host, user, password string, port int) (client *sftpClient, err error) {
	switch {
	case `` == strings.TrimSpace(host),
		`` == strings.TrimSpace(user),
		`` == strings.TrimSpace(password),
		0 >= port || port > 65535:
		return nil, errors.New("Invalid parameters")
	}

	client = &sftpClient{
		host:     host,
		user:     user,
		password: password,
		port:     port,
	}

	if err = client.connect(); nil != err {
		return nil, err
	}
	return client, nil
}

// Upload file to sftp server
func (sc *sftpClient) Put(localFile, remoteFile string) (err error) {
	srcFile, err := os.Open(localFile)
	if err != nil {
		return
	}
	defer srcFile.Close()

	// Make remote directories recursion
	parent := filepath.Dir(remoteFile)
	path := string(filepath.Separator)
	dirs := strings.Split(parent, path)
	for _, dir := range dirs {
		path = filepath.Join(path, dir)
		_ = sc.Mkdir(path)
	}

	dstFile, err := sc.Create(remoteFile)
	if err != nil {
		return
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return
}

func (sc *sftpClient) connect() (err error) {
	config := &ssh.ClientConfig{
		User:            sc.user,
		Auth:            []ssh.AuthMethod{ssh.Password(sc.password)},
		Timeout:         30 * time.Second,
		HostKeyCallback: trustedHostKeyCallback(sshTurstedKey),
	}

	// connet to ssh
	addr := fmt.Sprintf("%s:%d", sc.host, sc.port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}

	// create sftp client
	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	sc.Client = client

	return nil
}

// SSH Key-strings
func trustedHostKeyCallback(trustedKey string) ssh.HostKeyCallback {

	if trustedKey == "" {
		return func(_ string, _ net.Addr, k ssh.PublicKey) error {
			log.Printf("WARNING: SSH-key verification is *NOT* in effect: to fix, add this trustedKey: %q", keyString(k))
			return nil
		}
	}

	return func(_ string, _ net.Addr, k ssh.PublicKey) error {
		ks := keyString(k)
		if trustedKey != ks {
			return fmt.Errorf("SSH-key verification: expected %q but got %q", trustedKey, ks)
		}

		return nil
	}
}

func keyString(k ssh.PublicKey) string {
	return k.Type() + " " + base64.StdEncoding.EncodeToString(k.Marshal())
}
