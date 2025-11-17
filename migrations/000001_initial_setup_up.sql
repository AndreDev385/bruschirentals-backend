-- Create neighborhoods table
CREATE TABLE neighborhoods (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

-- Create buildings table
CREATE TABLE buildings (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    neighborhood_id UUID NOT NULL REFERENCES neighborhoods(id) ON DELETE CASCADE,
    address TEXT NOT NULL
);

-- Create index on neighborhood_id for better query performance
CREATE INDEX idx_buildings_neighborhood_id ON buildings(neighborhood_id);