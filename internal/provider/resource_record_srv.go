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
	_ resource.Resource                = &resourceRecordSRV{}
	_ resource.ResourceWithConfigure   = &resourceRecordSRV{}
	_ resource.ResourceWithImportState = &resourceRecordSRV{}
)

func NewResourceRecordSRV() resource.Resource {
	return &resourceRecordSRV{}
}

type resourceRecordSRV struct {
	client *regru.Client
}

type resourceRecordSRVModel struct {
	Zone     types.String `tfsdk:"zone"`
	Service  types.String `tfsdk:"service"`
	Target   types.String `tfsdk:"target"`
	Port     types.Int32  `tfsdk:"port"`
	Weight   types.Int32  `tfsdk:"weight"`
	Priority types.Int32  `tfsdk:"priority"`
}

func (r *resourceRecordSRV) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_record_srv"
}

func (r *resourceRecordSRV) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *resourceRecordSRV) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: md.ResourceRecordSRV,
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
			"service": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Passed to reg.ru as `service` without normalization. Whitespace is rejected; for example, `_sip._udp`.",
				Validators: []validator.String{
					validators.StringWithoutWhitespace(),
				},
			},
			"target": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: md.ResourceRecordSRVTarget,
				Validators: []validator.String{
					validators.StringWithoutWhitespace(),
					validators.CanonicalName(),
				},
			},
			"port": schema.Int32Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Non-negative integer passed to reg.ru as `port`.",
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			"weight": schema.Int32Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				MarkdownDescription: "Non-negative integer passed to reg.ru as `weight`. The provider sends `0` when this attribute is omitted.",
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
			"priority": schema.Int32Attribute{
				Optional: true,
				PlanModifiers: []planmodifier.Int32{
					int32planmodifier.RequiresReplace(),
				},
				MarkdownDescription: md.RecordPriority,
				Validators: []validator.Int32{
					int32validator.AtLeast(0),
				},
			},
		},
	}
}

func (r *resourceRecordSRV) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan resourceRecordSRVModel

	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone     = plan.Zone.ValueString()
		service  = plan.Service.ValueString()
		target   = plan.Target.ValueString()
		port     = int(plan.Port.ValueInt32())
		weight   = int(plan.Weight.ValueInt32())
		priority = int(plan.Priority.ValueInt32())
	)

	err := r.client.CreateRecordSRV(zone, service, target, port, weight, priority)
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating reg.ru SRV record",
			fmt.Sprintf(
				"Could not create reg.ru SRV record %s -> %d %d %s (%d) in zone %s: %s",
				service,
				weight,
				port,
				target,
				priority,
				zone,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, plan)...)
}

func (r *resourceRecordSRV) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state resourceRecordSRVModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone     = state.Zone.ValueString()
		service  = state.Service.ValueString()
		target   = state.Target.ValueString()
		port     = int(state.Port.ValueInt32())
		weight   = int(state.Weight.ValueInt32())
		priority = int(state.Priority.ValueInt32())
	)

	_, err := r.client.GetRecordSRV(zone, service, target, port, weight, priority)
	if err != nil {
		var notFoundError regru.RegruRecordNotFoundError
		if errors.As(err, &notFoundError) {
			res.State.RemoveResource(ctx)
			return
		}

		res.Diagnostics.AddError(
			"Error reading reg.ru SRV record",
			fmt.Sprintf(
				"Could not read reg.ru SRV record %s -> %d %d %s (%d) in zone %s: %s",
				service,
				weight,
				port,
				target,
				priority,
				zone,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}

func (r *resourceRecordSRV) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	res.Diagnostics.AddError(
		"Error updating reg.ru SRV record",
		"Terraform called Update for a resource that only supports replacement. Please report this issue to the provider developers.",
	)
}

func (r *resourceRecordSRV) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state resourceRecordSRVModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone     = state.Zone.ValueString()
		service  = state.Service.ValueString()
		target   = state.Target.ValueString()
		port     = int(state.Port.ValueInt32())
		weight   = int(state.Weight.ValueInt32())
		priority = int(state.Priority.ValueInt32())
	)

	err := r.client.DeleteRecordSRV(zone, service, target, port, weight, priority)
	if err != nil {
		res.Diagnostics.AddError(
			"Error deleting reg.ru SRV record",
			fmt.Sprintf(
				"Could not delete reg.ru SRV record %s -> %d %d %s (%d) in zone %s: %s",
				service,
				weight,
				port,
				target,
				priority,
				zone,
				err,
			),
		)
	}
}

func (r *resourceRecordSRV) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {

	d := strings.SplitN(req.ID, ",", 6)

	if len(d) < 6 {
		res.Diagnostics.AddError(
			"Error parsing import ID for SRV record",
			fmt.Sprintf("Provided ID %s does not match expected format 'zone,service,priority,weight,port,target'", req.ID),
		)
		return
	}

	var (
		zone    = d[0]
		service = d[1]
		target  = d[5]
	)

	port, portParseErr := strconv.Atoi(d[4])
	if portParseErr != nil {
		res.Diagnostics.AddError(
			"Error parsing import ID for SRV record",
			fmt.Sprintf("Failed to parse port %q as int: %s", d[4], portParseErr),
		)
		return
	}

	weight, weightParseErr := strconv.Atoi(d[3])
	if weightParseErr != nil {
		res.Diagnostics.AddError(
			"Error parsing import ID for SRV record",
			fmt.Sprintf("Failed to parse weight %q as int: %s", d[3], weightParseErr),
		)
		return
	}

	priority, priorityParseErr := strconv.Atoi(d[2])
	if priorityParseErr != nil {
		res.Diagnostics.AddError(
			"Error parsing import ID for SRV record",
			fmt.Sprintf("Failed to parse priority %q as int: %s", d[2], priorityParseErr),
		)
		return
	}

	_, err := r.client.GetRecordSRV(zone, service, target, port, weight, priority)
	if err != nil {
		res.Diagnostics.AddError(
			"Error reading reg.ru SRV record",
			fmt.Sprintf(
				"Could not read reg.ru SRV record %s -> %d %d %s (%d) in zone %s: %s",
				service,
				weight,
				port,
				target,
				priority,
				zone,
				err,
			),
		)
		return
	}

	state := resourceRecordSRVModel{
		Zone:     types.StringValue(zone),
		Service:  types.StringValue(service),
		Target:   types.StringValue(target),
		Port:     types.Int32Value(int32(port)),
		Weight:   types.Int32Value(int32(weight)),
		Priority: types.Int32Value(int32(priority)),
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}
