package main

import (
  "github.com/hashicorp/terraform/helper/schema"
  "github.com/hashicorp/terraform/helper/hashcode"

  "github.com/umweltdk/teamcity/types"

  "log"
)

func resourceParameter() *schema.Resource {
  return &schema.Resource{
    Schema: map[string]*schema.Schema{
      "name": &schema.Schema{
        Type:     schema.TypeString,
        Required: true,
      },
      "type": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
      "label": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
      "description": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
      // Text type options
      "validation_mode": &schema.Schema{
        Type:     schema.TypeString,
        Default:  "any",
        Optional: true,
      },
      // Checkbox type options
      "checked_value": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
      "unchecked_value": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
      // Select options
      "allow_multiple": &schema.Schema{
        Type:     schema.TypeBool,
        Optional: true,
      },
      "value_separator": &schema.Schema{
        Type:     schema.TypeString,
        Optional: true,
      },
      "items": &schema.Schema{
        Type:     schema.TypeList,
        Optional: true,
        Elem: &schema.Schema{Type: schema.TypeMap},
      },
    },
  }
}

func parameterHash(v interface{}) int {
  m := v.(map[string]interface{})
  return hashcode.String(m["name"].(string))
}

func parametersToDefinition(parameters types.Parameters) *schema.Set {
  ret := schema.NewSet(parameterHash, []interface{}{})
  for name, parameter := range parameters {
    param := make(map[string]interface{})
    if parameter.Spec != nil {
      spec := *parameter.Spec
      param["label"] = spec.Label
      param["description"] = spec.Description

      typeName := spec.Type.TypeName()
      param["type"] = typeName
      if typeName == "text" {
        param["validation_mode"] = spec.Type.(types.TextType).ValidationMode
      } else if typeName == "checkbox" {
        param["checked_value"] = spec.Type.(types.CheckboxType).Checked
        param["unchecked_value"] = spec.Type.(types.CheckboxType).Unchecked
      } else if typeName == "select" {
        param["allow_multiple"] = spec.Type.(types.SelectType).AllowMultiple
        param["value_separator"] = spec.Type.(types.SelectType).ValueSeparator
      }
    }
    param["name"] = name
    log.Printf("[INFO] Parameter %q\n", param)
    ret.Add(param)
  }
  return ret
}

func parameterValues(parameters types.Parameters) map[string]interface{} {
  ret := make(map[string]interface{})
  for name, parameter := range parameters {
    ret[name] = parameter.Value
  }  
  return ret
}

func definitionToParameters(parameters schema.Set) types.Parameters {
  ret := make(types.Parameters)
  for _, v := range parameters.List() {
    param := v.(map[string]interface{})
    parameter := types.Parameter{Spec: nil}
    if param["type"].(string) != "" || param["label"].(string) != "" || param["description"].(string) != "" {
      var tp types.ParameterType
      if param["type"].(string) == "text" {
        tp = &types.TextType{
          ValidationMode: param["validation_mode"].(string),
        }
      } else if param["type"].(string) == "password" {
        tp = &types.PasswordType{}
      } else if param["type"].(string) == "checkbox" {
        tp = &types.CheckboxType{
          Checked:   param["checked_value"].(string),
          Unchecked: param["unchecked_value"].(string),
        }
      } else if param["type"].(string) == "select" {
        tp = &types.SelectType{
          AllowMultiple: param["allow_multiple"].(bool),
          ValueSeparator: param["value_separator"].(string),
        }
      } else {
        tp = &types.TextType{"any"}
      }
      parameter.Spec = &types.ParameterSpec{
        Label:       param["label"].(string),
        Description: param["description"].(string),
        Type:        tp,
      }
      log.Printf("Parameter %s => %q", param["name"].(string), parameter)
    }
    ret[param["name"].(string)] = parameter
  }
  return ret
}
