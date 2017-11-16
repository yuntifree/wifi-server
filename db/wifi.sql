use wifi;

CREATE TABLE IF NOT EXISTS users
(
    uid     int unsigned NOT NULL AUTO_INCREMENT,
    username    varchar(36) NOT NULL,
    token       varchar(36) NOT NULL DEFAULT '',
    ctime       datetime NOT NULL DEFAULT '2017-11-01',
    -- 对应wifi_account
    wifi        int unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY(uid),
    UNIQUE KEY(username)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS wifi_account
(
    id      int unsigned NOT NULL AUTO_INCREMENT,
    phone   varchar(16) NOT NULL,
    -- 0x1:试用  0x2:开通
    bitmap  int unsigned NOT NULL DEFAULT 0,
    ctime   datetime NOT NULL DEFAULT '2017-11-01',
    etime   datetime NOT NULL DEFAULT '2017-11-01',
    PRIMARY KEY(id),
    UNIQUE KEY(phone),
    KEY(etime)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS park
(
    id      int unsigned NOT NULL AUTO_INCREMENT,
    name    varchar(128) NOT NULL DEFAULT '',
    address varchar(512) NOT NULL DEFAULT '',
    online  tinyint unsigned NOT NULL DEFAULT 0,
    deleted tinyint unsigned NOT NULL DEFAULT 0,
    ctime   datetime NOT NULL DEFAULT '2017-11-01',
    PRIMARY KEY(id)
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

CREATE TABLE IF NOT EXISTS phone_code
(
    id      bigint unsigned NOT NULL AUTO_INCREMENT,
    phone   varchar(16) NOT NULL,
    uid     int unsigned NOT NULL DEFAULT 0,
    code    int unsigned NOT NULL DEFAULT 0,
    used    tinyint unsigned NOT NULL DEFAULT 0,
    ctime   datetime NOT NULL DEFAULT '2017-11-01',
    stime   datetime NOT NULL DEFAULT '2017-11-01',
    etime   datetime NOT NULL DEFAULT '2017-11-01',
    PRIMARY KEY(id),
    KEY(phone),
    KEY(ctime)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS trial
(
    id      int unsigned NOT NULL AUTO_INCREMENT,
    wid     int unsigned NOT NULL,
    ctime   datetime NOT NULL DEFAULT '2017-11-01',
    etime   datetime NOT NULL DEFAULT '2017-11-01',
    PRIMARY KEY(id),
    UNIQUE KEY(wid)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS items
(
    id      int unsigned NOT NULL AUTO_INCREMENT,
    title   varchar(256) NOT NULL,
    price   int unsigned NOT NULL DEFAULT 0,
    online  tinyint unsigned NOT NULL DEFAULT 0,
    deleted tinyint unsigned NOT NULL DEFAULT 0,
    ctime   datetime NOT NULL DEFAULT '2017-11-01',
    PRIMARY KEY(id)
) ENGINE = InnoDB;

CREATE TABLE IF NOT EXISTS orders
(
    id      bigint unsigned NOT NULL AUTO_INCREMENT,
    oid     varchar(64) NOT NULL,
    uid     int unsigned NOT NULL DEFAULT 0,
    wid     int unsigned NOT NULL DEFAULT 0,
    item    int unsigned NOT NULL DEFAULT 0,
    price   int unsigned NOT NULL DEFAULT 0,
    fee     int unsigned NOT NULL DEFAULT 0,
    prepayid    varchar(64) NOT NULL DEFAULT '',
    ctime   datetime NOT NULL DEFAULT '2017-11-01',
    ftime   datetime NOT NULL DEFAULT '2017-11-01',
    -- status 0:未支付 1:支付成功
    status  tinyint unsigned NOT NULL DEFAULT 0,
    PRIMARY KEY(id),
    UNIQUE KEY(oid)
) ENGINE = InnoDB;
