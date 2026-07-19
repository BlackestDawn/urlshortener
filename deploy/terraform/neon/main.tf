# The project's default branch is prod. Its role/database/connection
# details are all computed attributes on neon_project itself.
resource "neon_project" "this" {
  name                      = var.app_name
  region_id                 = var.region_id
  pg_version                = var.pg_version
  # default value for retention is too high for free accounts
  history_retention_seconds = 21600
}

# Branching is copy-on-write: staging starts with prod's data, roles and
# database already cloned onto it (same user/password) — no separate
# neon_role/neon_database needed here.
resource "neon_branch" "staging" {
  project_id = neon_project.this.id
  parent_id  = neon_project.this.default_branch_id
  name       = "staging"
}

resource "neon_endpoint" "staging" {
  project_id = neon_project.this.id
  branch_id  = neon_branch.staging.id
  type       = "read_write"

  # Nobody's paying for idle staging compute — suspend fast.
  # Only available for paid plans
  # suspend_timeout_seconds = 60
}
