package apigen

/**
* @api {get} /pokka/:id pokka
* @apiName 获取指定Pokka
* @apiVersion 0.1.0
* @apiGroup Pokka
* @apiDescription 这是描述信息，可以有多行。
* @apiExample {curl} 接口示例:
* curl -i http://localhost/pokka/4711
* @apiHeader {String} access-key 请求头必须携带字段access-key
* @apiHeaderExample {json} 头部示例:
* {
* "access-key": "按照约定加密方式产生的token=="
* }
*
* @apiSuccess (200) {String} firstname 姓氏
* @apiSuccess (200) {String} lastname 名称
*
* @apiSuccessExample {json} 成功的响应:
* HTTP/1.1 200 OK
* {
* "firstname": "John",
* "lastname": "Doe"
* }
*
 */
