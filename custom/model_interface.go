package custom

/*
 * @Author: xiaozuhui
 * @Date: 2022-11-09 11:20:31
 * @LastEditors: xiaozuhui
 * @LastEditTime: 2022-11-09 13:28:25
 * @Description:
 */

type IModel interface {
	ModelToEntity() interface{}
}

// DBs
/**
 * @description: 将多个model批量转为entity
 * @param {[]IModel} ms
 * @return {[]interface{}}
 * @author: xiaozuhui
 */
func DBs(ms []IModel) []interface{} {
	var es = make([]interface{}, 0)
	for _, m := range ms {
		es = append(es, m.ModelToEntity())
	}
	return es
}
