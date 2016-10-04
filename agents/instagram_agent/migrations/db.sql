CREATE TABLE "counters" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "created" DATETIME DEFAULT (datetime('now')),
  "username" CHAR(50) NOT NULL,
  "media" INTEGER NOT NULL,
  "followed_by" INTEGER NOT NULL,
  "follows" INTEGER NOT NULL
);
