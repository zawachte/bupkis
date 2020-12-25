package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/containerd/containerd/remotes"
	ctypes "github.com/docker/cli/cli/config/types"
)

// Common errors
var (
	ErrNotLoggedIn = errors.New("not logged in")
)

// Client provides authentication operations for remotes.
type Client interface {
	// Login logs in to a remote server identified by the hostname.
	Login(ctx context.Context, hostname, username, secret string, insecure bool) error
	// Logout logs out from a remote server identified by the hostname.
	// Logout(ctx context.Context, hostname string) error
	// Resolver returns a new authenticated resolver.
	Resolver(ctx context.Context, client *http.Client, plainHTTP bool) (remotes.Resolver, error)
	Credential(hostname string) (string, string, error)
	GetAllCredentials() (map[string]ctypes.AuthConfig, error)
}
