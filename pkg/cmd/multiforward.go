package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	errNoConfig    = fmt.Errorf("no config is currently set, either ~/.multiforward.yaml or ./.multiforward.yaml")
	errUnspecified = fmt.Errorf("something went wrong. do you have your config file in ~/.multiforward.yaml or ./.multiforward.yaml")
)

// NewCmdNamespace provides a cobra command wrapping NamespaceOptions
func NewCmdNamespace(streams genericclioptions.IOStreams) *cobra.Command {
	// o := NewNamespaceOptions(streams)

	cmd := &cobra.Command{
		Use:          "ns [new-namespace] [flags]",
		Short:        "View or set the current namespace",
		Example:      "kubectl multiforward",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) (err error) {

			logger := logrus.New()

			err = Configure()
			if err != nil {
				return errNoConfig
			}

			logger.Info("yeah")

			return nil
		},
	}

	// Might need this
	// cmd.Flags().BoolVar(&o.listNamespaces, "list", o.listNamespaces, "if true, print the list of all namespaces in the current KUBECONFIG")
	// o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

func Configure() (err error) {
	viper.SetConfigName(".multiforward")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return errNoConfig
		} else {
			return errUnspecified
		}
	}

	return
}
