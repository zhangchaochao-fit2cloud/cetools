package ssh

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestScp(t *testing.T) {
	con, err := buildClientConfig()
	if err != nil {
		t.Fatalf("Couldn't build the client configuration: %s", err)
	}

	run, err := con.Run("ls /tmp")
	if err != nil {
		t.Fatalf("Couldn't build the client configuration: %s", err)
	}
	f, _ := os.Open("/opt/conf/app.yaml")
	err = con.Upload(context.Background(), f, "/tmp/app.yaml", "0777")
	if err != nil {
		t.Fatalf("Couldn't build the client configuration: %s", err)
	}
	fmt.Println(run)
}

func buildClientConfig() (*ConnInfo, error) {
	con := &ConnInfo{
		User:     "root",
		Addr:     "10.1.10.10",
		Port:     22,
		AuthMode: "password",
		Password: "",
	}

	return con, nil
}
