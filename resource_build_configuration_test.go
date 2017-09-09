package main

import (
  "fmt"
  //"log"
  //"regexp"
  //"strings"

  //"math/rand"
  "reflect"
  "testing"
  //"time"

  //"github.com/hashicorp/terraform/helper/acctest"
  "github.com/hashicorp/terraform/helper/resource"
  "github.com/hashicorp/terraform/terraform"

  "github.com/umweltdk/teamcity/teamcity"
  "github.com/umweltdk/teamcity/types"
)

var testAccBuildConfig = `
resource "teamcity_build_configuration" "bar" {
  project = "Single"
  name = "bar"
  parameter {
    name = "env.MUH"
    type = "password"
  }
  parameter {
    name = "env.TEST"
    type = "text"
    validation_mode = "not_empty"
    label = "Test framework"
    description = "Name of the test framework to use"
  }
  parameter_values {
    "env.TEST" = "Hello"
  }
  step {
    type = "simpleRunner"
    name = "Hell0"
    properties = {
      "script.content" = "npm run install"
      "teamcity.step.mode" = "default"
      "use.custom.script" = "true"
    }
  }
  attached_vcs_root {
    vcs_root = "Single_HttpsGithubComUmweltdkDockerNodeGit"
    checkout_rules = "+:refs/heads/master\n+:refs/heads/develop"
  }
}`

func TestAccBuildConfig_basic(t *testing.T) {
  var v types.BuildConfiguration

  resource.Test(t, resource.TestCase{
    PreCheck:     func() { testAccPreCheck(t) },
    Providers:    testAccProviders,
    CheckDestroy: testAccCheckBuildConfigDestroy,
    Steps: []resource.TestStep{
      resource.TestStep{
        Config: testAccBuildConfig,
        Check: resource.ComposeTestCheckFunc(
          testAccCheckBuildConfigExists("teamcity_build_configuration.bar", &v),
          testAccCheckBuildConfigAttributes(&v),
          testAccCheckBuildConfigSteps(&v, &types.BuildSteps{
            types.BuildStep{
              Type: "simpleRunner",
              Name: "Hell0",
              Properties: types.Properties{
                "script.content": "npm run install",
                "teamcity.step.mode": "default",
                "use.custom.script": "true",
              },
            },
          }),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "project", "Single"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "name", "bar"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "template", ""),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "description", ""),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "step.0.type", "simpleRunner"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "step.0.name", "Hell0"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "step.0.properties.script.content", "npm run install"),
          testAccCheckParameter("teamcity_build_configuration.bar", "env.MUH", types.ParameterSpec{
              Type: types.PasswordType{},
            }),
          testAccCheckParameter("teamcity_build_configuration.bar", "env.TEST", types.ParameterSpec{
              Label: "Test framework",
              Description: "Name of the test framework to use",
              Type: types.TextType{"not_empty"},
            }),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameter.#", "2"),
          testAccCheckAttachedRoot("teamcity_build_configuration.bar",
            "Single_HttpsGithubComUmweltdkDockerNodeGit",
            "+:refs/heads/master\n+:refs/heads/develop"),
        ),
      },
    },
  })
}

var testAccBuildConfigProjectParameter = `
resource "teamcity_project" "parent" {
  name = "ConfProject"
  parameter {
    name = "env.CLOVER"
    type = "text"
    validation_mode = "any"
  }
  parameter {
    name = "env.GROVER"
    type = "text"
    validation_mode = "any"
  }
  parameter_values = {
    "env.OVER" = "Parent"
  }
}
resource "teamcity_build_configuration" "bar" {
  project = "${teamcity_project.parent.id}"
  name = "Bar"
  parameter {
    name = "env.OVER"
    type = "checkbox"
    checked_value = "Hello"
  }
  parameter {
    name = "env.PLOVER"
    type = "checkbox"
    checked_value = "Hello"
  }
  parameter_values {
    "env.OVER" = "Owner"
  }
}`

