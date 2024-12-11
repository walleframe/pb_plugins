
protoc --go_out=./sql/ --go_opt=paths=source_relative ./sql.proto
protoc --go_out=./redis/ --go_opt=paths=source_relative ./redis.proto
protoc --go_out=./wrpc/ --go_opt=paths=source_relative ./wrpc.proto
