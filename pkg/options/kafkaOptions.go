package options

import "github.com/Dieg657/Poc.Go.Kafka/pkg/enums"

type KafkaOptions struct {
	Brokers           string                 `json:"brokers"`
	GroupId           string                 `json:"groupId"`
	Offset            enums.AutoOffsetReset  `json:"offset"`
	UserName          string                 `json:"username"`
	Password          string                 `json:"password"`
	SecurityProtocol  enums.SecurityProtocol `json:"securityProtocol"`
	SaslMechanism     enums.SaslMechanims    `json:"saslMechanism"`
	SchemaRegistry    SchemaRegistryOptions  `json:"schemaRegistry"`
	EnableIdempotence bool                   `json:"enableIdempotence"`
	RequestTimeout    int                    `json:"timeout"`
}

type SchemaRegistryOptions struct {
	Url                        string                           `json:"url"`
	BasicAuthUser              string                           `json:"basicAuthUser"`
	BasicAuthSecret            string                           `json:"basicAuthSecret"`
	AutoRegisterSchemas        bool                             `json:"autoRegisterSchemas"`
	RequestTimeout             int                              `json:"timeout"`
	BasicAuthCredentialsSource enums.BasicAuthCredentialsSource `json:"basicAuthCredentialsSource"`
}
