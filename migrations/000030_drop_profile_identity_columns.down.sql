ALTER TABLE campus_profiles ADD COLUMN full_name VARCHAR(150) DEFAULT '';
ALTER TABLE campus_profiles ADD COLUMN email VARCHAR(255);
CREATE UNIQUE INDEX campus_profiles_email_key ON campus_profiles (email);
