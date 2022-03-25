package store

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	redis "github.com/go-redis/redis/v8"
	pgx "github.com/jackc/pgx/v4"
)

var (
	storeService = &StorageService{}
	ctx          = context.Background()
)

type StorageService struct {
	redisClient        *redis.Client
	postgresConnection *pgx.Conn
	cacheDuration      time.Duration
}

func InitializeStore() *StorageService {
	postregsConnection, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_CONN_STRING"))
	if err != nil {
		panic(fmt.Sprintf("Error initializing PostgreSQL: %v", err))
	}

	redisDatabase, _ := strconv.Atoi(os.Getenv("REDIS_DATABSE"))
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDatabase,
	})

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Error initializing Redis: %v", err))
	}

	cacheDuration, err := strconv.Atoi(os.Getenv("REDIS_CACHE"))
	if err != nil {
		panic(fmt.Sprintf("Failed parsing default TTL | Error: %v", err))
	}

	storeService.redisClient = redisClient
	storeService.postgresConnection = postregsConnection
	storeService.cacheDuration = time.Duration(cacheDuration * 3600000000000)

	return storeService
}

func SaveUrlMapping(id int64, shortUrl string, completeUrl string, fallBackUrl string, ttl string) {
	AddUrlToCache(shortUrl, completeUrl, ttl)

	_, err := storeService.postgresConnection.Exec(ctx, "INSERT INTO urls(id, short_url, complete_url) VALUES ($1, $2, $3, $4)", id, shortUrl, completeUrl, fallBackUrl)

	if err != nil {
		panic(fmt.Sprintf("Failed saving key url into Postgres | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, completeUrl))
	}

	fmt.Printf("Saved short URL: %s - Complete URL: %s\n", shortUrl, completeUrl)
}

func ReactivateUrl(shortUrl string, ttl string) error {
	var completeUrl string
	err := storeService.postgresConnection.QueryRow(context.Background(), "SELECT complete_url FROM urls where short_url = $1", shortUrl).Scan(&completeUrl)
	if err != nil || completeUrl == "" {
		return fmt.Errorf("Failed getting completed_url from SQL | Error: %v - shortUrl: %s\n", err, shortUrl)
	}
	AddUrlToCache(shortUrl, completeUrl, ttl)
	return nil
}

func DeactivateUrl(shortUrl string) {
	err := storeService.redisClient.Del(ctx, shortUrl).Err()
	if err != nil {
		panic(fmt.Sprintf("Failed deleting key url into Redis | Error: %v - shortUrl: %s\n", err, shortUrl))
	}
}

func AddUrlToCache(shortUrl string, completeUrl string, ttl string) {
	secondsTTL := storeService.cacheDuration
	if ttl != "NA" {
		tmp, err := strconv.Atoi(ttl)
		if err != nil {
			panic(fmt.Sprintf("Failed parsing short URL TTL | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, completeUrl))
		}

		secondsTTL = time.Duration(tmp * 3600000000000) // Convert 1h to nanoseconds
	}

	err := storeService.redisClient.Set(ctx, shortUrl, completeUrl, secondsTTL).Err()
	if err != nil {
		panic(fmt.Sprintf("Failed saving key url into Redis | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, completeUrl))
	}
}

func RetrieveCompleteUrl(shortUrl string) (string, error) {
	result, err := storeService.redisClient.Get(ctx, shortUrl).Result()

	if err != nil {
		log.SetPrefix(fmt.Sprintf("Failed to retrieve complete URL | Error: %v - shortUrl: %s\n", err, shortUrl))
		return "", err
	}

	return result, nil
}

func RetrieveFallbackUrl(shortUrl string) (string, error) {
	var fallbackUrl string
	err := storeService.postgresConnection.QueryRow(context.Background(), "SELECT fallback_url FROM urls where short_url = $1", shortUrl).Scan(&fallbackUrl)
	if err != nil || fallbackUrl == "" {
		return "", fmt.Errorf("Failed getting fallback_url from SQL | Error: %v - shortUrl: %s\n", err, shortUrl)
	}

	return fallbackUrl, nil
}

func GetNextId() int64 {
	var id int64
	err := storeService.postgresConnection.QueryRow(context.Background(), "SELECT NEXTVAL(pg_get_serial_sequence('urls', 'id'))").Scan(&id)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get next ID: %v\n", err)
	}

	return id
}

func UpdateLink(shortUrl string, headers string, sourceIp string) {
	_, err := storeService.postgresConnection.Exec(ctx, "INSERT INTO clicks(short_url, headers, source_ip) VALUES ($1, $2::JSONB, $3)", shortUrl, headers, sourceIp)

	if err != nil {
		panic(fmt.Sprintf("Failed to updated URL | Error: %v - shortUrl: %s\n", err, shortUrl))
	}
}
