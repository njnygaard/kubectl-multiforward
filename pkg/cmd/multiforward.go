package cmd

import (
	"fmt"

	"github.com/njnygaard/kubectl-multiforward/forward"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	errGroupNotSpecified = "please specify a group found in your config (for you, one of these choices: \"%s\")"
	errNoConfig          = fmt.Errorf("no config is currently set, either ~/.multiforward.yaml or ./.multiforward.yaml")
	errMalformedConfig   = fmt.Errorf("could not unmarshal config found in either ~/.multiforward.yaml or ./.multiforward.yaml")
	errUnspecified       = fmt.Errorf("something went wrong. do you have your config file in ~/.multiforward.yaml or ./.multiforward.yaml")
)

type Config struct {
	Groups []Group
}

type Group struct {
	Name     string
	Services []Service
}

type Service struct {
	DisplayName string
	LocalPort   uint
	ServicePort uint
	Namespace   string
	Name        string
	Protocol    string
}

// NewCmdNamespace provides a cobra command wrapping NamespaceOptions
func NewCmdNamespace(streams genericclioptions.IOStreams) *cobra.Command {
	// o := NewNamespaceOptions(streams)

	cmd := &cobra.Command{
		Use:   "multiforward <group>",
		Short: "forward to services specified in ~/.multiforward.yaml or .multiforward.yaml",
		Long: `
An example format for your local config file (~/.multiforward.yaml or .multiforward.yaml) can be found here: https://github.com/njnygaard/kubectl-multiforward/blob/master/.multiforward.example.yaml

Or you can reference the following YAML structure where 'displayName' is a human-readable handle and 'name' and 'namespace' identify the Kubernetes resource.

groups:
  - name: production
    services:
      - displayName: Alertmanager
        localPort: 29093
        servicePort: 9093
        namespace: monitoring-prometheus
        name: alertmanager-operated
        protocol: http
      - displayName: Prometheus
        localPort: 29090
        servicePort: 9090
        namespace: monitoring-prometheus
        name: prometheus-operated
        protocol: http
  - name: staging
    services:
      - displayName: Elasticsearch
        localPort: 29200
        servicePort: 9200
        namespace: monitoring-eck
        name: elasticsearch-es-http
        protocol: https
`,
		Example:      "kubectl multiforward [group]",
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) (err error) {

			var config Config

			err = Configure()
			if err != nil {
				return errNoConfig
			}

			err = viper.Unmarshal(&config)
			if err != nil {
				return errMalformedConfig
			}

			var groupNames string
			for i, v := range config.Groups {
				if i == 0 {
					groupNames = v.Name
					continue
				}
				groupNames += ", " + v.Name
			}

			if len(args) == 0 || args[0] == "" {
				c.Help()
				fmt.Println()
				return fmt.Errorf(errGroupNotSpecified, groupNames)
			}

			var found bool
			for _, v := range config.Groups {
				if v.Name == args[0] {
					found = true
				}
			}

			if !found {
				return fmt.Errorf(errGroupNotSpecified, groupNames)
			}

			var serviceGroup []Service

			for _, v := range config.Groups {
				if v.Name == args[0] {
					serviceGroup = v.Services
				}
			}

			var serviceMapping = map[string]forward.ServiceMapping{}

			for _, v := range serviceGroup {
				var mapping forward.ServiceMapping
				mapping.Identifier = v.Name
				mapping.Namespace = v.Namespace
				mapping.ServicePort = int(v.ServicePort)
				mapping.LocalPort = int(v.LocalPort)
				mapping.Protocol = v.Protocol
				serviceMapping[v.DisplayName] = mapping
			}

			forward.Forward(serviceMapping)

			return nil
		},
	}

	cmd.Flags().BoolP("help", "h", false, fmt.Sprintf("help for %s, contains config locations and sample configuration", cmd.Name()))

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
