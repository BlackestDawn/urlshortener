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
  ]
}

resource "google_project_service" "this" {
  for_each           = toset(local.apis)
  project            = var.project_id
  service            = each.value
  disable_on_destroy = false
}

# Shared by every app in this project.
resource "google_artifact_registry_repository" "apps" {
  location      = var.region
  repository_id = "apps"
  format        = "DOCKER"
  description   = "Container images for personal apps"
  depends_on    = [google_project_service.this]
}

# One-time per GCP project: WIF pool + provider, reused by every app/repo
# below. Adding a new app does not touch these.
resource "google_iam_workload_identity_pool" "github_actions" {
  workload_identity_pool_id = "github-actions"
  display_name              = "GitHub Actions"
  depends_on                = [google_project_service.this]
}

resource "google_iam_workload_identity_pool_provider" "github" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.github_actions.workload_identity_pool_id
  workload_identity_pool_provider_id = "github"
  display_name                       = "GitHub OIDC"

  attribute_mapping = {
    "google.subject"             = "assertion.sub"
    "attribute.repository"       = "assertion.repository"
    "attribute.repository_owner" = "assertion.repository_owner"
  }

  # Restrict token exchange to repos owned by you, not just any GitHub repo.
  attribute_condition = "assertion.repository_owner == '${var.github_org}'"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
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
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github_actions.name}/attribute.repository/${var.github_org}/${var.github_repo}"
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
