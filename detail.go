package lrts_download

import (
	"encoding/json"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gocolly/colly"
	"github.com/spf13/cast"
	"net/url"
)

type DetailOptions struct {
}

const (
	urlPathBookDetail = "/yyting/page/bookDetailPage.action"
)

var (
	downloadDetailHttpQuery = url.Values{
		"bookId": []string{""},
		"token":  []string{"OqzlvCxt2i_P1SZKF6GjFg**_lK0uCQpm5tN-P6XdFZYawCDKSgeC4anU"},
		"imei":   []string{"MDI6MDA6MDA6MDA6MDA6MDA="},
		"nwt":    []string{"1"},
		"q":      []string{"2150"},
	}
)

type detailBookDetailResult struct {
	AdvertControlType  int    `json:"advertControlType"`
	Announcer          string `json:"announcer"`
	ApiStatus          int    `json:"apiStatus"`
	Author             string `json:"author"`
	BaseEntityId       int    `json:"baseEntityId"`
	BaseEntityType     int    `json:"baseEntityType"`
	BestCover          string `json:"bestCover"`
	CantDown           int    `json:"cantDown"`
	CantListen         int    `json:"cantListen"`
	CommentControlType int    `json:"commentControlType"`
	CommentCount       int    `json:"commentCount"`
	CommentMean        string `json:"commentMean"`
	Cover              string `json:"cover"`
	Desc               string `json:"desc"`
	DownPrice          int    `json:"downPrice"`
	Download           int    `json:"download"`
	EstimatedSections  int    `json:"estimatedSections"`
	ExtInfo            string `json:"extInfo"`
	ExtraInfos         []struct {
		Content string `json:"content"`
		Title   string `json:"title"`
	} `json:"extraInfos"`
	FatherTypeId   int    `json:"fatherTypeId"`
	FatherTypeName string `json:"fatherTypeName"`
	FeeType        int    `json:"feeType"`
	FreeEndTime    int    `json:"freeEndTime"`
	Id             int    `json:"id"`
	IsCollected    int    `json:"isCollected"`
	IsLike         int    `json:"isLike"`
	IsSend         int    `json:"isSend"`
	Labels         []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	} `json:"labels"`
	Length      int    `json:"length"`
	ListenPrice int    `json:"listenPrice"`
	Msg         string `json:"msg"`
	Name        string `json:"name"`
	OrgId       int    `json:"orgId"`
	PayType     int    `json:"payType"`
	Play        int    `json:"play"`
	PriceInfo   struct {
		Activitys            []interface{} `json:"activitys"`
		ApiStatus            int           `json:"apiStatus"`
		Buys                 string        `json:"buys"`
		CanUnlock            int           `json:"canUnlock"`
		CanUseTicket         int           `json:"canUseTicket"`
		ChargeGiftLabel      string        `json:"chargeGiftLabel"`
		ChoosePrice          int           `json:"choosePrice"`
		DeadlineTime         int           `json:"deadlineTime"`
		DiscountPrice        int           `json:"discountPrice"`
		Discounts            []interface{} `json:"discounts"`
		EntityId             int           `json:"entityId"`
		EntityType           int           `json:"entityType"`
		EstimatedSections    int           `json:"estimatedSections"`
		Frees                string        `json:"frees"`
		HasFreeListenCard    int           `json:"hasFreeListenCard"`
		IsConsumerUser       int           `json:"isConsumerUser"`
		LimitAmountTicket    string        `json:"limitAmountTicket"`
		Msg                  string        `json:"msg"`
		OfflineTime          int           `json:"offlineTime"`
		PayType              int           `json:"payType"`
		Price                int           `json:"price"`
		PriceType            int           `json:"priceType"`
		PriceUnit            string        `json:"priceUnit"`
		ResIds               interface{}   `json:"resIds"`
		RollAdUnlock         int           `json:"rollAdUnlock"`
		Sections             int           `json:"sections"`
		ShowDeadlineTime     int           `json:"showDeadlineTime"`
		Status               int           `json:"status"`
		Strategy             int           `json:"strategy"`
		TicketLimit          interface{}   `json:"ticketLimit"`
		UnlockEndTime        int           `json:"unlockEndTime"`
		UnlockLeftSectionNum int           `json:"unlockLeftSectionNum"`
		UnlockMaxSectionNum  int           `json:"unlockMaxSectionNum"`
		UsedTicket           interface{}   `json:"usedTicket"`
		VipExclusive         int           `json:"vipExclusive"`
		VipMinimumPrice      int           `json:"vipMinimumPrice"`
		WaitOffline          int           `json:"waitOffline"`
	} `json:"priceInfo"`
	RankingInfo           string `json:"rankingInfo"`
	RankingTarget         string `json:"rankingTarget"`
	ReceiveResourceUpdate int    `json:"receiveResourceUpdate"`
	RefId                 int    `json:"refId"`
	Rewarded              int    `json:"rewarded"`
	RollAdUnlock          int    `json:"rollAdUnlock"`
	Sections              int    `json:"sections"`
	ShowFreeEndTime       int    `json:"showFreeEndTime"`
	Sort                  int    `json:"sort"`
	State                 int    `json:"state"`
	Status                int    `json:"status"`
	Strategy              int    `json:"strategy"`
	Tags                  []struct {
		BgColor string `json:"bgColor"`
		Name    string `json:"name"`
		Type    int    `json:"type"`
	} `json:"tags"`
	TmeId      int         `json:"tmeId"`
	TmeType    int         `json:"tmeType"`
	TtsRef     interface{} `json:"ttsRef"`
	TtsType    int         `json:"ttsType"`
	Type       string      `json:"type"`
	TypeId     int         `json:"typeId"`
	Update     string      `json:"update"`
	UpdateTime int64       `json:"updateTime"`
	User       interface{} `json:"user"`
	Users      []struct {
		Cover     string `json:"cover"`
		Desc      string `json:"desc"`
		Flag      int    `json:"flag"`
		IsFollow  int    `json:"isFollow"`
		IsV       int    `json:"isV"`
		NickName  string `json:"nickName"`
		UserId    int    `json:"userId"`
		UserState int    `json:"userState"`
	} `json:"users"`
}

