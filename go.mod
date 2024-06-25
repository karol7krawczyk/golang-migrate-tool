module github.com/Karol7Krawczyk/golang-migrate

go 1.22.4

require (
	github.com/Karol7Krawczyk/golang-migrate/migrations/config v0.0.0-00010101000000-000000000000
	github.com/Karol7Krawczyk/golang-migrate/migrations/db v0.0.0-00010101000000-000000000000
	github.com/Karol7Krawczyk/golang-migrate/migrations/handlers v0.0.0-00010101000000-000000000000
	github.com/go-sql-driver/mysql v1.8.1
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.22
)

require filippo.io/edwards25519 v1.1.0 // indirect

replace github.com/Karol7Krawczyk/golang-migrate/migrations/config => ./migrations/config

replace github.com/Karol7Krawczyk/golang-migrate/migrations/db => ./migrations/db

replace github.com/Karol7Krawczyk/golang-migrate/migrations/handlers => ./migrations/handlers
