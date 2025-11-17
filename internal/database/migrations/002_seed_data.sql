-- FleetPass Seed Data
-- Test data for development

-- Insert test user (password: admin123, bcrypt hash)
INSERT INTO users (id, email, password_hash, role, first_name, last_name, is_active) VALUES
('00000000-0000-0000-0000-000000000001', 'admin@fleetpass.com', '$2a$10$YourBcryptHashHere', 'admin', 'Admin', 'User', true);

-- Insert organizations
INSERT INTO organizations (id, name, slug, is_active) VALUES
('10000000-0000-0000-0000-000000000001', 'Premier Auto Group', 'premier-auto', true),
('10000000-0000-0000-0000-000000000002', 'Elite Motors', 'elite-motors', true),
('10000000-0000-0000-0000-000000000003', 'Luxury Fleet Rentals', 'luxury-fleet', true);

-- Insert locations for Premier Auto Group
INSERT INTO locations (id, organization_id, name, address_line1, city, state, zip_code, country, phone, email, is_active) VALUES
('20000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', 'Downtown Showroom', '123 Main Street', 'Los Angeles', 'CA', '90001', 'USA', '(555) 123-4567', 'downtown@premierauto.com', true),
('20000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000001', 'Airport Location', '456 Airport Blvd', 'Los Angeles', 'CA', '90045', 'USA', '(555) 123-4568', 'airport@premierauto.com', true);

-- Insert locations for Elite Motors
INSERT INTO locations (id, organization_id, name, address_line1, city, state, zip_code, country, phone, email, is_active) VALUES
('20000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000002', 'Elite Motors Beverly Hills', '789 Rodeo Drive', 'Beverly Hills', 'CA', '90210', 'USA', '(555) 234-5678', 'bh@elitemotors.com', true);

-- Insert locations for Luxury Fleet
INSERT INTO locations (id, organization_id, name, address_line1, city, state, zip_code, country, phone, email, is_active) VALUES
('20000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000003', 'Luxury Fleet Santa Monica', '321 Ocean Ave', 'Santa Monica', 'CA', '90401', 'USA', '(555) 345-6789', 'sm@luxuryfleet.com', true);

-- Insert vehicles for Premier Auto Group - Downtown
INSERT INTO vehicles (
    id, organization_id, location_id, vin, make, model, year, trim,
    color_exterior, color_interior, condition, status, mileage,
    body_style, transmission, drivetrain, fuel_type, engine,
    mpg_city, mpg_highway, seats, doors, stock_number,
    description, daily_rate, weekly_rate, monthly_rate,
    features, images
) VALUES
(
    '30000000-0000-0000-0000-000000000001',
    '10000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000001',
    '1HGCM82633A001001',
    'Toyota',
    'Camry',
    2024,
    'XLE Premium',
    'Midnight Black Metallic',
    'Macadamia Leather',
    'new',
    'available',
    15,
    'Sedan',
    '8-Speed Automatic',
    'FWD',
    'Gasoline',
    '2.5L 4-Cylinder',
    28,
    39,
    5,
    4,
    'TOY24001',
    'Experience luxury and efficiency in this brand new 2024 Toyota Camry XLE Premium. Perfect for business or leisure travel.',
    89.99,
    549.99,
    1899.99,
    '["Leather Seats", "Sunroof", "Navigation System", "Backup Camera", "Blind Spot Monitor", "Lane Departure Warning", "Adaptive Cruise Control", "Heated Front Seats", "Wireless Phone Charging", "Apple CarPlay", "Android Auto", "Premium Audio System", "Dual-Zone Climate Control"]'::jsonb,
    '["https://images.unsplash.com/photo-1621007947382-bb3c3994e3fb?w=800", "https://images.unsplash.com/photo-1617654112368-307921291f42?w=800"]'::jsonb
),
(
    '30000000-0000-0000-0000-000000000002',
    '10000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000001',
    '5YFBURHE1HP001002',
    'Toyota',
    'Corolla',
    2024,
    'LE',
    'Super White',
    'Black Fabric',
    'new',
    'available',
    8,
    'Sedan',
    'CVT Automatic',
    'FWD',
    'Gasoline',
    '2.0L 4-Cylinder',
    31,
    40,
    5,
    4,
    'TOY24002',
    'Fuel-efficient and reliable 2024 Toyota Corolla. Perfect for city driving and daily commutes.',
    59.99,
    389.99,
    1299.99,
    '["Backup Camera", "Lane Departure Warning", "Adaptive Cruise Control", "Apple CarPlay", "Android Auto", "Automatic Climate Control", "LED Headlights"]'::jsonb,
    '["https://images.unsplash.com/photo-1623869675781-80aa31bfa4e6?w=800"]'::jsonb
),
(
    '30000000-0000-0000-0000-000000000003',
    '10000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000001',
    '5TDJKRFH8LS001003',
    'Toyota',
    'Highlander',
    2024,
    'XLE AWD',
    'Blueprint',
    'Black Leather',
    'new',
    'available',
    25,
    'SUV',
    '8-Speed Automatic',
    'AWD',
    'Gasoline',
    '3.5L V6',
    21,
    29,
    8,
    4,
    'TOY24003',
    'Spacious and family-friendly 2024 Toyota Highlander with three rows of seating. Perfect for group travel.',
    119.99,
    749.99,
    2599.99,
    '["Third Row Seating", "Leather Seats", "Sunroof", "Navigation System", "Backup Camera", "Blind Spot Monitor", "Power Liftgate", "Heated Seats", "Dual-Zone Climate Control", "Apple CarPlay", "Android Auto", "Premium Audio"]'::jsonb,
    '["https://images.unsplash.com/photo-1609521263047-f8f205293f24?w=800"]'::jsonb
);

