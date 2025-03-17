module github.com/sawyer523/gocqlx

go 1.24.1

require (
	github.com/google/go-cmp v0.7.0
	github.com/scylladb/gocqlx/v3 v3.0.1
)

require (
	github.com/gocql/gocql v1.7.0 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/scylladb/go-reflectx v1.0.1 // indirect
	golang.org/x/sync v0.12.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)

replace github.com/gocql/gocql => github.com/scylladb/gocql v1.14.4
