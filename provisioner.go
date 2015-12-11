package main

import ( // {{{
  "fmt"
  "log"
  "time"
  "os"

  "github.com/mitchellh/packer/common"
  "github.com/mitchellh/packer/helper/config"
  "github.com/mitchellh/packer/packer"
  "github.com/mitchellh/packer/packer/plugin"
  "github.com/mitchellh/packer/template/interpolate"
) // }}}

type Config struct { // {{{
  common.PackerConfig `mapstructure:",squash"`

  // Either:
  // 1. Duration: we wait for a duration, simple
  RawDuration time.Duration `mapstructure:"duration"`

  // 2. Until: until a script/shell is true
  //    "until": { "type": "powershell", "inline": [ "if ($ready) { exit 0 } else { exit 1 }" ] },
  //    "until": { "type": "powershell", "inline": "if ($ready) { exit 0 } else { exit 1 }" },
  Until map[string]string

  // 3. While: while a script/shell is false
  //    "while": { "type": "powershell", "inline": [ "if ($ready) { exit 0 } else { exit 1 }" ] },

  // between 2 tests, sleep for the given duration
  //    "sleep": "1m",
  Sleep time.Duration `mapstructure:"sleep"`

  // and run the test a maximum of times
  //    "tries": "60",
  Tries int

  // when we get a success, execute a script/shell (as soon as we get out of the loop)
  //    "on_success": { "type": "powershell", "inline": [ "Write-Output 'Software is installed!'" ] }
  //    "on_success": { "type": "powershell", "inline": "Write-Output 'Software is installed!'" }
  OnSuccess map[string]string `mapstructure:"on_success"`

  // After all tries failed, execute a script/shell
  //    "on_failure": { "type": "powershell", "script": "./scripts/Backup-Logs.ps1" },
  OnFailure map[string]string `mapstructure:"on_failure"`

  // Display a Message before starting
  Message string

  // The configuration template
  ctx interpolate.Context
} // }}}

type Provisioner struct { // {{{
  config Config
} // }}}

func (p *Provisioner) Prepare(raw ...interface{}) error { // {{{
  err := config.Decode(&p.config, &config.DecodeOpts {
          Interpolate: true,
          InterpolateFilter: &interpolate.RenderFilter {
            Exclude: []string {
              "execute_command",
            },
           },
         }, raw...)

  if err != nil {
    return err
  }

  var errs *packer.MultiError

  if p.config.RawDuration != "" {
    p.config.duration, err = time.ParseDuration(p.config.RawDuration)
    if err != nil {
      errs = packer.MultiErrorAppend(errs, fmt.Errorf("Failed parsing duration: %s", err))
    }
  }

  if p.config.Sleep == 0 {
    p.config.Sleep = 1 * time.Second
  }

  if errs != nil && len(errs.Errors) > 0 {
    return errs
  }

  log.Printf("Prepare: sleep: %s", p.config.sleep)
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
