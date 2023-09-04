# Project decisions

- Standard Library: do without known web frameworks like Echo or Gin
  in order to showcase knowledge over the language's standard library;
- API Documentation: use Swagger for its nice UI and to showcase usage
  of a well known tool (in the context);
- Provisioning: use Docker/Docker Compose to facilitate the project's
  portability through virtualization;

# Implemetation decisions

- A simple three layered "Main/Entrypoint -> REST API -> Database"
  architecture sufices for the solution here;
- A way to uniquely identify password cards is required in order to
  edit and delete a specifc card. To keep it simple, the decision was
  to use a small and easy to handle UUID library available at:
  https://pkg.go.dev/github.com/rs/xid;

