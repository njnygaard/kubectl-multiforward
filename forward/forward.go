package forward

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

/*************************************/
/**** Inspiration ********************/
/*************************************/
// https://github.com/gianarb/kube-port-forward

type PortForwardAServiceRequest struct {
	RestConfig  *rest.Config
	Service     v1.Service
	LocalPort   int
	ServicePort int
	Streams     genericclioptions.IOStreams
	StopCh      <-chan struct{}
	ReadyCh     chan struct{}
}

type ServiceMapping struct {
	Port       int
	Namespace  string
	Identifier string
	Protocol   string
}

func Forward(services map[string]ServiceMapping) {

	config, err := config()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	stream := genericclioptions.IOStreams{
		// Typically the forwarding is noisy, we can quiet that here.
		// In:     os.Stdin,
		// Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	// Create a signal sniffing channel and notify on SIGINT or SIGTERM
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	var syncGroup sync.WaitGroup
	syncGroup.Add(1)
	go func() {
		// sig := <-sigs
		<-sigs
		fmt.Println()
		fmt.Print("Bye...")
		syncGroup.Done()
		os.Exit(0)

	}()

	var readyChannels []chan struct{}
	var servicesNotFound []string

	for _, mapping := range services {

		readyCh := make(chan struct{})
		readyChannels = append(readyChannels, readyCh)
		stopCh := make(<-chan struct{})

		go func(srv string, m ServiceMapping, r chan struct{}, s <-chan struct{}) (err error) {

			var clientset *kubernetes.Clientset
			var opts metav1.GetOptions
			var svc *v1.Service

			clientset, err = kubernetes.NewForConfig(config)
			if err != nil {
				return
			}

			svc, err = clientset.CoreV1().Services(m.Namespace).Get(context.TODO(), srv, opts)
			if err != nil {
				fmt.Println(err)
				close(r)
				servicesNotFound = append(servicesNotFound, srv)
				return
			}

			err = PortForwardAService(config, PortForwardAServiceRequest{
				RestConfig:  config,
				Service:     *svc,
				LocalPort:   m.Port,
				ServicePort: m.Port,
				Streams:     stream,
				StopCh:      s,
				ReadyCh:     r,
			})
			if err != nil {
				fmt.Println(err)
				return
			}

			return

		}(mapping.Identifier, mapping, readyCh, stopCh)

	}

	for i := range readyChannels {
		for range readyChannels[i] {
			<-readyChannels[i]
		}
	}

	printTable(services, servicesNotFound)
	syncGroup.Wait()

}

func printTable(services map[string]ServiceMapping, servicesNotFound []string) {
	t := table.NewWriter()

	for service, mapping := range services {
		if !contains(servicesNotFound, mapping.Identifier) {
			t.AppendRow(table.Row{service, fmt.Sprintf("%s://localhost:%d", mapping.Protocol, mapping.Port)})
		}
	}

	t.SetStyle(table.StyleLight)
	colorBOnW := text.Colors{text.BgWhite, text.FgBlack}

	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Service", Colors: text.Colors{text.FgYellow}, ColorsHeader: colorBOnW},
		{Name: "URL", Colors: text.Colors{text.FgHiRed}, ColorsHeader: colorBOnW},
	})

	t.SetCaption("Monitoring Resources... ^C to exit\n")
	fmt.Println(t.Render())

	t.SetColumnConfigs([]table.ColumnConfig{})
}

// You have to portforward directly to a pod, the service is the abstraction.
func PortForwardAService(config *rest.Config, req PortForwardAServiceRequest) (err error) {

	var podName string
	var pods *v1.PodList
	var listOptions metav1.ListOptions
	var set labels.Set
	var dialer httpstream.Dialer
	var fw *portforward.PortForwarder
	var path string
	var hostIP string
	var clientset *kubernetes.Clientset

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return
	}

	set = labels.Set(req.Service.Spec.Selector)
	listOptions = metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	pods, err = clientset.CoreV1().Pods(req.Service.Namespace).List(context.TODO(), listOptions)

	for _, pod := range pods.Items {
		podName = pod.Name
	}

	if podName == "" || err != nil {
		err = fmt.Errorf("could not locate pod for service: %s", req.Service.Name)
		return
	}

	path = fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", req.Service.Namespace, podName)
	hostIP = strings.TrimLeft(req.RestConfig.Host, "htps:/")

	transport, upgrader, err := spdy.RoundTripperFor(req.RestConfig)
	if err != nil {
		return err
	}

	dialer = spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})
	fw, err = portforward.New(dialer, []string{fmt.Sprintf("%d:%d", req.LocalPort, req.ServicePort)}, req.StopCh, req.ReadyCh, req.Streams.Out, req.Streams.ErrOut)
	if err != nil {
		return err
	}

	return fw.ForwardPorts()
}

/*************************************/
/**** Helpers and Initialization *****/
/*************************************/

func config() (*rest.Config, error) {

	var kubeconfig *string

	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	// use the current context in kubeconfig
	return clientcmd.BuildConfigFromFlags("", *kubeconfig)
}

func homeDir() string {

	if h := os.Getenv("HOME"); h != "" {
		return h
	}

	return os.Getenv("USERPROFILE") // windows
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
