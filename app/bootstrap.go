package app

import (
	"context"

	"github.com/defval/di"
	"github.com/spf13/cobra"
)

type Bootstrapper interface {
	Bootstrap(a *AppContainer) error
}

type BootstrapFunc func(a *AppContainer) error

func (f BootstrapFunc) Bootstrap(c *AppContainer) error {
	return f(c)
}

type AppContainer struct {
	*di.Container
}

func Boot(bs ...Bootstrapper) (*AppContainer, error) {
	c, err := di.New()
	if err != nil {
		return nil, err
	}

	ac := &AppContainer{c}

	for _, b := range bs {
		if err := b.Bootstrap(ac); err != nil {
			return nil, err
		}
	}

	return &AppContainer{c}, nil
}

func (c *AppContainer) RunE(runE interface{}) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		c.ProvideValue(cmd)
		c.ProvideValue(args)
		c.ProvideValue(cmd.Context(), di.As(new(context.Context)))
		return c.Invoke(runE)
	}
}
