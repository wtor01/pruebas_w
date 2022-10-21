openapi:
	@./scripts/openapi.sh inventory cmd/inventory/http http
	@./scripts/openapi.sh auth cmd/auth/http http
	@./scripts/openapi.sh validations cmd/validations/http/validations validations
	@./scripts/openapi.sh master_tables/geographic cmd/master_tables/http/geographic geographic
	@./scripts/openapi.sh master_tables/calendar cmd/master_tables/http/calendar calendar
	@./scripts/openapi.sh master_tables/festive_days cmd/master_tables/http/festive_days festive_days
	@./scripts/openapi.sh master_tables/seasons cmd/master_tables/http/seasons seasons
	@./scripts/openapi.sh master_tables/tariff cmd/master_tables/http/tariff tariff
	@./scripts/openapi.sh master_tables/generate_calendars cmd/master_tables/http/generate_calendars generate_calendars
	@./scripts/openapi.sh billing_measures/scheduler cmd/billing_measures/http/scheduler scheduler
	@./scripts/openapi.sh billing_measures/dashboard cmd/billing_measures/http/dashboard dashboard

	@./scripts/openapi.sh billing_measures/self_consumption cmd/billing_measures/http/self_consumption self_consumption

	@./scripts/openapi.sh gross_measures/dashboard_stats cmd/gross_measures/http/dashboard_stats dashboard_stats
	@./scripts/openapi.sh gross_measures/dashboard cmd/gross_measures/http/dashboard dashboard
	@./scripts/openapi.sh process_measures/dashboard cmd/process_measures/http/dashboard dashboard
	@./scripts/openapi.sh process_measures/dashboard_process_measures_stats cmd/process_measures/http/dashboard_stats stats
	@./scripts/openapi.sh process_measures/closure cmd/process_measures/http/closures closures
	@./scripts/openapi.sh process_measures/scheduler cmd/process_measures/http/scheduler scheduler
	@./scripts/openapi.sh aggregations/config cmd/aggregations/http/config config
	@./scripts/openapi.sh aggregations/features cmd/aggregations/http/features	features
	@./scripts/openapi.sh aggregations/aggregations cmd/aggregations/http/aggregations aggregations

	@./scripts/openapi.sh gross_measures/smarkia cmd/gross_measures/http/smarkia smarkia
	@./scripts/openapi.sh gross_measures/supply_point cmd/gross_measures/http/supply_point supply_point




coverage:
	@go test ./... -cover -coverprofile=coverage.out
	@go tool cover -html=coverage.out

openapi-merge:
	@go run cmd/openapi/merge/main.go \
		api/openapi/api.yaml \
		api/openapi/validations.yaml \
		api/openapi/auth.yaml \
		api/openapi/dashboard.yaml \
		api/openapi/inventory.yaml \
		api/openapi/master_tables/geographic.yaml \
		api/openapi/master_tables/festive_days.yaml \
		api/openapi/master_tables/seasons.yaml \
		api/openapi/master_tables/tariff.yaml \
		api/openapi/master_tables/calendar.yaml \
		api/openapi/master_tables/generate_calendars.yaml \
		api/openapi/billing_measures/scheduler.yaml \
		api/openapi/billing_measures/dashboard.yaml \
		api/openapi/billing_measures/self_consumption.yaml \
		api/openapi/process_measures/dashboard.yaml \
		api/openapi/process_measures/scheduler.yaml \
		api/openapi/process_measures/dashboard_process_measures_stats.yaml \
		api/openapi/aggregations/config.yaml \
		api/openapi/aggregations/features.yaml \
		api/openapi/gross_measures/dashboard.yaml \
		api/openapi/gross_measures/dashboard_stats.yaml \
		api/openapi/aggregations/aggregations.yaml \
		api/openapi/gross_measures/smarkia.yaml \
		api/openapi/gross_measures/supply_point.yaml

openapi-merge-master_tables:
	@go run cmd/openapi/merge/main.go \
		api/openapi/master_tables/api.yaml \
		api/openapi/master_tables/geographic.yaml \
		api/openapi/master_tables/festive_days.yaml \
		api/openapi/master_tables/seasons.yaml \
		api/openapi/master_tables/tariff.yaml \
		api/openapi/master_tables/calendar.yaml \
		api/openapi/master_tables/generate_calendars.yaml

openapi-merge-billing-measures:
	@go run cmd/openapi/merge/main.go \
		api/openapi/billing_measures/api.yaml \
		api/openapi/billing_measures/scheduler.yaml \
		api/openapi/billing_measures/dashboard.yaml \
		api/openapi/billing_measures/self_consumption.yaml



openapi-merge-process-measures:
	@go run cmd/openapi/merge/main.go \
		api/openapi/process_measures/api.yaml \
		api/openapi/process_measures/scheduler.yaml \
		api/openapi/process_measures/dashboard.yaml \
		api/openapi/process_measures/measures.yaml \
		api/openapi/process_measures/dashboard_process_measures_stats.yaml

openapi-merge-gross-measures:
	@go run cmd/openapi/merge/main.go \
		api/openapi/gross_measures/api.yaml \
		api/openapi/gross_measures/dashboard_stats.yaml \
		api/openapi/gross_measures/dashboard.yaml \
		api/openapi/gross_measures/supply_point.yaml \
		api/openapi/gross_measures/smarkia.yaml



openapi-merge-aggregations:
	@go run cmd/openapi/merge/main.go \
		api/openapi/aggregations/api.yaml \
		api/openapi/aggregations/features.yaml \
		api/openapi/aggregations/config.yaml \
		api/openapi/aggregations/aggregations.yaml
