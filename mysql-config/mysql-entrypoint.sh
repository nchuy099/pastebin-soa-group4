#!/bin/bash
cp /etc/mysql/my.cnf /etc/mysql/conf.d/mysql.cnf
chmod 644 /etc/mysql/conf.d/mysql.cnf
exec docker-entrypoint.sh "$@"
