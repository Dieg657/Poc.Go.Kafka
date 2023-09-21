package setup

import (
	"fmt"
	"os"

	"github.com/Dieg657/Poc.Go.Kafka/pkg/options"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
)

type ISchemaRegistrySetup interface {
	New(config *options.KafkaOptions) error
	GetAvroSerializer() *avro.SpecificSerializer
	GetAvroDeserializer() *avro.SpecificDeserializer
	GetJsonSerializer() *jsonschema.Serializer
	GetJsonDeserializer() *jsonschema.Deserializer
}

type SchemaRegistrySetup struct {
	schemaregistry           *schemaregistry.Client
	avroSpecificSerializer   *avro.SpecificSerializer
	avroSpecificDeserializer *avro.SpecificDeserializer
	jsonSerializer           *jsonschema.Serializer
	jsonDeserializer         *jsonschema.Deserializer
}

func (registry *SchemaRegistrySetup) New(options *options.KafkaOptions) error {
	configuration := &schemaregistry.Config{
		SchemaRegistryURL:          options.SchemaRegistry.Url,
		RequestTimeoutMs:           options.RequestTimeout,
		ConnectionTimeoutMs:        5000,
		BasicAuthCredentialsSource: string(options.SchemaRegistry.BasicAuthCredentialsSource),
	}

	if options.SchemaRegistry.BasicAuthUser != "" && options.SchemaRegistry.BasicAuthSecret != "" {
		configuration.BasicAuthUserInfo = fmt.Sprintf("%s:%s", options.SchemaRegistry.BasicAuthUser, options.SchemaRegistry.BasicAuthSecret)
	} else {
		configuration.BasicAuthUserInfo = ""
	}

	schemaRegistry, err := schemaregistry.NewClient(configuration)

	if err != nil {
		return err
	}

	registry.schemaregistry = &schemaRegistry
	registry.setSerializers(options)
	registry.setDeserializers(options)

	return nil
}

func (registry *SchemaRegistrySetup) GetAvroSerializer() *avro.SpecificSerializer {
	return registry.avroSpecificSerializer
}

func (registry *SchemaRegistrySetup) GetAvroDeserializer() *avro.SpecificDeserializer {
	return registry.avroSpecificDeserializer
}

func (registry *SchemaRegistrySetup) GetJsonSerializer() *jsonschema.Serializer {
	return registry.jsonSerializer
}

func (registry *SchemaRegistrySetup) GetJsonDeserializer() *jsonschema.Deserializer {
	return registry.jsonDeserializer
}

func (registry *SchemaRegistrySetup) setSerializers(options *options.KafkaOptions) {
	registry.
		configAvroSerializer(options).
		configJsonSerializer(options)
}

func (registry *SchemaRegistrySetup) setDeserializers(options *options.KafkaOptions) {
	registry.
		configAvroDeserializer(options).
		configJsonDeserializer(options)
}

func (registry *SchemaRegistrySetup) configAvroSerializer(options *options.KafkaOptions) *SchemaRegistrySetup {
	configSerializer := avro.NewSerializerConfig()
	configSerializer.AutoRegisterSchemas = options.SchemaRegistry.AutoRegisterSchemas
	avroSerializer, err := avro.NewSpecificSerializer(*registry.schemaregistry, serde.ValueSerde, configSerializer)

	if err != nil {
		fmt.Println("Error on create AVRO Schema Serializer")
		os.Exit(-1)
	}

	registry.avroSpecificSerializer = avroSerializer

	return registry
}

func (registry *SchemaRegistrySetup) configAvroDeserializer(options *options.KafkaOptions) *SchemaRegistrySetup {
	configDeserializer := avro.NewDeserializerConfig()
	avroDeserializer, err := avro.NewSpecificDeserializer(*registry.schemaregistry, serde.ValueSerde, configDeserializer)

	if err != nil {
		fmt.Println("Error on create AVRO Schema Deserializer")
		os.Exit(-1)
	}

	registry.avroSpecificDeserializer = avroDeserializer

	return registry
}

func (registry *SchemaRegistrySetup) configJsonSerializer(options *options.KafkaOptions) *SchemaRegistrySetup {
	configSerializer := jsonschema.NewSerializerConfig()
	configSerializer.AutoRegisterSchemas = options.SchemaRegistry.AutoRegisterSchemas
	jsonSerializer, err := jsonschema.NewSerializer(*registry.schemaregistry, serde.ValueSerde, configSerializer)

	if err != nil {
		fmt.Println("Error on create JSON Schema Serializer")
		os.Exit(-1)
	}

	registry.jsonSerializer = jsonSerializer

	return registry
}

func (registry *SchemaRegistrySetup) configJsonDeserializer(options *options.KafkaOptions) *SchemaRegistrySetup {
	configDeserializer := jsonschema.NewDeserializerConfig()
	jsonDeserializer, err := jsonschema.NewDeserializer(*registry.schemaregistry, serde.ValueSerde, configDeserializer)

	if err != nil {
		fmt.Println("Error on create JSON Schema Deserializer")
		os.Exit(-1)
	}

	registry.jsonDeserializer = jsonDeserializer

	return registry
}
