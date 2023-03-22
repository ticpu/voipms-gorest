package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	voipms "github.com/ticpu/voipms-gorest/v1"
)

type options struct {
	Username   string
	ApiKey     string
	ApiUrl     string
	ApiTimeout time.Duration
}

var opts options
var vms *voipms.VoIpMsApi

func main() {
	var err error
	rootCmd := &cobra.Command{
		Use:   "voipms",
		Short: "CLI for VoIP.ms API",
		Run:   help,
	}

	opts.Username = os.Getenv("VOIPMS_USERNAME")
	opts.ApiKey = os.Getenv("VOIPMS_API_KEY")
	opts.ApiUrl = os.Getenv("VOIPMS_API_URL")
	apiTimeout := os.Getenv("VOIPMS_API_TIMEOUT")
	if len(apiTimeout) > 0 {
		if opts.ApiTimeout, err = time.ParseDuration(apiTimeout); err != nil {
			log.Fatalf("invalid duration %s: %v", apiTimeout, err)
		}
	} else {
		opts.ApiTimeout = 2 * time.Second
	}

	rootCmd.PersistentFlags().StringVarP(&opts.Username, "username", "u", opts.Username, "VoIP.ms account email address")
	rootCmd.PersistentFlags().StringVarP(&opts.ApiKey, "api-key", "p", opts.ApiKey, "VoIP.ms API key")
	rootCmd.PersistentFlags().StringVar(&opts.ApiUrl, "api-url", opts.ApiUrl, "VoIP.ms API URL")
	rootCmd.PersistentFlags().DurationVar(&opts.ApiTimeout, "api-timeout", opts.ApiTimeout, "Timeout for HTTP requests, defaults to 2")

	if len(opts.Username) == 0 || len(opts.ApiKey) == 0 {
		log.Fatalln("username and API key are both required")
	}

	if opts.ApiUrl != "" {
		vms = &voipms.VoIpMsApi{
			ApiUsername: opts.Username,
			ApiPassword: opts.ApiKey,
			ApiUrl:      opts.ApiUrl,
			ApiTimeout:  opts.ApiTimeout,
		}
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
		Use:   "getdidinfo DID [CLIENT]",
		Short: "Get a list of DIDs",
		Args:  cobra.RangeArgs(1, 2),
		Run:   getDidInfo,
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "getalldidsinfo [DID] [CLIENT]",
		Short: "Get a list of DIDs",
		Args:  cobra.NoArgs,
		Run:   getAllDidsInfo,
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
		err error
		did *voipms.DIDInfo
	)

	if len(args) == 2 {
		did, err = vms.GetDidInfo(args[1], args[0])
	} else if len(args) == 1 {
		did, err = vms.GetDidInfo("", args[0])
	}

	if err != nil {
		log.Fatalf("error while fetching did info: %v", err)
	}

	fmt.Printf("%v", *did)
}

func getAllDidsInfo(_ *cobra.Command, _ []string) {
	var (
		err  error
		dids *voipms.GetDidInfoResponse
	)

	dids, err = vms.GetAllDidInfo()

	if err != nil {
		log.Fatalf("error while fetching did info: %v", err)
	}

	fmt.Printf("%v", dids.DIDs)
}

func getDidInfoForClient(_ *cobra.Command, args []string) {
	var (
		err  error
		did  *voipms.DIDInfo
		dids *voipms.GetDidInfoResponse
	)

	if len(args) == 2 {
		if did, err = vms.GetDidInfo(args[0], args[1]); err == nil {
			dids = &voipms.GetDidInfoResponse{DIDs: []voipms.DIDInfo{*did}}
		}
	} else if len(args) == 1 {
		dids, err = vms.GetAllClientDidInfo(args[0])
	} else {
		log.Fatalf("need at least a client for this command: %v", err)
	}

	if err != nil {
		log.Fatalf("error while fetching did info: %v", err)
	}

	fmt.Printf("%v", dids.DIDs)
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
