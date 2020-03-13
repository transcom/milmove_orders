package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/transcom/mymove/pkg/cli"
)

const (
	// CertPathFlag is the path to the certificate to use for TLS
	CertPathFlag string = "certpath"
	// KeyPathFlag is the path to the key to use for TLS
	KeyPathFlag string = "keypath"
	// HostnameFlag is the hostname to connect to
	HostnameFlag string = "hostname"
	// PortFlag is the port to connect to
	PortFlag string = "port"
	// InsecureFlag indicates that TLS verification and validation can be skipped
	InsecureFlag string = "insecure"
)

// initRootFlags initializes flags relating to the orders api
func initRootFlags(flag *pflag.FlagSet) {
	cli.InitCACFlags(flag)
	cli.InitVerboseFlags(flag)

	flag.String(CertPathFlag, "./config/tls/devlocal-mtls.cer", "Path to the public cert")
	flag.String(KeyPathFlag, "./config/tls/devlocal-mtls.key", "Path to the private key")
	flag.String(HostnameFlag, cli.HTTPOrdersServerNameLocal, "The hostname to connect to")
	flag.Int(PortFlag, cli.MutualTLSPort, "The port to connect to")
	flag.Bool(InsecureFlag, false, "Skip TLS verification and validation")
}

// CheckRootConfig checks the validity of the orders api flags
func CheckRootConfig(v *viper.Viper) error {
	err := cli.CheckCAC(v)
	if err != nil {
		return err
	}

	err = cli.CheckVerbose(v)
	if err != nil {
		return err
	}

	if (v.GetString(CertPathFlag) != "" && v.GetString(KeyPathFlag) == "") || (v.GetString(CertPathFlag) == "" && v.GetString(KeyPathFlag) != "") {
		return fmt.Errorf("Both TLS certificate and key paths must be provided")
	}

	return nil
}

func main() {
	root := cobra.Command{
		Use:   "orders-api-client [flags]",
		Short: "Orders API client",
		Long:  "Orders API client",
	}
	initRootFlags(root.PersistentFlags())

	getOrdersCommand := &cobra.Command{
		Use:          "get-orders",
		Short:        "Gets Orders",
		Long:         "Gets Orders and all of its Revisions by UUID",
		RunE:         getOrders,
		SilenceUsage: true,
	}
	initGetOrdersFlags(getOrdersCommand.Flags())
	root.AddCommand(getOrdersCommand)

	postRevisionsCommand := &cobra.Command{
		Use:          "post-revisions",
		Short:        "Post Revisions",
		Long:         "Post Revisions to orders from a CSV file",
		RunE:         postRevisions,
		SilenceUsage: true,
	}
	initPostRevisionsFlags(postRevisionsCommand.Flags())
	root.AddCommand(postRevisionsCommand)

	completionCommand := &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long:  "To install completion scripts run:\n\norders-api-client completion > /usr/local/etc/bash_completion.d/orders-api-client",
		RunE: func(cmd *cobra.Command, args []string) error {
			return root.GenBashCompletion(os.Stdout)
		},
	}
	root.AddCommand(completionCommand)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
