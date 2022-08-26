PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS repositories
(
    id VARCHAR(20),
    url VARCHAR(2048) NOT NULL,
    branch VARCHAR(255) NOT NULL,
    polling_interval VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL,
    CONSTRAINT repositories_pk PRIMARY KEY (id),
    CONSTRAINT repositories_url_unique UNIQUE (url)
);

CREATE TABLE IF NOT EXISTS builds
(
    id VARCHAR(20),
    repo_id VARCHAR(20),
    status INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    CONSTRAINT builds_pk PRIMARY KEY (id),
    CONSTRAINT builds_repository_id_fk FOREIGN KEY (repo_id) REFERENCES repositories (id) ON DELETE CASCADE,
    CONSTRAINT builds_status_check CHECK (status IN (0, 1, 2, 3))
);

CREATE TABLE IF NOT EXISTS commits
(
    build_id VARCHAR(20),
    hash VARCHAR(40),
    CONSTRAINT commits_build_id_fk FOREIGN KEY (build_id) REFERENCES builds (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS logs
(
    build_id VARCHAR(20),
    data VARCHAR,
    CONSTRAINT logs_build_id_fk FOREIGN KEY (build_id) REFERENCES builds (id) ON DELETE CASCADE
);
