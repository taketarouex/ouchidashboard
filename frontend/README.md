# ouchidashboard frontend

## Install

`npm install`

## Test

`npm run test`

## Run Dev Server

`npm run dev`

## E2E Test

``` bash
docker build -f frontend/Dockerfile.test . -t ouchidashboard_frontend_e2e
docker run -it --cap-add=SYS_ADMIN --rm -v `pwd`:/home/pptruser ouchidashboard_frontend_e2e bash -c "cd frontend; npm run test:e2e"
```
