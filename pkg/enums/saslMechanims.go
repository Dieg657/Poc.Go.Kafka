package enums

type SaslMechanims string

const (
	GssApi      SaslMechanims = "GSSAPI"
	Plain       SaslMechanims = "PLAIN"
	ScramSHA256 SaslMechanims = "SCRAM-SHA-256"
	ScramSHA512 SaslMechanims = "SCRAM-SHA-2512"
	OAuthBearer SaslMechanims = "OAUTHBEARER"
)
