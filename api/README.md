# API Documentation

This directory contains the OpenAPI specification for the Okra API.

## Files

- `openapi.yaml` - OpenAPI 3.0.3 specification file

## Viewing the API Documentation

### Using Swagger UI

You can view the API documentation using Swagger UI:

1. **Online Swagger Editor**: Visit [https://editor.swagger.io/](https://editor.swagger.io/) and paste the contents of `openapi.yaml`

2. **Local Swagger UI**: Use Docker to run Swagger UI locally:
   ```bash
   docker run -p 8080:8080 -e SWAGGER_JSON=/api/openapi.yaml -v $(pwd)/api:/api swaggerapi/swagger-ui
   ```
   Then visit http://localhost:8080

3. **Using swaggo/swag** (for Go integration): Install and use swag to generate docs:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

### Using Redoc

You can also use Redoc for a beautiful documentation interface:

```bash
npx @redocly/cli preview-docs api/openapi.yaml
```

## Generating Client SDKs

You can generate client SDKs for various languages using [OpenAPI Generator](https://openapi-generator.tech/):

```bash
# Install OpenAPI Generator
npm install @openapitools/openapi-generator-cli -g

# Generate Go client
openapi-generator-cli generate -i api/openapi.yaml -g go -o clients/go

# Generate TypeScript client
openapi-generator-cli generate -i api/openapi.yaml -g typescript-axios -o clients/typescript

# Generate Python client
openapi-generator-cli generate -i api/openapi.yaml -g python -o clients/python
```

## API Endpoints

- `GET /health` - Health check endpoint
- `GET /hello` - Hello world endpoint
- `POST /demo` - Create a demo document in Firestore

## Updating the Specification

When adding new endpoints or modifying existing ones:

1. Update `openapi.yaml` with the new endpoint definition
2. Ensure request/response schemas are properly documented
3. Add examples for better developer experience
4. Update this README if needed

