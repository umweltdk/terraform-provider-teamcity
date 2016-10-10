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
  "github.com/hashicorp/terraform/helper/hashcode"
  "github.com/hashicorp/terraform/helper/resource"
  "github.com/hashicorp/terraform/terraform"

  "github.com/umweltdk/teamcity/teamcity"
  "github.com/umweltdk/teamcity/types"
)

var testAccProject = `
resource "teamcity_project" "bar" {
  parent = "Single"
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
    "env.MUH" = "Hush Hush"
    "env.TEST" = "Hello"
  }
}`

func TestAccProject_basic(t *testing.T) {
  var v types.Project

  resource.Test(t, resource.TestCase{
    PreCheck:     func() { testAccPreCheck(t) },
    Providers:    testAccProviders,
    CheckDestroy: testAccCheckProjectDestroy,
    Steps: []resource.TestStep{
      resource.TestStep{
        Config: testAccProject,
        Check: resource.ComposeTestCheckFunc(
          testAccCheckProjectExists("teamcity_project.bar", &v),
          testAccCheckProjectAttributes(&v),
          testAccCheckProjectParameters(&v, &types.Parameters{
            "env.TEST": types.Parameter{
              Value: "Hello",
              Spec: &types.ParameterSpec{
                Label: "Test framework",
                Description: "Name of the test framework to use",
                Type: types.TextType{"not_empty"},
              },
            },
            "env.MUH": types.Parameter{
              Spec: &types.ParameterSpec{
                Type: types.PasswordType{},
              },
            },
          }),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", "parent", "Single"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", "name", "bar"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", "description", ""),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", "parameter_values.env.TEST", "Hello"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.MUH", "name"), "env.MUH"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.MUH", "type"), "password"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.TEST", "name"), "env.TEST"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.TEST", "type"), "text"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.TEST", "validation_mode"), "not_empty"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.TEST", "label"), "Test framework"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.TEST", "description"), "Name of the test framework to use"),
        ),
      },
    },
  })
}

var testAccProjectParentParameter = `
resource "teamcity_project" "parent" {
  name = "Parent"
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
resource "teamcity_project" "bar" {
  parent = "${teamcity_project.parent.id}"
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

var testAccProjectParentParameterUpdate = `
resource "teamcity_project" "parent" {
  name = "Parent"
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
resource "teamcity_project" "bar" {
  parent = "${teamcity_project.parent.id}"
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

func TestAccProject_parentParameters(t *testing.T) {
  var v types.Project

  resource.Test(t, resource.TestCase{
    PreCheck:     func() { testAccPreCheck(t) },
    Providers:    testAccProviders,
    CheckDestroy: testAccCheckProjectDestroy,
    Steps: []resource.TestStep{
      resource.TestStep{
        Config: testAccProjectParentParameter,
        Check: resource.ComposeTestCheckFunc(
          testAccCheckProjectExists("teamcity_project.bar", &v),
          testAccCheckProjectAttributes(&v),
          testAccCheckProjectParameters(&v, &types.Parameters{
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
            "teamcity_project.bar", "parent", "Parent"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", "name", "Bar"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", "description", ""),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", "parameter_values.env.OVER", "Owner"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.OVER", "name"), "env.OVER"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.OVER", "type"), "checkbox"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.OVER", "checked_value"), "Hello"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.OVER", "unchecked_value"), ""),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.PLOVER", "name"), "env.PLOVER"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.PLOVER", "type"), "checkbox"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.PLOVER", "checked_value"), "Hello"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.PLOVER", "unchecked_value"), ""),
        ),
      },
      resource.TestStep{
        Config: testAccProjectParentParameterUpdate,
        Check: resource.ComposeTestCheckFunc(
          testAccCheckProjectExists("teamcity_project.bar", &v),
          testAccCheckProjectAttributes(&v),
          testAccCheckProjectParameters(&v, &types.Parameters{
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
            "teamcity_project.bar", "parent", "Parent"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", "name", "Bar"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", "description", ""),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", "parameter_values.env.OVER", "Owner"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", "parameter_values.env.PLOVER", "Owner"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.OVER", "name"), "env.OVER"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.OVER", "type"), "checkbox"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.OVER", "checked_value"), "Hello"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.OVER", "unchecked_value"), ""),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.MOVER", "name"), "env.MOVER"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.MOVER", "type"), "checkbox"),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.MOVER", "checked_value"), ""),
          resource.TestCheckResourceAttr(
            "teamcity_project.bar", param("env.MOVER", "unchecked_value"), "Hello"),
        ),
      },
    },
  })
}

func testAccCheckProjectDestroy(s *terraform.State) error {
  client := testAccProvider.Meta().(*teamcity.Client)

  for _, rs := range s.RootModule().Resources {
    if rs.Type != "teamcity_project" {
      continue
    }

    // Try to find the Group
    var err error
    project, err := client.GetProject(rs.Primary.ID)

    if err == nil && project == nil {
      continue
    }

    if err == nil {
      return fmt.Errorf("Project still exists")
    }

    return err
  }

  return nil
}

func testAccCheckProjectAttributes(v *types.Project) resource.TestCheckFunc {
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

func testAccCheckProjectParameters(v *types.Project, e *types.Parameters) resource.TestCheckFunc {
  return func(s *terraform.State) error {
    expected := *e
    if !reflect.DeepEqual(v.Parameters, expected) {
      return fmt.Errorf("bad parameters: %q %q", v.Parameters, expected)
    }

    return nil
  }
}

func testAccCheckProjectExists(n string, v *types.Project) resource.TestCheckFunc {
  return func(s *terraform.State) error {
    rs, ok := s.RootModule().Resources[n]
    if !ok {
      return fmt.Errorf("Not found: %s", n)
    }

    if rs.Primary.ID == "" {
      return fmt.Errorf("No Project ID is set")
    }

    client := testAccProvider.Meta().(*teamcity.Client)

    project, err := client.GetProject(rs.Primary.ID)
    if err != nil {
      return err
    }

    if project == nil {
      return fmt.Errorf("Project not found")
    }

    *v = *project

    return nil
  }
}

func param(name string, parameter string) string {
  return fmt.Sprintf("parameters.%d.%s", hashcode.String(name), parameter)
}
