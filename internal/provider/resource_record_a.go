package provider

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"terraform-provider-regru/internal/md"
	"terraform-provider-regru/internal/regru"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &resourceRecordA{}
	_ resource.ResourceWithConfigure   = &resourceRecordA{}
	_ resource.ResourceWithImportState = &resourceRecordA{}
)

func NewResourceRecordA() resource.Resource {
	return &resourceRecordA{}
}

type resourceRecordA struct {
	client *regru.Client
}

type resourceRecordAModel struct {
	Zone types.String `tfsdk:"zone"`
	Name types.String `tfsdk:"name"`
	IP   types.String `tfsdk:"ip"`
}

func (r *resourceRecordA) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_record_a"
}

func (r *resourceRecordA) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *resourceRecordA) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: md.ResourceRecordA,
		Attributes: map[string]schema.Attribute{
			"zone": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: md.ResourceRecordZone,
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Subdomain name, `*` for a wildcard, or `@` for the zone apex.",
			},
			"ip": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Target IPv4 value.",
			},
		},
	}
}

func (r *resourceRecordA) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan resourceRecordAModel

	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone = plan.Zone.ValueString()
		name = plan.Name.ValueString()
		ip   = plan.IP.ValueString()
	)

	err := r.client.CreateRecordA(zone, name, ip)
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating reg.ru A record",
			fmt.Sprintf(
				"Could not create reg.ru A record %s -> %s in zone %s: %s",
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

func (r *resourceRecordA) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state resourceRecordAModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone = state.Zone.ValueString()
		name = state.Name.ValueString()
		ip   = state.IP.ValueString()
	)

	_, err := r.client.GetRecordA(zone, name, ip)
	if err != nil {
		var notFoundError regru.RegruRecordNotFoundError
		if errors.As(err, &notFoundError) {
			res.State.RemoveResource(ctx)
			return
		}

		res.Diagnostics.AddError(
			"Error reading reg.ru A record",
			fmt.Sprintf(
				"Could not read reg.ru A record %s -> %s in zone %s: %s",
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

func (r *resourceRecordA) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	res.Diagnostics.AddError(
		"Error updating reg.ru A record",
		"Terraform called Update for a resource that only supports replacement. Please report this issue to the provider developers.",
	)
}

func (r *resourceRecordA) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state resourceRecordAModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone = state.Zone.ValueString()
		name = state.Name.ValueString()
		ip   = state.IP.ValueString()
	)

	err := r.client.DeleteRecordA(zone, name, ip)
	if err != nil {
		res.Diagnostics.AddError(
			"Error deleting reg.ru A record",
			fmt.Sprintf(
				"Could not delete reg.ru A record %s -> %s in zone %s: %s",
				name,
				ip,
				zone,
				err,
			),
		)
	}
}

func (r *resourceRecordA) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {

	d := strings.SplitN(req.ID, ",", 3)

	if len(d) < 3 {
		res.Diagnostics.AddError(
			"Error parsing import ID for A record",
			fmt.Sprintf("Provided ID %s does not match expected format 'zone,name,ip'", req.ID),
		)
		return
	}

	var (
		zone = d[0]
		name = d[1]
		ip   = d[2]
	)

	_, err := r.client.GetRecordA(zone, name, ip)
	if err != nil {
		res.Diagnostics.AddError(
			"Error reading reg.ru A record",
			fmt.Sprintf(
				"Could not read reg.ru A record %s -> %s in zone %s: %s",
				name,
				ip,
				zone,
				err,
			),
		)
		return
	}

	state := resourceRecordAModel{
		Zone: types.StringValue(zone),
		Name: types.StringValue(name),
		IP:   types.StringValue(ip),
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}
