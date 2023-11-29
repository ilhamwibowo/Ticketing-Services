-- Create seats table 
CREATE TABLE seats (
    id SERIAL,
    seat_number VARCHAR(20),
    event_id INT,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    PRIMARY KEY (seat_number, event_id)
);