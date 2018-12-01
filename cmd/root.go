package cmd

import (
	"os"
    "fmt"
    "github.com/spf13/cobra"
    impl "installer/imp"
)

var name string
var configPath string
var hostPkiPath string

var RootCmd = &cobra.Command{
	Use:   "install",
	Short: "Install is a very sample k8s installer",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(name) == 0 {
			cmd.Help()
			return
		}else{
			impl.Run(name,configPath,hostPkiPath)
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


func init(){
	RootCmd.Flags().StringVarP(&name, "name", "n", "", "current node  name")
	RootCmd.Flags().StringVarP(&configPath, "config", "c", "/root/installConfig/config.yml", "config file path")
	RootCmd.Flags().StringVarP(&hostPkiPath, "host-pki-path", "p","/etc/kubernetes/pki", "config file path")
}