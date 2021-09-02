package main

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"os"

	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/sample-cli-plugin/pkg/cmd"
)

func main() {
	flags := pflag.NewFlagSet("kubectl-ns", pflag.ExitOnError)
	pflag.CommandLine = flags

	root := cmd.NewCmdNamespace(genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
