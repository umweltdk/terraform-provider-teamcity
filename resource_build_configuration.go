package main

import (
    "errors"
    "fmt"

    "github.com/hashicorp/terraform/helper/schema"
    "github.com/hashicorp/terraform/helper/hashcode"

    "github.com/umweltdk/teamcity/teamcity"
    "github.com/umweltdk/teamcity/types"

    "log"
    "reflect"
)


func resourceBuildStep() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "type": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
            },
            "name": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
            },
            "properties": &schema.Schema{
                Type:     schema.TypeMap,
                Optional: true,
            },
        },
    }
}

func resourceAttachedVcsRoot() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "vcs_root": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
                ValidateFunc: teamcity.ValidateID,
            },
            "checkout_rules": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
            },
        },
    }
}

func resourceBuildConfiguration() *schema.Resource {
    return &schema.Resource{
        Create: resourceBuildConfigurationCreate,
        Read:   resourceBuildConfigurationRead,
        Update: resourceBuildConfigurationUpdate,
        Delete: resourceBuildConfigurationDelete,

        Schema: map[string]*schema.Schema{
            "project": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
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
            "template": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
                ValidateFunc: teamcity.ValidateID,
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
            "steps": &schema.Schema{
                Type:     schema.TypeList,
                Elem:     resourceBuildStep(),
                Optional: true,
            },
            "attached_vcs_roots": &schema.Schema{
                Type:     schema.TypeSet,
                Elem:     resourceAttachedVcsRoot(),
                Set:      attachedVcsRootHash,
                Optional: true,
            },
        },
    }
}

func resourceBuildTemplate() *schema.Resource {
    return &schema.Resource{
        Create: resourceBuildTemplateCreate,
        Read:   resourceBuildTemplateRead,
        Update: resourceBuildTemplateUpdate,
        Delete: resourceBuildTemplateDelete,

        Schema: map[string]*schema.Schema{
            "project": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
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
            "steps": &schema.Schema{
                Type:     schema.TypeList,
                Elem:     resourceBuildStep(),
                Optional: true,
            },
            "attached_vcs_roots": &schema.Schema{
                Type:     schema.TypeSet,
                Elem:     resourceAttachedVcsRoot(),
                Set:      attachedVcsRootHash,
                Optional: true,
            },
        },
    }
}

func resourceBuildConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
    return resourceBuildConfigurationCreateInternal(d, meta, false)
}

func resourceBuildConfigurationRead(d *schema.ResourceData, meta interface{}) error {
    return resourceBuildConfigurationReadInternal(d, meta, false)
}

func resourceBuildConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
    return resourceBuildConfigurationUpdateInternal(d, meta, false)
}

func resourceBuildConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
    return resourceBuildConfigurationDeleteInternal(d, meta, false)
}

func resourceBuildTemplateCreate(d *schema.ResourceData, meta interface{}) error {
    return resourceBuildConfigurationCreateInternal(d, meta, true)
}

func resourceBuildTemplateRead(d *schema.ResourceData, meta interface{}) error {
    return resourceBuildConfigurationReadInternal(d, meta, true)
}

func resourceBuildTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
    return resourceBuildConfigurationUpdateInternal(d, meta, true)
}

func resourceBuildTemplateDelete(d *schema.ResourceData, meta interface{}) error {
    return resourceBuildConfigurationDeleteInternal(d, meta, true)
}

/*
    ID             string         `json:"id,omitempty"`
    ProjectID      string         `json:"projectId"`
    TemplateFlag   bool           `json:"templateFlag"`
    Template       *TemplateID    `json:"template,omitempty"`
    Name           string         `json:"name"`
    Description    string         `json:"description,omitempty"`
    Settings       Properties     `json:"settings,omitempty"`
    Parameters     Properties     `json:"parameters,omitempty"`
    Steps          BuildSteps     `json:"steps,omitempty"`
    VcsRootEntries VcsRootEntries `json:"vcs-root-entries,omitempty"`
*/

