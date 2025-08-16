-- migrate:up
-- Create custom type for effect
CREATE TYPE "effect" AS ENUM ('allow', 'deny');

-- Create tables
CREATE TABLE "profiles" (
    "id" varchar(32) NOT NULL,
    "user_id" varchar(32) NOT NULL UNIQUE,
    "first_name" varchar(255) NOT NULL,
    "last_name" varchar(255) NOT NULL,
    "other_names" varchar(255),
    "bio" text,
    "phone" varchar(16),
    "created_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp with time zone,
    PRIMARY KEY ("id")
);

CREATE TABLE "users" (
    "id" varchar(32) NOT NULL,
    "email" varchar(255) NOT NULL UNIQUE,
    "username" varchar(128) UNIQUE,
    "password_hash" varchar(255),
    "is_superuser" boolean NOT NULL DEFAULT false,
    "has_roles" boolean NOT NULL DEFAULT false,
    "is_active" boolean NOT NULL DEFAULT true,
    "is_verified" boolean NOT NULL DEFAULT false,
    "password_reset_requested" boolean NOT NULL DEFAULT false,
    "password_change_on_login" boolean NOT NULL DEFAULT false,
    "suspended_until" timestamp with time zone,
    "banned_at" timestamp with time zone,
    "deactivate_at" timestamp with time zone,
    "created_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp with time zone,
    PRIMARY KEY ("id")
);

CREATE TABLE "actions" (
    "id" varchar(32) NOT NULL,
    "name" varchar(255) NOT NULL UNIQUE,
    "description" text,
    "created_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp with time zone,
    PRIMARY KEY ("id")
);

CREATE TABLE "resources" (
    "id" varchar(32) NOT NULL,
    "name" varchar(255) NOT NULL UNIQUE,
    "description" text,
    "created_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp with time zone,
    PRIMARY KEY ("id")
);

CREATE TABLE "roles" (
    "id" varchar(32) NOT NULL,
    "name" varchar(255) NOT NULL UNIQUE,
    "description" text,
    PRIMARY KEY ("id")
);

