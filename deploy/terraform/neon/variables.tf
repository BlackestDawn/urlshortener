variable "app_name" {
  type        = string
  description = "Name for the Neon project (one project per app)."
}

variable "region_id" {
  type        = string
  default     = "aws-eu-central-1"
  description = "Neon region — see https://neon.tech/docs/introduction/regions. Neon runs on AWS regions; aws-eu-central-1 (Frankfurt) is the closest to europe-west1."
}

variable "pg_version" {
  type    = number
  default = 17
}
