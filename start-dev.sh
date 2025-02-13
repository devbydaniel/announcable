docker compose -f compose-base.yml -f compose-dev.yml up -d 
docker compose logs -f compose-base.yml -f compose-dev.yml app -f
