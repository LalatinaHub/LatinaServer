package model

// ProxyNode represents a relay proxy node from external database.
// This is a generic representation that matches the database schema,
// suitable for conversion to various proxy protocols (Shadowsocks, VMess, VLESS, Trojan).
type ProxyNode struct {
	ID          int64
	Server      string
	IP          string
	ServerPort  int
	UUID        string
	Password    string
	Security    string
	AlterID     int
	Method      string
	Plugin      string
	PluginOpts  string
	Host        string
	TLS         bool
	Transport   string
	Path        string
	ServiceName string
	Insecure    bool
	SNI         string
	Remark      string
	ConnMode    string
	CountryCode string
	Region      string
	Org         string
	VPN         string
	Raw         string
}
