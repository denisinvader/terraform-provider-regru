package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
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
	_ resource.Resource                = &resourceRecordTXT{}
	_ resource.ResourceWithConfigure   = &resourceRecordTXT{}
	_ resource.ResourceWithImportState = &resourceRecordTXT{}
)

func NewResourceRecordTXT() resource.Resource {
	return &resourceRecordTXT{}
}

type resourceRecordTXT struct {
	client *regru.Client
}

type resourceRecordTXTModel struct {
	Zone types.String `tfsdk:"zone"`
	Name types.String `tfsdk:"name"`
	Text types.String `tfsdk:"text"`
}

func (r *resourceRecordTXT) Metadata(ctx context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_record_txt"
}

func (r *resourceRecordTXT) Configure(ctx context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
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

func (r *resourceRecordTXT) Schema(ctx context.Context, req resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: md.ResourceRecordTXT,
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
				MarkdownDescription: "Subdomain name.",
				Validators: []validator.String{
					validators.StringWithoutWhitespace(),
				},
			},
			"text": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				MarkdownDescription: md.ResourceRecordTXTText,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[^\s]+.*[^\s]+$`),
						"must not start or end with whitespace",
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[!"#$%&'()*+,\-./:;<=>?@0-9A-Za-z\[\\\]^_\x60{|}~ ]+$`),
						"must contain only ASCII characters",
					),
					stringvalidator.LengthAtMost(512),
				},
			},
		},
	}
}

func (r *resourceRecordTXT) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan resourceRecordTXTModel

	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone = plan.Zone.ValueString()
		name = plan.Name.ValueString()
		text = plan.Text.ValueString()
	)

	err := r.client.CreateRecordTXT(zone, name, text)
	if err != nil {
		res.Diagnostics.AddError(
			"Error creating reg.ru TXT record",
			fmt.Sprintf(
				"Could not create reg.ru TXT record %s -> %s in zone %s: %s",
				name,
				text,
				zone,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, plan)...)
}

func (r *resourceRecordTXT) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state resourceRecordTXTModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone = state.Zone.ValueString()
		name = state.Name.ValueString()
		text = state.Text.ValueString()
	)

	_, err := r.client.GetRecordTXT(zone, name, text)
	if err != nil {
		var notFoundError regru.RegruRecordNotFoundError
		if errors.As(err, &notFoundError) {
			res.State.RemoveResource(ctx)
			return
		}

		res.Diagnostics.AddError(
			"Error reading reg.ru TXT record",
			fmt.Sprintf(
				"Could not read reg.ru TXT record %s -> %s in zone %s: %s",
				name,
				text,
				zone,
				err,
			),
		)
		return
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}

func (r *resourceRecordTXT) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	res.Diagnostics.AddError(
		"Error updating reg.ru TXT record",
		"Terraform called Update for a resource that only supports replacement. Please report this issue to the provider developers.",
	)
}

func (r *resourceRecordTXT) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state resourceRecordTXTModel

	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	var (
		zone = state.Zone.ValueString()
		name = state.Name.ValueString()
		text = state.Text.ValueString()
	)

	err := r.client.DeleteRecordTXT(zone, name, text)
	if err != nil {
		res.Diagnostics.AddError(
			"Error deleting reg.ru TXT record",
			fmt.Sprintf(
				"Could not delete reg.ru TXT record %s -> %s in zone %s: %s",
				name,
				text,
				zone,
				err,
			),
		)
	}
}

func (r *resourceRecordTXT) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {

	d := strings.SplitN(req.ID, ",", 3)

	if len(d) < 3 {
		res.Diagnostics.AddError(
			"Error parsing import ID for TXT record",
			fmt.Sprintf("Provided ID %s does not match expected format 'zone,name,text'", req.ID),
		)
		return
	}

	var (
		zone = d[0]
		name = d[1]
		text = d[2]
	)

	_, err := r.client.GetRecordTXT(zone, name, text)
	if err != nil {
		res.Diagnostics.AddError(
			"Error reading reg.ru TXT record",
			fmt.Sprintf(
				"Could not read reg.ru TXT record %s -> %s in zone %s: %s",
				name,
				text,
				zone,
				err,
			),
		)
		return
	}

	state := resourceRecordTXTModel{
		Zone: types.StringValue(zone),
		Name: types.StringValue(name),
		Text: types.StringValue(text),
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}
