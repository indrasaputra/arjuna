set -euo pipefail

migrate -source "file://service/user/db/migrations" -database "postgresql://postgresuser:postgrespassword@postgres:5432/arjuna_user?sslmode=disable" -verbose up