var testAccBuildConfigProjectParameterUpdate = `
resource "teamcity_project" "parent" {
  name = "ConfProject"
  parameter {
    name = "env.CLOVER"
    type = "text"
    validation_mode = "any"
  }
  parameter {
    name = "env.PLOVER"
    type = "text"
    validation_mode = "any"
  }
  parameter_values = {
    "env.OVER" = "Parent"
    "env.PLOVER" = "Parent"
  }
}
resource "teamcity_build_configuration" "bar" {
  project = "${teamcity_project.parent.id}"
  name = "Bar"
  parameter {
    name = "env.OVER"
    type = "checkbox"
    checked_value = "Hello"
  }
  parameter {
    name = "env.MOVER"
    type = "checkbox"
    unchecked_value = "Hello"
  }
  parameter_values {
    "env.OVER" = "Owner"
    "env.PLOVER" = "Owner"
  }
}`

func TestAccBuildConfig_projectParameters(t *testing.T) {
  var v types.BuildConfiguration

  resource.Test(t, resource.TestCase{
    PreCheck:     func() { testAccPreCheck(t) },
    Providers:    testAccProviders,
    CheckDestroy: testAccCheckBuildConfigDestroy,
    IDRefreshName: "teamcity_build_configuration.bar",
    Steps: []resource.TestStep{
      resource.TestStep{
        Config: testAccBuildConfigProjectParameter,
        Check: resource.ComposeTestCheckFunc(
          testAccCheckBuildConfigExists("teamcity_build_configuration.bar", &v),
          testAccCheckBuildConfigAttributes(&v),
          testAccCheckBuildConfigParameters(&v, &types.Parameters{
            "env.OVER": types.Parameter{
              Value: "Owner",
              Spec: &types.ParameterSpec{
                Type: types.CheckboxType{"Hello", ""},
              },
            },
            "env.CLOVER": types.Parameter{
              Value: "",
              Spec: &types.ParameterSpec{
                Type: types.TextType{"any"},
              },
            },
            "env.GROVER": types.Parameter{
              Value: "",
              Spec: &types.ParameterSpec{
                Type: types.TextType{"any"},
              },
            },
            "env.PLOVER": types.Parameter{
              Value: "",
              Spec: &types.ParameterSpec{
                Type: types.CheckboxType{"Hello", ""},
              },
            },
          }),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "project", "ConfProject"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "name", "Bar"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "description", ""),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameter_values.env.OVER", "Owner"),
          testAccCheckParameter("teamcity_build_configuration.bar", "env.OVER", types.ParameterSpec{
              Type: types.CheckboxType{Checked: "Hello",},
            }),
          testAccCheckParameter("teamcity_build_configuration.bar", "env.PLOVER", types.ParameterSpec{
              Type: types.CheckboxType{Checked: "Hello",},
            }),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameter.#", "2"),
        ),
      },
      resource.TestStep{
        Config: testAccBuildConfigProjectParameterUpdate,
        Check: resource.ComposeTestCheckFunc(
          testAccCheckBuildConfigExists("teamcity_build_configuration.bar", &v),
          testAccCheckBuildConfigAttributes(&v),
          testAccCheckBuildConfigParameters(&v, &types.Parameters{
            "env.OVER": types.Parameter{
              Value: "Owner",
              Spec: &types.ParameterSpec{
                Type: types.CheckboxType{"Hello", ""},
              },
            },
            "env.MOVER": types.Parameter{
              Value: "",
              Spec: &types.ParameterSpec{
                Type: types.CheckboxType{"", "Hello"},
              },
            },
            "env.CLOVER": types.Parameter{
              Value: "",
              Spec: &types.ParameterSpec{
                Type: types.TextType{"any"},
              },
            },
            "env.PLOVER": types.Parameter{
              Value: "Owner",
              Spec: &types.ParameterSpec{
                Type: types.TextType{"any"},
              },
            },
          }),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "project", "ConfProject"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "name", "Bar"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "description", ""),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameter_values.env.OVER", "Owner"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameter_values.env.PLOVER", "Owner"),
          testAccCheckParameter("teamcity_build_configuration.bar", "env.OVER", types.ParameterSpec{
              Type: types.CheckboxType{Checked: "Hello",},
            }),
          testAccCheckParameter("teamcity_build_configuration.bar", "env.MOVER", types.ParameterSpec{
              Type: types.CheckboxType{Unchecked: "Hello",},
            }),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameter.#", "2"),
        ),
      },
    },
  })
}

