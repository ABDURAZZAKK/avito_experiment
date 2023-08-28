CREATE TABLE users (
    id         SERIAL        PRIMARY KEY,
    slug       VARCHAR(100)  NOT NULL
);


CREATE TABLE segments (
    slug        VARCHAR(150) PRIMARY KEY
);

CREATE TABLE users_segments (
    user_pk     INT          NOT NULL,
    segment_pk  VARCHAR(150) NOT NULL, 
    PRIMARY KEY (user_pk, segment_pk),
    FOREIGN KEY (user_pk)    REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (segment_pk) REFERENCES segments (slug) ON DELETE CASCADE
);


CREATE TABLE users_segments_stats (
    user_pk     INT          NOT NULL,
    segment_pk  VARCHAR(150) NOT NULL, 
    created_at  TIMESTAMP    NOT NULL,
    operation   VARCHAR(50)  NOT NULL
);