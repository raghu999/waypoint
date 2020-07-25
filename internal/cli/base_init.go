package cli

import (
	"errors"
	"fmt"

	clientpkg "github.com/hashicorp/waypoint/internal/client"
	configpkg "github.com/hashicorp/waypoint/internal/config"
	"github.com/hashicorp/waypoint/internal/serverclient"
)

// This file contains the various methods that are used to perform
// the Init call on baseCommand. They are broken down into individual
// smaller methods for readability but more importantly to power the
// "init" subcommand. This allows us to share as much logic as possible
// between Init and "init" to help ensure that "init" succeeding means that
// other commands will succeed as well.

// initConfig initializes the configuration.
func (c *baseCommand) initConfig(optional bool) (*configpkg.Config, error) {
	path, err := c.initConfigPath()
	if err != nil {
		return nil, err
	}

	if path == "" {
		if optional {
			return nil, nil
		}

		return nil, errors.New("A Waypoint configuration file is required but wasn't found.")
	}

	return c.initConfigLoad(path)
}

// initConfigPath returns the configuration path to load.
func (c *baseCommand) initConfigPath() (string, error) {
	path, err := configpkg.FindPath("", "")
	if err != nil {
		return "", fmt.Errorf("Error looking for a Waypoint configuration: %s", err)
	}

	return path, nil
}

// initConfigLoad loads the configuration at the given path.
func (c *baseCommand) initConfigLoad(path string) (*configpkg.Config, error) {
	var cfg configpkg.Config
	return &cfg, cfg.LoadPath(path)
}

// initClient initializes the client.
func (c *baseCommand) initClient() (*clientpkg.Project, error) {
	// Start building our client options
	opts := []clientpkg.Option{
		clientpkg.WithLogger(c.Log),
		clientpkg.WithClientConnect(
			serverclient.FromContext(c.contextStorage, ""),
			serverclient.FromEnv(),
		),
		clientpkg.WithProjectRef(c.refProject),
		clientpkg.WithWorkspaceRef(c.refWorkspace),
		clientpkg.WithLabels(c.flagLabels),
	}
	if !c.flagRemote {
		opts = append(opts, clientpkg.WithLocal())
	}

	// Create our client
	return clientpkg.New(c.Ctx, opts...)
}