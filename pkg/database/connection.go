package db

import (
	"fmt"
	"os"
	"time"

	"github.com/yorukot/go-template/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB is the global GORM database instance
var DB *gorm.DB

//-----------------------------------------------------------------------------
// Database Types and Configuration
//-----------------------------------------------------------------------------

// Supported database types
const (
	PostgreSQL = "postgres" // PostgreSQL database
	MySQL      = "mysql"    // MySQL database
	MariaDB    = "mariadb"  // MariaDB database (uses MySQL driver)
	SQLite     = "sqlite"   // SQLite database
)

// init initializes the database connection when the package is imported
// Database configuration is read from environment variables:
//   - DATABASE_TYPE: Type of database to connect to (default: postgres)
//   - DATABASE_HOST: Database server hostname or IP
//   - DATABASE_PORT: Database server port
//   - DATABASE_USER: Database username
//   - DATABASE_PASSWORD: Database password
//   - DATABASE_DBNAME: Database name
//   - DATABASE_SSLMODE: SSL mode for PostgreSQL
//   - DATABASE_PATH: File path for SQLite database
//   - DATABASE_MAX_IDLE_CONNS: Maximum number of idle connections
//   - DATABASE_MAX_OPEN_CONNS: Maximum number of open connections
//   - DATABASE_CONN_MAX_LIFETIME: Maximum lifetime of connections in minutes
func init() {
	// Get database type from environment
	dbType := os.Getenv("DATABASE_TYPE")
	if dbType == "" {
		dbType = PostgreSQL // Default to PostgreSQL if not specified
		logger.Log.Sugar().Info("DATABASE_TYPE not set, defaulting to PostgreSQL")
	}

	var err error

	// Initialize the appropriate database based on type
	switch dbType {
	case PostgreSQL:
		logger.Log.Sugar().Info("Initializing PostgreSQL connection")
		DB, err = initPostgreSQL()
	case MySQL, MariaDB:
		logger.Log.Sugar().Infof("Initializing %s connection", dbType)
		DB, err = initMySQL()
	case SQLite:
		logger.Log.Sugar().Info("Initializing SQLite connection")
		DB, err = initSQLite()
	default:
		logger.Log.Sugar().Fatalf("Unsupported database type: %s", dbType)
	}

	if err != nil {
		logger.Log.Sugar().Fatalf("Failed to connect to %s database: %v", dbType, err)
	}

	// Configure connection pool
	configureConnectionPool(dbType)

	logger.Log.Sugar().Infof("Successfully connected to %s database", dbType)
}

// configureConnectionPool sets up the database connection pool parameters
func configureConnectionPool(dbType string) {
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Log.Sugar().Fatalf("Failed to get SQL DB instance: %v", err)
	}

	// Get connection pool settings from environment or use defaults
	maxIdleConns := getEnvIntWithDefault("DATABASE_MAX_IDLE_CONNS", 10)
	maxOpenConns := getEnvIntWithDefault("DATABASE_MAX_OPEN_CONNS", 100)
	connMaxLifetime := time.Duration(getEnvIntWithDefault("DATABASE_CONN_MAX_LIFETIME", 30)) * time.Minute

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	// Check connection
	if err = sqlDB.Ping(); err != nil {
		logger.Log.Sugar().Fatalf("Failed to ping %s database: %v", dbType, err)
	}
}

//-----------------------------------------------------------------------------
// Database Initialization Functions
//-----------------------------------------------------------------------------

// initPostgreSQL initializes a PostgreSQL connection
// Required env vars: DATABASE_HOST, DATABASE_USER, DATABASE_PASSWORD,
// DATABASE_DBNAME, DATABASE_PORT, DATABASE_SSLMODE
func initPostgreSQL() (*gorm.DB, error) {
	host := os.Getenv("DATABASE_HOST")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_DBNAME")
	port := os.Getenv("DATABASE_PORT")
	sslmode := os.Getenv("DATABASE_SSLMODE")

	// Validate required parameters
	if host == "" || user == "" || dbname == "" || port == "" {
		return nil, fmt.Errorf("missing required PostgreSQL connection parameters")
	}

	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// initMySQL initializes a MySQL or MariaDB connection
// Required env vars: DATABASE_HOST, DATABASE_USER, DATABASE_PASSWORD,
// DATABASE_DBNAME, DATABASE_PORT
func initMySQL() (*gorm.DB, error) {
	host := os.Getenv("DATABASE_HOST")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_DBNAME")
	port := os.Getenv("DATABASE_PORT")

	// Validate required parameters
	if host == "" || user == "" || dbname == "" || port == "" {
		return nil, fmt.Errorf("missing required MySQL connection parameters")
	}

	// Build connection string
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local",
		user, password, host, port, dbname,
	)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

// initSQLite initializes a SQLite connection
// Required env vars: DATABASE_PATH (or uses default "database.db")
func initSQLite() (*gorm.DB, error) {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "database.db" // Default SQLite database file
		logger.Log.Sugar().Info("DATABASE_PATH not set, using default: database.db")
	}

	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}

//-----------------------------------------------------------------------------
// Utility Functions
//-----------------------------------------------------------------------------

// getEnvIntWithDefault gets an integer from environment variable or returns the default value
func getEnvIntWithDefault(key string, defaultVal int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultVal
	}

	var intVal int
	_, err := fmt.Sscanf(value, "%d", &intVal)
	if err != nil {
		logger.Log.Sugar().Warnf("Invalid value for %s: %s, using default: %d", key, value, defaultVal)
		return defaultVal
	}

	return intVal
}

// GetDB returns the GORM DB instance
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase closes the database connection
func CloseDatabase() {
	sqlDB, err := DB.DB()
	if err != nil {
		logger.Log.Sugar().Fatal("Failed to get SQL DB instance:", err)
	}

	if err := sqlDB.Close(); err != nil {
		logger.Log.Sugar().Fatal("Failed to close database connection:", err)
	}
	logger.Log.Info("Successfully disconnected from database")
}
