package DxValue

import (
	"testing"
	"io/ioutil"
	"fmt"
	"github.com/suiyunonghen/DxCommonLib"
	"unsafe"
	"time"
	"bytes"
)


type mTest struct{
	bB	string
	bMstr bool
}

type ATest struct {
	A  int
	B  int
	mc map[string]int
}

func TestDxRecord_Escape(t *testing.T) {
	/*str := `{"ctrlpath":"C:\\frank\\pvt_new\\te\"mp\\\\\"lateopt\""}`
	str = `["asdf","C:\\frank\\pvt_new\\te\"mp\\\\\"lateopt\"","asdf","\"\"\\\\"]`
	vc := NewArray()// NewRecord()
	vc.JsonParserFromByte(([]byte)(str),true,true)
	fmt.Println(vc.ToString())*/
	//str := "{\"id\":\"00\",\"output\":\"ntripsvr://0000@58.49.94.210:2103/WUH9\",\"inputidx\":1}"

	lvReq := NewRecord()
	lvReq.SetString("msg", "HTTP/1.0 401 Unauthorized \r\n WWW-Authenticate: Basic realm=\"RTCM32_GGB\"")
	str:=string(lvReq.BytesWithSort(false))
	fmt.Println(str)
	fmt.Println(lvReq.AsString("msg",""))

	str = `{"msg":"\r\n123"}`

	json := NewRecord()
	json.ClearValue(true)
	json.JsonParserFromByte([]byte(str), true, false)
	arr := json.NewArray("testArray",false)


	json.JsonParserFromByte([]byte(`{"testArray":["/Date(1558844627000)/"],"msg":"\r\n123"}`),true,true)

	fmt.Println(json.AsBytes("msg"))
	arr = json.AsArray("testArray")
	fmt.Println(arr.AsDateTime(0,0).ToTime())
	fmt.Println(json.AsString("msg",""))
	fmt.Println(json.ToString())
}

func TestDxRecord_ForcePathRecord(t *testing.T) {
	json := NewRecord()
	s := json.ForcePathArray("data")
	s.Append("hello")
	s = json.ForcePathArray("data")
	s.Append("word")
	fmt.Println(json.ToString())
	return
	rc := json.ForcePathRecord("test.gg.mm")
	rc.SetString("Name","不得an")
	fmt.Println(json.ToString())
	rc = json.ForcePathRecord("test.gg")
	fmt.Println(rc.ToString())

	arr := rc.ForcePathArray("mm.Array")
	arr.SetString(0,"ASdf")
	arr.SetString(1,"ceshia")
	fmt.Println(rc.ToString())
}

func TestDxRecord_BytesWithSort(t *testing.T) {
	lvResp := NewRecord()
	lvAcc := lvResp.ForcePathRecord("account")
	lvAcc.SetInt("max", 34)
	lvAcc.SetInt("used", 42)
	fmt.Println(string(lvResp.Bytes(false)))
	str :=string( lvResp.BytesWithSort(false))
	fmt.Println(str)

}


func Test_Record(t *testing.T)  {
	r := NewRecord()
	r.JsonParserFromByte([]byte(`{
"sys":{
  "sd":"c:\"\\t\"est"},"gg":"asdfasdf","value":23,"age":232.24,
"Node1":{
	\"234234\":"as\td\\\"faf"
}`),true,false)
	fmt.Println(r.String())

	mA := &ATest{A:123,B:234}
	mA.mc = make(map[string]int)
	mA.mc["saf"]=23443
	fmt.Println(mA)
	fmt.Println(uintptr(unsafe.Pointer(mA)) )
	mB := &ATest{A:3423,B:23434}
	*mA = *mB
	fmt.Println(uintptr(unsafe.Pointer(mA)) )
	fmt.Println(mA)
}

