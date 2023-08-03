package authenticate

import "time"

type Config struct {
	// CallbackHost is external host used for redirect uri
	CallbackHost string `yaml:"callback_host" mapstructure:"callback_host" default:"http://localhost:7400/v1beta1/auth/callback"`

	// OIDCCallbackHost is external host used for redirect uri
	// Deprecated: use CallbackHost instead
	OIDCCallbackHost string `yaml:"oidc_callback_host" mapstructure:"oidc_callback_host" default:"http://localhost:7400/v1beta1/auth/oidc/callback"`

	OIDCConfig map[string]OIDCConfig `yaml:"oidc_config" mapstructure:"oidc_config"`
	Session    SessionConfig         `yaml:"session" mapstructure:"session"`
	Token      TokenConfig           `yaml:"token" mapstructure:"token"`
	MailOTP    MailOTPConfig         `yaml:"mail_otp" mapstructure:"mail_otp"`
	MailLink   MailLinkConfig        `yaml:"mail_link" mapstructure:"mail_link"`
}

type TokenConfig struct {
	// Path to rsa key file, it can contain more than one key as a json array
	// jwt will be signed by first key, but will be tried to be decoded by all matching key ids, this helps in key rotation.
	// If not provided, access token will not be generated
	RSAPath string `yaml:"rsa_path" mapstructure:"rsa_path"`

	// Issuer uniquely identifies the service that issued the token
	// a good example could be fully qualified domain name
	Issuer string `yaml:"iss" mapstructure:"iss" default:"frontier"`

	// Validity is the duration for which the token is valid
	Validity time.Duration `yaml:"validity" mapstructure:"validity" default:"1h"`
}

type SessionConfig struct {
	HashSecretKey  string `mapstructure:"hash_secret_key" yaml:"hash_secret_key" default:"hash-secret-should-be-32-chars--"`
	BlockSecretKey string `mapstructure:"block_secret_key" yaml:"block_secret_key" default:"block-secret-should-be-32-chars-"`
	Domain         string `mapstructure:"domain" yaml:"domain" default:""`
}

type OIDCConfig struct {
	ClientID     string        `yaml:"client_id" mapstructure:"client_id"`
	ClientSecret string        `yaml:"client_secret" mapstructure:"client_secret"`
	IssuerUrl    string        `yaml:"issuer_url" mapstructure:"issuer_url"`
	Validity     time.Duration `yaml:"validity" mapstructure:"validity" default:"15m"`
}

type MailOTPConfig struct {
	Subject  string        `yaml:"subject" mapstructure:"subject" default:"Frontier Login - OTP"`
	Body     string        `yaml:"body" mapstructure:"body" default:"Please copy/paste the One Time Password in login form.<h2>{{.Otp}}</h2>This code will expire in 10 minutes."`
	Validity time.Duration `yaml:"validity" mapstructure:"validity" default:"10m"`
}

type MailLinkConfig struct {
	Subject  string        `yaml:"subject" mapstructure:"subject" default:"Frontier Login - One time link"`
	Body     string        `yaml:"body" mapstructure:"body" default:"Click on the following link or copy/paste the url in browser to login.<h3><a href='{{.Link}}' target='_blank'>Login</a></h3>Address: {{.Link}} <br>This link will expire in 10 minutes."`
	Validity time.Duration `yaml:"validity" mapstructure:"validity" default:"10m"`
}