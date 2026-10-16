package provider

import (
	"context"
	"fmt"
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
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *regru.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *dataSourceZone) Schema(_ context.Context, req datasource.SchemaRequest, res *datasource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				// FIXME: write description
				Description: "",
			},
			"service_type": schema.StringAttribute{
				Computed: true,
				// FIXME: write description
				Description: "",
			},
			"service_id": schema.StringAttribute{
				Computed: true,
				// FIXME: write description
				Description: "",
			},
			"soa": schema.SingleNestedAttribute{
				Computed: true,
				// FIXME: write description
				Description: "",
				Attributes: map[string]schema.Attribute{
					"ttl": schema.StringAttribute{
						Computed: true,
						// FIXME: write description
						Description: "",
					},
					"minimum_ttl": schema.StringAttribute{
						Computed: true,
						// FIXME: write description
						Description: "",
					},
				},
			},
			"records": schema.ListNestedAttribute{
				Computed: true,
				// FIXME: write description
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed: true,
							// FIXME: write description
							Description: "",
						},
						"name": schema.StringAttribute{
							Computed: true,
							// FIXME: write description
							Description: "",
						},
						"priority": schema.Int32Attribute{
							Computed: true,
							// FIXME: write description
							Description: "",
						},
						"content": schema.StringAttribute{
							Computed: true,
							// FIXME: write description
							Description: "",
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
		res.Diagnostics.AddError("Unable to read reg.ru zone", err.Error())
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
