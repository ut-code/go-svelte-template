package handler

import (
	"encoding/json"
	"net/http"

	"web-go-template/internal/db"
)

type Handler struct {
	queries *db.Queries
}

func New(queries *db.Queries) *Handler {
	return &Handler{queries: queries}
}

func (h *Handler) todayWord(w http.ResponseWriter, r *http.Request) (db.Word, bool) {
	word, err := h.queries.GetTodayWord(r.Context())
	if err != nil {
		http.Error(w, "no word set for today", http.StatusNotFound)
		return word, false
	}
	return word, true
}

func (h *Handler) GetTodayWord(w http.ResponseWriter, r *http.Request) {
	word, ok := h.todayWord(w, r)
	if !ok {
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":     word.ID,
		"length": len(word.Word),
		"date":   word.Date.Time.Format("2006-01-02"),
	})
}

type GuessRequest struct {
	Guess string `json:"guess"`
}

func (h *Handler) GuessWord(w http.ResponseWriter, r *http.Request) {
	word, ok := h.todayWord(w, r)
	if !ok {
		return
	}

	var req GuessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Guess) != len(word.Word) {
		http.Error(w, "guess must be the same length as the word", http.StatusBadRequest)
		return
	}

	result := evaluateGuess(word.Word, req.Guess)
	correct := req.Guess == word.Word

	writeJSON(w, http.StatusOK, map[string]any{
		"result":  result,
		"correct": correct,
	})
}

type SubmitScoreRequest struct {
	PlayerName string `json:"player_name"`
	Attempts   int32  `json:"attempts"`
	Solved     bool   `json:"solved"`
}

func (h *Handler) SubmitScore(w http.ResponseWriter, r *http.Request) {
	word, ok := h.todayWord(w, r)
	if !ok {
		return
	}

	var req SubmitScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	score, err := h.queries.CreateScore(r.Context(), db.CreateScoreParams{
		PlayerName: req.PlayerName,
		WordID:     word.ID,
		Attempts:   req.Attempts,
		Solved:     req.Solved,
	})
	if err != nil {
		http.Error(w, "failed to save score", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, score)
}

func (h *Handler) GetScores(w http.ResponseWriter, r *http.Request) {
	word, ok := h.todayWord(w, r)
	if !ok {
		return
	}

	scores, err := h.queries.GetScoresByWord(r.Context(), word.ID)
	if err != nil {
		http.Error(w, "failed to get scores", http.StatusInternalServerError)
		return
	}

	if scores == nil {
		scores = []db.Score{}
	}

	writeJSON(w, http.StatusOK, scores)
}

// evaluateGuess returns a slice of hints for each character.
// "correct" = right letter, right position
// "present" = right letter, wrong position
// "absent"  = letter not in word
func evaluateGuess(answer, guess string) []string {
	answerRunes := []rune(answer)
	guessRunes := []rune(guess)
	n := len(answerRunes)
	result := make([]string, n)
	used := make([]bool, n)

	for i := 0; i < n; i++ {
		if guessRunes[i] == answerRunes[i] {
			result[i] = "correct"
			used[i] = true
		}
	}

	for i := 0; i < n; i++ {
		if result[i] == "correct" {
			continue
		}
		found := false
		for j := 0; j < n; j++ {
			if !used[j] && guessRunes[i] == answerRunes[j] {
				result[i] = "present"
				used[j] = true
				found = true
				break
			}
		}
		if !found {
			result[i] = "absent"
		}
	}

	return result
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
