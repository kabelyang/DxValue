/*
DxValue的Record记录集对象
可以用来序列化反序列化Json,MsgPack等，并提供一系列的操作函数
Autor: 不得闲
QQ:75492895
 */
package DxValue

import (
	"unsafe"
	"bytes"
	"reflect"
	"github.com/suiyunonghen/DxCommonLib"
	"strings"
	"math"
	"strconv"
	"io/ioutil"
	"io"
	"bufio"
	"os"
)

/******************************************************
*  DxRecord
******************************************************/
type(
		DxRecord		struct{
		DxBaseValue
		fRecords		map[string]*DxBaseValue
		PathSplitChar	byte
		}
)


func (r *DxRecord)ClearValue()  {
	if r.fRecords != nil{
		for _,v := range r.fRecords{
			v.ClearValue()
			v.fParent = nil
		}
	}
	if r.fRecords == nil || len(r.fRecords) > 0{
		r.fRecords = make(map[string]*DxBaseValue,32)
	}
}

func (r *DxRecord)splitPathFields(charrune rune) bool {
	return charrune == rune(r.PathSplitChar)
}

func (r *DxRecord)NewRecord(keyName string)(rec *DxRecord)  {
	if value,ok := r.fRecords[keyName];ok && value != nil{
		if value.fValueType == DVT_Record{
			rec = (*DxRecord)(unsafe.Pointer(value))
			rec.ClearValue()
			rec.fParent = &r.DxBaseValue
			return
		}
		value.fParent = nil
	}
	rec = new(DxRecord)
	rec.fValueType = DVT_Record
	rec.PathSplitChar = r.PathSplitChar
	rec.fRecords = make(map[string]*DxBaseValue,32)
	r.fRecords[keyName] = &rec.DxBaseValue
	rec.fParent = &r.DxBaseValue
	return
}

func (r *DxRecord)Find(keyName string)*DxBaseValue  {
	if v,ok := r.fRecords[keyName];ok{
		return v
	}
	return nil
}



func (r *DxRecord)ForcePath(path string,v interface{}) {
	fields := strings.FieldsFunc(path,r.splitPathFields)
	vlen := len(fields)
	if vlen == 0{
		return
	}
	rec := r
	for i := 0;i<vlen - 1;i++{
		vbase := rec.Find(fields[i])
		if vbase != nil && vbase.fValueType == DVT_Record{
			rec = (*DxRecord)(unsafe.Pointer(vbase))
		}else{
			rec = rec.NewRecord(fields[i])
		}
	}
	rec.SetValue(fields[vlen - 1],v)
}

func (r *DxRecord)NewArray(keyName string)(arr *DxArray)  {
	if value,ok := r.fRecords[keyName];ok && value != nil{
		if value.fValueType == DVT_Array{
			arr = (*DxArray)(unsafe.Pointer(value))
			arr.ClearValue()
			arr.fParent = &r.DxBaseValue
			return
		}
		value.fParent = nil
	}
	arr = new(DxArray)
	arr.fValueType = DVT_Array
	arr.fParent = &r.DxBaseValue
	r.fRecords[keyName] = &arr.DxBaseValue
	return
}

func (r *DxRecord)SetInt(KeyName string,v int)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_Int:
			(*DxIntValue)(unsafe.Pointer(value)).fvalue = v
			return
		case DVT_Int32:
			if v <= math.MaxInt32 && v >= math.MinInt32{
				(*DxInt32Value)(unsafe.Pointer(value)).fvalue = int32(v)
				return
			}
		case DVT_Int64:
			(*DxInt64Value)(unsafe.Pointer(value)).fvalue = int64(v)
			return
		}
		value.fParent = nil
	}
	var m DxIntValue
	m.fvalue = v
	m.fValueType = DVT_Int
	m.fParent = &r.DxBaseValue
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetInt32(KeyName string,v int32)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_Int:
			(*DxIntValue)(unsafe.Pointer(value)).fvalue = int(v)
			return
		case DVT_Int32:
			(*DxInt32Value)(unsafe.Pointer(value)).fvalue = v
			return
		case DVT_Int64:
			(*DxInt64Value)(unsafe.Pointer(value)).fvalue = int64(v)
			return
		}
		value.fParent = nil
	}
	var m DxInt32Value
	m.fvalue = v
	m.fValueType = DVT_Int32
	m.fParent = &r.DxBaseValue
	r.fRecords[KeyName] = &m.DxBaseValue
}


