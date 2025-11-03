FROM postgres:15-alpine
LABEL maintainer="ritwik"

# create empty init dir in image (no host files required)
RUN mkdir -p /docker-entrypoint-initdb.d/

EXPOSE 5432