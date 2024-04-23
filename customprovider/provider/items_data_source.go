package provider

import (
	"context"
	"terraform-provider-customprovider/client"
	"terraform-provider-customprovider/model"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &itemsDataSource{} // check interface compliance

func NewItemsDataSource() datasource.DataSource {
	return &itemsDataSource{}
}

type itemsDataSource struct{}

func (d *itemsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_items"
}

func (d *itemsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"items": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *itemsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state model.ItemsTFSDK

	// fetch items from API
	items, err := client.ReadItems()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read CustomAPI Items",
			err.Error(),
		)
		return
	}

	// map response body to model
	for _, item := range items.Items {
		itemState := model.ItemTFSDK{
			ID:   types.StringValue(item.ID),
			Name: types.StringValue(item.Name),
		}
		state.Items = append(state.Items, itemState)
	}

	// set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
