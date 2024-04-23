package provider

import (
	"context"
	"terraform-provider-customprovider/client"
	"terraform-provider-customprovider/model"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &itemResource{} // check interface compliance

func NewItemResource() resource.Resource {
	return &itemResource{}
}

type itemResource struct{}

func (r *itemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_item"
}

func (r *itemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *itemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// retrieve values from plan
	var plan model.ItemTFSDK
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// create new item
	item, err := client.CreateItem(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Custom API item",
			"Could not item, unexpected error: "+err.Error(),
		)
		return
	}

	// map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(item.ID)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *itemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// get current state
	var state model.ItemTFSDK
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// get refreshed order value from CustomAPI
	item, err := client.ReadItem(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading CustomAPI Item",
			"Unable to Read CustomAPI Item: "+err.Error(),
		)
		return
	}

	// overwrite items with refreshed state
	if item == nil {
		// resource does not exist anymore
		resp.State.RemoveResource(ctx)
		return
	}
	state = model.ItemTFSDK{
		ID:   types.StringValue(item.ID),
		Name: types.StringValue(item.Name),
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *itemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// retrieve values from plan
	var plan model.ItemTFSDK
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve values from state
	var state model.ItemTFSDK
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// update existing item
	item, err := client.UpdateItem(state.ID.ValueString(), plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating CustomAPI item",
			"Could not update item, unexpected error: "+err.Error(),
		)
		return
	}

	// remark: fetching updated item can be skipped as put returns item

	// update resource state with updated items and timestamp
	plan.ID = types.StringValue(item.ID)
	plan.Name = types.StringValue(item.Name)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *itemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// retrieve values from state
	var state model.ItemTFSDK
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// delete existing order
	err := client.DeleteItem(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting CustomAPI Order",
			"Could not delete item, unexpected error: "+err.Error(),
		)
		return
	}
}
