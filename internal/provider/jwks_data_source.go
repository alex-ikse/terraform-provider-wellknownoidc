// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &JWKSDataSource{}

func NewJWKSDataSource() datasource.DataSource {
	return &JWKSDataSource{}
}

// JWKSDataSource defines the data source implementation.
type JWKSDataSource struct {
}

// JWKSDataSourceModel describes the data source data model.
type JWKSDataSourceModel struct {
	ID      types.String `tfsdk:"id"`
	JWKSURI types.String `tfsdk:"jwks_uri"`
	Keys    types.List   `tfsdk:"keys"`
	PEM     types.List   `tfsdk:"pem"`
}

func (d *JWKSDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jwks"
}

func (d *JWKSDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the jwksation generator and the language server.
		MarkdownDescription: "JWKS data source",

		Attributes: map[string]schema.Attribute{
			"jwks_uri": schema.StringAttribute{
				MarkdownDescription: "OIDC JWKS URI",
				Required:            true,
			},
			"keys": schema.ListAttribute{
				MarkdownDescription: "JWKS Keys as JSON String",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"pem": schema.ListAttribute{
				MarkdownDescription: "JWKS Keys as PEM format",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique ID of this JWKS",
				Computed:            true,
			},
		},
	}
}

func (d *JWKSDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func convertJkwsToPem(key []byte) ([]byte, error) {
	set, err := jwk.ParseKey(key)
	if err != nil {
		return []byte{}, err
	}

	pem, err := jwk.EncodePEM(set)
	if err != nil {
		return []byte{}, err
	}
	return pem, nil
}

func (d *JWKSDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data JWKSDataSourceModel
	var diagnostics diag.Diagnostics
	var jwksDocument *JWKSDocument

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	r, err := http.Get(data.JWKSURI.ValueString())
	if err != nil {
		tflog.Error(ctx, "Unable to get the JWKS document")
		resp.Diagnostics.AddError("Unable to get the JWKS document", "Provided URL is not valid or not available.")
		return
	}
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		tflog.Error(ctx, "Unable to read the JWKS document")
		resp.Diagnostics.AddError("Unable to read the JWKS document", "Document found is not readable.")
		return
	}

	err = json.Unmarshal(b, &jwksDocument)
	if err != nil {
		tflog.Error(ctx, "Unable to decode the JWKS document")
		resp.Diagnostics.AddError("Unable to decode the JWKS document", "Document found is not in the right JSON format.")
		return
	}

	var keysAsString []string
	var pemAsString []string
	for _, k := range jwksDocument.Keys {

		key, err := json.Marshal(k)
		if err != nil {
			tflog.Error(ctx, "Unable to encode the JWKS key as string")
			resp.Diagnostics.AddError("Unable to encode the JWKS key as string", "The JSON key can't be converted to string.")
		}
		keysAsString = append(keysAsString, string(key))
		pem, err := convertJkwsToPem(key)
		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("Unable to convert JWKS keys as PEM: %v", err))
			resp.Diagnostics.AddError("Unable to convert JWKS keys as PEM", "The JSON key can't be converted to PEM.")
		}
		pemAsString = append(pemAsString, strings.TrimSpace(string(pem)))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.ID = data.JWKSURI
	data.Keys, diagnostics = types.ListValueFrom(ctx, types.StringType, keysAsString)
	resp.Diagnostics.Append(diagnostics...)

	data.PEM, diagnostics = types.ListValueFrom(ctx, types.StringType, pemAsString)
	resp.Diagnostics.Append(diagnostics...)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