func resourceBuildConfigurationCreateInternal(d *schema.ResourceData, meta interface{}, template bool) error {
    client := meta.(*teamcity.Client)

    projectID := d.Get("project").(string)
    name := d.Get("name").(string)
    d.Partial(true)
    templateID := ""
    if !template {
        templateID = d.Get("template").(string)
    }
    config := types.BuildConfiguration{
        ProjectID: projectID,
        TemplateFlag: template,
        TemplateID: types.TemplateId(templateID),
        Name: name,
        Description: d.Get("description").(string),
        Steps: resourceBuildSteps(d.Get("steps").([]interface{})),
    }
    
    if err := client.CreateBuildConfiguration(&config); err != nil {
        return err
    }
    id := config.ID
    d.SetId(id)
    d.SetPartial("project")
    d.SetPartial("name")
    d.SetPartial("description")
    if !template {
        d.SetPartial("template")
    }
    d.SetPartial("steps")

    var project_parameters types.Parameters
    if project, err := client.GetProject(projectID); err != nil {
        return err
    } else {
        project_parameters = project.Parameters
    }
    template_parameters := make(types.Parameters)
    if templateID != "" {
        if teamplate_config, err := client.GetBuildConfiguration(templateID); err != nil {
            return err
        } else {
            template_parameters = teamplate_config.Parameters
        }
    }

    parameters := definitionToParameters(*d.Get("parameters").(*schema.Set))
    for name, _ := range parameters {
        if project_parameter, ok := project_parameters[name]; ok && project_parameter.Spec != nil {
            return fmt.Errorf("Can't redefine project parameter %s", name)
        }
        if template_parameter, ok := template_parameters[name]; ok && template_parameter.Spec != nil {
            return fmt.Errorf("Can't redefine template parameter %s", name)
        }
    }
    for name, v := range d.Get("parameter_values").(map[string]interface{}) {
        value := v.(string)
        parameter, ok := parameters[name]
        if !ok {
            if parameter, ok = project_parameters[name]; !ok {
                if parameter, ok = template_parameters[name]; !ok {
                    parameter = types.Parameter{
                        Value: value,
                    }
                }
            }
        }
        parameter.Value = value
        parameters[name] = parameter
        log.Printf("Parameter value %s => %s", name, parameter.Value)
    }
    log.Printf("Replace Parameters value %q", parameters)
    if err := client.ReplaceAllBuildConfigurationParameters(id, &parameters); err != nil {
        return err
    }
    d.SetPartial("parameter_values")
    d.SetPartial("parameters")

    for _, root := range resourceAttachedVcsRoots(*d.Get("attached_vcs_roots").(*schema.Set)) {
        err := client.AttachBuildConfigurationVcsRoot(id, &root)
        if err != nil {
            return err
        }
    }
    d.SetPartial("attached_vcs_roots")

    d.Partial(false)
    return nil
}

