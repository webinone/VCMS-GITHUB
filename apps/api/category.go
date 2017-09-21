package api

import (
	"VCMS/apps/handler"
	apiModel "VCMS/apps/models/api"
	rdbModel "VCMS/apps/models/rdb"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	//"strings"
	//"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
)

type CategoryAPI struct {
	request CategoryRequest
}

type CategoryRequest struct {
	CategoryName string `validate:"required" json:"name"`
	CategoryDesc string `json:"desc"`
	ParentIdx    int64  `validate:"required" json:"parent_idx"`
	UseYN        string `validate:"required" json:"use_yn"`
}

// Tree에서 사용하는 용도
type CategoryTreeDBResult struct {
	Idx          int64                   `json:"id"`
	TenantID     string                  `json:"tenant_id"`
	CategoryId   string                  `json:"category_id"`
	Name         string                  `json:"name"`
	Desc         string                  `json:"desc"`
	ParentIdx    int64                   `json:"parent_idx"`
	UseYN        string                  `json:"use_yn"`
	PathId       string                  `json:"path_id"`
	PathName     string                  `json:"path_name"`
	Level        int64                   `json:"level"`
	Expanded     bool                    `json:"expanded"`
	ChildCount   int64                   `json:"child_count"`
	ContentCount int64                   `json:"content_count"`
	Children     []*CategoryTreeDBResult `json:"children"`
}

// Tree에서 사용하는 용도
type CategoryTreeViewResult struct {
	Name     string                    `json:"label"`
	Expanded bool                      `json:"expanded"`
	Data     CategoryDataViewResult    `json:"data"`
	Children []*CategoryTreeViewResult `json:"children"`
}

type CategoryDataViewResult struct {
	Idx          int64  `json:"id"`
	TenantID     string `json:"tenant_id"`
	CategoryId   string `json:"category_id"`
	Desc         string `json:"desc"`
	ParentIdx    int64  `json:"parent_idx"`
	UseYN        string `json:"use_yn"`
	PathId       string `json:"path_id"`
	PathName     string `json:"path_name"`
	Level        int64  `json:"level"`
	ChildCount   int64  `json:"child_count"`
	ContentCount int64  `json:"content_count"`
}

func (this *CategoryTreeViewResult) Size() int {
	var size int = len(this.Children)
	for _, c := range this.Children {
		size += c.Size()
	}
	return size
}

func (this *CategoryTreeViewResult) Add(nodes ...*CategoryTreeViewResult) bool {
	var size = this.Size()
	for _, n := range nodes {
		if n.Data.ParentIdx == this.Data.Idx {
			this.Children = append(this.Children, n)
		} else {
			for _, c := range this.Children {
				if c.Add(n) {
					break
				}
			}
		}
	}
	return this.Size() == size+len(nodes)
}

// 카테고리 리스트 조회
func (api CategoryAPI) GetCategories(c echo.Context) error {

	claims := apiModel.GetJWTClaims(c)

	categories := []CategoryTreeDBResult{}

	tx := c.Get("Tx").(*gorm.DB)

	tx.Raw(`
		SELECT  A.idx,
			A.tenant_id,
			A.category_id,
			A.path_id,
			A.path_name,
			A.name,
			A.desc,
			A.parent_idx,
			A.level,
			A.created_at,
			A.updated_at,
			A.updated_id,
			A.use_yn,
			(
				SELECT COUNT(*) FROM TB_CATEGORY
				WHERE parent_idx = A.idx
			) AS child_count,
			  (
			    CASE
				WHEN parent_idx = 0
				  THEN (
				    SELECT COUNT(*)
				    FROM TB_CONTENT
				    WHERE tenant_id = ?
				    AND deleted_at IS NULL
				  )
				  ELSE
				  (
				    SELECT COUNT(*)
				    FROM TB_CONTENT
				    WHERE tenant_id = ?
				    AND category_id = A.category_id
				    AND deleted_at IS NULL
				  )
			    END
			  ) AS content_count
			FROM
			(
				SELECT 	hi.idx,
					hi.tenant_id,
					hi.category_id,
					FN_CATEGORY_GETPATHID(hi.idx) path_id,
					FN_CATEGORY_GETPATHNAME(hi.idx) path_name,
					hi.name,
					hi.desc,
					parent_idx,
					level,
					created_at,
					updated_at,
					updated_id,
					use_yn
			FROM
			(
				SELECT FN_CATEGORY_CONNECTBY_IDX(idx) AS idx, @level AS level
					FROM (
						SELECT 	@start_with := 0,
							@idx := @start_with,
							@level := 0
						) vars, TB_CATEGORY
						WHERE @idx IS NOT NULL
						) ho
						JOIN TB_CATEGORY hi
						ON hi.idx = ho.idx
						AND hi.tenant_id = ?
			) A
			ORDER BY A.path_id
		`, claims.TenantId, claims.TenantId, claims.TenantId).Scan(&categories)

	var root *CategoryTreeViewResult = &CategoryTreeViewResult{}

	data := []*CategoryTreeViewResult{}

	for _, category := range categories {

		s := append(data, &CategoryTreeViewResult{
			Name: category.Name,
			Data: CategoryDataViewResult{
				Idx:          category.Idx,
				TenantID:     category.TenantID,
				CategoryId:   category.CategoryId,
				Desc:         category.Desc,
				ParentIdx:    category.ParentIdx,
				UseYN:        category.UseYN,
				PathId:       category.PathId,
				PathName:     category.PathName,
				Level:        category.Level,
				ChildCount:   category.ChildCount,
				ContentCount: category.ContentCount,
			},
			Expanded: true,
			Children: nil,
		})
		data = s
	}

	fmt.Println(root.Add(data...), root.Size())
	//bytes, _ := json.MarshalIndent(root, "", "\t") //formated output
	////bytes, _ := json.Marshal(root)
	//fmt.Println(string(bytes))

	return handler.APIResultHandler(c, true, http.StatusOK, root)
}

