package qeycloak

import (
	"fmt"

	"github.com/netascode/go-aci"
	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/tidwall/gjson"
)

// This is the type that implements the Rest resource methods.
type Rest struct{}

// The following statement is not required. It is a type assertion to indicate to Go that Command
// implements the following interfaces. If the function signature doesn't match or isn't implemented,
// we get nice compile time errors at this location.

var _ = (infer.Annotated)((*Rest)(nil))
var _ = (infer.CustomResource[RestInputs, RestOutputs])((*Rest)(nil))
var _ = (infer.CustomDiff[RestInputs, RestOutputs])((*Rest)(nil))
var _ = (infer.CustomRead[RestInputs, RestOutputs])((*Rest)(nil))
var _ = (infer.CustomUpdate[RestInputs, RestOutputs])((*Rest)(nil))
var _ = (infer.CustomDelete[RestOutputs])((*Rest)(nil))
var _ = (infer.ExplicitDependencies[RestInputs, RestOutputs])((*Rest)(nil))

// Implementing Annotate lets you provide descriptions and default values for resources and they will
// be visible in the provider's schema and the generated SDKs.
func (c *Rest) Annotate(a infer.Annotator) {
	a.Describe(&c, "Manages ACI Model Objects via REST API calls. This resource can only manage a single API object.")
}

// These are the inputs (or arguments) to a Rest resource.
type RestInputs struct {
	// The field tags are used to provide metadata on the schema representation.
	// pulumi:"optional" specifies that a field is optional. This must be a pointer.
	// provider:"replaceOnChanges" specifies that the resource will be replaced if the field changes.
	Dn        string             `pulumi:"dn"`
	ClassName string             `pulumi:"class_name"`
	Content   *map[string]string `pulumi:"content,optional"`
	Children  *[]Child           `pulumi:"children,optional"`
}

type Child struct {
	Rn        string             `pulumi:"rn"`
	ClassName string             `pulumi:"class_name"`
	Content   *map[string]string `pulumi:"content,optional"`
}

// Annotate lets you provide descriptions and default values for fields and they will
// be visible in the provider's schema and the generated SDKs.
func (c *RestInputs) Annotate(a infer.Annotator) {
	a.Describe(&c.Dn, "Distinguished name of object being managed including its relative name, e.g. uni/tn-EXAMPLE_TENANT.")
	a.Describe(&c.ClassName, "Which class object is being created. (Make sure there is no colon in the classname)")
	a.Describe(&c.Content, "Map of key-value pairs those needed to be passed to the Model object as parameters. Make sure the key name matches the name with the object parameter in ACI.")
	a.Describe(&c.Children, "List of child objects to be created. Each child object must have a unique relative name.")
}

func (c *Child) Annotate(a infer.Annotator) {
	a.Describe(&c.Rn, "Relative name of child object.")
	a.Describe(&c.ClassName, "Which class object is being created. (Make sure there is no colon in the classname)")
	a.Describe(&c.Content, "Map of key-value pairs those needed to be passed to the Model object as parameters. Make sure the key name matches the name with the object parameter in ACI.")
}

// These are the outputs (or properties) of a Rest resource.
type RestOutputs struct {
	RestInputs
}

// This is the Create method. This will be run on every Rest resource creation.
func (c *Rest) Create(ctx p.Context, name string, input RestInputs, preview bool) (string, RestOutputs, error) {
	state := RestOutputs{RestInputs: input}
	id := input.Dn

	// If in preview, don't run the command.
	if preview {
		return id, state, nil
	}

	client := infer.GetConfig[Config](ctx).Client
	body := aci.Body{}
	if input.Content != nil {
		for k, v := range *input.Content {
			body = body.Set(input.ClassName+".attributes."+k, v)
		}
	}
	if input.Children != nil {
		for _, child := range *input.Children {
			if child.Content != nil {
				childBody := aci.Body{}
				for ck, cv := range *child.Content {
					childBody = childBody.Set(child.ClassName+".attributes."+ck, cv)
				}
				body = body.SetRaw(input.ClassName+".children.-1", childBody.Str)
			}
		}
	}
	_, err := client.Post(id, body.Str)

	return id, state, err
}

