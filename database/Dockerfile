FROM postgres:12.1-alpine

# Copy the database initialization script:
COPY init.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]