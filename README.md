# gRPC PostgreSQL CRUD Sample

Step 1. Create Proto file, define messages and services

Step 2. Specify appropriate inject tag comments on top of fields for custom gorm tags

Step 3. Generate pg.go file

Step 4. Inject Gorm Tag to skip xxx fields protoc-go-inject-tag -input=./postgrescrud/postgrescrud.pb.go -XXX_skip=gorm

Step 5. Generate Reverse proxy, pb.gw.go file

Step 6. Success!
