package entities

type Node struct {
	Flag      string
	Memo      string
	Host      string
	PortRange [2]int
	Protocol  string // http | socks5
	User      string
	Pass      string
}
