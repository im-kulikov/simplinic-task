BEGIN;

-- Table Definition
CREATE TABLE "public"."schemes" (
    "id" SERIAL,
    "version" integer,
    "tags" jsonb,
    "data" jsonb,
    "created_at" timestamp DEFAULT NOW(),
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

-- Index Definition
CREATE INDEX schemes_tags ON public.schemes USING gin (tags);
CREATE INDEX schemes_data ON public.schemes USING gin (data);
CREATE INDEX schemes_version ON public.schemes USING btree (version);
CREATE INDEX schemes_created_at ON public.schemes USING btree (created_at);
CREATE INDEX schemes_created_at_desc ON public.schemes USING btree (created_at DESC);
CREATE INDEX schemes_deleted_at ON public.schemes USING btree (deleted_at);
CREATE INDEX schemes_deleted_at_desc ON public.schemes USING btree (deleted_at DESC);

-- search indexes
CREATE INDEX schemes_tags_version_actual ON public.schemes USING btree (tags, version) WHERE deleted_at ISNULL;
CREATE INDEX schemes_tags_version_deleted ON public.schemes USING btree (tags, version, data) WHERE deleted_at ISNULL;

COMMIT;
