CREATE SCHEMA demo;

ALTER SCHEMA demo OWNER TO demo;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;

CREATE TABLE demo."user" (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    email text NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE demo.identity_provider (
    id text NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);

CREATE TABLE demo.identity (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    identity_provider_id text NOT NULL,
    user_id uuid NOT NULL,
    external_id text NOT NULL,
    created_at timestamp without time zone DEFAULT now(),
    updated_at timestamp without time zone DEFAULT now()
);
