CREATE TABLE url (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    short_url TEXT UNIQUE NOT NULL,
    complete_url TEXT NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

CREATE TABLE clicks (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    short_url TEXT NOT NULL,
    headers JSONB,
    clicked_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW(),
    source_ip TEXT,
    CONSTRAINT fk_short_url
        FOREIGN KEY(short_url) REFERENCES url(short_url)
);