func (r *DxRecord)SetInt64(KeyName string,v int64)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_Int:
			if DxCommonLib.IsAmd64 || v <= math.MaxInt32 && v >= math.MinInt32{
				(*DxIntValue)(unsafe.Pointer(value)).fvalue = int(v)
				return
			}
		case DVT_Int32:
			if v <= math.MaxInt32 && v >= math.MinInt32{
				(*DxInt32Value)(unsafe.Pointer(value)).fvalue = int32(v)
				return
			}
			return
		case DVT_Int64:
			(*DxInt64Value)(unsafe.Pointer(value)).fvalue = v
			return
		}
		value.fParent = nil
	}
	var m DxInt64Value
	m.fvalue = v
	m.fValueType = DVT_Int64
	m.fParent = &r.DxBaseValue
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetBool(KeyName string,v bool)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Bool {
			(*DxBoolValue)(unsafe.Pointer(value)).fvalue = v
			return
		}
		value.fParent = nil
	}
	var m DxBoolValue
	m.fvalue = v
	m.fParent = &r.DxBaseValue
	m.fValueType = DVT_Bool
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetNull(KeyName string)  {
	r.fRecords[KeyName] = nil
}



func (r *DxRecord)SetFloat(KeyName string,v float32)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Float{
			(*DxFloatValue)(unsafe.Pointer(value)).fvalue = v
			return
		}else if value.fValueType == DVT_Double{
			(*DxDoubleValue)(unsafe.Pointer(value)).fvalue = float64(v)
			return
		}
		value.fParent = nil
	}
	var m DxFloatValue
	m.fvalue = v
	m.fParent = &r.DxBaseValue
	m.fValueType = DVT_Float
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetDouble(KeyName string,v float64)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Double{
			(*DxDoubleValue)(unsafe.Pointer(value)).fvalue = v
			return
		}else if value.fValueType == DVT_Float{
			if v <= math.MaxFloat32 && v >= math.MinInt32{
				(*DxFloatValue)(unsafe.Pointer(value)).fvalue = float32(v)
				return
			}
		}
		value.fParent = nil
	}
	var m DxDoubleValue
	m.fvalue = v
	m.fParent = &r.DxBaseValue
	m.fValueType = DVT_Double
	r.fRecords[KeyName] = &m.DxBaseValue
}



func (r *DxRecord)SetString(KeyName string,v string)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Double{
			(*DxStringValue)(unsafe.Pointer(value)).fvalue = v
			return
		}
		value.fParent = nil
	}
	var m DxStringValue
	m.fvalue = v
	m.fValueType = DVT_String
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)SetBinary(KeyName string,v []byte,reWrite bool)  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Binary{
			bv := (*DxBinaryValue)(unsafe.Pointer(value))
			if reWrite{
				bv.SetBinary(v,false)
			}else{
				bv.Append(v)
			}
			return
		}
		value.fParent = nil
	}
	var m DxBinaryValue
	m.Append(v)
	m.fParent = &r.DxBaseValue
	m.fValueType = DVT_Binary
	r.fRecords[KeyName] = &m.DxBaseValue
}

func (r *DxRecord)AsBytes(keyName string)[]byte  {
	if value,ok := r.fRecords[keyName];ok && value != nil{
		bt,_ := value.AsBytes()
		return bt
	}
	return nil
}

