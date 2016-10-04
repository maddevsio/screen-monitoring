CREATE TABLE "counters" (
  "created" DATETIME DEFAULT (datetime('now')),
  "media" integer NOT NULL,
  "followed_by" integer NOT NULL,
  "follows" integer NOT NULL
);
