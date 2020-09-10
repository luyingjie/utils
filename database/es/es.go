package es

import (
	// "context"
	// . "encoding/json"
	"strings"
	"time"
	"utils/error"

	// "fmt"

	"github.com/golibs/uuid"

	// "github.com/olivere/elastic"
	"gopkg.in/olivere/elastic.v3"
)

//连接ES
func connect(esUrl, esIndex string) *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetURL(esUrl),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second))
	if err != nil {
		error.Try(2000, 3, "utils/database/es/es/connect/NewClient", err)
		return nil
	}

	exists, exErr := client.IndexExists(esIndex).Do()
	if exErr != nil {
		error.Try(2000, 3, "utils/database/es/es/connect/IndexExists", exErr)
		return nil
	}
	if !exists {
		createIndex, creErr := client.CreateIndex(esIndex).Do()
		if creErr != nil {
			error.Try(2000, 3, "utils/database/es/es/connect/CreateIndex", creErr)
			return nil
		}
		if !createIndex.Acknowledged {
			return nil
		}
	}

	return client
}

//保存
func Save(esUrl, esIndex, esType string, data map[string]interface{}) {
	client := connect(esUrl, esIndex)
	defer client.Stop()
	_, err := client.Index().
		Index(esIndex).
		Type(esType).
		BodyJson(data).
		Refresh(true).
		Do()

	if err != nil {
		error.Try(2000, 3, "utils/database/es/es/Save/Index", err)
	}
}

//uuid在什么地方生成需要研究。1.全局生成，贯穿整个Task。2.局部生成，返回传参。
//将整个库往其他库移动的保存情况要去掉模型中的id。
func SaveById(esUrl, esIndex, esType, id string, data map[string]interface{}) string {
	// var id UUID = uuid.Rand()
	// fmt.Println(id.Hex())
	// fmt.Println(id.Raw())

	// id1, err := uuid.FromStr("1870747d-b26c-4507-9518-1ca62bc66e5d")
	// id2 := uuid.MustFromStr("1870747db26c450795181ca62bc66e5d")
	// fmt.Println(id1 == id2) // true

	if id == "" {
		// id = uuid.Rand().Hex()
		id = strings.Replace(uuid.Rand().Hex(), "-", "", -1)
	}

	client := connect(esUrl, esIndex)
	defer client.Stop()
	_, err := client.Index().
		Index(esIndex).
		Type(esType).
		Id(id).
		BodyJson(data).
		Refresh(true).
		Do()

	if err != nil {
		error.Try(2000, 3, "utils/database/es/es/SaveById/Index", err)
	}
	return id
}

//按Id查询
func SearchById(esUrl, esIndex, esType, id string) *elastic.GetResult {
	client := connect(esUrl, esIndex)
	defer client.Stop()
	searchResult, err := client.Get().
		Index(esIndex).
		Type(esType).
		Id(id).
		Do()

	if err != nil {
		error.Try(2000, 3, "utils/database/es/es/SearchById/Get", err)
		return nil
	}

	return searchResult
}

// func setOperation(name string, ope map[string]interface{},bq *elastic.BoolQuery) *elastic.BoolQuery {
// 		// q = q.Must(NewTermQuery("tag", "wow"))
// 		// q = q.MustNot(NewRangeQuery("age").From(10).To(20))
// 		// q = q.Filter(NewTermQuery("account", "1"))
// 		// q = q.Should(NewTermQuery("tag", "sometag"), NewTermQuery("tag", "sometagtag"))
// 		// must :: 多个查询条件的完全匹配,相当于 and。
// 		// must_not :: 多个查询条件的相反匹配，相当于 not。
// 		// should :: 至少有一个查询条件匹配, 相当于 or。
// 	if ope != nil && ope[name] != nil

// 	return bq
// }

