package opencc

import (
	"fmt"
	"testing"
)

//
func Test_opencc(t *testing.T) {
	cc, err := NewOpenCC("s2twp")
	if err != nil {
		fmt.Println(err)
		return
	}
	nText := cc.ConvertText(`保税工厂声明：本书为无限小说网(txt53.com)以下作品内容之版权与本站无任何关系`)

	fmt.Println(nText)

}
