data "wellknownoidc_document" "example" {
  discovery_url = "https://account.google.com"
}

data "wellknownoidc_jwks" "example" {
  jwks_uri = data.wellknownoidc_document.example.jwks_uri
}
