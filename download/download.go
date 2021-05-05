package download

import (
	"bilibili_download/config"
	"bilibili_download/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var line = 1
var curLine = 1

type DownloadManager struct {
	info config.RoomInfo
	DLTotal int64
	DLUrl string
	downStatic int
	savefile *os.File
	reader *Reader
	flvHeader []byte
	flvMetaData []byte
	printLine int
}

func MakeNewManager(info config.RoomInfo) *DownloadManager {
	x:=&DownloadManager{
		info:    info,
		DLTotal: 0,
		DLUrl: "",
		downStatic:0,
		printLine: line,
	}
	line++
	return x
}

func (d *DownloadManager) Run() {
	//ticker先执行一次
	if d.info.CheckTime>10 {
		d.checkBegin()
	}
	ti:=time.NewTicker(time.Duration(d.info.CheckTime)*time.Second)
	for true {
		select {
		case <-ti.C:
			d.checkBegin()
		}
	}
}

func (d *DownloadManager)checkBegin()  {
	if d.downStatic!=1 {
		//开启时间检测模式
		if d.info.NeedTicker {
			now:=time.Now().Format("1504")
			now_time,_:=strconv.Atoi(now)
			if d.info.EndTime < d.info.BeginTime {
				d.info.EndTime+=24
			}
			if now_time>d.info.BeginTime&&now_time<d.info.EndTime {
				//在时间内.
				go func() {
					err:=d.begin()
					if err != nil {
						d.print("开启任务失败:",err.Error())
					}
				}()
			}else{
				d.print("不在自定义时间内")
			}
		}else{
			go func() {
				err:=d.begin()
				if err != nil {
					d.print("开启任务失败:",err.Error())
				}
			}()
		}
	}
}

func (d *DownloadManager) begin() error{
	req,err:=http.NewRequest("GET",fmt.Sprintf("https://live.bilibili.com/%s",d.info.RoomNum),nil)
	if err != nil {
		return err
	}
	for key, value := range header {
		req.Header.Set(key,value)
	}
	req.Header.Set("referer",fmt.Sprintf("https://live.bilibili.com/%s",d.info.RoomNum))
	req.Header.Set("cookie",config.Config().Cookie)
	client:=http.Client{}
	rep,err:=client.Do(req)
	if err != nil {
		return err
	}
	if rep.StatusCode!=200 {
		return errors.New(rep.Status)
	}
	defer rep.Body.Close()
	h5,err:=ioutil.ReadAll(rep.Body)
	if err != nil {
		return err
	}
	goqu,err:=goquery.NewDocumentFromReader(bytes.NewReader(h5))
	if err != nil {
		panic(err.Error())
	}
	goqu.Find("script").Each(func(i int, selection *goquery.Selection) {
		se:=selection.Text()
		if strings.HasPrefix(se,"window.__NEPTUNE_IS_MY_WAIFU__") {
			ddd:=strings.Replace(se,"window.__NEPTUNE_IS_MY_WAIFU__=","",-1)
			d.formatDownloadUrl(ddd)
		}
	})
	return nil
}

func (d *DownloadManager) formatDownloadUrl(text string)  {
	bb:=model.RoomInfo{}
	json.Unmarshal([]byte(text),&bb)
	if bb.RoomInitRes.Data.LiveStatus!=1 {
		d.print(d.info.RoomNum,"未开播")
		return
	}
	if len(bb.RoomInitRes.Data.PlayurlInfo.Playurl.Stream)>=1 {
		steam:=bb.RoomInitRes.Data.PlayurlInfo.Playurl.Stream[0]
		if len(steam.Format)>=1 {
			format:=steam.Format[0]
			if len(format.Codec)>=1 {
				codec:=format.Codec[0]
				if len(codec.UrlInfo)>=1 {
					urlinfo:=codec.UrlInfo[0]
					url:=fmt.Sprintf("%s%s%s",urlinfo.Host,codec.BaseUrl,urlinfo.Extra)
					//go d.changeFile()
					d.changeFile()
					go d.DownloadFileProgress(url)
				}
			}
		}
	}
}

type Reader struct {
	io.Reader
	Current int64
	RoomNum string
	ChangeFile int64
	PrintLine int
}

func (r *Reader) Read(p []byte) (n int, err error){
	n, err = r.Reader.Read(p)
	r.Current += int64(n)
	if curLine > r.PrintLine {
		fmt.Printf("\033[%dA\033[100D\033[K",r.PrintLine-curLine)
	}else if curLine<r.PrintLine{
		fmt.Printf("\033[%dB\033[100D\033[K",curLine-r.PrintLine)
	}
	curLine=r.PrintLine
	//curLine=r.PrintLine
	fmt.Printf("\r直播间%s:共计下载:%.4fMB---%d字节\033[K",r.RoomNum,float64(r.Current/1000)/1000.0,r.Current)
	if r.ChangeFile!=0&& r.Current/1024/1024> r.ChangeFile {
		//fmt.Println("退出")
		runtime.Goexit()
		return 0,errors.New("NEXT")
	}
	return
}

func (d *DownloadManager)DownloadFileProgress(url string) {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer func() {_ = r.Body.Close()}()
	//fip:=d.checkPath()
	//f, err := os.Create(path.Join(fip,time.Now().Format("20060102150405.flv")))
	//if err != nil {
	//	panic(err)
	//}
	defer func() {
		_ = d.savefile.Close()
	}()
	reader := &Reader{
		Reader: r.Body,
		RoomNum: d.info.RoomNum,
		ChangeFile: d.info.SaveSpace,
		PrintLine:d.printLine,
	}
	d.downStatic=1
	d.reader=reader
	defer func() {d.downStatic=0}()
	_, err = io.Copy(d.savefile, reader)
	if err!=nil&&err.Error()=="NEXT" {
		go d.begin()
	}
	if err != nil {
		d.print(fmt.Sprintf("下载失败:%s",err.Error()))
	}else{
		d.print("下载完毕/更换文件")
	}
	//d.saveToFile()
	//d.print("下载完毕/更换文件")
}

func (d *DownloadManager) checkPath() string {
	filePath:=path.Join(config.Config().Dir,d.info.RoomNum)
	if config.Config().Dir=="." {
		wd,_:=os.Getwd()
		filePath=path.Join(wd,d.info.RoomNum)
	}
	info,err:=os.Stat(filePath)
	if err != nil {
		os.MkdirAll(filePath,os.ModePerm)
	}else if !info.IsDir() {
		d.print("%s不是一个目录,退出")
		return ""
	}
	return filePath
}

func (d *DownloadManager) changeFile()  {
	//ti:=time.NewTicker(time.Minute)
	//for true {
		//select {
		//case <-ti.C:
		//	if d.savefile!=nil {
		//		//paths:=d.savefile.
		//		_,filename:=path.Split(d.savefile.Name())
		//		name:=strings.Replace(filename,".flv","",-1)
		//		fileTime,err:=time.Parse("20060102150405",name)
		//		if err != nil {
		//			fmt.Println("文件解析错误")
		//			continue
		//		}
		//		if fileTime.Add(time.Second*time.Duration(d.info.SaveTime)).Before(time.Now()) {
		//			continue
		//		}
		//		//continue
		//	}
			//if d.downStatic==0 {
			//	return
			//}
			//d.savefile.Close()
			var err error
			d.savefile,err=os.Create(path.Join(d.checkPath(),time.Now().Format("20060102150405.flv")))
			if err != nil {
				d.print("文件创建失败:%s",err.Error())
				return
			}
			//d.savefile.Write(d.flvHeader)
	//d.savefile.Write(d.flvMetaData)
	//d.print("更新文件")
			//time.Sleep(time.Second*time.Duration(d.info.SaveTime))
		//}
	//}
}

func (d *DownloadManager) saveToFile()  {
	size := 32 * 1024
	buf:=make([]byte,size)
	buff:=new(bytes.Buffer)
	//wg:=&sync.WaitGroup{}
	//wg.Add(1)
	//go func() {
	//
	//}()
	go func() {
		//wg.Wait()
	RESTART:
		var (
			meta [9]byte
			data []byte
		)
		_, err := buff.Read(meta[:])
		if  err != nil &&err!=io.EOF{
			panic(err)
			return
		}
		if err==io.EOF || meta[0]==0{
			goto RESTART
		}
		if string(meta[:4]) != "FLV\x01"  {
			//fmt.Println(n)
			panic(errors.New("!=flv"))
			return
		}
		l := bytes_int32(meta[5:])
		data = make([]byte, l+4)
		if _, err := buff.Read(data[9:]); err != nil {
			panic(err)
			return
		}
		copy(data, meta[:])
		d.flvHeader=data
		//d.savefile.Write(data)
		//wg.Done()
		d.print("头部读取结束")
		last:=make([]byte,11)
		for true {
			if d.downStatic==0 {
				return
			}
			var (
				meta []byte
				data []byte
			)
			type Tag struct {
				Head []byte
				Data []byte
			}
			meta = make([]byte, 11)
			if last!=nil&&last[0]!=0 {
				meta=last
				last=nil
				//copy()
			}else{
				_, err := buff.Read(meta)
				if  err != nil&&err!=io.EOF {
					return
				}
				if err==io.EOF {
					continue
				}
			}

			switch meta[0] & 31 {
			case 8, 9, 18:
				l := bytes_int24(meta[1:])
				data = make([]byte, l+4)
				if buff.Len()<l+4 {
					//fmt.Println("不足!继续等待")
					last=meta
					time.Sleep(time.Second)
					continue
				}
				if i, _ := buff.Read(data); i != int(l+4) {
					return
				}
				//fmt.Println(Tag{meta, data[:l]})
				//fmt.Printf("新的tag")
				//if meta[0]==0x9 {
				//	//meta[]
				//	//fmt.Println("%s",meta[0])
				//	cc:=meta[4:8]
				//	//fmt.Println(cc)
				//	Timestamp := uint32(cc[3])<<32 + uint32(cc[0])<<16 + uint32(cc[1])<<8 + uint32(cc[2])
				//	//fmt.Println(Timestamp)
				//	err:=d.f.WriteVideoTag(BytesCombine(meta,data),Timestamp)
				//	if err != nil {
				//		panic(err.Error())
				//	}
				//}
				//if meta[0]==0x8 {
				//	cc:=meta[4:8]
				//	Timestamp := uint32(cc[3])<<32 + uint32(cc[0])<<16 + uint32(cc[1])<<8 + uint32(cc[2])
				//	err:=d.f.WriteAudioTag(BytesCombine(meta,data),Timestamp)
				//	if err != nil {
				//		panic(err.Error())
				//	}
				//}
				if meta[0]==0x12 {
					//fmt.Println(18)
					d.flvMetaData=BytesCombine(meta,data)
					_,err:=d.savefile.Write(BytesCombine(d.flvHeader,d.flvMetaData))
					if err != nil {
						panic(err.Error())
					}
					d.savefile.Sync()
					//d.f.ChangeHeaderBytes(d.flvMetaData)
				}else{
					_,err:=d.savefile.Write(BytesCombine(meta,data))
					if err != nil {
						panic(err.Error())
					}
					d.savefile.Sync()
				}
			default:
				return
			}
		}
	}()
	for true {
		n,err:=d.reader.Read(buf)
		if err != nil && err!=io.EOF{
			d.print("reader 出错:",err.Error())
			break
		}
		//d.savefile.Write(buf[0:n])
		buff.Write(buf[0:n])
		if err==io.EOF {
			d.print("出现EOF")
			return
		}
		if err.Error()=="NEXT" {
			go d.begin()
			return
		}
	}
}

func (d *DownloadManager) print(s ...string) {
	if curLine > d.printLine {
		fmt.Printf("\033[%dA\033[K\033[%dD",d.printLine-curLine,100)
	}else if curLine<d.printLine{
		fmt.Printf("\033[%dB\033[K\033[%dD",curLine-d.printLine,100)
	}
	curLine=d.printLine
	fmt.Print(d.info.RoomNum,":",s)
}


func bytes_int32(x []byte) int {
	return int(x[0])<<24 + int(x[1])<<16 + int(x[2])<<8 + int(x[3])
}


func bytes_int24(x []byte) int {
	return int(x[0])<<16 + int(x[1])<<8 + int(x[2])
}

func BytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}
