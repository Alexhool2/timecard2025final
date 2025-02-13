#!/bin/bash
set -e

# Se o diretório de dados estiver vazio (primeira inicialização), substitui o template
if [ ! -d "/var/lib/mysql/mysql" ]; then
    echo "Substituindo variáveis no init.sql.template..."
    envsubst < /docker-entrypoint-initdb.d/init.sql.template > /docker-entrypoint-initdb.d/init.sql
fi

# Chama o entrypoint original do MySQL (cujo caminho é /usr/local/bin/docker-entrypoint.sh)
exec /usr/local/bin/docker-entrypoint.sh "$@"