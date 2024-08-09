CREATE TABLE IF NOT EXISTS realty (
    id BIGSERIAL PRIMARY KEY,
    name TEXT,
    address1 TEXT,
    address2 TEXT,
    postal_code TEXT,
    lat DOUBLE PRECISION,
    lng DOUBLE PRECISION,
    title TEXT,
    featured_status TEXT,
    city_name TEXT,
    photo_count INTEGER,
    photo_url TEXT,
    raw_property_type TEXT,
    property_type TEXT,
    updated TIMESTAMP,
    rent_range INTEGER[],
    beds_range INTEGER[],
    baths_range INTEGER[],
    dimensions_range INTEGER[]
);

CREATE INDEX idx_city ON realty (city_name);
