package handler

import (
	"encoding/json" // ★ 変更点: JSONを扱うために追加
	"log"           // ★ 変更点: エラーログを出力するために追加
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ★ 変更点: ここから下を修正

	// 1. レスポンス用の構造体を作成し、Messageに "OK" をセット
	resp := model.HealthzResponse{Message: "OK"}

	// 2. レスポンスをJSON形式に変換してブラウザに送信
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// 3. もしエラーが発生したら、ターミナルにログを出力
		log.Println("ERROR:", err)
	}
}
