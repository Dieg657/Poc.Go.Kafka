package enums

type SecurityProtocol string

const (
	Plaintext    SecurityProtocol = "PLAIN"
	SslPlaintext SecurityProtocol = "SSL_PLAINTEXT"
	Ssl          SecurityProtocol = "SSL"
	SaslSsl      SecurityProtocol = "SASL_SSL"
)
