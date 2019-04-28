go test . -cover -coverprofile="coverage.out"; 
go tool cover -html="coverage.out";
Remove-Item coverage.out;