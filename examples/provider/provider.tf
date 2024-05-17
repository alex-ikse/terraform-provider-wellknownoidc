provider "wellknownoidc" {
}

terraform {
  required_providers {
    wellknownoidc = {
      source = "alex-ikse/wellknownoidc"
    }
  }
}