-- Insert vehicles for Premier Auto Group - Airport
INSERT INTO vehicles (
    id, organization_id, location_id, vin, make, model, year, trim,
    color_exterior, color_interior, condition, status, mileage,
    body_style, transmission, drivetrain, fuel_type, engine,
    mpg_city, mpg_highway, seats, doors, stock_number,
    description, daily_rate, weekly_rate, monthly_rate,
    features, images
) VALUES
(
    '30000000-0000-0000-0000-000000000004',
    '10000000-0000-0000-0000-000000000001',
    '20000000-0000-0000-0000-000000000002',
    '1N4BL4BV8MN001004',
    'Nissan',
    'Altima',
    2023,
    'SV',
    'Pearl White',
    'Charcoal Cloth',
    'used',
    'available',
    12500,
    'Sedan',
    'CVT Automatic',
    'FWD',
    'Gasoline',
    '2.5L 4-Cylinder',
    28,
    39,
    5,
    4,
    'NIS23001',
    'Well-maintained 2023 Nissan Altima with low mileage. Comfortable and economical choice.',
    54.99,
    349.99,
    1199.99,
    '["Backup Camera", "Blind Spot Warning", "Apple CarPlay", "Android Auto", "Keyless Entry", "Dual-Zone Climate Control"]'::jsonb,
    '["https://images.unsplash.com/photo-1617886322207-c7d580144e87?w=800"]'::jsonb
);

-- Insert vehicles for Elite Motors
INSERT INTO vehicles (
    id, organization_id, location_id, vin, make, model, year, trim,
    color_exterior, color_interior, condition, status, mileage,
    body_style, transmission, drivetrain, fuel_type, engine,
    mpg_city, mpg_highway, seats, doors, stock_number,
    description, daily_rate, weekly_rate, monthly_rate,
    features, images
) VALUES
(
    '30000000-0000-0000-0000-000000000005',
    '10000000-0000-0000-0000-000000000002',
    '20000000-0000-0000-0000-000000000003',
    'WBA5A5C50GG001005',
    'BMW',
    '5 Series',
    2024,
    '540i xDrive',
    'Alpine White',
    'Cognac Vernasca Leather',
    'new',
    'available',
    50,
    'Sedan',
    '8-Speed Automatic',
    'AWD',
    'Gasoline',
    '3.0L Turbo I6',
    23,
    32,
    5,
    4,
    'BMW24001',
    'Luxurious 2024 BMW 5 Series with cutting-edge technology and refined comfort. Make a statement.',
    199.99,
    1299.99,
    4499.99,
    '["Premium Leather", "Panoramic Sunroof", "Head-Up Display", "Harman Kardon Audio", "Wireless Charging", "Gesture Control", "Adaptive LED Headlights", "Massaging Seats", "360 Camera", "Parking Assist", "Apple CarPlay", "Android Auto"]'::jsonb,
    '["https://images.unsplash.com/photo-1555215695-3004980ad54e?w=800"]'::jsonb
),
(
    '30000000-0000-0000-0000-000000000006',
    '10000000-0000-0000-0000-000000000002',
    '20000000-0000-0000-0000-000000000003',
    'WBA5A5C51GG001006',
    'Mercedes-Benz',
    'E-Class',
    2024,
    'E 450 4MATIC',
    'Obsidian Black',
    'Bengal Red/Black Leather',
    'new',
    'available',
    75,
    'Sedan',
    '9-Speed Automatic',
    'AWD',
    'Gasoline',
    '3.0L Turbo I6',
    22,
    31,
    5,
    4,
    'MBZ24001',
    'Sophisticated 2024 Mercedes-Benz E-Class offering the perfect blend of luxury and performance.',
    219.99,
    1449.99,
    4999.99,
    '["Nappa Leather", "Panoramic Roof", "Burmester Audio", "MBUX Interface", "Digital Cockpit", "Wireless Charging", "AMG Styling", "Multi-Contour Seats", "360 Camera", "Active Parking Assist", "Apple CarPlay", "Android Auto"]'::jsonb,
    '["https://images.unsplash.com/photo-1618843479313-40f8afb4b4d8?w=800"]'::jsonb
);

