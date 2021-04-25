gen:
	GOOS=linux GOARCH=amd64 go build -o SC_report main.go

upload:
	scp ./SC_report cent:/root/My/go_shu 