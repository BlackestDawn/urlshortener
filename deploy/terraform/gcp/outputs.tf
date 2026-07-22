output "artifact_registry_repo" {
  value       = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.apps.repository_id}"
  description = "Prefix for image tags, e.g. <this>/urlshortener:<git-sha>"
}

output "workload_identity_provider" {
  value       = google_iam_workload_identity_pool_provider.github.name
  description = "Value for the GCP_WORKLOAD_IDENTITY_PROVIDER GitHub Actions secret."
}

output "deployer_service_accounts" {
  value       = google_service_account.deployer.email
  description = "Value for the GCP_DEPLOYER_SA GitHub Actions secret."
}

output "runtime_service_accounts" {
  value       = google_service_account.runtime.email
  description = "Value for the GCP_RUNTIME_SA GitHub Actions secret."
}

output "secret_ids" {
  value       = { for k, s in google_secret_manager_secret.db_url : k => s.secret_id }
  description = "Secret Manager secret IDs — populate versions with `gcloud secrets versions add <id> --data-file=-` using the Neon stack's outputs."
}
