/*
Quoting an identifier also makes it case-sensitive, whereas unquoted names are always folded to lower case.
For example, the identifiers FOO, foo, and "foo" are considered the same by PostgreSQL, but "Foo" and "FOO" are different from these three and each other. 
(The folding of unquoted names to lower case in PostgreSQL is incompatible with the SQL standard, which says that unquoted names should be folded to upper case. 
Thus, foo should be equivalent to "FOO" not "foo" according to the standard. 
If you want to write portable applications you are advised to always quote a particular name or never quote it.)
*/
CREATE TABLE IF NOT EXISTS Person
(
    pid      INT PRIMARY KEY,
    Name     VARCHAR,
    Email    VARCHAR UNIQUE,
    Password VARCHAR
);
CREATE TABLE IF NOT EXISTS Member
(
    SavedTags VARCHAR UNIQUE,
    Help      BOOLEAN
) INHERITS (Person);

CREATE TABLE IF NOT EXISTS Club
(
    Bio  VARCHAR,
    Help BOOLEAN,
    Size INT,
    Tags VARCHAR
) INHERITS (Person);

CREATE TABLE IF NOT EXISTS Chat
(
    cid      INT PRIMARY KEY,
    DateTime timestamp,
    Log      VARCHAR
);

CREATE TABLE IF NOT EXISTS Event
(
    eid         INT PRIMARY KEY,
    DateTime    timestamp,
    Description VARCHAR,
    Location    VARCHAR,
    Fee         FLOAT
);

CREATE TABLE IF NOT EXISTS SwippedRight
(
    userId INT,
    clubId INT,
    CONSTRAINT User_FK FOREIGN KEY (userId) REFERENCES Member (pid),
    CONSTRAINT Club_FK FOREIGN KEY (clubId) REFERENCES Club (pid),
    PRIMARY KEY (userId, clubId)
);

CREATE TABLE IF NOT EXISTS UChats
(
    userId INT,
    chatId INT,
    CONSTRAINT User_FK FOREIGN KEY (userId) REFERENCES Member (pid),
    CONSTRAINT Chat_FK FOREIGN KEY (ChatId) REFERENCES Chat (cid),
    PRIMARY KEY (userId, chatId)
);

CREATE TABLE IF NOT EXISTS CChats
(
    clubId INT,
    chatId INT,
    CONSTRAINT Club_FK FOREIGN KEY (clubId) REFERENCES Club (pid),
    CONSTRAINT Chat_FK FOREIGN KEY (ChatId) REFERENCES Chat (cid),
    PRIMARY KEY (clubId, chatId)
);

CREATE TABLE IF NOT EXISTS Attend
(
    userId  INT,
    eventId INT,
    CONSTRAINT User_FK FOREIGN KEY (userId) REFERENCES Member (pid),
    CONSTRAINT Event_FK FOREIGN KEY (eventId) REFERENCES Event (eid),
    PRIMARY KEY (userId, eventId)
);

CREATE TABLE IF NOT EXISTS Host
(
    clubId  INT,
    eventId INT,
    CONSTRAINT Club_FK FOREIGN KEY (ClubId) REFERENCES Club (pid),
    CONSTRAINT Event_FK FOREIGN KEY (eventId) REFERENCES Event (eid),
    PRIMARY KEY (clubId, eventId)
);