CREATE TABLE IF NOT EXISTS realty (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    address1 TEXT NOT NULL,
    address2 TEXT NOT NULL,
    postal_code TEXT NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    lng DOUBLE PRECISION NOT NULL,
    title TEXT NOT NULL,
    featured_status TEXT NOT NULL,
    city_name TEXT NOT NULL,
    photo_count INTEGER NOT NULL,
    photo_url TEXT NOT NULL,
    raw_property_type TEXT NOT NULL,
    property_type TEXT NOT NULL,
    updated TIMESTAMP NOT NULL,
    rent_range INTEGER[] NOT NULL,
    beds_range INTEGER[] NOT NULL,
    baths_range INTEGER[] NOT NULL,
    dimensions_range INTEGER[] NOT NULL
);

CREATE INDEX idx_city ON realty (city_name);
