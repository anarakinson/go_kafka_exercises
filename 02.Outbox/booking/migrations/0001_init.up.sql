BEGIN;

-- таблица booking
CREATE TABLE IF NOT EXISTS Booking (
    booking_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    event_id UUID,
    created_at TIMESTAMP DEFAULT NOW()
);

-- таблица outbox
CREATE TABLE IF NOT EXISTS Outbox (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid() REFERENCES Booking(event_id), 
    event_type VARCHAR(250), 
    payload JSON NOT NULL, 
    processed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()    
);

-- индексы
CREATE INDEX IF NOT EXISTS index_outbox_created ON Outbox(created_at)
CREATE INDEX IF NOT EXISTS index_booking_created ON Booking(created_at)

COMMIT;