package model

import "github.com/hashicorp/terraform-plugin-framework/types"

// JSON for HTTP
type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TFSDK for Terrafrom internally
type ItemsTFSDK struct {
	Items []ItemTFSDK `tfsdk:"items"`
}

type ItemTFSDK struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}
