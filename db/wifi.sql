use wifi;

CREATE TABLE IF NOT EXISTS users
(
    uid     int unsigned NOT NULL AUTO_INCREMENT,
    username    varchar(36) NOT NULL,
    phone       varchar(16) NOT NULL DEFAULT '',
    token       varchar(36) NOT NULL DEFAULT '',
    ctime       datetime NOT NULL DEFAULT '2017-11-01',
    PRIMARY KEY(uid),
    UNIQUE KEY(username)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS banner
(
    id      bigint unsigned NOT NULL AUTO_INCREMENT,
    img     varchar(128) NOT NULL DEFAULT '',
    dst     varchar(256) NOT NULL DEFAULT '',
    online  tinyint unsigned NOT NULL DEFAULT 0,
    deleted tinyint unsigned NOT NULL DEFAULT 0,
    ctime   datetime NOT NULL DEFAULT '2017-11-01',
    PRIMARY KEY(id),
    KEY(ctime)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS feedback
(
    id      bigint unsigned NOT NULL AUTO_INCREMENT,
    phone   varchar(16) NOT NULL DEFAULT '',
    content varchar(1024) NOT NULL DEFAULT '',
    ctime   datetime NOT NULL DEFAULT '2017-11-01',
    PRIMARY KEY(id),
    KEY(ctime)
) ENGINE = InnoDB;
