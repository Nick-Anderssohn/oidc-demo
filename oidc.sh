#!/usr/bin/env bash

install_goose() {
    echo "Installing Goose..."
    go install github.com/pressly/goose/v3/cmd/goose@v3.24.1
    if [ $? -eq 0 ]; then
        echo "Goose installed successfully."
    else
        echo "Failed to install Goose."
        exit 1
    fi
}

install_sqlc() {
    echo "Installing sqlc..."
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.29.0
    if [ $? -eq 0 ]; then
        echo "sqlc installed successfully."
    else
        echo "Failed to install sqlc."
        exit 1
    fi
}

first_time_setup() {
    install_goose
    install_sqlc
}

create_sql_migration() {
    cd db/goose/migrations
    goose create $1 sql
    cd ../../..
}

db_migrate() {
    echo "Running database migrations..."
    cd db/goose
    goose up
    cd ../..
}

sqlcgen() {
    cd internal/sqlc
    sqlc generate
    cd ../..
}

db_dump_schema() {
    docker exec -e PGPASSWORD='demo' db pg_dump -s -U demo -d demo > db/schema.sql
    echo "dumped schema to db/schema.sql"
}

db_reset() {
    echo "Resetting database..."
    docker compose down db
    docker volume rm oidc-demo_db_data
    docker compose up -d db

    # Wait a few seconds for the database to start
    echo "Waiting for the database to start..."
    sleep 3

    db_migrate

    echo "DB reset complete"
}

db_start() {
    echo "Starting database using Docker Compose..."
    docker compose up -d db
    if [ $? -eq 0 ]; then
        echo "Database started successfully."
    else
        echo "Failed to start the database."
        exit 1
    fi
}


db_stop() {
    echo "Stopping database using Docker Compose..."
    docker-compose stop db
    if [ $? -eq 0 ]; then
        echo "Database stopped successfully."
    else
        echo "Failed to stop the database."
        exit 1
    fi
}

build_frontend() {
    cd frontend/oidc-demo
    npm run build
    cd ../..
    rm -rf cmd/server/static
    mkdir cmd/server/static
    cp -r frontend/oidc-demo/dist/* cmd/server/static/
}

run_api() {
    build_frontend

    docker compose up --build api
}

run() {
    db_start

    # Wait a few seconds for the database to start
    echo "Waiting for the database to start..."
    sleep 3

    db_migrate
    run_api
}

case $1 in
    first_time_setup)
        first_time_setup
        ;;
    create_sql_migration)
        create_sql_migration $2
        ;;
    db_migrate)
        db_migrate
        ;;
    sqlcgen)
        sqlcgen
        ;;
    db_dump_schema)
        db_dump_schema
        ;;
    db_reset)
        db_reset
        ;;
    db_start)
        db_start
        ;;
    db_stop)
        db_stop
        ;;
    build_frontend)
        build_frontend
        ;;
    run_api)
        run_api
        ;;
    run)
        run
        ;;
    *)
        echo "Usage: $0 {first_time_setup|create_sql_migration|db_migrate|sqlcgen|db_dump_schema|db_reset|db_start|db_stop|build_frontend|run_api|run}"
        exit 1
esac