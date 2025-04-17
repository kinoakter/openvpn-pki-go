CREATE TABLE IF NOT EXISTS ca (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    certificate TEXT NOT NULL,
    private_key TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS client_certificates (
    id UUID PRIMARY KEY,
    common_name TEXT NOT NULL,
    certificate TEXT NOT NULL,
    private_key TEXT NOT NULL,
    tls_crypt_v2_key TEXT,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT now()
);