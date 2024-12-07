package telegram

import (
	"testing"
)

func TestDecodeTelegramInitData(t *testing.T) {
	t.Run("success decode telegram init data", func(t *testing.T) {
		telegramInitDataText := "query_id=AAF42zIVAAAAAHjbMhWPTccE&user=%7B%22id%22%3A355654520%2C%22first_name%22%3A%22Sergey%22%2C%22last_name%22%3A%22Konar%22%2C%22username%22%3A%22sergey_konar%22%2C%22language_code%22%3A%22ru%22%2C%22is_premium%22%3Atrue%2C%22allows_write_to_pm%22%3Atrue%2C%22photo_url%22%3A%22https%3A%5C%2F%5C%2Ft.me%5C%2Fi%5C%2Fuserpic%5C%2F320%5C%2FUKloFp3wkmk8Uz2Z74Z6lufjlDKJIjl7eHuMjaRzCpI.svg%22%7D&auth_date=1732181719&signature=fzsMe4hquuUM85C9YukEGStFeTJKJkBYe3caJhimZExWDOKjuivsB-0rPcEHs_6lPIlATa7DLUqtM0qIeOmtDQ&hash=fb2ba5f31223c7829fe5d5de4faa6d7614d01866a300c16c887eeb6843f4453f"
		initData := InitData{
			Token: "11",
		}
		_, err := initData.Decode(telegramInitDataText)
		if err != nil {
			t.Error("fail to decode", err)
		}
	})
}
