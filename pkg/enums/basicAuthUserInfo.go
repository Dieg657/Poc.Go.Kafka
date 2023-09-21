package enums

type BasicAuthCredentialsSource string

const (
	UserInfo    BasicAuthCredentialsSource = "USER_INFO"
	SaslInherit BasicAuthCredentialsSource = "SASL_INHERIT"
)
