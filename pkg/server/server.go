package server

type HubServer interface {
	Run()
}

type TunnelServer interface {
	Run() error
}