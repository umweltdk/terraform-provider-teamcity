package main

import (
  "fmt"
  "github.com/hashicorp/terraform/helper/hashcode"
  "github.com/hashicorp/terraform/helper/resource"

  "github.com/umweltdk/teamcity/types"
)

func testAccCheckParameter(name,paramName string, ps types.ParameterSpec) resource.TestCheckFunc {
  hk := hashcode.String(fmt.Sprintf("%s=%s", paramName, ps))
  checked_value := ""
  unchecked_value := ""
  validation_mode := ""
  allow_multiple := "false"
  value_separator := ""
  if ps.Type.TypeName() == "checkbox" {
    v := ps.Type.(types.CheckboxType)
    checked_value = v.Checked
    unchecked_value = v.Unchecked
  } else if ps.Type.TypeName() == "text" {
    v := ps.Type.(types.TextType)
    validation_mode = v.ValidationMode
  } else if ps.Type.TypeName() == "select" {
    v := ps.Type.(types.SelectType)
    if v.AllowMultiple {
      allow_multiple = "true"
    } else {
      allow_multiple = "false"
    }
    value_separator = v.ValueSeparator
  }
  return resource.ComposeTestCheckFunc(
    resource.TestCheckResourceAttr(name, fmt.Sprintf("parameter.%d.name", hk), paramName),
    resource.TestCheckResourceAttr(name, fmt.Sprintf("parameter.%d.label", hk), ps.Label),
    resource.TestCheckResourceAttr(name, fmt.Sprintf("parameter.%d.description", hk), ps.Description),
    resource.TestCheckResourceAttr(name, fmt.Sprintf("parameter.%d.type", hk), ps.Type.TypeName()),
    resource.TestCheckResourceAttr(name, fmt.Sprintf("parameter.%d.checked_value", hk), checked_value),
    resource.TestCheckResourceAttr(name, fmt.Sprintf("parameter.%d.unchecked_value", hk), unchecked_value),
    resource.TestCheckResourceAttr(name, fmt.Sprintf("parameter.%d.validation_mode", hk), validation_mode),
    resource.TestCheckResourceAttr(name, fmt.Sprintf("parameter.%d.allow_multiple", hk), allow_multiple),
    resource.TestCheckResourceAttr(name, fmt.Sprintf("parameter.%d.value_separator", hk), value_separator),
  )
}