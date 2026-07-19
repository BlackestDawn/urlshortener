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

variable "apps" {
  type = map(object({
    github_repo = string
  }))
  description = "One entry per app deployed into this project. The map key is used as the app name (service account ids, secret names, Artifact Registry image name). Add a new entry here for each future project instead of writing new Terraform."
}
