FROM postgres:18-trixie
LABEL authors="gestalt"

#ENV POSTGRES_USER=appuser
#ENV POSTGRES_PASSWORD=password
#ENV POSTGRES_DB=metadatadb

EXPOSE 5432

COPY init/ /docker-entrypoint-initdb.d/
