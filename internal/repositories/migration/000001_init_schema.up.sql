PRAGMA foreign_keys = ON;

CREATE TABLE users (
    userid INTEGER PRIMARY KEY NOT NULL,
    email VARCHAR(255) NOT NULL,
    pass VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT false
);

CREATE TABLE tasks (
    taskid INTEGER PRIMARY KEY NOT NULL,
    text VARCHAR(255) NOT NULL,
    date Date NOT NULL,
    done BOOLEAN DEFAULT false,
    userid INTEGER,
    FOREIGN KEY(userid) REFERENCES users(userid)
);

CREATE TABLE tags (
    tagid INTEGER PRIMARY KEY NOT NULL,
    value VARCHAR(255),
    taskid INTEGER,
    FOREIGN KEY(taskid) REFERENCES tasks(taskid)
);