func (r *DxRecord)Bytes()[]byte  {
	var buffer bytes.Buffer
	buffer.WriteByte('{')
	isFirst := true
	for k,v := range r.fRecords{
		if !isFirst{
			buffer.WriteString(`,"`)
		}else{
			isFirst = false
			buffer.WriteByte('"')
		}
		buffer.WriteString(k)
		buffer.WriteString(`":`)
		if v != nil{
			vt := v.fValueType
			if vt == DVT_String || vt == DVT_Binary{
				buffer.WriteByte('"')
			}
			buffer.WriteString(v.ToString())
			if vt == DVT_String || vt == DVT_Binary{
				buffer.WriteByte('"')
			}
		}else{
			buffer.WriteString("null")
		}
	}
	buffer.WriteByte('}')
	return buffer.Bytes()
}

func (r *DxRecord)findPathNode(path string)(rec *DxRecord,keyName string)  {
	fields := strings.FieldsFunc(path,r.splitPathFields)
	vlen := len(fields)
	if vlen == 0{
		return nil,""
	}
	rec = r
	for i := 0;i < vlen - 1;i++{
		vbase := rec.Find(fields[i])
		if vbase != nil && vbase.fValueType == DVT_Record{
			rec = (*DxRecord)(unsafe.Pointer(vbase))
		}else{
			return nil,""
		}
	}
	return rec,fields[vlen - 1]
}

func (r *DxRecord)AsBytesByPath(Path string)[]byte  {
	rec,keyName := r.findPathNode(Path)
	if rec != nil {
		if keyName != ""{
			return rec.AsBytes(keyName)
		}
	}
	return nil
}

func getBaseType(vt reflect.Type)reflect.Kind  {
	if vt.Kind() == reflect.Ptr{
		return getBaseType(vt.Elem())
	}
	return vt.Kind()
}

func getRealValue(v *reflect.Value)*reflect.Value  {
	if !v.IsValid(){
		return nil
	}
	if v.Kind() == reflect.Ptr{
		if !v.IsNil(){
			va := v.Elem()
			return getRealValue(&va)
		}else{
			return nil
		}
	}
	return v
}


func (r *DxRecord)SetRecordValue(keyName string,v *DxRecord) {
	if v != nil && v.fParent != nil {
		panic("Must Set A Single Record(no Parent)")
	}
	if value, ok := r.fRecords[keyName]; ok && value != nil {
		if  value.fValueType == DVT_Record {
			nrec := (*DxRecord)(unsafe.Pointer(value))
			nrec.fParent = nil
			*nrec = *v
			nrec.fParent = &r.DxBaseValue
			return
		}
		value.fParent = nil
	}
	if v != nil {
		r.fRecords[keyName] = &v.DxBaseValue
		v.fParent = &r.DxBaseValue
	}
}

func (r *DxRecord)SetArray(KeyName string,v *DxArray)  {
	if v != nil && v.fParent != nil {
		panic("Must Set A Single Array(no Parent)")
	}
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Array{
			arr := (*DxArray)(unsafe.Pointer(value))
			arr.fParent = nil
			*arr = *v
			arr.fParent = &r.DxBaseValue
			return
		}
		value.fParent = nil
	}
	if v!=nil{
		r.fRecords[KeyName] = &v.DxBaseValue
		v.fParent = &r.DxBaseValue
	}
}

func (r *DxRecord)AsBaseValue(keyName string)*DxBaseValue{
	if r.fRecords != nil{
		return r.fRecords[keyName]
	}
	return nil
}