//查询
//操作符目前只支持All指定，不支持单独逐个指定，数据结构保持扩展支持。
func Search(esUrl, esIndex, esType string, searchModel *models.RequestESModel) *elastic.SearchResult { //*elastic.SearchHits {
	client := connect(esUrl, esIndex)
	defer client.Stop()
	searchCommand := client.Search().
		Index(esIndex).
		Type(esType)

	boolQuery := elastic.NewBoolQuery()
	ope := "AND"
	//判断操作符，暂时不需要提取出来为方法。 目前还不全面，后面设计后再抽公共方法。
	if searchModel.Operation != nil && searchModel.Operation["All"] != nil {
		ope = searchModel.Operation["All"]["All"]
	}

	if searchModel.Term != nil {
		opeTerm := ""
		if searchModel.Operation != nil && searchModel.Operation["Term"] != nil {
			opeTerm = searchModel.Operation["Term"]["All"]
		}
		for key, val := range searchModel.Term {
			if key == "Id" {
				key = "_id"
			}
			//这里使用string的话Count的查询情况可能会不对，只是目前没有按Count查询。
			if key != "" && val != "" {
				//这里验证下是不是必须按小写查询，在ES6下转小写查询完全造成。
				//但是ES2下面需要测试，目前看Null的查询时不对的，需要替换成小写，为了避免歧义，换成了root。最后换成了null，现在不区分大小写，所以都是正常的,除了部分影响太大的还是沿用大写。
				// 这里因为命名有点混乱，有2个地方，一个是模块名称，因为影响较大。二是生成的id里面有大写，目前来看之前的ES6未关闭分词的版本因为是大写数据，传入的是小写所以查询正常，需要全部测试。
				// 这里的主要问题是有没有关闭分词，没关闭分词的话。比如Null，用Null查询不到，必须用null才可以，所有的大写字母都会有问题。
				// 应该是所有的都被分词转成了小写，可以在查询时将值转成小写，但是这样用户查询也无法区分大小写。关闭分词是最好的办法。
				// 这里按兼容性配置还调用。
				termQuery := elastic.NewTermQuery(key, val)

				switch opeTerm {
				case "AND":
					boolQuery.Must(termQuery)
				case "OR":
					boolQuery.Should(termQuery)
				case "":
					switch ope {
					case "AND":
						boolQuery.Must(termQuery)
					case "OR":
						boolQuery.Should(termQuery)
					}
				}
			}
		}
	}
	if searchModel.QueryString != nil {
		for key, val := range searchModel.QueryString {
			if key != "" && val != "" {
				QueryString := elastic.NewQueryStringQuery(val)
				QueryString.DefaultField(key)
				// searchCommand = searchCommand.Query(QueryString)
				switch ope {
				case "AND":
					boolQuery.Must(QueryString)
				case "OR":
					boolQuery.Should(QueryString)
				}
			}
		}
	}

	if searchModel.Wildcard != nil {
		for key, val := range searchModel.Wildcard {
			if key != "" && val != "" {
				Wildcard := elastic.NewWildcardQuery(key, val)
				switch ope {
				case "AND":
					boolQuery.Must(Wildcard)
				case "OR":
					boolQuery.Should(Wildcard)
				}
			}
		}
	}
	//这里有Bug,Gte这些方法接受的是interface而非string，将传值直接赋值会报错nil，必须直接给“”的字符串形式，这里需要研究。
	if searchModel.Range != nil {
		for key, m := range searchModel.Range {
			if key != "" {
				// _key := key
				// var gte,lte interface{}
				rangeQuery := elastic.NewRangeQuery(key)
				for o, val := range m {
					// fmt.Println(reflect.TypeOf(val).String())
					// fmt.Println(reflect.TypeOf("2018-03-14").String())
					if o != "" && val != "" {
						switch o {
						case "Gte":
							rangeQuery = rangeQuery.Gte(val) //Gte("2018-03-14") //rangeQuery.Gte(val)
						case "Gt":
							rangeQuery = rangeQuery.Gt(val)
						case "Lte":
							rangeQuery = rangeQuery.Lte(val) //Lte("2018-03-14") //rangeQuery.Lte(val)
						case "Lt":
							rangeQuery = rangeQuery.Lt(val)

							// rangeQuery := elastic.NewRangeQuery("SystemTime")
							// 	Gte("2018-03-15").
							// 	Lte("now")
							// 	// TimeZone("+1:00")
							// searchCommand = searchCommand.Query(rangeQuery)
						}
					}
				}

				// rangeQuery := elastic.NewRangeQuery(_key).
				// 			Gte(gte).Lte(lte)
				//Linux 下HH:mm:ss报错，需要hh:mm:ss windows则相反,后测试HH:mm:ss
				//在win和linux下都正常，估计之前不正常和格式有关。
				//.TimeZone("+8:00")是不行的，数据条数不同应该是一位UTC时间的问题，应该需要调整-8小时
				// searchCommand = searchCommand.Query(rangeQuery.Format("yyyy-MM-dd HH:mm:ss"))
				switch ope {
				case "AND":
					boolQuery.Must(rangeQuery.Format("yyyy-MM-dd HH:mm:ss"))
				case "OR":
					boolQuery.Should(rangeQuery.Format("yyyy-MM-dd HH:mm:ss"))
				}
			}
		}
	}
	if searchModel.Sort != nil {
		for key, val := range searchModel.Sort {
			if key != "" {
				//这里多个排序形式需要研究下
				searchCommand = searchCommand.Sort(key, val)
			}
		}
	}

	if searchModel.Aggs != nil {
		for key, val := range searchModel.Aggs {
			agg := elastic.NewTermsAggregation().Field(val["Field"])
			searchCommand = searchCommand.Aggregation(key, agg)
		}
	}

	searchCommand = searchCommand.Query(boolQuery)
	if searchModel.Limit != 0 {
		searchCommand = searchCommand.From(searchModel.Offset).Size(searchModel.Limit)
	}

	searchResult, err := searchCommand.Pretty(true).Do()
	if err != nil {
		error.Try(2000, 3, "utils/database/es/es/Search/Do", err)
		return nil
	}

	// return searchResult.Hits
	return searchResult
}

