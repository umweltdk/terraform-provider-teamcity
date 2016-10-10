package main

import (
    "fmt"

    "github.com/hashicorp/terraform/helper/schema"

    "github.com/umweltdk/teamcity/teamcity"
    "github.com/umweltdk/teamcity/types"

    "log"
    "reflect"
)

func resourceProject() *schema.Resource {
    return &schema.Resource{
        Create: resourceProjectCreate,
        Read:   resourceProjectRead,
        Update: resourceProjectUpdate,
        Delete: resourceProjectDelete,

        Schema: map[string]*schema.Schema{
            "parent": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
                ForceNew: true,
                ValidateFunc: teamcity.ValidateID,
            },
            "name": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "description": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
            },
            "parameters": &schema.Schema{
                Type:     schema.TypeSet,
                Elem:     resourceParameter(),
                Set:      parameterHash,
                Optional: true,
            },
            "parameter_values": &schema.Schema{
                Type:     schema.TypeMap,
                Optional: true,
            },
        },
    }
}

/*
    ID                  string              `json:"id,omitempty"`
    Name                string              `json:"name"`
    Description         string              `json:"description,omitempty"`
    ParentProjectID     ProjectId           `json:"parentProject,omitempty"`
*/

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*teamcity.Client)

    d.Partial(true)
    parent := d.Get("parent").(string)
    name := d.Get("name").(string)
    project := types.Project{
        ParentProjectID: types.ProjectId(parent),
        Name: name,
        Description: d.Get("description").(string),
    }
    err := client.CreateProject(&project)
    if err != nil {
        return err
    }
    id := project.ID
    d.SetId(id)
    d.SetPartial("parent")
    d.SetPartial("name")
    d.SetPartial("description")

    if parent == "" {
        parent = "_Root"
    }
    var parent_parameters types.Parameters
    if parent_project, err := client.GetProject(parent); err != nil {
        return err
    } else {
        parent_parameters = parent_project.Parameters
    }

    parameters := definitionToParameters(*d.Get("parameters").(*schema.Set))
    for name, _ := range parameters {
        if parent_parameter, ok := parent_parameters[name]; ok && parent_parameter.Spec != nil {
            return fmt.Errorf("Can't redefine parent parameter %s", name)
        }
    }
    for name, v := range d.Get("parameter_values").(map[string]interface{}) {
        value := v.(string)
        parameter, ok := parameters[name]
        if !ok {
            if parameter, ok = parent_parameters[name]; !ok {
                parameter = types.Parameter{
                    Value: value,
                }
            }
        }
        parameter.Value = value
        parameters[name] = parameter
        log.Printf("Parameter value %s => %s", name, parameter.Value)
    }
    log.Printf("Replace Parameters value %q", parameters)
    if err := client.ReplaceAllProjectParameters(id, &parameters); err != nil {
        return err
    }
    d.SetPartial("parameter_values")
    d.SetPartial("parameters")

    d.Partial(false)
    return nil
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
    log.Printf("Reading project resource %q", d.Id())
    client := meta.(*teamcity.Client)
    project, err := client.GetProject(d.Id())
    if err != nil {
        return err
    }
    if project == nil {
        d.SetId("")
        return nil
    }

    parent := project.ParentProjectID
    if parent == "_Root" {
        parent = ""
    } 
    d.Set("parent", parent)
    d.Set("name", project.Name)
    d.Set("description", project.Description)

    var parent_parameters types.Parameters
    if parent_project, err := client.GetProject(string(project.ParentProjectID)); err != nil {
        return err
    } else {
        parent_parameters = parent_project.Parameters
    }
    parameters := project.Parameters
    values := make(map[string]interface{})
    current := d.Get("parameter_values").(map[string]interface{})
    for name, parameter := range project.Parameters {
        if parent_parameter, ok := parent_parameters[name]; ok {
            if parent_parameter.Value != parameter.Value {
                values[name] = parameter.Value
            }
            if parent_parameter.Spec != nil || parameter.Spec == nil {
                delete(parameters, name)
            }
        } else {
            if parameter.Spec == nil {
                delete(parameters, name)
            }
            pwt := types.PasswordType{}
            if parameter.Value != "" {
                values[name] = parameter.Value
            } else if parameter.Spec != nil && parameter.Spec.Type == pwt {
                if value, ok := current[name]; ok && value != "" {
                    values[name] = value
                }
            }
        }
    }
    d.Set("parameters", parametersToDefinition(parameters))
    d.Set("parameter_values", values)

    return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*teamcity.Client)

    id := d.Id()
    d.Partial(true)
    if d.HasChange("parameters") || d.HasChange("parameter_values") {
        parent := d.Get("parent").(string)
        if parent == "" {
            parent = "_Root"
        }
        var parent_parameters types.Parameters
        if parent_project, err := client.GetProject(parent); err != nil {
            return err
        } else {
            parent_parameters = parent_project.Parameters
        }

        o, n := d.GetChange("parameters")
        parameters := definitionToParameters(*n.(*schema.Set))
        old := definitionToParameters(*o.(*schema.Set))
        replace_parameters := make(types.Parameters)
        delete_parameters := old
        for name, parameter := range parameters {
            if parent_parameter, ok := parent_parameters[name]; ok && parent_parameter.Spec != nil {
                return fmt.Errorf("Can't redefine parent parameter %s", name)
            }
            if !reflect.DeepEqual(parameter, old[name]) {
                replace_parameters[name] = parameter
            }
            delete(delete_parameters, name)
        }
        for name, v := range d.Get("parameter_values").(map[string]interface{}) {
            value := v.(string)
            parameter, ok := parameters[name]
            if !ok {
                if parameter, ok = parent_parameters[name]; !ok {
                    parameter = types.Parameter{
                        Value: value,
                    }
                }
            }
            parameter.Value = value
            replace_parameters[name] = parameter
        }
        for name, _ := range delete_parameters {
            if err := client.DeleteProjectParameter(id, name); err != nil {
                return err
            }
        }
        for name, parameter := range replace_parameters {
            if err := client.ReplaceProjectParameter(id, name, &parameter); err != nil {
                return err
            }
        }
        d.SetPartial("parameter_values")
        d.SetPartial("parameters")
    }

    return nil
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*teamcity.Client)
    return client.DeleteProject(d.Id())
}