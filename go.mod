// This is a generated file. Do not edit directly.

module githhub.com/njnygaard/kubectl-multiforward

go 1.16

replace (
	k8s.io/api => k8s.io/api v0.0.0-20210825040442-f20796d02069
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20210825040238-74be3b88bedb
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20210825042947-c992623183f8
	k8s.io/client-go => k8s.io/client-go v0.0.0-20210825040738-3dc80a3333cd
)

require (
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	k8s.io/cli-runtime v0.22.1
	k8s.io/client-go v1.5.2
	k8s.io/sample-cli-plugin v0.22.1
)
