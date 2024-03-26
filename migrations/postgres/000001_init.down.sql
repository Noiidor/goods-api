BEGIN;

DROP TABLE IF EXISTS "public"."goods";

DROP TABLE IF EXISTS "public"."projects";

DROP TABLE IF EXISTS "public"."schema_migrations";

DROP SEQUENCE IF EXISTS "goods_id_seq";
DROP SEQUENCE IF EXISTS "goods_priority_seq";
DROP SEQUENCE IF EXISTS "projects_id_seq";

DROP FUNCTION IF EXISTS "next_priority";

COMMIT;