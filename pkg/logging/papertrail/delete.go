package papertrail

import (
	"io"

	"github.com/fastly/cli/pkg/cmd"
	"github.com/fastly/cli/pkg/compute/manifest"
	"github.com/fastly/cli/pkg/config"
	"github.com/fastly/cli/pkg/text"
	"github.com/fastly/go-fastly/v3/fastly"
)

// DeleteCommand calls the Fastly API to delete a Papertrail logging endpoint.
type DeleteCommand struct {
	cmd.Base
	manifest       manifest.Data
	Input          fastly.DeletePapertrailInput
	serviceVersion cmd.OptionalServiceVersion
	autoClone      cmd.OptionalAutoClone
}

// NewDeleteCommand returns a usable command registered under the parent.
func NewDeleteCommand(parent cmd.Registerer, globals *config.Data) *DeleteCommand {
	var c DeleteCommand
	c.Globals = globals
	c.manifest.File.SetOutput(c.Globals.Output)
	c.manifest.File.Read(manifest.Filename)
	c.CmdClause = parent.Command("delete", "Delete a Papertrail logging endpoint on a Fastly service version").Alias("remove")
	c.RegisterServiceVersionFlag(cmd.ServiceVersionFlagOpts{
		Dst: &c.serviceVersion.Value,
	})
	c.RegisterAutoCloneFlag(cmd.AutoCloneFlagOpts{
		Action: c.autoClone.Set,
		Dst:    &c.autoClone.Value,
	})
	c.CmdClause.Flag("name", "The name of the Papertrail logging object").Short('n').Required().StringVar(&c.Input.Name)
	c.RegisterServiceIDFlag(&c.manifest.Flag.ServiceID)
	return &c
}

// Exec invokes the application logic for the command.
func (c *DeleteCommand) Exec(in io.Reader, out io.Writer) error {
	serviceID, serviceVersion, err := cmd.ServiceDetails(cmd.ServiceDetailsOpts{
		AutoCloneFlag:      c.autoClone,
		Client:             c.Globals.Client,
		Manifest:           c.manifest,
		Out:                out,
		ServiceVersionFlag: c.serviceVersion,
		VerboseMode:        c.Globals.Flag.Verbose,
	})
	if err != nil {
		c.Globals.ErrLog.Add(err)
		return err
	}

	c.Input.ServiceID = serviceID
	c.Input.ServiceVersion = serviceVersion.Number

	if err := c.Globals.Client.DeletePapertrail(&c.Input); err != nil {
		c.Globals.ErrLog.Add(err)
		return err
	}

	text.Success(out, "Deleted Papertrail logging endpoint %s (service %s version %d)", c.Input.Name, c.Input.ServiceID, c.Input.ServiceVersion)
	return nil
}
