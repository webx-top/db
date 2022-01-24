go install github.com/webx-top/db/cmd/dbgenerator@latest
dbgenerator -h 127.0.0.1 -d nging -p root -o dbschema -match "^(nging_vhost_group|nging_vhost|official_film_item|official_film_type)$" -charset utf8mb4
