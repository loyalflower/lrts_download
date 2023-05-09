package lrts_download

import (
	"encoding/json"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gocolly/colly"
	"github.com/spf13/cast"
	"net/url"
	"time"
)

type DetailOptions struct {
	CommonOptions
}

const (
	urlPathBookDetail   = "/yyting/page/bookDetailPage.action"
	urlPathAblumnDetail = "/yyting/page/ablumnDetailPage.action"
)

var (
	errorNotFound = errors.New("not found")

	downloadDetailBookHttpQuery = url.Values{
		"bookId": []string{""},
		"token":  []string{defaultToken},
		"imei":   []string{defaultIMEI},
		"nwt":    []string{"1"},
		"q":      []string{"2150"},
	}

	downloadDetailAblumnHttpQuery = url.Values{
		"ablumnId": []string{""},
		"token":    []string{defaultToken},
		"imei":     []string{defaultIMEI},
		"nwt":      []string{"1"},
		"q":        []string{"2150"},
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

type detailBookResponse struct {
	commonApiResponse
	Data struct {
		BookDetail    detailBookDetailResult   `json:"bookDetail"`
		FolderList    []detailFolderListResult `json:"folderList"`
		RecommendList []detailRecommendList    `json:"recommendList"`
		RelateGroup   interface{}              `json:"relateGroup"`
	} `json:"data"`
}

type detailAblumnResponse struct {
	commonApiResponse
	Data struct {
		AblumnDetail struct {
			Ablumn struct {
				AdvertControlType  int           `json:"advertControlType"`
				AlbumType          int           `json:"albumType"`
				Announcer          string        `json:"announcer"`
				Author             string        `json:"author"`
				BaseEntityId       int           `json:"baseEntityId"`
				BaseEntityType     int           `json:"baseEntityType"`
				BestCover          string        `json:"bestCover"`
				CommentControlType int           `json:"commentControlType"`
				CommentCount       int           `json:"commentCount"`
				CommentMean        float64       `json:"commentMean"`
				Cover              string        `json:"cover"`
				Description        string        `json:"description"`
				EstimatedSections  int           `json:"estimatedSections"`
				ExtraInfos         []interface{} `json:"extraInfos"`
				FreeEndTime        int           `json:"freeEndTime"`
				Id                 int           `json:"id"`
				IsCollected        int           `json:"isCollected"`
				IsLike             int           `json:"isLike"`
				IsSend             int           `json:"isSend"`
				Labels             interface{}   `json:"labels"`
				Length             int           `json:"length"`
				Name               string        `json:"name"`
				OriginCover        string        `json:"originCover"`
				PlayCount          int           `json:"playCount"`
				PriceInfo          struct {
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
				Rewarded              int    `json:"rewarded"`
				RollAdUnlock          int    `json:"rollAdUnlock"`
				Sections              int    `json:"sections"`
				ShowFreeEndTime       int    `json:"showFreeEndTime"`
				Sort                  int    `json:"sort"`
				Source                int    `json:"source"`
				State                 int    `json:"state"`
				Strategy              int    `json:"strategy"`
				Tags                  []struct {
					BgColor string `json:"bgColor"`
					Name    string `json:"name"`
					Type    int    `json:"type"`
				} `json:"tags"`
				TmeId      int    `json:"tmeId"`
				TypeId     int    `json:"typeId"`
				TypeName   string `json:"typeName"`
				UpdateTime int64  `json:"updateTime"`
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
			} `json:"ablumn"`
			ApiStatus int    `json:"apiStatus"`
			Msg       string `json:"msg"`
			Status    int    `json:"status"`
			User      struct {
				Cover     string `json:"cover"`
				Desc      string `json:"desc"`
				Flag      int    `json:"flag"`
				IsFollow  int    `json:"isFollow"`
				IsV       int    `json:"isV"`
				NickName  string `json:"nickName"`
				UserId    int    `json:"userId"`
				UserState int    `json:"userState"`
			} `json:"user"`
		} `json:"ablumnDetail"`
		RecommendList []struct {
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
		} `json:"recommendList"`
		RelateGroup interface{} `json:"relateGroup"`
	} `json:"data"`
}

func Detail(bookId string, options DetailOptions) error {
	query := downloadInitQuery(bookId, downloadDetailBookHttpQuery)

	c := colly.NewCollector()

	// Find and visit all links
	c.OnResponse(func(resp *colly.Response) {
		err := detailOnResponse(resp)
		if errors.Is(err, errorNotFound) && resp.Request.URL.Path == urlPathBookDetail {
			albumQuery := downloadInitQuery("", downloadDetailAblumnHttpQuery)
			albumQuery.Set("ablumnId", bookId)
			urlEntity := initDetailUrl(domainDownload+urlPathAblumnDetail, albumQuery)
			if urlEntity == nil {
				return
			}
			c.Request("GET", urlEntity.String(), nil, nil, downloadHttpHeaders)
		}
	})

	c.OnError(func(response *colly.Response, err error) {
		fmt.Printf("Request URL: %s \nfailed with response: %s\nError: %s\n", response.Request.URL, string(response.Body), err)
	})

	urlEntity := initDetailUrl(domainDownload+urlPathBookDetail, query)
	if urlEntity == nil {
		return errors.New("请求创建失败")
	}

	return c.Request("GET", urlEntity.String(), nil, nil, downloadHttpHeaders)
}

func initDetailUrl(urlPath string, values url.Values) *url.URL {
	urlEntity, parseErr := url.Parse(urlPath)
	if parseErr != nil {
		return nil
	}
	values.Set("sc", paramsMd5(values, urlEntity.Path, secret))
	urlEntity.RawQuery = values.Encode()
	return urlEntity
}

func detailOnResponse(resp *colly.Response) error {
	if resp.StatusCode != 200 {
		return errors.New("search response error")
	}
	tableData := [][]string{}
	if resp.Request.URL.Path == urlPathBookDetail {
		result := detailBookResponse{}
		err := json.Unmarshal(resp.Body, &result)
		if err != nil {
			return err
		}

		if result.Data.BookDetail.Name == "" {
			return errorNotFound
		}

		tableData = append(tableData, []string{"书名", result.Data.BookDetail.Name})
		tableData = append(tableData, []string{"book_id", cast.ToString(result.Data.BookDetail.Id)})
		tableData = append(tableData, []string{"作者", result.Data.BookDetail.Author})
		tableData = append(tableData, []string{"主播", result.Data.BookDetail.Announcer})
		tableData = append(tableData, []string{"描述", result.Data.BookDetail.Desc})
		tableData = append(tableData, []string{"类型", result.Data.BookDetail.FatherTypeName + " > " + result.Data.BookDetail.Type})
		tableData = append(tableData, []string{"更新日期", result.Data.BookDetail.Update})
		tableData = append(tableData, []string{"集数", cast.ToString(result.Data.BookDetail.Sections)})
	} else {
		result := detailAblumnResponse{}
		err := json.Unmarshal(resp.Body, &result)
		if err != nil {
			return err
		}
		if result.Status != 0 {
			return errors.New(result.Msg)
		}

		if result.Data.AblumnDetail.Ablumn.Name == "" {
			return errorNotFound
		}

		tableData = append(tableData, []string{"书名", result.Data.AblumnDetail.Ablumn.Name})
		tableData = append(tableData, []string{"book_id", cast.ToString(result.Data.AblumnDetail.Ablumn.Id)})
		tableData = append(tableData, []string{"作者", result.Data.AblumnDetail.Ablumn.Author})
		tableData = append(tableData, []string{"主播", result.Data.AblumnDetail.Ablumn.Announcer})
		tableData = append(tableData, []string{"描述", result.Data.AblumnDetail.Ablumn.Description})
		tableData = append(tableData, []string{"类型", result.Data.AblumnDetail.Ablumn.TypeName})
		tableData = append(tableData, []string{"更新日期", time.UnixMilli(result.Data.AblumnDetail.Ablumn.UpdateTime).Format(time.DateOnly)})
		tableData = append(tableData, []string{"集数", cast.ToString(result.Data.AblumnDetail.Ablumn.Sections)})
	}
	outputTable(nil, tableData)
	return nil
}
