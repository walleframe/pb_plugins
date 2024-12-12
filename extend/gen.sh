
protoc --go_out=./mysql/ --go_opt=paths=source_relative ./mysql.proto
protoc --go_out=./redis/ --go_opt=paths=source_relative ./redis.proto
protoc --go_out=./wrpc/ --go_opt=paths=source_relative ./wrpc.proto
