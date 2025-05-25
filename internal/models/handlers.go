package models

// @Description Basic book information response
// @Example {"id":1,"title":"Война и мир","author":"Л. Н. Толстой","price":1300}
type BookResponse struct {
	ID     uint   `json:"id"     example:"1"`
	Title  string `json:"title"  example:"Война и мир"`
	Author string `json:"author" example:"Л. Н. Толстой"`
	Price  uint   `json:"price"  example:"1300"`
}

// @Description User profile with books list
// @Example {"username":"Wladim1r","total_books":2,"books":[{"id":1,"title":"Война и мир","author":"Л. Н. Толстой","price":1300}]}
type UserBooksResponse struct {
	Username   string         `json:"username"    example:"Wladim1r"`
	TotalBooks int            `json:"total_books" example:"2"`
	Books      []BookResponse `json:"books"`
}

// @Description List of users with their books
// @Example {"data":[{"username":"Wladim1r","total_books":1,"books":[{"id":1,"title":"Война и мир","author":"Л. Н. Толстой","price":1300}]}]}
type UsersBooksResponse struct {
	Data []UserBooksResponse `json:"data"`
}

// @Description Book creation/update request
// @Example {"title":"Война и мир","author":"Л. Н. Толстой","price":1300}
type BookRequest struct {
	Title  string `json:"title"  binding:"required" example:"Война и мир"`
	Author string `json:"author" binding:"required" example:"Л. Н. Толстой"`
	Price  uint   `json:"price"  binding:"required" example:"1300"`
}

// @Description Detailed book response
// @Example {"data":{"id":1,"title":"Война и мир","author":"Л. Н. Толстой","price":1300,"user_id":1}}
type GetBook struct {
	Data Book `json:"data"`
}

// @Description Books response metadata
// @Example {"total":5,"user_id":1}
type MetaBook struct {
	Total  int  `json:"total"   example:"5"`
	UserID uint `json:"user_id" example:"1"`
}

// @Description Paginated books response
// @Example {"data":[{"id":1,"title":"Война и мир","author":"Л. Н. Толстой","price":1300,"user_id":1}],"meta":{"total":1,"user_id":1}}
type GetBooks struct {
	Data []Book   `json:"data"`
	Meta MetaBook `json:"meta"`
}

// @Description Default error response
// @Example {"error":"User not found"}
type ErrorResponse struct {
	Error string `json:"error" example:"error description"`
}

// @Description Default successfully response
// @Example {"message":"User created successfully"}
type SuccessResponse struct {
	Message string `json:"message" example:"message description"`
}