var testAccBuildConfigProjectTemplateParameter = `
resource "teamcity_project" "parent" {
  name = "ConfProjectTemplate"
  parameter {
    name = "env.CLOVER"
    type = "text"
    validation_mode = "any"
  }
  parameter {
    name = "env.GROVER"
    type = "text"
    validation_mode = "any"
  }
  parameter_values = {
    "env.OVER" = "Parent"
  }
}
resource "teamcity_build_template" "far" {
  project = "${teamcity_project.parent.id}"
  name = "Far"
  parameter {
    name = "env.TCLOVER"
    type = "text"
    validation_mode = "any"
  }
  parameter {
    name = "env.TGROVER"
    type = "text"
    validation_mode = "any"
  }
  parameter_values = {
    "env.TOVER" = "Template"
  }
}
resource "teamcity_build_configuration" "bar" {
  project = "${teamcity_project.parent.id}"
  name = "Bar"
  template = "${teamcity_build_template.far.id}"
  parameter {
    name = "env.OVER"
    type = "checkbox"
    checked_value = "Hello"
  }
  parameter {
    name = "env.TOVER"
    type = "checkbox"
    checked_value = "Hello"
  }
  parameter {
    name = "env.PLOVER"
    type = "checkbox"
    checked_value = "Hello"
  }
  parameter {
    name = "env.TPLOVER"
    type = "checkbox"
    checked_value = "Hello"
  }
  parameter_values {
    "env.OVER" = "Owner"
    "env.TOVER" = "Owner"
  }
}`

var testAccBuildConfigProjectTemplateParameterUpdate = `
resource "teamcity_project" "parent" {
  name = "ConfProjectTemplate"
  parameter {
    name = "env.CLOVER"
    type = "text"
    validation_mode = "any"
  }
  parameter {
    name = "env.PLOVER"
    type = "text"
    validation_mode = "any"
  }
  parameter_values = {
    "env.OVER" = "Parent"
    "env.PLOVER" = "Parent"
  }
}
resource "teamcity_build_template" "far" {
  project = "${teamcity_project.parent.id}"
  name = "Far"
  parameter {
    name = "env.TCLOVER"
    type = "text"
    validation_mode = "any"
  }
  parameter {
    name = "env.TPLOVER"
    type = "text"
    validation_mode = "any"
  }
  parameter_values = {
    "env.TOVER" = "Template"
    "env.TPLOVER" = "Template"
  }
}
resource "teamcity_build_configuration" "bar" {
  project = "${teamcity_project.parent.id}"
  name = "Bar"
  template = "${teamcity_build_template.far.id}"
  parameter {
    name = "env.OVER"
    type = "checkbox"
    checked_value = "Hello"
  }
  parameter {
    name = "env.TOVER"
    type = "checkbox"
    checked_value = "Hello"
  }
  parameter {
    name = "env.MOVER"
    type = "checkbox"
    unchecked_value = "Hello"
  }
  parameter_values {
    "env.OVER" = "Owner"
    "env.TOVER" = "Owner"
    "env.PLOVER" = "Owner"
    "env.TPLOVER" = "Owner"
  }
}`

