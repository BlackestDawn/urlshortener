terraform {
  required_version = ">= 1.7.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 6.0"
    }
  }

  # Local state to start. Once you have a GCS bucket for state (create it
  # once, by hand, in whichever project you use for shared tooling), switch
  # to a remote backend so state isn't only on your laptop:
  #
  # backend "gcs" {
  #   bucket = "your-terraform-state-bucket"
  #   prefix = "gcp-apps"
  # }
}

provider "google" {
  project = var.project_id
  region  = var.region
}
