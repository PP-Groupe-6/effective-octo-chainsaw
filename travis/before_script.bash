#!/usr/bin/env bash
set -eux

if [ "${PGVERSION-}" != "" ]
then
  # The tricky test user, below, has to actually exist so that it can be used in a test
  # of aclitem formatting. It turns out aclitems cannot contain non-existing users/roles.
  psql -U postgres -c 'create database prix_banque_test;'
  psql -U postgres -c "create user dev SUPERUSER PASSWORD 'dev';"
  psql -U postgres -d prix_banque_test -f ./scripts/create.sql
  psql -U postgres -d prix_banque_test -f ./scripts/populate.sql
fi