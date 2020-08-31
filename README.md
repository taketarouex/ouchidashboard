# ouchi-dashboard

display and store data related to `ouchi` life.

## Description

This collects data related `ouchi` life from `nature remo`
and store them in `Firestore`

## Test

### Unit Test

`make test`

### Integration Test

`docker-compose run --service-ports -d firestore`

`FIRESTORE_EMULATOR_HOST=localhost:8812 GCP_PROJECT="test" FIRESTORE_DOC_PATH="test" make integration_test`

## CI

### test

[github actions](..github/workflows/test.yml)

## CD

[Continuous deployment from Git using Cloud Build](https://cloud.google.com/run/docs/continuous-deployment-with-cloud-build?hl=ja#new-service)

## Todo

- [ ] e2e test
- [ ] bad request
- [ ] code ci from github
