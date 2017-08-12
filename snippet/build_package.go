// sytemendian project sytemendian.go
package snippet

import (
	"bytes"
	"encoding/binary"
	"log"

	"net"
	"unsafe"
)

// tmpp header
type TMPPHeader struct {
	Len        uint32
	Version    uint32
	CmdID      uint32
	CmdStatus  uint32
	SequenceNo uint32
}

type TagHeader struct {
	Type uint32
	// indicate the length of the field Data in Tag
	Length uint32
}

//sizeof(tag) = sizeof(tagheader) + tagheader.Length
type Tag struct {
	TagHeader
	Data string
}

func SizeofTag(t Tag) uint32 {

	l := uint32(unsafe.Sizeof(t.TagHeader)) + uint32(t.Length)
	log.Println("sizeof tag:", l, "tag conent:", t)
	return l
}

// append tag to []byte
// should pass pointer to slice since we want modify the input parameter buf
func AppendTag(buf *[]byte, t Tag) {
	// append content of t.Type
	// need [:] to change array to slice
	typeb := (*[4]byte)(unsafe.Pointer(&t.Type))[:]
	*buf = append(*buf, typeb...)

	// append content of t.Length
	lenb := (*[4]byte)(unsafe.Pointer(&t.Length))[:]
	*buf = append(*buf, lenb...)

	// append content of t.Data
	*buf = append(*buf, []byte(t.Data)...)

	// append the '\0' as the last byte
	*buf = append(*buf, []byte{0}...)

	log.Println("appendTag buf size:", len(*buf))
}

func BuildConn() (*net.TCPConn, error) {

	addr, err := net.ResolveTCPAddr("tcp", "10.121.129.15:8100")
	if err != nil {
		log.Println("failed to resolve tcp ")
		return nil, err
	}

	c, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Println("faield to dial tcp")
		return nil, err
	}

	return c, nil
}

func TestBuildPackage() {

	var header TMPPHeader
	header.CmdID = 0x00000001 | 0x00006000
	header.SequenceNo = 1
	header.Version = 1
	header.CmdStatus = 1

	marr := []string{"13671517560", "21457368424652335413579781245085", "02709", "yilinliu"}

	header_size := unsafe.Sizeof(header)
	log.Println("debug header size: ", header_size)

	taglist := make([]Tag, 0)

	// calc the content size to fill the Len field of header
	content_size := uint32(0)

	for index, e := range marr {
		var t Tag
		t.Type = uint32(index + 1)
		t.Length = uint32(len(e)) + 1
		t.Data = e
		l := SizeofTag(t)
		content_size += l
		taglist = append(taglist, t)
	}

	// fill the whole length of header
	header.Len = uint32(header_size) + uint32(content_size)
	log.Println("content_lenght:", content_size, "whole length with header:", header.Len)

	// create buffer
	buf := new(bytes.Buffer)

	// create header into buffer
	binary.Write(buf, binary.BigEndian, header)
	for _, e := range buf.Bytes() {
		log.Printf("%02x", e)
	}
	log.Println("header len:", len(buf.Bytes()))
	fbuf := buf.Bytes()[:]
	log.Print("before write body:", fbuf, ", fbuf size:", len(fbuf))

	// append content to buffer
	for _, t := range taglist {
		log.Println("write tag:", t, ", tag size:", SizeofTag(t))

		AppendTag(&fbuf, t)

	}
	// dump fbuf
	log.Println("fbuf size:", len(fbuf))

	for _, e := range fbuf {
		log.Printf("%02x", e)
	}

	c, err := BuildConn()
	if err != nil {

		log.Println(err.Error())
		return
	}

	n, err := c.Write(fbuf)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("send bytes: ", n)
}
