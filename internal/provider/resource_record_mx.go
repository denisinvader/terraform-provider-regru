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
	_ resource.Resource                = &resourceRecordMX{}
	_ resource.ResourceWithConfigure   = &resourceRecordMX{}
	_ resource.ResourceWithImportState = &resourceRecordMX{}
)

func NewResourceRecordMX() resource.Resource {
	return &resourceRecordMX{}
}

type resourceRecordMX struct {
	client *regru.Client
}

type resourceRecordMXModel struct {
	Zone       types.String `tfsdk:"zone"`
	Name       types.String `tfsdk:"name"`
	MailServer types.String `tfsdk:"mail_server"`
	Priority   types.Int32  `tfsdk:"priority"`
}

func (r *resourceRecordMX) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_record_mx"
}

func (r *resourceRecordMX) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *resourceRecordMX) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: md.ResourceRecordMX,
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
				MarkdownDescription: "Subdomain name or `@` for the zone apex.",
				Validators: []validator.String{
					validators.StringWithoutWhitespace(),
				},
			},
			"mail_server": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: md.ResourceRecordMXMailServer,
				Validators: []validator.String{
					validators.StringWithoutWhitespace(),
					validators.CanonicalName(),
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

func (r *resourceRecordMX) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan resourceRecordMXModel

	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone       = plan.Zone.ValueString()
		name       = plan.Name.ValueString()
		mailServer = plan.MailServer.ValueString()
		priority   = int(plan.Priority.ValueInt32())
	)

	err := r.client.CreateRecordMX(zone, name, mailServer, priority)
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating reg.ru MX record",
			fmt.Sprintf(
				"Could not create reg.ru MX record %s -> %s (%d) in zone %s: %s",
				name,
				mailServer,
				priority,
				zone,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, plan)...)
}

func (r *resourceRecordMX) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state resourceRecordMXModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone       = state.Zone.ValueString()
		name       = state.Name.ValueString()
		mailServer = state.MailServer.ValueString()
		priority   = int(state.Priority.ValueInt32())
	)

	_, err := r.client.GetRecordMX(zone, name, mailServer, priority)
	if err != nil {
		var notFoundError regru.RegruRecordNotFoundError
		if errors.As(err, &notFoundError) {
			res.State.RemoveResource(ctx)
			return
		}

		res.Diagnostics.AddError(
			"Error reading reg.ru MX record",
			fmt.Sprintf(
				"Could not read reg.ru MX record %s -> %s (%d) in zone %s: %s",
				name,
				mailServer,
				priority,
				zone,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}

func (r *resourceRecordMX) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	res.Diagnostics.AddError(
		"Error updating reg.ru MX record",
		"Terraform called Update for a resource that only supports replacement. Please report this issue to the provider developers.",
	)
}

func (r *resourceRecordMX) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state resourceRecordMXModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone       = state.Zone.ValueString()
		name       = state.Name.ValueString()
		mailServer = state.MailServer.ValueString()
		priority   = int(state.Priority.ValueInt32())
	)

	err := r.client.DeleteRecordMX(zone, name, mailServer, priority)
	if err != nil {
		res.Diagnostics.AddError(
			"Error deleting reg.ru MX record",
			fmt.Sprintf(
				"Could not delete reg.ru MX record %s -> %s (%d) in zone %s: %s",
				name,
				mailServer,
				priority,
				zone,
				err,
			),
		)
	}
}

func (r *resourceRecordMX) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {

	d := strings.SplitN(req.ID, ",", 4)

	if len(d) < 4 {
		res.Diagnostics.AddError(
			"Error parsing import ID for MX record",
			fmt.Sprintf("Provided ID %s does not match expected format 'zone,name,priority,mail_server'", req.ID),
		)
		return
	}

	var (
		zone       = d[0]
		name       = d[1]
		mailServer = d[3]
	)

	priority, parseErr := strconv.Atoi(d[2])
	if parseErr != nil {
		res.Diagnostics.AddError(
			"Error parsing import ID for MX record",
			fmt.Sprintf(
				"Failed to parse priority %q as int: %s",
				d[2],
				parseErr,
			),
		)
		return
	}

	_, err := r.client.GetRecordMX(zone, name, mailServer, int(priority))
	if err != nil {
		res.Diagnostics.AddError(
			"Error reading reg.ru MX record",
			fmt.Sprintf(
				"Could not read reg.ru MX record %s -> %s (%d) in zone %s: %s",
				name,
				mailServer,
				zone,
				err,
			),
		)
		return
	}

	state := resourceRecordMXModel{
		Zone:       types.StringValue(zone),
		Name:       types.StringValue(name),
		MailServer: types.StringValue(mailServer),
		Priority:   types.Int32Value(int32(priority)),
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}
