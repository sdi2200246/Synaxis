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
    ('d0000000-0000-0000-0000-000000000003', 'ΟΑΚΑ', 'Κηφισίας 37', 'Maroussi', 'Greece', 38.0368, 23.7875, 70000),
    ('d0000000-0000-0000-0000-000000000004', 'Μέγαρο Μουσικής', 'Βασιλίσσης Σοφίας', 'Athens', 'Greece', 37.9756, 23.7492, 1960),
    ('d0000000-0000-0000-0000-000000000005', 'Τεχνόπολη', 'Πειραιώς 100', 'Athens', 'Greece', 37.9779, 23.7114, 3000),
    ('d0000000-0000-0000-0000-000000000006', 'ΣΕΦ', 'Εθνάρχου Μακαρίου', 'Piraeus', 'Greece', 37.9400, 23.6650, 14000);

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

-- user (password: user123)
INSERT INTO "user" (id, username, password_hash, first_name, last_name, email, phone, address, city, country, tax_id, role, status)
VALUES (
    'a0000000-0000-0000-0000-000000000002',
    'user',
    '$2a$10$nrb/rZ08vT5Eky1bMXxSd.6qTPepKX0YdHqYpA2iE3ZYUyC1qEBbq',
    'John', 'Doe',
    'john@example.com',
    '2109876543',
    'Ermou 10', 'Athens', 'Greece',
    '123456789',
    'user', 'approved'
);

-- user (password: user124)
INSERT INTO "user" (id, username, password_hash, first_name, last_name, email, phone, address, city, country, tax_id, role, status)
VALUES (
    'a0000000-0000-0000-0000-000000000003',
    'pending_user',
    '$2a$10$dEr0Gjk0HO7u9uvxDX0Zlep14G3kOWd9oM6Sh5qWFZpvhqApClt6q',
    'Maria', 'Papadopoulou',
    'maria@example.com',
    '2107654321',
    'Stadiou 15', 'Athens', 'Greece',
    '987654321',
    'user', 'pending'
);

-- ============================================================
-- EVENTS (mix of statuses, venues, dates, organizers)
-- ============================================================

-- Published events (should appear in search)
INSERT INTO event (id, organizer_id, venue_id, title, event_type, status, description, capacity, start_datetime, end_datetime) VALUES
    ('e0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000001', 'Athens Jazz Night', 'Concert', 'PUBLISHED', 'An evening of world-class jazz under the stars at the ancient Odeon.', 500, '2026-06-15 20:00:00', '2026-06-15 23:30:00'),
    ('e0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000002', 'Hamlet - Modern Retelling', 'Play', 'PUBLISHED', 'A bold contemporary take on Shakespeare''s classic tragedy.', 800, '2026-07-10 19:00:00', '2026-07-10 22:00:00'),
    ('e0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000003', 'Athens Marathon Training Camp', 'Training', 'PUBLISHED', 'A weekend bootcamp to prepare runners for the Athens Authentic Marathon.', 200, '2026-05-20 07:00:00', '2026-05-20 14:00:00'),
    ('e0000000-0000-0000-0000-000000000004', 'a0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000004', 'Tech Startups Conference 2026', 'Conference', 'PUBLISHED', 'Leading founders and investors discuss the future of Greek tech.', 1000, '2026-09-05 09:00:00', '2026-09-05 18:00:00'),
    ('e0000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000005', 'Athens Electronic Music Festival', 'Festival', 'PUBLISHED', 'Three days of electronic music across five stages at Technopolis.', 3000, '2026-08-01 18:00:00', '2026-08-03 04:00:00'),
    ('e0000000-0000-0000-0000-000000000006', 'a0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000006', 'Greek Basketball Cup Final', 'Match', 'PUBLISHED', 'The 2026 Greek Basketball Cup Final at the Peace and Friendship Stadium.', 10000, '2026-06-28 19:30:00', '2026-06-28 22:00:00'),
    ('e0000000-0000-0000-0000-000000000007', 'a0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000001', 'Classical Piano Recital', 'Concert', 'PUBLISHED', 'Award-winning pianist performs Chopin and Rachmaninoff.', 400, '2026-07-22 20:30:00', '2026-07-22 22:30:00'),
    ('e0000000-0000-0000-0000-000000000008', 'a0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000004', 'React & Go Workshop', 'Workshop', 'PUBLISHED', 'Hands-on full-stack workshop building a real app with React and Go.', 60, '2026-06-10 10:00:00', '2026-06-10 17:00:00'),
    ('e0000000-0000-0000-0000-000000000009', 'a0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000005', 'Street Food Festival Athens', 'Festival', 'PUBLISHED', 'Taste dishes from 40 vendors across the Mediterranean.', 2000, '2026-05-30 12:00:00', '2026-05-30 23:00:00'),
    ('e0000000-0000-0000-0000-000000000010', 'a0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000006', 'Piraeus Rock Night', 'Concert', 'PUBLISHED', 'Greek and international rock bands live at SEF.', 8000, '2026-07-18 20:00:00', '2026-07-19 01:00:00');

-- Draft events (should NOT appear in search)
INSERT INTO event (id, organizer_id, venue_id, title, event_type, status, description, capacity, start_datetime, end_datetime) VALUES
    ('e0000000-0000-0000-0000-000000000011', 'a0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000002', 'Secret Comedy Show', 'Comedy', 'DRAFT', 'A surprise lineup of Greece''s best comedians.', 300, '2026-08-15 21:00:00', '2026-08-15 23:30:00'),
    ('e0000000-0000-0000-0000-000000000012', 'a0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000003', 'Yoga in the Park', 'Workshop', 'DRAFT', 'Morning yoga sessions at the Olympic park.', 100, '2026-06-01 07:00:00', '2026-06-01 09:00:00');

