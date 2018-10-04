package main

// The PrivOps interface captures the privileged operations that hil-vpnd
// needs to perform. In a production setup, we always do this by calling the
// 'hil-vpn-privop' command (see PrivOpsCmd), but this interface allows us
// to test more easily.
type PrivOps interface {
	CreateVPN(name string, vlanNo uint16, portNo uint16) (string, error)
	StartVPN(name string) error
	StopVPN(name string) error
	DeleteVPN(name string) error
	ListVPNs() ([]string, error)
}

// An implementation of PrivOps that calls the 'hil-vpn-privop' command.
type PrivOpsCmd struct{}

func (PrivOpsCmd) CreateVPN(name string, vlanNo uint16, portNo uint16) (string, error) {
	panic("TODO")
}
func (PrivOpsCmd) StartVPN(name string) error  { panic("TODO") }
func (PrivOpsCmd) StopVPN(name string) error   { panic("TODO") }
func (PrivOpsCmd) DeleteVPN(name string) error { panic("TODO") }
func (PrivOpsCmd) ListVPNs() ([]string, error) { panic("TODO") }
