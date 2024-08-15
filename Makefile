run/api:
	@go run ./cmd/api -db-dsn=${DSN} -smtp-username=${SMTP_USERNAME} -smtp-password=${SMTP_PASSWORD_DEMO}
