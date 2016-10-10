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
  parameters = [
    {
      name = "env.MUH"
      type = "password"
    },
    {
      name = "env.TEST"
      type = "text"
      validation_mode = "not_empty"
      label = "Test framework"
      description = "Name of the test framework to use"
    }
  ]
  parameter_values {
    "env.TEST" = "Hello"
  }
  steps = [
    {
      type = "simpleRunner"
      name = "Hell0"
      properties = {
        "script.content" = "npm run install"
        "teamcity.step.mode" = "default"
        "use.custom.script" = "true"
      }
    }
  ]
  attached_vcs_roots = [
    {
      vcs_root = "Single_HttpsGithubComUmweltdkDockerNodeGit"
      checkout_rules = "+:refs/heads/master\n+:refs/heads/develop"
    }
  ]
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
            "teamcity_build_configuration.bar", "steps.0.type", "simpleRunner"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "steps.0.name", "Hell0"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "steps.0.properties.script.content", "npm run install"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameters.2280284500.name", "env.MUH"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameters.2280284500.type", "password"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameters.1250194707.name", "env.TEST"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameters.1250194707.type", "text"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameters.1250194707.validation_mode", "not_empty"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameters.1250194707.label", "Test framework"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", "parameters.1250194707.description", "Name of the test framework to use"),
        ),
      },
    },
  })
}

var testAccBuildConfigProjectParameter = `
resource "teamcity_project" "parent" {
  name = "ConfProject"
  parameters = [
    {
      name = "env.CLOVER"
      type = "text"
      validation_mode = "any"
    },
    {
      name = "env.GROVER"
      type = "text"
      validation_mode = "any"
    }
  ]
  parameter_values = {
    "env.OVER" = "Parent"
  }
}
resource "teamcity_build_configuration" "bar" {
  project = "${teamcity_project.parent.id}"
  name = "Bar"
  parameters = [
    {
      name = "env.OVER"
      type = "checkbox"
      checked_value = "Hello"
    },
    {
      name = "env.PLOVER"
      type = "checkbox"
      checked_value = "Hello"
    }
  ]
  parameter_values {
    "env.OVER" = "Owner"
  }
}`

var testAccBuildConfigProjectParameterUpdate = `
resource "teamcity_project" "parent" {
  name = "ConfProject"
  parameters = [
    {
      name = "env.CLOVER"
      type = "text"
      validation_mode = "any"
    },
    {
      name = "env.PLOVER"
      type = "text"
      validation_mode = "any"
    }
  ]
  parameter_values = {
    "env.OVER" = "Parent"
    "env.PLOVER" = "Parent"
  }
}
resource "teamcity_build_configuration" "bar" {
  project = "${teamcity_project.parent.id}"
  name = "Bar"
  parameters = [
    {
      name = "env.OVER"
      type = "checkbox"
      checked_value = "Hello"
    },
    {
      name = "env.MOVER"
      type = "checkbox"
      unchecked_value = "Hello"
    }
  ]
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
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "name"), "env.OVER"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "type"), "checkbox"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "checked_value"), "Hello"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "unchecked_value"), ""),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.PLOVER", "name"), "env.PLOVER"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.PLOVER", "type"), "checkbox"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.PLOVER", "checked_value"), "Hello"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.PLOVER", "unchecked_value"), ""),
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
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "name"), "env.OVER"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "type"), "checkbox"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "checked_value"), "Hello"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "unchecked_value"), ""),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.MOVER", "name"), "env.MOVER"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.MOVER", "type"), "checkbox"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.MOVER", "checked_value"), ""),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.MOVER", "unchecked_value"), "Hello"),
        ),
      },
    },
  })
}

var testAccBuildConfigProjectTemplateParameter = `
resource "teamcity_project" "parent" {
  name = "ConfProjectTemplate"
  parameters = [
    {
      name = "env.CLOVER"
      type = "text"
      validation_mode = "any"
    },
    {
      name = "env.GROVER"
      type = "text"
      validation_mode = "any"
    }
  ]
  parameter_values = {
    "env.OVER" = "Parent"
  }
}
resource "teamcity_build_template" "far" {
  project = "${teamcity_project.parent.id}"
  name = "Far"
  parameters = [
    {
      name = "env.TCLOVER"
      type = "text"
      validation_mode = "any"
    },
    {
      name = "env.TGROVER"
      type = "text"
      validation_mode = "any"
    }
  ]
  parameter_values = {
    "env.TOVER" = "Template"
  }
}
resource "teamcity_build_configuration" "bar" {
  project = "${teamcity_project.parent.id}"
  name = "Bar"
  template = "${teamcity_build_template.far.id}"
  parameters = [
    {
      name = "env.OVER"
      type = "checkbox"
      checked_value = "Hello"
    },
    {
      name = "env.TOVER"
      type = "checkbox"
      checked_value = "Hello"
    },
    {
      name = "env.PLOVER"
      type = "checkbox"
      checked_value = "Hello"
    },
    {
      name = "env.TPLOVER"
      type = "checkbox"
      checked_value = "Hello"
    }
  ]
  parameter_values {
    "env.OVER" = "Owner"
    "env.TOVER" = "Owner"
  }
}`

var testAccBuildConfigProjectTemplateParameterUpdate = `
resource "teamcity_project" "parent" {
  name = "ConfProjectTemplate"
  parameters = [
    {
      name = "env.CLOVER"
      type = "text"
      validation_mode = "any"
    },
    {
      name = "env.PLOVER"
      type = "text"
      validation_mode = "any"
    }
  ]
  parameter_values = {
    "env.OVER" = "Parent"
    "env.PLOVER" = "Parent"
  }
}
resource "teamcity_build_template" "far" {
  project = "${teamcity_project.parent.id}"
  name = "Far"
  parameters = [
    {
      name = "env.TCLOVER"
      type = "text"
      validation_mode = "any"
    },
    {
      name = "env.TPLOVER"
      type = "text"
      validation_mode = "any"
    }
  ]
  parameter_values = {
    "env.TOVER" = "Template"
    "env.TPLOVER" = "Template"
  }
}
resource "teamcity_build_configuration" "bar" {
  project = "${teamcity_project.parent.id}"
  name = "Bar"
  template = "${teamcity_build_template.far.id}"
  parameters = [
    {
      name = "env.OVER"
      type = "checkbox"
      checked_value = "Hello"
    },
    {
      name = "env.TOVER"
      type = "checkbox"
      checked_value = "Hello"
    },
    {
      name = "env.MOVER"
      type = "checkbox"
      unchecked_value = "Hello"
    }
  ]
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
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "name"), "env.OVER"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "type"), "checkbox"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "checked_value"), "Hello"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "unchecked_value"), ""),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.PLOVER", "name"), "env.PLOVER"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.PLOVER", "type"), "checkbox"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.PLOVER", "checked_value"), "Hello"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.PLOVER", "unchecked_value"), ""),
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
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "name"), "env.OVER"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "type"), "checkbox"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "checked_value"), "Hello"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.OVER", "unchecked_value"), ""),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.MOVER", "name"), "env.MOVER"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.MOVER", "type"), "checkbox"),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.MOVER", "checked_value"), ""),
          resource.TestCheckResourceAttr(
            "teamcity_build_configuration.bar", param("env.MOVER", "unchecked_value"), "Hello"),
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

