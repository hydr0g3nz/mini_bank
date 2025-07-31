package dto

import "time"

// VerifyRequest represents the input data for verifying a top-up request
type VerifyRequest struct {
	UserID        uint    `json:"user_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
}

// VerifyResponse represents the output data for verifying a top-up request
type VerifyResponse struct {
	TransactionID uint      `json:"transaction_id"`
	UserID        uint      `json:"user_id"`
	Amount        float64   `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	Status        string    `json:"status"`
	ExpiresAt     time.Time `json:"expires_at"`
}

// ConfirmRequest represents the input data for confirming a top-up transaction
type ConfirmRequest struct {
	TransactionID uint `json:"transaction_id"`
}

// ConfirmResponse represents the output data for confirming a top-up transaction
type ConfirmResponse struct {
	TransactionID uint    `json:"transaction_id"`
	UserID        uint    `json:"user_id"`
	Amount        float64 `json:"amount"`
	Status        string  `json:"status"`
	Balance       float64 `json:"balance"`
}

// RegisterRequest represents user registration request DTO
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=candidate company_hr admin"`
}

// LoginRequest represents user login request DTO
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateProfileRequest represents user profile update request DTO
type UpdateProfileRequest struct {
	Email string `json:"email,omitempty" validate:"omitempty,email"`
}

// ChangePasswordRequest represents password change request DTO
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID            int        `json:"id"`
	Email         string     `json:"email"`
	Role          string     `json:"role"`
	IsActive      bool       `json:"is_active"`
	EmailVerified bool       `json:"email_verified"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
}

// LoginResponse represents login response with token
type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

// UsersListResponse represents paginated users list response
type UsersListResponse struct {
	Users  []*UserResponse `json:"users"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

// CreateCategoryRequest represents category creation request DTO
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

// UpdateCategoryRequest represents category update request DTO
type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
}

// CategoryResponse represents category data in responses
type CategoryResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// CategoriesListResponse represents categories list response
type CategoriesListResponse struct {
	Categories []*CategoryResponse `json:"categories"`
	Total      int                 `json:"total"`
}

// CreateMenuItemRequest represents menu item creation request DTO
type CreateMenuItemRequest struct {
	CategoryID  int     `json:"category_id" validate:"required,gt=0"`
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price" validate:"required,gte=0"`
}

// UpdateMenuItemRequest represents menu item update request DTO
type UpdateMenuItemRequest struct {
	CategoryID  int     `json:"category_id" validate:"required,gt=0"`
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price" validate:"required,gte=0"`
}

// MenuItemResponse represents menu item data in responses
type MenuItemResponse struct {
	ID          int               `json:"id"`
	CategoryID  int               `json:"category_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	CreatedAt   time.Time         `json:"created_at"`
	Category    *CategoryResponse `json:"category,omitempty"`
}

// MenuItemListResponse represents menu items list response
type MenuItemListResponse struct {
	Items  []*MenuItemResponse `json:"items"`
	Total  int                 `json:"total"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
}

// Order DTOs
type CreateOrderRequest struct {
	TableID int `json:"table_id" validate:"required,gt=0"`
}

type UpdateOrderRequest struct {
	Status string `json:"status" validate:"required,oneof=open closed"`
}

type OrderResponse struct {
	ID        int            `json:"id"`
	TableID   int            `json:"table_id"`
	Status    string         `json:"status"`
	Notes     string         `json:"notes,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	ClosedAt  *time.Time     `json:"closed_at,omitempty"`
	Table     *TableResponse `json:"table,omitempty"`
}

type OrderWithItemsResponse struct {
	ID        int                  `json:"id"`
	TableID   int                  `json:"table_id"`
	Status    string               `json:"status"`
	Notes     string               `json:"notes,omitempty"`
	Items     []*OrderItemResponse `json:"items"`
	Total     float64              `json:"total"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	ClosedAt  *time.Time           `json:"closed_at,omitempty"`
	Table     *TableResponse       `json:"table,omitempty"`
}

type OrderListResponse struct {
	Orders []*OrderResponse `json:"orders"`
	Total  int              `json:"total"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
}

// Order Item DTOs
type AddOrderItemRequest struct {
	OrderID  int    `json:"order_id" validate:"required,gt=0"`
	ItemID   int    `json:"item_id" validate:"required,gt=0"`
	Quantity int    `json:"quantity" validate:"required,gt=0"`
	Notes    string `json:"notes,omitempty"`
}

type UpdateOrderItemRequest struct {
	Quantity int    `json:"quantity" validate:"required,gt=0"`
	Notes    string `json:"notes,omitempty"`
}

type OrderItemResponse struct {
	ID        int               `json:"id"`
	OrderID   int               `json:"order_id"`
	ItemID    int               `json:"item_id"`
	Quantity  int               `json:"quantity"`
	UnitPrice float64           `json:"unit_price"`
	Subtotal  float64           `json:"subtotal"`
	Notes     string            `json:"notes,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	MenuItem  *MenuItemResponse `json:"menu_item,omitempty"`
}

type OrderTotalResponse struct {
	OrderID   int                  `json:"order_id"`
	Items     []*OrderItemResponse `json:"items"`
	Total     float64              `json:"total"`
	ItemCount int                  `json:"item_count"`
}

// Table DTOs
type CreateTableRequest struct {
	TableNumber int `json:"table_number" validate:"required,gt=0"`
	Seating     int `json:"seating" validate:"gte=0"`
}

type UpdateTableRequest struct {
	TableNumber int `json:"table_number" validate:"required,gt=0"`
	Seating     int `json:"seating" validate:"gte=0"`
}

type TableResponse struct {
	ID          int       `json:"id"`
	TableNumber int       `json:"table_number"`
	QRCode      string    `json:"qr_code"`
	Seating     int       `json:"seating"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TableListResponse struct {
	Tables []*TableResponse `json:"tables"`
	Total  int              `json:"total"`
}

// QR Code DTOs
type QRCodeScanResponse struct {
	TableID      int            `json:"table_id"`
	Table        *TableResponse `json:"table"`
	HasOpenOrder bool           `json:"has_open_order"`
	OpenOrder    *OrderResponse `json:"open_order,omitempty"`
}

// Payment DTOs
type ProcessPaymentRequest struct {
	OrderID int     `json:"order_id" validate:"required,gt=0"`
	Amount  float64 `json:"amount" validate:"required,gt=0"`
	Method  string  `json:"method" validate:"required,oneof=cash credit_card wallet"`
}

type PaymentResponse struct {
	ID        int            `json:"id"`
	OrderID   int            `json:"order_id"`
	Amount    float64        `json:"amount"`
	Method    string         `json:"method"`
	Reference string         `json:"reference,omitempty"`
	PaidAt    time.Time      `json:"paid_at"`
	Order     *OrderResponse `json:"order,omitempty"`
}

type PaymentListResponse struct {
	Payments []*PaymentResponse `json:"payments"`
	Total    int                `json:"total"`
	Limit    int                `json:"limit"`
	Offset   int                `json:"offset"`
}
