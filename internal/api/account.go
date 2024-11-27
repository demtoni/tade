package api

import (
	"context"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/demtoni/tade/internal/database"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	yoocommon "github.com/rvinnie/yookassa-sdk-go/yookassa/common"
	yoopayment "github.com/rvinnie/yookassa-sdk-go/yookassa/payment"
	"golang.org/x/crypto/bcrypt"
)

var (
	validPassword = regexp.MustCompile(`^[[:print:]]{8,72}$`)
	validUsername = regexp.MustCompile(`^[[:word:]]{1,32}$`)
)

func (s *Server) setCookie(id int64, w http.ResponseWriter, r *http.Request) error {
	session, _ := s.store.Get(r, "session")
	session.Options.MaxAge = 259200
	session.Options.Secure = false
	session.Options.SameSite = http.SameSiteLaxMode
	session.Options.Domain = s.config.Domain
	session.Values["id"] = id
	return session.Save(r, w)
}

type RegistrationRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	InviteCode string `json:"invite"`
}

func (r *RegistrationRequest) Bind(_ *http.Request) error {
	if !validUsername.MatchString(r.Username) {
		return errors.New(ErrorBadUsername)
	}
	if !validPassword.MatchString(r.Password) {
		return errors.New(ErrorBadPassword)
	}
	return nil
}

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	data := &RegistrationRequest{}

	if err := render.Bind(r, data); err != nil {
		s.SendError(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}

	u, err := s.queries.GetUserByName(r.Context(), data.Username)
	if u.ID != 0 {
		s.SendError(w, r, nil, http.StatusBadRequest, ErrorUsernameTaken)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
	if err != nil {
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	inv, err := s.queries.GetInvite(r.Context(), data.InviteCode)
	if err != nil || inv.Used > 0 {
		s.SendError(w, r, nil, http.StatusBadRequest, ErrorBadInvite)
		return
	}

	if err := s.queries.UseInvite(r.Context(), inv.ID); err != nil {
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	id, err := s.queries.CreateUser(r.Context(), database.CreateUserParams{
		Name:         data.Username,
		PasswordHash: string(hash),
		Balance:      int64(0),
		Invites:      int64(0),
	})
	if err != nil {
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	if s.setCookie(id, w, r); err != nil {
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	render.Status(r, http.StatusCreated)
}

func (s *Server) AuthCtx(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.store.Get(r, "session")

		id, ok := session.Values["id"]
		if !ok {
			s.SendError(w, r, nil, http.StatusUnauthorized, ErrorUnauthorized)
			return
		}

		u, err := s.queries.GetUser(r.Context(), id.(int64))
		if err != nil {
			s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
			return
		}

		ctx := context.WithValue(r.Context(), "user", &u)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *LoginRequest) Bind(_ *http.Request) error {
	if r.Username == "" || r.Password == "" {
		return errors.New(ErrorEmptyField)
	}
	return nil
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	data := &LoginRequest{}

	if err := render.Bind(r, data); err != nil {
		s.SendError(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}

	u, err := s.queries.GetUserByName(r.Context(), data.Username)
	if err != nil {
		s.SendError(w, r, err, http.StatusNotFound, ErrorUserNotFound)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(data.Password)); err != nil {
		s.SendError(w, r, err, http.StatusNotFound, ErrorUserNotFound)
		return
	}

	if err := s.setCookie(u.ID, w, r); err != nil {
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	render.Status(r, http.StatusOK)
}

type PasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (r *PasswordRequest) Bind(_ *http.Request) error {
	if !validPassword.MatchString(r.NewPassword) {
		return errors.New(ErrorBadPassword)
	}
	return nil
}

func (s *Server) ChangePassword(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user").(*database.User)

	data := &PasswordRequest{}
	if err := render.Bind(r, data); err != nil {
		s.SendError(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(data.OldPassword)); err != nil {
		s.SendError(w, r, nil, http.StatusBadRequest, ErrorPasswordsDontMatch)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), 10)
	if err != nil {
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	err = s.queries.UpdatePassword(r.Context(), database.UpdatePasswordParams{
		ID:           u.ID,
		PasswordHash: string(hash),
	})
	if err != nil {
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	render.Status(r, http.StatusOK)
}

const (
	TransactionInProcess = "in_process"
	TransactionCanceled  = "canceled"
	TransactionCompleted = "completed"
)

type BalanceRequest struct {
	Amount        int    `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	ReturnURL     string `json:"return_url"`
}

type BalanceResponse struct {
	PaymentUrl string `json:"payment_url"`
}

func (r *BalanceResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (r *BalanceRequest) Bind(_ *http.Request) error {
	if r.Amount < 0 {
		return errors.New(ErrorBadAmount)
	}
	if r.PaymentMethod == "" || r.ReturnURL == "" {
		return errors.New(ErrorEmptyField)
	}
	return nil
}

func (s *Server) AddBalance(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user").(*database.User)

	data := &BalanceRequest{}
	if err := render.Bind(r, data); err != nil {
		s.SendError(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}

	payment, err := s.kassa.CreatePayment(&yoopayment.Payment{
		Amount: &yoocommon.Amount{
			Value:    strconv.Itoa(data.Amount),
			Currency: "RUB",
		},
		PaymentMethod: yoopayment.PaymentMethodType(data.PaymentMethod),
		Confirmation: yoopayment.Redirect{
			Type:      "redirect",
			ReturnURL: data.ReturnURL,
		},
	})
	if err != nil {
		s.SendError(w, r, nil, http.StatusBadRequest, err.Error())
		return
	}

	paymentUrl, err := s.kassa.ParsePaymentLink(payment)
	if err != nil {
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}
	_, err = s.queries.CreateTransaction(r.Context(), database.CreateTransactionParams{
		PaymentID: payment.ID,
		Amount:    int64(data.Amount),
		Status:    TransactionInProcess,
		Timestamp: payment.CreatedAt.Unix(),
		Url:       paymentUrl,
		UserID:    u.ID,
	})
	if err != nil {
		// TODO: cancel payment if for some reason we can't save transaction to the db
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	render.Render(w, r, &BalanceResponse{paymentUrl})
}

const transactionTimeLimit int64 = 3600

func (s *Server) ProcessPayments() error {
	_, err := s.queries.CancelExpiredTransactions(context.TODO(), time.Now().Unix()-transactionTimeLimit)
	if err != nil {
		return err
	}

	var cursor string
	for {
		captureBatch, err := s.kassa.FindPayments(&yoopayment.PaymentListFilter{
			Cursor: cursor,
			Status: yoopayment.WaitingForCapture,
		})
		if err != nil {
			return err
		}

		for _, payment := range captureBatch.Items {
			p, err := s.kassa.CapturePayment(&payment)
			if err != nil {
				return err
			}
			meta, err := s.queries.UpdateTransaction(context.TODO(), database.UpdateTransactionParams{
				Status:    TransactionCompleted,
				PaymentID: p.ID,
			})
			if err != nil {
				return err
			}
			u, err := s.queries.GetUser(context.TODO(), meta.UserID)
			if err != nil {
				return err
			}
			if err := s.queries.UpdateBalance(context.TODO(), database.UpdateBalanceParams{
				Balance: u.Balance + meta.Amount,
				ID:      u.ID,
			}); err != nil {
				return err
			}
		}
		if cursor == "" {
			break
		}
	}
	return nil
}

type TransactionResponse struct {
	URL       string `json:"url"`
	Amount    int64  `json:"amount"`
	Timestamp int64  `json:"timestamp"`
	Status    string `json:"status"`
}

func (r *TransactionResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewTransactionListResponse(transactions *[]database.Transaction) []render.Renderer {
	list := []render.Renderer{}
	for _, trans := range *transactions {
		list = append(list, &TransactionResponse{
			trans.Url, trans.Amount, trans.Timestamp, trans.Status,
		})
	}
	return list
}

func (s *Server) GetTransactionList(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user").(*database.User)

	transactions, err := s.queries.ListTransactions(r.Context(), u.ID)
	if err != nil {
		s.SendError(w, r, nil, http.StatusNotFound, ErrorNoTransactions)
		return
	}

	render.RenderList(w, r, NewTransactionListResponse(&transactions))
}

type UserInfoResponse struct {
	Username string `json:"username"`
	Balance  int64  `json:"balance"`
	Invites  int64  `json:"invites"`
}

func (r *UserInfoResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (s *Server) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user").(*database.User)
	render.Render(w, r, &UserInfoResponse{u.Name, u.Balance, u.Invites})
}

type InviteResponse struct {
	Code string `json:"invite_code"`
}

func (r *InviteResponse) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func NewInviteListResponse(invites *[]database.Invite) []render.Renderer {
	list := []render.Renderer{}
	for _, invite := range *invites {
		list = append(list, &InviteResponse{invite.Code})
	}
	return list
}

func (s *Server) ListInvites(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user").(*database.User)

	invites, err := s.queries.GetUserInvites(r.Context(), u.ID)
	if err != nil {
		s.SendError(w, r, nil, http.StatusNotFound, ErrorNoUnusedInvites)
		return
	}

	render.RenderList(w, r, NewInviteListResponse(&invites))
}

func (s *Server) GenerateInvite(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user").(*database.User)

	if u.Invites == 0 {
		s.SendError(w, r, nil, http.StatusForbidden, ErrorNoInvites)
		return
	}

	inviteCode := uuid.New().String()
	if err := s.queries.CreateInvite(r.Context(), database.CreateInviteParams{
		Code:   inviteCode,
		Used:   int64(0),
		UserID: u.ID,
	}); err != nil {
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	if err := s.queries.UpdateUserInvites(r.Context(), u.ID); err != nil {
		s.SendError(w, r, err, http.StatusInternalServerError, ErrorInternal)
		return
	}

	render.Status(r, http.StatusCreated)
}
