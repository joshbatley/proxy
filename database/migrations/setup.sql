CREATE TABLE IF NOT EXISTS collection(
  id INTEGER NOT NULL PRIMARY KEY,
  name TEXT NOT NULL,
  UNIQUE(name)
);

INSERT OR IGNORE INTO collection (name) VALUES ("DEFAULT");

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
  allowCors INTEGER
  FOREIGN KEY(collection) REFERENCES collection(id)
);
