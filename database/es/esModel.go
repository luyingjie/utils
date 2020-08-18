package es

import (
	"encoding/json"
	"strings"

	// "github.com/olivere/elastic"
	"gopkg.in/olivere/elastic.v3"
)

//API请求返回体模型
type RequestESModel struct {
	Limit       int
	Offset      int
	Operation   map[string]map[string]string
	Term        map[string]string
	QueryString map[string]string
	Wildcard    map[string]string
	Sort        map[string]bool
	Range       map[string]map[string]string
	Aggs        map[string]map[string]string
}

type HitModel struct {
	Score  *float64
	Id     string
	Source interface{}
}

type ResponseESModel struct {
	Total        int
	Hits         []HitModel
	Aggregations map[string][]map[string]interface{}
}

//mapper
func (*ResponseESModel) SearchHitsToResponseModel(sr *elastic.SearchResult) *ResponseESModel {
	model := new(ResponseESModel)
	if sr.Hits != nil {
		model.Total = int(sr.Hits.TotalHits)

		hitsModel := []HitModel{}
		for _, hit := range sr.Hits.Hits {
			hitModel := new(HitModel)
			hitModel.Id = hit.Id
			hitModel.Score = hit.Score
			hitModel.Source = hit.Source
			// hitsModel[i] = hitModel
			hitsModel = append(hitsModel, *hitModel)
		}
		model.Hits = hitsModel
		// 其实sr.Aggregations可以直接map[string]interface{}丢出去，但是那样就没有结果模型转化，冗余字段也很多，但是运行速度会快一点。
		if sr.Aggregations != nil {
			aggsModel := map[string][]map[string]interface{}{}
			for name, aggs := range sr.Aggregations {
				t := map[string]interface{}{}
				json.Unmarshal(*aggs, &t)
				model := []map[string]interface{}{}
				for _, buck := range t["buckets"].([]interface{}) {
					m := map[string]interface{}{}
					m["Key"] = buck.(map[string]interface{})["key"]
					m["Count"] = buck.(map[string]interface{})["doc_count"]
					model = append(model, m)
				}
				aggsModel[name] = model
			}
			model.Aggregations = aggsModel
		}
	}
	return model
}

func (*ResponseESModel) GetResultToResponseModel(gr *elastic.GetResult) *ResponseESModel {
	model := new(ResponseESModel)
	model.Total = 1

	hitsModel := []HitModel{}

	hitModel := new(HitModel)
	hitModel.Id = gr.Id
	hitModel.Source = gr.Source
	// hitsModel[i] = hitModel
	hitsModel = append(hitsModel, *hitModel)

	model.Hits = hitsModel
	return model
}

// SetESSelectValue : ES查询条件的判断
func SetESSelectValue(val string) (string, string) {
	// strings.Contains("helloogo", "hello")
	// strings.Replace(license, "\\n", "", -1)

	t := "Wildcard" // QueryString
	is := strings.ContainsAny(val, "AND&OR")
	if is {
		t = "QueryString"
		val = SetESQueryStringValue(val)
	}

	return t, val
}

// SetESQueryStringValue : 设置查询空格
func SetESQueryStringValue(val string) string {
	val = strings.Replace(val, " AND ", "AND", -1)
	val = strings.Replace(val, " OR ", "OR", -1)
	val = strings.Replace(val, " ", "\\\\ ", -1)
	val = strings.Replace(val, "AND", " AND ", -1)
	val = strings.Replace(val, "OR", " OR ", -1)
	return val
}
