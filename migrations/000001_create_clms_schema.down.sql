-- Drop indexes 
DROP INDEX IF EXISTS idx_reservation_book_status;
DROP INDEX IF EXISTS idx_fine_member_paid;
DROP INDEX IF EXISTS idx_loans_member_return;
DROP INDEX IF EXISTS idx_copy_barcode;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_books_isbn;

-- Drop tables in reverse order of creation
DROP TABLE IF EXISTS Fine;
DROP TABLE IF EXISTS Loans;
DROP TABLE IF EXISTS Reservation;
DROP TABLE IF EXISTS Copy;
DROP TABLE IF EXISTS BookGenre;
DROP TABLE IF EXISTS BookAuthor;
DROP TABLE IF EXISTS Staff;
DROP TABLE IF EXISTS Member;
DROP TABLE IF EXISTS Branch;
DROP TABLE IF EXISTS Books;
DROP TABLE IF EXISTS Genre;
DROP TABLE IF EXISTS Author;
DROP TABLE IF EXISTS Tokens;
DROP TABLE IF EXISTS Users;