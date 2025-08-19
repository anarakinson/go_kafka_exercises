BEGIN;

DROP INDEX IF EXISTS index_booking_created;
DROP INDEX IF EXISTS index_outbox_created;

DROP TABLE IF EXISTS Outbox;
DROP TABLE IF EXISTS Booking;

COMMIT;