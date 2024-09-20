package config

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}
type InventorySrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}
type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}
type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}
type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}
type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type ServerConfig struct {
	Name             string             `mapstructure:"name" json:"name"`
	Tags             []string           `mapstructure:"tags" json:"tags"`
	MysqlInfo        MysqlConfig        `mapstructure:"mysql" json:"mysql"`
	ConsulInfo       ConsulConfig       `mapstructure:"consul" json:"consul"`
	RedisInfo        RedisConfig        `mapstructure:"redis" json:"redis"`
	Host             string             `mapstructure:"host" json:"host"`
	GoodsSrvInfo     GoodsSrvConfig     `mapstructure:"goods_srv" json:"goods_srv"`
	InventorySrvInfo InventorySrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
	JaegerInfo       JaegerConfig       `mapstructure:"jaeger" json:"jaeger"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      uint64 `mapstructure:"port" json:"port"`
	DataId    string `mapstructure:"dataid" json:"dataid"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
	Group     string `mapstructure:"group" json:"group"`
}