func TestDxRecord_SetRecordValue(t *testing.T) {
	vc := NewRecord()
	vcc := vc.NewRecord("testc",true)
	fmt.Println(vcc.AsStringByPath("testc.gg.asdf",""))
	vcc.SetString("BB","Asdf")
	fmt.Println(vc.String())
	mb := NewRecord()
	mb.SetInt("gg",123)
	vc.SetRecordValue("testc",mb)
	fmt.Println(vc.String())
}

func TestDxRecord_Extract(t *testing.T) {
	vc := NewRecord()
	vcc := vc.NewRecord("test",true)
	vcc.SetString("Name","不得闲")
	vcc.SetString("Sex","男")
	vcc.SetInt("Age",32)
	vc.SetGoTime("now",time.Now())
	fmt.Println(vc.String())

	vcc.Extract()
	fmt.Println(vcc.String())
	fmt.Println(vc.String())
}

func TestDxRecord_Masharl(t *testing.T)  {
	vc := NewRecord()
	arr := NewArray()
	fmt.Println(arr)
	vcc := vc.NewRecord("testc",true)
	vcc.SetString("BB","Asdf")
	bt,err := Marshal(vc)
	if err == nil{
		mr := NewRecord()
		Unmarshal(bt,mr)
		fmt.Println(mr.ToString())
	}
}

func TestDxRecord_JsonParserFromByte(t *testing.T) {
	buf, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	rc := NewRecord()
	_,err = rc.JsonParserFromByte(buf,true,false)
	if err != nil{
		fmt.Println("Parser Error: ",err)
	}
	fmt.Println(rc.ToString())
	rc.SaveMsgPackFile("DataProxy.config.msgPack")
}

func TestParserTime(t *testing.T)  {
	fmt.Println(time.Now().Format("2006-01-02T15:04:05Z"))
	fmt.Println(time.Parse("2006-01-02T15:04:05Z","2010-07-12T03:05:21Z"))
	at := DxCommonLib.ParserJsonTime("/Date(1402384458000)/")
	fmt.Println(at.ToTime().Format("2006-01-02T15:04:05Z"))
	at = DxCommonLib.ParserJsonTime("/Date(1224043200000+0800)/")
	fmt.Println(at.ToTime().Format("2006-01-02T15:04:05Z"))
}


func TestDxRecord_AsBool(t *testing.T) {
	rc := NewRecord()
	rc.JsonParserFromByte([]byte(`{"BoolValue":  true  ,"object":{"objBool":  false  }}`),false,false)
	fmt.Println("BoolValue=",rc.AsBool("BoolValue",false))
	fmt.Println("object.objBool=",rc.AsBoolByPath("object.objBool",true))
}

func TestDxRecord_AsArray(t *testing.T) {
	buf, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	rc := NewRecord()
	_,err = rc.JsonParserFromByte(buf,true,false)
	if err != nil{
		fmt.Println("Parser Error: ",err)
	}
	array := rc.AsArray("list")
	if array != nil{

		for i := 0;i<array.Length();i++{
			fmt.Print("Array Index ",i,"=")
			switch array.VaueTypeByIndex(0) {
			case DVT_Record:
				rc := array.AsRecord(i)
				if i == 1{
					fmt.Println("input.remark=",rc.AsStringByPath("input.remark","unkonwn"))
				}
				fmt.Println(rc.ToString())
			}
		}
	}
}


func TestEscapStr(t *testing.T){
	stb := []byte(`{"id":"001", "data":"$GPGGA,093805.00,2255.48843,N,\"11401.10693,E,1,23,0.6,10.527,M,0.000,M,0.0,0001*49"}`)
	rec := NewRecord()
	rec.JsonParserFromByte(stb,true,false)
	fmt.Println(rec.String())

 	fmt.Println(DxCommonLib.EscapeJsonStr(`测试"fasdf""`))

	fmt.Println(DxCommonLib.ParserEscapeStr([]byte(`\u6D4B\u8BD5\"fasdf\"\"`)))

}

