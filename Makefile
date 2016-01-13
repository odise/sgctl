build:
	go build sgctl.go

get:
	go get -u github.com/aws/aws-sdk-go/...
	go get -u github.com/spf13/cobra

clean:
	rm -f sgctl
