package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	voipms "github.com/ticpu/voipms-gorest/v1"
)

type options struct {
	Username string
	ApiKey   string
	ApiUrl   string
}

var opts options
var vms *voipms.VoIpMsApi

func main() {
	rootCmd := &cobra.Command{
		Use:   "voipms",
		Short: "CLI for VoIP.ms API",
		Run:   help,
	}

	opts.Username = os.Getenv("VOIPMS_USERNAME")
	opts.ApiKey = os.Getenv("VOIPMS_API_KEY")
	opts.ApiUrl = os.Getenv("VOIPMS_API_URL")

	rootCmd.PersistentFlags().StringVarP(&opts.Username, "username", "u", opts.Username, "VoIP.ms account email address")
	rootCmd.PersistentFlags().StringVarP(&opts.ApiKey, "api-key", "p", opts.ApiKey, "VoIP.ms API key")
	rootCmd.PersistentFlags().StringVar(&opts.ApiUrl, "api-url", opts.ApiUrl, "VoIP.ms API URL")

	if len(opts.Username) == 0 || len(opts.ApiKey) == 0 {
		log.Fatalln("username and API key are both required")
	}

	if opts.ApiUrl != "" {
		vms = voipms.NewVoIpMsClientWithUrl(opts.Username, opts.ApiKey, opts.ApiUrl)
	} else {
		vms = voipms.NewVoIpMsClient(opts.Username, opts.ApiKey)
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "setdidpop DID POP",
		Short: "Set the pop value for a DID",
		Args:  cobra.ExactArgs(2),
		Run:   setDidPop,
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "getserversinfo [POP]",
		Short: "Get information about VoIP.ms servers",
		Run:   getServersInfo,
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "getregistrationstatus DID",
		Short: "Get registration status for a DID",
		Args:  cobra.ExactArgs(1),
		Run:   getRegistrationStatus,
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func help(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func setDidPop(cmd *cobra.Command, args []string) {
	did := args[0]
	pop := args[1]
	if response, err := vms.SetDidPopByHostname(did, pop); err != nil {
		log.Fatalf("error setting pop to %v, %v", pop, err)
	} else {
		log.Printf("success: %v", response.Success)
	}
}

func printServerInfo(server *voipms.ServerInfo) {
	recommended := ""

	if server.ServerRecommended == true {
		recommended = "✅"
	}

	fmt.Printf("%d %s (%s : %s) %s\n", server.ServerPOP, server.ServerName, server.ServerHostname, server.ServerIP, recommended)
}

func getServersInfo(cmd *cobra.Command, args []string) {

	if len(args) == 1 {
		var (
			server   *voipms.ServerInfo
			popName  = args[0]
			pop, err = strconv.Atoi(popName)
		)

		if err == nil {
			server, err = vms.GetServersInfoForPop(pop)
		} else {
			server, err = vms.GetServersInfoForPopHostname(popName)
		}

		printServerInfo(server)
	} else {
		var server voipms.ServerInfo

		if response, err := vms.GetServersInfo(); err != nil {
			log.Fatalf("error getting servers info %v", err)
		} else {
			if len(response.Servers) == 0 {
				fmt.Println("no server listed")
			} else {
				for _, server = range response.Servers {
					printServerInfo(&server)
				}
			}
		}
	}
}

func getRegistrationStatus(cmd *cobra.Command, args []string) {
	did := args[0]
	if response, err := vms.GetRegistrationStatus(did); err != nil {
		log.Fatalf("error getting registration status %v", err)
	} else {
		log.Printf("success: %v", response)
	}
}
