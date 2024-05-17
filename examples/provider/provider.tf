provider "wellknownoidc" {
}

terraform {
  required_providers {
    wellknownoidc = {
      source = "registry.terraform.io/alex-ikse/wellknownoidc"
    }
  }
}
