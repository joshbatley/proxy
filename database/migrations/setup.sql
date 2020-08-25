CREATE TABLE IF NOT EXISTS collection(
  id INTEGER NOT NULL PRIMARY KEY,
  friendlyname TEXT NOT NULL,
  UNIQUE(friendlyname)
);

INSERT OR IGNORE INTO collection (friendlyname) VALUES ("DEFAULT");

CREATE TABLE IF NOT EXISTS cache(
  id INTEGER NOT NULL PRIMARY KEY,
  collection INTEGER NOT NULL,
  url TEXT NOT NULL,
  headers TEXT,
  body BLOB,
  status INTEGER,
  method TEXT,
  datetime INTEGER,
  FOREIGN KEY(collection) REFERENCES collection(id)
);

CREATE TABLE IF NOT EXISTS rules(
  id INTEGER NOT NULL PRIMARY KEY,
  collection INTEGER NOT NULL,
  pattern TEXT NOT NULL,
  cache INTEGER,
  expiry INTEGER,
  offlinecache INTEGER,
  FOREIGN KEY(collection) REFERENCES collection(id)
);
