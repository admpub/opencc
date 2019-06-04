package opencc

import (
	"fmt"
	"testing"
)

//
func Test_opencc(t *testing.T) {
	cc, err := NewOpenCC(S2TWP)
	if err != nil {
		fmt.Println(err)
		return
	}
	nText := cc.ConvertText(`迪拜（阿拉伯语：دبي，英语：Dubai），是阿拉伯联合酋长国人口最多的城市，位于波斯湾东南海岸，迪拜也是组成阿联酋七个酋长国之一——迪拜酋长国的首都。`)
	fmt.Println(nText)
}
