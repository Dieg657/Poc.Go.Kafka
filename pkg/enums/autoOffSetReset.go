package enums

type AutoOffsetReset string

const (
	Earliest AutoOffsetReset = "earliest"
	Latest   AutoOffsetReset = "latest"
	Error    AutoOffsetReset = "error"
)