func (r *DxRecord)SetValue(keyName string,v interface{})  {
	if v == nil{
		r.fRecords[keyName] = nil
		return
	}
	switch value := v.(type) {
	case int: r.SetInt(keyName,value)
	case int32: r.SetInt32(keyName,value)
	case int64: r.SetInt64(keyName,value)
	case int8: r.SetInt(keyName,int(value))
	case uint8: r.SetInt(keyName,int(value))
	case int16: r.SetInt(keyName,int(value))
	case uint16: r.SetInt(keyName,int(value))
	case uint32: r.SetInt(keyName,int(value))
	case *int: r.SetInt(keyName,*value)
	case *int32: r.SetInt32(keyName,*value)
	case *int64: r.SetInt64(keyName,*value)
	case *int8: r.SetInt(keyName,int(*value))
	case *uint8: r.SetInt(keyName,int(*value))
	case *int16: r.SetInt(keyName,int(*value))
	case *uint16: r.SetInt(keyName,int(*value))
	case *uint32: r.SetInt(keyName,int(*value))
	case string: r.SetString(keyName,value)
	case []byte: r.SetBinary(keyName,value,true)
	case *[]byte: r.SetBinary(keyName,*value,true)
	case bool: r.SetBool(keyName,value)
	case *bool: r.SetBool(keyName,*value)
	case *string: r.SetString(keyName,*value)
	case float32: r.SetFloat(keyName,value)
	case float64: r.SetDouble(keyName,value)
	case *float32: r.SetFloat(keyName,*value)
	case *float64: r.SetDouble(keyName,*value)
	case *DxRecord: r.SetRecordValue(keyName,value)
	case DxRecord: r.SetRecordValue(keyName,&value)
	case DxArray:  r.SetArray(keyName,&value)
	case *DxArray: r.SetArray(keyName,value)
	case DxInt64Value: r.SetInt64(keyName,value.fvalue)
	case *DxInt64Value: r.SetInt64(keyName,value.fvalue)
	case DxInt32Value: r.SetInt32(keyName,value.fvalue)
	case *DxInt32Value: r.SetInt32(keyName,value.fvalue)
	case DxFloatValue: r.SetFloat(keyName,value.fvalue)
	case *DxFloatValue: r.SetFloat(keyName,value.fvalue)
	case DxDoubleValue: r.SetDouble(keyName,value.fvalue)
	case *DxDoubleValue: r.SetDouble(keyName,value.fvalue)
	case DxBoolValue: r.SetBool(keyName,value.fvalue)
	case *DxBoolValue: r.SetBool(keyName,value.fvalue)
	case DxIntValue: r.SetInt(keyName,value.fvalue)
	case *DxIntValue: r.SetInt(keyName,value.fvalue)
	case DxStringValue: r.SetString(keyName,value.fvalue)
	case *DxStringValue: r.SetString(keyName,value.fvalue)
	case DxBinaryValue:  r.SetBinary(keyName,value.Bytes(),true)
	case *DxBinaryValue:  r.SetBinary(keyName,value.Bytes(),true)
	default:
		reflectv := reflect.ValueOf(v)
		rv := getRealValue(&reflectv)
		if rv == nil{
			if _,ok := r.fRecords[keyName];!ok{
				r.fRecords[keyName] = nil
			}
			return
		}
		switch rv.Kind(){
		case reflect.Struct:
			rec := r.NewRecord(keyName)
			rtype := rv.Type()
			for i := 0;i < rtype.NumField();i++{
				sfield := rtype.Field(i)
				fv := rv.Field(i)
				fieldvalue := getRealValue(&fv)
				if fieldvalue != nil{
					switch fieldvalue.Kind() {
					case reflect.Int,reflect.Uint32:
						rec.SetInt(sfield.Name,int(fieldvalue.Int()))
					case reflect.Bool:
						rec.SetBool(sfield.Name,fieldvalue.Bool())
					case reflect.Int64:
						rec.SetInt64(sfield.Name,fieldvalue.Int())
					case reflect.Int32,reflect.Int8,reflect.Int16,reflect.Uint8,reflect.Uint16:
						rec.SetInt32(sfield.Name,int32(fieldvalue.Int()))
					case reflect.Float32:
						rec.SetFloat(sfield.Name,float32(fieldvalue.Float()))
					case reflect.Float64:
						rec.SetDouble(sfield.Name,fieldvalue.Float())
					case reflect.String:
						rec.SetString(sfield.Name,fieldvalue.String())
					default:
						if fieldvalue.CanInterface(){
							rec.SetValue(sfield.Name,fieldvalue.Interface())
						}
					}
				}
			}
		case reflect.Map:
			rec := r.NewRecord(keyName)
			mapkeys := rv.MapKeys()
			if len(mapkeys) == 0{
				return
			}
			kv := mapkeys[0]
			if getBaseType(kv.Type()) != reflect.String{
				panic("Invalidate Record Key")
			}
			rvalue := rv.MapIndex(mapkeys[0])
			//获得Value类型
			valueKind := getBaseType(rvalue.Type())
			for _,kv = range mapkeys{
				rvalue = rv.MapIndex(kv)
				prvalue := getRealValue(&rvalue)
				if prvalue != nil{
					switch valueKind {
					case reflect.Int,reflect.Uint32:
						rec.SetInt(kv.String(),int(prvalue.Int()))
					case reflect.Bool:
						rec.SetBool(kv.String(),prvalue.Bool())
					case reflect.Int64:
						rec.SetInt64(kv.String(),prvalue.Int())
					case reflect.Int32,reflect.Int8,reflect.Int16,reflect.Uint8,reflect.Uint16:
						rec.SetInt32(kv.String(),int32(prvalue.Int()))
					case reflect.Float32:
						rec.SetFloat(kv.String(),float32(prvalue.Float()))
					case reflect.Float64:
						rec.SetDouble(kv.String(),prvalue.Float())
					case reflect.String:
						rec.SetString(kv.String(),prvalue.String())
					default:
						if prvalue.CanInterface(){
							rec.SetValue(kv.String(),prvalue.Interface())
						}
					}
				}
			}
		case reflect.Slice,reflect.Array:
			arr := r.NewArray(keyName)
			vlen := rv.Len()
			for i := 0;i< vlen;i++{
				av := rv.Index(i)
				arrvalue := getRealValue(&av)
				switch arrvalue.Kind() {
				case reflect.Int,reflect.Uint32:
					arr.SetInt(i,int(arrvalue.Int()))
				case reflect.Bool:
					arr.SetBool(i,arrvalue.Bool())
				case reflect.Int64:
					arr.SetInt64(i,arrvalue.Int())
				case reflect.Int32,reflect.Int8,reflect.Int16,reflect.Uint8,reflect.Uint16:
					arr.SetInt32(i,int32(arrvalue.Int()))
				case reflect.Float32:
					arr.SetFloat(i,float32(arrvalue.Float()))
				case reflect.Float64:
					arr.SetDouble(i,arrvalue.Float())
				case reflect.String:
					arr.SetString(i,arrvalue.String())
				default:
					if arrvalue.CanInterface(){
						arr.SetValue(i,arrvalue.Interface())
					}
				}
			}
		}
	}
}


