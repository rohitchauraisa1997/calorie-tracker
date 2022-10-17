# Calorie-tracker app using react as frontend and gin as the backend server.
## To Setup the Calorie Tracker in your local environment.
### Docker compose to bring the calorie tracker app up
> docker compose up -d

### to access into different containers/services of docker compose 
> docker compose ps
### If all services are up and running well the app comes up at localhost:3000 so access
> http://localhost:3000
### to access into different containers/services of docker compose 

> docker compose exec mongodb mongosh calorie-tracker-db -u calorie-user -p calorie123

> docker compose exec backend sh

> docker compose exec frontend sh