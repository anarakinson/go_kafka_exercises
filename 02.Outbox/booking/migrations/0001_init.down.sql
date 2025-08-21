BEGIN;

DROP INDEX IF EXISTS index_booking_created;
DROP INDEX IF EXISTS index_outbox_created;
DROP TABLE IF EXISTS index_outbox_processed_created;

DROP TABLE IF EXISTS Booking;
DROP TABLE IF EXISTS Outbox;

COMMIT;