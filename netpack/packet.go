package netpack

//TODO server password transfer protocol and other message
//Packet
type Packet struct {
	Token     string `json:"token"`
	HostName  string `json:"host_name"`
	OS        string `json:"os"`
	SocksAuth bool   `json:"socks_auth"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
}
