package opencc

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/admpub/opencc/data"
)

const (
	// Traditional Chinese (Hong Kong standard) to Simplified Chinese
	HK2S = "hk2s"
	// Simplified Chinese to Traditional Chinese (Taiwan standard, with phrases)
	S2HK = "s2hk"
	// Simplified Chinese to Traditional Chinese
	S2T = "s2t"
	// Simplified Chinese to Traditional Chinese (Taiwan standard)
	S2TW = "s2tw"
	// Simplified Chinese to Traditional Chinese (Taiwan standard, with phrases)
	S2TWP = "s2twp"
	// Traditional Chinese to Traditional Chinese (Hong Kong standard)
	T2HK = "t2hk"
	// Traditional Chinese to Simplified Chinese
	T2S = "t2s"
	// Traditional Chinese to Traditional Chinese (Taiwan standard)
	T2TW = "t2tw"
	// Traditional Chinese (Taiwan standard) to Simplified Chinese
	TW2S = "tw2s"
	// Traditional Chinese (Taiwan standard) to Simplified Chinese (with phrases)
	TW2SP = "tw2sp"
)

type FileOCD string

type Dict struct {
	Type   string  `json:"type"`
	File   FileOCD `json:"file"`
	Dicts  []*Dict `json:"dicts"`
	CfgMap map[string][]string
	maxLen int //最大长度
	minLen int //最小长度
}

type Segmentation struct {
	Type string `json:"type"`
	Dict *Dict  `json:"dict"`
}

type ConversionChain struct {
	Dict *Dict `json:"dict"`
}

type Config struct {
	Name            string             `json:"name"`
	Segmentation    Segmentation       `json:"segmentation"`
	ConversionChain []*ConversionChain `json:"conversion_chain"`
}

func (c *Config) init() error {
	var err error
	for _, cv := range c.ConversionChain {
		err = cv.init()
		if err != nil {
			return err
		}
	}
	return nil
}
func (cv *ConversionChain) init() error {
	return cv.Dict.init()
}

func (d *Dict) init() (err error) {
	if len(d.File) > 0 {
		d.CfgMap, d.maxLen, d.minLen, err = d.File.readFile()
		//fmt.Println("File = ", string(d.File), d.maxLen, d.minLen)
		if err != nil {
			return err
		}
	}
	if len(d.Dicts) > 0 {
		for _, childDict := range d.Dicts {
			err = childDict.init()
			if err != nil {
				return
			}
		}
	}
	return nil
}

func (fo *FileOCD) readFile() (map[string][]string, int, int, error) {
	f, err := data.Asset("dictionary/" + string(*fo))
	if err != nil {
		return nil, 0, 0, err
	}
	cfgMap := make(map[string][]string)
	buf := bufio.NewReader(bytes.NewBuffer(f))
	max := 0
	min := 0
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			return cfgMap, max, min, nil
		}
		fields := strings.Fields(line)
		if len(fields) > 1 {
			if len([]rune(fields[0])) > max {
				max = len([]rune(fields[0]))
			}
			if min <= 0 || len([]rune(fields[0])) < min {
				min = len([]rune(fields[0]))
			}
			cfgMap[fields[0]] = fields[1:]
		}
	}
	return cfgMap, max, min, nil
}

// =============================================================
func (c *Config) convertText(text string) string {
	for _, cv := range c.ConversionChain {
		text = cv.convertText(text)
	}
	return text
}
func (c *ConversionChain) convertText(text string) string {
	return c.Dict.convertTextWithMap(text)
}

func (d *Dict) convertTextWithMap(text string) string {
	newText := text
	runes := []rune(text)
	//
	if d.CfgMap != nil {
		if len(runes) < d.minLen {
			return text
		}
		//
		maxL := d.maxLen
		if maxL > len(runes) {
			maxL = len(runes)
		}
		//
		for i := maxL; i >= d.minLen; i-- {
			for j := 0; j <= len(runes)-i; j++ {
				if i == 0 || j+i > len(runes) {
					continue
				}
				old := string(runes[j : j+i])
				if newStr, ok := d.CfgMap[old]; ok {
					newText = strings.Replace(newText, old, newStr[0], 1)
					j = j + i - 1
				}
			}
		}
	}
	//
	for _, cd := range d.Dicts {
		newText = cd.convertTextWithMap(newText)
	}
	return newText
}
