package DxMsgPack

import (
	"time"
	"reflect"
	"github.com/suiyunonghen/DxValue/Coders"
	"fmt"
)

var(
	valueEncoders []Coders.EncoderFunc
)

func init()  {
	valueEncoders = []Coders.EncoderFunc{
		reflect.Bool:          encodeBoolValue,
		reflect.Int:           encodeInt64Value,
		reflect.Int8:          encodeInt64Value,
		reflect.Int16:         encodeInt64Value,
		reflect.Int32:         encodeInt64Value,
		reflect.Int64:         encodeInt64Value,
		reflect.Uint:          encodeInt64Value,
		reflect.Uint8:         encodeInt64Value,
		reflect.Uint16:        encodeInt64Value,
		reflect.Uint32:        encodeInt64Value,
		reflect.Uint64:        encodeInt64Value,
		reflect.Float32:       encodeFloat32Value,
		reflect.Float64:       encodeFloat64Value,
		reflect.Complex64:     encodeUnsupportedValue,
		reflect.Complex128:    encodeUnsupportedValue,
		reflect.Array:         encodeArrayValue,
		reflect.Chan:          encodeUnsupportedValue,
		reflect.Func:          encodeUnsupportedValue,
		reflect.Interface:     encodeInterfaceValue,
		reflect.Ptr:           encodeUnsupportedValue,
		reflect.Slice:         encodeSliceValue,
		reflect.String:        encodeStringValue,
		reflect.Struct:        encodeStructValue,
		reflect.UnsafePointer: encodeUnsupportedValue,
	}
}

func encodeBoolValue(encoder Coders.Encoder,v reflect.Value) error  {
	return encoder.(*MsgPackEncoder).EncodeBool(v.Bool())
}

func encodeInt64Value(encoder Coders.Encoder,v reflect.Value) error  {
	return encoder.(*MsgPackEncoder).EncodeInt(v.Int())
}

func encodeFloat32Value(encoder Coders.Encoder,v reflect.Value) error  {
	return encoder.(*MsgPackEncoder).EncodeFloat(float32(v.Float()))
}

func encodeFloat64Value(encoder Coders.Encoder,v reflect.Value) error  {
	return encoder.(*MsgPackEncoder).EncodeDouble(v.Float())
}

func encodeUnsupportedValue(encoder Coders.Encoder,v reflect.Value) error  {
	return fmt.Errorf("msgpack: Encode(unsupported %s)", v.Type())
}

func encodeStringValue(encoder Coders.Encoder,v reflect.Value) error  {
	return encoder.(*MsgPackEncoder).EncodeString(v.String())
}

func encodeInterfaceValue(encoder Coders.Encoder, v reflect.Value) error {
	if v.IsNil() {
		return encoder.(*MsgPackEncoder).WriteByte(0xc0) //null
	}
	v = v.Elem()
	if encodeFunc := encoder.(*MsgPackEncoder).GetEncoderFunc(v.Type());encodeFunc !=nil{
		return encodeFunc(encoder,v)
	}else if v.CanInterface(){
		return encoder.EncodeStand(v.Interface())
	}else {
		return  fmt.Errorf("msgpack: Encode(unsupported %s)", v.Type())
	}
}


func encodeStructValue(encoder Coders.Encoder, strct reflect.Value) error {
	structFields := Coders.Fields(strct.Type())
	mapLen := structFields.Len(msgPackName)
	msgEncoder := encoder.(*MsgPackEncoder)
	err := msgEncoder.EncodeMapLen(mapLen)
	if err != nil {
		return err
	}
	for i:=0;i<structFields.Len("");i++{
		f := structFields.Field(i)
		MarshName := f.MarshalName(msgPackName)
		if MarshName != "-" {
			err = msgEncoder.EncodeString(MarshName)
			if err!=nil{
				return err
			}
			f.EncodeValue(encoder,strct)
		}
	}

	return nil
}

