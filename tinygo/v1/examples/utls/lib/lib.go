package lib

import (
	tls "github.com/refraction-networking/utls"
)

//tinyjson:json
type TLSConfig struct {
	NextProtos                  []string          `json:"next_protos"`
	ApplicationSettings         map[string][]byte `json:"application_settings"`
	ServerName                  string            `json:"server_name"`
	InsecureSkipVerify          bool              `json:"insecure_skip_verify"`
	InsecureSkipTimeVerify      bool              `json:"insecure_skip_time_verify"`
	OmitEmptyPsk                bool              `json:"omit_empty_psk"`
	InsecureServerNameToVerify  string            `json:"insecure_server_name_to_verify"`
	SessionTicketsDisabled      bool              `json:"session_tickets_disabled"`
	PQSignatureSchemesEnabled   bool              `json:"pq_signature_schemes_enabled"`
	DynamicRecordSizingDisabled bool              `json:"dynamic_record_sizing_disabled"`
	ECHConfigs                  []byte            `json:"ech_configs"`
}

//tinyjson:json
type Configurables struct {
	TLSConfig     *TLSConfig `json:"tls_config"`      // will be converted to tls.Config
	ClientHelloID string     `json:"client_hello_id"` // will be converted to tls.ClientHelloID
}

func (c *Configurables) GetTLSConfig() *tls.Config {
	config := &tls.Config{
		NextProtos:                  c.TLSConfig.NextProtos,
		ApplicationSettings:         c.TLSConfig.ApplicationSettings,
		ServerName:                  c.TLSConfig.ServerName,
		InsecureSkipVerify:          c.TLSConfig.InsecureSkipVerify,
		InsecureSkipTimeVerify:      c.TLSConfig.InsecureSkipTimeVerify,
		OmitEmptyPsk:                c.TLSConfig.OmitEmptyPsk,
		InsecureServerNameToVerify:  c.TLSConfig.InsecureServerNameToVerify,
		SessionTicketsDisabled:      c.TLSConfig.SessionTicketsDisabled,
		PQSignatureSchemesEnabled:   c.TLSConfig.PQSignatureSchemesEnabled,
		DynamicRecordSizingDisabled: c.TLSConfig.DynamicRecordSizingDisabled,
	}

	echConfigs, err := tls.UnmarshalECHConfigs(c.TLSConfig.ECHConfigs)
	if err == nil { // otherwise do we need to return an error or just ignore it?
		config.ECHConfigs = echConfigs
	}

	return config
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
	default:
		panic("unknown client hello id")
	}
}