//按Id删除,补全查询删除。
func DeleteById(esUrl, esIndex, esType, id string) {
	client := connect(esUrl, esIndex)
	defer client.Stop()
	_, err := client.Delete().
		Index(esIndex).
		Type(esType).
		Id(id).
		Refresh(true).
		Do()

	if err != nil {
		error.Try(2000, 3, "utils/database/es/es/DeleteById/Delete", err)
	}
}

// 按条件删除
func Delete(esUrl, esIndex, esType string, searchModel *models.RequestESModel) {
	client := connect(esUrl, esIndex)
	defer client.Stop()
	searchCommand := client.DeleteByQuery().
		Index(esIndex).
		Type(esType)

	boolQuery := elastic.NewBoolQuery()
	ope := "AND"
	if searchModel.Operation != nil && searchModel.Operation["All"] != nil {
		ope = searchModel.Operation["All"]["All"]
	}

	if searchModel.Term != nil {
		opeTerm := ""
		if searchModel.Operation != nil && searchModel.Operation["Term"] != nil {
			opeTerm = searchModel.Operation["Term"]["All"]
		}
		for key, val := range searchModel.Term {
			if key == "Id" {
				key = "_id"
			}
			//这里使用string的话Count的查询情况可能会不对，只是目前没有按Count查询。
			if key != "" && val != "" {
				termQuery := elastic.NewTermQuery(key, val)

				switch opeTerm {
				case "AND":
					boolQuery.Must(termQuery)
				case "OR":
					boolQuery.Should(termQuery)
				case "":
					switch ope {
					case "AND":
						boolQuery.Must(termQuery)
					case "OR":
						boolQuery.Should(termQuery)
					}
				}
			}
		}
	}

	if searchModel.Wildcard != nil {
		for key, val := range searchModel.Wildcard {
			if key != "" && val != "" {
				Wildcard := elastic.NewWildcardQuery(key, val)
				switch ope {
				case "AND":
					boolQuery.Must(Wildcard)
				case "OR":
					boolQuery.Should(Wildcard)
				}
			}
		}
	}

	if searchModel.Range != nil {
		for key, m := range searchModel.Range {
			if key != "" {
				// _key := key
				// var gte,lte interface{}
				rangeQuery := elastic.NewRangeQuery(key)
				for o, val := range m {
					// fmt.Println(reflect.TypeOf(val).String())
					// fmt.Println(reflect.TypeOf("2018-03-14").String())
					if o != "" && val != "" {
						switch o {
						case "Gte":
							rangeQuery = rangeQuery.Gte(val) //Gte("2018-03-14") //rangeQuery.Gte(val)
						case "Gt":
							rangeQuery = rangeQuery.Gt(val)
						case "Lte":
							rangeQuery = rangeQuery.Lte(val) //Lte("2018-03-14") //rangeQuery.Lte(val)
						case "Lt":
							rangeQuery = rangeQuery.Lt(val)

							// rangeQuery := elastic.NewRangeQuery("SystemTime")
							// 	Gte("2018-03-15").
							// 	Lte("now")
							// 	// TimeZone("+1:00")
							// searchCommand = searchCommand.Query(rangeQuery)
						}
					}
				}

				switch ope {
				case "AND":
					boolQuery.Must(rangeQuery.Format("yyyy-MM-dd HH:mm:ss"))
				case "OR":
					boolQuery.Should(rangeQuery.Format("yyyy-MM-dd HH:mm:ss"))
				}
			}
		}
	}

	_, err := searchCommand.Query(boolQuery).Do()

	if err != nil {
		error.Try(2000, 3, "utils/database/es/es/Delete/Do", err)
	}

}

//按Id修改,这里其实支持直接插入map[string]interface{},可以和前面的Save方法合并。分开是可以寻求其他修改文档的方法。
func UpdateById(esUrl, esIndex, esType, id string, data map[string]interface{}) {
	client := connect(esUrl, esIndex)
	defer client.Stop()
	_, err := client.Index().
		Index(esIndex).
		Type(esType).
		Id(id).
		BodyJson(data).
		// BodyString(json.MapToString(data)).
		Refresh(true).
		Do()

	if err != nil {
		error.Try(2000, 3, "utils/database/es/es/UpdateById/Index", err)
	}
}
