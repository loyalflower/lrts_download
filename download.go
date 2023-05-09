package lrts_download

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"time"
)

type DownloadOptions struct {
	Token  string
	Output string // 输出目录
}

type listEntity struct {
	Buy           int         `json:"buy"`
	CanUnlock     int         `json:"canUnlock"`
	CantListen    int         `json:"cantListen"`
	DownPrice     float32     `json:"downPrice"`
	FatherTypeId  int         `json:"fatherTypeId"`
	FeeType       int         `json:"feeType"`
	HasAiLrc      int         `json:"hasAiLrc"`
	HasLyric      int         `json:"hasLyric"`
	Id            int64       `json:"id"`
	LastModify    int64       `json:"lastModify"`
	Length        int         `json:"length"`
	ListenPrice   float32     `json:"listenPrice"`
	Name          string      `json:"name"`
	OnlineTime    interface{} `json:"onlineTime"`
	Path          string      `json:"path"` // 收听地址
	PayType       int         `json:"payType"`
	Plays         int         `json:"plays"`
	Section       int         `json:"section"`
	SectionId     string      `json:"sectionId"`
	Size          int         `json:"size"`
	State         int         `json:"state"`
	Strategy      int         `json:"strategy"`
	TmeId         int         `json:"tmeId"`
	TypeId        int         `json:"typeId"`
	TypeName      string      `json:"typeName"`
	UnlockEndTime int         `json:"unlockEndTime"`
}

type listResponse struct {
	ApiStatus int          `json:"apiStatus"`
	BookId    int          `json:"bookId"`
	List      []listEntity `json:"list"`
	Msg       string       `json:"msg"`
	Sections  int          `json:"sections"`
	Status    int          `json:"status"`
	UserType  int          `json:"userType"`
}

