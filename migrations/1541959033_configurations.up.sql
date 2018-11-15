BEGIN;

-- Table Definition
CREATE TABLE "public"."configs" (
    "id" SERIAL,
    "scheme_id" integer REFERENCES "schemes",
    "created_at" timestamp DEFAULT NOW(),
    "deleted_at" timestamp DEFAULT NULL,
    PRIMARY KEY ("id")
);

CREATE TABLE "public"."config_versions" (
    "config_id" integer REFERENCES "configs" ON DELETE CASCADE,
    "scheme_id" integer REFERENCES "schemes" ON DELETE CASCADE,
    "version" integer DEFAULT 1,
    "tags" jsonb DEFAULT NULL,
    "data" jsonb DEFAULT NULL,
    "created_at" timestamp DEFAULT NOW(),
    PRIMARY KEY ("config_id", "version")
);

-- Index Definition
CREATE INDEX configs__created_at ON public.configs USING btree (created_at);
CREATE INDEX configs__created_at_desc ON public.configs USING btree (created_at DESC);
CREATE INDEX configs__deleted_at ON public.configs USING btree (deleted_at);
CREATE INDEX configs__deleted_at_desc ON public.configs USING btree (deleted_at DESC);

CREATE INDEX config_versions__tags ON public.config_versions USING gin (tags);
CREATE INDEX config_versions__data ON public.config_versions USING gin (data);
CREATE INDEX config_versions__version ON public.config_versions USING btree (version);
CREATE INDEX config_versions__version_desc ON public.config_versions USING btree (version DESC);
CREATE INDEX config_versions__created_at ON public.config_versions USING btree (created_at);
CREATE INDEX config_versions__created_at_desc ON public.config_versions USING btree (created_at DESC);

CREATE INDEX configs__actual ON public.configs USING btree (deleted_at) WHERE deleted_at ISNULL;
CREATE INDEX configs__deleted ON public.configs USING btree (deleted_at) WHERE deleted_at IS NOT NULL;

CREATE INDEX config_versions__tags_version ON public.config_versions USING btree (tags, version);
CREATE INDEX config_versions__tags_version_desc ON public.config_versions USING btree (tags, version DESC);
CREATE INDEX config_versions__tags_version_data ON public.config_versions USING btree (tags, version, data);

COMMIT;