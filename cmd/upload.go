package cmd

import (
	"fmt"
	"github.com/begulsen/golang-sftp/pkg/ftpManager"
	"time"

	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "A brief description of your command",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println("Sftp Upload started ", time.Now().String())

		ftpClient, err := ftpManager.NewConn(flags.host, flags.username, flags.password, flags.port)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = ftpClient.Put("test.xml", "deneme.xml")
		if err != nil {
			fmt.Println(err)
			return err
		} else {
			fmt.Println("Sftp Upload finished ", time.Now().String())
		}
		return nil
	},
}

var flags struct {
	host     string
	username string
	password string
	port     int
	verbose  bool
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.PersistentFlags().StringVarP(&flags.host, "hostname", "H", "localhost", "Sftp Service hostname")
	uploadCmd.PersistentFlags().StringVarP(&flags.username, "username", "u", "username", "Sftp Service username")
	uploadCmd.PersistentFlags().StringVarP(&flags.password, "password", "p", "password", "Sftp Service password")
	uploadCmd.PersistentFlags().IntVarP(&flags.port, "port", "P", 443, "Sftp Service port")
}