func encodeArrayValue(encoder Coders.Encoder,v reflect.Value)(err error)  {
	arlen := v.Len()
	msgEncoder := encoder.(*MsgPackEncoder)
	switch {
	case arlen < 16: //1001XXXX|    N objects
		err = msgEncoder.WriteByte(byte(CodeFixedArrayLow) | byte(arlen))
	case arlen <= Max_map16_len:  //0xdc  |YYYYYYYY|YYYYYYYY|    N objects
		err = msgEncoder.WriteUint16(uint16(arlen),CodeArray16)
	default:
		if arlen > Max_map32_len{
			arlen = Max_map32_len
		}
		err = msgEncoder.WriteUint32(uint32(arlen),CodeArray32)
	}
	if err!=nil{
		return
	}
	for i := 0;i< arlen;i++ {
		av := v.Index(i)
		arrvalue := Coders.GetRealValue(&av)
		if encodeFunc := msgEncoder.GetEncoderFunc(arrvalue.Type());encodeFunc !=nil{
			err = encodeFunc(encoder,*arrvalue)
		}else if arrvalue.CanInterface(){
			err = encoder.EncodeStand(arrvalue.Interface())
		}
		if err!=nil{
			return
		}
	}
	return nil
}

func encodeSliceValue(encoder Coders.Encoder,v reflect.Value)(err error)  {
	if v.IsNil() {
		return encoder.(*MsgPackEncoder).WriteByte(0xc0) //null
	}else{
		return encodeArrayValue(encoder,v)
	}
}

func encodeErrorValue(e Coders.Encoder, v reflect.Value) error {
	if v.IsNil() {
		return e.(*MsgPackEncoder).WriteByte(0xc0) //null
	}
	return e.(*MsgPackEncoder).EncodeString(v.Interface().(error).Error())
}

func (coder *MsgPackEncoder)EncodeCustom()error  {
	return nil
}

func grow(b []byte, n int) []byte {
	if cap(b) >= n {
		return b[:n]
	}
	b = b[:cap(b)]
	b = append(b, make([]byte, n-len(b))...)
	return b
}

func encodeByteArrayValue(e Coders.Encoder, v reflect.Value)(err error) {
	btlen := v.Len()
	encoder := e.(*MsgPackEncoder)
	switch {
	case btlen <= Max_str8_len:
		encoder.buf[0] = byte(0xc4)
		encoder.buf[1] = byte(btlen)
		_,err = encoder.w.Write(encoder.buf[:2])
	case btlen <= Max_str16_len:
		err = encoder.WriteUint16(uint16(btlen),CodeBin16)
	default:
		if btlen > Max_str32_len{
			btlen = Max_str32_len
		}
		err = encoder.WriteUint32(uint32(btlen),CodeBin32)
	}
	if err!=nil{
		return
	}

	if v.CanAddr() {
		b := v.Slice(0, btlen).Bytes()
		_,err = encoder.w.Write(b)
		return
	}

	encoder.buf = grow(encoder.buf, btlen)
	reflect.Copy(reflect.ValueOf(encoder.buf), v)
	_,err = encoder.w.Write(encoder.buf)
	return
}

func encodeMapStringStringValue(e Coders.Encoder, v reflect.Value) error {
	encoder := e.(*MsgPackEncoder)
	if v.IsNil() {
		return encoder.WriteByte(0xc0) //null
	}

	if err := encoder.EncodeMapLen(v.Len()); err != nil {
		return err
	}

	m := v.Convert(Coders.MapStringStringType).Interface().(map[string]string)
	for mk, mv := range m {
		if err := encoder.EncodeString(mk); err != nil {
			return err
		}
		if err := encoder.EncodeString(mv); err != nil {
			return err
		}
	}

	return nil
}

func encodeMapIntStringValue(e Coders.Encoder, v reflect.Value) error {
	encoder := e.(*MsgPackEncoder)
	if v.IsNil() {
		return encoder.WriteByte(0xc0) //null
	}

	if err := encoder.EncodeMapLen(v.Len()); err != nil {
		return err
	}

	m := v.Convert(Coders.MapIntStringType).Interface().(map[int]string)
	for mk, mv := range m {
		if err := encoder.EncodeInt(int64(mk)); err != nil {
			return err
		}
		if err := encoder.EncodeString(mv); err != nil {
			return err
		}
	}

	return nil
}

func encodeMapInt64StringValue(e Coders.Encoder, v reflect.Value) error {
	encoder := e.(*MsgPackEncoder)
	if v.IsNil() {
		return encoder.WriteByte(0xc0) //null
	}

	if err := encoder.EncodeMapLen(v.Len()); err != nil {
		return err
	}

	m := v.Convert(Coders.MapIntStringType).Interface().(map[int64]string)
	for mk, mv := range m {
		if err := encoder.EncodeInt(mk); err != nil {
			return err
		}
		if err := encoder.EncodeString(mv); err != nil {
			return err
		}
	}

	return nil
}

