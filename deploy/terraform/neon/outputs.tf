output "prod_connection_uri" {
  value       = neon_project.this.connection_uri
  description = "Push into Secret Manager: terraform output -raw prod_connection_uri | gcloud secrets versions add <app>-prod-database-url --data-file=-"
  sensitive   = true
}

output "staging_connection_uri" {
  value       = "postgres://${neon_project.this.database_user}:${urlencode(neon_project.this.database_password)}@${neon_endpoint.staging.host}/${neon_project.this.database_name}?sslmode=require"
  description = "Push into Secret Manager: terraform output -raw staging_connection_uri | gcloud secrets versions add <app>-staging-database-url --data-file=-"
  sensitive   = true
}
