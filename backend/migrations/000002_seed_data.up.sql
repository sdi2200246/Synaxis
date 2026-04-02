-- categories
INSERT INTO category (id, name, parent_id) VALUES
    ('c0000000-0000-0000-0000-000000000001', 'Music', NULL),
    ('c0000000-0000-0000-0000-000000000002', 'Sports', NULL),
    ('c0000000-0000-0000-0000-000000000003', 'Theatre', NULL),
    ('c0000000-0000-0000-0000-000000000004', 'Conference', NULL),
    ('c0000000-0000-0000-0000-000000000005', 'Workshop', NULL),
    ('c0000000-0000-0000-0000-000000000006', 'Festival', NULL);

-- venues
INSERT INTO venue (id, name, address, city, country, latitude, longitude, capacity) VALUES
    ('d0000000-0000-0000-0000-000000000001', 'Ωδείο Ηρώδου Αττικού', 'Διονυσίου Αρεοπαγίτου', 'Athens', 'Greece', 37.9704, 23.7245, 5000),
    ('d0000000-0000-0000-0000-000000000002', 'Θέατρο Παλλάς', 'Βουκουρεστίου 5', 'Athens', 'Greece', 37.9792, 23.7351, 1200),
    ('d0000000-0000-0000-0000-000000000003', 'ΟΑΚΑ', 'Κηφισίας 37', 'Maroussi', 'Greece', 38.0368, 23.7875, 70000);

-- admin user (password: admin123)
INSERT INTO "user" (id, username, password_hash, first_name, last_name, email, phone, address, city, country, tax_id, role, status)
VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'admin',
    '$2a$10$sdSP3MKw.c0SXe.sxQC6HOkt9Tsf47ClWlaq3wa3dt1U9g3aIxB.K',
    'Admin', 'User',
    'admin@synaxis.com',
    '2101234567',
    'Admin St 1', 'Athens', 'Greece',
    '000000000',
    'admin', 'approved'
);