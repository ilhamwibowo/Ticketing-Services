-- Create events table
CREATE TABLE events (
    event_id SERIAL PRIMARY KEY,
    event_name VARCHAR(100)
);

-- Inserting events
INSERT INTO events (event_name) VALUES 
    ('Concert A'),
    ('Conference X'),
    ('Theater Show B');