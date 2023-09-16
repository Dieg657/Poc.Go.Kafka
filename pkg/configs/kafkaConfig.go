package config

type AutoOffsetReset string

const (
	Earliest AutoOffsetReset = "earliest"
	Latest   AutoOffsetReset = "latest"
	Error    AutoOffsetReset = "error"
)

type KafkaConfig struct {
	Brokers        string
	GroupId        string
	Offset         AutoOffsetReset
	UserName       string
	Password       string
	SchemaRegistry SchemaRegistryConfig
}

type SchemaRegistryConfig struct {
	Url                string
	BasicAuthUser      string
	BasicAuthSecret    string
	AutoRegisterSchema bool
}
