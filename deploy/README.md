# Deployment

```
merge to main ──▶ build & push image ──▶ migrate staging DB ──▶ deploy Cloud Run staging
                                                                        │
                                                            (manual approval)
                                                                        ▼
                                                    promote the *same* image ──▶ migrate prod DB ──▶ deploy Cloud Run prod
```

- **Compute**: Google Cloud Run, one service (`urlshortener-staging`, `urlshortener-prod`) per environment.
- **Database**: [Neon](https://neon.tech) serverless Postgres — prod and a copy-on-write staging branch.
- **Auth**: GitHub Actions authenticates to GCP via Workload Identity Federation — no long-lived service account keys.
- **Provisioning**: `deploy/terraform/gcp` and `deploy/terraform/neon`. The GCP stack references a shared Workload Identity Pool and Artifact Registry repo, provisioned once by a separate bootstrap stack ([various-terraform/gcp-bootstrap](https://github.com/BlackestDawn/various-terraform)) and reused (via Terraform data sources, not recreated) by every app deployed into this GCP project.
- **Migrations**: run as a CI step (`make migrate-up`) against the target environment's DB before each deploy, not automatically on boot.

## One-time bootstrap

> **Prerequisite**: the shared Workload Identity Pool/provider (`github-actions`) and the `apps` Artifact Registry repo must already exist in the target GCP project before running anything below — `deploy/terraform/gcp` only reads them via data sources, it doesn't create them. Apply [various-terraform/gcp-bootstrap](https://github.com/BlackestDawn/various-terraform) first if this is a new GCP project.

1. **GCP stack** — creates this app's deployer/runtime service accounts, IAM bindings, and empty Secret Manager secret containers.

   ```bash
   cd deploy/terraform/gcp
   cp terraform.tfvars.example terraform.tfvars   # fill in project_id, github_org
   terraform init
   terraform apply
   ```

2. **Neon stack** — creates the Neon project plus a staging branch.

   ```bash
   cd deploy/terraform/neon
   export NEON_API_KEY=...                        # https://console.neon.tech/app/settings/api-keys
   cp terraform.tfvars.example terraform.tfvars
   terraform init
   terraform apply
   ```

3. **Populate secrets** — Terraform only creates the secret containers, not their values.

   ```bash
   cd deploy/terraform/neon
   terraform output -raw prod_connection_uri    | gcloud secrets versions add urlshortener-prod-database-url    --data-file=-
   terraform output -raw staging_connection_uri | gcloud secrets versions add urlshortener-staging-database-url --data-file=-
   ```

4. **GitHub repo secrets** (Settings → Secrets and variables → Actions):

   From `deploy/terraform/gcp` outputs:
   - `GCP_PROJECT_ID`
   - `GCP_WORKLOAD_IDENTITY_PROVIDER` → `terraform output workload_identity_provider`
   - `GCP_DEPLOYER_SA` → `terraform output deployer_service_accounts`
   - `GCP_RUNTIME_SA` → `terraform output runtime_service_accounts`

5. **GitHub environment** — create an environment named `production` (Settings → Environments) with a required reviewer, so `deploy-prod.yml` always pauses for manual approval.

6. **First deploy** — push to `main` (or run "Deploy to staging" manually) to create the staging Cloud Run service, then run "Deploy to production" manually once staging looks good. The domain mapping in the next step needs these services to already exist.

7. **Domain mappings** — Cloud Run custom domains aren't in Terraform; one-time CLI step per service:

   ```bash
   PROJECT=your-gcp-project-id
   REGION=your-gcp-region
   DOMAIN=your-domain

   gcloud beta run domain-mappings create --service=urlshortener-prod    --domain=shortener.$DOMAIN         --region=$REGION --project=$PROJECT
   gcloud beta run domain-mappings create --service=urlshortener-staging --domain=shortener-staging.$DOMAIN --region=$REGION --project=$PROJECT
   ```

   Then create the matching DNS record in whichever Cloud DNS zone hosts `$DOMAIN` for each — Cloud Run subdomain mappings resolve via a CNAME to `ghs.googlehosted.com.`:

   ```bash
   PROJECT=your-gcp-project-id
   DOMAIN=your-domain
   ZONE=your-cloud-dns-zone-name   # e.g. `gcloud dns managed-zones list --project=$PROJECT`

   gcloud dns record-sets create shortener.$DOMAIN. \
     --zone=$ZONE \
     --type=CNAME \
     --ttl=300 \
     --rrdatas=ghs.googlehosted.com. \
     --project=$PROJECT
   ```

   Repeat for `shortener-staging`.
