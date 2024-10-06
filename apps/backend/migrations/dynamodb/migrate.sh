#!/bin/sh
check_table_command="aws dynamodb describe-table --table-name voluntrips --endpoint-url http://localhost:4566"
create_table_command="aws dynamodb create-table --cli-input-json file://migrations/dynamodb/schema.json --endpoint-url http://localhost:4566"
$check_table_command || {
    echo "Table does not exist, creating it."
    $create_table_command
}
