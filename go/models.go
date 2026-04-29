package tesote

import "encoding/json"

// Account represents a bank account in v1 and v2 (identical schema).
type Account struct {
	ID              string       `json:"id"`
	Name            string       `json:"name"`
	Data            AccountData  `json:"data"`
	Bank            AccountBank  `json:"bank"`
	LegalEntity     AccountLegal `json:"legal_entity"`
	TesoteCreatedAt string       `json:"tesote_created_at"`
	TesoteUpdatedAt string       `json:"tesote_updated_at"`
}

// AccountData holds the per-account metadata block.
type AccountData struct {
	MaskedAccountNumber          string  `json:"masked_account_number"`
	Currency                     string  `json:"currency"`
	TransactionsDataCurrentAsOf  *string `json:"transactions_data_current_as_of"`
	BalanceDataCurrentAsOf       *string `json:"balance_data_current_as_of"`
	CustomUserProvidedIdentifier *string `json:"custom_user_provided_identifier"`
	BalanceCents                 *string `json:"balance_cents,omitempty"`
	AvailableBalanceCents        *string `json:"available_balance_cents,omitempty"`
}

// AccountBank is the embedded bank reference.
type AccountBank struct {
	Name string `json:"name"`
}

// AccountLegal is the embedded legal-entity reference.
type AccountLegal struct {
	ID        *string `json:"id"`
	LegalName *string `json:"legal_name"`
}

// PagePagination is the page-based pagination block (v1 accounts, v2 accounts).
type PagePagination struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
	TotalCount  int `json:"total_count"`
}

// CursorPagination is the cursor-based pagination block (transactions index).
type CursorPagination struct {
	HasMore  bool   `json:"has_more"`
	PerPage  int    `json:"per_page"`
	AfterID  string `json:"after_id"`
	BeforeID string `json:"before_id"`
}

// AccountListResponse is the wire shape for GET /accounts.
type AccountListResponse struct {
	Total      int            `json:"total"`
	Accounts   []Account      `json:"accounts"`
	Pagination PagePagination `json:"pagination"`
}

// Transaction is the v1 transaction schema (also returned by GET /v2/transactions/{id}).
type Transaction struct {
	ID                    string                `json:"id"`
	Status                string                `json:"status"`
	Data                  TransactionData       `json:"data"`
	TesoteImportedAt      string                `json:"tesote_imported_at"`
	TesoteUpdatedAt       string                `json:"tesote_updated_at"`
	TransactionCategories []TransactionCategory `json:"transaction_categories"`
	Counterparty          *Counterparty         `json:"counterparty"`
}

// TransactionData is the inner data block for Transaction.
type TransactionData struct {
	AmountCents         int64   `json:"amount_cents"`
	Currency            string  `json:"currency"`
	Description         string  `json:"description"`
	TransactionDate     string  `json:"transaction_date"`
	CreatedAt           *string `json:"created_at"`
	CreatedAtDate       *string `json:"created_at_date"`
	Note                *string `json:"note"`
	ExternalServiceID   *string `json:"external_service_id"`
	RunningBalanceCents *int64  `json:"running_balance_cents,omitempty"`
}

