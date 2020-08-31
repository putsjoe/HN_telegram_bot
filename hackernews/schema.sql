
CREATE TABLE posts(
    ID          integer primary key,
    Title       text not null,
    Score       integer not null,
    Comments    integer not null,
    URL         text not null,
    read        integer DEFAULT 0
);

CREATE TABLE user(
    ID          integer primary key,
    userID      integer not null,
    postID      integer not null
);