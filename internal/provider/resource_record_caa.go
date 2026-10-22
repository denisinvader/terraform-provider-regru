package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"terraform-provider-regru/internal/md"
	"terraform-provider-regru/internal/provider/validators"
	"terraform-provider-regru/internal/regru"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &resourceRecordCAA{}
	_ resource.ResourceWithConfigure   = &resourceRecordCAA{}
	_ resource.ResourceWithImportState = &resourceRecordCAA{}
)

func NewResourceRecordCAA() resource.Resource {
	return &resourceRecordCAA{}
}

type resourceRecordCAA struct {
	client *regru.Client
}

type resourceRecordCAAModel struct {
	Zone  types.String `tfsdk:"zone"`
	Name  types.String `tfsdk:"name"`
	Flag  types.String `tfsdk:"flag"`
	Tag   types.String `tfsdk:"tag"`
	Value types.String `tfsdk:"value"`
}

func (r *resourceRecordCAA) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_record_caa"
}

func (r *resourceRecordCAA) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*regru.Client)
	if !ok {
		res.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected *regru.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *resourceRecordCAA) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: md.ResourceRecordCAA,
		Attributes: map[string]schema.Attribute{
			"zone": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: md.RecordZone,
				Validators: []validator.String{
					validators.StringWithoutWhitespace(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Subdomain name, `*` for a wildcard, or `@` for the zone apex.",
				Validators: []validator.String{
					validators.StringWithoutWhitespace(),
				},
			},
			"flag": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: md.ResourceRecordCAAFlag,
				Validators: []validator.String{
					stringvalidator.OneOf("0", "128"),
				},
			},
			"tag": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: md.ResourceRecordCAATag,
				Validators: []validator.String{
					stringvalidator.OneOf("issue", "issuewild", "iodef"),
				},
			},
			"value": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: md.ResourceRecordCAAValue,
				Validators:          []validator.String{},
			},
		},
	}
}

func (r *resourceRecordCAA) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan resourceRecordCAAModel

	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone  = plan.Zone.ValueString()
		name  = plan.Name.ValueString()
		flag  = plan.Flag.ValueString()
		tag   = plan.Tag.ValueString()
		value = plan.Value.ValueString()
	)

	err := r.client.CreateRecordCAA(zone, name, flag, tag, value)
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating reg.ru CAA record",
			fmt.Sprintf(
				"Could not create reg.ru CAA record %s -> %s %s %q in zone %s: %s",
				name,
				flag,
				tag,
				value,
				zone,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, plan)...)
}

func (r *resourceRecordCAA) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state resourceRecordCAAModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone  = state.Zone.ValueString()
		name  = state.Name.ValueString()
		flag  = state.Flag.ValueString()
		tag   = state.Tag.ValueString()
		value = state.Value.ValueString()
	)

	_, err := r.client.GetRecordCAA(zone, name, flag, tag, value)
	if err != nil {
		var notFoundError regru.RegruRecordNotFoundError
		if errors.As(err, &notFoundError) {
			res.State.RemoveResource(ctx)
			return
		}

		res.Diagnostics.AddError(
			"Error reading reg.ru CAA record",
			fmt.Sprintf(
				"Could not read reg.ru CAA record %s -> %s %s %q in zone %s: %s",
				name,
				flag,
				tag,
				value,
				zone,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}

func (r *resourceRecordCAA) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	res.Diagnostics.AddError(
		"Error updating reg.ru CAA record",
		"Terraform called Update for a resource that only supports replacement. Please report this issue to the provider developers.",
	)
}

func (r *resourceRecordCAA) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state resourceRecordCAAModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone  = state.Zone.ValueString()
		name  = state.Name.ValueString()
		flag  = state.Flag.ValueString()
		tag   = state.Tag.ValueString()
		value = state.Value.ValueString()
	)

	err := r.client.DeleteRecordCAA(zone, name, flag, tag, value)
	if err != nil {
		res.Diagnostics.AddError(
			"Error deleting reg.ru CAA record",
			fmt.Sprintf(
				"Could not delete reg.ru CAA record %s -> %s %s %q in zone %s: %s",
				name,
				flag,
				tag,
				value,
				zone,
				err,
			),
		)
	}
}

func (r *resourceRecordCAA) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {

	d := strings.SplitN(req.ID, ",", 5)

	if len(d) < 5 {
		res.Diagnostics.AddError(
			"Error parsing import ID for CAA record",
			fmt.Sprintf("Provided ID %s does not match expected format 'zone,name,flag,tag,value'", req.ID),
		)
		return
	}

	var (
		zone  = d[0]
		name  = d[1]
		flag  = d[2]
		tag   = d[3]
		value = d[4]
	)

	_, err := r.client.GetRecordCAA(zone, name, flag, tag, value)
	if err != nil {
		res.Diagnostics.AddError(
			"Error reading reg.ru CAA record",
			fmt.Sprintf(
				"Could not read reg.ru CAA record %s -> %s %s %q in zone %s: %s",
				name,
				flag,
				tag,
				value,
				zone,
				err,
			),
		)
		return
	}

	state := resourceRecordCAAModel{
		Zone:  types.StringValue(zone),
		Name:  types.StringValue(name),
		Flag:  types.StringValue(flag),
		Tag:   types.StringValue(tag),
		Value: types.StringValue(value),
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}
