package lib

import (
	"crypto/x509"
	"log"
	"os"
	"strings"

	tls "github.com/refraction-networking/utls"
)

//tinyjson:json
type TLSConfig struct {
	NextProtos                  []string          `json:"next_protos"`
	ApplicationSettings         map[string][]byte `json:"application_settings"`
	ServerName                  string            `json:"server_name"`
	InsecureSkipVerify          bool              `json:"insecure_skip_verify"` // if not set, host must supply a root CA certificate via root_ca or root_ca_dirs
	InsecureSkipTimeVerify      bool              `json:"insecure_skip_time_verify"`
	OmitEmptyPsk                bool              `json:"omit_empty_psk"`
	InsecureServerNameToVerify  string            `json:"insecure_server_name_to_verify"`
	SessionTicketsDisabled      bool              `json:"session_tickets_disabled"`
	PQSignatureSchemesEnabled   bool              `json:"pq_signature_schemes_enabled"`
	DynamicRecordSizingDisabled bool              `json:"dynamic_record_sizing_disabled"`
	ECHConfigs                  []byte            `json:"ech_configs"`

	RootCAPath string   `json:"root_ca_path"` // if set, will be parsed as a x509 Root CA certificate
	RootCADirs []string `json:"root_ca_dirs"` // if non-empty, all x509 certs found in the directory specified by the list will be used to verify the host
}

func (tlsConf *TLSConfig) GetConfig() *tls.Config {
	conf := &tls.Config{
		NextProtos:                  tlsConf.NextProtos,
		ApplicationSettings:         tlsConf.ApplicationSettings,
		ServerName:                  tlsConf.ServerName,
		InsecureSkipVerify:          tlsConf.InsecureSkipVerify,
		InsecureSkipTimeVerify:      tlsConf.InsecureSkipTimeVerify,
		OmitEmptyPsk:                tlsConf.OmitEmptyPsk,
		InsecureServerNameToVerify:  tlsConf.InsecureServerNameToVerify,
		SessionTicketsDisabled:      tlsConf.SessionTicketsDisabled,
		PQSignatureSchemesEnabled:   tlsConf.PQSignatureSchemesEnabled,
		DynamicRecordSizingDisabled: tlsConf.DynamicRecordSizingDisabled,
	}

	echConfigs, err := tls.UnmarshalECHConfigs(tlsConf.ECHConfigs)
	if err == nil { // otherwise do we need to return an error or just ignore it?
		conf.ECHConfigs = echConfigs
	}

	if !tlsConf.InsecureSkipVerify {
		rootCAs, err := tlsConf.loadRootCAs()
		if err == nil {
			conf.RootCAs = rootCAs
		} else {
			panic("failed to load root CAs: " + err.Error())
		}
	}

	return conf
}

// loadRootCAs loads the root CA certificates from the RootCA and RootCADirs fields.
//
// Derived from crypto/x509.loadSystemRoots
func (tlsConf *TLSConfig) loadRootCAs() (*x509.CertPool, error) {
	roots := x509.NewCertPool()

	// if RootCA is set, use it as the first cert
	var files []string
	if tlsConf.RootCAPath != "" {
		log.Printf("UTLSClientWrappingTransport: loading root CA certificate set via config: %s\n", tlsConf.RootCAPath)
		files = append(files, tlsConf.RootCAPath)
	}

	if f := os.Getenv(certFileEnv); f != "" {
		log.Printf("UTLSClientWrappingTransport: loading root CA certificate set via ENV: %s\n", tlsConf.RootCAPath)
		files = append(files, f)
	}

	var firstErr error
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err == nil {
			roots.AppendCertsFromPEM(data)
			break
		}
		if firstErr == nil && !os.IsNotExist(err) {
			firstErr = err
		}
	}

	var dirs []string = tlsConf.RootCADirs
	if d := os.Getenv(certDirEnv); d != "" {
		// OpenSSL and BoringSSL both use ":" as the SSL_CERT_DIR separator.
		// See:
		//  * https://golang.org/issue/35325
		//  * https://www.openssl.org/docs/man1.0.2/man1/c_rehash.html
		dirs = strings.Split(d, ":")
	}

	for _, directory := range dirs {
		fis, err := readUniqueDirectoryEntries(directory)
		if err != nil {
			if firstErr == nil && !os.IsNotExist(err) {
				firstErr = err
			}
			continue
		}
		for _, fi := range fis {
			data, err := os.ReadFile(directory + "/" + fi.Name())
			if err == nil {
				roots.AppendCertsFromPEM(data)
			}
		}
	}

	return roots, firstErr
}

//tinyjson:json
type Configurables struct {
	TLSConfig                *TLSConfig `json:"tls_config"`                 // will be converted to tls.Config
	ClientHelloID            string     `json:"client_hello_id"`            // will be converted to tls.ClientHelloID
	InternalBufferSize       int        `json:"internal_buffer_size"`       // will be used to allocate internal temporary buffer
	BackgroundWorkerFairness bool       `json:"background_worker_fairness"` // if true, use fairWorker, otherwise use unfairWorker
}

func (c *Configurables) GetTLSConfig() *tls.Config {
	return c.TLSConfig.GetConfig()
}

func (c *Configurables) GetClientHelloID() tls.ClientHelloID {
	switch c.ClientHelloID {
	case "HelloChrome_Auto", "HelloChrome", "Chrome", "chrome":
		return tls.HelloChrome_Auto
	case "HelloEdge_Auto", "HelloEdge", "Edge", "edge":
		return tls.HelloEdge_Auto
	case "HelloFirefox_Auto", "HelloFirefox", "Firefox", "firefox":
		return tls.HelloFirefox_Auto
	case "HelloSafari_Auto", "HelloSafari", "Safari", "safari":
		return tls.HelloSafari_Auto
	case "HelloRandomized", "Randomized", "randomized", "Random", "random":
		return tls.HelloRandomized
	case "HelloRandomizedALPN", "RandomizedALPN", "randomized_alpn", "RandomALPN", "random_alpn":
		return tls.HelloRandomizedALPN
	case "HelloRandomizedNoALPN", "RandomizedNoALPN", "randomized_no_alpn", "RandomNoALPN", "random_no_alpn":
		return tls.HelloRandomizedNoALPN
	case "HelloGolang", "Golang", "golang", "Go", "go", "crypto/tls", "Default", "default", "":
		return tls.HelloGolang
	default:
		panic("unknown client hello id")
	}
}
