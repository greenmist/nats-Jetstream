# nats-Jetstream
build :- docker compose build
        docker compose build --no-cache
start :- docker compose up

Clean up previous runs :-
docker compose down
docker system prune -f

down :- docker compose down -v
