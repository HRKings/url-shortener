CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_url TEXT UNIQUE NOT NULL,
    complete_url TEXT NOT NULL,
    fallback_url TEXT,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

CREATE TABLE clicks (
    id SERIAL PRIMARY KEY,
    short_url TEXT NOT NULL,
    headers JSONB,
    clicked_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    source_ip TEXT,
    CONSTRAINT fk_short_url
        FOREIGN KEY(short_url) REFERENCES urls(short_url)
);