// 카테고리 한건 조회
func (api CategoryAPI) GetCategory(c echo.Context) error {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	claims := apiModel.GetJWTClaims(c)

	claims.TenantId = "livetest"

	category := &rdbModel.Category{}

	tx := c.Get("Tx").(*gorm.DB)

	if tx.Preload("Tenant").Where("idx = ? and use_yn = ?",
		idx, "Y").Find(category).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "CATEGORY NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, category)

}

// 카테고리 패스 조회
func (api CategoryAPI) GetCategoryPaths(c echo.Context) error {

	claims := apiModel.GetJWTClaims(c)

	categories := []CategoryTreeDBResult{}

	tx := c.Get("Tx").(*gorm.DB)

	tx.Raw(`
		SELECT *
			FROM
			(
				SELECT 	hi.idx,
					hi.tenant_id,
					hi.category_id,
					FN_CATEGORY_GETPATHID(hi.idx) path_id,
					FN_CATEGORY_GETPATHNAME(hi.idx) path_name,
					hi.name,
					hi.desc,
					parent_idx,
					level,
					created_at,
					updated_at,
					updated_id,
					use_yn,
					(
						SELECT COUNT(*) FROM TB_CATEGORY
						WHERE parent_idx = hi.idx
					) AS child_count
			FROM
			(
				SELECT FN_CATEGORY_CONNECTBY_IDX(idx) AS idx, @level AS level
					FROM (
						SELECT 	@start_with := 0,
							@idx := @start_with,
							@level := 0
						) vars, TB_CATEGORY
						WHERE @idx IS NOT NULL
						) ho
						JOIN TB_CATEGORY hi
						ON hi.idx = ho.idx
						AND hi.tenant_id = ?
			) A
			ORDER BY A.path_id`, claims.TenantId).Scan(&categories)

	//var root *CategoryTreeDBResult = &CategoryTreeDBResult{}

	//data := []*CategoryTreeDBResult{}

	//for _, category := range categories {
	//
	//	s := append(data, &CategoryTreeDBResult{
	//		Idx:        category.Idx,
	//		TenantID:   category.TenantID,
	//		CategoryId: category.CategoryId,
	//		Name:       category.Name,
	//		Desc:       category.Desc,
	//		ParentIdx:  category.ParentIdx,
	//		UseYN:      category.UseYN,
	//		PathId:     category.PathId,
	//		PathName:   category.PathName,
	//		Level:      category.Level,
	//		ChildCount: category.ChildCount,
	//		Children:   nil,
	//	})
	//	data = s
	//}
	//
	//fmt.Println(root.Add(data...), root.Size())
	//bytes, _ := json.MarshalIndent(root, "", "\t") //formated output
	////bytes, _ := json.Marshal(root)
	//fmt.Println(string(bytes))

	return handler.APIResultHandler(c, true, http.StatusOK, categories)
}

// 카테고리 등록
func (api CategoryAPI) PostCategory(c echo.Context) error {

	claims := apiModel.GetJWTClaims(c)

	payload := &api.request
	c.Bind(payload)

	if err := c.Validate(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	category := &rdbModel.Category{
		TenantId:   claims.TenantId,
		Name:       payload.CategoryName,
		CategoryId: uuid.NewV4().String(),
		Desc:       payload.CategoryDesc,
		UseYN:      payload.UseYN,
		ParentIdx:  payload.ParentIdx,
	}

	// 기존에 등록되어 있는 Tenant인지 체크 한다.
	tx := c.Get("Tx").(*gorm.DB)
	if !tx.Where("category_id = ? ",
		category.CategoryId).Find(category).RecordNotFound() {
		return echo.NewHTTPError(http.StatusConflict, "Already exists CATEGORY ")
	}

	tx.Create(category)

	return handler.APIResultHandler(c, true, http.StatusCreated, category)
}

// 카테고리 수정
func (api CategoryAPI) PutCategory(c echo.Context) error {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	payload := &api.request
	c.Bind(payload)

	tx := c.Get("Tx").(*gorm.DB)

	category := &rdbModel.Category{}

	if tx.Model(category).Where("idx = ? ", idx).
		Updates(
			rdbModel.Category{Name: payload.CategoryName,
				Desc:      payload.CategoryDesc,
				ParentIdx: payload.ParentIdx,
				UseYN:     payload.UseYN}).RowsAffected == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "CATEGORY NOT FOUND")
	}

	return handler.APIResultHandler(c, true, http.StatusOK, "Update OK")

}

// 카테고리 삭제
func (api CategoryAPI) DeleteCategory(c echo.Context) error {

	idx, _ := strconv.ParseInt(c.Param("idx"), 0, 64)

	category := &rdbModel.Category{}
	tx := c.Get("Tx").(*gorm.DB)

	if tx.Where("idx = ? ",
		idx).Find(category).RecordNotFound() {
		return echo.NewHTTPError(http.StatusNotFound, "CATEGORY NOT FOUND")
	}

	// 카테고리에 등록된 컨텐츠들이 존재한다면 삭제 불가...
	contents := []rdbModel.Content{}

	logrus.Debug(tx.Where("category_id = ?", category.CategoryId).Find(&contents).RecordNotFound())

	tx.Where("category_id = ?", category.CategoryId).Find(&contents)

	if len(contents) > 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "CONTENTS_EXIST")
	}

	tx.Delete(category, "idx = ?", idx)

	return handler.APIResultHandler(c, true, http.StatusOK, "Delete OK")
}