func TestAccBuildConfig_projectTemplateParameters(t *testing.T) {
  var v types.BuildConfiguration

  resource.Test(t, resource.TestCase{
    PreCheck:     func() { testAccPreCheck(t) },
    Providers:    testAccProviders,
    CheckDestroy: testAccCheckBuildConfigDestroy,
    IDRefreshName: "teamcity_build_configuration.bar",
    Steps: []resource.TestStep{
      resource.TestStep{
        Config: testAccBuildConfigProjectTemplateParameter,
        Check: resource.ComposeTestCheckFunc(
          testAccCheckBuildConfigExists("teamcity_build_configuration.bar", &v),
          testAccCheckBuildConfigAttributes(&v),
          testAccCheckBuildConfigParameters(&v, &types.Parameters{
            "env.OVER": types.Parameter{
              Value: "Owner",
              Spec: &types.ParameterSpec{
                Type: types.CheckboxType{"Hello", ""},
              },
            },
            "env.TOVER": types.Parameter{
              Value: "Owner",
              Spec: &types.ParameterSpec{
                Type: types.CheckboxType{"Hello", ""},
              },
            },
            "env.CLOVER": types.Parameter{
              Value: "",
              Spec: &types.ParameterSpec{
                Type: types.TextType{"any"},
              },
            },
            "env.TCLOVER": types.Parameter{
              Value: "",
              Spec: &types.ParameterSpec{
                Type: types.TextType{"any"},
              },
            },
            "env.GROVER": types.Parameter{
              Value: "",
              Spec: &types.ParameterSpec{
                Type: types.TextType{"any"},
              },
            },
            "env.TGROVER": types.Parameter{
              Value: "",
              Spec: &types.ParameterSpec{
                Type: types.TextType{"any"},
              },
            },
            "env.PLOVER": types.Parameter{
              Value: "",
              Spec: &types.ParameterSpec{
                Type: types.CheckboxType{"Hello", ""},
              },
            },
            "env.TPLOVER": types.Parameter{
              Value: "",
              Spec: &types.ParameterSpec{
                Type: types.CheckboxType{"Hello", ""},
              },
            },
          }),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "project", "ConfProjectTemplate"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "name", "Bar"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "description", ""),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameter_values.env.OVER", "Owner"),
          testAccCheckParameter("teamcity_build_configuration.bar", "env.OVER", types.ParameterSpec{
              Type: types.CheckboxType{Checked: "Hello",},
            }),
          testAccCheckParameter("teamcity_build_configuration.bar", "env.TOVER", types.ParameterSpec{
              Type: types.CheckboxType{Checked: "Hello",},
            }),
          testAccCheckParameter("teamcity_build_configuration.bar", "env.PLOVER", types.ParameterSpec{
              Type: types.CheckboxType{Checked: "Hello",},
            }),
          testAccCheckParameter("teamcity_build_configuration.bar", "env.TPLOVER", types.ParameterSpec{
              Type: types.CheckboxType{Checked: "Hello",},
            }),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameter.#", "4"),
        ),
      },
      resource.TestStep{
        Config: testAccBuildConfigProjectTemplateParameterUpdate,
        Check: resource.ComposeTestCheckFunc(
          testAccCheckBuildConfigExists("teamcity_build_configuration.bar", &v),
          testAccCheckBuildConfigAttributes(&v),
          testAccCheckBuildConfigParameters(&v, &types.Parameters{
            "env.OVER": types.Parameter{
              Value: "Owner",
              Spec: &types.ParameterSpec{
                Type: types.CheckboxType{"Hello", ""},
              },
            },
            "env.TOVER": types.Parameter{
              Value: "Owner",
              Spec: &types.ParameterSpec{
                Type: types.CheckboxType{"Hello", ""},
              },
            },
            "env.MOVER": types.Parameter{
              Value: "",
              Spec: &types.ParameterSpec{
                Type: types.CheckboxType{"", "Hello"},
              },
            },
            "env.CLOVER": types.Parameter{
              Value: "",
              Spec: &types.ParameterSpec{
                Type: types.TextType{"any"},
              },
            },
            "env.TCLOVER": types.Parameter{
              Value: "",
              Spec: &types.ParameterSpec{
                Type: types.TextType{"any"},
              },
            },
            "env.PLOVER": types.Parameter{
              Value: "Owner",
              Spec: &types.ParameterSpec{
                Type: types.TextType{"any"},
              },
            },
            "env.TPLOVER": types.Parameter{
              Value: "Owner",
              Spec: &types.ParameterSpec{
                Type: types.TextType{"any"},
              },
            },
          }),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "project", "ConfProjectTemplate"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "name", "Bar"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "description", ""),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameter_values.env.OVER", "Owner"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameter_values.env.PLOVER", "Owner"),
          testAccCheckParameter("teamcity_build_configuration.bar", "env.OVER", types.ParameterSpec{
              Type: types.CheckboxType{Checked: "Hello",},
            }),
          testAccCheckParameter("teamcity_build_configuration.bar", "env.TOVER", types.ParameterSpec{
              Type: types.CheckboxType{Checked: "Hello",},
            }),
          testAccCheckParameter("teamcity_build_configuration.bar", "env.MOVER", types.ParameterSpec{
              Type: types.CheckboxType{Unchecked: "Hello",},
            }),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameter.#", "3"),
        ),
      },
    },
  })
}

