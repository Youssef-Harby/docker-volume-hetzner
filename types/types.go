// Package types contains shared types for the docker-volume-hetzner plugin
package types

// ResizeRequest represents a request to resize a volume
type ResizeRequest struct {
	Name    string
	Options map[string]string
}