// TransactionCategory tags a transaction with a category.
type TransactionCategory struct {
	Name                 string  `json:"name"`
	ExternalCategoryCode *string `json:"external_category_code"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}

// Counterparty is the embedded counterparty reference on a transaction.
type Counterparty struct {
	Name string `json:"name"`
}

// TransactionListResponse is the wire shape for GET .../transactions.
type TransactionListResponse struct {
	Total        int              `json:"total"`
	Transactions []Transaction    `json:"transactions"`
	Pagination   CursorPagination `json:"pagination"`
}

// SyncTransaction is the flattened, Plaid-compatible v2 sync entry.
type SyncTransaction struct {
	TransactionID          string   `json:"transaction_id"`
	AccountID              string   `json:"account_id"`
	Amount                 float64  `json:"amount"`
	ISOCurrencyCode        string   `json:"iso_currency_code"`
	UnofficialCurrencyCode string   `json:"unofficial_currency_code"`
	Date                   string   `json:"date"`
	Datetime               *string  `json:"datetime"`
	Name                   string   `json:"name"`
	MerchantName           *string  `json:"merchant_name"`
	Pending                bool     `json:"pending"`
	Category               []string `json:"category"`
	RunningBalanceCents    *int64   `json:"running_balance_cents,omitempty"`
}

// RemovedTransaction is the entry shape in the sync response's "removed" array.
type RemovedTransaction struct {
	TransactionID string `json:"transaction_id"`
	AccountID     string `json:"account_id"`
}

// TransactionSyncResponse is the wire shape for POST .../transactions/sync.
type TransactionSyncResponse struct {
	Added      []SyncTransaction    `json:"added"`
	Modified   []SyncTransaction    `json:"modified"`
	Removed    []RemovedTransaction `json:"removed"`
	NextCursor *string              `json:"next_cursor"`
	HasMore    bool                 `json:"has_more"`
}

// TransactionOrder is the v2 transaction order resource.
type TransactionOrder struct {
	ID                string                  `json:"id"`
	Status            string                  `json:"status"`
	Amount            float64                 `json:"amount"`
	Currency          string                  `json:"currency"`
	Description       string                  `json:"description"`
	Reference         *string                 `json:"reference"`
	ExternalReference *string                 `json:"external_reference"`
	IdempotencyKey    *string                 `json:"idempotency_key"`
	BatchID           *string                 `json:"batch_id"`
	ScheduledFor      *string                 `json:"scheduled_for"`
	ApprovedAt        *string                 `json:"approved_at"`
	SubmittedAt       *string                 `json:"submitted_at"`
	CompletedAt       *string                 `json:"completed_at"`
	FailedAt          *string                 `json:"failed_at"`
	CancelledAt       *string                 `json:"cancelled_at"`
	SourceAccount     OrderSourceAccount      `json:"source_account"`
	Destination       OrderDestination        `json:"destination"`
	Fee               *OrderFee               `json:"fee"`
	ExecutionStrategy *string                 `json:"execution_strategy"`
	TesoteTransaction *OrderTesoteTransaction `json:"tesote_transaction"`
	LatestAttempt     *OrderLatestAttempt     `json:"latest_attempt"`
	CreatedAt         string                  `json:"created_at"`
	UpdatedAt         string                  `json:"updated_at"`
}

// OrderSourceAccount is the source-account reference on a TransactionOrder.
type OrderSourceAccount struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	PaymentMethodID string `json:"payment_method_id"`
}

// OrderDestination is the destination reference on a TransactionOrder.
type OrderDestination struct {
	PaymentMethodID  string `json:"payment_method_id"`
	CounterpartyID   string `json:"counterparty_id"`
	CounterpartyName string `json:"counterparty_name"`
}

// OrderFee is the optional fee block on a TransactionOrder.
type OrderFee struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// OrderTesoteTransaction links the bank-side transaction once executed.
type OrderTesoteTransaction struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// OrderLatestAttempt is the most-recent execution attempt summary.
type OrderLatestAttempt struct {
	ID                string  `json:"id"`
	Status            string  `json:"status"`
	AttemptNumber     int     `json:"attempt_number"`
	ExternalReference *string `json:"external_reference"`
	SubmittedAt       *string `json:"submitted_at"`
	CompletedAt       *string `json:"completed_at"`
	ErrorCode         *string `json:"error_code"`
	ErrorMessage      *string `json:"error_message"`
}

// OffsetEnvelope is the standard offset-based list envelope (sync_sessions, orders, methods).
type OffsetEnvelope struct {
	HasMore bool `json:"has_more"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
}

// TransactionOrderListResponse wraps GET .../transaction_orders.
type TransactionOrderListResponse struct {
	Items   []TransactionOrder `json:"items"`
	HasMore bool               `json:"has_more"`
	Limit   int                `json:"limit"`
	Offset  int                `json:"offset"`
}

// PaymentMethod is the v2 payment method resource.
type PaymentMethod struct {
	ID            string                     `json:"id"`
	MethodType    string                     `json:"method_type"`
	Currency      string                     `json:"currency"`
	Label         *string                    `json:"label"`
	Details       map[string]any             `json:"details"`
	Verified      bool                       `json:"verified"`
	VerifiedAt    *string                    `json:"verified_at"`
	LastUsedAt    *string                    `json:"last_used_at"`
	Counterparty  *PaymentMethodCounterparty `json:"counterparty"`
	TesoteAccount *PaymentMethodAccount      `json:"tesote_account"`
	CreatedAt     string                     `json:"created_at"`
	UpdatedAt     string                     `json:"updated_at"`
}

