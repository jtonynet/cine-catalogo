# 4. Integration Tests "Happy Path"

Date: 2024-01-14

## Status

Accepted

## Context

The project started with a focus on simplicity, and TDD was not adopted due to excitement in development, aiming for quick initial deliveries and learning. As the project progressed, a set of tests became indispensable for security and automation. We weighed the options between a suite of unit or integration tests and the standard library or third-party library [Testify](https://github.com/stretchr/testify).

Reference material:
- The book [Test-Driven Development in Go](https://www.amazon.com.br/Test-Driven-Development-practical-idiomatic-real-world/dp/1803247878) was used as a guide.
- https://medium.com/nerd-for-tech/testing-rest-api-in-go-with-testify-and-mockery-c31ea2cc88f9
- https://forum.golangbridge.org/t/how-to-test-gin-gonic-handler-function-within-a-function/33334/2

Issues regarding "context passing" bugs in integration tests:
- https://github.com/gin-gonic/gin/issues/1292
- https://github.com/gin-gonic/gin/pull/2803
- https://github.com/gin-gonic/gin/issues/2778
- https://github.com/gin-gonic/gin/issues/2816


## Decision

During the research, it became evident that in a two-tier architecture in such a simple CRUD, the integration testing approach would be the best option, as there is not much specialized logic of reduced size that justifies unit testing. By using integration testing at the moment, we ensure the "happy path" of the API, and we are already preparing everything for future iterations to evolve with error tests and corner cases of the application.

We also decided to use the [Testify](https://github.com/stretchr/testify) library due to its good adoption in the community and features compared to the standard GO library, which, in the version used, lacks certain asserts.

We tested the routes in a scenario closest to real-world integration testing. We did not use the testing context provided by GIN because the standard route approach worked as expected in our tests.

## Consequences

With ample supporting material and widespread community adoption, Testify ensures security in the implementation of integration tests. We will not test the smallest unit of our system because it does not require that level of complexity, and it makes sense to validate the routes. We expect to reduce cognitive load by using the library and increase the level of security in deployments.