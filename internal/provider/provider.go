// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure WellKnownOIDCProvider satisfies various provider interfaces.
var _ provider.Provider = &WellKnownOIDCProvider{}

// WellKnownOIDCProvider defines the provider implementation.
type WellKnownOIDCProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// WellKnownOIDCProviderModel describes the provider data model.
type WellKnownOIDCProviderModel struct {
}

func (p *WellKnownOIDCProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "wellknownoidc"
	resp.Version = p.version
}

func (p *WellKnownOIDCProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{},
	}
}

func (p *WellKnownOIDCProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data WellKnownOIDCProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = ""
	resp.ResourceData = ""
}

func (p *WellKnownOIDCProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *WellKnownOIDCProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDocumentDataSource,
		NewJWKSDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &WellKnownOIDCProvider{
			version: version,
		}
	}
}
