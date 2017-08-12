### build package example
we have an example here to indicate how to build tcp package with the format of "header + body", this example tells us the way to use binary.BigEndian to produce header, use append to concate the slice of object.


<a href="#anchor_build_header"> 1. build the header </a><br>
<a href="#anchor_create_package"> 2. create the package </a><br>
<a href="#anchor_write_object"> 3. write the object to buffer </a><br>
<a href="#anchor_append_body"> 4. append body to buffer </a><br>
<a href="#anchor_system_endian"> 5. system endian </a><br>
<a href="#anchor_gbk_codec"> 6. gbk codec </a><br>

#### <a name="anchor_build_header">1 How to build header</a>
Generally we build package using the format of "header+body"

>
    type TMPPHeader struct {
	    Len        uint32
	    Version    uint32
	    CmdID      uint32
	    CmdStatus  uint32
	    SequenceNo uint32
	}

the 'Len' indicates the length of the whole package.
> 	
	header.Len = header_size + content_size

unsafe.Sizeof tells the size of header.

In C, we use zero-length array to indicate one buffer for the consistence of struct, like
> 	
	struct  test {
		int a;
		int b;
		char data[0];
	};

on win32-plat, the sizeof(test) = 8(4+4),but in golang, we have to use this data structure to represent this

>
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

define SizeofTag to calculate the size of Tag
> 	
	func SizeofTag(t Tag) uint32 {
		l := uint32(unsafe.Sizeof(t.TagHeader)) + uint32(t.Length)
		log.Println("sizeof tag:", l, "tag conent:", t)
		return l
	}

	
#### <a name="anchor_create_package">2 How to create package</a>
bytes package in golang offers the typical byte operation functions.
> 	
	b := new(bytes.Buffer)


#### <a name="anchor_write_object">3 How to write object to buffer</a>
binary package could help.
> 	
	// create header into buffer
	err := binary.Write(buf, binary.BigEndian, header)
	if err != nil {
		log.Println("failed to write header, err:", err.Error())
	}
	
	
#### <a name="anchor_append_body">4 Append body to buffer</a>
Assuming we have a taglist, and try to append the tags to buffer allocated before.
>	
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

one trivial poinit here is append "[]byte{0}" as the last zero-byte of the string to buffer.


#### <a name="anchor_system_endian">5 System Endian</a>
While create header, generaly we use binary.BigEndian to format the header to buffer, while append the body, we prefer to use LittleEndian.

#### <a name="anchor_gbk_codec"> 6 GBK Codec</a>
the package mahonia offers the way to transform gbk to utf8. see codec.go

