-- Campus tenants are contract-enabled Practiq institutions. The authoritative
-- school and membership records remain in practiq-be; school_id is therefore
-- intentionally not a cross-database foreign key.
CREATE TABLE campus_tenants (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  school_id UUID NOT NULL UNIQUE,
  status VARCHAR(20) NOT NULL DEFAULT 'active'
    CHECK (status IN ('active', 'suspended', 'closed')),
  activated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deactivated_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Kept nullable in this first schema migration because a second migration
-- promotes it to NOT NULL only after every write path has been tenantized.
ALTER TABLE courses ADD COLUMN tenant_id UUID REFERENCES campus_tenants(id);
CREATE INDEX idx_courses_tenant_id ON courses(tenant_id);
ALTER TABLE courses DROP CONSTRAINT IF EXISTS courses_slug_key;
CREATE UNIQUE INDEX courses_tenant_slug_key
  ON courses(tenant_id, slug) WHERE tenant_id IS NOT NULL;

ALTER TABLE calendar_events ADD COLUMN tenant_id UUID REFERENCES campus_tenants(id);
CREATE INDEX idx_calendar_events_tenant_id ON calendar_events(tenant_id);

ALTER TABLE conversations ADD COLUMN tenant_id UUID REFERENCES campus_tenants(id);
CREATE INDEX idx_conversations_tenant_id ON conversations(tenant_id);

ALTER TABLE notifications ADD COLUMN tenant_id UUID REFERENCES campus_tenants(id);
CREATE INDEX idx_notifications_tenant_user_created
  ON notifications(tenant_id, user_id, created_at DESC);
