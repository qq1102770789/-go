package config

type GoodsSrvConfig struct {
	Host string `mapstructure:"host" json:"host"` //viper
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}
type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}
type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"` //viper
	Port int    `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name        string         `mapstructure:"name" json:"name"`
	Host        string         `mapstructure:"host" json:"host"`
	Tags        []string       `mapstructure:"tags" json:"tags"`
	Port        int            `mapstructure:"port" json:"port"`
	UserSrvInfo GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	JWTInfo     JWTConfig      `mapstructure:"jwt" json:"jwt"`
	ConsulInfo  ConsulConfig   `mapstructure:"consul" json:"consul"`
	JaegerInfo  JaegerConfig   `mapstructure:"jaeger" json:"jaeger"`
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