CREATE TABLE "roles_users" (
    "id" varchar(32) NOT NULL,
    "role_id" varchar(32) NOT NULL,
    "user_id" varchar(32) NOT NULL,
    "assigned_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

CREATE TABLE "scopes" (
    "id" varchar(32) NOT NULL,
    "name" varchar(255) NOT NULL,
    "description" text,
    "created_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp with time zone,
    PRIMARY KEY ("id")
);

CREATE TABLE "permissions" (
    "id" varchar(32) NOT NULL,
    "name" varchar(255) NOT NULL UNIQUE,
    "description" text,
    "action_id" varchar(32) NOT NULL,
    "resource_id" varchar(32) NOT NULL,
    "scope_id" varchar(32),
    "effect" effect NOT NULL DEFAULT 'allow',
    "created_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp with time zone,
    PRIMARY KEY ("id")
);

CREATE TABLE "roles_permissions" (
    "id" varchar(32) NOT NULL,
    "role_id" varchar(32) NOT NULL,
    "permission_id" varchar(32) NOT NULL,
    "assigned_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "assigned_by" varchar(32),
    PRIMARY KEY ("id")
);

CREATE TABLE "oauth_providers" (
    "id" varchar(32) NOT NULL,
    "name" varchar(255) NOT NULL UNIQUE,
    "auth_url" varchar(1024) NOT NULL,
    "token_url" varchar(1024) NOT NULL,
    "user_info_url" varchar(1024) NOT NULL,
    "is_active" boolean NOT NULL DEFAULT true,
    "created_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp with time zone,
    PRIMARY KEY ("id")
);

CREATE TABLE "oauth_accounts" (
    "id" varchar(32) NOT NULL,
    "user_id" varchar(32) NOT NULL,
    "provider_id" varchar(32) NOT NULL,
    "provider_user_id" varchar(255) NOT NULL,
    "provider_username" varchar(255),
    "provider_email" varchar(255),
    "access_token" text,
    "refresh_token" text,
    "token_expires_at" timestamp with time zone,
    "raw_user_data" jsonb,
    "is_primary" boolean NOT NULL DEFAULT false,
    "created_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp with time zone,
    PRIMARY KEY ("id")
);

CREATE TABLE "users_permissions" (
    "id" varchar(32) NOT NULL,
    "user_id" varchar(32) NOT NULL,
    "permission_id" varchar(32) NOT NULL,
    "assigned_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "assigned_by" varchar(32),
    PRIMARY KEY ("id")
);

CREATE TABLE "audit_logs" (
    "id" varchar(32) NOT NULL,
    "user_id" varchar(32),
    "session_id" varchar(64),
    "action" varchar(255) NOT NULL,
    "resource_type" varchar(255),
    "resource_id" varchar(32),
    "previous_state" jsonb,
    "updated_state" jsonb,
    "ip_address" inet,
    "user_agent" text,
    "request_id" varchar(64),
    "success" boolean NOT NULL DEFAULT false,
    "error_code" varchar(64),
    "error_message" text,
    "metadata" jsonb,
    "created_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

-- Triggers
CREATE TRIGGER set_profiles_updated_at
BEFORE UPDATE ON "profiles"
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_users_updated_at
BEFORE UPDATE ON "users"
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_actions_updated_at
BEFORE UPDATE ON "actions"
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_resources_updated_at
BEFORE UPDATE ON "resources"
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_permissions_updated_at
BEFORE UPDATE ON "permissions"
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_oauth_providers_updated_at
BEFORE UPDATE ON "oauth_providers"
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_oauth_accounts_updated_at
BEFORE UPDATE ON "oauth_accounts"
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_scopes_updated_at
BEFORE UPDATE ON "scopes"
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

-- Indexes
-- Note: Indexes for UNIQUE constraints in the table definition are created automatically.
-- The explicit CREATE UNIQUE INDEX statements are commented out as they are redundant.

-- users
-- CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- CREATE UNIQUE INDEX "idx_users_username" ON "users" ("username");
CREATE INDEX "idx_users_active" ON "users" ("is_active");
CREATE INDEX "idx_users_verified" ON "users" ("is_verified");
CREATE INDEX "idx_users_active_verified" ON "users" ("is_active", "is_verified");

-- actions
-- CREATE UNIQUE INDEX "idx_actions_name" ON "actions" ("name");

-- resources
-- CREATE UNIQUE INDEX "idx_resources_name" ON "resources" ("name");

-- roles
-- CREATE UNIQUE INDEX "idx_roles_name" ON "roles" ("name");

-- roles_users
CREATE UNIQUE INDEX "idx_roles_users_role_id_user_id" ON "roles_users" ("role_id", "user_id");

-- permissions
-- CREATE UNIQUE INDEX "idx_permissions_name" ON "permissions" ("name");
CREATE UNIQUE INDEX "idx_permissions_action_id_resource_id_effect" ON "permissions" ("action_id", "resource_id", "effect");
CREATE UNIQUE INDEX "idx_permissions_action_resource_scope_effect" ON "permissions" ("action_id", "resource_id", "scope_id", "effect");

-- roles_permissions
CREATE UNIQUE INDEX "idx_roles_permissions_role_id_permission_id" ON "roles_permissions" ("role_id", "permission_id");
CREATE INDEX "idx_roles_permissions_lookup" ON "roles_permissions" ("permission_id");

-- oauth_accounts
CREATE UNIQUE INDEX "idx_oauth_accounts_user_id_provider_id" ON "oauth_accounts" ("user_id", "provider_id");
CREATE UNIQUE INDEX "idx_oauth_accounts_provider_id_provider_user_id" ON "oauth_accounts" ("provider_id", "provider_user_id");
CREATE INDEX "idx_oauth_accounts_user_id" ON "oauth_accounts" ("user_id");

-- users_permissions
CREATE UNIQUE INDEX "idx_users_permissions_user_id_permission_id" ON "users_permissions" ("user_id", "permission_id");
CREATE INDEX "idx_users_permissions_lookup" ON "users_permissions" ("permission_id");

-- audit_logs
CREATE INDEX "idx_audit_logs_user_created" ON "audit_logs" ("user_id", "created_at");
CREATE INDEX "idx_audit_logs_resource" ON "audit_logs" ("resource_type", "resource_id");
CREATE INDEX "idx_audit_logs_action_created" ON "audit_logs" ("action", "created_at");
CREATE INDEX "idx_audit_logs_ip_created" ON "audit_logs" ("ip_address", "created_at");
CREATE INDEX "idx_audit_logs_session" ON "audit_logs" ("session_id") WHERE "session_id" IS NOT NULL;

-- Foreign key constraints
ALTER TABLE "profiles" ADD CONSTRAINT "fk_profiles_user_id_users_id" FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE;
ALTER TABLE "roles_users" ADD CONSTRAINT "fk_roles_users_role_id_roles_id" FOREIGN KEY("role_id") REFERENCES "roles"("id") ON DELETE CASCADE;
ALTER TABLE "roles_users" ADD CONSTRAINT "fk_roles_users_user_id_users_id" FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE;
ALTER TABLE "permissions" ADD CONSTRAINT "fk_permissions_action_id_actions_id" FOREIGN KEY("action_id") REFERENCES "actions"("id") ON DELETE CASCADE;
ALTER TABLE "permissions" ADD CONSTRAINT "fk_permissions_resource_id_resources_id" FOREIGN KEY("resource_id") REFERENCES "resources"("id") ON DELETE CASCADE;
ALTER TABLE "permissions" ADD CONSTRAINT "fk_permissions_scope_id_scopes_id" FOREIGN KEY("scope_id") REFERENCES "scopes"("id") ON DELETE SET NULL;
ALTER TABLE "roles_permissions" ADD CONSTRAINT "fk_roles_permissions_role_id_roles_id" FOREIGN KEY("role_id") REFERENCES "roles"("id") ON DELETE CASCADE;
ALTER TABLE "roles_permissions" ADD CONSTRAINT "fk_roles_permissions_permission_id_permissions_id" FOREIGN KEY("permission_id") REFERENCES "permissions"("id") ON DELETE CASCADE;
ALTER TABLE "roles_permissions" ADD CONSTRAINT "fk_roles_permissions_assigned_by_users_id" FOREIGN KEY("assigned_by") REFERENCES "users"("id") ON DELETE SET NULL;
ALTER TABLE "oauth_accounts" ADD CONSTRAINT "fk_oauth_accounts_user_id_users_id" FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE;
ALTER TABLE "oauth_accounts" ADD CONSTRAINT "fk_oauth_accounts_provider_id_oauth_providers_id" FOREIGN KEY("provider_id") REFERENCES "oauth_providers"("id") ON DELETE CASCADE;
ALTER TABLE "users_permissions" ADD CONSTRAINT "fk_users_permissions_user_id_users_id" FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE;
ALTER TABLE "users_permissions" ADD CONSTRAINT "fk_users_permissions_permission_id_permissions_id" FOREIGN KEY("permission_id") REFERENCES "permissions"("id") ON DELETE CASCADE;
ALTER TABLE "users_permissions" ADD CONSTRAINT "fk_users_permissions_assigned_by_users_id" FOREIGN KEY("assigned_by") REFERENCES "users"("id") ON DELETE SET NULL;
ALTER TABLE "audit_logs" ADD CONSTRAINT "fk_audit_logs_user_id_users_id" FOREIGN KEY("user_id") REFERENCES "users"("id") ON DELETE CASCADE;


-- migrate:down
-- Drop foreign key constraints in reverse order of creation
ALTER TABLE "audit_logs" DROP CONSTRAINT "fk_audit_logs_user_id_users_id";
ALTER TABLE "users_permissions" DROP CONSTRAINT "fk_users_permissions_assigned_by_users_id";
ALTER TABLE "users_permissions" DROP CONSTRAINT "fk_users_permissions_permission_id_permissions_id";
ALTER TABLE "users_permissions" DROP CONSTRAINT "fk_users_permissions_user_id_users_id";
ALTER TABLE "oauth_accounts" DROP CONSTRAINT "fk_oauth_accounts_provider_id_oauth_providers_id";
ALTER TABLE "oauth_accounts" DROP CONSTRAINT "fk_oauth_accounts_user_id_users_id";
ALTER TABLE "roles_permissions" DROP CONSTRAINT "fk_roles_permissions_assigned_by_users_id";
ALTER TABLE "roles_permissions" DROP CONSTRAINT "fk_roles_permissions_permission_id_permissions_id";
ALTER TABLE "roles_permissions" DROP CONSTRAINT "fk_roles_permissions_role_id_roles_id";
ALTER TABLE "permissions" DROP CONSTRAINT "fk_permissions_scope_id_scopes_id";
ALTER TABLE "permissions" DROP CONSTRAINT "fk_permissions_resource_id_resources_id";
ALTER TABLE "permissions" DROP CONSTRAINT "fk_permissions_action_id_actions_id";
ALTER TABLE "roles_users" DROP CONSTRAINT "fk_roles_users_user_id_users_id";
ALTER TABLE "roles_users" DROP CONSTRAINT "fk_roles_users_role_id_roles_id";
ALTER TABLE "profiles" DROP CONSTRAINT "fk_profiles_user_id_users_id";

-- Drop triggers
DROP TRIGGER IF EXISTS set_profiles_updated_at ON "profiles";
DROP TRIGGER IF EXISTS set_users_updated_at ON "users";
DROP TRIGGER IF EXISTS set_actions_updated_at ON "actions";
DROP TRIGGER IF EXISTS set_resources_updated_at ON "resources";
DROP TRIGGER IF EXISTS set_permissions_updated_at ON "permissions";
DROP TRIGGER IF EXISTS set_oauth_providers_updated_at ON "oauth_providers";
DROP TRIGGER IF EXISTS set_oauth_accounts_updated_at ON "oauth_accounts";
DROP TRIGGER IF EXISTS set_scopes_updated_at ON "scopes";

-- Drop indexes
DROP INDEX IF EXISTS "idx_audit_logs_session";
DROP INDEX IF EXISTS "idx_audit_logs_ip_created";
DROP INDEX IF EXISTS "idx_audit_logs_action_created";
DROP INDEX IF EXISTS "idx_audit_logs_resource";
DROP INDEX IF EXISTS "idx_audit_logs_user_created";
DROP INDEX IF EXISTS "idx_users_permissions_lookup";
DROP INDEX IF EXISTS "idx_users_permissions_user_id_permission_id";
DROP INDEX IF EXISTS "idx_oauth_accounts_user_id";
DROP INDEX IF EXISTS "idx_oauth_accounts_provider_id_provider_user_id";
DROP INDEX IF EXISTS "idx_oauth_accounts_user_id_provider_id";
DROP INDEX IF EXISTS "idx_roles_permissions_lookup";
DROP INDEX IF EXISTS "idx_roles_permissions_role_id_permission_id";
DROP INDEX IF EXISTS "idx_permissions_action_resource_scope_effect";
DROP INDEX IF EXISTS "idx_permissions_action_id_resource_id_effect";
-- DROP INDEX IF EXISTS "idx_permissions_name"; -- Handled by UNIQUE constraint
DROP INDEX IF EXISTS "idx_roles_users_role_id_user_id";
-- DROP INDEX IF EXISTS "idx_roles_name"; -- Handled by UNIQUE constraint
-- DROP INDEX IF EXISTS "idx_resources_name"; -- Handled by UNIQUE constraint
-- DROP INDEX IF EXISTS "idx_actions_name"; -- Handled by UNIQUE constraint
DROP INDEX IF EXISTS "idx_users_active_verified";
DROP INDEX IF EXISTS "idx_users_verified";
DROP INDEX IF EXISTS "idx_users_active";
-- DROP INDEX IF EXISTS "idx_users_username"; -- Handled by UNIQUE constraint
-- DROP INDEX IF EXISTS "idx_users_email"; -- Handled by UNIQUE constraint

-- Drop tables in reverse order of creation
DROP TABLE IF EXISTS "audit_logs";
DROP TABLE IF EXISTS "users_permissions";
DROP TABLE IF EXISTS "oauth_accounts";
DROP TABLE IF EXISTS "oauth_providers";
DROP TABLE IF EXISTS "roles_permissions";
DROP TABLE IF EXISTS "permissions";
DROP TABLE IF EXISTS "scopes";
DROP TABLE IF EXISTS "roles_users";
DROP TABLE IF EXISTS "roles";
DROP TABLE IF EXISTS "resources";
DROP TABLE IF EXISTS "actions";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "profiles";

-- Drop custom type
DROP TYPE IF EXISTS "effect";
