# postgres/Dockerfile

# 1) Start from the official Postgres 15 image
FROM postgres:15

# 2) Install the Debian package that provides pgvector for Postgres 15
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
         postgresql-15-pgvector \
    && rm -rf /var/lib/apt/lists/*

# 3) Copy our init script so that, on first container startup, the "vector" extension is created
COPY init-vector.sql /docker-entrypoint-initdb.d/

# 4) Expose the default PostgreSQL port
EXPOSE 5432

# (The default entrypoint from `postgres:15` will run any *.sql in /docker-entrypoint-initdb.d/ before starting.)
