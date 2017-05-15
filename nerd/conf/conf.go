//Package conf gives the CLI access to the nerd config file. By default this config file is
//~/.nerd/config.json, but the location can be changed using SetLocation().
//
//All read and write operation to the config file should go through the Read() and Write() functions.
//This way we can keep an in-memory representation of the config (in the global conf variable) for fast read.
package conf

import (
	"time"
)

//AuthConfig contains config details with respect to the authentication server.
type AuthConfig struct {
	APIEndpoint      string `json:"api_endpoint"`
	PublicKey        string `json:"public_key"`
	ClientID         string `json:"client_id"`
	OAuthSuccessURL  string `json:"oauth_success_url"`
	OAuthLocalServer string `json:"nerd_oauth_localserver"`
}

//CredentialsConfig contains oauth and jwt credentials
type CredentialsConfig struct {
	OAuth OAuthConfig `json:"oauth,omitempty"`
	JWT   JWTConfig   `json:"jwt,omitempty"`
}

//OAuthConfig contians oauth credentials
type OAuthConfig struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiration   time.Time `json:"expiration"`
	Scope        string    `json:"scope"`
	TokenType    string    `json:"token_type"`
}

//JWTConfig contains JWT credentials
type JWTConfig struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

//CurrentProjectConfig contains details of the current working project.
type CurrentProjectConfig struct {
	Name      string `json:"current_project"`
	AWSRegion string `json:"aws_region"`
}
