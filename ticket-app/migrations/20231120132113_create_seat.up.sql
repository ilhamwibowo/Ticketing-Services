-- Create seats table 
CREATE TABLE seats (
    id SERIAL,
    seat_number VARCHAR(20),
    event_id INT,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    PRIMARY KEY (seat_number, event_id)
);

-- Inserting seats for Concert A
INSERT INTO seats (seat_number, event_id) VALUES 
    ('A1', 1),
    ('A2', 1),
    ('A3', 1),
    ('B1', 1),
    ('B2', 1);

-- Inserting seats for Conference X
INSERT INTO seats (seat_number, event_id) VALUES 
    ('MainHall1', 2),
    ('MainHall2', 2),
    ('MainHall3', 2),
    ('WorkshopA1', 2),
    ('WorkshopA2', 2);

-- Inserting seats for Theater Show B
INSERT INTO seats (seat_number, event_id) VALUES 
    ('FrontRow1', 3),
    ('FrontRow2', 3),
    ('Balcony1', 3),
    ('Balcony2', 3),
    ('Balcony3', 3);