-- Cancelled event (should NOT appear in search)
INSERT INTO event (id, organizer_id, venue_id, title, event_type, status, description, capacity, start_datetime, end_datetime) VALUES
    ('e0000000-0000-0000-0000-000000000013', 'a0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000001', 'Summer Opera Gala', 'Opera', 'CANCELLED', 'Cancelled due to scheduling conflicts.', 600, '2026-07-01 20:00:00', '2026-07-01 23:00:00');

-- ============================================================
-- EVENT CATEGORIES (M:N)
-- ============================================================
INSERT INTO eventcategory (event_id, category_id) VALUES
    ('e0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000001'),  -- Jazz Night → Music
    ('e0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000006'),  -- Jazz Night → Festival
    ('e0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000003'),  -- Hamlet → Theatre
    ('e0000000-0000-0000-0000-000000000003', 'c0000000-0000-0000-0000-000000000002'),  -- Marathon → Sports
    ('e0000000-0000-0000-0000-000000000004', 'c0000000-0000-0000-0000-000000000004'),  -- Tech Conf → Conference
    ('e0000000-0000-0000-0000-000000000005', 'c0000000-0000-0000-0000-000000000001'),  -- Electronic Fest → Music
    ('e0000000-0000-0000-0000-000000000005', 'c0000000-0000-0000-0000-000000000006'),  -- Electronic Fest → Festival
    ('e0000000-0000-0000-0000-000000000006', 'c0000000-0000-0000-0000-000000000002'),  -- Basketball → Sports
    ('e0000000-0000-0000-0000-000000000007', 'c0000000-0000-0000-0000-000000000001'),  -- Piano Recital → Music
    ('e0000000-0000-0000-0000-000000000008', 'c0000000-0000-0000-0000-000000000005'),  -- React Workshop → Workshop
    ('e0000000-0000-0000-0000-000000000008', 'c0000000-0000-0000-0000-000000000004'),  -- React Workshop → Conference
    ('e0000000-0000-0000-0000-000000000009', 'c0000000-0000-0000-0000-000000000006'),  -- Food Fest → Festival
    ('e0000000-0000-0000-0000-000000000010', 'c0000000-0000-0000-0000-000000000001'), -- Rock Night → Music
    ('e0000000-0000-0000-0000-000000000011', 'c0000000-0000-0000-0000-000000000003'), -- Comedy (draft) → Theatre
    ('e0000000-0000-0000-0000-000000000012', 'c0000000-0000-0000-0000-000000000005'), -- Yoga (draft) → Workshop
    ('e0000000-0000-0000-0000-000000000013', 'c0000000-0000-0000-0000-000000000001'); -- Opera (cancelled) → Music

-- ============================================================
-- TICKET TYPES (varying prices for filter testing)
-- ============================================================
INSERT INTO tickettype (id, event_id, name, price, quantity, available) VALUES
    ('f0000000-0000-0000-0000-000000000001', 'e0000000-0000-0000-0000-000000000001', 'General', 25.00, 300, 300),
    ('f0000000-0000-0000-0000-000000000002', 'e0000000-0000-0000-0000-000000000001', 'VIP', 80.00, 200, 200),
    ('f0000000-0000-0000-0000-000000000003', 'e0000000-0000-0000-0000-000000000002', 'Standard', 35.00, 600, 600),
    ('f0000000-0000-0000-0000-000000000004', 'e0000000-0000-0000-0000-000000000002', 'Premium', 60.00, 200, 200),
    ('f0000000-0000-0000-0000-000000000005', 'e0000000-0000-0000-0000-000000000003', 'Runner', 15.00, 200, 200),
    ('f0000000-0000-0000-0000-000000000006', 'e0000000-0000-0000-0000-000000000004', 'Early Bird', 0.00, 300, 300),
    ('f0000000-0000-0000-0000-000000000007', 'e0000000-0000-0000-0000-000000000004', 'Regular', 50.00, 500, 500),
    ('f0000000-0000-0000-0000-000000000008', 'e0000000-0000-0000-0000-000000000004', 'VIP', 150.00, 200, 200),
    ('f0000000-0000-0000-0000-000000000009', 'e0000000-0000-0000-0000-000000000005', 'Day Pass', 30.00, 1500, 1500),
    ('f0000000-0000-0000-0000-000000000010', 'e0000000-0000-0000-0000-000000000005', '3-Day Pass', 70.00, 1500, 1500),
    ('f0000000-0000-0000-0000-000000000011', 'e0000000-0000-0000-0000-000000000006', 'Upper Tier', 20.00, 6000, 6000),
    ('f0000000-0000-0000-0000-000000000012', 'e0000000-0000-0000-0000-000000000006', 'Courtside', 120.00, 4000, 4000),
    ('f0000000-0000-0000-0000-000000000013', 'e0000000-0000-0000-0000-000000000007', 'Standard', 40.00, 400, 400),
    ('f0000000-0000-0000-0000-000000000014', 'e0000000-0000-0000-0000-000000000008', 'Participant', 100.00, 60, 60),
    ('f0000000-0000-0000-0000-000000000015', 'e0000000-0000-0000-0000-000000000009', 'Entry', 5.00, 2000, 2000),
    ('f0000000-0000-0000-0000-000000000016', 'e0000000-0000-0000-0000-000000000010', 'Standing', 25.00, 5000, 5000),
    ('f0000000-0000-0000-0000-000000000017', 'e0000000-0000-0000-0000-000000000010', 'Seated', 45.00, 3000, 3000);