module vsys.dbhelper

go 1.20

require (
	github.com/gorilla/mux v1.8.1
	vsys.commons v0.0.0
)

require (
	github.com/go-sql-driver/mysql v1.8.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	gorm.io/driver/mysql v1.5.4
	gorm.io/gorm v1.25.7
)

replace vsys.commons => ../vsys.commons
