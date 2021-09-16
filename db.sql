-- ----------------------------
-- Sequence structure for tgd2users_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."tgd2users_id_seq";
CREATE SEQUENCE "public"."tgd2users_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for tgd2users
-- ----------------------------
DROP TABLE IF EXISTS "public"."tgd2users";
CREATE TABLE "public"."tgd2users" (
  "id" int8 NOT NULL DEFAULT nextval('tgd2users_id_seq'::regclass),
  "idtelegr" int8 NOT NULL,
  "tgnick" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "d2nick" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."tgd2users_id_seq"
OWNED BY "public"."tgd2users"."id";
SELECT setval('"public"."tgd2users_id_seq"', 7, true);

-- ----------------------------
-- Primary Key structure for table tgd2users
-- ----------------------------
ALTER TABLE "public"."tgd2users" ADD CONSTRAINT "tgd2users_pkey" PRIMARY KEY ("id");
