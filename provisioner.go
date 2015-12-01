package main

import ( // {{{
  "fmt"
  "log"
  "time"
  "os"

  "github.com/mitchellh/packer/common"
  "github.com/mitchellh/packer/packer"
  "github.com/mitchellh/packer/packer/plugin"
) // }}}

type config struct { // {{{
  common.PackerConfig `mapstructure:",squash"`

  // Either:
  // 1. Duration: we wait for a duration, simple
  RawDuration string `mapstructure:"duration"`
  duration time.Duration

  // 2. Until: until a script/shell is true
  //    "until": { "type": "powershell", "inline": [ "if ($ready) { exit 0 } else { exit 1 }" ] },

  // 3. While: while a script/shell is false
  //    "while": { "type": "powershell", "inline": [ "if ($ready) { exit 0 } else { exit 1 }" ] },

  // between 2 tests, sleep for the given duration
  //    "sleep": "1m",
  Sleep string `mapstructure:"sleep"`
  sleep time.Duration

  // and run the test a maximum of times
  //    "tries": "60",
  Tries int

  // when we get a success, execute a script/shell (as soon as we get out of the loop)
  //    "on_success": { "type": "powershell", "inline": [ "Write-Output 'Software is installed!'" ] }

  // After all tries failed, execute a script/shell
  //    "on_failure": { "type": "powershell", "script": "./scripts/Backup-Logs.ps1" },

  // Display a Message before starting
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
