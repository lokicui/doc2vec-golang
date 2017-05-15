package neuralnet

// NOTE: THIS FILE WAS PRODUCED BY THE
// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)
// DO NOT EDIT

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *TNeuralNetImpl) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zxhx uint32
	zxhx, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zxhx > 0 {
		zxhx--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Syn0":
			var zlqf uint32
			zlqf, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Syn0) >= int(zlqf) {
				z.Syn0 = (z.Syn0)[:zlqf]
			} else {
				z.Syn0 = make([]TVector, zlqf)
			}
			for zxvk := range z.Syn0 {
				var zdaf uint32
				zdaf, err = dc.ReadArrayHeader()
				if err != nil {
					return
				}
				if cap(z.Syn0[zxvk]) >= int(zdaf) {
					z.Syn0[zxvk] = (z.Syn0[zxvk])[:zdaf]
				} else {
					z.Syn0[zxvk] = make(TVector, zdaf)
				}
				for zbzg := range z.Syn0[zxvk] {
					z.Syn0[zxvk][zbzg], err = dc.ReadFloat32()
					if err != nil {
						return
					}
				}
			}
		case "Dsyn0":
			var zpks uint32
			zpks, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Dsyn0) >= int(zpks) {
				z.Dsyn0 = (z.Dsyn0)[:zpks]
			} else {
				z.Dsyn0 = make([]TVector, zpks)
			}
			for zbai := range z.Dsyn0 {
				var zjfb uint32
				zjfb, err = dc.ReadArrayHeader()
				if err != nil {
					return
				}
				if cap(z.Dsyn0[zbai]) >= int(zjfb) {
					z.Dsyn0[zbai] = (z.Dsyn0[zbai])[:zjfb]
				} else {
					z.Dsyn0[zbai] = make(TVector, zjfb)
				}
				for zcmr := range z.Dsyn0[zbai] {
					z.Dsyn0[zbai][zcmr], err = dc.ReadFloat32()
					if err != nil {
						return
					}
				}
			}
		case "Syn1":
			var zcxo uint32
			zcxo, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Syn1) >= int(zcxo) {
				z.Syn1 = (z.Syn1)[:zcxo]
			} else {
				z.Syn1 = make([]TVector, zcxo)
			}
			for zajw := range z.Syn1 {
				var zeff uint32
				zeff, err = dc.ReadArrayHeader()
				if err != nil {
					return
				}
				if cap(z.Syn1[zajw]) >= int(zeff) {
					z.Syn1[zajw] = (z.Syn1[zajw])[:zeff]
				} else {
					z.Syn1[zajw] = make(TVector, zeff)
				}
				for zwht := range z.Syn1[zajw] {
					z.Syn1[zajw][zwht], err = dc.ReadFloat32()
					if err != nil {
						return
					}
				}
			}
		case "Syn1neg":
			var zrsw uint32
			zrsw, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Syn1neg) >= int(zrsw) {
				z.Syn1neg = (z.Syn1neg)[:zrsw]
			} else {
				z.Syn1neg = make([]TVector, zrsw)
			}
			for zhct := range z.Syn1neg {
				var zxpk uint32
				zxpk, err = dc.ReadArrayHeader()
				if err != nil {
					return
				}
				if cap(z.Syn1neg[zhct]) >= int(zxpk) {
					z.Syn1neg[zhct] = (z.Syn1neg[zhct])[:zxpk]
				} else {
					z.Syn1neg[zhct] = make(TVector, zxpk)
				}
				for zcua := range z.Syn1neg[zhct] {
					z.Syn1neg[zhct][zcua], err = dc.ReadFloat32()
					if err != nil {
						return
					}
				}
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
func (z *TNeuralNetImpl) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 4
	// write "Syn0"
	err = en.Append(0x84, 0xa4, 0x53, 0x79, 0x6e, 0x30)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Syn0)))
	if err != nil {
		return
	}
	for zxvk := range z.Syn0 {
		err = en.WriteArrayHeader(uint32(len(z.Syn0[zxvk])))
		if err != nil {
			return
		}
		for zbzg := range z.Syn0[zxvk] {
			err = en.WriteFloat32(z.Syn0[zxvk][zbzg])
			if err != nil {
				return
			}
		}
	}
	// write "Dsyn0"
	err = en.Append(0xa5, 0x44, 0x73, 0x79, 0x6e, 0x30)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Dsyn0)))
	if err != nil {
		return
	}
	for zbai := range z.Dsyn0 {
		err = en.WriteArrayHeader(uint32(len(z.Dsyn0[zbai])))
		if err != nil {
			return
		}
		for zcmr := range z.Dsyn0[zbai] {
			err = en.WriteFloat32(z.Dsyn0[zbai][zcmr])
			if err != nil {
				return
			}
		}
	}
	// write "Syn1"
	err = en.Append(0xa4, 0x53, 0x79, 0x6e, 0x31)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Syn1)))
	if err != nil {
		return
	}
	for zajw := range z.Syn1 {
		err = en.WriteArrayHeader(uint32(len(z.Syn1[zajw])))
		if err != nil {
			return
		}
		for zwht := range z.Syn1[zajw] {
			err = en.WriteFloat32(z.Syn1[zajw][zwht])
			if err != nil {
				return
			}
		}
	}
	// write "Syn1neg"
	err = en.Append(0xa7, 0x53, 0x79, 0x6e, 0x31, 0x6e, 0x65, 0x67)
	if err != nil {
		return err
	}
	err = en.WriteArrayHeader(uint32(len(z.Syn1neg)))
	if err != nil {
		return
	}
	for zhct := range z.Syn1neg {
		err = en.WriteArrayHeader(uint32(len(z.Syn1neg[zhct])))
		if err != nil {
			return
		}
		for zcua := range z.Syn1neg[zhct] {
			err = en.WriteFloat32(z.Syn1neg[zhct][zcua])
			if err != nil {
				return
			}
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TNeuralNetImpl) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 4
	// string "Syn0"
	o = append(o, 0x84, 0xa4, 0x53, 0x79, 0x6e, 0x30)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Syn0)))
	for zxvk := range z.Syn0 {
		o = msgp.AppendArrayHeader(o, uint32(len(z.Syn0[zxvk])))
		for zbzg := range z.Syn0[zxvk] {
			o = msgp.AppendFloat32(o, z.Syn0[zxvk][zbzg])
		}
	}
	// string "Dsyn0"
	o = append(o, 0xa5, 0x44, 0x73, 0x79, 0x6e, 0x30)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Dsyn0)))
	for zbai := range z.Dsyn0 {
		o = msgp.AppendArrayHeader(o, uint32(len(z.Dsyn0[zbai])))
		for zcmr := range z.Dsyn0[zbai] {
			o = msgp.AppendFloat32(o, z.Dsyn0[zbai][zcmr])
		}
	}
	// string "Syn1"
	o = append(o, 0xa4, 0x53, 0x79, 0x6e, 0x31)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Syn1)))
	for zajw := range z.Syn1 {
		o = msgp.AppendArrayHeader(o, uint32(len(z.Syn1[zajw])))
		for zwht := range z.Syn1[zajw] {
			o = msgp.AppendFloat32(o, z.Syn1[zajw][zwht])
		}
	}
	// string "Syn1neg"
	o = append(o, 0xa7, 0x53, 0x79, 0x6e, 0x31, 0x6e, 0x65, 0x67)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Syn1neg)))
	for zhct := range z.Syn1neg {
		o = msgp.AppendArrayHeader(o, uint32(len(z.Syn1neg[zhct])))
		for zcua := range z.Syn1neg[zhct] {
			o = msgp.AppendFloat32(o, z.Syn1neg[zhct][zcua])
		}
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TNeuralNetImpl) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zdnj uint32
	zdnj, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zdnj > 0 {
		zdnj--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Syn0":
			var zobc uint32
			zobc, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Syn0) >= int(zobc) {
				z.Syn0 = (z.Syn0)[:zobc]
			} else {
				z.Syn0 = make([]TVector, zobc)
			}
			for zxvk := range z.Syn0 {
				var zsnv uint32
				zsnv, bts, err = msgp.ReadArrayHeaderBytes(bts)
				if err != nil {
					return
				}
				if cap(z.Syn0[zxvk]) >= int(zsnv) {
					z.Syn0[zxvk] = (z.Syn0[zxvk])[:zsnv]
				} else {
					z.Syn0[zxvk] = make(TVector, zsnv)
				}
				for zbzg := range z.Syn0[zxvk] {
					z.Syn0[zxvk][zbzg], bts, err = msgp.ReadFloat32Bytes(bts)
					if err != nil {
						return
					}
				}
			}
		case "Dsyn0":
			var zkgt uint32
			zkgt, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Dsyn0) >= int(zkgt) {
				z.Dsyn0 = (z.Dsyn0)[:zkgt]
			} else {
				z.Dsyn0 = make([]TVector, zkgt)
			}
			for zbai := range z.Dsyn0 {
				var zema uint32
				zema, bts, err = msgp.ReadArrayHeaderBytes(bts)
				if err != nil {
					return
				}
				if cap(z.Dsyn0[zbai]) >= int(zema) {
					z.Dsyn0[zbai] = (z.Dsyn0[zbai])[:zema]
				} else {
					z.Dsyn0[zbai] = make(TVector, zema)
				}
				for zcmr := range z.Dsyn0[zbai] {
					z.Dsyn0[zbai][zcmr], bts, err = msgp.ReadFloat32Bytes(bts)
					if err != nil {
						return
					}
				}
			}
		case "Syn1":
			var zpez uint32
			zpez, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Syn1) >= int(zpez) {
				z.Syn1 = (z.Syn1)[:zpez]
			} else {
				z.Syn1 = make([]TVector, zpez)
			}
			for zajw := range z.Syn1 {
				var zqke uint32
				zqke, bts, err = msgp.ReadArrayHeaderBytes(bts)
				if err != nil {
					return
				}
				if cap(z.Syn1[zajw]) >= int(zqke) {
					z.Syn1[zajw] = (z.Syn1[zajw])[:zqke]
				} else {
					z.Syn1[zajw] = make(TVector, zqke)
				}
				for zwht := range z.Syn1[zajw] {
					z.Syn1[zajw][zwht], bts, err = msgp.ReadFloat32Bytes(bts)
					if err != nil {
						return
					}
				}
			}
		case "Syn1neg":
			var zqyh uint32
			zqyh, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Syn1neg) >= int(zqyh) {
				z.Syn1neg = (z.Syn1neg)[:zqyh]
			} else {
				z.Syn1neg = make([]TVector, zqyh)
			}
			for zhct := range z.Syn1neg {
				var zyzr uint32
				zyzr, bts, err = msgp.ReadArrayHeaderBytes(bts)
				if err != nil {
					return
				}
				if cap(z.Syn1neg[zhct]) >= int(zyzr) {
					z.Syn1neg[zhct] = (z.Syn1neg[zhct])[:zyzr]
				} else {
					z.Syn1neg[zhct] = make(TVector, zyzr)
				}
				for zcua := range z.Syn1neg[zhct] {
					z.Syn1neg[zhct][zcua], bts, err = msgp.ReadFloat32Bytes(bts)
					if err != nil {
						return
					}
				}
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
func (z *TNeuralNetImpl) Msgsize() (s int) {
	s = 1 + 5 + msgp.ArrayHeaderSize
	for zxvk := range z.Syn0 {
		s += msgp.ArrayHeaderSize + (len(z.Syn0[zxvk]) * (msgp.Float32Size))
	}
	s += 6 + msgp.ArrayHeaderSize
	for zbai := range z.Dsyn0 {
		s += msgp.ArrayHeaderSize + (len(z.Dsyn0[zbai]) * (msgp.Float32Size))
	}
	s += 5 + msgp.ArrayHeaderSize
	for zajw := range z.Syn1 {
		s += msgp.ArrayHeaderSize + (len(z.Syn1[zajw]) * (msgp.Float32Size))
	}
	s += 8 + msgp.ArrayHeaderSize
	for zhct := range z.Syn1neg {
		s += msgp.ArrayHeaderSize + (len(z.Syn1neg[zhct]) * (msgp.Float32Size))
	}
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TVector) DecodeMsg(dc *msgp.Reader) (err error) {
	var zzpf uint32
	zzpf, err = dc.ReadArrayHeader()
	if err != nil {
		return
	}
	if cap((*z)) >= int(zzpf) {
		(*z) = (*z)[:zzpf]
	} else {
		(*z) = make(TVector, zzpf)
	}
	for zjpj := range *z {
		(*z)[zjpj], err = dc.ReadFloat32()
		if err != nil {
			return
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z TVector) EncodeMsg(en *msgp.Writer) (err error) {
	err = en.WriteArrayHeader(uint32(len(z)))
	if err != nil {
		return
	}
	for zrfe := range z {
		err = en.WriteFloat32(z[zrfe])
		if err != nil {
			return
		}
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z TVector) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	o = msgp.AppendArrayHeader(o, uint32(len(z)))
	for zrfe := range z {
		o = msgp.AppendFloat32(o, z[zrfe])
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TVector) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var ztaf uint32
	ztaf, bts, err = msgp.ReadArrayHeaderBytes(bts)
	if err != nil {
		return
	}
	if cap((*z)) >= int(ztaf) {
		(*z) = (*z)[:ztaf]
	} else {
		(*z) = make(TVector, ztaf)
	}
	for zgmo := range *z {
		(*z)[zgmo], bts, err = msgp.ReadFloat32Bytes(bts)
		if err != nil {
			return
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z TVector) Msgsize() (s int) {
	s = msgp.ArrayHeaderSize + (len(z) * (msgp.Float32Size))
	return
}
