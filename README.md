# Journey2 Events Service

## Build, run and test

Run unit-tests

```
make test
```

Run integration tests (The app needs to be running!)
```
make integration_test
```

Run the project

```
make run
```

Run project with live updates while developing

```
gowatch
```

## Environment Variables

```
JOURNEY2_PORT=:8600
```

## API

```
POST /action

{
  "aid" : "9735965b-e1cb-4d7f-adb9-a4adf457f61a",
  "uid" : "ceb2a540-48c7-40ec-bc22-24ffd54d880d",
  "act" : "act_complete_trial",
  "par" : ""
}
```

Pre-defined actions:
- act_land_on_site - just open the page
- act_complete_trial - minimal interaction that explains what site is about
- act_begin_signup - click on signup link and sees the signup form
- act_complete_signup - completes signup
- act_payment - makes a single payment