type albumListResponse struct {
	ApiStatus int         `json:"apiStatus"`
	Count     int         `json:"count"`
	Detail    interface{} `json:"detail"`
	List      []struct {
		AudioId                int         `json:"audioId"`
		BaseEntityId           int         `json:"baseEntityId"`
		BaseEntityType         int         `json:"baseEntityType"`
		BestCover              string      `json:"bestCover"`
		Buy                    int         `json:"buy"`
		CanUnlock              int         `json:"canUnlock"`
		Cover                  string      `json:"cover"`
		HasAiLrc               int         `json:"hasAiLrc"`
		HasLyric               int         `json:"hasLyric"`
		Length                 int         `json:"length"`
		Name                   string      `json:"name"`
		OnlineTime             int64       `json:"onlineTime"`
		Path                   string      `json:"path"`
		PayType                int         `json:"payType"`
		PlayCount              int         `json:"playCount"`
		Reason                 string      `json:"reason"`
		RefId                  int         `json:"refId"`
		Section                int         `json:"section"`
		SectionId              string      `json:"sectionId"`
		Size                   int         `json:"size"`
		SrcCommentControlType  int         `json:"srcCommentControlType"`
		SrcEntityCommentCount  int         `json:"srcEntityCommentCount"`
		SrcEntityCover         string      `json:"srcEntityCover"`
		SrcEntityDesc          string      `json:"srcEntityDesc"`
		SrcEntityId            int         `json:"srcEntityId"`
		SrcEntityName          string      `json:"srcEntityName"`
		SrcEntityPlayCount     interface{} `json:"srcEntityPlayCount"`
		SrcEntityRankingInfo   string      `json:"srcEntityRankingInfo"`
		SrcEntityRankingList   interface{} `json:"srcEntityRankingList"`
		SrcEntityRankingTarget string      `json:"srcEntityRankingTarget"`
		SrcEntityTtsRefId      interface{} `json:"srcEntityTtsRefId"`
		SrcEntityTtsType       interface{} `json:"srcEntityTtsType"`
		SrcEntityUserCover     string      `json:"srcEntityUserCover"`
		SrcEntityUserFollow    interface{} `json:"srcEntityUserFollow"`
		SrcEntityUserId        interface{} `json:"srcEntityUserId"`
		SrcEntityUserName      string      `json:"srcEntityUserName"`
		SrcId                  int         `json:"srcId"`
		SrcRewarded            int         `json:"srcRewarded"`
		SrcSection             int         `json:"srcSection"`
		SrcSectionId           string      `json:"srcSectionId"`
		SrcType                int         `json:"srcType"`
		Strategy               int         `json:"strategy"`
		TmeId                  int         `json:"tmeId"`
		TypeId                 int         `json:"typeId"`
		TypeName               string      `json:"typeName"`
		UnlockEndTime          int         `json:"unlockEndTime"`
	} `json:"list"`
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

type albumEntityListResponse struct {
	ApiStatus int `json:"apiStatus"`
	List      []struct {
		Attach    string `json:"attach"`
		Id        int    `json:"id"`
		Md5Code   int    `json:"md5Code"`
		Path      string `json:"path"`
		PathMeta  string `json:"pathMeta"`
		Section   int    `json:"section"`
		SectionId string `json:"sectionId"`
		Type      int    `json:"type"`
	} `json:"list"`
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

const (
	secret = "iJ0DgxmdC83#I&j@iwg"

	domainDownload         = "http://dapis.mting.info"
	urlPathBookList        = "/yyting/bookclient/ClientGetBookResource.action"
	urlPathAlbumList       = "/yyting/snsresource/getAblumnAudios.action"
	urlPathAlbumEntityList = "/yyting/gateway/entityPath.action"

	ctxBookIdFieldName    = "bookId"
	ctxAudioNameFieldName = "audioName"
	ctxOrderSortFieldName = "orderSort"
	// 只是临时存储用的，还不涉及下载
	ctxMidBookIdFieldName    = "midBookId"
	ctxMidAudioNameFieldName = "midAudioName"

	requestTimeout = 30 * time.Second
	randomDelay    = 1 * time.Second

	defaultToken = "OqzlvCxt2i_P1SZKF6GjFg**_lK0uCQpm5tN-P6XdFZYawCDKSgeC4anU"
	defaultIMEI  = "MDI6MDA6MDA6MDA6MDA6DDA="
)

var (
	downloadHttpHeaders = http.Header{
		"User-Agent":    []string{"Android13/yyting/vivo/V2133A/ch_uc_beta/167/Android"},
		"Referer":       []string{"yytingting.com"},
		"ClientVersion": []string{"6.3.4.0"},
	}
	downloadListHttpQuery = url.Values{
		"bookId":   []string{""},
		"pageNum":  []string{"1"},
		"pageSize": []string{"50"},
		"sortType": []string{"0"},
		"token":    []string{defaultToken},
		"imei":     []string{defaultIMEI},
		"nwt":      []string{"1"},
		"q":        []string{"1930"},
	}

	downloadAlbumListHttpQuery = url.Values{
		"ablumnId": []string{""},
		"pageNum":  []string{"1"},
		"pageSize": []string{"10000"},
		"sortType": []string{"0"},
		"token":    []string{defaultToken},
		"imei":     []string{defaultIMEI},
		"nwt":      []string{"1"},
		"q":        []string{"2205"},
	}

	downloadAlbumEntityHttpQuery = url.Values{
		"entityId":   []string{""},
		"entityType": []string{"2"},
		"opType":     []string{"0"},
		"sections":   []string{"[]"}, // sections id.  [7463598,7463599]
		"type":       []string{"0"},
		"token":      []string{defaultToken},
		"imei":       []string{defaultIMEI},
		"nwt":        []string{"1"},
		"q":          []string{"2224"},
	}
)

func Download(bookId string, options DownloadOptions) error {
	c := colly.NewCollector(colly.AllowURLRevisit())
	c.SetRequestTimeout(requestTimeout)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: randomDelay,
	})

	// create a request queue with 2 consumer threads
	q, _ := queue.New(
		10, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000}, // Use default queue storage
	)

	c.OnRequest(func(r *colly.Request) {
		if name := r.Ctx.Get(ctxAudioNameFieldName); name != "" {
			fmt.Printf("Downloading %s\n", name)
		}
		if name := r.Ctx.Get(ctxMidAudioNameFieldName); name != "" {
			fmt.Printf("Search Resource %s\n", name)
		}
	})
	c.OnResponse(func(resp *colly.Response) {
		r := downloadOnResponse(resp, options)
		if r != nil {
			for _, v := range r {
				q.AddRequest(v)
			}
		}
	})
	c.OnError(func(response *colly.Response, err error) {
		fmt.Println("Request URL:", response.Request.URL, "failed with response:", string(response.Body), "\nError:", err)
	})

	q.AddRequest(downloadInitRequest(domainDownload+urlPathBookList, downloadInitQuery(bookId, downloadListHttpQuery)))
	return q.Run(c)
}

func downloadInitQuery(bookId string, urlValues url.Values) url.Values {
	httpQuery := make(map[string][]string, len(urlValues))
	for k, v := range urlValues {
		httpQuery[k] = v
	}
	if bookId != "" {
		httpQuery["bookId"] = []string{bookId}
	}
	return httpQuery
}

func downloadInitRequest(uri string, values url.Values) *colly.Request {
	u, _ := url.Parse(uri)
	if values != nil {
		values.Set("sc", paramsMd5(values, u.Path, secret))
		u.RawQuery = values.Encode()
	}
	r := &colly.Request{
		URL:     u,
		Method:  "GET",
		Ctx:     colly.NewContext(),
		Headers: &downloadHttpHeaders,
	}
	return r
}

