locals {
  apis = [
    "run.googleapis.com",
    "artifactregistry.googleapis.com",
    "iamcredentials.googleapis.com",
    "secretmanager.googleapis.com",
    "sts.googleapis.com",
  ]

  envs = ["staging", "prod"]

  deployer_roles = [
    "roles/run.admin",
    "roles/artifactregistry.writer",
    "roles/iam.serviceAccountUser",
    "roles/secretmanager.secretAccessor",
  ]
}

# --- Shared by every app in this project ------------------------------------
#
# The WIF pool/provider and the "apps" Artifact Registry repo are owned by a
# separate bootstrap stack (not tied to any individual app repo) and shared
# by every app deployed into this GCP project. They're referenced here
# read-only via data sources rather than recreated, since this stack has its
# own local Terraform state and re-declaring these as resources here would
# try to create already-existing, uniquely-named GCP objects and fail with
# 409s.

data "google_iam_workload_identity_pool" "github_actions" {
  workload_identity_pool_id = "github-actions"
}

data "google_iam_workload_identity_pool_provider" "github" {
  workload_identity_pool_id          = data.google_iam_workload_identity_pool.github_actions.workload_identity_pool_id
  workload_identity_pool_provider_id = "github"
}

data "google_artifact_registry_repository" "apps" {
  location      = var.region
  repository_id = "apps"
}

# API enablement is idempotent (enabling an already-enabled API just
# succeeds), so it's safe for this stack to also declare it rather than
# assume the bootstrap stack's state will always be the one managing it.
resource "google_project_service" "this" {
  for_each           = toset(local.apis)
  project            = var.project_id
  service            = each.value
  disable_on_destroy = false
}

# --- Per app ---------------------------------------------------------------

# Used by GitHub Actions to build/push/deploy.
resource "google_service_account" "deployer" {
  account_id   = "${var.github_repo}-deployer"
  display_name = "${var.github_repo} GitHub Actions deployer"
}

# Used as the running Cloud Run service's own identity.
resource "google_service_account" "runtime" {
  account_id   = "${var.github_repo}-run"
  display_name = "${var.github_repo} Cloud Run runtime identity"
}

resource "google_project_iam_member" "deployer" {
  for_each = toset(local.deployer_roles)
  project  = var.project_id
  role     = each.value
  member   = "serviceAccount:${google_service_account.deployer.email}"
}

resource "google_project_iam_member" "runtime_secret_access" {
  project  = var.project_id
  role     = "roles/secretmanager.secretAccessor"
  member   = "serviceAccount:${google_service_account.runtime.email}"
}

# Let GitHub Actions, but only from the matching repo, impersonate this
# app's deployer SA — no downloaded JSON keys.
resource "google_service_account_iam_member" "deployer_wif_binding" {
  service_account_id = google_service_account.deployer.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${data.google_iam_workload_identity_pool.github_actions.name}/attribute.repository/${var.github_org}/${var.github_repo}"
}

# Secret containers only — the actual connection string is pushed in
# separately (see deploy/terraform/neon output + gcloud secrets versions
# add), so it never lands in this stack's state.
resource "google_secret_manager_secret" "db_url" {
  for_each  = toset(local.envs)
  secret_id = "${var.github_repo}-${each.value}-database-url"

  replication {
    auto {}
  }

  depends_on = [google_project_service.this]
}
