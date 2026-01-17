-- Roles
INSERT INTO roles (name, description) VALUES
('admin', 'Administrator with full access'),
('staff', 'Normal staff member');

--Branches
INSERT INTO branches (email, name, coordinator_name, contact_number, established_on, aashram_area, created_by) 
VALUES ('delhi_branch@djjs.org', 'Delhi Ashram', 'Swami Vivek', '9876543210', '2005-04-12', 1500.75, 'system'), 
('mumbai_branch@djjs.org', 'Mumbai Ashram', 'Swami Anand', '9123456789', '2010-09-30', 2200.50, 'system');

--Areas
INSERT INTO areas (branch_id, district_id, district_coverage, area_name, area_coverage, created_by, updated_by) 
VALUES (1, gen_random_uuid(), 75.5, 'North Zone', 60.2, 'admin', 'system'), 
(2, gen_random_uuid(), 85.0, 'South Zone', 70.8, 'system', 'system');

--Users
INSERT INTO users (name, email, contact_number, password, role_id, is_deleted, created_on, created_by) 
VALUES 
('Admin User', 'admin@example.com', '9876543210', '$2a$10$oKmRjWC3/0ESYAkmOcNdCeHXJryFuDAF3HuTp5vVElsUERgJf7956', 1, FALSE, NOW(), 'system'),
('Staff User', 'staff1@example.com', '9123456780', '$2a$10$TwNFNphjv/1s4WwXPl0l8udRyTdbEM5p32yfKXeW3U4kE6klY6tTa', 2, FALSE, NOW(), 'system'),

--Countries
INSERT INTO countries (name) VALUES ('India'), ('USA');

--States
INSERT INTO states (name, country_id) VALUES
('Andhra Pradesh', 1),
('Arunachal Pradesh', 1),
('Assam', 1),
('Bihar', 1),
('Chhattisgarh', 1),
('Goa', 1),
('Gujarat', 1),
('Haryana', 1),
('Himachal Pradesh', 1),
('Jharkhand', 1),
('Karnataka', 1),
('Kerala', 1),
('Madhya Pradesh', 1),
('Maharashtra', 1),
('Manipur', 1),
('Meghalaya', 1),
('Mizoram', 1),
('Nagaland', 1),
('Odisha', 1),
('Punjab', 1),
('Rajasthan', 1),
('Sikkim', 1),
('Tamil Nadu', 1),
('Telangana', 1),
('Tripura', 1),
('Uttar Pradesh', 1),
('Uttarakhand', 1),
('West Bengal', 1),
('Andaman and Nicobar Islands', 1),
('Chandigarh', 1),
('Dadra and Nagar Haveli and Daman and Diu', 1),
('Delhi', 1),
('Jammu and Kashmir', 1),
('Ladakh', 1),
('Lakshadweep', 1),
('Puducherry', 1),
('New York',2);

-- Insert values into media_coverage_type
INSERT INTO media_coverage_type (media_type) VALUES ('Social'), ('News Paper');

-- Insert values into promotion_material_type
INSERT INTO promotion_material_type (material_type) VALUES ('Flex'), ('Online Ads');

-- -- Made dummy event for foreign key relation
-- INSERT INTO event_details (event_name) VALUES ('Ram Katha'), ('Peace Program'), ('Krishna Janmastmi');

-- Event Types
INSERT INTO event_types (name) VALUES ('Spiritual'),('Cultural'),('Peace Procession'),('Peace Assembly'),('Fixed Program'),('Others');

-- Event Categories (FK → event_types)
INSERT INTO event_categories (name, event_type_id) VALUES ('Ram Katha',1),('Krishna Katha',1),('Bhajan Sandhya',1),('Meditation Event',1);

INSERT INTO districts (id, name, state_id, country_id) VALUES
(1, 'Mathura', 1, 1),
(2, 'Agra', 1, 1),
(3, 'Mumbai Suburban', 2, 1),
(4, 'Pune', 2, 1),
(5, 'Bangalore Urban', 3, 1),
(6, 'Mysore', 3, 1),
(7, 'Lucknow', 1, 1),
(8, 'Nagpur', 2, 1),
(9, 'Gulbarga', 3, 1),
(10, 'Varanasi', 1, 1);

INSERT INTO cities (id, name, state_id) VALUES
(1, 'Mathura City', 1),
(2, 'Agra City', 1),
(3, 'Mumbai', 2),
(4, 'Pune City', 2),
(5, 'Bangalore', 3),
(6, 'Mysore City', 3),
(7, 'Lucknow City', 1),
(8, 'Nagpur City', 2),
(9, 'Gulbarga City', 3),
(10, 'Varanasi City', 1);

-- Insert Orators master data
INSERT INTO orators (name, created_by) VALUES
('Swami Ji', 'system'),
('Swami Vivek', 'system'),
('Swami Anand', 'system'),
('Swami Prem', 'system'),
('Swami Shanti', 'system'),
('Dr. Anya Sharma', 'system'),
('Dr. Rohan Verma', 'system'),
('Swami Krishna', 'system'),
('Swami Radha', 'system'),
('Swami Gyan', 'system');

-- Insert Languages master data
INSERT INTO languages (name, code, created_by) VALUES
('Hindi', 'hi', 'system'),
('English', 'en', 'system'),
('Sanskrit', 'sa', 'system'),
('Gujarati', 'gu', 'system'),
('Marathi', 'mr', 'system'),
('Bengali', 'bn', 'system'),
('Tamil', 'ta', 'system'),
('Telugu', 'te', 'system'),
('Kannada', 'kn', 'system'),
('Malayalam', 'ml', 'system'),
('Punjabi', 'pa', 'system'),
('Urdu', 'ur', 'system'),
('Odia', 'or', 'system'),
('Assamese', 'as', 'system');

-- Insert Prefixes master data
INSERT INTO prefixes (name, created_by) VALUES
('Mr.', 'system'),
('Mrs.', 'system'),
('Ms.', 'system'),
('Dr.', 'system'),
('Prof.', 'system'),
('Shri', 'system'),
('Smt.', 'system'),
('Shrimati', 'system'),
('Kumari', 'system'),
('Swami', 'system');