type detailFolderListResult struct {
	CollectionCount int           `json:"collectionCount"`
	EntityCount     int           `json:"entityCount"`
	EntityList      []interface{} `json:"entityList"`
	FolderId        int64         `json:"folderId"`
	HeadPic         string        `json:"headPic"`
	Name            string        `json:"name"`
	NickName        string        `json:"nickName"`
}

type detailRecommendList struct {
	Announcer      string `json:"announcer"`
	Cover          string `json:"cover"`
	Desc           string `json:"desc"`
	Hot            int    `json:"hot"`
	Id             int    `json:"id"`
	Name           string `json:"name"`
	PayType        int    `json:"payType"`
	RecReason      string `json:"recReason"`
	ShortRecReason string `json:"shortRecReason"`
	Strategy       int    `json:"strategy"`
	Tags           []struct {
		BgColor string `json:"bgColor"`
		Name    string `json:"name"`
		Type    int    `json:"type"`
	} `json:"tags"`
	Type int `json:"type"`
}

type detailResponse struct {
	ApiStatus int `json:"apiStatus"`
	Data      struct {
		BookDetail    detailBookDetailResult   `json:"bookDetail"`
		FolderList    []detailFolderListResult `json:"folderList"`
		RecommendList []detailRecommendList    `json:"recommendList"`
		RelateGroup   interface{}              `json:"relateGroup"`
	} `json:"data"`
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

func Detail(bookId string, options DetailOptions) error {
	query := downloadInitQuery(bookId, downloadDetailHttpQuery)

	c := colly.NewCollector()

	// Find and visit all links
	c.OnResponse(func(resp *colly.Response) {
		detailOnResponse(resp)
	})

	c.OnError(func(response *colly.Response, err error) {
		fmt.Printf("Request URL: %s \nfailed with response: %s\nError: %s\n", response.Request.URL, string(response.Body), err)
	})

	query.Set("sc", paramsMd5(query, urlPathBookDetail, secret))

	urlEntity, parseErr := url.Parse(domainDownload + urlPathBookDetail)
	if parseErr != nil {
		return parseErr
	}
	urlEntity.RawQuery = query.Encode()

	return c.Request("GET", urlEntity.String(), nil, nil, downloadHttpHeaders)
}

func detailOnResponse(resp *colly.Response) error {
	if resp.StatusCode != 200 {
		return errors.New("search response error")
	}
	result := detailResponse{}
	err := json.Unmarshal(resp.Body, &result)
	if err != nil {
		return err
	}
	tableData := [][]string{}
	tableData = append(tableData, []string{"书名", result.Data.BookDetail.Name})
	tableData = append(tableData, []string{"book_id", cast.ToString(result.Data.BookDetail.Id)})
	tableData = append(tableData, []string{"作者", result.Data.BookDetail.Author})
	tableData = append(tableData, []string{"主播", result.Data.BookDetail.Announcer})
	tableData = append(tableData, []string{"描述", result.Data.BookDetail.Desc})
	tableData = append(tableData, []string{"类型", result.Data.BookDetail.FatherTypeName + " > " + result.Data.BookDetail.Type})
	tableData = append(tableData, []string{"更新日期", result.Data.BookDetail.Update})

	outputTable(nil, tableData)
	return nil
}
