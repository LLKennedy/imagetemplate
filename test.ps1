go test ./... -cover -coverprofile="coverage.out"; 
go tool cover -html="coverage.out";
go tool cover -func="coverage.out";
Remove-Item coverage.out;