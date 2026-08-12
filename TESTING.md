# Testing

First be sure to provide your ADS target NetId. You can do that either in the shell or via an env file that the integration tests load automatically.

Using PowerShell:
```pwsh
$env:ADS_TARGET_NET_ID = "199.4.42.250.1.1"
```

Using an env file:
```dotenv
ADS_TARGET_NET_ID=199.4.42.250.1.1
ADS_ROUTER_HOST=127.0.0.1
```

Integration tests look for the first existing file in this order:
- the path in `ADS_TEST_ENV_FILE`
- `test/integration/.env.integration`
- `test/integration/.env`
- `.env.integration` at the repo root
- `.env` at the repo root

- Run all unit tests across the repo (excludes integration build-tag tests):
```pwsh
go test ./...
```

- Run all tests (including tests with the integration build tag):
```pwsh
go test -v -tags=integration ./...
```

- Run all tests in the integration package:
```pwsh
go test -v -tags=integration ./test/integration
```

- Run a single top-level test in the integration package:
```pwsh
go test -v -tags=integration ./test/integration -run '^TestReadAllAccessPaths$'
```

- Run a specific subtest (e.g., sint_max inside TestStaticSeed):
```pwsh
go test -v -tags=integration ./test/integration -run '^TestStaticSeed/sint_max$'
```

- List all tests in the package:
```pwsh
go test -list . -tags=integration ./test/integration
```
