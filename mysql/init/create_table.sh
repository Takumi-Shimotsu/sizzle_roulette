#!/bin/sh

CMD_MYSQL="mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE}"
$CMD_MYSQL -e "create table sizzle (
    id int(10)  AUTO_INCREMENT NOT NULL primary key,
    name varchar(255) NOT NULL,
    taste varchar(255),
    texture varchar(255),
    information varchar(255),
    price int(10)
    );"
#$CMD_MYSQL -e  "insert into article values (1, '記事1', '記事1です。');"
#$CMD_MYSQL -e  "insert into article values (2, '記事2', '記事2です。');"