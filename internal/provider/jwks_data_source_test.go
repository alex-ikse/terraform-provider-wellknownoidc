// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccJWKSDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + fmt.Sprintf(testAccJWKSDataSourceConfig, issuer),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.wellknownoidc_jwks.test", "keys.#", "2"),
				),
			},
		},
	})
}

const testAccJWKSDataSourceConfig = `
data "wellknownoidc_document" "test" {
	discovery_url = "%s"
}

data "wellknownoidc_jwks" "test" {
	jwks_uri = data.wellknownoidc_document.test.jwks_uri
}
`
