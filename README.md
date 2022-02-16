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

## GET `/:short_url_id`

Redirects the user to the original URL

# SQL Structure

_A SQL script to create the tables can be found under the `assets` directory in the root of this repository_
![SQL Diagram](assets/sql_diagram.png)
