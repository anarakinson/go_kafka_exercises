BEGIN;

-- таблица outbox
CREATE TABLE IF NOT EXISTS Outbox (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    event_type VARCHAR(250), 
    payload JSON NOT NULL, 
    processed_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT NOW()    
);

-- таблица booking
CREATE TABLE IF NOT EXISTS Booking (
    booking_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    event_id UUID REFERENCES Outbox(event_id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- индексы по времени
CREATE INDEX IF NOT EXISTS index_outbox_created ON Outbox(created_at);
CREATE INDEX IF NOT EXISTS index_booking_created ON Booking(created_at);

-- индекс для запроса поиска необработанных событий
CREATE INDEX IF NOT EXISTS index_outbox_processed_created 
ON Outbox(processed_at, created_at) 
WHERE processed_at IS NULL;

COMMIT;