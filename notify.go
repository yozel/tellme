package tellme

import (
	"fmt"
	"net/http"
	"net/url"
)

func SendNotification(telegramToken string, chatID int64, message string) error {
	fmt.Printf("Sending notification to %d: %s", chatID, message)
	// // return nil
	url := url.URL{
		Scheme: "https",
		Host:   "api.telegram.org",
		Path:   fmt.Sprintf("bot%s/sendMessage", telegramToken),
		RawQuery: url.Values{
			"chat_id":    []string{fmt.Sprintf("%d", chatID)},
			"text":       []string{message},
			"parse_mode": []string{"MarkdownV2"},
		}.Encode(),
	}
	res, err := http.Get(url.String())
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api returned status code %d", res.StatusCode)
	}
	return nil
}
