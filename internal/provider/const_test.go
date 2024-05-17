package provider

const (
	// issuer = "https://container.googleapis.com/v1/projects/cloud-services-i2jm/locations/europe-west4/clusters/gcp-gke-cnp-cloudservices2"
	issuer = "https://accounts.google.com"
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the wellknownoidc client is properly configured.
	providerConfig = `
terraform {
  required_providers {
    wellknownoidc = {
      source = "registry.terraform.io/alex-ikse/wellknownoidc"
    }
  }
}

provider "wellknownoidc" {
}
`
)
