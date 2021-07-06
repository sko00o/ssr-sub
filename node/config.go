package node

// Config set for ssr config
type Config struct {
	ID            string `json:"_"`
	Server        string `json:"server"`
	ServerPort    int    `json:"server_port"`
	Method        string `json:"method"`
	Protocol      string `json:"protocol"`
	ProtocolParam string `json:"protocol_param"`
	OBFS          string `json:"obfs"`
	OBFSParam     string `json:"obfs_param"`
	Password      string `json:"password"`
	Remarks       string `json:"remarks"`
	Group         string `json:"group"`
}

type CheckConfig struct {
	TCPTimeout string `yaml:"timeout"`
	Not        string `yaml:"not"`
}
