package client

import (
	"bytes"
	"encoding/json"
	"path"

	"github.com/nerdalize/nerd/nerd/client"
)

const (
	//MetadataObjectKey is the key of the object that contains an a dataset's metadata.
	MetadataObjectKey = "metadata"
)

//MetadataExists checks if the metadata object exists.
func (c *Client) MetadataExists(bucket, root string) (bool, error) {
	return c.Exists(bucket, path.Join(root, MetadataObjectKey))
}

//MetadataUpload uploads a dataset's metadata.
func (c *Client) MetadataUpload(bucket, root string, m *Metadata) error {
	dat, err := json.Marshal(m)
	if err != nil {
		return client.NewError("failed to encode metadata", err)
	}
	err = c.Upload(bucket, path.Join(root, MetadataObjectKey), bytes.NewReader(dat))
	if err != nil {
		return client.NewError("failed to upload index file", err)
	}
	return nil
}

//MetadataDownload downloads a dataset's metadata.
func (c *Client) MetadataDownload(bucket, root string) (*Metadata, error) {
	r, err := c.Download(bucket, path.Join(root, MetadataObjectKey))
	if err != nil {
		return nil, client.NewError("failed to download metadata", err)
	}
	dec := json.NewDecoder(r)
	metadata := &Metadata{}
	err = dec.Decode(metadata)
	if err != nil {
		return nil, client.NewError("failed to decode metadata", err)
	}
	return metadata, nil
}
