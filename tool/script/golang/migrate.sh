set -euo pipefail

migrate -source "file://service/user/db/migrations" -database "postgresql://postgresuser:postgrespassword@postgres:5432/arjuna_user?sslmode=disable" -verbose up
migrate -source "file://service/auth/db/migrations" -database "postgresql://postgresuser:postgrespassword@postgres:5432/arjuna_auth?sslmode=disable" -verbose up
migrate -source "file://service/transaction/db/migrations" -database "postgresql://postgresuser:postgrespassword@postgres:5432/arjuna_transaction?sslmode=disable" -verbose up
