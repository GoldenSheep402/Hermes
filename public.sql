/*
 Navicat Premium Data Transfer

 Source Server         : localhost_5432
 Source Server Type    : PostgreSQL
 Source Server Version : 170002 (170002)
 Source Host           : localhost:5432
 Source Catalog        : hermes
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 170002 (170002)
 File Encoding         : 65001

 Date: 07/04/2025 01:07:06
*/


-- ----------------------------
-- Table structure for binds
-- ----------------------------
DROP TABLE IF EXISTS "public"."binds";
CREATE TABLE "public"."binds" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "uid" text COLLATE "pg_catalog"."default" NOT NULL,
  "platform" text COLLATE "pg_catalog"."default" NOT NULL,
  "open_id" text COLLATE "pg_catalog"."default" NOT NULL,
  "attr" jsonb
)
;
ALTER TABLE "public"."binds" OWNER TO "hermes";

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS "public"."casbin_rule";
CREATE TABLE "public"."casbin_rule" (
  "id" int8 NOT NULL DEFAULT nextval('casbin_rule_id_seq'::regclass),
  "ptype" varchar(100) COLLATE "pg_catalog"."default",
  "v0" varchar(100) COLLATE "pg_catalog"."default",
  "v1" varchar(100) COLLATE "pg_catalog"."default",
  "v2" varchar(100) COLLATE "pg_catalog"."default",
  "v3" varchar(100) COLLATE "pg_catalog"."default",
  "v4" varchar(100) COLLATE "pg_catalog"."default",
  "v5" varchar(100) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."casbin_rule" OWNER TO "hermes";

-- ----------------------------
-- Table structure for categories
-- ----------------------------
DROP TABLE IF EXISTS "public"."categories";
CREATE TABLE "public"."categories" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "name" text COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."categories" OWNER TO "hermes";

-- ----------------------------
-- Table structure for files
-- ----------------------------
DROP TABLE IF EXISTS "public"."files";
CREATE TABLE "public"."files" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "torrent_id" text COLLATE "pg_catalog"."default",
  "length" int8,
  "path" text COLLATE "pg_catalog"."default",
  "path_utf8" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."files" OWNER TO "hermes";

-- ----------------------------
-- Table structure for group_membership_metadata
-- ----------------------------
DROP TABLE IF EXISTS "public"."group_membership_metadata";
CREATE TABLE "public"."group_membership_metadata" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "membership_id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "group_metadata_id" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "key" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "value" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "type" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "description" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "order" int8 NOT NULL
)
;
ALTER TABLE "public"."group_membership_metadata" OWNER TO "hermes";

-- ----------------------------
-- Table structure for group_memberships
-- ----------------------------
DROP TABLE IF EXISTS "public"."group_memberships";
CREATE TABLE "public"."group_memberships" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "uid" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "gid" char(26) COLLATE "pg_catalog"."default" NOT NULL
)
;
ALTER TABLE "public"."group_memberships" OWNER TO "hermes";

-- ----------------------------
-- Table structure for group_metadata
-- ----------------------------
DROP TABLE IF EXISTS "public"."group_metadata";
CREATE TABLE "public"."group_metadata" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "gid" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "key" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "value" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "type" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "description" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "order" int8 NOT NULL
)
;
ALTER TABLE "public"."group_metadata" OWNER TO "hermes";

-- ----------------------------
-- Table structure for groups
-- ----------------------------
DROP TABLE IF EXISTS "public"."groups";
CREATE TABLE "public"."groups" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "name" text COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."groups" OWNER TO "hermes";

-- ----------------------------
-- Table structure for innet_trackers
-- ----------------------------
DROP TABLE IF EXISTS "public"."innet_trackers";
CREATE TABLE "public"."innet_trackers" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "address" text COLLATE "pg_catalog"."default",
  "enable" bool
)
;
ALTER TABLE "public"."innet_trackers" OWNER TO "hermes";

-- ----------------------------
-- Table structure for metadata
-- ----------------------------
DROP TABLE IF EXISTS "public"."metadata";
CREATE TABLE "public"."metadata" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "category_id" text COLLATE "pg_catalog"."default",
  "order" int8,
  "key" text COLLATE "pg_catalog"."default",
  "type" text COLLATE "pg_catalog"."default",
  "description" text COLLATE "pg_catalog"."default",
  "value" text COLLATE "pg_catalog"."default",
  "default_value" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."metadata" OWNER TO "hermes";

-- ----------------------------
-- Table structure for peers
-- ----------------------------
DROP TABLE IF EXISTS "public"."peers";
CREATE TABLE "public"."peers" (
  "peer_id" text COLLATE "pg_catalog"."default",
  "ip" text COLLATE "pg_catalog"."default",
  "port" int8,
  "last_seen" timestamptz(6),
  "status" int8
)
;
ALTER TABLE "public"."peers" OWNER TO "hermes";

