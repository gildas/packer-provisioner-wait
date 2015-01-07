package wait

import (
  "testing"
  "github.com/mitchellh/packer/packer"
)

func testConfig() map[string]interface{} {
  return map[string]interface{} {
    "duration": []interface{} { "20s" },
    "message":  []interface{} { "Waiting for 20 seconds" },
  }
}

func TestImplementsProvisioner(t *testing.T) {
  var raw interface{}

  raw = &Provisioner{}

  if _, ok := raw.(packer.Provisioner); !ok {
    t.Fatalf("Interface packer.Provisioner is not implemented")
  }
}

func TestPrepare_InvalidKey(t *testing.T) {
  var p Provisioner

  config := testConfig()
  config["i_should_not_be_valid"] = true

  err := p.Prepare(config)
  if err == nil {
    t.Fatalf("should have error")
  }
}

func TestPrepare_Defaults(t *testing.T) {
  var p Provisioner

  config := map[string]interface{}{ "duration": "40s" }

  if err := p.Prepare(config); err != nil {
    t.Fatalf("Error: %s", err)
  }
}