func encodeMapStringInterfaceValue(e Coders.Encoder, v reflect.Value) error {
	encoder := e.(*MsgPackEncoder)
	if v.IsNil() {
		return encoder.WriteByte(0xc0) //null
	}

	if err := encoder.EncodeMapLen(v.Len()); err != nil {
		return err
	}

	m := v.Convert(Coders.MapStringInterfaceType).Interface().(map[string]interface{})
	for mk, mv := range m {
		if err := encoder.EncodeString(mk); err != nil {
			return err
		}
		if err := encoder.EncodeStand(mv); err != nil {
			return err
		}
	}

	return nil
}

func encodeMapIntInterfaceValue(e Coders.Encoder, v reflect.Value) error {
	encoder := e.(*MsgPackEncoder)
	if v.IsNil() {
		return encoder.WriteByte(0xc0) //null
	}

	if err := encoder.EncodeMapLen(v.Len()); err != nil {
		return err
	}

	m := v.Convert(Coders.MapStringInterfaceType).Interface().(map[int]interface{})
	for mk, mv := range m {
		if err := encoder.EncodeInt(int64(mk)); err != nil {
			return err
		}
		if err := encoder.EncodeStand(mv); err != nil {
			return err
		}
	}

	return nil
}

func encodeMapInt64InterfaceValue(e Coders.Encoder, v reflect.Value) error {
	encoder := e.(*MsgPackEncoder)
	if v.IsNil() {
		return encoder.WriteByte(0xc0) //null
	}

	if err := encoder.EncodeMapLen(v.Len()); err != nil {
		return err
	}

	m := v.Convert(Coders.MapStringInterfaceType).Interface().(map[int64]interface{})
	for mk, mv := range m {
		if err := encoder.EncodeInt(mk); err != nil {
			return err
		}
		if err := encoder.EncodeStand(mv); err != nil {
			return err
		}
	}

	return nil
}

func (encoder *MsgPackEncoder)EncodeMapLen(maplen int)(err error){
	if maplen <= Max_fixmap_len{   //fixmap
		err = encoder.WriteByte(0x80 | byte(maplen))
	}else if maplen <= Max_map16_len{
		//写入长度
		err = encoder.WriteUint16(uint16(maplen),CodeMap16)
	}else{
		if maplen > Max_map32_len{
			maplen = Max_map32_len
		}
		err = encoder.WriteUint32(uint32(maplen),CodeMap32)
	}
	return
}

func (coder *MsgPackEncoder)GetEncoderFunc(typ reflect.Type) Coders.EncoderFunc {
	kind := typ.Kind()

	if typ == errorType {
		return encodeErrorValue
	}

	switch kind {
	case reflect.Ptr:
		return func(e Coders.Encoder, v reflect.Value) error {
			if v.IsNil() {
				return e.(*MsgPackEncoder).WriteByte(0xc0) //null
			}
			encoderFunc := coder.GetEncoderFunc(typ.Elem())
			return encoderFunc(e, v.Elem())
		}
	case reflect.Slice:
		if typ.Elem().Kind() == reflect.Uint8 {
			return func(encoder Coders.Encoder, value reflect.Value) error {
				return encoder.(*MsgPackEncoder).EncodeBinary(value.Bytes())
			}
		}
	case reflect.Array:
		if typ.Elem().Kind() == reflect.Uint8 {
			return encodeByteArrayValue
		}
	case reflect.Map:
		switch typ.Key() {
		case Coders.StringType:
			switch typ.Elem(){
			case Coders.StringType:
				return encodeMapStringStringValue
			case Coders.InterfaceType:
				return encodeMapStringInterfaceValue
			}
		case Coders.Int64Type:
			switch typ.Elem(){
			case Coders.StringType:
				return encodeMapInt64StringValue
			case Coders.InterfaceType:
				return encodeMapInt64InterfaceValue
			}
		case Coders.IntType:
			switch typ.Elem(){
			case Coders.StringType:
				return encodeMapIntStringValue
			case Coders.InterfaceType:
				return encodeMapIntInterfaceValue
			}
		}

	}
	result := valueEncoders[kind]
	if result == nil{
		result = encodeUnsupportedValue
	}
	return result
}

