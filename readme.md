# EduPage2 public backend server
This is a rewrite of the current (private) backend server for EduPage2.
This version is written to be easier to extend and most importantly, faster.
While the core functionality remains the same, this version implements a new routing system, which will make older versions of the EduPage2 app stop working, once this version is fully adopted.

## What does 'public backend' mean?
Public backend in this case means, that the EduPage2 proxy/caching (this) will be publicly available for anyone to use. The server has automatically generated OpenAPI documentation which will be available once we migrate from the old API server.