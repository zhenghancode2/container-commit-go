package imageutil

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// PushOptions contains options for pushing an image
type PushOptions struct {
	// Authentication options
	Username string
	Password string
	// Push options
	Insecure    bool
	Concurrency int
}

// DefaultPushOptions returns default options for pushing
func DefaultPushOptions() *PushOptions {
	return &PushOptions{
		Concurrency: 3,
		Insecure:    true,
	}
}

// PushImage pushes an image to a specified repository
func PushImage(ctx context.Context, img v1.Image, destRef string, opts *PushOptions) error {
	if opts == nil {
		opts = DefaultPushOptions()
	}

	// Parse the destination reference
	ref, err := name.ParseReference(destRef)
	if err != nil {
		return fmt.Errorf("parsing reference %q: %w", destRef, err)
	}

	// Setup authentication
	auth, err := setupAuthentication(ref, opts)
	if err != nil {
		return fmt.Errorf("setting up authentication: %w", err)
	}

	// Configure remote options
	remoteOpts := configureRemoteOptions(opts, auth)

	// Push the image
	if err := remote.Write(ref, img, remoteOpts...); err != nil {
		return fmt.Errorf("writing image to remote %q: %w", ref.Name(), err)
	}

	return nil
}

// setupAuthentication configures authentication for the given reference
func setupAuthentication(ref name.Reference, opts *PushOptions) (authn.Authenticator, error) {
	if opts.Username != "" && opts.Password != "" {
		return authn.FromConfig(authn.AuthConfig{
			Username: opts.Username,
			Password: opts.Password,
		}), nil
	}

	// Use default keychain which checks environment variables and docker config
	return authn.DefaultKeychain.Resolve(ref.Context())
}

// configureRemoteOptions sets up the remote options based on the provided push options
func configureRemoteOptions(opts *PushOptions, auth authn.Authenticator) []remote.Option {
	var options []remote.Option

	options = append(options, remote.WithAuth(auth))

	if opts.Insecure {
		options = append(options, remote.WithTransport(http.DefaultTransport))
	}

	if opts.Concurrency > 0 {
		options = append(options, remote.WithJobs(opts.Concurrency))
	}

	return options
}
