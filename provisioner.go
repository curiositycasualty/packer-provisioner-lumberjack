package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

// template for interpolating commands ala chef provisioner(s)
type ExecuteTemplate struct {
	Sudo bool
}

type Config struct {

	// print the paths of would-be-truncated files
	PrintOnly bool `mapstructure:"print_only"`

	// to sudo or not to sudo, that is the question
	PreventSudo bool `mapstructure:"prevent_sudo"`

	// // paths to exclude from truncation
	ExcludePaths []string `mapstructure:"exclude_paths"`

	// base `find` command
	BaseCommand string `mapstructure:"base_command"`

	ctx interpolate.Context
}

type Provisioner struct {
	config Config
}

// as required for a Provisioner Packer plugin
func (p *Provisioner) Prepare(raws ...interface{}) error {

	//
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter:  &interpolate.RenderFilter{},
	}, raws...)
	if err != nil {
		return err
	}

	return nil
}

// as required for a Provisioner Packer plugin
func (p *Provisioner) Provision(ctx context.Context, ui packer.Ui, comm packer.Communicator) error {
	ui.Say("Running truncate provisioner...")

	var find_command strings.Builder

	// if not specified in config, supply the default
	if p.config.BaseCommand == "" {
		find_command.WriteString("{{if .Sudo}}sudo {{end}} find / -name \"*.log\"")
	} else {
		find_command.WriteString(string(p.config.BaseCommand))
	}

	// append ' -a -not -path ' to the command once per excluded path
	for _, path := range p.config.ExcludePaths {
		if _, err := find_command.WriteString(fmt.Sprintf(" -a -not -path %s ", path)); err != nil {
			return fmt.Errorf("Error preparing exclusion paths in shell command: %s", err)
		}
	}

	// print message or append xargs/truncate command
	if p.config.PrintOnly {
		ui.Message("Printing paths of (rather than truncating) effected files...")
	} else {
		find_command.WriteString(" | xargs {{if .Sudo}}sudo {{end}} truncate -s 0")
	}

	p.config.ctx.Data = &ExecuteTemplate{
		Sudo: !p.config.PreventSudo,
	}

	// render the command using interpolation helpers
	rendered_command, err := interpolate.Render(find_command.String(), &p.config.ctx)
	if err != nil {
		return fmt.Errorf("Error preparing interpolated shell command: %s", err)
	}

	// print the interpolated command
	ui.Say(fmt.Sprintf("Running: %s...", rendered_command))

	// use the communicator to exec the rendered command
	cmd := &packer.RemoteCmd{Command: rendered_command}
	if err := cmd.RunWithUi(ctx, comm, ui); err != nil {
		return err
	}

	// return if the rendered command tails
	if cmd.ExitStatus() != 0 {
		return fmt.Errorf("Non-zero exit status: %d", cmd.ExitStatus())
	}

	return nil
}
