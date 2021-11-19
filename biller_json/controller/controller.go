package controller

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
)

var mockNumber = map[string]map[string]string{
	"inquiry": {
		"0812345678901": "success",
		"0812345678902": "failed",
		"0812345678911": "success",
		"0812345678912": "success",
		"0812345678921": "success",
		"0812345678922": "success",
		"0812345678923": "success",
	},
	"purchase": {
		"0812345678911": "success",
		"0812345678912": "faile",
		"0812345678921": "pending",
		"0812345678922": "pending",
		"0812345678923": "pending",
	},
	"advice": {
		"0812345678921": "success",
		"0812345678922": "failed",
		"0812345678923": "pednding",
	},
}

type (
	SignatureRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	BaseResponse struct {
		ErrorCode string `json:"error_code,omitempty"`
		Message   string `json:"message,omitempty"`
	}

	SignatureResponse struct {
		BaseResponse
		Signature string `json:"signature,omitempty"`
	}

	InquiryRequest struct {
		CustNo      string `json:"cust_no"`
		ProductCode string `json:"product_code"`
	}

	InquiryResponse struct {
		BaseResponse
		CustNo      string `json:"cust_no,omitempty"`
		CustName    string `json:"cust_name,omitempty"`
		ProductCode string `json:"product_code,omitempty"`
		Amount      string `json:"amount,omitempty"`
		Period      string `json:"period,omitempty"`
		RefNo       string `json:"ref_no,omitempty"`
	}

	PaymentRequest struct {
		CustNo      string `json:"cust_no"`
		ProductCode string `json:"product_code"`
		RefNo       string `json:"ref_no"`
		TrxIDReff   string `json:"trx_id_reff"`
	}

	PaymentResponse struct {
		BaseResponse
		CustNo      string `json:"cust_no,omitempty"`
		CustName    string `json:"cust_name,omitempty"`
		ProductCode string `json:"product_code,omitempty"`
		Amount      string `json:"amount,omitempty"`
		Period      string `json:"period,omitempty"`
		RefNo       string `json:"ref_no,omitempty"`
		TrxIDReff   string `json:"trx_id_reff,omitempty"`
	}

	AdviceRequest struct {
		CustNo    string `json:"cust_no"`
		RefNo     string `json:"ref_no"`
		TrxIDReff string `json:"trx_id_reff"`
	}

	AdviceResponse struct {
		BaseResponse
		CustNo      string `json:"cust_no,omitempty"`
		CustName    string `json:"cust_name,omitempty"`
		ProductCode string `json:"product_code,omitempty"`
		Amount      string `json:"amount,omitempty"`
		Period      string `json:"period,omitempty"`
		SN          string `json:"sn,omitempty"`
		TrxIDReff   string `json:"trx_id_reff,omitempty"`
	}
)

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		headerSign := strings.Split(header, " ")
		if len(headerSign) < 2 || headerSign[0] != "Bearer" {
			return c.String(http.StatusUnauthorized, "failed header")
		}

		signature, err := base64.StdEncoding.DecodeString(headerSign[1])
		if err != nil {
			return c.String(http.StatusUnauthorized, "wrong signature")
		}

		timeToCheck, err := time.Parse(time.RFC3339Nano, string(signature))
		if err != nil {
			return c.String(http.StatusUnauthorized, "wrong time")
		}

		if time.Since(timeToCheck).Minutes() >= 5 {
			return c.String(http.StatusUnauthorized, "expired")
		}

		return next(c)
	}
}

func Signature(c echo.Context) error {
	var req SignatureRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "failed read body")
	}

	if !(req.Username == "alterra" && req.Password == "rahasia") {
		return c.String(http.StatusUnauthorized, "unauthorize")
	}

	sign := time.Now().Format(time.RFC3339Nano)
	signature := base64.StdEncoding.EncodeToString([]byte(sign))

	response := SignatureResponse{
		BaseResponse: BaseResponse{
			ErrorCode: "100",
			Message:   "success",
		},
		Signature: signature,
	}

	return c.JSON(http.StatusOK, response)
}

func Inquiry(c echo.Context) error {
	var req InquiryRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "failed read body")
	}

	mockType := mockNumber["inquiry"][req.CustNo]

	var response InquiryResponse

	switch mockType {
	case "success":
		response.ErrorCode = "100"
		response.Message = "Success"
		response.Amount = "55000"
		response.CustName = "Andi"
		response.CustNo = req.CustNo
		response.Period = time.Now().Format("200601")
		response.ProductCode = req.ProductCode
		response.RefNo = "SN1234567890 inquiry"
	default:
		response.ErrorCode = "400"
		response.Message = "Failed"
	}

	return c.JSON(http.StatusOK, response)
}

func Purchase(c echo.Context) error {
	var req PaymentRequest
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "failed read body")
	}

	mockType := mockNumber["purchase"][req.CustNo]

	var response PaymentResponse

	switch mockType {
	case "success":
		response.ErrorCode = "100"
		response.Message = "Success"
		response.Amount = "55000"
		response.CustName = "Andi"
		response.CustNo = req.CustNo
		response.Period = time.Now().Format("200601")
		response.ProductCode = req.ProductCode
		response.RefNo = "SN1234567890 purchase"
		response.TrxIDReff = req.TrxIDReff
	case "pending":
		response.ErrorCode = "68"
		response.Message = "Pending"
	default:
		response.ErrorCode = "400"
		response.Message = "Failed"
	}

	return c.JSON(http.StatusOK, response)
}

func Advice(c echo.Context) error {
	var req AdviceResponse
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, "failed read body")
	}

	mockType := mockNumber["advice"][req.CustNo]

	var response AdviceResponse

	switch mockType {
	case "success":
		response.ErrorCode = "100"
		response.Message = "Success"
		response.Amount = "55000"
		response.CustName = "Andi"
		response.CustNo = req.CustNo
		response.Period = time.Now().Format("200601")
		response.ProductCode = req.ProductCode
		response.SN = "SN1234567890 advice"
		response.TrxIDReff = req.TrxIDReff
	case "pending":
		response.ErrorCode = "68"
		response.Message = "Pending"
	default:
		response.ErrorCode = "400"
		response.Message = "Failed"
	}

	return c.JSON(http.StatusOK, response)
}

func sendCallback() {

}
