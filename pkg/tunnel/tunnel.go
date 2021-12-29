package tunnel

type Tunnel interface {
	Listen() error
}