-- ----------------------------
-- Table structure for settings
-- ----------------------------
DROP TABLE IF EXISTS "public"."settings";
CREATE TABLE "public"."settings" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "peer_expire_time" int8 DEFAULT 60,
  "smtp_enable" bool DEFAULT false,
  "smtp_host" text COLLATE "pg_catalog"."default",
  "smtp_port" int8,
  "smtp_user" text COLLATE "pg_catalog"."default",
  "smtp_pass" text COLLATE "pg_catalog"."default",
  "smtp_send_name" text COLLATE "pg_catalog"."default",
  "smtp_send_addr" text COLLATE "pg_catalog"."default",
  "register_enable" bool DEFAULT true,
  "login_enable" bool DEFAULT true,
  "publish_enable" bool DEFAULT true
)
;
ALTER TABLE "public"."settings" OWNER TO "hermes";

-- ----------------------------
-- Table structure for single_sums
-- ----------------------------
DROP TABLE IF EXISTS "public"."single_sums";
CREATE TABLE "public"."single_sums" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "uid" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "torrent_id" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "upload" int8 NOT NULL,
  "download" int8 NOT NULL,
  "is_finish" bool NOT NULL
)
;
ALTER TABLE "public"."single_sums" OWNER TO "hermes";

-- ----------------------------
-- Table structure for subnets
-- ----------------------------
DROP TABLE IF EXISTS "public"."subnets";
CREATE TABLE "public"."subnets" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "cidr" text COLLATE "pg_catalog"."default",
  "is_allow" bool
)
;
ALTER TABLE "public"."subnets" OWNER TO "hermes";

-- ----------------------------
-- Table structure for sums
-- ----------------------------
DROP TABLE IF EXISTS "public"."sums";
CREATE TABLE "public"."sums" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "uid" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "real_upload" int8 NOT NULL,
  "real_download" int8 NOT NULL,
  "add_upload" int8 NOT NULL,
  "add_download" int8 NOT NULL
)
;
ALTER TABLE "public"."sums" OWNER TO "hermes";

-- ----------------------------
-- Table structure for torrent_metadata
-- ----------------------------
DROP TABLE IF EXISTS "public"."torrent_metadata";
CREATE TABLE "public"."torrent_metadata" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "torrent_id" text COLLATE "pg_catalog"."default",
  "category_id" text COLLATE "pg_catalog"."default",
  "metadata_id" text COLLATE "pg_catalog"."default",
  "value" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."torrent_metadata" OWNER TO "hermes";

-- ----------------------------
-- Table structure for torrent_statuses
-- ----------------------------
DROP TABLE IF EXISTS "public"."torrent_statuses";
CREATE TABLE "public"."torrent_statuses" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "torrent_id" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "upload_count" int8 NOT NULL DEFAULT 0,
  "upload_sum" int8 NOT NULL DEFAULT 0,
  "download_count" int8 NOT NULL DEFAULT 0,
  "download_sum" int8 NOT NULL DEFAULT 0,
  "seeding_count" int8 NOT NULL DEFAULT 0
)
;
ALTER TABLE "public"."torrent_statuses" OWNER TO "hermes";

-- ----------------------------
-- Table structure for torrents
-- ----------------------------
DROP TABLE IF EXISTS "public"."torrents";
CREATE TABLE "public"."torrents" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "category_id" text COLLATE "pg_catalog"."default",
  "info_hash" text COLLATE "pg_catalog"."default",
  "creator_id" text COLLATE "pg_catalog"."default",
  "is_single_file" bool,
  "announce" text COLLATE "pg_catalog"."default",
  "created_by" text COLLATE "pg_catalog"."default",
  "creation_date" timestamptz(6),
  "comment" text COLLATE "pg_catalog"."default",
  "name" text COLLATE "pg_catalog"."default",
  "name_utf8" text COLLATE "pg_catalog"."default",
  "length" int8,
  "md5sum" text COLLATE "pg_catalog"."default",
  "pieces" bytea,
  "piece_length" int8,
  "private" bool,
  "source" text COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."torrents" OWNER TO "hermes";

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
  "id" char(26) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6),
  "deleted_at" timestamptz(6),
  "name" text COLLATE "pg_catalog"."default",
  "salt" text COLLATE "pg_catalog"."default" NOT NULL,
  "password" text COLLATE "pg_catalog"."default" NOT NULL,
  "limit" int8,
  "key" text COLLATE "pg_catalog"."default",
  "is_admin" bool NOT NULL
)
;
ALTER TABLE "public"."users" OWNER TO "hermes";

-- ----------------------------
-- Indexes structure for table binds
-- ----------------------------
CREATE INDEX "idx_binds_deleted_at" ON "public"."binds" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_binds_open_id" ON "public"."binds" USING btree (
  "open_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);
