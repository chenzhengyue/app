module app

go 1.13

replace comm => ./comm

replace logger => ./logger

require (
	comm v0.0.0-00010101000000-000000000000
	github.com/chenzhengyue/logger v0.0.0-20200417201020-0122311c54bb
	logger v0.0.0-00010101000000-000000000000
)
