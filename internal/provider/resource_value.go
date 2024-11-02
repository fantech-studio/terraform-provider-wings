package provider

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"fantech.dev/terraform-provider-wings/internal/model"
)

var (
	_ resource.Resource = &ValueResource{}
)

func NewValueResource() resource.Resource {
	return &ValueResource{}
}

type ValueResource struct {
	c *config
}

type (
	valueResource struct {
		ID             types.String             `tfsdk:"id"`
		ValueID        types.String             `tfsdk:"value_id"`
		Description    types.String             `tfsdk:"description"`
		Enabled        types.Bool               `tfsdk:"enabled"`
		DefaultVariant types.String             `tfsdk:"default_variant"`
		Bool           []valueResourceBool      `tfsdk:"bool"`
		Int            []valueResourceInt       `tfsdk:"int"`
		String         []valueResourceString    `tfsdk:"string"`
		Object         []valueResourceObject    `tfsdk:"object"`
		Targeting      []valueResourceTargeting `tfsdk:"targeting"`
		Test           []valueResourceTest      `tfsdk:"test"`
	}

	valueResourceBool struct {
		Variant types.String `tfsdk:"variant"`
		Value   types.Bool   `tfsdk:"value"`
	}

	valueResourceInt struct {
		Variant types.String `tfsdk:"variant"`
		Value   types.Int64  `tfsdk:"value"`
	}

	valueResourceString struct {
		Variant types.String `tfsdk:"variant"`
		Value   types.String `tfsdk:"value"`
	}

	valueResourceObject struct {
		Variant   types.String             `tfsdk:"variant"`
		Value     types.String             `tfsdk:"value"`
		Transform []valueResourceTransform `tfsdk:"transform"`
	}

	valueResourceTargeting struct {
		Variant types.String `tfsdk:"variant"`
		Expr    types.String `tfsdk:"expr"`
	}

	valueResourceTest struct {
		Variables types.String `tfsdk:"variables"`
		Expected  types.String `tfsdk:"expected"`
	}

	valueResourceTransform struct {
		Expr types.String `tfsdk:"expr"`
	}
)

func (v *ValueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	value, err := v.c.GetValue(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error get value", err.Error())
		return
	}

	state := valueState(value)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (v *ValueResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_value"
}

func (v *ValueResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Wings value resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Computed ID.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"value_id": schema.StringAttribute{
				Description: "The ID of this Value.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"enabled": schema.BoolAttribute{
				Required: true,
			},
			"default_variant": schema.StringAttribute{
				Required: true,
			},
		},
		Blocks: map[string]schema.Block{
			"bool": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"variant": schema.StringAttribute{
							Required: true,
						},
						"value": schema.BoolAttribute{
							Required: true,
						},
					},
				},
			},
			"string": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"variant": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"object": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"variant": schema.StringAttribute{
							Required: true,
						},
						"value": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.RegexMatches(
									regexp.MustCompile(`^{.*}$`),
									"Must be map object, not array",
								),
							},
						},
					},
					Blocks: map[string]schema.Block{
						"transform": schema.ListNestedBlock{
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"expr": schema.StringAttribute{
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"int": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"variant": schema.StringAttribute{
							Required: true,
						},
						"value": schema.Int64Attribute{
							Required: true,
						},
					},
				},
			},
			"targeting": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"variant": schema.StringAttribute{
							Required: true,
						},
						"expr": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"test": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"variables": schema.StringAttribute{
							Required: true,
						},
						"expected": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

func (v *valueResource) value() (*model.Value, error) {
	variants := model.Variants{}
	for _, val := range v.Bool {
		variants[val.Variant.ValueString()] = model.ValueEvaluation{
			Bool: &model.Bool{
				Value: val.Value.ValueBool(),
			},
		}
	}
	for _, val := range v.String {
		variants[val.Variant.ValueString()] = model.ValueEvaluation{
			String: &model.String{
				Value: val.Value.ValueString(),
			},
		}
	}
	for _, val := range v.Object {
		m := make(map[string]any)
		err := json.Unmarshal([]byte(val.Value.ValueString()), &m)
		if err != nil {
			return nil, err
		}
		transforms := make([]*model.ValueTransform, 0, len(val.Transform))
		for _, t := range val.Transform {
			transforms = append(transforms, &model.ValueTransform{
				Expr: t.Expr.ValueString(),
			})
		}
		variants[val.Variant.ValueString()] = model.ValueEvaluation{
			Object: &model.Object{
				Value:      m,
				Transforms: transforms,
			},
		}
	}
	for _, val := range v.Int {
		variants[val.Variant.ValueString()] = model.ValueEvaluation{
			Int: &model.Int{
				Value: val.Value.ValueInt64(),
			},
		}
	}

	rules := make([]model.ValueTargetingRule, 0, len(v.Targeting))
	for _, t := range v.Targeting {
		rules = append(rules, model.ValueTargetingRule{
			Variant: t.Variant.ValueString(),
			Expr:    t.Expr.ValueString(),
		})
	}

	tests := make([]*model.EvaluationTest, 0, len(v.Test))
	for _, t := range v.Test {
		m := make(map[string]any)
		err := json.Unmarshal([]byte(t.Variables.ValueString()), &m)
		if err != nil {
			return nil, err
		}
		tests = append(tests, &model.EvaluationTest{
			Variables: m,
			Expected:  t.Expected.ValueString(),
		})
	}
	value := &model.Value{
		ID:             v.ID.ValueString(),
		Enabled:        v.Enabled.ValueBool(),
		Description:    v.Description.ValueString(),
		DefaultVariant: v.DefaultVariant.ValueString(),
		Variants:       variants,
		Targeting: model.Targeting{
			Rules: rules,
		},
		Tests: tests,
	}
	return value, nil
}

