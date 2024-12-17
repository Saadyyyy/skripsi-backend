package chat_bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type ChatBotHandler struct{}

func NewChatBotHandler() *ChatBotHandler {
	return &ChatBotHandler{}
}

func (h *ChatBotHandler) ChatSoalHandler(c echo.Context) error {
	soalID := c.Param("id")
	apiURL := fmt.Sprintf("http://localhost:8080/soal/detail?soal_id=%s", soalID)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create request"})
	}

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing token"})
	}
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch soal from API"})
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read response body"})
	}

	var soalResponse map[string]interface{}
	if err := json.Unmarshal(body, &soalResponse); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse API response"})
	}

	data, ok := soalResponse["data"].(map[string]interface{})
	if !ok || data == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid API response structure"})
	}

	soal, soalOk := data["Soal"].(string)
	jawabanA, jawabanAOk := data["JawabanA"].(string)
	jawabanB, jawabanBOk := data["JawabanB"].(string)
	jawabanC, jawabanCOk := data["JawabanC"].(string)
	jawabanD, jawabanDOk := data["JawabanD"].(string)
	jawabanE, jawabanEOk := data["JawabanE"].(string)
	jawabanBenar, jawabanBenarOk := data["JawabanBenar"].(string)

	if !soalOk || !jawabanAOk || !jawabanBOk || !jawabanCOk || !jawabanDOk || !jawabanEOk || !jawabanBenarOk {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Incomplete soal data"})
	}

	// Susun prompt dengan jawaban benar terlebih dahulu, baru penjelasan soal
	prompt := fmt.Sprintf(
		"Pilihlah Jawaban yang benar terlebih dahulu dari opsi A, B, C, D, dan E tanpa menjelaskan jawaban yang salah: %s\n\nJelaskan soal berikut tanpa menjelaskan jawaban yang salah:\n%s\n\nA. %s\nB. %s\nC. %s\nD. %s\nE. %s",
		jawabanBenar, soal, jawabanA, jawabanB, jawabanC, jawabanD, jawabanE,
	)

	chatGPTResponse, err := callOpenAI(prompt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to communicate with OpenAI API"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		// "soal":       data,
		"penjelasan": chatGPTResponse,
	})
}

func callOpenAI(prompt string) (string, error) {

	err := godotenv.Load(".env")
	if err != nil {
		return "", fmt.Errorf("Error loading .env file")
	}

	openAIURL := "https://api.openai.com/v1/chat/completions"
	openAIKey, exists := os.LookupEnv("OPENAI_API_KEY")
	if !exists || openAIKey == "" {
		return "", fmt.Errorf("OpenAI API --key is missing")
	}
	// fmt.Println("Using OpenAI API Key:", openAIKey)

	requestBody, _ := json.Marshal(map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 1.0,
		"max_tokens":  150,
	})

	req, err := http.NewRequest("POST", openAIURL, ioutil.NopCloser(bytes.NewReader(requestBody)))
	if err != nil {
		return "", fmt.Errorf("failed to create OpenAI request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openAIKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to communicate with OpenAI API: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read OpenAI response body: %v", err)
	}

	fmt.Println("Response from OpenAI:", string(body))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse OpenAI response: %v", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("unexpected OpenAI response format: %v", result)
	}

	text, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract text from choices: %v", choices)
	}

	return text, nil
}
