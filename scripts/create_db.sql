CREATE TABLE task (
    taskid INT NOT NULL,
    text VARCHAR(255) NOT NULL,
    date Date NOT NULL,
    done BOOLEAN DEFAULT false,
    PRIMARY KEY (taskid)
);

CREATE TABLE tag (
    tagid INT,
    value VARCHAR(255),
    taskid INT,
    PRIMARY KEY (tagid),
    FOREIGN KEY (taskid) REFERENCES task(taskid)
);