CREATE INDEX "idx_binds_uid" ON "public"."binds" USING btree (
  "uid" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table binds
-- ----------------------------
ALTER TABLE "public"."binds" ADD CONSTRAINT "uni_binds_open_id" UNIQUE ("open_id");

-- ----------------------------
-- Primary Key structure for table binds
-- ----------------------------
ALTER TABLE "public"."binds" ADD CONSTRAINT "binds_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table casbin_rule
-- ----------------------------
CREATE UNIQUE INDEX "idx_casbin_rule" ON "public"."casbin_rule" USING btree (
  "ptype" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v0" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v1" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v2" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v3" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v4" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST,
  "v5" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table casbin_rule
-- ----------------------------
ALTER TABLE "public"."casbin_rule" ADD CONSTRAINT "casbin_rule_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table categories
-- ----------------------------
CREATE INDEX "idx_categories_deleted_at" ON "public"."categories" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table categories
-- ----------------------------
ALTER TABLE "public"."categories" ADD CONSTRAINT "categories_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table files
-- ----------------------------
CREATE INDEX "idx_files_deleted_at" ON "public"."files" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table files
-- ----------------------------
ALTER TABLE "public"."files" ADD CONSTRAINT "files_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table group_membership_metadata
-- ----------------------------
CREATE INDEX "idx_group_membership_metadata_deleted_at" ON "public"."group_membership_metadata" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table group_membership_metadata
-- ----------------------------
ALTER TABLE "public"."group_membership_metadata" ADD CONSTRAINT "group_membership_metadata_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table group_memberships
-- ----------------------------
CREATE INDEX "idx_group_memberships_deleted_at" ON "public"."group_memberships" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table group_memberships
-- ----------------------------
ALTER TABLE "public"."group_memberships" ADD CONSTRAINT "group_memberships_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table group_metadata
-- ----------------------------
CREATE INDEX "idx_group_metadata_deleted_at" ON "public"."group_metadata" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table group_metadata
-- ----------------------------
ALTER TABLE "public"."group_metadata" ADD CONSTRAINT "group_metadata_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table groups
-- ----------------------------
CREATE INDEX "idx_groups_deleted_at" ON "public"."groups" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table groups
-- ----------------------------
ALTER TABLE "public"."groups" ADD CONSTRAINT "groups_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table innet_trackers
-- ----------------------------
CREATE INDEX "idx_innet_trackers_deleted_at" ON "public"."innet_trackers" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table innet_trackers
-- ----------------------------
ALTER TABLE "public"."innet_trackers" ADD CONSTRAINT "innet_trackers_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table metadata
-- ----------------------------
CREATE INDEX "idx_metadata_deleted_at" ON "public"."metadata" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table metadata
-- ----------------------------
ALTER TABLE "public"."metadata" ADD CONSTRAINT "metadata_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table settings
-- ----------------------------
CREATE INDEX "idx_settings_deleted_at" ON "public"."settings" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table settings
-- ----------------------------
ALTER TABLE "public"."settings" ADD CONSTRAINT "settings_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table single_sums
-- ----------------------------
CREATE INDEX "idx_single_sums_deleted_at" ON "public"."single_sums" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table single_sums
-- ----------------------------
ALTER TABLE "public"."single_sums" ADD CONSTRAINT "single_sums_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table subnets
-- ----------------------------
CREATE INDEX "idx_subnets_deleted_at" ON "public"."subnets" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table subnets
-- ----------------------------
ALTER TABLE "public"."subnets" ADD CONSTRAINT "subnets_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table sums
-- ----------------------------
CREATE INDEX "idx_sums_deleted_at" ON "public"."sums" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table sums
-- ----------------------------
ALTER TABLE "public"."sums" ADD CONSTRAINT "sums_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table torrent_metadata
-- ----------------------------
CREATE INDEX "idx_torrent_metadata_deleted_at" ON "public"."torrent_metadata" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table torrent_metadata
-- ----------------------------
ALTER TABLE "public"."torrent_metadata" ADD CONSTRAINT "torrent_metadata_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table torrent_statuses
-- ----------------------------
CREATE INDEX "idx_torrent_statuses_deleted_at" ON "public"."torrent_statuses" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "torrent_id" ON "public"."torrent_statuses" USING btree (
  "torrent_id" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table torrent_statuses
-- ----------------------------
ALTER TABLE "public"."torrent_statuses" ADD CONSTRAINT "torrent_statuses_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table torrents
-- ----------------------------
CREATE INDEX "idx_torrents_deleted_at" ON "public"."torrents" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table torrents
-- ----------------------------
ALTER TABLE "public"."torrents" ADD CONSTRAINT "torrents_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table users
-- ----------------------------
CREATE INDEX "idx_users_deleted_at" ON "public"."users" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");
