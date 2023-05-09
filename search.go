package lrts_download

import (
	"encoding/json"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gocolly/colly"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cast"
	"net/url"
	"os"
	"strings"
)

type SearchOptions struct {
	CommonOptions
	Type string // 搜索类型，对应 searchOption 参数
}

const (
	urlPathBookSearch = "/yyting/search/searchBatch.action"

	searchTypeBook   = "书籍"
	searchTypeFolder = "听单"
	searchTypeAlbum  = "节目"
	searchTypeRead   = "看书"
	searchTypeAnchor = "主播"
)

var (
	downloadSearchHttpQuery = url.Values{
		"keyWord":      []string{""},
		"searchOption": []string{"1,2,3,4,5"}, // 书籍 节目 看书 主播 听单
		"pageSize":     []string{"5"},
		"type":         []string{"0"},
		"token":        []string{defaultToken},
		"imei":         []string{defaultIMEI},
		"nwt":          []string{"1"},
		"q":            []string{"2150"},
	}
)

// 听单
type searchFolderResult struct {
	Count   int `json:"count"`
	HasNext int `json:"hasNext"`
	List    []struct {
		CollectCount int         `json:"collectCount"`
		Cover        string      `json:"cover"`
		EntityCount  int         `json:"entityCount"`
		EntityList   interface{} `json:"entityList"`
		Id           int64       `json:"id"`
		Name         string      `json:"name"`
		NickName     string      `json:"nickName"`
		OverallPos   int         `json:"overallPos"`
		Pt           int         `json:"pt"`
		SourceType   int         `json:"sourceType"`
		Url          string      `json:"url"`
		UserId       int         `json:"userId"`
	} `json:"list"`
}

// 节目
type searchAlbumResult struct {
	Count   int `json:"count"`
	HasNext int `json:"hasNext"`
	List    []struct {
		BaseEntityId   int         `json:"baseEntityId"`
		BaseEntityType int         `json:"baseEntityType"`
		CommentCount   int         `json:"commentCount"`
		Cover          string      `json:"cover"`
		CreateTime     int64       `json:"createTime"`
		Description    string      `json:"description"`
		EntityType     int         `json:"entityType"`
		Flag           int         `json:"flag"`
		H5Url          string      `json:"h5Url"`
		Id             int         `json:"id"`
		IsH5Book       interface{} `json:"isH5Book"`
		Name           string      `json:"name"`
		NickName       string      `json:"nickName"`
		PayFree        int         `json:"payFree"`
		PayType        int         `json:"payType"`
		PlayCount      int         `json:"playCount"`
		RankingInfo    string      `json:"rankingInfo"`
		RankingTarget  string      `json:"rankingTarget"`
		RecReason      string      `json:"recReason"`
		RollAdUnlock   int         `json:"rollAdUnlock"`
		Sections       int         `json:"sections"`
		ShortRecReason string      `json:"shortRecReason"`
		Source         int         `json:"source"`
		SourceType     int         `json:"sourceType"`
		Tags           []struct {
			BgColor string `json:"bgColor"`
			Name    string `json:"name"`
			Type    int    `json:"type"`
		} `json:"tags"`
		TypeName   string `json:"typeName"`
		UpdateTime int64  `json:"updateTime"`
		UserId     int    `json:"userId"`
		UserState  int    `json:"userState"`
	} `json:"list"`
}

// 看书
type searchReadResult struct {
	Count   int `json:"count"`
	HasNext int `json:"hasNext"`
	List    []struct {
		Author       string        `json:"author"`
		ContentState int           `json:"contentState"`
		Cover        string        `json:"cover"`
		Desc         string        `json:"desc"`
		Id           int           `json:"id"`
		Name         string        `json:"name"`
		OverallPos   int           `json:"overallPos"`
		PayType      int           `json:"payType"`
		Pt           int           `json:"pt"`
		RecReason    string        `json:"recReason"`
		SourceType   int           `json:"sourceType"`
		Tags         []interface{} `json:"tags"`
		Type         string        `json:"type"`
		Url          string        `json:"url"`
	} `json:"list"`
}

