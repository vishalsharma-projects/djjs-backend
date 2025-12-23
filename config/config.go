package config

import (
	"context"
    "fmt"
    "log"
    "net/url"
    "os"
    "strconv"
    "time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

	"github.com/followCode/djjs-event-reporting-backend/app/models"
)

// Legacy GORM connection (for existing code)
var DB *gorm.DB

// New pgx connection pool (for auth system)
var AuthDB *pgxpool.Pool

// Redis client (for rate limiting)
var RedisClient *redis.Client

// JWT Configuration
var JWTSecret []byte

// JWT Token Configuration
var JWTTTL time.Duration = 10 * time.Minute
var JWTIssuer string
var JWTAudience string

// Token Configuration
var TokenPepper []byte
var RefreshTokenTTL time.Duration = 30 * 24 * time.Hour // 30 days
var VerificationTTL time.Duration = 30 * time.Minute
var PasswordResetTTL time.Duration = 30 * time.Minute

// Cookie Configuration
var CookieSecure bool
var CookieSameSite string = "Lax"
var CookiePath string = "/" // Changed from "/auth/refresh" to "/" so cookie is available for all API requests

// Security Configuration
var RequireEmailVerified bool
var FrontendOrigin string
var TrustProxy bool

// Rate Limiting Configuration
var RateLimitLoginPerIP int = 5
var RateLimitLoginPerEmail int = 3
var RateLimitForgotPasswordPerIP int = 3
var RateLimitForgotPasswordPerEmail int = 2
var RateLimitWindow time.Duration = 15 * time.Minute

func LoadJWTSecret() {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Fatal("JWT_SECRET is not set in environment")
    }
    JWTSecret = []byte(secret)
}

func ConnectDB() {
    dbUser := os.Getenv("POSTGRES_USER") 
    dbPass := os.Getenv("POSTGRES_PASSWORD")
    dbName := os.Getenv("POSTGRES_DB")
    dbPort := os.Getenv("PG_PORT")
    dbHost := os.Getenv("POSTGRES_HOST")

    // Validate required environment variables
    if dbHost == "" {
        log.Fatal("POSTGRES_HOST is required in .env or environment variables")
    }
    if dbUser == "" {
        log.Fatal("POSTGRES_USER is required in .env or environment variables")
    }
    if dbPass == "" {
        log.Fatal("POSTGRES_PASSWORD is required in .env or environment variables")
    }
    if dbName == "" {
        log.Fatal("POSTGRES_DB is required in .env or environment variables")
    }
    if dbPort == "" {
        dbPort = "5432" // Default PostgreSQL port
    }


    // URL encode password and other components to handle special characters like @, #, etc.
    // Using connection URI format which handles special characters more reliably
    encodedUser := url.QueryEscape(dbUser)
    encodedPassword := url.QueryEscape(dbPass)
    encodedDBName := url.QueryEscape(dbName)
    encodedHost := url.QueryEscape(dbHost)

    // Build connection URI with connection timeout for remote databases
    // Format: postgres://user:password@host:port/dbname?sslmode=disable&connect_timeout=10
    dsn := fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=10",
        encodedUser, encodedPassword, encodedHost, dbPort, encodedDBName,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to DB:", err)
    }

    // Configure connection pool for better performance and scalability
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal("Failed to get underlying sql.DB:", err)
    }

    // SetMaxIdleConns sets the maximum number of connections in the idle connection pool
    sqlDB.SetMaxIdleConns(10)
    
    // SetMaxOpenConns sets the maximum number of open connections to the database
    sqlDB.SetMaxOpenConns(100)
    
    // SetConnMaxLifetime sets the maximum amount of time a connection may be reused
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    // Set connection timeout for establishing new connections
    sqlDB.SetConnMaxIdleTime(5 * time.Minute)

    DB = db
    log.Println("Database connection pool configured successfully")
}

func AutoMigrate() {
    DB.AutoMigrate(&models.Role{}, &models.User{})
}

