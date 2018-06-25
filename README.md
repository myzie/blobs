# Blobs

Datastore for binary objects with properties. Backed by any S3 compatible
object store along with Postgresql for queries and storing object properties.

## API

Each object has a unique path within its operating context. Context is somewhat
similar to an S3 bucket in this way.

Objects may have an associated JSON `properties` object. The properties objects
are currently schemaless however schemas and validations may be added as
options, somewhat like Firebase. Fast searches on values within the JSON are
easily added via Postgresql's secondary indexes on JSONB.

Using a POST request, clients may upload an object and associated properties in
a single request. The object or the properties may be provided individually as
well.

Behind the scenes, objects are stored in S3 in a fixed bucket with a key
corresponding to the specified object `path`. The `path` field cannot be updated
after the object is created. As a future enhancement, consider a `move` or
`copy` operation to switch to a new path assuming it is not already in use.

When the object is downloaded an alternate name can be set via the HTTP header:
`Content-Disposition: attachment; filename=FILENAME`.

Potentially the `name` field on the object, stored in Postgres, can be set
automatically upon download. If `name` is blank, the basename from the `path`
will be used.

Authentication is via JWTs and potentially API keys in the future. The JWT may
optionally include the authorized `context`. Although, if a simple authorization
scheme is implemented here that may not be desired.

It is an external choice how to allocate paths and contexts to best suit the
application at hand.
