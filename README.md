# blobs

Authenticated object store with a multi-provider object storage backing.

## API

Blobs are stored at a logical path.

## Internal Operating Principles

 * On upload, the `name` field determines the blob name and extension.
   Subsequently, the extension remains fixed but the name can be changed
   including removing the extension.
 * Blob properties are set with a `PUT` after the blob is uploaded or in
   the upload `POST` request.
 * Objects are stored internally at `<bucket>/<blobid>/<blobid>.<ext>`.
   The name is stored in the database and in the `ContentDisposition` field.
   This makes a rename a quick operation.
