package api


type APIResult struct {
	Success 	bool 		`json:"success" bson:"success"`
	ResultCode  	int  		`json:"resultCode" bson:"resultCode"`
	ResultData 	interface{} 	`json:"resultData" bson:"resultData"`
}