func downloadOnResponse(resp *colly.Response, options DownloadOptions) []*colly.Request {
	if resp.StatusCode != 200 {
		return nil
	}
	if resp.Request.URL.Path == urlPathBookList {
		list := listResponse{}
		err := json.Unmarshal(resp.Body, &list)
		if err != nil {
			return nil
		}
		if len(list.List) == 0 {
			// 如果是第一页就没有数据
			// 尝试更换其他下载方式
			if resp.Request.URL.Query().Get("pageNum") == "1" {
				q := downloadInitQuery("", downloadAlbumListHttpQuery)
				q.Set("ablumnId", resp.Request.URL.Query().Get("bookId"))
				r := downloadInitRequest(
					domainDownload+urlPathAlbumList,
					q,
				)
				return []*colly.Request{r}
			}
			return nil
		}
		r := make([]*colly.Request, len(list.List))
		for i, entity := range list.List {
			r[i] = downloadInitRequest(entity.Path, nil)
			r[i].Ctx.Put(ctxAudioNameFieldName, entity.Name)
			r[i].Ctx.Put(ctxBookIdFieldName, cast.ToString(list.BookId))
			r[i].Ctx.Put(ctxOrderSortFieldName, cast.ToString(entity.Section))
		}
		// 继续拉取列表
		if len(list.List) == cast.ToInt(resp.Request.URL.Query().Get("pageSize")) {
			values := resp.Request.URL.Query()
			values.Set("pageNum", cast.ToString(cast.ToInt(values.Get("pageNum"))+1))
			r = append(r, downloadInitRequest(domainDownload+urlPathBookList, values))
		}
		return r
	} else if resp.Request.URL.Path == urlPathAlbumList {
		albumList := albumListResponse{}
		err := json.Unmarshal(resp.Body, &albumList)
		if err != nil {
			return nil
		}
		r := make([]*colly.Request, len(albumList.List))
		for i, entity := range albumList.List {
			q := downloadInitQuery("", downloadAlbumEntityHttpQuery)
			q.Set("entityId", resp.Request.URL.Query().Get("ablumnId"))
			q.Set("sections", "["+entity.SectionId+"]")
			r[i] = downloadInitRequest(domainDownload+urlPathAlbumEntityList, q)
			r[i].Ctx.Put(ctxMidAudioNameFieldName, entity.Name)
			r[i].Ctx.Put(ctxMidBookIdFieldName, q.Get("entityId"))
			r[i].Ctx.Put(ctxOrderSortFieldName, cast.ToString(entity.Section))
		}
		return r
	} else if resp.Request.URL.Path == urlPathAlbumEntityList {
		albumEntityList := albumEntityListResponse{}
		err := json.Unmarshal(resp.Body, &albumEntityList)
		if err != nil {
			return nil
		}
		r := make([]*colly.Request, len(albumEntityList.List))
		for i, entity := range albumEntityList.List {
			r[i] = downloadInitRequest(entity.Path, nil)
			r[i].Ctx.Put(ctxAudioNameFieldName, resp.Ctx.Get(ctxMidAudioNameFieldName))
			r[i].Ctx.Put(ctxBookIdFieldName, resp.Ctx.Get(ctxMidBookIdFieldName))
		}
		return r
	} else {
		saveAudio(resp, options)
	}

	return nil
}

func saveAudio(resp *colly.Response, options DownloadOptions) error {
	name := resp.Request.Ctx.Get(ctxAudioNameFieldName)
	bookId := resp.Request.Ctx.Get(ctxBookIdFieldName)
	if name == "" || bookId == "" {
		return errors.New("saveAudio failed, name or bookId is empty")
	}
	name = resp.Request.Ctx.Get(ctxOrderSortFieldName) + "_" + name + path.Ext(resp.Request.URL.Path)
	saveDir, dirErr := os.Getwd()
	if dirErr != nil {
		return dirErr
	}
	dirName := bookId
	if options.Output != "" {
		dirName = options.Output
	}
	saveDir = filepath.Join(saveDir, dirName)
	if _, existsErr := os.Stat(saveDir); os.IsNotExist(existsErr) {
		mkErr := os.MkdirAll(saveDir, os.ModePerm)
		if mkErr != nil {
			return mkErr
		}
	}
	savePath := filepath.Join(saveDir, name)
	return os.WriteFile(savePath, resp.Body, os.ModePerm)
}

func paramsMd5(query url.Values, urlPath, secret string) string {
	keys := lo.Keys(query)
	sort.Strings(keys)
	queryString := ""
	for i := range keys {
		if keys[i] == "sc" {
			continue
		}
		queryString += keys[i] + "=" + query.Get(keys[i]) + "&"
	}
	if queryString != "" {
		queryString = queryString[:len(queryString)-1]
	}
	// 将字符串转换为字节切片
	data := []byte(urlPath + "?" + queryString + secret)

	// 计算字符串的 md5 值
	hash := md5.Sum(data)

	// 将字节数组转换为十六进制字符串
	return hex.EncodeToString(hash[:])
}
