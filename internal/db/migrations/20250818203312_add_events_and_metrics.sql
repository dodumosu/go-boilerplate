-- migrate:up
CREATE TABLE "events" (
    "id" varchar(32) NOT NULL,
    "name" varchar(255) NOT NULL UNIQUE,
    "description" text,
    "created_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp with time zone,
    PRIMARY KEY ("id")
);
CREATE INDEX IF NOT EXISTS "idx_events_name_trgm" ON "events" USING gin ("name" gin_trgm_ops);

CREATE TRIGGER set_timestamp BEFORE UPDATE ON "events" FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE "metrics" (
    "id" varchar(32) NOT NULL,
    "event_id" varchar(32) NOT NULL,
    "dimensions" jsonb NOT NULL,
    "date" date NOT NULL,
    "count" bigint NOT NULL,
    "created_at" timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp with time zone,
    PRIMARY KEY ("id")
);
CREATE INDEX IF NOT EXISTS "idx_metrics_dimensions_gin" ON "metrics" USING gin ("dimensions");
CREATE INDEX IF NOT EXISTS "idx_metrics_date" ON "metrics" ("date");
CREATE INDEX IF NOT EXISTS "idx_metrics_event_id" ON "metrics" ("event_id");

CREATE TRIGGER set_timestamp BEFORE UPDATE ON "metrics" FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE "public"."suspensions" (
    "id" character varying(32) NOT NULL,
    "suspended_user_id" character varying(32) NOT NULL,
    "suspending_user_id" character varying(32) NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    PRIMARY KEY ("id")
);

CREATE INDEX "idx_suspensions_suspended_user_id"
ON "public"."suspensions" ("suspended_user_id");
CREATE INDEX "idx_suspensions_suspending_user_id"
ON "public"."suspensions" ("suspending_user_id");
CREATE INDEX "idx_suspensions_created_at"
ON "public"."suspensions" ("created_at");
CREATE UNIQUE INDEX "idx_suspensions_suspended_suspender"
ON "public"."suspensions" ("suspended_user_id", "suspending_user_id");


CREATE TABLE "public"."bans" (
    "id" character varying(32) NOT NULL,
    "banned_user_id" character varying(32) NOT NULL,
    "banning_user_id" character varying(32) NOT NULL,
    "created_at" timestamp with time zone NOT NULL,
    PRIMARY KEY ("id")
);

CREATE INDEX "idx_bans_banned_user_id"
ON "public"."bans" ("banned_user_id");
CREATE INDEX "idx_bans_banning_user_id"
ON "public"."bans" ("banning_user_id");
CREATE INDEX "idx_bans_created_at"
ON "public"."bans" ("created_at");
CREATE UNIQUE INDEX "idx_bans_banned_banner"
ON "public"."bans" ("banned_user_id", "banning_user_id");

ALTER TABLE "metrics"
ADD CONSTRAINT "fk_metrics_event_id" FOREIGN KEY("event_id") REFERENCES "events"("id") ON DELETE CASCADE;

ALTER TABLE "public"."suspensions"
ADD CONSTRAINT "fk_suspensions_suspended_user_id_users_id" FOREIGN KEY("suspended_user_id") REFERENCES "public"."users"("id");

ALTER TABLE "public"."suspensions"
ADD CONSTRAINT "fk_suspensions_suspending_user_id_users_id" FOREIGN KEY("suspending_user_id") REFERENCES "public"."users"("id");

ALTER TABLE "public"."bans"
ADD CONSTRAINT "fk_bans_banned_user_id_users_id" FOREIGN KEY("banned_user_id") REFERENCES "public"."users"("id");

ALTER TABLE "public"."bans"
ADD CONSTRAINT "fk_bans_banning_user_id_users_id" FOREIGN KEY("banning_user_id") REFERENCES "public"."users"("id");

-- migrate:down

DROP TABLE IF EXISTS "bans" CASCADE;
DROP TABLE IF EXISTS "suspensions" CASCADE;
DROP TABLE IF EXISTS "metrics" CASCADE;
DROP TABLE IF EXISTS "events" CASCADE;
