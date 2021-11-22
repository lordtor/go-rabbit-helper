package go_rabbit_helper

type Rabbit struct {
	Server    string  `json:"server" yaml:"server"`
	Host      string  `json:"host" yaml:"host"`
	Port      string  `json:"port" yaml:"port"`
	Username  string  `json:"username" yaml:"username"`
	Password  string  `json:"-" yaml:"password"`
	Data      DataSet `json:"-" yaml:"data" `
	Consumer  bool    `json:"consumer" yaml:"consumer"`
	Publisher bool    `json:"publisher" yaml:"publisher"`
}

type DataSet struct {
	Queue        string `json:"queue" yaml:"queue"`
	Exchange     string `json:"exchange" yaml:"exchange"`
	ExchangeType string `json:"exchange_type" yaml:"exchange_type"`
	RoutingKey   string `json:"routingKey" yaml:"routingKey"`
}
