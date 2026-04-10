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
('Old Ben Tree', '9789768111159', 'Cubola Productions', 2002, 8, 'A beloved Belizean folk tale and children’s story.'),
('Time and the River', '9780435989569', 'Heinemann', 2007, 14, 'A historical novel focused on a young Belizean woman grappling with slavery and the struggles of early 19th-century Belizean society.'),
('Pataki Full: Seven Belizean Short Stories', '9789768111128', 'Cubola Productions', 1991, 12, 'A collection of short stories by Sir Colville Young capturing the rich oral traditions, humor, and cultural nuances of Belize.'),
('In Times Like These', '9780435989279', 'Heinemann', 1991, 15, 'A novel set during the turbulent period of Belize''s path to independence, dealing with political and personal choices.'),
('Characters & Caricatures in Belizean Folklore', '9789768111166', 'Belize UNESCO Commission', 1991, 8, 'A deep dive into the fascinating characters of Belizean folklore, like Tata Duende, Sisimito, and La Llorona.'),
('Shots from the Heartland', '9789768111401', 'Cubola Productions', 2007, 16, 'A thought-provoking collection of essays by Evan X Hyde exploring the social and political landscape of Belize.'),
('Warlords and Maize Men', '9789768111227', 'Cubola Productions', 1999, 12, 'A comprehensive and fascinating guide to the major Maya archaeological sites scattered across Belize.'),
('Ping Wing Juk Me', '9789768111050', 'Cubola Productions', 1998, 14, 'Six Belizean plays exploring different facets of everyday life, history, and culture in the country.');

-- 4. Create Physical Copies (All Available)
INSERT INTO Copy (BranchID, BookID, Barcode, Status) 
VALUES 
(1, 1, 'BNL-00001', 'Available'), -- Beka Lamb
(1, 2, 'BNL-00002', 'Available'), -- San Joaquin
(1, 3, 'BNL-00003', 'Available'), -- On Heroes
(1, 4, 'BNL-00004', 'Available'), -- Old Ben Tree
(1, 5, 'BNL-00005', 'Available'), -- Time and the River
(1, 6, 'BNL-00006', 'Available'), -- Pataki Full
(1, 7, 'BNL-00007', 'Available'), -- In Times Like These
(1, 8, 'BNL-00008', 'Available'), -- Characters & Caricatures
(1, 9, 'BNL-00009', 'Available'), -- Shots from the Heartland
(1, 10, 'BNL-00010', 'Available'), -- Warlords and Maize Men
(1, 11, 'BNL-00011', 'Available'); -- Ping Wing Juk Me

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