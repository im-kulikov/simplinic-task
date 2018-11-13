BEGIN;

-- Table Definition
CREATE TABLE "public"."configs" (
    "id" SERIAL,
    "version" integer,
    "schemes_id" integer REFERENCES "schemes",
    "tags" jsonb,
    "data" jsonb,
    "created_at" timestamp DEFAULT NOW(),
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

-- Index Definition
CREATE INDEX configs_tags ON public.configs USING gin (tags);
CREATE INDEX configs_data ON public.configs USING gin (data);
CREATE INDEX configs_version ON public.configs USING btree (version);
CREATE INDEX configs_created_at ON public.configs USING btree (created_at);
CREATE INDEX configs_created_at_desc ON public.configs USING btree (created_at DESC);
CREATE INDEX configs_deleted_at ON public.configs USING btree (deleted_at);
CREATE INDEX configs_deleted_at_desc ON public.configs USING btree (deleted_at DESC);

-- search indexes
CREATE INDEX configs_tags_version_schemes_id_actual ON public.configs USING btree (tags, version, schemes_id) WHERE deleted_at ISNULL;
CREATE INDEX configs_tags_version_schemes_id_deleted ON public.configs USING btree (tags, version, schemes_id, data) WHERE deleted_at IS NOT NULL;


COMMIT;