// LoadAuthConfig loads configuration for the new auth system (pgx + Redis)
func LoadAuthConfig() error {
	// Load JWT secret (required)
	jwtSecretStr := os.Getenv("JWT_SECRET")
	if jwtSecretStr == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	JWTSecret = []byte(jwtSecretStr)

	// Load Token Pepper (required)
	pepperStr := os.Getenv("TOKEN_PEPPER")
	if pepperStr == "" {
		return fmt.Errorf("TOKEN_PEPPER is required")
	}
	TokenPepper = []byte(pepperStr)

	// Database URL for pgx
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Build from components
		dbUser := os.Getenv("POSTGRES_USER")
		dbPass := os.Getenv("POSTGRES_PASSWORD")
		dbHost := os.Getenv("POSTGRES_HOST")
		dbPort := os.Getenv("PG_PORT")
		if dbPort == "" {
			dbPort = "5432"
		}
		dbName := os.Getenv("POSTGRES_DB")

		if dbUser == "" || dbPass == "" || dbHost == "" || dbName == "" {
			return fmt.Errorf("database configuration missing: need DATABASE_URL or POSTGRES_* variables")
		}

		encodedUser := url.QueryEscape(dbUser)
		encodedPassword := url.QueryEscape(dbPass)
		dbURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", encodedUser, encodedPassword, dbHost, dbPort, dbName)
	}

	// Connect to PostgreSQL with pgx
	ctx := context.Background()
	var err error
	AuthDB, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := AuthDB.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Connect to Redis (optional - rate limiting will be disabled if not available)
	redisURL := os.Getenv("REDIS_URL")
	if redisURL != "" {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Printf("Warning: Failed to parse Redis URL: %v (rate limiting will be disabled)", err)
		} else {
			RedisClient = redis.NewClient(opt)
			// Test Redis connection (non-fatal)
			if err := RedisClient.Ping(ctx).Err(); err != nil {
				log.Printf("Warning: Redis connection failed: %v (rate limiting will be disabled)", err)
				RedisClient = nil
			} else {
				log.Println("Redis connected successfully")
			}
		}
	} else {
		log.Println("Redis not configured (REDIS_URL not set) - rate limiting will be disabled")
	}

	// JWT TTL (optional, default 10 min)
	if ttlStr := os.Getenv("JWT_TTL"); ttlStr != "" {
		if ttl, err := time.ParseDuration(ttlStr); err == nil {
			JWTTTL = ttl
		}
	}

	// JWT Issuer/Audience
	JWTIssuer = os.Getenv("JWT_ISSUER")
	if JWTIssuer == "" {
		JWTIssuer = "djjs-backend"
	}
	JWTAudience = os.Getenv("JWT_AUDIENCE")
	if JWTAudience == "" {
		JWTAudience = "djjs-frontend"
	}

	// Cookie settings
	CookieSecure = os.Getenv("COOKIE_SECURE") != "false"
	if sameSite := os.Getenv("COOKIE_SAME_SITE"); sameSite != "" {
		CookieSameSite = sameSite
	}
	if path := os.Getenv("COOKIE_PATH"); path != "" {
		CookiePath = path
	}

	// Security settings
	RequireEmailVerified = os.Getenv("REQUIRE_EMAIL_VERIFIED") == "true"
	FrontendOrigin = os.Getenv("FRONTEND_ORIGIN")
	if FrontendOrigin == "" {
		// Only default to localhost if explicitly in debug mode
		// In production, this should be set via environment variable
		if os.Getenv("GIN_MODE") == "debug" {
			FrontendOrigin = "http://localhost:4200"
		}
	}
	TrustProxy = os.Getenv("TRUST_PROXY") == "true"

	// Rate limiting (optional overrides)
	if val := os.Getenv("RATE_LIMIT_LOGIN_PER_IP"); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			RateLimitLoginPerIP = n
		}
	}
	if val := os.Getenv("RATE_LIMIT_LOGIN_PER_EMAIL"); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			RateLimitLoginPerEmail = n
		}
	}
	if val := os.Getenv("RATE_LIMIT_WINDOW"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			RateLimitWindow = d
		}
	}

	log.Println("Auth configuration loaded successfully")
	return nil
}