func (encoder *MsgPackEncoder)encodeInterfaceArr(arr []interface{})(err error)  {
	arlen := len(arr)
	switch {
	case arlen < 16: //1001XXXX|    N objects
		err = encoder.WriteByte(byte(CodeFixedArrayLow) | byte(arlen))
	case arlen <= Max_map16_len:  //0xdc  |YYYYYYYY|YYYYYYYY|    N objects
		encoder.WriteUint16(uint16(arlen),CodeArray16)
	default:
		if arlen > Max_map32_len{
			arlen = Max_map32_len
		}
		encoder.WriteUint32(uint32(arlen),CodeArray32)
	}

	for i := 0;i <= arlen - 1;i++{
		if arr[i] == nil{
			err = encoder.WriteByte(0xc0) //null
		}else{
			err = encoder.EncodeStand(arr[i])
		}
		if err != nil{
			return
		}
	}
	return err
}


func (encoder *MsgPackEncoder)EncodeStand(v interface{})(error)  {
	switch value := v.(type) {
	case *string:
		return encoder.EncodeString(*value)
	case string:
		return encoder.EncodeString(value)
	case *[]interface{}:
		if value!=nil{
			return encoder.WriteByte(byte(CodeNil))
		}else{
			return encoder.encodeInterfaceArr(*value)
		}
	case []interface{}:
		return encoder.encodeInterfaceArr(value)
	case *time.Time:
		return encoder.EncodeTime(*value)
	case time.Time:
		return encoder.EncodeTime(value)
	case *int8:
		return encoder.EncodeInt(int64(*value))
	case int8:
		return encoder.EncodeInt(int64(value))
	case *int16:
		return encoder.EncodeInt(int64(*value))
	case int16:
		return encoder.EncodeInt(int64(value))
	case *int32:
		return encoder.EncodeInt(int64(*value))
	case int32:
		return encoder.EncodeInt(int64(value))
	case *int64:
		return encoder.EncodeInt(*value)
	case int64:
		return encoder.EncodeInt(value)
	case *uint8:
		return encoder.EncodeInt(int64(*value))
	case uint8:
		return encoder.EncodeInt(int64(value))
	case *uint16:
		return encoder.EncodeInt(int64(*value))
	case uint16:
		return encoder.EncodeInt(int64(value))
	case *uint32:
		return encoder.EncodeInt(int64(*value))
	case uint32:
		return encoder.EncodeInt(int64(value))
	case *uint64:
		return encoder.EncodeInt(int64(*value))
	case uint64:
		return encoder.EncodeInt(int64(value))
	case *float32:
		return encoder.EncodeFloat(*value)
	case *float64:
		return encoder.EncodeDouble(*value)
	case float32:
		return encoder.EncodeFloat(value)
	case float64:
		return encoder.EncodeDouble(value)
	case *bool:
		return encoder.EncodeBool(*value)
	case bool:
		return encoder.EncodeBool(value)
	case *[]byte:
		return encoder.EncodeBinary(*value)
	case *map[string]interface{}:
		return encoder.encodeStrMapFunc(value)
	case map[string]interface{}:
		return encoder.encodeStrMapFunc(&value)
	case *map[int]interface{}:
		return encoder.encodeIntMapFunc(value)
	case map[int]interface{}:
		return encoder.encodeIntMapFunc(&value)
	case *map[int64]interface{}:
		return encoder.encodeInt64MapFunc(value)
	case map[int64]interface{}:
		return encoder.encodeInt64MapFunc(&value)
	case *map[string]string:
		return encoder.encodeStrStrMapFunc(value)
	case map[string]string:
		return encoder.encodeStrStrMapFunc(&value)
	case *time.Duration:
		return encoder.EncodeInt(int64(*value))
	default:
		v := reflect.ValueOf(value)
		rv := Coders.GetRealValue(&v)
		if rv == nil{
			return encoder.WriteByte(0xc0) //null
		}
		switch rv.Kind(){
		case reflect.Struct:
			if rv.Type() == Coders.TimeType{
				return encoder.EncodeTime(rv.Interface().(time.Time))
			}
			structFields := Coders.Fields(rv.Type())
			mapLen := structFields.Len(msgPackName)
			err := encoder.EncodeMapLen(mapLen)
			if err != nil {
				return err
			}
			for i:=0;i<structFields.Len("");i++{
				f := structFields.Field(i)
				MarshName := f.MarshalName(msgPackName)
				if MarshName != "-" {
					if err = encoder.EncodeString(MarshName);err!=nil{
						return err
					}
					if err = f.EncodeValue(encoder,*rv);err!=nil{
						return err
					}

				}
			}
		case reflect.Map:
		case reflect.Slice,reflect.Array:

		}
	}
	return nil
}
