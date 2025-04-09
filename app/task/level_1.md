The repository should contain the following

- [x] README.md file explaining the purpose of the repo, along with local setup instructions.

- [x] Explicitly maintaining dependencies in a file ex (pom.xml, build.gradle, go.mod, requirements.txt, etc).

- [x] Makefile to build and run the REST API locally.

- [x] Ability to run DB schema migrations to create the student table.

> Note since I used the docker-compose for the db I used the volume mounted feature to do db migration

- [x] Config (such as database URL) should not be hard-coded in the code and should be passed through environment variables.

- [ ] Postman collection for the APIs.

### API expectations

- [x] Support API versioning (e.g., api/v1/<resource>).

- [x] Using proper HTTP verbs for different operations.

- [x] API should emit meaningful logs with appropriate log levels.

> Used Logrus for the logging 

- [x] API should have a /healthcheck endpoint.

- [ ] Unit tests for different endpoints.