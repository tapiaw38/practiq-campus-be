-- A grade belongs to the version of the work the teacher actually read.
-- Without a version the teacher could open v1, the student resubmit as v2,
-- and the grade land on v2 — content nobody reviewed.
ALTER TABLE submissions ADD COLUMN version INT NOT NULL DEFAULT 1;
