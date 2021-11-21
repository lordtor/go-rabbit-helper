package go_rabbit_helper

type Rabbit struct {
	Server    string  `yaml:"server"`
	Host      string  `yaml:"host"`
	Port      string  `yaml:"port"`
	Username  string  `yaml:"username"`
	Password  string  `yaml:"password" json:"-"`
	Data      DataSet `yaml:"data" json:"-"`
	Consumer  bool    `yaml:"consumer"`
	Publisher bool    `yaml:"publisher"`
}

type DataSet struct {
	Queue        string `yaml:"queue"`
	Exchange     string `yaml:"exchange"`
	ExchangeType string `yaml:"exchange_type"`
	RoutingKey   string `yaml:"routingKey"`
}
