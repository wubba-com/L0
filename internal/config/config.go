package cfg

type ConfigDB struct {
	Username, Password, Host, Port, DB string
	MaxAttempts                        int
}
