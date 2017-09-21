package api

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/Sirupsen/logrus"
	"VCMS/apps/libs"
	"github.com/jinzhu/gorm"
	"VCMS/apps/handler"
	rdbModel 	"VCMS/apps/models/rdb"
	apiModel 	"VCMS/apps/models/api"
	appConfig 	"VCMS/apps/config"
	"strconv"
)

type UserAPI struct {
	request UserRequest
}

type UserRequest struct {
	TenantId	string	`validate:"required" json:"tenant_id"`
	UserId 		string	`validate:"required" json:"user_id"`
	UserName	string	`validate:"required" json:"name"`
	UserRole	string	`validate:"required" json:"role"`
	TelNo		string	`json:"tel_no"`
	Email		string	`validate:"email" json:"email"`
	Password	string	`json:"password"`
	UseYN		string 	`validate:"required" json:"use_yn"`
}

// 사용자 등록
func (api UserAPI) PostUser (c echo.Context) error  {

	payload := &api.request
	c.Bind(payload)

	if err   := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	enc_pass := libs.Sha256Encoding(payload.Password)

	claims 	 := apiModel.GetJWTClaims(c)

	logrus.Debug("#### TenantId : ", claims.TenantId)

	user := &rdbModel.User{
		TenantID	:payload.TenantId,
		UserId		:payload.UserId,
		Name		:payload.UserName,
		Role		:payload.UserRole,
		Telno		:payload.TelNo,
		Email		:payload.Email,
		Password	:enc_pass,
		UseYN		:"Y",
		UpdatedId	:claims.UserId,
	}

	tx := c.Get("Tx").(*gorm.DB)

	if !tx.Where("tenant_id = ? and user_id = ?",
		user.TenantID, user.UserId ).Find(user).RecordNotFound() {
		return echo.NewHTTPError(http.StatusConflict, "Already exists User ")
	}

	tx.Create(user)

	return handler.APIResultHandler(c, true, http.StatusCreated, user)
}

// 사용자 한건 조회
func (api UserAPI) GetUserByIdx (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	user := &rdbModel.User{}

	//var count = 0
	tx := c.Get("Tx").(*gorm.DB)

	if tx.Preload("Tenant").Where("idx = ? ", idx ).Find(user).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "User NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, user)

}

func (api UserAPI) GetUserByUserId (c echo.Context) error  {

	user_id, _ := strconv.ParseInt(c.Param("user_id"), 0, 64)

	user := &rdbModel.User{}

	//var count = 0
	tx := c.Get("Tx").(*gorm.DB)

	if tx.Preload("Tenant").Where("user_id = ? ", user_id ).Find(user).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "User NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, user)

}



// 사용자 리스트 조회
func (api UserAPI) GetUsers (c echo.Context) error  {

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	start_date		:= c.QueryParam("start_date")
	end_date		:= c.QueryParam("end_date")

	tenant_id		:= c.QueryParam("tenant_id")
	role			:= c.QueryParam("role");

	search_string 		:= c.QueryParam("search_string")

	sort_name		:= c.QueryParam("sort_name")
	sort_dir		:= c.QueryParam("sort_dir")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### start_date : ", start_date)
	logrus.Debug("##### end_date : ", end_date)
	logrus.Debug("##### search_string : ", search_string)

	tenants := []rdbModel.User{}
	tx := c.Get("Tx").(*gorm.DB)

	//if tenant_id == "" {
	//	return echo.NewHTTPError(http.StatusBadRequest, "TENANT_ID Required Parameter")
	//}

	if tenant_id != "" {
		tx = tx.Where("tenant_id = ? ", tenant_id)
	}

	if role != "" {
		tx = tx.Where("role = ? ", role)
	}

	if search_string != "" {
		tx = tx.Where("user_id LIKE ? or name LIKE ? ", "%" + search_string + "%", "%" + search_string + "%")
	}

	if start_date != "" {
		tx = tx.Where("created_at BETWEEN ? AND ?", start_date, end_date)
	}

	if sort_name == "" {
		sort_name = "idx"
		sort_dir  = "DESC"
	}

	var count = 0
	tx.Find(&tenants).Count(&count)

	tx.Preload("Tenant").Order(sort_name +" "+sort_dir).Offset(offset).Limit(limit).Find(&tenants)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{"total_count": count,
			"rows": tenants})

}

// 사용자 수정
func (api UserAPI) PutUser (c echo.Context) error  {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	payload := &api.request
	c.Bind(payload)

	enc_pass := libs.Sha256Encoding(payload.Password)

	if err   := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tx := c.Get("Tx").(*gorm.DB)

	user := &rdbModel.User{}

	if payload.Password != "" {
		if tx.Model(user).Where("idx = ? ", idx).
			Updates(rdbModel.User{Name: payload.UserName,
			Role:payload.UserRole,
			Telno:payload.TelNo,
			Email:payload.Email,
			Password:enc_pass,
			UseYN:payload.UseYN,
		}).RowsAffected == 0 {

			return echo.NewHTTPError(http.StatusNotFound, "User NOT FOUND")
		}
	} else {
		if tx.Model(user).Where("idx = ? ", idx).
			Updates(rdbModel.User{Name: payload.UserName,
			Role:payload.UserRole,
			Telno:payload.TelNo,
			Email:payload.Email,
			UseYN:payload.UseYN,
		}).RowsAffected == 0 {

			return echo.NewHTTPError(http.StatusNotFound, "User NOT FOUND")
		}
	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Update OK")

}


// 사용자 삭제
func (api UserAPI) DeleteUser (c echo.Context) error  {

	id, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	user := &rdbModel.User{}

	tx := c.Get("Tx").(*gorm.DB)

	if tx.Delete(user, "idx = ?", id).RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "User NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
}


// 사용자 리스트 조회
func (api UserAPI) GetHelloTUsers (c echo.Context) error  {

	offset, _ 		:= strconv.Atoi(c.QueryParam("offset"))
	limit, _ 		:= strconv.Atoi(c.QueryParam("limit"))

	search_string 		:= c.QueryParam("search_string")

	logrus.Debug("##### offset : ", offset)
	logrus.Debug("##### limit : ", limit)
	logrus.Debug("##### search_string : ", search_string)

	hellotUsers := []rdbModel.MemberDefault{}

	tx, err := gorm.Open(appConfig.Config.RDB[1].Product, appConfig.Config.RDB[1].ConnectString)
	if err != nil {
		panic(err)
	}
	defer tx.Close()
	tx.LogMode(appConfig.Config.RDB[1].Debug)

	if search_string != "" {
		//tx = tx.Where("member_id LIKE ? or member_name LIKE ? ", "%" + search_string + "%", "%" + search_string + "%")
		tx = tx.Where("member_id LIKE ? ", "%" + search_string + "%")
	}

	var count = 0
	tx.Find(&hellotUsers).Count(&count)

	tx.Offset(offset).Limit(limit).Find(&hellotUsers)

	return handler.APIResultHandler(c, true, http.StatusOK,
		map[string]interface{}{"total_count": count,
			"rows": hellotUsers})

}