// This is the Diff method.
func (c *Rest) Diff(ctx p.Context, id string, olds RestOutputs, news RestInputs) (p.DiffResponse, error) {
	diffs := make(map[string]p.PropertyDiff)
	hasChanges := false
	if olds.Dn != news.Dn {
		diffs["dn"] = p.PropertyDiff{Kind: p.UpdateReplace}
		hasChanges = true
	}
	if olds.ClassName != news.ClassName {
		diffs["class_name"] = p.PropertyDiff{Kind: p.UpdateReplace}
		hasChanges = true
	}
	if olds.Content != nil {
		for k := range *olds.Content {
			if _, ok := (*news.Content)[k]; !ok {
				diffs["content."+k] = p.PropertyDiff{Kind: p.Delete}
				hasChanges = true
			} else if (*olds.Content)[k] != (*news.Content)[k] {
				diffs["content."+k] = p.PropertyDiff{Kind: p.Update}
				hasChanges = true
			}
		}
	}
	if news.Content != nil {
		for k := range *news.Content {
			if _, ok := (*olds.Content)[k]; !ok {
				diffs["content."+k] = p.PropertyDiff{Kind: p.Add}
				hasChanges = true
			}
		}
	}
	if olds.Children != nil {
		for i, child := range *olds.Children {
			rn := child.Rn
			if news.Children != nil {
				for _, newChild := range *news.Children {
					if rn == newChild.Rn {
						if child.ClassName != newChild.ClassName {
							diffs[fmt.Sprintf("children[%v].class_name", i)] = p.PropertyDiff{Kind: p.UpdateReplace}
							hasChanges = true
						}
						if child.Content != nil {
							for ck := range *child.Content {
								if _, ok := (*newChild.Content)[ck]; !ok {
									diffs[fmt.Sprintf("children[%v].content.%v", i, ck)] = p.PropertyDiff{Kind: p.Delete}
									hasChanges = true
								} else if (*child.Content)[ck] != (*newChild.Content)[ck] {
									diffs[fmt.Sprintf("children[%v].content.%v", i, ck)] = p.PropertyDiff{Kind: p.Update}
									hasChanges = true
								}
							}
						}
						if newChild.Content != nil {
							for ck := range *newChild.Content {
								if _, ok := (*child.Content)[ck]; !ok {
									diffs[fmt.Sprintf("children[%v].content.%v", i, ck)] = p.PropertyDiff{Kind: p.Add}
									hasChanges = true
								}
							}
						}
					}
				}
			}
		}
	}
	if news.Children != nil {
		for _, child := range *news.Children {
			rn := child.Rn
			found := false
			if olds.Children != nil {
				for _, oldChild := range *olds.Children {
					if rn == oldChild.Rn {
						found = true
						break
					}
				}
			}
			if !found {
				diffs["children"] = p.PropertyDiff{Kind: p.Add}
				hasChanges = true
			}
		}
	}
	resp := p.DiffResponse{DeleteBeforeReplace: true, HasChanges: hasChanges, DetailedDiff: diffs}
	return resp, nil
}

// WireDependencies is relevant to secrets handling. This method indicates which Inputs
// the Outputs are derived from. If an output is derived from a secret input, the output
// will be a secret.

// This naive implementation conveys that every output is derived from its respective inputs.
func (r *Rest) WireDependencies(f infer.FieldSelector, args *RestInputs, state *RestOutputs) {

	dnInput := f.InputField(&args.Dn)
	classNameInput := f.InputField(&args.ClassName)
	contentInput := f.InputField(&args.Content)
	childrenInput := f.InputField(&args.Children)

	f.OutputField(&state.Dn).DependsOn(dnInput)
	f.OutputField(&state.ClassName).DependsOn(classNameInput)
	f.OutputField(&state.Content).DependsOn(contentInput)
	f.OutputField(&state.Children).DependsOn(childrenInput)
}

// This is the Read method.
func (c *Rest) Read(ctx p.Context, id string, input RestInputs, state RestOutputs) (string, RestInputs, RestOutputs, error) {

	client := infer.GetConfig[Config](ctx).Client
	var resp gjson.Result
	var err error
	if input.Children != nil {
		resp, err = client.GetDn(id, aci.Query("rsp-subtree", "children"))
	} else {
		resp, err = client.GetDn(id)
	}
	content := *state.Content
	if input.Content != nil {
		for k := range *input.Content {
			content[k] = resp.Get(input.ClassName + ".attributes." + k).String()
		}
	}
	if input.Children != nil {
		var children []Child
		for _, child := range *input.Children {
			stateChild := Child{Rn: child.Rn, ClassName: child.ClassName}
			childContent := make(map[string]string)
			if child.Content != nil {
				rChildren := resp.Get(input.ClassName + ".children")
				rChildren.ForEach(func(k, v gjson.Result) bool {
					if v.Get(child.ClassName+".attributes.rn").String() == child.Rn {
						for ck := range *child.Content {
							childContent[ck] = v.Get(child.ClassName + ".attributes." + ck).String()
						}
						stateChild.Content = &childContent
						return false // stop iterating
					}
					return true // keep iterating
				})
			}
			children = append(children, stateChild)
		}
		state.Children = &children
	}

	return id, input, state, err
}

// The Update method will be run on every update.
func (c *Rest) Update(ctx p.Context, id string, olds RestOutputs, news RestInputs, preview bool) (RestOutputs, error) {
	state := RestOutputs{RestInputs: news}

	// If in preview, don't run the command.
	if preview {
		return state, nil
	}

	client := infer.GetConfig[Config](ctx).Client
	body := aci.Body{}
	if news.Content != nil {
		for k, v := range *news.Content {
			body = body.Set(news.ClassName+".attributes."+k, v)
		}
	}
	if news.Children != nil {
		for _, child := range *news.Children {
			if child.Content != nil {
				childBody := aci.Body{}
				for ck, cv := range *child.Content {
					childBody = childBody.Set(child.ClassName+".attributes."+ck, cv)
				}
				body = body.SetRaw(news.ClassName+".children.-1", childBody.Str)
			}
		}
	}
	_, err := client.Post(id, body.Str)

	return state, err
}

// The Delete method will run when the resource is deleted.
func (c *Rest) Delete(ctx p.Context, id string, props RestOutputs) error {
	client := infer.GetConfig[Config](ctx).Client
	_, err := client.DeleteDn(id)
	return err
}
