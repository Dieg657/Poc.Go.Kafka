package enums

type MessageSerialization int

const (
	JsonStandardSerialization MessageSerialization = 1
	JsonApiSerialization      MessageSerialization = 2
	AvroSerialization         MessageSerialization = 3
)

type MessageDeserialization int

const (
	JsonStandardDeserialization MessageDeserialization = 1
	JsonApiDeserialization      MessageDeserialization = 2
	AvroDeserialization         MessageDeserialization = 3
)
