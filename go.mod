// This is a generated file. Do not edit directly.

module github.com/njnygaard/kubectl-multiforward

go 1.16

replace (
	k8s.io/api => k8s.io/api v0.0.0-20210825040442-f20796d02069
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20210825040238-74be3b88bedb
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20210825042947-c992623183f8
	k8s.io/client-go => k8s.io/client-go v0.0.0-20210825040738-3dc80a3333cd
)

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/go-openapi/strfmt v0.20.2 // indirect
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	k8s.io/api v0.0.0-20210825040442-f20796d02069
	k8s.io/apimachinery v0.0.0-20210825040238-74be3b88bedb
	k8s.io/cli-runtime v0.22.1
	k8s.io/client-go v1.5.2
)
