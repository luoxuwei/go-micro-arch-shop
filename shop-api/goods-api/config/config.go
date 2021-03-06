package config

type SentinelConfig struct {
	Resource string `mapstructure:"resource" json:"resource"`
	Strategy int32    `mapstructure:"strategy" json:"strategy"`
	Behavior int32 `mapstructure:"behavior" json:"behavior"`
	Threshold float64 `mapstructure:"threshold" json:"threshold"`
	Interval uint32 `mapstructure:"interval" json:"interval"`
}

type JaegerConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type GoodsSrvConfig struct {
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
	Name        string        `mapstructure:"name" json:"name"`
	Host        string        `mapstructure:"host" json:"host"`
	Tags        []string      `mapstructure:"tags" json:"tags"`
	Port        int           `mapstructure:"port" json:"port"`
	JWTInfo     JWTConfig     `mapstructure:"jwt" json:"jwt"`
	GoodsSrvInfo GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`
	RedisInfo   RedisConfig   `mapstructure:"redis" json:"redis"`
	JaegerInfo  JaegerConfig   `mapstructure:"jaeger" json:"jaeger"`
	SentinelInfo SentinelConfig `mapstructure:"sentinel" json:"sentinel"`
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
