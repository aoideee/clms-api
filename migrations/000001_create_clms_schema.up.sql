-- 1. Independent tables first
CREATE TABLE Author (
    AuthorID SERIAL PRIMARY KEY,
    FirstName VARCHAR(100) NOT NULL,
    LastName VARCHAR(100) NOT NULL,
    Nationality VARCHAR(100),
    Biography TEXT
);

CREATE TABLE Genre (
    GenreID SERIAL PRIMARY KEY,
    GenreName VARCHAR(100) NOT NULL
);

CREATE TABLE Books (
    BookID SERIAL PRIMARY KEY,
    Title VARCHAR(255) NOT NULL,
    ISBN VARCHAR(13) NOT NULL,
    Publisher VARCHAR(150),
    PublicationYear INT,
    MinimumAge INT,
    Description TEXT
);

CREATE TABLE Branch (
    BranchID SERIAL PRIMARY KEY,
    BranchName VARCHAR(100) NOT NULL,
    Address TEXT,
    PhoneNumber VARCHAR(20),
    Email VARCHAR(100)
);

CREATE TABLE Member (
    MemberID SERIAL PRIMARY KEY,
    FirstName VARCHAR(100) NOT NULL,
    LastName VARCHAR(100) NOT NULL,
    DOB DATE NOT NULL,
    PhoneNumber VARCHAR(20),
    Email VARCHAR(100) NOT NULL,
    Address TEXT,
    AccountStatus VARCHAR(20) NOT NULL DEFAULT 'Active'
);

-- 2. Tables with foreign key dependencies
CREATE TABLE BookAuthor (
    AuthorID INT NOT NULL REFERENCES Author(AuthorID) ON DELETE CASCADE,
    BookID INT NOT NULL REFERENCES Books(BookID) ON DELETE CASCADE,
    PRIMARY KEY (AuthorID, BookID)
);

CREATE TABLE BookGenre (
    GenreID INT NOT NULL REFERENCES Genre(GenreID) ON DELETE CASCADE,
    BookID INT NOT NULL REFERENCES Books(BookID) ON DELETE CASCADE,
    PRIMARY KEY (GenreID, BookID)
);

CREATE TABLE Copy (
    CopyID SERIAL PRIMARY KEY,
    BranchID INT NOT NULL REFERENCES Branch(BranchID) ON DELETE CASCADE,
    BookID INT NOT NULL REFERENCES Books(BookID) ON DELETE CASCADE,
    Barcode VARCHAR(50) NOT NULL,
    Status VARCHAR(50) NOT NULL
);

CREATE TABLE Reservation (
    ReservationID SERIAL PRIMARY KEY,
    BookID INT NOT NULL REFERENCES Books(BookID) ON DELETE CASCADE,
    MemberID INT NOT NULL REFERENCES Member(MemberID) ON DELETE CASCADE,
    DateReserved TIMESTAMP NOT NULL DEFAULT NOW(),
    Status VARCHAR(50) NOT NULL
);

-- 3. Transactional tables with multiple dependencies
CREATE TABLE Loans (
    LoanID SERIAL PRIMARY KEY,
    CopyID INT NOT NULL REFERENCES Copy(CopyID) ON DELETE RESTRICT,
    MemberID INT NOT NULL REFERENCES Member(MemberID) ON DELETE CASCADE,
    CheckoutDate TIMESTAMP NOT NULL DEFAULT NOW(),
    DueDate DATE NOT NULL,
    ReturnDate TIMESTAMP
);

CREATE TABLE Fine (
    FineID SERIAL PRIMARY KEY,
    LoanID INT NOT NULL REFERENCES Loans(LoanID) ON DELETE CASCADE,
    MemberID INT NOT NULL REFERENCES Member(MemberID) ON DELETE CASCADE,
    FineType VARCHAR(50) NOT NULL,
    Amount DECIMAL(10,2) NOT NULL,
    PaidStatus BOOLEAN NOT NULL DEFAULT FALSE
);

-- 4. Apply Critical Indexes
-- Unique Indexes for exact lookups
CREATE UNIQUE INDEX idx_books_isbn ON Books(ISBN);
CREATE UNIQUE INDEX idx_members_email ON Member(Email);
CREATE UNIQUE INDEX idx_copy_barcode ON Copy(Barcode);

-- Composite Indexes for filtering logic
CREATE INDEX idx_loans_member_return ON Loans(MemberID, ReturnDate);
CREATE INDEX idx_fine_member_paid ON Fine(MemberID, PaidStatus);

-- Filtered Index for Hold Queue logic
-- In PostgreSQL, a filtered index uses a WHERE clause. Adjust 'Pending' to match your actual application statuses.
CREATE INDEX idx_reservation_book_status ON Reservation(BookID, Status) WHERE Status = 'Pending';