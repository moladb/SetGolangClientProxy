package snippet

import (
	"github.com/axgle/mahonia"
)

type Codec struct {
}

func (c *Codec) ConvertUTF8ToGBK(s string) string {

	// use NewEncoder
	gbkCoder := mahonia.NewEncoder("gbk")

	// 	Encoder's ConvertString will convert UTF8 to the gbk
	return gbkCoder.ConvertString(s)

}

func (c *Codec) ConvertGBKToUTF8(s string) string {

	// decoder will decode the character from the gbk and convert to utf-8
	gbkDecoder := mahonia.NewDecoder("gbk")
	return gbkDecoder.ConvertString(s)
}
