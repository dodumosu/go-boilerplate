\restrict NBf5Up9jfSUinL8dNecS40rt5kcmaJez72vwMiS7yVj3jkPzdxbod656YranTnb

-- Dumped from database version 15.14 (Ubuntu 15.14-1.pgdg24.04+1)
-- Dumped by pg_dump version 15.14 (Ubuntu 15.14-1.pgdg24.04+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: fuzzystrmatch; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS fuzzystrmatch WITH SCHEMA public;


--
-- Name: EXTENSION fuzzystrmatch; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION fuzzystrmatch IS 'determine similarities and distance between strings';


--
-- Name: ltree; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS ltree WITH SCHEMA public;


--
-- Name: EXTENSION ltree; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION ltree IS 'data type for hierarchical tree-like structures';


--
-- Name: pg_trgm; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;


--
-- Name: EXTENSION pg_trgm; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pg_trgm IS 'text similarity measurement and index searching based on trigrams';


--
-- Name: effect; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.effect AS ENUM (
    'allow',
    'deny'
);


--
-- Name: trigger_set_timestamp(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.trigger_set_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: actions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.actions (
    id character varying(32) NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


--
-- Name: audit_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.audit_logs (
    id character varying(32) NOT NULL,
    user_id character varying(32),
    session_id character varying(64),
    action character varying(255) NOT NULL,
    resource_type character varying(255),
    resource_id character varying(32),
    previous_state jsonb,
    updated_state jsonb,
    ip_address inet,
    user_agent text,
    request_id character varying(64),
    success boolean DEFAULT false NOT NULL,
    error_code character varying(64),
    error_message text,
    metadata jsonb,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: bans; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.bans (
    id character varying(32) NOT NULL,
    banned_user_id character varying(32) NOT NULL,
    banning_user_id character varying(32) NOT NULL,
    created_at timestamp with time zone NOT NULL
);


--
-- Name: events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.events (
    id character varying(32) NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


--
-- Name: metrics; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.metrics (
    id character varying(32) NOT NULL,
    event_id character varying(32) NOT NULL,
    dimensions jsonb NOT NULL,
    date date NOT NULL,
    count bigint NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


--
-- Name: oauth_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.oauth_accounts (
    id character varying(32) NOT NULL,
    user_id character varying(32) NOT NULL,
    provider_id character varying(32) NOT NULL,
    provider_user_id character varying(255) NOT NULL,
    provider_username character varying(255),
    provider_email character varying(255),
    access_token text,
    refresh_token text,
    token_expires_at timestamp with time zone,
    raw_user_data jsonb,
    is_primary boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


--
-- Name: oauth_providers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.oauth_providers (
    id character varying(32) NOT NULL,
    name character varying(255) NOT NULL,
    auth_url character varying(1024) NOT NULL,
    token_url character varying(1024) NOT NULL,
    user_info_url character varying(1024) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


--
-- Name: permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.permissions (
    id character varying(32) NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    action_id character varying(32) NOT NULL,
    resource_id character varying(32) NOT NULL,
    scope_id character varying(32),
    effect public.effect DEFAULT 'allow'::public.effect NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


--
-- Name: profiles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.profiles (
    id character varying(32) NOT NULL,
    user_id character varying(32) NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    other_names character varying(255),
    bio text,
    phone character varying(16),
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


--
-- Name: resources; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.resources (
    id character varying(32) NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


--
-- Name: roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.roles (
    id character varying(32) NOT NULL,
    name character varying(255) NOT NULL,
    description text
);


--
-- Name: roles_permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.roles_permissions (
    id character varying(32) NOT NULL,
    role_id character varying(32) NOT NULL,
    permission_id character varying(32) NOT NULL,
    assigned_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    assigned_by character varying(32)
);


--
-- Name: roles_users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.roles_users (
    id character varying(32) NOT NULL,
    role_id character varying(32) NOT NULL,
    user_id character varying(32) NOT NULL,
    assigned_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: scopes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.scopes (
    id character varying(32) NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


--
-- Name: suspensions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.suspensions (
    id character varying(32) NOT NULL,
    suspended_user_id character varying(32) NOT NULL,
    suspending_user_id character varying(32) NOT NULL,
    created_at timestamp with time zone NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id character varying(32) NOT NULL,
    email character varying(255) NOT NULL,
    username character varying(128),
    password_hash character varying(255),
    is_superuser boolean DEFAULT false NOT NULL,
    has_roles boolean DEFAULT false NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    is_verified boolean DEFAULT false NOT NULL,
    password_reset_requested boolean DEFAULT false NOT NULL,
    password_change_on_login boolean DEFAULT false NOT NULL,
    suspended_until timestamp with time zone,
    banned_at timestamp with time zone,
    deactivate_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone
);


--
-- Name: users_permissions; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users_permissions (
    id character varying(32) NOT NULL,
    user_id character varying(32) NOT NULL,
    permission_id character varying(32) NOT NULL,
    assigned_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    assigned_by character varying(32)
);


--
-- Name: actions actions_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.actions
    ADD CONSTRAINT actions_name_key UNIQUE (name);


--
-- Name: actions actions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.actions
    ADD CONSTRAINT actions_pkey PRIMARY KEY (id);


--
-- Name: audit_logs audit_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_pkey PRIMARY KEY (id);


--
-- Name: bans bans_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bans
    ADD CONSTRAINT bans_pkey PRIMARY KEY (id);


--
-- Name: events events_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_name_key UNIQUE (name);


--
-- Name: events events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);


--
-- Name: metrics metrics_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.metrics
    ADD CONSTRAINT metrics_pkey PRIMARY KEY (id);


--
-- Name: oauth_accounts oauth_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.oauth_accounts
    ADD CONSTRAINT oauth_accounts_pkey PRIMARY KEY (id);


--
-- Name: oauth_providers oauth_providers_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.oauth_providers
    ADD CONSTRAINT oauth_providers_name_key UNIQUE (name);


--
-- Name: oauth_providers oauth_providers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.oauth_providers
    ADD CONSTRAINT oauth_providers_pkey PRIMARY KEY (id);


--
-- Name: permissions permissions_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_name_key UNIQUE (name);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: profiles profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT profiles_pkey PRIMARY KEY (id);


--
-- Name: profiles profiles_user_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT profiles_user_id_key UNIQUE (user_id);


--
-- Name: resources resources_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_name_key UNIQUE (name);


--
-- Name: resources resources_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.resources
    ADD CONSTRAINT resources_pkey PRIMARY KEY (id);


--
-- Name: roles roles_name_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_name_key UNIQUE (name);


--
-- Name: roles_permissions roles_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles_permissions
    ADD CONSTRAINT roles_permissions_pkey PRIMARY KEY (id);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: roles_users roles_users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles_users
    ADD CONSTRAINT roles_users_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: scopes scopes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.scopes
    ADD CONSTRAINT scopes_pkey PRIMARY KEY (id);


--
-- Name: suspensions suspensions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suspensions
    ADD CONSTRAINT suspensions_pkey PRIMARY KEY (id);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users_permissions users_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users_permissions
    ADD CONSTRAINT users_permissions_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: idx_audit_logs_action_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_audit_logs_action_created ON public.audit_logs USING btree (action, created_at);


--
-- Name: idx_audit_logs_ip_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_audit_logs_ip_created ON public.audit_logs USING btree (ip_address, created_at);


--
-- Name: idx_audit_logs_resource; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_audit_logs_resource ON public.audit_logs USING btree (resource_type, resource_id);


--
-- Name: idx_audit_logs_session; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_audit_logs_session ON public.audit_logs USING btree (session_id) WHERE (session_id IS NOT NULL);


--
-- Name: idx_audit_logs_user_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_audit_logs_user_created ON public.audit_logs USING btree (user_id, created_at);


--
-- Name: idx_bans_banned_banner; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_bans_banned_banner ON public.bans USING btree (banned_user_id, banning_user_id);


--
-- Name: idx_bans_banned_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_bans_banned_user_id ON public.bans USING btree (banned_user_id);


--
-- Name: idx_bans_banning_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_bans_banning_user_id ON public.bans USING btree (banning_user_id);


--
-- Name: idx_bans_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_bans_created_at ON public.bans USING btree (created_at);


--
-- Name: idx_events_name_trgm; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_events_name_trgm ON public.events USING gin (name public.gin_trgm_ops);


--
-- Name: idx_metrics_date; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_metrics_date ON public.metrics USING btree (date);


--
-- Name: idx_metrics_dimensions_gin; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_metrics_dimensions_gin ON public.metrics USING gin (dimensions);


--
-- Name: idx_metrics_event_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_metrics_event_id ON public.metrics USING btree (event_id);


--
-- Name: idx_oauth_accounts_provider_id_provider_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_oauth_accounts_provider_id_provider_user_id ON public.oauth_accounts USING btree (provider_id, provider_user_id);


--
-- Name: idx_oauth_accounts_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_oauth_accounts_user_id ON public.oauth_accounts USING btree (user_id);


--
-- Name: idx_oauth_accounts_user_id_provider_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_oauth_accounts_user_id_provider_id ON public.oauth_accounts USING btree (user_id, provider_id);


--
-- Name: idx_permissions_action_id_resource_id_effect; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_permissions_action_id_resource_id_effect ON public.permissions USING btree (action_id, resource_id, effect);


--
-- Name: idx_permissions_action_resource_scope_effect; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_permissions_action_resource_scope_effect ON public.permissions USING btree (action_id, resource_id, scope_id, effect);


--
-- Name: idx_roles_permissions_lookup; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_roles_permissions_lookup ON public.roles_permissions USING btree (permission_id);


--
-- Name: idx_roles_permissions_role_id_permission_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_roles_permissions_role_id_permission_id ON public.roles_permissions USING btree (role_id, permission_id);


--
-- Name: idx_roles_users_role_id_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_roles_users_role_id_user_id ON public.roles_users USING btree (role_id, user_id);


--
-- Name: idx_suspensions_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_suspensions_created_at ON public.suspensions USING btree (created_at);


--
-- Name: idx_suspensions_suspended_suspender; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_suspensions_suspended_suspender ON public.suspensions USING btree (suspended_user_id, suspending_user_id);


--
-- Name: idx_suspensions_suspended_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_suspensions_suspended_user_id ON public.suspensions USING btree (suspended_user_id);


--
-- Name: idx_suspensions_suspending_user_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_suspensions_suspending_user_id ON public.suspensions USING btree (suspending_user_id);


--
-- Name: idx_users_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_active ON public.users USING btree (is_active);


--
-- Name: idx_users_active_verified; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_active_verified ON public.users USING btree (is_active, is_verified);


--
-- Name: idx_users_permissions_lookup; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_permissions_lookup ON public.users_permissions USING btree (permission_id);


--
-- Name: idx_users_permissions_user_id_permission_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX idx_users_permissions_user_id_permission_id ON public.users_permissions USING btree (user_id, permission_id);


--
-- Name: idx_users_verified; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_users_verified ON public.users USING btree (is_verified);


--
-- Name: actions set_actions_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_actions_updated_at BEFORE UPDATE ON public.actions FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: oauth_accounts set_oauth_accounts_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_oauth_accounts_updated_at BEFORE UPDATE ON public.oauth_accounts FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: oauth_providers set_oauth_providers_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_oauth_providers_updated_at BEFORE UPDATE ON public.oauth_providers FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: permissions set_permissions_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_permissions_updated_at BEFORE UPDATE ON public.permissions FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: profiles set_profiles_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_profiles_updated_at BEFORE UPDATE ON public.profiles FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: resources set_resources_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_resources_updated_at BEFORE UPDATE ON public.resources FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: scopes set_scopes_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_scopes_updated_at BEFORE UPDATE ON public.scopes FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: events set_timestamp; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON public.events FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: metrics set_timestamp; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_timestamp BEFORE UPDATE ON public.metrics FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: users set_users_updated_at; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER set_users_updated_at BEFORE UPDATE ON public.users FOR EACH ROW EXECUTE FUNCTION public.trigger_set_timestamp();


--
-- Name: audit_logs fk_audit_logs_user_id_users_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT fk_audit_logs_user_id_users_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: bans fk_bans_banned_user_id_users_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bans
    ADD CONSTRAINT fk_bans_banned_user_id_users_id FOREIGN KEY (banned_user_id) REFERENCES public.users(id);


--
-- Name: bans fk_bans_banning_user_id_users_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.bans
    ADD CONSTRAINT fk_bans_banning_user_id_users_id FOREIGN KEY (banning_user_id) REFERENCES public.users(id);


--
-- Name: metrics fk_metrics_event_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.metrics
    ADD CONSTRAINT fk_metrics_event_id FOREIGN KEY (event_id) REFERENCES public.events(id) ON DELETE CASCADE;


--
-- Name: oauth_accounts fk_oauth_accounts_provider_id_oauth_providers_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.oauth_accounts
    ADD CONSTRAINT fk_oauth_accounts_provider_id_oauth_providers_id FOREIGN KEY (provider_id) REFERENCES public.oauth_providers(id) ON DELETE CASCADE;


--
-- Name: oauth_accounts fk_oauth_accounts_user_id_users_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.oauth_accounts
    ADD CONSTRAINT fk_oauth_accounts_user_id_users_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: permissions fk_permissions_action_id_actions_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT fk_permissions_action_id_actions_id FOREIGN KEY (action_id) REFERENCES public.actions(id) ON DELETE CASCADE;


--
-- Name: permissions fk_permissions_resource_id_resources_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT fk_permissions_resource_id_resources_id FOREIGN KEY (resource_id) REFERENCES public.resources(id) ON DELETE CASCADE;


--
-- Name: permissions fk_permissions_scope_id_scopes_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT fk_permissions_scope_id_scopes_id FOREIGN KEY (scope_id) REFERENCES public.scopes(id) ON DELETE SET NULL;


--
-- Name: profiles fk_profiles_user_id_users_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.profiles
    ADD CONSTRAINT fk_profiles_user_id_users_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: roles_permissions fk_roles_permissions_assigned_by_users_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles_permissions
    ADD CONSTRAINT fk_roles_permissions_assigned_by_users_id FOREIGN KEY (assigned_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: roles_permissions fk_roles_permissions_permission_id_permissions_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles_permissions
    ADD CONSTRAINT fk_roles_permissions_permission_id_permissions_id FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- Name: roles_permissions fk_roles_permissions_role_id_roles_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles_permissions
    ADD CONSTRAINT fk_roles_permissions_role_id_roles_id FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: roles_users fk_roles_users_role_id_roles_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles_users
    ADD CONSTRAINT fk_roles_users_role_id_roles_id FOREIGN KEY (role_id) REFERENCES public.roles(id) ON DELETE CASCADE;


--
-- Name: roles_users fk_roles_users_user_id_users_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles_users
    ADD CONSTRAINT fk_roles_users_user_id_users_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: suspensions fk_suspensions_suspended_user_id_users_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suspensions
    ADD CONSTRAINT fk_suspensions_suspended_user_id_users_id FOREIGN KEY (suspended_user_id) REFERENCES public.users(id);


--
-- Name: suspensions fk_suspensions_suspending_user_id_users_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.suspensions
    ADD CONSTRAINT fk_suspensions_suspending_user_id_users_id FOREIGN KEY (suspending_user_id) REFERENCES public.users(id);


--
-- Name: users_permissions fk_users_permissions_assigned_by_users_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users_permissions
    ADD CONSTRAINT fk_users_permissions_assigned_by_users_id FOREIGN KEY (assigned_by) REFERENCES public.users(id) ON DELETE SET NULL;


--
-- Name: users_permissions fk_users_permissions_permission_id_permissions_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users_permissions
    ADD CONSTRAINT fk_users_permissions_permission_id_permissions_id FOREIGN KEY (permission_id) REFERENCES public.permissions(id) ON DELETE CASCADE;


--
-- Name: users_permissions fk_users_permissions_user_id_users_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users_permissions
    ADD CONSTRAINT fk_users_permissions_user_id_users_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict NBf5Up9jfSUinL8dNecS40rt5kcmaJez72vwMiS7yVj3jkPzdxbod656YranTnb


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20250807163655'),
    ('20250807180211'),
    ('20250818203312');
