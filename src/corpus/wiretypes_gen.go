package corpus

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *TCorpusImpl) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zcua uint32
	zcua, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zcua > 0 {
		zcua--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Words":
			var zxhx uint32
			zxhx, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Words) >= int(zxhx) {
				z.Words = (z.Words)[:zxhx]
			} else {
				z.Words = make(TWordItemSlice, zxhx)
			}
			for zxvk := range z.Words {
				err = z.Words[zxvk].DecodeMsg(dc)
				if err != nil {
					return
				}
			}
		case "Word2Idx":
			var zlqf uint32
			zlqf, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Word2Idx == nil && zlqf > 0 {
				z.Word2Idx = make(map[string]int32, zlqf)
			} else if len(z.Word2Idx) > 0 {
				for key, _ := range z.Word2Idx {
					delete(z.Word2Idx, key)
				}
			}
			for zlqf > 0 {
				zlqf--
				var zbzg string
				var zbai int32
				zbzg, err = dc.ReadString()
				if err != nil {
					return
				}
				zbai, err = dc.ReadInt32()
				if err != nil {
					return
				}
				z.Word2Idx[zbzg] = zbai
			}
		case "Doc2WordsIdx":
			var zdaf uint32
			zdaf, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Doc2WordsIdx) >= int(zdaf) {
				z.Doc2WordsIdx = (z.Doc2WordsIdx)[:zdaf]
			} else {
				z.Doc2WordsIdx = make([][]int32, zdaf)
			}
			for zcmr := range z.Doc2WordsIdx {
				var zpks uint32
				zpks, err = dc.ReadArrayHeader()
				if err != nil {
					return
				}
				if cap(z.Doc2WordsIdx[zcmr]) >= int(zpks) {
					z.Doc2WordsIdx[zcmr] = (z.Doc2WordsIdx[zcmr])[:zpks]
				} else {
					z.Doc2WordsIdx[zcmr] = make([]int32, zpks)
				}
				for zajw := range z.Doc2WordsIdx[zcmr] {
					z.Doc2WordsIdx[zcmr][zajw], err = dc.ReadInt32()
					if err != nil {
						return
					}
				}
			}
		case "Doc2Idx":
			var zjfb uint32
			zjfb, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			if z.Doc2Idx == nil && zjfb > 0 {
				z.Doc2Idx = make(map[string]int32, zjfb)
			} else if len(z.Doc2Idx) > 0 {
				for key, _ := range z.Doc2Idx {
					delete(z.Doc2Idx, key)
				}
			}
			for zjfb > 0 {
				zjfb--
				var zwht string
				var zhct int32
				zwht, err = dc.ReadString()
				if err != nil {
					return
				}
				zhct, err = dc.ReadInt32()
				if err != nil {
					return
				}
				z.Doc2Idx[zwht] = zhct
			}
		case "MinReduce":
			z.MinReduce, err = dc.ReadInt32()
			if err != nil {
				return
			}
		case "MinCnt":
			z.MinCnt, err = dc.ReadInt32()
			if err != nil {
				return
			}
		case "WordsCnt":
			z.WordsCnt, err = dc.ReadInt()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *TCorpusImpl) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 7
	// write "Words"
	err = en.Append(0x87, 0xa5, 0x57, 0x6f, 0x72, 0x64, 0x73)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Words)))
	if err != nil {
		return
	}
	for zxvk := range z.Words {
		err = z.Words[zxvk].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	// write "Word2Idx"
	err = en.Append(0xa8, 0x57, 0x6f, 0x72, 0x64, 0x32, 0x49, 0x64, 0x78)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.Word2Idx)))
	if err != nil {
		return
	}
	for zbzg, zbai := range z.Word2Idx {
		err = en.WriteString(zbzg)
		if err != nil {
			return
		}
		err = en.WriteInt32(zbai)
		if err != nil {
			return
		}
	}
	// write "Doc2WordsIdx"
	err = en.Append(0xac, 0x44, 0x6f, 0x63, 0x32, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x49, 0x64, 0x78)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Doc2WordsIdx)))
	if err != nil {
		return
	}
	for zcmr := range z.Doc2WordsIdx {
		err = en.WriteArrayHeader(uint32(len(z.Doc2WordsIdx[zcmr])))
		if err != nil {
			return
		}
		for zajw := range z.Doc2WordsIdx[zcmr] {
			err = en.WriteInt32(z.Doc2WordsIdx[zcmr][zajw])
			if err != nil {
				return
			}
		}
	}
	// write "Doc2Idx"
	err = en.Append(0xa7, 0x44, 0x6f, 0x63, 0x32, 0x49, 0x64, 0x78)
	if err != nil {
		return err
	}
	err = en.WriteMapHeader(uint32(len(z.Doc2Idx)))
	if err != nil {
		return
	}
	for zwht, zhct := range z.Doc2Idx {
		err = en.WriteString(zwht)
		if err != nil {
			return
		}
		err = en.WriteInt32(zhct)
		if err != nil {
			return
		}
	}
	// write "MinReduce"
	err = en.Append(0xa9, 0x4d, 0x69, 0x6e, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteInt32(z.MinReduce)
	if err != nil {
		return
	}
	// write "MinCnt"
	err = en.Append(0xa6, 0x4d, 0x69, 0x6e, 0x43, 0x6e, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteInt32(z.MinCnt)
	if err != nil {
		return
	}
	// write "WordsCnt"
	err = en.Append(0xa8, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x43, 0x6e, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteInt(z.WordsCnt)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TCorpusImpl) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 7
	// string "Words"
	o = append(o, 0x87, 0xa5, 0x57, 0x6f, 0x72, 0x64, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Words)))
	for zxvk := range z.Words {
		o, err = z.Words[zxvk].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	// string "Word2Idx"
	o = append(o, 0xa8, 0x57, 0x6f, 0x72, 0x64, 0x32, 0x49, 0x64, 0x78)
	o = msgp.AppendMapHeader(o, uint32(len(z.Word2Idx)))
	for zbzg, zbai := range z.Word2Idx {
		o = msgp.AppendString(o, zbzg)
		o = msgp.AppendInt32(o, zbai)
	}
	// string "Doc2WordsIdx"
	o = append(o, 0xac, 0x44, 0x6f, 0x63, 0x32, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x49, 0x64, 0x78)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Doc2WordsIdx)))
	for zcmr := range z.Doc2WordsIdx {
		o = msgp.AppendArrayHeader(o, uint32(len(z.Doc2WordsIdx[zcmr])))
		for zajw := range z.Doc2WordsIdx[zcmr] {
			o = msgp.AppendInt32(o, z.Doc2WordsIdx[zcmr][zajw])
		}
	}
	// string "Doc2Idx"
	o = append(o, 0xa7, 0x44, 0x6f, 0x63, 0x32, 0x49, 0x64, 0x78)
	o = msgp.AppendMapHeader(o, uint32(len(z.Doc2Idx)))
	for zwht, zhct := range z.Doc2Idx {
		o = msgp.AppendString(o, zwht)
		o = msgp.AppendInt32(o, zhct)
	}
	// string "MinReduce"
	o = append(o, 0xa9, 0x4d, 0x69, 0x6e, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65)
	o = msgp.AppendInt32(o, z.MinReduce)
	// string "MinCnt"
	o = append(o, 0xa6, 0x4d, 0x69, 0x6e, 0x43, 0x6e, 0x74)
	o = msgp.AppendInt32(o, z.MinCnt)
	// string "WordsCnt"
	o = append(o, 0xa8, 0x57, 0x6f, 0x72, 0x64, 0x73, 0x43, 0x6e, 0x74)
	o = msgp.AppendInt(o, z.WordsCnt)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TCorpusImpl) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zcxo uint32
	zcxo, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zcxo > 0 {
		zcxo--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Words":
			var zeff uint32
			zeff, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Words) >= int(zeff) {
				z.Words = (z.Words)[:zeff]
			} else {
				z.Words = make(TWordItemSlice, zeff)
			}
			for zxvk := range z.Words {
				bts, err = z.Words[zxvk].UnmarshalMsg(bts)
				if err != nil {
					return
				}
			}
		case "Word2Idx":
			var zrsw uint32
			zrsw, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.Word2Idx == nil && zrsw > 0 {
				z.Word2Idx = make(map[string]int32, zrsw)
			} else if len(z.Word2Idx) > 0 {
				for key, _ := range z.Word2Idx {
					delete(z.Word2Idx, key)
				}
			}
			for zrsw > 0 {
				var zbzg string
				var zbai int32
				zrsw--
				zbzg, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				zbai, bts, err = msgp.ReadInt32Bytes(bts)
				if err != nil {
					return
				}
				z.Word2Idx[zbzg] = zbai
			}
		case "Doc2WordsIdx":
			var zxpk uint32
			zxpk, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Doc2WordsIdx) >= int(zxpk) {
				z.Doc2WordsIdx = (z.Doc2WordsIdx)[:zxpk]
			} else {
				z.Doc2WordsIdx = make([][]int32, zxpk)
			}
			for zcmr := range z.Doc2WordsIdx {
				var zdnj uint32
				zdnj, bts, err = msgp.ReadArrayHeaderBytes(bts)
				if err != nil {
					return
				}
				if cap(z.Doc2WordsIdx[zcmr]) >= int(zdnj) {
					z.Doc2WordsIdx[zcmr] = (z.Doc2WordsIdx[zcmr])[:zdnj]
				} else {
					z.Doc2WordsIdx[zcmr] = make([]int32, zdnj)
				}
				for zajw := range z.Doc2WordsIdx[zcmr] {
					z.Doc2WordsIdx[zcmr][zajw], bts, err = msgp.ReadInt32Bytes(bts)
					if err != nil {
						return
					}
				}
			}
		case "Doc2Idx":
			var zobc uint32
			zobc, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			if z.Doc2Idx == nil && zobc > 0 {
				z.Doc2Idx = make(map[string]int32, zobc)
			} else if len(z.Doc2Idx) > 0 {
				for key, _ := range z.Doc2Idx {
					delete(z.Doc2Idx, key)
				}
			}
			for zobc > 0 {
				var zwht string
				var zhct int32
				zobc--
				zwht, bts, err = msgp.ReadStringBytes(bts)
				if err != nil {
					return
				}
				zhct, bts, err = msgp.ReadInt32Bytes(bts)
				if err != nil {
					return
				}
				z.Doc2Idx[zwht] = zhct
			}
		case "MinReduce":
			z.MinReduce, bts, err = msgp.ReadInt32Bytes(bts)
			if err != nil {
				return
			}
		case "MinCnt":
			z.MinCnt, bts, err = msgp.ReadInt32Bytes(bts)
			if err != nil {
				return
			}
		case "WordsCnt":
			z.WordsCnt, bts, err = msgp.ReadIntBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *TCorpusImpl) Msgsize() (s int) {
	s = 1 + 6 + msgp.ArrayHeaderSize
	for zxvk := range z.Words {
		s += z.Words[zxvk].Msgsize()
	}
	s += 9 + msgp.MapHeaderSize
	if z.Word2Idx != nil {
		for zbzg, zbai := range z.Word2Idx {
			_ = zbai
			s += msgp.StringPrefixSize + len(zbzg) + msgp.Int32Size
		}
	}
	s += 13 + msgp.ArrayHeaderSize
	for zcmr := range z.Doc2WordsIdx {
		s += msgp.ArrayHeaderSize + (len(z.Doc2WordsIdx[zcmr]) * (msgp.Int32Size))
	}
	s += 8 + msgp.MapHeaderSize
	if z.Doc2Idx != nil {
		for zwht, zhct := range z.Doc2Idx {
			_ = zhct
			s += msgp.StringPrefixSize + len(zwht) + msgp.Int32Size
		}
	}
	s += 10 + msgp.Int32Size + 7 + msgp.Int32Size + 9 + msgp.IntSize
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TWordItem) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zema uint32
	zema, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zema > 0 {
		zema--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Cnt":
			z.Cnt, err = dc.ReadInt32()
			if err != nil {
				return
			}
		case "Point":
			var zpez uint32
			zpez, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Point) >= int(zpez) {
				z.Point = (z.Point)[:zpez]
			} else {
				z.Point = make([]int32, zpez)
			}
			for zsnv := range z.Point {
				z.Point[zsnv], err = dc.ReadInt32()
				if err != nil {
					return
				}
			}
		case "Code":
			var zqke uint32
			zqke, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Code) >= int(zqke) {
				z.Code = (z.Code)[:zqke]
			} else {
				z.Code = make([]bool, zqke)
			}
			for zkgt := range z.Code {
				z.Code[zkgt], err = dc.ReadBool()
				if err != nil {
					return
				}
			}
		case "Word":
			z.Word, err = dc.ReadString()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *TWordItem) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "Cnt"
	err = en.Append(0x84, 0xa3, 0x43, 0x6e, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteInt32(z.Cnt)
	if err != nil {
		return
	}
	// write "Point"
	err = en.Append(0xa5, 0x50, 0x6f, 0x69, 0x6e, 0x74)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Point)))
	if err != nil {
		return
	}
	for zsnv := range z.Point {
		err = en.WriteInt32(z.Point[zsnv])
		if err != nil {
			return
		}
	}
	// write "Code"
	err = en.Append(0xa4, 0x43, 0x6f, 0x64, 0x65)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Code)))
	if err != nil {
		return
	}
	for zkgt := range z.Code {
		err = en.WriteBool(z.Code[zkgt])
		if err != nil {
			return
		}
	}
	// write "Word"
	err = en.Append(0xa4, 0x57, 0x6f, 0x72, 0x64)
	if err != nil {
		return err
	}
	err = en.WriteString(z.Word)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TWordItem) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "Cnt"
	o = append(o, 0x84, 0xa3, 0x43, 0x6e, 0x74)
	o = msgp.AppendInt32(o, z.Cnt)
	// string "Point"
	o = append(o, 0xa5, 0x50, 0x6f, 0x69, 0x6e, 0x74)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Point)))
	for zsnv := range z.Point {
		o = msgp.AppendInt32(o, z.Point[zsnv])
	}
	// string "Code"
	o = append(o, 0xa4, 0x43, 0x6f, 0x64, 0x65)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Code)))
	for zkgt := range z.Code {
		o = msgp.AppendBool(o, z.Code[zkgt])
	}
	// string "Word"
	o = append(o, 0xa4, 0x57, 0x6f, 0x72, 0x64)
	o = msgp.AppendString(o, z.Word)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TWordItem) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zqyh uint32
	zqyh, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zqyh > 0 {
		zqyh--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Cnt":
			z.Cnt, bts, err = msgp.ReadInt32Bytes(bts)
			if err != nil {
				return
			}
		case "Point":
			var zyzr uint32
			zyzr, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Point) >= int(zyzr) {
				z.Point = (z.Point)[:zyzr]
			} else {
				z.Point = make([]int32, zyzr)
			}
			for zsnv := range z.Point {
				z.Point[zsnv], bts, err = msgp.ReadInt32Bytes(bts)
				if err != nil {
					return
				}
			}
		case "Code":
			var zywj uint32
			zywj, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Code) >= int(zywj) {
				z.Code = (z.Code)[:zywj]
			} else {
				z.Code = make([]bool, zywj)
			}
			for zkgt := range z.Code {
				z.Code[zkgt], bts, err = msgp.ReadBoolBytes(bts)
				if err != nil {
					return
				}
			}
		case "Word":
			z.Word, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *TWordItem) Msgsize() (s int) {
	s = 1 + 4 + msgp.Int32Size + 6 + msgp.ArrayHeaderSize + (len(z.Point) * (msgp.Int32Size)) + 5 + msgp.ArrayHeaderSize + (len(z.Code) * (msgp.BoolSize)) + 5 + msgp.StringPrefixSize + len(z.Word)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TWordItemSlice) DecodeMsg(dc *msgp.Reader) (err error) {
	var zrfe uint32
	zrfe, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zrfe) {
		(*z) = (*z)[:zrfe]
	} else {
		(*z) = make(TWordItemSlice, zrfe)
	}
	for zzpf := range *z {
		err = (*z)[zzpf].DecodeMsg(dc)
		if err != nil {
			return
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z TWordItemSlice) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zgmo := range z {
		err = z[zgmo].EncodeMsg(en)
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z TWordItemSlice) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zgmo := range z {
		o, err = z[zgmo].MarshalMsg(o)
		if err != nil {
			return
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TWordItemSlice) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var zeth uint32
	zeth, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(zeth) {
		(*z) = (*z)[:zeth]
	} else {
		(*z) = make(TWordItemSlice, zeth)
	}
	for ztaf := range *z {
		bts, err = (*z)[ztaf].UnmarshalMsg(bts)
		if err != nil {
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z TWordItemSlice) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize
	for zsbz := range z {
		s += z[zsbz].Msgsize()
	}
	return
}
