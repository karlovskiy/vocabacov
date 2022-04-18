CREATE TABLE IF NOT EXISTS PHRASES_NEW(
    ID INTEGER PRIMARY KEY,
    LANG TEXT NOT NULL,
    PHRASE TEXT NOT NULL,
    STATUS TEXT DEFAULT 'ACTIVE',
    CREATION_DATE DATETIME DEFAULT CURRENT_TIMESTAMP,
    MODIFICATION_DATE DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO PHRASES_NEW(LANG, PHRASE) select LANG,PHRASE from PHRASES;

DROP TABLE PHRASES;

ALTER TABLE PHRASES_NEW RENAME TO PHRASES;