set -e

STORAGE_PATH=$1

FILE=$2
FILE_NAME=$(basename "${FILE}")
FILE_SIZE=$(wc -c <"${FILE}" | tr -d '[:space:]')

curl -s                                   \
  -F "path=${STORAGE_PATH}"               \
  -F "size=${FILE_SIZE}"                  \
  -F "file=@${FILE}"                      \
  -X POST                                 \
  -H "Authorization: Bearer $BLOBS_TOKEN" \
  -H "Content-Type: multipart/form-data"  \
  http://localhost:8080/blobs/ | jq '.'
