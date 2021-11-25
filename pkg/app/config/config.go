package config

import "crypto/tls"

type EdgeHubConfig struct {
	HubServerAddr          string
	ProxyServerAddr        string
	SecureProxyServerAddr  string
	DummyProxyServerAddr   string
	DummySecureProxyServerAddr string
	EnableDummyIf          bool
	HubAgentDummyIfName    string
	TLSConfig              *tls.Config
}
