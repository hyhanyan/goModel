package mail

import (
	"fmt"
	"testing"
)

func TestMail(t *testing.T) {
	err := SendMail("2222xxxxxx@qq.com", "object", "test")
	if err != nil {
		fmt.Println(err.Error())
	}
}