func valueState(v *model.Value) *valueResource {
	var (
		bools   []valueResourceBool
		strs    []valueResourceString
		objects []valueResourceObject
		ints    []valueResourceInt
	)

	for k, val := range v.Variants {
		if val.Bool != nil {
			if bools == nil {
				bools = make([]valueResourceBool, 0, len(v.Variants))
			}
			bools = append(bools, valueResourceBool{
				Variant: types.StringValue(k),
				Value:   types.BoolValue(val.Bool.Value),
			})
		}
		if val.String != nil {
			if strs == nil {
				strs = make([]valueResourceString, 0, len(v.Variants))
			}
			strs = append(strs, valueResourceString{
				Variant: types.StringValue(k),
				Value:   types.StringValue(val.String.Value),
			})
		}
		if val.Object != nil {
			if objects == nil {
				objects = make([]valueResourceObject, 0, len(v.Variants))
			}
			b, _ := json.Marshal(val.Object.Value)
			transforms := make([]valueResourceTransform, 0, len(val.Object.Transforms))
			for _, t := range val.Object.Transforms {
				transforms = append(transforms, valueResourceTransform{
					Expr: types.StringValue(t.Expr),
				})
			}
			objects = append(objects, valueResourceObject{
				Variant:   types.StringValue(k),
				Value:     types.StringValue(string(b)),
				Transform: transforms,
			})
		}
		if val.Int != nil {
			if ints == nil {
				ints = make([]valueResourceInt, 0, len(v.Variants))
			}
			ints = append(ints, valueResourceInt{
				Variant: types.StringValue(k),
				Value:   types.Int64Value(val.Int.Value),
			})
		}
	}

	targeting := make([]valueResourceTargeting, 0, len(v.Targeting.Rules))
	for _, t := range v.Targeting.Rules {
		targeting = append(targeting, valueResourceTargeting{
			Variant: types.StringValue(t.Variant),
			Expr:    types.StringValue(t.Expr),
		})
	}

	tests := make([]valueResourceTest, 0, len(v.Tests))
	for _, t := range v.Tests {
		b, _ := json.Marshal(t.Variables)
		tests = append(tests, valueResourceTest{
			Variables: types.StringValue(string(b)),
			Expected:  types.StringValue(t.Expected),
		})
	}

	return &valueResource{
		ID:             types.StringValue(v.ID),
		ValueID:        types.StringValue(v.ID),
		Description:    types.StringValue(v.Description),
		Enabled:        types.BoolValue(v.Enabled),
		DefaultVariant: types.StringValue(v.DefaultVariant),
		Bool:           bools,
		String:         strs,
		Object:         objects,
		Int:            ints,
		Targeting:      targeting,
		Test:           tests,
	}
}

func (v *ValueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan valueResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	value, err := plan.value()
	if err != nil {
		resp.Diagnostics.AddError("Error creating value", "Invalid Attribute(s): "+err.Error())
		return
	}

	value, err = v.c.CreateValue(ctx, value)
	if err != nil {
		resp.Diagnostics.AddError("Error creating value", err.Error())
		return
	}

	plan.ID = types.StringValue(value.ID)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (v *ValueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state valueResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (v *ValueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan valueResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	value, err := plan.value()
	if err != nil {
		resp.Diagnostics.AddError("Error updating value", "Invalid Attribute(s): "+err.Error())
		return
	}

	_, err = v.c.UpdateValue(ctx, value)
	if err != nil {
		resp.Diagnostics.AddError("Error updating value", err.Error())
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (v *ValueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state valueResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := v.c.DeleteValue(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting value", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func (v *ValueResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	v.c = req.ProviderData.(*config)
}