// PaymentMethodCounterparty is the counterparty link on a PaymentMethod.
type PaymentMethodCounterparty struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PaymentMethodAccount is the tesote_account link on a PaymentMethod.
type PaymentMethodAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PaymentMethodListResponse wraps GET /v2/payment_methods.
type PaymentMethodListResponse struct {
	Items   []PaymentMethod `json:"items"`
	HasMore bool            `json:"has_more"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
}

// SyncSession is the v2 sync session resource.
type SyncSession struct {
	ID                 string            `json:"id"`
	Status             string            `json:"status"`
	StartedAt          string            `json:"started_at"`
	CompletedAt        *string           `json:"completed_at"`
	TransactionsSynced int               `json:"transactions_synced"`
	AccountsCount      int               `json:"accounts_count"`
	Error              *SyncSessionError `json:"error"`
	Performance        *SyncPerformance  `json:"performance"`
}

// SyncSessionError describes a failed sync session.
type SyncSessionError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// SyncPerformance contains optional performance metrics.
type SyncPerformance struct {
	TotalDuration   float64 `json:"total_duration"`
	ComplexityScore float64 `json:"complexity_score"`
	SyncSpeedScore  float64 `json:"sync_speed_score"`
}

// SyncSessionListResponse wraps GET /v2/accounts/{id}/sync_sessions.
type SyncSessionListResponse struct {
	SyncSessions []SyncSession `json:"sync_sessions"`
	Limit        int           `json:"limit"`
	Offset       int           `json:"offset"`
	HasMore      bool          `json:"has_more"`
}

// AccountSyncResponse is the 202 response from POST /v2/accounts/{id}/sync.
type AccountSyncResponse struct {
	Message       string `json:"message"`
	SyncSessionID string `json:"sync_session_id"`
	Status        string `json:"status"`
	StartedAt     string `json:"started_at"`
}

// BulkResultEntry is one account's slice of the bulk transactions response.
type BulkResultEntry struct {
	AccountID    string           `json:"account_id"`
	Transactions []Transaction    `json:"transactions"`
	Pagination   CursorPagination `json:"pagination"`
}

// BulkTransactionsResponse wraps POST /v2/transactions/bulk.
type BulkTransactionsResponse struct {
	BulkResults []BulkResultEntry `json:"bulk_results"`
}

// SearchTransactionsResponse wraps GET /v2/transactions/search.
type SearchTransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
	Total        int           `json:"total"`
}

// BatchSummary is the response for GET /v2/accounts/{id}/batches/{batch_id}.
type BatchSummary struct {
	BatchID          string             `json:"batch_id"`
	TotalOrders      int                `json:"total_orders"`
	TotalAmountCents int64              `json:"total_amount_cents"`
	AmountCurrency   string             `json:"amount_currency"`
	Statuses         map[string]int     `json:"statuses"`
	BatchStatus      string             `json:"batch_status"`
	CreatedAt        string             `json:"created_at"`
	Orders           []TransactionOrder `json:"orders"`
}

// BatchCreateResponse wraps POST /v2/accounts/{id}/batches.
type BatchCreateResponse struct {
	BatchID string             `json:"batch_id"`
	Orders  []TransactionOrder `json:"orders"`
	Errors  []json.RawMessage  `json:"errors"`
}

// BatchApproveResponse wraps POST /v2/accounts/{id}/batches/{id}/approve.
type BatchApproveResponse struct {
	Approved int `json:"approved"`
	Failed   int `json:"failed"`
}

// BatchSubmitResponse wraps POST /v2/accounts/{id}/batches/{id}/submit.
type BatchSubmitResponse struct {
	Enqueued int `json:"enqueued"`
	Failed   int `json:"failed"`
}

// BatchCancelResponse wraps POST /v2/accounts/{id}/batches/{id}/cancel.
type BatchCancelResponse struct {
	Cancelled int               `json:"cancelled"`
	Skipped   int               `json:"skipped"`
	Errors    []json.RawMessage `json:"errors"`
}

// StatusResponse is the response from GET /status and GET /v2/status.
type StatusResponse struct {
	Status        string `json:"status"`
	Authenticated bool   `json:"authenticated"`
}

// WhoamiResponse is the response from GET /whoami and GET /v2/whoami.
type WhoamiResponse struct {
	Client WhoamiClient `json:"client"`
}

// WhoamiClient identifies the API key's owner.
type WhoamiClient struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
