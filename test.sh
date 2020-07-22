curl -i --location --request POST '127.0.0.1:8080/fizzbuzz' \
--header 'Content-Type: application/json' \
--data-raw '{
        "str1": "fizz",
        "limit": 20,
        "int1": 3,
        "int2": 6,
        "str2":"buzz"
}'
