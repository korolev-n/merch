#!/bin/bash

# файл: login.targets
# Подготовим logins.txt заранее
rm -f login.targets

for i in $(seq -w 0 99999); do
  echo "POST http://localhost:8080/api/auth" >> login.targets
  echo "Content-Type: application/json" >> login.targets
  echo "" >> login.targets
  echo "{\"username\": \"user_$i\", \"password\": \"Password123\"}" >> login.targets
  echo "" >> login.targets
done
