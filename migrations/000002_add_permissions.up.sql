CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);

-- This is the linking table (Many-to-Many relationship)
CREATE TABLE IF NOT EXISTS users_permissions (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
    PRIMARY KEY (user_id, permission_id)
);

-- Insert the default permissions for the CLMS
INSERT INTO permissions (code)
VALUES 
    ('books:read'),     -- Anyone can view the book catalog
    ('books:write'),    -- Librarians/Managers can add or delete books
    ('members:read'),   -- Librarians/Managers can view citizen profiles
    ('members:write'),  -- Librarians/Managers can create/edit citizen profiles
    ('staff:write');    -- Only Managers can create other staff accounts