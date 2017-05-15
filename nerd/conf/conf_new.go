package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

type Conf struct {
	location string
	m        *sync.Mutex
}

type ConfigSnapshot struct {
	Auth            AuthConfig           `json:"auth"`
	Credentials     CredentialsConfig    `json:"credentials"`
	EnableLogging   bool                 `json:"enable_logging"`
	CurrentProject  CurrentProjectConfig `json:"current_project"`
	NerdAPIEndpoint string               `json:"nerd_api_endpoint"`
}

func NewConf(loc string) *Conf {
	return &Conf{
		location: loc,
		m:        &sync.Mutex{},
	}
}

//GetDefaultLocation sets the location to ~/.nerd/config.json
func GetDefaultLocation() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get home dir")
	}
	return filepath.Join(dir, ".nerd", "config.json"), nil
}

type ConfInterface interface {
	SetLocation(loc string)
	Defaults() ConfigSnapshot
	Read() (*ConfigSnapshot, error)
	WriteJWT(jwt, refreshToken string) error
	WriteOAuth(accessToken, refreshToken string, expiration time.Time, scope, tokenType string) error
}

func (c *Conf) readFile(base *ConfigSnapshot) error {
	content, err := ioutil.ReadFile(c.location)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "failed to open config file")
	}
	err = json.Unmarshal(content, base)
	if err != nil {
		return errors.Wrap(err, "failed to parse config file")
	}
	return nil
}

func (c *Conf) write(cs *ConfigSnapshot) error {
	f, err := os.Create(c.location)
	if err != nil {
		return errors.Wrap(err, "failed to create/open config file")
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "\t")
	err = enc.Encode(cs)
	if err != nil {
		return errors.Wrap(err, "failed to encode json")
	}
	return nil
}

func (c *Conf) Defaults() ConfigSnapshot {
	return ConfigSnapshot{
		Auth: AuthConfig{
			APIEndpoint:      "http://auth.nerdalize.com",
			OAuthLocalServer: "localhost:9876",
			OAuthSuccessURL:  "https://cloud.nerdalize.com",
			ClientID:         "GuoeRJLYOXzVa9ydPjKi83lCctWtXpNHuiy46Yux",
			PublicKey: `-----BEGIN PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEAkYbLnam4wo+heLlTZEeh1ZWsfruz9nk
kyvc4LwKZ8pez5KYY76H1ox+AfUlWOEq+bExypcFfEIrJkf/JXa7jpzkOWBDF9Sa
OWbQHMK+vvUXieCJvCc9Vj084ABwLBgX
-----END PUBLIC KEY-----`,
		},
		EnableLogging: false,
		CurrentProject: CurrentProjectConfig{
			Name:      "projectx",
			AWSRegion: "eu-west-1",
		},
		NerdAPIEndpoint: "https://batch.nerdalize.com/v1",
	}
}

func (c *Conf) SetLocation(loc string) {
	c.m.Lock()
	defer c.m.Unlock()
	c.location = loc
}

func (c *Conf) Read() (*ConfigSnapshot, error) {
	c.m.Lock()
	defer c.m.Unlock()
	conf := c.Defaults()
	err := c.readFile(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func (c *Conf) WriteJWT(jwt, refreshToken string) error {
	c.m.Lock()
	defer c.m.Unlock()
	conf := &ConfigSnapshot{}
	err := c.readFile(conf)
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}
	conf.Credentials.JWT.Token = jwt
	conf.Credentials.JWT.RefreshToken = refreshToken
	err = c.write(conf)
	if err != nil {
		return errors.Wrap(err, "failed to write config")
	}
	return nil
}

func (c *Conf) WriteOAuth(accessToken, refreshToken string, expiration time.Time, scope, tokenType string) error {
	c.m.Lock()
	defer c.m.Unlock()
	conf := &ConfigSnapshot{}
	err := c.readFile(conf)
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}
	conf.Credentials.OAuth.AccessToken = accessToken
	conf.Credentials.OAuth.RefreshToken = refreshToken
	conf.Credentials.OAuth.Expiration = expiration
	conf.Credentials.OAuth.Scope = scope
	conf.Credentials.OAuth.TokenType = tokenType
	err = c.write(conf)
	if err != nil {
		return errors.Wrap(err, "failed to write config")
	}
	return nil
}