func resourceBuildConfigurationReadInternal(d *schema.ResourceData, meta interface{}, template bool) error {
    log.Printf("Reading resource %q", d.Id())
    client := meta.(*teamcity.Client)
    config, err := client.GetBuildConfiguration(d.Id())
    if err != nil {
        return err
    }
    if config == nil || template != config.TemplateFlag {
        d.SetId("")
        return nil
    }

    log.Printf("Reading resource %q\n%q", d.Id(), d.Get("parameters"))
    d.Set("project", config.ProjectID)
    d.Set("name", config.Name)
    d.Set("description", config.Description)
    if !template {
        d.Set("template", config.TemplateID)
    }

    steps := make([]map[string]interface{}, 0)
    for _, step := range config.Steps {
        v := make(map[string]interface{})
        v["type"] = step.Type
        if step.Name != "" {
            v["name"] = step.Name
        }
        properties := make(map[string]interface{})
        for name, prop := range step.Properties {
            properties[name] = prop
        }
        if len(properties) > 0 {
            v["properties"] = properties
        }
        steps = append(steps, v)
    }
    log.Printf("[INFO] Steps %q\n", steps)
    d.Set("steps", steps)


    var project_parameters types.Parameters
    if project, err := client.GetProject(string(config.ProjectID)); err != nil {
        return err
    } else {
        project_parameters = project.Parameters
    }
    templateID := string(config.TemplateID)
    template_parameters := make(types.Parameters)
    if templateID != "" {
        if teamplate_config, err := client.GetBuildConfiguration(templateID); err != nil {
            return err
        } else {
            template_parameters = teamplate_config.Parameters
        }
    }
    parameters := config.Parameters
    values := make(map[string]interface{})
    current := d.Get("parameter_values").(map[string]interface{})
    for name, parameter := range config.Parameters {
        if project_parameter, ok := project_parameters[name]; ok {
            if project_parameter.Value != parameter.Value {
                values[name] = parameter.Value
            }
            if project_parameter.Spec != nil || parameter.Spec == nil {
                delete(parameters, name)
            }
        } else if template_parameter, ok := template_parameters[name]; ok {
            if template_parameter.Value != parameter.Value {
                values[name] = parameter.Value
            }
            if template_parameter.Spec != nil || parameter.Spec == nil {
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

func resourceBuildSteps(steps []interface{}) types.BuildSteps {
    tcSteps := make(types.BuildSteps, 0)
    for _, s := range steps {
        step := s.(map[string]interface{})
        typeName := step["type"].(string)
        name := step["name"].(string)
        properties := step["properties"].(map[string]interface{})
        actualProps := make(types.Properties)
        for k, v := range properties {
            actualProps[k] = v.(string)
        }

        tcSteps = append(tcSteps, types.BuildStep{
            Type: typeName,
            Name: name,
            Properties: actualProps,
        })
    }

    return tcSteps
}

func resourceAttachedVcsRoots(vcsRoots schema.Set) types.VcsRootEntries {
    tcRoots := make(types.VcsRootEntries, 0)
    for _, s := range vcsRoots.List() {
        entry := s.(map[string]interface{})
        vcsRoot := entry["vcs_root"].(string)
        rules := entry["checkout_rules"].(string)

        tcRoots = append(tcRoots, types.VcsRootEntry{
            VcsRootID: types.VcsRootId(vcsRoot),
            CheckoutRules: rules,
        })
    }

    return tcRoots
}

func resourceBuildConfigurationUpdateInternal(d *schema.ResourceData, meta interface{}, template bool) error {
    client := meta.(*teamcity.Client)
    //var err error
    id := d.Id()
    d.Partial(true)
    if d.HasChange("parameters") {
       projectID := d.Get("project").(string)
        var project_parameters types.Parameters
        if project, err := client.GetProject(projectID); err != nil {
            return err
        } else {
            project_parameters = project.Parameters
        }

        o, n := d.GetChange("parameters")
        parameters := definitionToParameters(*n.(*schema.Set))
        old := definitionToParameters(*o.(*schema.Set))
        replace_parameters := make(types.Parameters)
        delete_parameters := old
        for name, parameter := range parameters {
            if project_parameter, ok := project_parameters[name]; ok && project_parameter.Spec != nil {
                return fmt.Errorf("Can't redefine project parameter %s", name)
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
                if parameter, ok = project_parameters[name]; !ok {
                    parameter = types.Parameter{
                        Value: value,
                    }
                }
            }
            parameter.Value = value
            replace_parameters[name] = parameter
        }
        for name, _ := range delete_parameters {
            if err := client.DeleteBuildConfigurationParameter(id, name); err != nil {
                return err
            }
        }
        for name, parameter := range replace_parameters {
            if err := client.ReplaceBuildConfigurationParameter(id, name, &parameter); err != nil {
                return err
            }
        }
        d.SetPartial("parameter_values")
        d.SetPartial("parameters")
    }

    if d.HasChange("attached_vcs_roots") {
        old, n := d.GetChange("attached_vcs_roots")
        existing := make(map[types.VcsRootId]bool)

        for _, root := range resourceAttachedVcsRoots(*n.(*schema.Set)) {
            err := client.AttachBuildConfigurationVcsRoot(id, &root)
            if err != nil {
                return err
            }
            existing[root.VcsRootID] = true
        }
        for _, root := range resourceAttachedVcsRoots(*old.(*schema.Set)) {
            if !existing[root.VcsRootID] {
                err := client.DetachBuildConfigurationVcsRoot(id, string(root.VcsRootID))
                if err != nil {
                    return err
                }
            }
        }

        d.SetPartial("attached_vcs_roots")
    }
    if d.HasChange("description") {
        return errors.New("Description update not supported")
        /*
        if err := client.SetBuildConfigurationDescription(d.Id(), d.Get("description").(string)); err != nil {
            return err
        }
        d.SetPartial("description")
        */
    }
    if !template && d.HasChange("template") {
        return errors.New("Template update not supported")
        /*
        if err := client.SetBuildConfigurationTemplate(d.Id(), d.Get("template").(string)); err != nil {
            return err
        }
        d.SetPartial("template")
        */
    }
    if d.HasChange("steps") {
        return errors.New("Build steps updating not supported")
        /*
        if err := client.ReplaceAllBuildConfigurationSteps(d.Id(), resourceBuildSteps(d.Get("steps").([]interface{}))); err != nil {
            return err
        }
        d.SetPartial("steps")
        */
    }

    d.Partial(false)
    return nil
}

func resourceBuildConfigurationDeleteInternal(d *schema.ResourceData, meta interface{}, template bool) error {
    client := meta.(*teamcity.Client)
    return client.DeleteBuildConfiguration(d.Id())
}

func attachedVcsRootHash(v interface{}) int {
    m := v.(map[string]interface{})
    return hashcode.String(m["vcs_root"].(string))
}
