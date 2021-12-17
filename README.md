# URL Shortener

This is a simple URL shortener written in Go, it features a Redis cache with a 3-day TTL and a click logger with persistent storage in PostgreSQL

# Running

Create a `.env` file with the necessary values (with can be found in the example), and execute the `run.sh`

# SQL Structure

![SQL Diagram](assets/sql_diagram.png)