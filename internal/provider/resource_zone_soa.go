package provider

import (
	"context"
	"fmt"
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
	_ resource.Resource                = &resourceZoneSoa{}
	_ resource.ResourceWithConfigure   = &resourceZoneSoa{}
	_ resource.ResourceWithImportState = &resourceZoneSoa{}
)

func NewResourceZoneSoa() resource.Resource {
	return &resourceZoneSoa{}
}

type resourceZoneSoa struct {
	client *regru.Client
}

type resourceZoneSoaModel struct {
	Zone       types.String `tfsdk:"zone"`
	Ttl        types.String `tfsdk:"ttl"`
	MinimumTtl types.String `tfsdk:"minimum_ttl"`
}

func (r *resourceZoneSoa) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_zone_soa"
}

func (r *resourceZoneSoa) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *resourceZoneSoa) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: md.ResourceZoneSoa,
		Attributes: map[string]schema.Attribute{
			"zone": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Domain passed to reg.ru as `domain_name`. Changing it replaces this Terraform resource.",
				Validators: []validator.String{
					validators.StringWithoutWhitespace(),
				},
			},
			"ttl": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: md.SoaTtl,
				Validators: []validator.String{
					validators.StringWithoutWhitespace(),
				},
			},
			"minimum_ttl": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: md.SoaMinimumTtl,
				Validators: []validator.String{
					validators.StringWithoutWhitespace(),
				},
			},
		},
	}
}

func (r *resourceZoneSoa) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan resourceZoneSoaModel

	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone       = plan.Zone.ValueString()
		ttl        = plan.Ttl.ValueString()
		minimumTtl = plan.MinimumTtl.ValueString()
	)

	err := r.client.UpdateZoneSoa(zone, ttl, minimumTtl)
	if err != nil {
		res.Diagnostics.AddError(
			"Error updating reg.ru zone SOA",
			fmt.Sprintf(
				"Could not update reg.ru zone %s SOA with values ttl=%s and minimum_ttl=%s: %s",
				zone,
				ttl,
				minimumTtl,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, plan)...)
}

func (r *resourceZoneSoa) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state resourceZoneSoaModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	zone, err := r.client.GetZoneInfo(state.Zone.ValueString())
	if err != nil {
		res.Diagnostics.AddError(
			"Unable to read reg.ru zone",
			fmt.Sprintf("Could not read zone %q: %s", state.Zone.ValueString(), err),
		)
		return
	}

	state.Ttl = types.StringValue(zone.Soa.Ttl)
	state.MinimumTtl = types.StringValue(zone.Soa.MinimumTtl)

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}

func (r *resourceZoneSoa) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var plan resourceZoneSoaModel

	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone       = plan.Zone.ValueString()
		ttl        = plan.Ttl.ValueString()
		minimumTtl = plan.MinimumTtl.ValueString()
	)

	err := r.client.UpdateZoneSoa(zone, ttl, minimumTtl)
	if err != nil {
		res.Diagnostics.AddError(
			"Error updating reg.ru zone SOA",
			fmt.Sprintf(
				"Could not update reg.ru zone %s SOA with values ttl=%s and minimum_ttl=%s: %s",
				zone,
				ttl,
				minimumTtl,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, plan)...)
}

func (r *resourceZoneSoa) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
}

func (r *resourceZoneSoa) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	zone, err := r.client.GetZoneInfo(req.ID)
	if err != nil {
		res.Diagnostics.AddError(
			"Unable to read reg.ru zone",
			fmt.Sprintf("Could not read zone %q: %s", req.ID, err),
		)
	}

	state := resourceZoneSoaModel{
		Zone:       types.StringValue(zone.Name),
		Ttl:        types.StringValue(zone.Soa.Ttl),
		MinimumTtl: types.StringValue(zone.Soa.MinimumTtl),
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}