func TestDxValue_JsonParserFromByte(t *testing.T) {
	var v DxValue
	buf, err := ioutil.ReadFile("DataProxy.config.json")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	index,err := v.JsonParserFromByte(buf,true,false)
	if err != nil{
		fmt.Println("Parser Error: ",err," index=",index)
	}else{
		switch v.ValueType() {
		case DVT_Record:
			fmt.Println("Is Json Object: ",(*DxRecord)(unsafe.Pointer(v.fValue)).ToString())
		case DVT_Array:
			fmt.Println("Is Json Array: ",(*DxArray)(unsafe.Pointer(v.fValue)).ToString())
		}
	}
}


func TestDxRecord_SaveJsonFile(t *testing.T) {
	rec := NewRecord()
	rec.SetInt("Age",-12)
	rec.SetString("Name","suiyunonghen")
	rec.SetValue("Home",map[string]interface{}{
		"Addres": "湖北武汉",
		"Code":"430000",
		"Peoples":4,
	})
	rec.SetDouble("Double",234234234.4564564)
	rec.SetFloat("Float",-34.534)
	rec.SetValue("Now",time.Now())
	//rec.SaveJsonFile("d:\\testJson.json",true)
	//rec.SaveMsgPackFile("d:\\msgpack.bin")
	fmt.Println(rec.ToString())
}

func TestMsgPackDecode(t *testing.T)  {
	bt, err := ioutil.ReadFile("d:\\1.bin")
	if err != nil {
		fmt.Println("ReadFile Err:",err)
		return
	}
	coder := NewDecoder(bytes.NewReader(bt))
	rec := NewRecord()
	if err := coder.Decode(&rec.DxBaseValue);err!=nil{
		fmt.Println("Error；",err)
	}
	fmt.Println(rec.AsDateTime("create_time",0).ToTime())
	fmt.Println(rec.ToString())
}

func TestDxRecord_LoadMsgPackFile(t *testing.T) {
	rec := NewRecord()
	if err := rec.LoadMsgPackFile("test.Msgpack");err!=nil{
		fmt.Println("Error；",err)
	}else{
		fmt.Println(DxCommonLib.FastByte2String(rec.BytesWithSort(false)))
	}
}


func TestDxRecord_AsString(t *testing.T) {
	rc := NewRecord()
	rc.JsonParserFromByte([]byte(`{"StringValue":   "String    3"  ,  "object":  {"objStr":"  ObjStr1  ",  "ObjName"  :  "   I  nnerObje  ct  "}}`),false,false)
	fmt.Println("StringValue=",rc.AsString("StringValue",""))
	fmt.Println("object.objStr=",rc.AsStringByPath("object.objStr",""))
	fmt.Println("object.ObjName=",rc.AsStringByPath("object.ObjName",""))
}

func TestDxRecord_SetIntRecordValue(t *testing.T) {
	rc := NewRecord()
	inarc := NewIntKeyRecord()
	inarc.SetInt(2,23)
	inarc.SetValue(23,"DxSoft")
	rc.SetIntRecordValue("IntRecord",inarc)
	fmt.Println(rc.Contains("IntRecord.23"))
	if !rc.Contains("IntRecord.ConVertName"){
		rc.ForcePath("IntRecord.ConVertName","Record")
		fmt.Println(rc.Contains("IntRecord.ConVertName"))
	}
	fmt.Println(rc.ToString())
}

func TestDxIniDecoder_Decode(t *testing.T) {
	bt, _ := ioutil.ReadFile("E:\\Delphi\\Leigod\\BoheBin\\BoHe.ini")
	if bt[0] == 0xEF && bt[1] == 0xBB && bt[2] == 0xBF { //UTF-8
		bt = bt[3:]
	}
	buffer := bytes.NewBuffer(bt)
	decoder := NewIniDecoder(buffer,DxCommonLib.File_Code_Utf8)
	r := NewRecord()
	decoder.Decode(r)
	fmt.Println(r.ToString())
}