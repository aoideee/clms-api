-- 1. Wipe everything and reset IDs to 1
TRUNCATE Fine, Loans, Copy, Books, Member, Branch RESTART IDENTITY CASCADE;

-- 2. Insert the National Heritage Library Branch
INSERT INTO Branch (BranchName, Address, PhoneNumber, Email) 
VALUES ('National Heritage Library', 'Belmopan, Cayo, Belize', '501-822-3300', 'nhl@bnlsis.org');

-- 3. Insert Iconic Belizean Literature
INSERT INTO Books (Title, ISBN, Publisher, PublicationYear, MinimumAge, Description)
VALUES 
('Beka Lamb', '9780435988326', 'Heinemann', 1982, 12, 'A classic story of a young girl growing up in pre-independence Belize.'),
('The Festival of San Joaquin', '9780435989484', 'Heinemann', 1997, 16, 'A powerful novel exploring the struggles and resilience of a woman in rural Belize.'),
('On Heroes, Lizards and Passion', '9789768111104', 'Cubola Productions', 1988, 14, 'An anthology of short stories reflecting the diverse culture of Belize.'),
('Old Ben Tree', '9789768111159', 'Cubola Productions', 2002, 8, 'A beloved Belizean folk tale and children’s story.');

-- 4. Create Physical Copies (All Available)
INSERT INTO Copy (BranchID, BookID, Barcode, Status) 
VALUES 
(1, 1, 'BNL-00001', 'Available'), -- Beka Lamb
(1, 2, 'BNL-00002', 'Available'), -- San Joaquin
(1, 3, 'BNL-00003', 'Available'), -- On Heroes
(1, 4, 'BNL-00004', 'Available'); -- Old Ben Tree

-- 5. Register our "VIP" Member
WITH InsertedUser AS (
    INSERT INTO Users (FirstName, LastName, Email, Role, Activated)
    VALUES ('Evan', 'Hyde', 'evan.hyde@example.com', 'Member', false)
    ON CONFLICT (Email) DO UPDATE SET FirstName = EXCLUDED.FirstName
    RETURNING UserID
)
INSERT INTO Member (UserID, DOB, PhoneNumber, Address, AccountStatus)
SELECT UserID, '1947-04-30', '501-822-1234', 'Belize City', 'Active' 
FROM InsertedUser;