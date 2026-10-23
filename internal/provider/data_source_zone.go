package provider

import (
	"context"
	"fmt"
	"terraform-provider-regru/internal/md"
	"terraform-provider-regru/internal/regru"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &dataSourceZone{}
	_ datasource.DataSourceWithConfigure = &dataSourceZone{}
)

func NewDataSourceZone() datasource.DataSource {
	return &dataSourceZone{}
}

type dataSourceZone struct {
	client *regru.Client
}

type dataSourceZoneModel struct {
	Name        types.String                `tfsdk:"name"`
	ServiceType types.String                `tfsdk:"service_type"`
	ServiceID   types.String                `tfsdk:"service_id"`
	Soa         *dataSourceZoneSoaModel     `tfsdk:"soa"`
	Records     []dataSourceZoneRecordModel `tfsdk:"records"`
}

type dataSourceZoneSoaModel struct {
	Ttl        types.String `tfsdk:"ttl"`
	MinimumTtl types.String `tfsdk:"minimum_ttl"`
}

type dataSourceZoneRecordModel struct {
	Type     types.String `tfsdk:"type"`
	Name     types.String `tfsdk:"name"`
	Priority types.Int32  `tfsdk:"priority"`
	Content  types.String `tfsdk:"content"`
}

func (d *dataSourceZone) Metadata(_ context.Context, req datasource.MetadataRequest, res *datasource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_zone"
}

func (d *dataSourceZone) Configure(_ context.Context, req datasource.ConfigureRequest, res *datasource.ConfigureResponse) {
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

	d.client = client
}

func (d *dataSourceZone) Schema(_ context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		MarkdownDescription: "Reads the zone's SOA settings, reg.ru service details, and current DNS records.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Zone name, e.g. `example.com`.",
			},
			"service_type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Internal reg.ru service type. This is normally `domain`.",
			},
			"service_id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Numeric identifier of the internal reg.ru service associated with the domain.",
			},
			"soa": schema.SingleNestedAttribute{
				Computed:            true,
				MarkdownDescription: "Start of Authority (SOA) record settings.",
				Attributes: map[string]schema.Attribute{
					"ttl": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: md.SoaTtl,
					},
					"minimum_ttl": schema.StringAttribute{
						Computed:            true,
						MarkdownDescription: md.SoaMinimumTtl,
					},
				},
			},
			"records": schema.ListNestedAttribute{
				Computed:            true,
				MarkdownDescription: "List of DNS records currently configured for the zone.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Record type, such as `A`, `AAAA`, `CAA`, `CNAME`, `HTTPS`, `MX`, `NS`, `SRV`, or `TXT`.",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Subdomain name, `*` for a wildcard, or `@` for the zone apex.",
						},
						"priority": schema.Int32Attribute{
							Computed:            true,
							MarkdownDescription: "Record priority.",
						},
						"content": schema.StringAttribute{
							Computed:            true,
							MarkdownDescription: "Record content: an IPv4 address for `A` records, an IPv6 address for `AAAA` records, a target domain for `CNAME` records, and so on.",
						},
					},
				},
			},
		},
	}
}

func (d *dataSourceZone) Read(ctx context.Context, req datasource.ReadRequest, res *datasource.ReadResponse) {
	var state dataSourceZoneModel

	res.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	zone, err := d.client.GetZoneInfo(state.Name.ValueString())
	if err != nil {
		res.Diagnostics.AddError(
			"Unable to read reg.ru zone",
			fmt.Sprintf("Could not read zone %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	state.Name = types.StringValue(zone.Name)
	state.ServiceID = types.StringValue(zone.ServiceID)
	state.ServiceType = types.StringValue(zone.ServiceType)
	state.Soa = &dataSourceZoneSoaModel{
		Ttl:        types.StringValue(zone.Soa.Ttl),
		MinimumTtl: types.StringValue(zone.Soa.MinimumTtl),
	}
	for _, record := range zone.Records {
		state.Records = append(state.Records, dataSourceZoneRecordModel{
			Type:     types.StringValue(record.Type),
			Name:     types.StringValue(record.Name),
			Priority: types.Int32Value(int32(record.Priority)),
			Content:  types.StringValue(record.Content),
		})
	}

	res.Diagnostics.Append(res.State.Set(ctx, state)...)
}
