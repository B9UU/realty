CREATE TABLE IF NOT EXISTS realty (
    id bigserial PRIMARY KEY,
    listing_type text NOT NULL,
    promo_type text NOT NULL,
    url text NOT NULL,
    project_name text NOT NULL,
    display_address text NOT NULL
);
