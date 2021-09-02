package cmd

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	namespaceExample = `
	# view the current namespace in your KUBECONFIG
	%[1]s ns

	# view all of the namespaces in use by contexts in your KUBECONFIG
	%[1]s ns --list

	# switch your current-context to one that contains the desired namespace
	%[1]s ns foo
`

	// errNoContext = fmt.Errorf("no context is currently set, use %q to select a new one", "kubectl config use-context <context>")
)

// NamespaceOptions provides information required to update
// the current context on a user's KUBECONFIG
// type NamespaceOptions struct {
// 	configFlags *genericclioptions.ConfigFlags

// 	resultingContext     *api.Context
// 	resultingContextName string

// 	userSpecifiedCluster   string
// 	userSpecifiedContext   string
// 	userSpecifiedAuthInfo  string
// 	userSpecifiedNamespace string

// 	rawConfig      api.Config
// 	listNamespaces bool
// 	args           []string

// 	genericclioptions.IOStreams
// }

// NewNamespaceOptions provides an instance of NamespaceOptions with default values
// func NewNamespaceOptions(streams genericclioptions.IOStreams) *NamespaceOptions {
// 	return &NamespaceOptions{
// 		configFlags: genericclioptions.NewConfigFlags(true),

// 		IOStreams: streams,
// 	}
// }

// NewCmdNamespace provides a cobra command wrapping NamespaceOptions
func NewCmdNamespace(streams genericclioptions.IOStreams) *cobra.Command {
	// o := NewNamespaceOptions(streams)

	cmd := &cobra.Command{
		Use:          "ns [new-namespace] [flags]",
		Short:        "View or set the current namespace",
		Example:      fmt.Sprintf(namespaceExample, "kubectl"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {

			logger := logrus.New()

			logger.Info("yeah")

			return nil
		},
	}

	// Might need this
	// cmd.Flags().BoolVar(&o.listNamespaces, "list", o.listNamespaces, "if true, print the list of all namespaces in the current KUBECONFIG")
	// o.configFlags.AddFlags(cmd.Flags())

	return cmd
}
