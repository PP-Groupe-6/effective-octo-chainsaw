#!/usr/bin/env bash
set -eux

if [ "${PGVERSION-}" != "" ]
then
  psql -U postgres -c 'create database prix_banque_test;'
  psql -U postgres -c "create user dev SUPERUSER PASSWORD 'dev';"
  psql -U postgres -d prix_banque_test -f ./scripts/create.sql
  psql -U postgres -d prix_banque_test -f ./scripts/populate.sql
fi