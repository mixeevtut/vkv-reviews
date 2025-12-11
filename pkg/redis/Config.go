package redis

type Settings struct {
	Host     string
	Port     uint16
	User     string
	Pass     string
	Stream   string
	Group    string
	Consumer string
}
