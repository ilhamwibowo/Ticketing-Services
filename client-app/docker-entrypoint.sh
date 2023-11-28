#!/bin/bash

echo "Making Database Migrations"
python manage.py makemigrations

echo "Apply Database Migrations"
python manage.py migrate

echo "Starting Server"
python manage.py runserver 0.0.0.0:8000