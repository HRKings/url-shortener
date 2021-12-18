package store

import (
	"context"
	"fmt"
	redis "github.com/go-redis/redis/v8"
	pgx "github.com/jackc/pgx/v4"
	"os"
	"strconv"
	"time"
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

	cacheDuration, _ := strconv.Atoi(os.Getenv("REDIS_CACHE"))

	storeService.redisClient = redisClient
	storeService.postgresConnection = postregsConnection
	storeService.cacheDuration = time.Duration(cacheDuration)

	return storeService
}

func SaveUrlMapping(id int, shortUrl string, completeUrl string) {
	err := storeService.redisClient.Set(ctx, shortUrl, completeUrl, storeService.cacheDuration).Err()

	if err != nil {
		panic(fmt.Sprintf("Failed saving key url into Redis | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, completeUrl))
	}

	_, err = storeService.postgresConnection.Exec(ctx, "INSERT INTO urls(id, short_url, complete_url) VALUES ($1, $2, $3)", id, shortUrl, completeUrl)

	if err != nil {
		panic(fmt.Sprintf("Failed saving key url into Postgres | Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, completeUrl))
	}

	fmt.Printf("Saved short URL: %s - Complete URL: %s\n", shortUrl, completeUrl)
}

func RetrieveCompleteUrl(shortUrl string) string {
	result, err := storeService.redisClient.Get(ctx, shortUrl).Result()

	if err != nil {
		panic(fmt.Sprintf("Failed to retrieve complete URL | Error: %v - shortUrl: %s\n", err, shortUrl))
	}

	return result
}

func GetNextId() int {
	var id int
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
