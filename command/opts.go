package command

import (
	"fmt"
	"net/url"
	"strings"
)

//NerdAPIOpts configure how the platform endpoint is reached
type NerdAPIOpts struct {
	NerdAPIScheme string `long:"api-scheme" default:"https" default-mask:"https" env:"NERD_API_SCHEME" description:"protocol for endpoint communication"`

	NerdAPIHostname string `long:"api-hostname" default:"platform.nerdalize.net" default-mask:"platform.nerdalize.net" env:"NERD_API_HOST" description:"hostname of the compute platform"`

	NerdAPIVersion string `long:"api-basepath" default:"v1" default-mask:"v1" env:"NERD_API_VERSION" description:"endpoint version as basepath"`
}

//URL returns a fully qualitied url on the platform endpoint
func (opts *NerdAPIOpts) URL(path string) (loc *url.URL, err error) {
	loc, err = url.Parse(fmt.Sprintf(
		"%s://%s/%s/%s",
		opts.NerdAPIScheme,
		opts.NerdAPIHostname,
		opts.NerdAPIVersion,
		strings.TrimLeft(path, "/"),
	))
	return loc, err
}