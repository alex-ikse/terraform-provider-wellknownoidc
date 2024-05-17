// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DocumentDataSource{}

func NewDocumentDataSource() datasource.DataSource {
	return &DocumentDataSource{}
}

// DocumentDataSource defines the data source implementation.
type DocumentDataSource struct {
}

// DocumentDataSourceModel describes the data source data model.
type DocumentDataSourceModel struct {
	ID                               types.String `tfsdk:"id"`
	DiscoveryURL                     types.String `tfsdk:"discovery_url"`
	Issuer                           types.String `tfsdk:"issuer"`
	JWKSURI                          types.String `tfsdk:"jwks_uri"`
	ResponseTypesSupported           types.List   `tfsdk:"response_types_supported"`
	SubjectTypesSupported            types.List   `tfsdk:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported types.List   `tfsdk:"id_token_signing_alg_values_supported"`
	ClaimsSupported                  types.List   `tfsdk:"claims_supported"`
	GrantTypes                       types.List   `tfsdk:"grant_types"`
}

func (d *DocumentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_document"
}

func (d *DocumentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Document data source",

		Attributes: map[string]schema.Attribute{
			"discovery_url": schema.StringAttribute{
				MarkdownDescription: "OpenID configuration discovery URL",
				Required:            true,
			},
			"issuer": schema.StringAttribute{
				MarkdownDescription: "URL using the https scheme with no query or fragment component that the OP asserts as its Issuer Identifier",
				Computed:            true,
			},
			"jwks_uri": schema.StringAttribute{
				MarkdownDescription: "URL of the OP's JSON Web Key Set [JWK] document. This contains the signing key(s) the RP uses to validate signatures from the OP",
				Computed:            true,
			},
			"response_types_supported": schema.ListAttribute{
				MarkdownDescription: "List of the OAuth 2.0 response_type values that this OP supports",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"subject_types_supported": schema.ListAttribute{
				MarkdownDescription: "List of the OAuth 2.0 response_mode values that this OP support",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"id_token_signing_alg_values_supported": schema.ListAttribute{
				MarkdownDescription: "List of the JWS signing algorithms (alg values) supported by the OP for the ID Token to encode the Claims in a JWT",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"claims_supported": schema.ListAttribute{
				MarkdownDescription: "List of the Claim Names of the Claims that the OpenID Provider MAY be able to supply values for",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"grant_types": schema.ListAttribute{
				MarkdownDescription: "List of the OAuth 2.0 Grant Type values that this OP supports",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique ID of this configuration",
				Computed:            true,
			},
		},
	}
}

func (d *DocumentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (d *DocumentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DocumentDataSourceModel
	var diagnostics diag.Diagnostics
	var discoveryDocument *DiscoveryDocument

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	r, err := http.Get(fmt.Sprintf("%s/.well-known/openid-configuration", data.DiscoveryURL.ValueString()))
	if err != nil {
		tflog.Error(ctx, "Unable to get the discovery document")
		resp.Diagnostics.AddError("Unable to get the discovery document", "Provided URL is not valid or not available.")
		return
	}
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		tflog.Error(ctx, "Unable to read the discovery document")
		resp.Diagnostics.AddError("Unable to read the discovery document", "Document found is not readable.")
		return
	}

	err = json.Unmarshal(b, &discoveryDocument)
	if err != nil {
		tflog.Error(ctx, "Unable to decode the discovery document")
		resp.Diagnostics.AddError("Unable to decode the discovery document", "Document found is not in the right JSON format.")
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.ID = data.DiscoveryURL
	data.Issuer = types.StringValue(discoveryDocument.Issuer)
	data.JWKSURI = types.StringValue(discoveryDocument.JWKSURI)

	data.ClaimsSupported, diagnostics = types.ListValueFrom(ctx, types.StringType, discoveryDocument.ClaimsSupported)
	resp.Diagnostics.Append(diagnostics...)
	data.ResponseTypesSupported, diagnostics = types.ListValueFrom(ctx, types.StringType, discoveryDocument.ResponseTypesSupported)
	resp.Diagnostics.Append(diagnostics...)
	data.SubjectTypesSupported, diagnostics = types.ListValueFrom(ctx, types.StringType, discoveryDocument.SubjectTypesSupported)
	resp.Diagnostics.Append(diagnostics...)
	data.IDTokenSigningAlgValuesSupported, diagnostics = types.ListValueFrom(ctx, types.StringType, discoveryDocument.IDTokenSigningAlgValuesSupported)
	resp.Diagnostics.Append(diagnostics...)
	data.GrantTypes, diagnostics = types.ListValueFrom(ctx, types.StringType, discoveryDocument.GrantTypes)
	resp.Diagnostics.Append(diagnostics...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
