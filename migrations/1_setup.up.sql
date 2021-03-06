CREATE TABLE IF NOT EXISTS Collections(
  ID INTEGER NOT NULL PRIMARY KEY,
  Name TEXT NOT NULL,
  HealthCheckURLs TEXT,
  UNIQUE(Name)
);

INSERT OR IGNORE INTO Collections (Name) VALUES ("DEFAULT");

CREATE TABLE IF NOT EXISTS Responses(
  ID TEXT NOT NULL PRIMARY KEY,
  EndpointID INTEGER NOT NULL,
  URL TEXT NOT NULL,
  Headers TEXT,
  Body BLOB,
  Status INTEGER,
  Method TEXT,
  DateTime INTEGER,
  UNIQUE(Method, Status, EndpointID),
  FOREIGN KEY(EndpointID) REFERENCES Endpoints(ID)
);

CREATE TABLE IF NOT EXISTS Endpoints(
  ID TEXT NOT NULL PRIMARY KEY,
  CollectionID INTEGER NOT NULL,
  PreferedStatus INTEGER NOT NULL,
  URL TEXT NOT NULL,
  Method TEXT NOT NULL,
  UNIQUE(URL, Method, CollectionID),
  FOREIGN KEY(CollectionID) REFERENCES Collections(ID)
);

CREATE TABLE IF NOT EXISTS Rules(
  ID INTEGER NOT NULL PRIMARY KEY,
  CollectionID INTEGER NOT NULL,
  Pattern BLOB NOT NULL,
  SaveResponse INTEGER,
  ForceCors INTEGER,
  Expiry INTEGER,
  SkipOffline INTEGER,
  DelayTime INTEGER,
  RemapRegex BLOB,
  FOREIGN KEY(CollectionID) REFERENCES Collections(ID)
);

INSERT OR IGNORE INTO rules (ID, CollectionID, Pattern, SaveResponse, ForceCors, Expiry, SkipOffline, DelayTime, RemapRegex) VALUES (1, 1, ".*", 1, 1, 3600, 1, 0, "")
