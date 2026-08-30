package search_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"
	"blogx_server/utils/sql"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type ArticleSearchRequest struct {
	common.PageInfo
	Tag  string `form:"tag"`
	Type int8   `form:"type"` // 0猜你喜欢 1最新发布 2最多回复 3最多点赞 4最多收藏
}

type ArticleBaseInfo struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Abstract string `json:"abstract"`
}

type ArticleListResponse struct {
	models.ArticleModel
	AdminTop      bool    `json:"adminTop"`
	CategoryTitle *string `json:"categoryTitle"`
	UserNickname  string  `json:"userNickname"`
	UserAvatar    string  `json:"userAvatar"`
}

func (SearchApi) ArticleSearchView(c *gin.Context) {
	cr := middleware.GetBind[ArticleSearchRequest](c)

	sortMap := map[int8]string{
		0: "_score",
		1: "created_at",
		2: "comment_count",
		3: "digg_count",
		4: "collect_count",
	}
	sortKey := sortMap[cr.Type]
	if sortKey == "" {
		res.FailWithMsg("搜索类型错误", c)
		return
	}

	query := elastic.NewBoolQuery()
	if cr.Keyword != "" {
		query.Should(
			elastic.NewMatchQuery("title", cr.Keyword),
			elastic.NewMatchQuery("abstract", cr.Keyword),
			elastic.NewMatchQuery("content", cr.Keyword),
		)
	}

	if cr.Tag != "" {
		query.Must(elastic.NewTermQuery("tag_list", cr.Tag))
	}

	// 只能查发布的文章
	query.Must(elastic.NewTermsQuery("status", 3))

	// 给后面查mysql构造条件用的
	var articleIDList []uint

	// 把管理员置顶的查出来
	var userIDList []uint
	var topArticleIDList []uint
	global.DB.Model(&models.UserModel{}).Where("role = ?", enum.AdminRole).Select("id").Scan(&userIDList)
	global.DB.Model(&models.UserTopArticleModel{}).Where("user_id in ?", userIDList).Select("article_id").Scan(&topArticleIDList)
	var articleTopMap = map[uint]bool{}
	if len(topArticleIDList) > 0 {
		var topArticleIDListAny []interface{}
		for _, u := range topArticleIDList {
			topArticleIDListAny = append(topArticleIDListAny, u)
			articleTopMap[u] = true
			articleIDList = append(articleIDList, u)
		}
		query.Should(elastic.NewTermsQuery("id", topArticleIDListAny...))
	}

	if cr.Type == 0 {
		// 只有猜你喜欢，才会把用户喜欢的标签带入查询
		claims, err := jwt.ParseTokenByGin(c)
		if err == nil && claims != nil {
			var userConf models.UserConfModel
			if err := global.DB.Take(&userConf, "user_id = ?", claims.UserID).Error; err != nil {
				res.FailWithMsg("用户配置不存在", c)
				return
			}
			if len(userConf.LikeTags) > 0 {
				tagQuery := elastic.NewBoolQuery()
				var tagAnyList []interface{}
				for _, tag := range userConf.LikeTags {
					tagAnyList = append(tagAnyList, tag)
				}
				tagQuery.Should(elastic.NewTermsQuery("tag_list", tagAnyList...))
				query.Must(tagQuery)
			}
		}
	}

	highlight := elastic.NewHighlight()
	highlight.Field("title")
	highlight.Field("abstract")

	result, err := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(query).
		From(cr.GetOffset()).
		Size(cr.GetLimit()).
		Highlight(highlight).
		Sort(sortKey, false).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		res.FailWithMsg("ES查询失败", c)
	}

	source, _ := query.Source()
	byteData, _ := json.Marshal(source)
	fmt.Println("es查询语句", string(byteData))

	count := result.Hits.TotalHits.Value
	fmt.Println("count-->", count)

	searchArticleMap := map[uint]ArticleBaseInfo{}
	for _, hit := range result.Hits.Hits {
		var art ArticleBaseInfo
		err := json.Unmarshal(hit.Source, &art)
		if err != nil {
			logrus.Warnf("解析失败 %s  %s", err, string(hit.Source))
			continue
		}
		if len(hit.Highlight["title"]) > 0 {
			art.Title = hit.Highlight["title"][0]
		}
		if len(hit.Highlight["abstract"]) > 0 {
			art.Abstract = hit.Highlight["abstract"][0]
		}
		if hit.Score != nil {
			fmt.Println(*hit.Score, art.Title, art.ID)
		}
		searchArticleMap[art.ID] = art
		articleIDList = append(articleIDList, art.ID)
	}

	where := global.DB.Where("")
	if len(articleIDList) > 0 {
		where = global.DB.Where("id in ?", articleIDList)
	}
	_list, _, _ := common.ListQuery(models.ArticleModel{}, common.Options{
		Where:        where,
		PreLoads:     []string{"CategoryModel", "UserModel"},
		DefaultOrder: sql.ConvertSliceOrderSql(articleIDList),
	})

	list := make([]ArticleListResponse, 0)
	for _, model := range _list {
		item := ArticleListResponse{
			ArticleModel: model,
			AdminTop:     articleTopMap[model.ID],
			UserNickname: model.UserModel.Nickname,
			UserAvatar:   model.UserModel.Avatar,
		}
		if item.CategoryModel != nil {
			item.CategoryTitle = &item.CategoryModel.Title
		}
		item.Title = searchArticleMap[model.ID].Title
		item.Abstract = searchArticleMap[model.ID].Abstract
		list = append(list, item)
	}

	res.OkWithList(list, int(count), c)
}
