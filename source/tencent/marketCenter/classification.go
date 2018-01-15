package marketCenter

import (
	"stock"
	"regexp"
	"encoding/json"
)

type Class []struct {
	//分类等级
	Level uint8
	//分类id
	Id string
	//分类名称
	Name string
	//首字母
	Initials string
	//子级分类详情
	SubClass Class
	//子级分类id列表
	LastSubClass map[string]string
}

//解析分类结构
type classEnJson struct {
	//子级分类id
	Chd []string `json:"chd"`
	//
	Clk string `json:"clk"`
	//分类名称首字母
	Fl string `json:"fl"`
	//分类id
	ID string `json:"id"`
	//分类等级
	Idx uint8 `json:"idx"`
	//上级分类id
	Pt interface{} `json:"pt"`
	//分类名称
	T string `json:"t"`
}

//读取分类
func ReaderClass() (Class , error) {
	classBody, err := downloadAllClass()
	if err != nil {
		return nil , err
	}
	classT , err := findAllClass(classBody)
	if err != nil {
		return nil , err
	}
	return classT , nil
}

//下载行情中心分类数据
func downloadAllClass() ([]byte, error) {
	out, err := stock.Get("http://stockapp.finance.qq.com/mstats/?mod=all")
	if err != nil {
		return nil, err
	}
	return out, nil
}

//查找行情中心所有有效分类
func findAllClass(b []byte) (Class, error) {
	//匹配json分类数据得到分类结构
	r := regexp.MustCompile(`menuList : (.*?),\n`)
	result := r.FindSubmatch(b)
	/**
	 * 优化json数据，让json达到golang.json库解析标准
	 */
	//去除对于id数据
	r = regexp.MustCompile(`"[a-z A-Z _ 0-9]+":{"id":`)
	result[1] = r.ReplaceAll(result[1], []byte("{\"id\":"))
	//替换json开始和结束括号 {{...},{...}} => [{...},{...}]
	r = regexp.MustCompile(`^{`)
	result[1] = r.ReplaceAll(result[1], []byte("["))
	//替换结束括号
	r = regexp.MustCompile(`}$`)
	result[1] = r.ReplaceAll(result[1], []byte("]"))
	//解析json
	var classJson []classEnJson
	err := json.Unmarshal(result[1], &classJson)
	if err != nil {
		return nil, err
	}
	//分类ID => classEnJson 用于快速查找定位分类详细信息
	idToClass := make(map[string]classEnJson, len(classJson))
	//记录顶级分类变量
	topLevelId := []string{}
	for _, v := range classJson {
		idToClass[v.ID] = v
		//记录顶级分类
		if v.Idx == 1 {
			topLevelId = append(topLevelId, v.ID)
		}
	}
	out := make(Class, len(topLevelId))
	LastClassNames , err := readerLastClassName(b)
	if err != nil {
		return nil , err
	}
	//根据顶级分类查找下级所有分类
	for k, v := range topLevelId {
		//第一级分类
		out[k].Level = 1
		out[k].Id = idToClass[v].ID
		out[k].Name = idToClass[v].T
		out[k].Initials = idToClass[v].Fl
		out[k].SubClass = make(Class, len(idToClass[v].Chd))
		for key, val := range idToClass[v].Chd {
			//第二集分类
			out[k].SubClass[key].Level = 2
			out[k].SubClass[key].Id = idToClass[val].ID
			out[k].SubClass[key].Name = idToClass[val].T
			out[k].SubClass[key].Initials = idToClass[val].Fl
			out[k].SubClass[key].LastSubClass = make(map[string]string, len(idToClass[val].Chd))
			//整理下级分类id=>name
			for _, value := range idToClass[val].Chd {
				tue, ok := idToClass[value]
				name := ""
				if ok {
					name = tue.T
				}else {
					name , ok = LastClassNames[value]
					if !ok {
						name = ""
					}
				}
				out[k].SubClass[key].LastSubClass[value] = name
			}
		}
		//整理下级分类id=>name
		out[k].LastSubClass = make(map[string]string, len(idToClass[v].Chd))
		for _, value := range idToClass[v].Chd {
			tue, ok := idToClass[value]
			name := ""
			if ok {
				name = tue.T
			}
			out[k].LastSubClass[value] = name
		}
	}
	return out, nil
}

//提取3级分类名称
func readerLastClassName(b []byte) (map[string]string, error) {
	r := regexp.MustCompile(`<div id="alllist">\n(.*?)+\n`)
	result := r.FindSubmatch(b)
	r = regexp.MustCompile(`<li><a\s+class="clk-mo-li"\s+id="a-l-(.+?)"\s+href="\?id=(.+?)"\s*.*?>(.+?)</a></li>`)
	names := r.FindAllSubmatch(result[0], -1)
	out := make(map[string]string, len(names))
	for _, v := range names {
		out[string(v[1])] = string(v[3])
	}
	//os.Exit(-1)
	return out, nil
}
