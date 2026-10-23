package provider

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"terraform-provider-regru/internal/md"
	"terraform-provider-regru/internal/provider/validators"
	"terraform-provider-regru/internal/regru"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &resourceRecordNS{}
	_ resource.ResourceWithConfigure   = &resourceRecordNS{}
	_ resource.ResourceWithImportState = &resourceRecordNS{}
)

func NewResourceRecordNS() resource.Resource {
	return &resourceRecordNS{}
}

type resourceRecordNS struct {
	client *regru.Client
}

type resourceRecordNSModel struct {
	Zone      types.String `tfsdk:"zone"`
	Name      types.String `tfsdk:"name"`
	DNSServer types.String `tfsdk:"dns_server"`
	Index     types.Int32  `tfsdk:"index"`
}

func (r *resourceRecordNS) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_record_ns"
}

func (r *resourceRecordNS) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *resourceRecordNS) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: md.ResourceRecordNS,
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
				MarkdownDescription: "Subdomain name or `*` for wildcard.",
				Validators: []validator.String{
					validators.StringWithoutWhitespace(),
				},
			},
			"dns_server": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: md.ResourceRecordNSDNSServer,
				Validators: []validator.String{
					validators.StringWithoutWhitespace(),
					validators.CanonicalName(),
				},
			},
			"index": schema.Int32Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				MarkdownDescription: md.ResourceRecordNSIndex,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
		},
	}
}

func (r *resourceRecordNS) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan resourceRecordNSModel

	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone      = plan.Zone.ValueString()
		name      = plan.Name.ValueString()
		dnsServer = plan.DNSServer.ValueString()
		index     = int(plan.Index.ValueInt32())
	)

	err := r.client.CreateRecordNS(zone, name, dnsServer, index)
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating reg.ru NS record",
			fmt.Sprintf(
				"Could not create reg.ru NS record %s -> %s (%d) in zone %s: %s",
				name,
				dnsServer,
				index,
				zone,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, plan)...)
}

func (r *resourceRecordNS) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state resourceRecordNSModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone      = state.Zone.ValueString()
		name      = state.Name.ValueString()
		dnsServer = state.DNSServer.ValueString()
		index     = int(state.Index.ValueInt32())
	)

	_, err := r.client.GetRecordNS(zone, name, dnsServer, index)
	if err != nil {
		var notFoundError regru.RegruRecordNotFoundError
		if errors.As(err, &notFoundError) {
			res.State.RemoveResource(ctx)
			return
		}

		res.Diagnostics.AddError(
			"Error reading reg.ru NS record",
			fmt.Sprintf(
				"Could not read reg.ru NS record %s -> %s (%d) in zone %s: %s",
				name,
				dnsServer,
				index,
				zone,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}

func (r *resourceRecordNS) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	res.Diagnostics.AddError(
		"Error updating reg.ru NS record",
		"Terraform called Update for a resource that only supports replacement. Please report this issue to the provider developers.",
	)
}

func (r *resourceRecordNS) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state resourceRecordNSModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone      = state.Zone.ValueString()
		name      = state.Name.ValueString()
		dnsServer = state.DNSServer.ValueString()
		index     = int(state.Index.ValueInt32())
	)

	err := r.client.DeleteRecordNS(zone, name, dnsServer, index)
	if err != nil {
		res.Diagnostics.AddError(
			"Error deleting reg.ru NS record",
			fmt.Sprintf(
				"Could not delete reg.ru NS record %s -> %s (%d) in zone %s: %s",
				name,
				dnsServer,
				index,
				zone,
				err,
			),
		)
	}
}

func (r *resourceRecordNS) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {

	d := strings.SplitN(req.ID, ",", 4)

	if len(d) < 4 {
		res.Diagnostics.AddError(
			"Error parsing import ID for NS record",
			fmt.Sprintf("Provided ID %s does not match expected format 'zone,name,index,dns_server'", req.ID),
		)
		return
	}

	var (
		zone      = d[0]
		name      = d[1]
		dnsServer = d[3]
	)

	index, parseErr := strconv.Atoi(d[2])
	if parseErr != nil {
		res.Diagnostics.AddError(
			"Error parsing import ID for NS record",
			fmt.Sprintf(
				"Failed to parse index %q as int: %s",
				d[2],
				parseErr,
			),
		)
		return
	}

	_, err := r.client.GetRecordNS(zone, name, dnsServer, int(index))
	if err != nil {
		res.Diagnostics.AddError(
			"Error reading reg.ru NS record",
			fmt.Sprintf(
				"Could not read reg.ru NS record %s -> %s (%d) in zone %s: %s",
				name,
				dnsServer,
				index,
				zone,
				err,
			),
		)
		return
	}

	state := resourceRecordNSModel{
		Zone:      types.StringValue(zone),
		Name:      types.StringValue(name),
		DNSServer: types.StringValue(dnsServer),
		Index:     types.Int32Value(int32(index)),
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}