-- Insert vehicles for Luxury Fleet
INSERT INTO vehicles (
    id, organization_id, location_id, vin, make, model, year, trim,
    color_exterior, color_interior, condition, status, mileage,
    body_style, transmission, drivetrain, fuel_type, engine,
    mpg_city, mpg_highway, seats, doors, stock_number,
    description, daily_rate, weekly_rate, monthly_rate,
    features, images
) VALUES
(
    '30000000-0000-0000-0000-000000000007',
    '10000000-0000-0000-0000-000000000003',
    '20000000-0000-0000-0000-000000000004',
    '5YJ3E1EA3MF001007',
    'Tesla',
    'Model 3',
    2024,
    'Long Range AWD',
    'Pearl White Multi-Coat',
    'Black Premium Interior',
    'new',
    'available',
    100,
    'Sedan',
    'Single-Speed Automatic',
    'AWD',
    'Electric',
    'Electric Motor',
    132,
    126,
    5,
    4,
    'TSL24001',
    'Cutting-edge 2024 Tesla Model 3 with exceptional range and performance. Experience the future of driving.',
    149.99,
    999.99,
    3499.99,
    '["Autopilot", "Premium Audio", "Glass Roof", "Wireless Charging", "Heated Seats", "Navigate on Autopilot", "Summon", "Auto Lane Change", "Dashcam", "Sentry Mode", "Mobile App Control", "Over-the-Air Updates"]'::jsonb,
    '["https://images.unsplash.com/photo-1560958089-b8a1929cea89?w=800"]'::jsonb
),
(
    '30000000-0000-0000-0000-000000000008',
    '10000000-0000-0000-0000-000000000003',
    '20000000-0000-0000-0000-000000000004',
    '1C4RJFBG3MC001008',
    'Jeep',
    'Wrangler',
    2024,
    'Rubicon 4xe',
    'Hydro Blue',
    'Black Leather',
    'new',
    'available',
    200,
    'SUV',
    '8-Speed Automatic',
    '4WD',
    'Plug-in Hybrid',
    '2.0L Turbo + Electric',
    20,
    21,
    5,
    4,
    'JEP24001',
    'Adventure-ready 2024 Jeep Wrangler Rubicon 4xe. Go anywhere with hybrid efficiency and legendary capability.',
    139.99,
    899.99,
    3199.99,
    '["Removable Roof", "Removable Doors", "Off-Road Package", "Rock Rails", "Winch", "LED Lighting", "Premium Audio", "Apple CarPlay", "Android Auto", "Forward Collision Warning", "Blind Spot Monitor"]'::jsonb,
    '["https://images.unsplash.com/photo-1606664515524-ed2f786a0bd6?w=800"]'::jsonb
);

-- Add one vehicle in maintenance status
UPDATE vehicles
SET status = 'maintenance'
WHERE id = '30000000-0000-0000-0000-000000000004';

-- Add one vehicle as rented
UPDATE vehicles
SET status = 'rented'
WHERE id = '30000000-0000-0000-0000-000000000002';
