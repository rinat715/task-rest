CREATE TABLE task (
    taskid INTEGER PRIMARY KEY NOT NULL,
    text VARCHAR(255) NOT NULL,
    date Date NOT NULL,
    done BOOLEAN DEFAULT false
);

CREATE TABLE tag (
    tagid INTEGER PRIMARY KEY NOT NULL,
    value VARCHAR(255),
    taskid INT,
    FOREIGN KEY (taskid) REFERENCES task(taskid)
);