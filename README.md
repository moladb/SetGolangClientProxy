### build package example

#### 1 How to build header
Generally we build package using the format of "header+body"

> // tmpp header
    type TMPPHeader struct {
	    Len        uint32
	    Version    uint32
	    CmdID      uint32
	    CmdStatus  uint32
	    SequenceNo uint32
}

the 'Len' indicates the length of the whole package.
> header.Len = header_size + content_size

unsafe.Sizeof tells the size of header.

In C, we use zero-length array to indicate one buffer for the consistence of struct, like
> struct  test {
    int a;
    int b;
    char data[0];
};

on win32-plat, the sizeof(test) = 8(4+4),but in golang, we have to use this data structure to represent this
> type TagHeader struct {
	Type uint32
	// indicate the length of the field Data in Tag
	Length uint32
   }
    //sizeof(tag) = sizeof(tagheader) + tagheader.Length
    type Tag struct {
	    TagHeader
	    Data string
    }

define SizeofTag to calculate the size of Tag
> func SizeofTag(t Tag) uint32 {
    l := uint32(unsafe.Sizeof(t.TagHeader)) + uint32(t.Length)
	log.Println("sizeof tag:", l, "tag conent:", t)
	return l
}

#### 2 How to create package
bytes package in golang offers the typical byte operation functions.
> b := new(bytes.Buffer)
#### how to write object to buffer
binary package could help.
> 	// create header into buffer
	err := binary.Write(buf, binary.BigEndian, header)
	if err != nil {
		log.Println("failed to write header, err:", err.Error())
	}
#### 3 append body to buffer
assuming we have a taglist, and try to append the tags to buffer allocated before.
> // append tag to []byte
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

one trivial poinit here is append "[]byte{0}" as the last zero-byte of the string to buffer.

####4 System Endian
while create header, generaly we use binary.BigEndian to format the header to buffer, while append the body, we prefer to use LittleEndian.