func testAccCheckBuildConfigDestroy(s *terraform.State) error {
  client := testAccProvider.Meta().(*teamcity.Client)

  for _, rs := range s.RootModule().Resources {
    if rs.Type != "teamcity_build_configuration" {
      continue
    }

    // Try to find the Group
    var err error
    config, err := client.GetBuildConfiguration(rs.Primary.ID)

    if err == nil && config == nil {
      continue
    }

    if err == nil {
      return fmt.Errorf("Build configuration still exists")
    }

    return err
  }

  return nil
}

func testAccCheckBuildConfigAttributes(v *types.BuildConfiguration) resource.TestCheckFunc {
  return func(s *terraform.State) error {

    /*
    if *v.Engine != "mysql" {
      return fmt.Errorf("bad engine: %#v", *v.Engine)
    }

    if *v.EngineVersion == "" {
      return fmt.Errorf("bad engine_version: %#v", *v.EngineVersion)
    }

    if *v.BackupRetentionPeriod != 0 {
      return fmt.Errorf("bad backup_retention_period: %#v", *v.BackupRetentionPeriod)
    }
    */

    return nil
  }
}

func testAccCheckBuildConfigParameters(v *types.BuildConfiguration, e *types.Parameters) resource.TestCheckFunc {
  return func(s *terraform.State) error {
    expected := *e
    if !reflect.DeepEqual(v.Parameters, expected) {
      return fmt.Errorf("bad parameters: %q %q", v.Parameters, expected)
    }

    return nil
  }
}

func testAccCheckBuildConfigSteps(v *types.BuildConfiguration, e *types.BuildSteps) resource.TestCheckFunc {
  return func(s *terraform.State) error {
    expected := *e
    if len(v.Steps) == len(expected) {
      for idx, step := range v.Steps {
        if expected[idx].ID == "" {
          expected[idx].ID = step.ID
        }
      }
    }
    if !reflect.DeepEqual(v.Steps, expected) {
      return fmt.Errorf("bad steps: %q %q", v.Steps, expected)
    }

    return nil
  }
}

func testAccCheckBuildConfigExists(n string, v *types.BuildConfiguration) resource.TestCheckFunc {
  return func(s *terraform.State) error {
    rs, ok := s.RootModule().Resources[n]
    if !ok {
      return fmt.Errorf("Not found: %s", n)
    }

    if rs.Primary.ID == "" {
      return fmt.Errorf("No DB Instance ID is set")
    }

    client := testAccProvider.Meta().(*teamcity.Client)

    config, err := client.GetBuildConfiguration(rs.Primary.ID)
    if err != nil {
      return err
    }

    if config == nil {
      return fmt.Errorf("Build configuration not found")
    }

    *v = *config

    return nil
  }
}

func testAccCheckAttachedRoot(name, vcs_root, checkout_rules string) resource.TestCheckFunc {
  kks := make(map[string]interface{}) 
  kks["vcs_root"] = vcs_root
  kks["checkout_rules"] = checkout_rules
  hk := attachedVcsRootValueHash(kks)
  return resource.ComposeTestCheckFunc(
    resource.TestCheckResourceAttr(name, fmt.Sprintf("attached_vcs_root.%d.vcs_root", hk), vcs_root),
    resource.TestCheckResourceAttr(name, fmt.Sprintf("attached_vcs_root.%d.checkout_rules", hk), checkout_rules),
  )
}