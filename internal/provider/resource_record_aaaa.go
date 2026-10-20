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
	_ resource.Resource                = &resourceRecordAAAA{}
	_ resource.ResourceWithConfigure   = &resourceRecordAAAA{}
	_ resource.ResourceWithImportState = &resourceRecordAAAA{}
)

func NewResourceRecordAAAA() resource.Resource {
	return &resourceRecordAAAA{}
}

type resourceRecordAAAA struct {
	client *regru.Client
}

type resourceRecordAAAAModel struct {
	Zone types.String `tfsdk:"zone"`
	Name types.String `tfsdk:"name"`
	IP   types.String `tfsdk:"ip"`
}

func (r *resourceRecordAAAA) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_record_aaaa"
}

func (r *resourceRecordAAAA) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *resourceRecordAAAA) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: md.ResourceRecordAAAA,
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
			"ip": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Target IPv6 address.",
				Validators: []validator.String{
					validators.StringWithoutWhitespace(),
				},
			},
		},
	}
}

func (r *resourceRecordAAAA) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan resourceRecordAAAAModel

	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone = plan.Zone.ValueString()
		name = plan.Name.ValueString()
		ip   = plan.IP.ValueString()
	)

	err := r.client.CreateRecordAAAA(zone, name, ip)
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating reg.ru AAAA record",
			fmt.Sprintf(
				"Could not create reg.ru AAAA record %s -> %s in zone %s: %s",
				name,
				ip,
				zone,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, plan)...)
}

func (r *resourceRecordAAAA) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state resourceRecordAAAAModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone = state.Zone.ValueString()
		name = state.Name.ValueString()
		ip   = state.IP.ValueString()
	)

	_, err := r.client.GetRecordAAAA(zone, name, ip)
	if err != nil {
		var notFoundError regru.RegruRecordNotFoundError
		if errors.As(err, &notFoundError) {
			res.State.RemoveResource(ctx)
			return
		}

		res.Diagnostics.AddError(
			"Error reading reg.ru AAAA record",
			fmt.Sprintf(
				"Could not read reg.ru AAAA record %s -> %s in zone %s: %s",
				name,
				ip,
				zone,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}

func (r *resourceRecordAAAA) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	res.Diagnostics.AddError(
		"Error updating reg.ru AAAA record",
		"Terraform called Update for a resource that only supports replacement. Please report this issue to the provider developers.",
	)
}

func (r *resourceRecordAAAA) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state resourceRecordAAAAModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone = state.Zone.ValueString()
		name = state.Name.ValueString()
		ip   = state.IP.ValueString()
	)

	err := r.client.DeleteRecordAAAA(zone, name, ip)
	if err != nil {
		res.Diagnostics.AddError(
			"Error deleting reg.ru AAAA record",
			fmt.Sprintf(
				"Could not delete reg.ru AAAA record %s -> %s in zone %s: %s",
				name,
				ip,
				zone,
				err,
			),
		)
	}
}

func (r *resourceRecordAAAA) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {

	d := strings.SplitN(req.ID, ",", 3)

	if len(d) < 3 {
		res.Diagnostics.AddError(
			"Error parsing import ID for AAAA record",
			fmt.Sprintf("Provided ID %s does not match expected format 'zone,name,ip'", req.ID),
		)
		return
	}

	var (
		zone = d[0]
		name = d[1]
		ip   = d[2]
	)

	_, err := r.client.GetRecordAAAA(zone, name, ip)
	if err != nil {
		res.Diagnostics.AddError(
			"Error reading reg.ru AAAA record",
			fmt.Sprintf(
				"Could not read reg.ru AAAA record %s -> %s in zone %s: %s",
				name,
				ip,
				zone,
				err,
			),
		)
		return
	}

	state := resourceRecordAAAAModel{
		Zone: types.StringValue(zone),
		Name: types.StringValue(name),
		IP:   types.StringValue(ip),
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}
