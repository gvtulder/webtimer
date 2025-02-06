package server

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *Command) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Key":
			z.Key, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Key")
				return
			}
		case "cmd":
			z.Command, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Command")
				return
			}
		case "sec":
			z.Seconds, err = dc.ReadInt64()
			if err != nil {
				err = msgp.WrapError(err, "Seconds")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z Command) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 3
	// write "Key"
	err = en.Append(0x83, 0xa3, 0x4b, 0x65, 0x79)
	if err != nil {
		return
	}
	err = en.WriteString(z.Key)
	if err != nil {
		err = msgp.WrapError(err, "Key")
		return
	}
	// write "cmd"
	err = en.Append(0xa3, 0x63, 0x6d, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.Command)
	if err != nil {
		err = msgp.WrapError(err, "Command")
		return
	}
	// write "sec"
	err = en.Append(0xa3, 0x73, 0x65, 0x63)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.Seconds)
	if err != nil {
		err = msgp.WrapError(err, "Seconds")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z Command) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 3
	// string "Key"
	o = append(o, 0x83, 0xa3, 0x4b, 0x65, 0x79)
	o = msgp.AppendString(o, z.Key)
	// string "cmd"
	o = append(o, 0xa3, 0x63, 0x6d, 0x64)
	o = msgp.AppendString(o, z.Command)
	// string "sec"
	o = append(o, 0xa3, 0x73, 0x65, 0x63)
	o = msgp.AppendInt64(o, z.Seconds)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *Command) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "Key":
			z.Key, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Key")
				return
			}
		case "cmd":
			z.Command, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Command")
				return
			}
		case "sec":
			z.Seconds, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Seconds")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z Command) Msgsize() (s int) {
	s = 1 + 4 + msgp.StringPrefixSize + len(z.Key) + 4 + msgp.StringPrefixSize + len(z.Command) + 4 + msgp.Int64Size
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TimerStateMessage) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "K":
			z.Key, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Key")
				return
			}
		case "a":
			z.Active, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "Active")
				return
			}
		case "b":
			z.Black, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "Black")
				return
			}
		case "r":
			z.Running, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "Running")
				return
			}
		case "c":
			z.Countdown, err = dc.ReadBool()
			if err != nil {
				err = msgp.WrapError(err, "Countdown")
				return
			}
		case "s":
			z.RemainingSeconds, err = dc.ReadInt64()
			if err != nil {
				err = msgp.WrapError(err, "RemainingSeconds")
				return
			}
		case "C":
			z.Clients, err = dc.ReadInt()
			if err != nil {
				err = msgp.WrapError(err, "Clients")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *TimerStateMessage) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 7
	// write "K"
	err = en.Append(0x87, 0xa1, 0x4b)
	if err != nil {
		return
	}
	err = en.WriteString(z.Key)
	if err != nil {
		err = msgp.WrapError(err, "Key")
		return
	}
	// write "a"
	err = en.Append(0xa1, 0x61)
	if err != nil {
		return
	}
	err = en.WriteBool(z.Active)
	if err != nil {
		err = msgp.WrapError(err, "Active")
		return
	}
	// write "b"
	err = en.Append(0xa1, 0x62)
	if err != nil {
		return
	}
	err = en.WriteBool(z.Black)
	if err != nil {
		err = msgp.WrapError(err, "Black")
		return
	}
	// write "r"
	err = en.Append(0xa1, 0x72)
	if err != nil {
		return
	}
	err = en.WriteBool(z.Running)
	if err != nil {
		err = msgp.WrapError(err, "Running")
		return
	}
	// write "c"
	err = en.Append(0xa1, 0x63)
	if err != nil {
		return
	}
	err = en.WriteBool(z.Countdown)
	if err != nil {
		err = msgp.WrapError(err, "Countdown")
		return
	}
	// write "s"
	err = en.Append(0xa1, 0x73)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.RemainingSeconds)
	if err != nil {
		err = msgp.WrapError(err, "RemainingSeconds")
		return
	}
	// write "C"
	err = en.Append(0xa1, 0x43)
	if err != nil {
		return
	}
	err = en.WriteInt(z.Clients)
	if err != nil {
		err = msgp.WrapError(err, "Clients")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TimerStateMessage) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 7
	// string "K"
	o = append(o, 0x87, 0xa1, 0x4b)
	o = msgp.AppendString(o, z.Key)
	// string "a"
	o = append(o, 0xa1, 0x61)
	o = msgp.AppendBool(o, z.Active)
	// string "b"
	o = append(o, 0xa1, 0x62)
	o = msgp.AppendBool(o, z.Black)
	// string "r"
	o = append(o, 0xa1, 0x72)
	o = msgp.AppendBool(o, z.Running)
	// string "c"
	o = append(o, 0xa1, 0x63)
	o = msgp.AppendBool(o, z.Countdown)
	// string "s"
	o = append(o, 0xa1, 0x73)
	o = msgp.AppendInt64(o, z.RemainingSeconds)
	// string "C"
	o = append(o, 0xa1, 0x43)
	o = msgp.AppendInt(o, z.Clients)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TimerStateMessage) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "K":
			z.Key, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Key")
				return
			}
		case "a":
			z.Active, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Active")
				return
			}
		case "b":
			z.Black, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Black")
				return
			}
		case "r":
			z.Running, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Running")
				return
			}
		case "c":
			z.Countdown, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Countdown")
				return
			}
		case "s":
			z.RemainingSeconds, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "RemainingSeconds")
				return
			}
		case "C":
			z.Clients, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Clients")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *TimerStateMessage) Msgsize() (s int) {
	s = 1 + 2 + msgp.StringPrefixSize + len(z.Key) + 2 + msgp.BoolSize + 2 + msgp.BoolSize + 2 + msgp.BoolSize + 2 + msgp.BoolSize + 2 + msgp.Int64Size + 2 + msgp.IntSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *VersionMessage) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "version":
			z.Version, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "Version")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z VersionMessage) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 1
	// write "version"
	err = en.Append(0x81, 0xa7, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteString(z.Version)
	if err != nil {
		err = msgp.WrapError(err, "Version")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z VersionMessage) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 1
	// string "version"
	o = append(o, 0x81, 0xa7, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.Version)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *VersionMessage) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "version":
			z.Version, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "Version")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z VersionMessage) Msgsize() (s int) {
	s = 1 + 8 + msgp.StringPrefixSize + len(z.Version)
	return
}
