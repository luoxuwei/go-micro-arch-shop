package config


type UserSrvConfig struct {
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

type AliyunConfig struct {
	ApiKey     string `mapstructure:"key" json:"key"`
	ApiSecrect string `mapstructure:"secrect" json:"secrect"`
}

type AliyunSmsConfig struct {
	SignName string `mapstructure:"sign_name" json:"sign_name"`
	TemplateCode string `mapstructure:"template_code" json:"template_code"`
}

type ServerConfig struct {
	Name        string        `mapstructure:"name" json:"name"`
	Host        string        `mapstructure:"host" json:"host"`
	Tags        []string      `mapstructure:"tags" json:"tags"`
	Port        int           `mapstructure:"port" json:"port"`
	JWTInfo     JWTConfig     `mapstructure:"jwt" json:"jwt"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`
	AliyunInfo  AliyunConfig  `mapstructure:"aliyun" json:"aliyun"`
    AliyunSmsInfo AliyunSmsConfig `mapstructure:"sms" json:"sms"`
	RedisInfo   RedisConfig   `mapstructure:"redis" json:"redis"`
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
