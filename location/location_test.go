package location

import (
	"fmt"
	"testing"
)

func TestGetIpgeolocationInfo(t *testing.T) {

	info, err := GetIpgeolocationInfo("128.14.140.254")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(info)

}
