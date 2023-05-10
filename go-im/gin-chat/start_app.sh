echo "--- updating api docs... ---"
./scripts/update_docs.sh
echo "--- starting app... ---"
go run $(pwd)/main.go