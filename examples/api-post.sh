set -e
FILE=$1
FILE_NAME=$(basename "${FILE}")
FILE_SIZE=$(wc -c <"${FILE}" | tr -d '[:space:]')
curl -s                                   \
  -F "name=${FILE_NAME}"                  \
  -F "size=${FILE_SIZE}"                  \
  -F "file=@${FILE}"                      \
  -X POST                                 \
  -H "Authorization: Bearer $BLOBS_TOKEN" \
  -H "Content-Type: multipart/form-data"  \
  http://localhost:8080/blobs/ | jq '.'
