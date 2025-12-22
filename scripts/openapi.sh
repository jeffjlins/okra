#!/bin/zsh

# Run Swagger UI to view OpenAPI documentation
# Uses port 8081 to avoid conflict with the API server (which runs on 8080)

docker run -p 8081:8080 -e SWAGGER_JSON=/api/openapi.yaml -v $(pwd)/api:/api swaggerapi/swagger-ui