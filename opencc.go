package opencc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"unicode"

	"github.com/admpub/opencc/data"
)

type OpenCC struct {
	conf *Config
}

// Supported conversions: s2t, t2s, s2tw, tw2s, s2hk, hk2s, s2twp, tw2sp, t2tw, t2hk
func NewOpenCC(conversions string) (*OpenCC, error) {
	body, err := data.Asset("config/" + conversions + ".json")
	if err != nil {
		return nil, err
	}
	var conf *Config
	err = json.Unmarshal(body, &conf)
	if err != nil {
		return nil, err
	}
	err = conf.init()
	if err != nil {
		return nil, err
	}
	//
	return &OpenCC{conf: conf}, nil
}

func (oc *OpenCC) Name() string {
	return oc.conf.Name
}

func (oc *OpenCC) ConvertFile(in io.Reader, out io.Writer) error {
	inReader := bufio.NewReader(in)
	for {
		lineText, err := inReader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		nLineText := oc.splitText(lineText)
		_, err = io.WriteString(out, nLineText)
		if err != nil {
			return err
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}

func (oc *OpenCC) ConvertText(text string) string {
	return oc.splitText(text)
}

func (oc *OpenCC) splitText(text string) string {
	prev := 0
	newText := bytes.NewBuffer(nil)
	for i, c := range text {
		if isPunctuations(c) {
			v := text[prev:i]
			tx := oc.convertString(v)
			newText.WriteString(tx)
			prev = i
		}
	}

	v := text[prev:]
	tx := oc.convertString(v)
	newText.WriteString(tx)

	return newText.String()
}

func (oc *OpenCC) convertString(text string) string {
	return oc.conf.convertText(text)
}

// 是否标点符号
func isPunctuations(character rune) bool {
	return unicode.In(character, unicode.Punct, unicode.Space)
}
