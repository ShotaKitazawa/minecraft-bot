package rcon

type RconClient interface {
	List() ([]string, error)
	WhitelistAdd(string) error
	WhitelistRemove(string) error
	WhitelistList() ([]string, error)
	DataGetEntity(string) (*User, error)
	Title(string) ([]string, error)
}
