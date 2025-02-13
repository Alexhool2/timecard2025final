FROM mysql:8.0

COPY mysql-init/ /docker-entrypoint-initdb.d/
RUN chmod +x /docker-entrypoint-initdb.d/init.sh