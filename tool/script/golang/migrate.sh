set -euo pipefail

atlas migrate apply --dir file://service/auth/db/migrations --url "postgresql://postgresuser:postgrespassword@postgres:5432/arjuna_auth?sslmode=disable"
atlas migrate apply --dir file://service/user/db/migrations --url "postgresql://postgresuser:postgrespassword@postgres:5432/arjuna_user?sslmode=disable"
atlas migrate apply --dir file://service/wallet/db/migrations --url "postgresql://postgresuser:postgrespassword@postgres:5432/arjuna_wallet?sslmode=disable"
atlas migrate apply --dir file://service/transaction/db/migrations --url "postgresql://postgresuser:postgrespassword@postgres:5432/arjuna_transaction?sslmode=disable"
