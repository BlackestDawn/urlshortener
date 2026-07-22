variable "project_id" {
  type        = string
  description = "GCP project shared by all your personal apps."
}

variable "region" {
  type        = string
  default     = "europe-west1"
  description = "Region for Cloud Run and Artifact Registry."
}

variable "github_org" {
  type        = string
  description = "GitHub org/username that owns the repos allowed to deploy into this project."
}

variable "github_repo" {
  type = string
  default = "urlshortener"
  description = "Github repository name."
}
