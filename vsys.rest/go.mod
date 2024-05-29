module vsys.rest

go 1.20

require (
	github.com/gorilla/mux v1.8.1
	vsys.commons v0.0.0
	vsys.dbhelper v0.0.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	gorm.io/driver/mysql v1.5.4 // indirect
	gorm.io/gorm v1.25.7 // indirect
)

replace vsys.dbhelper => ../vsys.dbhelper

replace vsys.commons => ../vsys.commons
