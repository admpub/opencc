package opencc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"unicode"

	"github.com/wzshiming/opencc/data"
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

//
func (oc *OpenCC) ConvertFile(in io.Reader, out io.Writer) error {
	inReader := bufio.NewReader(in)
	for {
		lineText, readErr := inReader.ReadString('\n')
		if readErr != nil && readErr != io.EOF {
			return readErr
		}
		nLineText, err := oc.splitText(lineText)
		if err != nil {
			return err
		}
		_, err = out.Write([]byte(nLineText))
		if err != nil {
			return err
		}
		if readErr == io.EOF {
			break
		}
	}
	return nil
}

//
func (oc *OpenCC) ConvertText(text string) (string, error) {
	return oc.splitText(text)
}

//
func (oc *OpenCC) splitText(text string) (string, error) {
	prev := 0
	newText := bytes.NewBuffer(nil)
	for i, c := range text {
		if isPunctuations(c) {
			v := text[prev:i]
			tx, err := oc.convertString(v)
			if err != nil {
				return text, err
			}
			newText.WriteString(tx)
			prev = i
		}
	}

	v := text[prev:]
	tx, err := oc.convertString(v)
	if err != nil {
		return text, err
	}
	newText.WriteString(tx)

	return newText.String(), nil
}

//
func (oc *OpenCC) convertString(text string) (string, error) {
	var err error
	if oc.conf == nil {
		return text, fmt.Errorf("no config")
	}
	text, err = oc.conf.convertText(text)
	if err != nil {
		return text, err
	}
	return text, nil
}

//是否标点符号
func isPunctuations(character rune) bool {
	return unicode.In(character, unicode.Punct, unicode.Space)
}