func (r *DxRecord)KeyValueType(KeyName string)DxValueType  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		return value.fValueType
	}
	return DVT_Null
}

func (r *DxRecord)AsInt32(KeyName string,defavalue int32)int32  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_Int: return int32((*DxIntValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Int32: return int32((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Int64: return int32((*DxInt64Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Bool:
			if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
				return 1
			}else{
				return 0
			}
		case DVT_Double:return int32((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return int32((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		default:
			panic("can not convert Type to int32")
		}
	}
	return defavalue
}

func (r *DxRecord)AsInt32ByPath(path string,defavalue int32)int32  {
	rec,keyName := r.findPathNode(path)
	if rec != nil && keyName != ""{
		return rec.AsInt32(keyName,defavalue)
	}
	return defavalue
}

func (r *DxRecord)AsIntByPath(path string,defavalue int)int  {
	rec,keyName := r.findPathNode(path)
	if rec != nil && keyName != ""{
		return rec.AsInt(keyName,defavalue)
	}
	return defavalue
}

func (r *DxRecord)AsInt(KeyName string,defavalue int)int  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_Int: return (*DxIntValue)(unsafe.Pointer(value)).fvalue
		case DVT_Int32: return int((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Int64: return int((*DxInt64Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Bool:
			if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
				return 1
			}else{
				return 0
			}
		case DVT_Double:return int((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return int((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		default:
			panic("can not convert Type to int")
		}
	}
	return defavalue
}

func (r *DxRecord)AsInt64ByPath(path string,defavalue int64)int64  {
	rec,keyName := r.findPathNode(path)
	if rec != nil && keyName != ""{
		return rec.AsInt64(keyName,defavalue)
	}
	return defavalue
}

func (r *DxRecord)AsInt64(KeyName string,defavalue int64)int64  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_Int: return int64((*DxIntValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Int32: return int64((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Int64: return int64((*DxInt64Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Bool:
			if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
				return 1
			}else{
				return 0
			}
		case DVT_Double:return int64((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return int64((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		default:
			panic("can not convert Type to int64")
		}
	}
	return defavalue
}

func (r *DxRecord)AsBoolByPath(path string,defavalue bool)bool  {
	rec,keyName := r.findPathNode(path)
	if rec != nil && keyName != ""{
		return rec.AsBool(keyName,defavalue)
	}
	return defavalue
}

func (r *DxRecord)AsBool(KeyName string,defavalue bool)bool  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_Int: return (*DxIntValue)(unsafe.Pointer(value)).fvalue != 0
		case DVT_Int32: return (*DxInt32Value)(unsafe.Pointer(value)).fvalue != 0
		case DVT_Int64: return (*DxInt64Value)(unsafe.Pointer(value)).fvalue != 0
		case DVT_Bool: return bool((*DxBoolValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Double:return float64((*DxDoubleValue)(unsafe.Pointer(value)).fvalue) != 0
		case DVT_Float:return float32((*DxFloatValue)(unsafe.Pointer(value)).fvalue) != 0
		default:
			panic("can not convert Type to Bool")
		}
	}
	return defavalue
}


func (r *DxRecord)AsFloatByPath(path string,defavalue float32)float32  {
	rec,keyName := r.findPathNode(path)
	if rec != nil && keyName != ""{
		return rec.AsFloat(keyName,defavalue)
	}
	return defavalue
}

func (r *DxRecord)AsFloat(KeyName string,defavalue float32)float32  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_Int: return float32((*DxIntValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Int32: return float32((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Int64: return float32((*DxInt64Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Bool:
		    if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
		    	return 1
			}
			return 0
		case DVT_Double:return float32((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return (*DxFloatValue)(unsafe.Pointer(value)).fvalue
		default:
			panic("can not convert Type to Float")
		}
	}
	return defavalue
}


func (r *DxRecord)AsDoubleByPath(path string,defavalue float64)float64  {
	rec,keyName := r.findPathNode(path)
	if rec != nil && keyName != ""{
		return rec.AsDouble(keyName,defavalue)
	}
	return defavalue
}

func (r *DxRecord)AsDouble(KeyName string,defavalue float64)float64  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		switch value.fValueType {
		case DVT_Int: return float64((*DxIntValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Int32: return float64((*DxInt32Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Int64: return float64((*DxInt64Value)(unsafe.Pointer(value)).fvalue)
		case DVT_Bool:
			if (*DxBoolValue)(unsafe.Pointer(value)).fvalue{
				return 1
			}
			return 0
		case DVT_Double:return float64((*DxDoubleValue)(unsafe.Pointer(value)).fvalue)
		case DVT_Float:return float64((*DxFloatValue)(unsafe.Pointer(value)).fvalue)
		default:
			panic("can not convert Type to Double")
		}
	}
	return defavalue
}

func (r *DxRecord)AsStringByPath(path string,defavalue string)string  {
	rec,keyName := r.findPathNode(path)
	if rec != nil && keyName != ""{
		return rec.AsString(keyName,defavalue)
	}
	return defavalue
}

func (r *DxRecord)AsString(KeyName string,defavalue string)string  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		return value.ToString()
	}
	return defavalue
}

func (r *DxRecord)AsRecordByPath(path string)*DxRecord  {
	rec,keyName := r.findPathNode(path)
	if rec != nil && keyName != ""{
		return rec.AsRecord(keyName)
	}
	return nil
}

func (r *DxRecord)AsRecord(KeyName string)*DxRecord  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Record{
			return (*DxRecord)(unsafe.Pointer(value))
		}
		panic("not Record Value")
	}
	return nil
}

func (r *DxRecord)AsArray(KeyName string)*DxArray  {
	if value,ok := r.fRecords[KeyName];ok && value != nil{
		if value.fValueType == DVT_Array{
			return (*DxArray)(unsafe.Pointer(value))
		}
		panic("not Array Value")
	}
	return nil
}

func (r *DxRecord)Length()int  {
	if r.fRecords != nil{
		return len(r.fRecords)
	}
	return 0
}

func (r *DxRecord)Contains(keyName string)bool  {
	if r.fRecords != nil{
		if vr,vk := r.findPathNode(keyName);vr!=nil && vr.fRecords != nil{
			_,ok := vr.fRecords[vk]
			return ok
		}
	}
	return false
}

func (r *DxRecord)Remove(keyOrPath string)  {
	if r.fRecords != nil{
		if vr,vk := r.findPathNode(keyOrPath);vr!=nil && vr.fRecords != nil{
			if v,ok := vr.fRecords[vk];ok{
				if v != nil{
					v.ClearValue()
				}
				delete(vr.fRecords,vk)
			}
		}
	}
}

func (r *DxRecord)Range(iteafunc func(keyName string,value *DxBaseValue)bool){
	if r.fRecords != nil && iteafunc!=nil{
		for k,v := range r.fRecords{
			if !iteafunc(k,v){
				return
			}
		}
	}
}

func (r *DxRecord)ToString()string  {
	return DxCommonLib.FastByte2String(r.Bytes())
}

func (r *DxRecord)parserValue(keyName string, b []byte,ConvertEscape bool)(parserlen int, err error)  {
	blen := len(b)
	i := 0
	valuestart := -1
	validCharIndex := -1
	startValue := false
	for i<blen {
		if !IsSpace(b[i]){
			switch b[i] {
			case ':':
				startValue = true
				//valuestart = i //自己记录有效的开始位置，和有效的结束位置，省去一个trim
			case '{':
				var rec DxRecord
				rec.PathSplitChar = r.PathSplitChar
				rec.fValueType = DVT_Record
				rec.fRecords = make(map[string]*DxBaseValue,32)
				if parserlen,err = rec.JsonParserFromByte(b[i:blen],ConvertEscape);err == nil{
					r.SetRecordValue(keyName,&rec)
				}
				parserlen+=2 //会多解析一个{
				return
			case '[':
				var arr DxArray
				arr.fValueType = DVT_Array
				if parserlen,err = arr.JsonParserFromByte(b[i:],ConvertEscape);err == nil{
					r.SetArray(keyName,&arr)
				}
				parserlen+=2
				return
			case ',','}':
				//bvalue := bytes.Trim(b[valuestart + 1:i]," \r\n\t")
				bvalue := b[valuestart: validCharIndex+1]
				if len(bvalue) == 0{
					return i,ErrInvalidateJson
				}
				if bytes.IndexByte(bvalue,'.') > -1{
					if vf,err := strconv.ParseFloat(DxCommonLib.FastByte2String(bvalue),64);err!=nil{
						return i,ErrInvalidateJson
					}else{
						r.SetDouble(keyName,vf)
					}
				}else {
					st := DxCommonLib.FastByte2String(bvalue)
					if st == "true" || strings.ToUpper(st) == "TRUE"{
						r.SetBool(keyName,true)
					}else if st == "false" || strings.ToUpper(st) == "FALSE"{
						r.SetBool(keyName,false)
					}else if st == "null" || strings.ToUpper(st) == "NULL"{
						r.SetNull(keyName)
					}else{
						if vf,err := strconv.Atoi(st);err!=nil{
							return i,ErrInvalidateJson
						}else{
							r.SetInt(keyName,vf)
						}
					}
				}
				return i,nil
			case '"': //string
				plen := bytes.IndexByte(b[i+1:blen],'"')
				if plen > -1{
					bvalue := b[i+1:plen+i+1]
					st := ""
					if ConvertEscape{
						st = DxCommonLib.ParserEscapeStr(bvalue)
					}else{
						st = DxCommonLib.FastByte2String(bvalue)
					}
					r.SetString(keyName,st)
					return plen + i + 2,nil
				}
				return i,ErrInvalidateJson
			default:
				if !startValue && valuestart == -1{
					return i,ErrInvalidateJson
				}
				if valuestart == -1{
					valuestart = i
					startValue = false
				}else{
					validCharIndex = i
				}
			}
		}
		i += 1
	}
	return blen,ErrInvalidateJson
}

func (r *DxRecord)LoadJsonFile(fileName string,ConvertEscape bool)error  {
	databytes, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		return err
	}
	if databytes[0] == 0xEF && databytes[1] == 0xBB && databytes[2] == 0xBF{//BOM
		databytes = databytes[3:]
	}
	_,err = r.JsonParserFromByte(databytes,ConvertEscape)
	return err
}

func (r *DxRecord)SaveJsonWriter(w io.Writer)error  {
	writer := bufio.NewWriter(w)
	err := writer.WriteByte('{')
	if err != nil{
		return err
	}
	isFirst := true
	for k,v := range r.fRecords{
		if !isFirst{
			_,err = writer.WriteString(`,"`)
		}else{
			isFirst = false
			err = writer.WriteByte('"')
		}
		if err != nil{
			return err
		}
		_,err = writer.WriteString(k)
		if err!=nil{
			return err
		}
		_, err = writer.WriteString(`":`)
		if err!=nil{
			return err
		}
		if v != nil{
			vt := v.fValueType
			if vt == DVT_String || vt == DVT_Binary{
				err = writer.WriteByte('"')
			}
			if err != nil{
				return err
			}
			_,err = writer.WriteString(v.ToString())
			if err == nil && (vt == DVT_String || vt == DVT_Binary){
				err = writer.WriteByte('"')
			}
		}else{
			_,err = writer.WriteString("null")
		}
		if err != nil{
			return err
		}
	}
	writer.WriteByte('}')
	err = writer.Flush()
	return err
}

func (r *DxRecord)SaveJsonFile(fileName string,BOMFile bool)error  {
	if file,err := os.OpenFile(fileName,os.O_CREATE | os.O_TRUNC,0644);err == nil{
		defer file.Close()
		if BOMFile{
			file.Write([]byte{0xEF,0xBB,0xBF})
		}
		return r.SaveJsonWriter(file)
	}else{
		return err
	}
}

func (r *DxRecord)LoadJsonReader(reader io.Reader)error  {
	return nil
}

func (r *DxRecord)JsonParserFromByte(JsonByte []byte,ConvertEscape bool)(parserlen int, err error)  {
	i := 0
	r.ClearValue()
	objStart := false
	keyStart := false
	btlen := len(JsonByte)
	plen := -1
	keyName := ""
	for i < btlen{
		if IsSpace(JsonByte[i]){
			i++
			continue
		}
		if !objStart && JsonByte[i] != '{' {
			return 0,ErrInvalidateJson
		}
		switch JsonByte[i]{
		case '{':
			objStart = true
			keyStart = true
		case '}':
			if keyStart{
				return i,ErrInvalidateJson
			}
			objStart = false
			return i,nil
		case '"': //keyName
			if keyStart{
				//获取string
				plen = bytes.IndexByte(JsonByte[i+1:btlen],'"')
				if plen > -1{
					keyName = DxCommonLib.FastByte2String(JsonByte[i+1:i+1+plen])
				}
				i += plen+2
				keyStart = false
				//解析Value
				if ilen,err := r.parserValue(keyName,JsonByte[i:btlen],ConvertEscape);err!=nil{
					return ilen + i,err
				}else{
					i += ilen
					continue
				}
			}
		case ',': //next key
			if keyStart{
				return i,ErrInvalidateJson
			}
			keyStart = true
		case ':': //value
			if keyStart || objStart{
				return i,ErrInvalidateJson
			}
		case '[':
			if objStart || keyStart{
				return i,ErrInvalidateJson
			}
		case ']':
			if keyStart || keyStart{
				return i,ErrInvalidateJson
			}
		default:

		}
		i+=1
	}
	return btlen,ErrInvalidateJson
}

func NewRecord()*DxRecord  {
	result := new(DxRecord)
	result.PathSplitChar = '.'
	result.fValueType = DVT_Record
	result.fRecords = make(map[string]*DxBaseValue,32)
	return result
}
