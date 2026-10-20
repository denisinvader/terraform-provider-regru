package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"terraform-provider-regru/internal/md"
	"terraform-provider-regru/internal/provider/validators"
	"terraform-provider-regru/internal/regru"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &resourceRecordCNAME{}
	_ resource.ResourceWithConfigure   = &resourceRecordCNAME{}
	_ resource.ResourceWithImportState = &resourceRecordCNAME{}
)

func NewResourceRecordCNAME() resource.Resource {
	return &resourceRecordCNAME{}
}

type resourceRecordCNAME struct {
	client *regru.Client
}

type resourceRecordCNAMEModel struct {
	Zone      types.String `tfsdk:"zone"`
	Name      types.String `tfsdk:"name"`
	Canonical types.String `tfsdk:"canonical"`
}

func (r *resourceRecordCNAME) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_record_cname"
}

func (r *resourceRecordCNAME) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *resourceRecordCNAME) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: md.ResourceRecordCNAME,
		Attributes: map[string]schema.Attribute{
			"zone": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: md.ResourceRecordZone,
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
			"canonical": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: md.ResourceRecordCNAMECanonical,
				Validators: []validator.String{
					validators.StringWithoutWhitespace(),
					validators.CanonicalName(),
				},
			},
		},
	}
}

func (r *resourceRecordCNAME) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan resourceRecordCNAMEModel

	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone      = plan.Zone.ValueString()
		name      = plan.Name.ValueString()
		canonical = plan.Canonical.ValueString()
	)

	err := r.client.CreateRecordCNAME(zone, name, canonical)
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating reg.ru CNAME record",
			fmt.Sprintf(
				"Could not create reg.ru CNAME record %s -> %s in zone %s: %s",
				name,
				canonical,
				zone,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, plan)...)
}

func (r *resourceRecordCNAME) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state resourceRecordCNAMEModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone      = state.Zone.ValueString()
		name      = state.Name.ValueString()
		canonical = state.Canonical.ValueString()
	)

	_, err := r.client.GetRecordCNAME(zone, name, canonical)
	if err != nil {
		var notFoundError regru.RegruRecordNotFoundError
		if errors.As(err, &notFoundError) {
			res.State.RemoveResource(ctx)
			return
		}

		res.Diagnostics.AddError(
			"Error reading reg.ru CNAME record",
			fmt.Sprintf(
				"Could not read reg.ru CNAME record %s -> %s in zone %s: %s",
				name,
				canonical,
				zone,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}

func (r *resourceRecordCNAME) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	res.Diagnostics.AddError(
		"Error updating reg.ru CNAME record",
		"Terraform called Update for a resource that only supports replacement. Please report this issue to the provider developers.",
	)
}

func (r *resourceRecordCNAME) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state resourceRecordCNAMEModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone      = state.Zone.ValueString()
		name      = state.Name.ValueString()
		canonical = state.Canonical.ValueString()
	)

	err := r.client.DeleteRecordCNAME(zone, name, canonical)
	if err != nil {
		res.Diagnostics.AddError(
			"Error deleting reg.ru CNAME record",
			fmt.Sprintf(
				"Could not delete reg.ru CNAME record %s -> %s in zone %s: %s",
				name,
				canonical,
				zone,
				err,
			),
		)
	}
}

func (r *resourceRecordCNAME) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {

	d := strings.SplitN(req.ID, ",", 3)

	if len(d) < 3 {
		res.Diagnostics.AddError(
			"Error parsing import ID for CNAME record",
			fmt.Sprintf("Provided ID %s does not match expected format 'zone,name,canonical'", req.ID),
		)
		return
	}

	var (
		zone      = d[0]
		name      = d[1]
		canonical = d[2]
	)

	_, err := r.client.GetRecordCNAME(zone, name, canonical)
	if err != nil {
		res.Diagnostics.AddError(
			"Error reading reg.ru CNAME record",
			fmt.Sprintf(
				"Could not read reg.ru CNAME record %s -> %s in zone %s: %s",
				name,
				canonical,
				zone,
				err,
			),
		)
		return
	}

	state := resourceRecordCNAMEModel{
		Zone:      types.StringValue(zone),
		Name:      types.StringValue(name),
		Canonical: types.StringValue(canonical),
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}
