module app

go 1.13

replace comm => ./src/comm

replace logger => ./src/logger

require (
	comm v0.0.0-00010101000000-000000000000
	github.com/chenzhengyue/logger v0.0.0-20200417201020-0122311c54bb
	github.com/go-sql-driver/mysql v1.5.0
	logger v0.0.0-00010101000000-000000000000
)
