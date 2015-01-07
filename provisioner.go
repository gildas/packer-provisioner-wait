package main

import ( // {{{
  "fmt"
  "time"
  "os"

  "github.com/mitchellh/packer/common"
  "github.com/mitchellh/packer/packer"
  "github.com/mitchellh/packer/packer/plugin"
) // }}}

type config struct { // {{{
  common.PackerConfig `mapstructure:",squash"`

  RawDuration string `mapstructure:"duration"`
  duration time.Duration

  // The Message to display in packer
  Message string

  // The configuration template
  tpl *packer.ConfigTemplate
} // }}}

type Provisioner struct { // {{{
  config config
} // }}}

func (p *Provisioner) Prepare(raws ...interface{}) error { // {{{
  md, err := common.DecodeConfig(&p.config, raws...)
  if err != nil {
    return err
  }

  p.config.tpl, err = packer.NewConfigTemplate()
  if err != nil {
    return err
  }
  p.config.tpl.UserVars = p.config.PackerUserVars

  // Accumulate any errors
  errs := common.CheckUnusedConfig(md)

  templates := map[string]*string{
    "message": &p.config.Message,
  }

  for n, ptr := range templates {
    var err error
    *ptr, err = p.config.tpl.Process(*ptr, nil)
    if err != nil {
      errs = packer.MultiErrorAppend(
        errs, fmt.Errorf("Error processing %s: %s", n, err))
    }
  }

  if p.config.RawDuration != "" {
    p.config.duration, err = time.ParseDuration(p.config.RawDuration)
    if err != nil {
      errs = packer.MultiErrorAppend(
        errs, fmt.Errorf("Failed parsing duration: %s", err))
    }
  }

  if errs != nil && len(errs.Errors) > 0 {
    return errs
  }

  return nil
} // }}}

func (p *Provisioner) Provision(ui packer.Ui, comm packer.Communicator) error { // {{{
  if p.config.Message != "" {
    ui.Say(p.config.Message)
  }

  time.Sleep(p.config.duration)

  return nil
} // }}}

func (p *Provisioner) Cancel() { // {{{
  os.Exit(0)
} // }}}

func main() { // {{{
  server, err := plugin.Server()
  if err != nil {
    panic(err)
  }

  server.RegisterProvisioner(new(Provisioner))
  server.Serve()
} // }}}
