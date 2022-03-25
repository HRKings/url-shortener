# URL Shortener

This is a simple URL shortener written in Go, it features a Redis cache with a 3-day TTL and a click logger with persistent storage in PostgreSQL

# Running

## Shell
Create a `.env` file with the necessary values (with can be found in the example), and execute the `run.sh`

## Docker

Using docker-compose, run the following
```bash
docker-compose -f docker-compose.yml -f docker-compose.run.yml up -d
```

# API Specification

_An OpenAPI specification can be found under the `assets` directory in the root of this repository_

The service only has two endpoints:

## POST `/`

Used to shorten an URL, accepts a JSON body containing a `url` property with the URL to be shortened

Example:
```json
{
  "url": "https://github.com/HRKings/url-shortener",
  "fallback": "https://example.com", // (Optional) Optional fallback URL for when the short link expires
  "ttl": "72" // (Optional) TTL in hours of the link
}
```

## GET `/:short_url`

Redirects the user to the original URL

## PUT `/:short_url`

Add the shortUrl into REDIS (will redirect again)

### PUT `/:short_url?ttl=EXPIRATION_IN_HOURS`

Optionally the query param `ttl` can be sent to set the cache expiration of this link

## DELETE `/:short_url`

Remove the shortUrl from REDIS (will not redirect anymore)

# SQL Structure

_A SQL script to create the tables can be found under the `assets` directory in the root of this repository_
![SQL Diagram](assets/sql_diagram.png)
