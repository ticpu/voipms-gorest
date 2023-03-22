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
		Args:  cobra.RangeArgs(0, 1),
		Run:   getServersInfo,
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "getregistrationstatus DID",
		Short: "Get registration status for a DID",
		Args:  cobra.ExactArgs(1),
		Run:   getRegistrationStatus,
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "getclients [CLIENT]",
		Short: "Get a list of clients",
		Args:  cobra.RangeArgs(0, 1),
		Run:   getClients,
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "getdidinfo [DID] [CLIENT]",
		Short: "Get a list of DIDs",
		Args:  cobra.RangeArgs(0, 2),
		Run:   getDidInfo,
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "getdidinfoforclient CLIENT [DID]",
		Short: "Get a list of DIDs for a client",
		Args:  cobra.RangeArgs(1, 2),
		Run:   getDidInfoForClient,
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func help(cmd *cobra.Command, _ []string) {
	_ = cmd.Help()
}

func setDidPop(_ *cobra.Command, args []string) {
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

func getClients(_ *cobra.Command, args []string) {
	var (
		err     error
		clients *voipms.BaseResponse
	)

	if len(args) == 1 {
		clients, err = vms.GetClientOneClient(args[0])
	} else {
		clients, err = vms.GetClients()
	}

	if err != nil {
		log.Fatalf("error while fetching clients: %v", err)
	}

	fmt.Printf("%v", clients)
}

func getDidInfo(_ *cobra.Command, args []string) {
	var (
		err  error
		dids *voipms.GetDidInfoResponse
	)

	if len(args) == 2 {
		dids, err = vms.GetDidInfo(args[1], args[0])
	} else if len(args) == 1 {
		dids, err = vms.GetDidInfo("", args[0])
	} else {
		dids, err = vms.GetAllDidInfo()
	}

	if err != nil {
		log.Fatalf("error while fetching did info: %v", err)
	}

	fmt.Printf("%v", dids)
}

func getDidInfoForClient(_ *cobra.Command, args []string) {
	var (
		err  error
		dids *voipms.GetDidInfoResponse
	)

	if len(args) == 2 {
		dids, err = vms.GetDidInfo(args[0], args[1])
	} else if len(args) == 1 {
		dids, err = vms.GetDidInfo(args[0], "")
	} else {
		log.Fatalf("need at least a client for this command: %v", err)
	}

	if err != nil {
		log.Fatalf("error while fetching did info: %v", err)
	}

	fmt.Printf("%v", dids)
}

func getServersInfo(_ *cobra.Command, args []string) {

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

func getRegistrationStatus(_ *cobra.Command, args []string) {
	did := args[0]
	if response, err := vms.GetRegistrationStatus(did); err != nil {
		log.Fatalf("error getting registration status %v", err)
	} else {
		log.Printf("success: %v", response)
	}
}
