package config


type UseropSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`
	Port   int    `mapstructure:"port" json:"port"`
	Expire int    `mapstructure:"expire" json:"expire"`
}

type ServerConfig struct {
	Name          string          `mapstructure:"name" json:"name"`
	Host          string          `mapstructure:"host" json:"host"`
	Tags          []string        `mapstructure:"tags" json:"tags"`
	Port          int             `mapstructure:"port" json:"port"`
	JWTInfo       JWTConfig       `mapstructure:"jwt" json:"jwt"`
	UseropSrvInfo UseropSrvConfig `mapstructure:"userop_srv" json:"userop_srv"`
	GoodsSrvInfo  UseropSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	ConsulInfo    ConsulConfig    `mapstructure:"consul" json:"consul"`
	RedisInfo     RedisConfig     `mapstructure:"redis" json:"redis"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
