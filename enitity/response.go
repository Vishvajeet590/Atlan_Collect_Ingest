package entity

/*type response struct {
	QuestionId int    `json:"question_id"`
	Response   string `json:"response"`
}
*/
type Responses struct {
	UserId   int            `json:"user_id"`
	Response map[int]string `json:"responses"`
}
