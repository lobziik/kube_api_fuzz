package main

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	kcapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

type settings struct {
	apiserverStdout string
	apiserverStderr string

	etcdStdout string
	etcdStderr string

	alsoLogToStdout bool
}

func main() {
	stgs := parseSettings()

	testEnv := &envtest.Environment{
		CRDDirectoryPaths:        []string{},
		ErrorIfCRDPathMissing:    true,
		AttachControlPlaneOutput: stgs.alsoLogToStdout,
		ControlPlane: envtest.ControlPlane{
			APIServer: &envtest.APIServer{
				Out: getWriter(stgs.apiserverStdout, stgs.alsoLogToStdout),
				Err: getWriter(stgs.apiserverStderr, stgs.alsoLogToStdout),
			},
			Etcd: &envtest.Etcd{
				Out: getWriter(stgs.etcdStdout, stgs.alsoLogToStdout),
				Err: getWriter(stgs.etcdStderr, stgs.alsoLogToStdout),
			},
		},
	}
	cfg, err := testEnv.Start()
	defer testEnv.Stop()
	if err != nil {
		panic(err)
	}

	fmt.Println("Envtest started")
	dumpPath, err := dumpCertsAndConfig(cfg)
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dumpPath)

	for {
		fmt.Println(fmt.Sprintf("still running, stuff in: %s, apiserver url: %s", dumpPath, cfg.Host))
		time.Sleep(10 * time.Second)
	}
}

func parseSettings() *settings {
	alsoToStdout := flag.Bool("alsoLogToStdout", false, "forward apiserver and etcd output to stdout")

	apiserverStdout := flag.String("apiserverStdout", "", "forward apiserver stdout to this file")
	apiserverStderr := flag.String("apiserverStderr", "", "forward apiserver stderr to this file")
	etcdStdout := flag.String("etcdStdout", "", "forward etcd stdout to this file")
	etcdStderr := flag.String("etcdStderr", "", "forward etcd stderr to this file")
	flag.Parse()

	stgs := &settings{
		alsoLogToStdout: *alsoToStdout,

		apiserverStdout: *apiserverStdout,
		apiserverStderr: *apiserverStderr,

		etcdStdout: *etcdStdout,
		etcdStderr: *etcdStderr,
	}

	fmt.Println(fmt.Sprintf("settings: %+v", stgs))

	return stgs
}

func getWriter(path string, alsoStdout bool) io.Writer {
	if path != "" {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
		if err != nil {
			panic(fmt.Errorf("can not create file %s: %w", path, err))
		}
		if alsoStdout {
			return io.MultiWriter(os.Stdout, file)
		}
		return file
	}

	if alsoStdout {
		return os.Stdout
	}
	return nil
}

func dumpCertsAndConfig(cfg *rest.Config) (string, error) {
	dir, err := os.MkdirTemp("", "envtest-fuzz-")
	if err != nil {
		return "", err
	}
	fmt.Println("temp dir: ", dir)

	kubeconfig := filepath.Join(dir, "kubeconfig")
	kubeconfigContent, err := kubeConfigFromREST(cfg)
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(kubeconfig, kubeconfigContent, 0666); err != nil {
		return "", err
	}

	clientKey := filepath.Join(dir, "clientKey.pem")
	if err := os.WriteFile(clientKey, cfg.KeyData, 0666); err != nil {
		return "", err
	}

	clientCert := filepath.Join(dir, "clientCert.pem")
	if err := os.WriteFile(clientCert, cfg.CertData, 0666); err != nil {
		return "", err
	}

	return dir, nil
}

// Copied from the envtest internals
func kubeConfigFromREST(cfg *rest.Config) ([]byte, error) {
	const (
		envtestName = "envtest"
	)

	kubeConfig := kcapi.NewConfig()
	protocol := "https"
	if !rest.IsConfigTransportTLS(*cfg) {
		protocol = "http"
	}

	// cfg.Host is a URL, so we need to parse it so we can properly append the API path
	baseURL, err := url.Parse(cfg.Host)
	if err != nil {
		return nil, fmt.Errorf("unable to interpret config's host value as a URL: %w", err)
	}

	kubeConfig.Clusters[envtestName] = &kcapi.Cluster{
		// TODO(directxman12): if client-go ever decides to expose defaultServerUrlFor(config),
		// we can just use that.  Note that this is not the same as the public DefaultServerURL,
		// which requires us to pass a bunch of stuff in manually.
		Server:                   (&url.URL{Scheme: protocol, Host: baseURL.Host, Path: cfg.APIPath}).String(),
		CertificateAuthorityData: cfg.CAData,
	}
	kubeConfig.AuthInfos[envtestName] = &kcapi.AuthInfo{
		// try to cover all auth strategies that aren't plugins
		ClientCertificateData: cfg.CertData,
		ClientKeyData:         cfg.KeyData,
		Token:                 cfg.BearerToken,
		Username:              cfg.Username,
		Password:              cfg.Password,
	}
	kcCtx := kcapi.NewContext()
	kcCtx.Cluster = envtestName
	kcCtx.AuthInfo = envtestName
	kubeConfig.Contexts[envtestName] = kcCtx
	kubeConfig.CurrentContext = envtestName

	contents, err := clientcmd.Write(*kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize kubeconfig file: %w", err)
	}
	return contents, nil
}
