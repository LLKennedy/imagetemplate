go test . ./components/... ./internal/filesystem ./render/... ./scaffold/... -cover -coverprofile="coverage.out"; 
# if ($LastExitCode -eq 0) {
	go tool cover -html="coverage.out";
	go tool cover -func="coverage.out";
# }
Remove-Item coverage.out;