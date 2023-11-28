-- Create the seat table
CREATE TABLE seats (
    id SERIAL PRIMARY KEY,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN'
);
