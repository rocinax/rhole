#!/bin/sh

# define version format regex
v_regex="^(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)(-(alpha|beta))?$"

# define version
version="${1}"

if !([[ "${version}" =~ ${v_regex} ]]); then
  echo 1
fi