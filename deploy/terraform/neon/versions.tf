terraform {
  required_version = ">= 1.7.0"

  required_providers {
    neon = {
      source  = "kislerdm/neon"
      version = "~> 0.14"
    }
  }

  # See the note in ../gcp/versions.tf about moving to a shared GCS backend.
}

provider "neon" {
  # Reads NEON_API_KEY from the environment by default — create one at
  # https://console.neon.tech/app/settings/api-keys and export it rather
  # than putting it in a .tf/.tfvars file.
}