// 书籍
type searchBookResult struct {
	Count   int `json:"count"`
	HasNext int `json:"hasNext"`
	List    []struct {
		Announcer       string      `json:"announcer"`
		Author          string      `json:"author"`
		BaseEntityId    int         `json:"baseEntityId"`
		BaseEntityType  int         `json:"baseEntityType"`
		CommentCount    int         `json:"commentCount"`
		CommentMean     string      `json:"commentMean"`
		Cover           string      `json:"cover"`
		Desc            string      `json:"desc"`
		DisplayOrder    int         `json:"displayOrder"`
		DownPrice       int         `json:"downPrice"`
		EntityId        int         `json:"entityId"`
		EntityType      int         `json:"entityType"`
		FeeTypeId       interface{} `json:"feeTypeId"`
		H5Url           string      `json:"h5Url"`
		Hot             int         `json:"hot"`
		Id              int         `json:"id"`
		IsH5Book        interface{} `json:"isH5Book"`
		LastUpdateTime  string      `json:"lastUpdateTime"`
		ListenPrice     int         `json:"listenPrice"`
		Name            string      `json:"name"`
		PayFree         int         `json:"payFree"`
		PayType         int         `json:"payType"`
		Plays           int         `json:"plays"`
		Price           int         `json:"price"`
		RankingInfo     string      `json:"rankingInfo"`
		RankingTarget   string      `json:"rankingTarget"`
		RecReason       string      `json:"recReason"`
		Sections        int         `json:"sections"`
		ShortRecReason  string      `json:"shortRecReason"`
		ShowFreeEndTime int         `json:"showFreeEndTime"`
		Sort            int         `json:"sort"`
		SourceType      int         `json:"sourceType"`
		State           int         `json:"state"`
		Strategy        int         `json:"strategy"`
		Tags            []struct {
			BgColor string `json:"bgColor"`
			Name    string `json:"name"`
			Type    int    `json:"type"`
		} `json:"tags"`
		TypeId     int    `json:"typeId"`
		TypeName   string `json:"typeName"`
		VipLibrary int    `json:"vipLibrary"`
	} `json:"list"`
}

type searchResponse struct {
	commonApiResponse
	Data struct {
		FolderResult searchFolderResult `json:"folderResult"`
		AlbumResult  searchAlbumResult  `json:"albumResult"`
		ReadResult   searchReadResult   `json:"readResult"`
		Point        string             `json:"point"`
		BookResult   searchBookResult   `json:"bookResult"`
	} `json:"data"`
}

func Search(keywords string, options SearchOptions) error {
	query := downloadInitQuery("", downloadSearchHttpQuery)
	query.Set("keyWord", keywords)

	c := colly.NewCollector()

	// Find and visit all links
	c.OnResponse(func(resp *colly.Response) {
		searchOnResponse(resp)
	})

	c.OnError(func(response *colly.Response, err error) {
		fmt.Printf("Request URL: %s \nfailed with response: %s\nError: %s\n", response.Request.URL, string(response.Body), err)
	})

	query.Set("sc", paramsMd5(query, urlPathBookSearch, secret))

	urlEntity, parseErr := url.Parse(domainDownload + urlPathBookSearch)
	if parseErr != nil {
		return parseErr
	}
	urlEntity.RawQuery = query.Encode()

	return c.Request("GET", urlEntity.String(), nil, nil, downloadHttpHeaders)
}

func searchOnResponse(resp *colly.Response) error {
	if resp.StatusCode != 200 {
		return errors.New("search response error")
	}
	result := searchResponse{}
	err := json.Unmarshal(resp.Body, &result)
	if err != nil {
		return err
	}
	tableData := [][]string{}
	tableHeader := []string{"分类", "书名", "book_id", "作者", "集数", "简介"}
	textLimit := func(text string) string {
		return string([]rune(strings.TrimSpace(text))[:25])
	}
	if len(result.Data.BookResult.List) > 0 {
		for _, v := range result.Data.BookResult.List {
			tableData = append(tableData, []string{searchTypeBook, v.Name, cast.ToString(v.Id), v.Announcer, cast.ToString(v.Sections), textLimit(v.RecReason)})
		}
		tableData = append(tableData, []string{"全部数量", cast.ToString(result.Data.BookResult.Count), "", "", "", ""})
		tableData = append(tableData, []string{""})
	}
	if len(result.Data.AlbumResult.List) > 0 {
		for _, v := range result.Data.AlbumResult.List {
			tableData = append(tableData, []string{searchTypeAlbum, v.Name, cast.ToString(v.Id), v.NickName, cast.ToString(v.Sections), textLimit(v.Description)})
		}
		tableData = append(tableData, []string{"全部数量", cast.ToString(result.Data.AlbumResult.Count), "", "", "", ""})
		tableData = append(tableData, []string{""})
	}
	if len(result.Data.ReadResult.List) > 0 {
		for _, v := range result.Data.ReadResult.List {
			tableData = append(tableData, []string{searchTypeRead, v.Name, cast.ToString(v.Id), v.Author, "", textLimit(v.RecReason)})
		}
		tableData = append(tableData, []string{"全部结果", cast.ToString(result.Data.ReadResult.Count), "", "", "", ""})
		tableData = append(tableData, []string{""})
	}
	if len(result.Data.FolderResult.List) > 0 {
		for _, v := range result.Data.FolderResult.List {
			tableData = append(tableData, []string{searchTypeFolder, v.Name, cast.ToString(v.Id), v.NickName, cast.ToString(v.EntityCount), ""})
		}
		tableData = append(tableData, []string{"全部结果", cast.ToString(result.Data.FolderResult.Count), "", "", "", ""})
		tableData = append(tableData, []string{""})
	}

	outputTable(tableHeader, tableData)

	return nil
}

func outputTable(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader(header)
	table.AppendBulk(data) // Add Bulk Data
	table.SetRowSeparator("-")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}
