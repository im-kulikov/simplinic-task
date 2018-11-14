BEGIN;

-- Table Definition
CREATE TABLE "public"."schemes" (
    "id" SERIAL,
    "created_at" timestamp DEFAULT NOW(),
    "deleted_at" timestamp DEFAULT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "public".scheme_versions(
    "scheme_id" integer REFERENCES "schemes" ON DELETE CASCADE,
    "version" integer DEFAULT 1,
    "tags" jsonb DEFAULT NULL,
    "data" jsonb DEFAULT NULL,
    "created_at" timestamp DEFAULT NOW(),
    PRIMARY KEY ("scheme_id", "version")
);

-- Index Definition
CREATE INDEX schemes__created_at ON public.schemes USING btree (created_at);
CREATE INDEX schemes__created_at_desc ON public.schemes USING btree (created_at DESC);
CREATE INDEX schemes__deleted_at ON public.schemes USING btree (deleted_at);
CREATE INDEX schemes__deleted_at_desc ON public.schemes USING btree (deleted_at DESC);

CREATE INDEX schemes__actual ON public.schemes USING btree (deleted_at) WHERE deleted_at ISNULL;
CREATE INDEX schemes__deleted ON public.schemes USING btree (deleted_at) WHERE deleted_at IS NOT NULL;

CREATE INDEX scheme_versions__data ON public.scheme_versions USING gin (data);
CREATE INDEX scheme_versions__tags ON public.scheme_versions USING gin (tags);
CREATE INDEX scheme_versions__version ON public.scheme_versions USING btree (version);
CREATE INDEX scheme_versions__created_at ON public.scheme_versions USING btree (created_at);
CREATE INDEX scheme_versions__created_at_desc ON public.scheme_versions USING btree (created_at DESC);

-- search indexes
CREATE INDEX scheme_versions__tags_version ON public.scheme_versions USING btree (tags, version);
CREATE INDEX scheme_versions__tags_version_data ON public.scheme_versions USING btree (tags, version, data);

COMMIT;
