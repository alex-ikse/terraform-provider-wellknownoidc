data "wellknownoidc_document" "example" {
  discovery_url = "https://container.googleapis.com/v1/projects/cloud-services-i2jm/locations/europe-west4/clusters/gcp-gke-cnp-cloudservices2"
}

data "wellknownoidc_jwks" "example" {
  jwks_uri = data.wellknownoidc_document.example.jwks_uri
}
