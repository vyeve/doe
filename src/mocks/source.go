package mocks

// Logger
//go:generate mockgen -destination=./mock-logger.go -package=mocks -mock_names=Logger=MockLogger doe/src/logger Logger

// Data
//go:generate mockgen -destination=./mock-repository.go -package=mocks -mock_names=Repository=MockRepository doe/src/data Repository
//go:generate mockgen -destination=./mock-sql-tx.go -package=mocks -mock_names=SQLTx=MockSQLTx doe/src/data SQLTx
//go:generate mockgen -destination=./mock-sql-db.go -package=mocks -mock_names=SQLTb=MockSQLDb doe/src/data SQLDb
