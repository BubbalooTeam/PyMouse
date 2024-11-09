package telegram

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
)

const (
	SYMBOL_FIRST_PAGE    = "« %d"
	SYMBOL_PREVIOUS_PAGE = "‹ %d"
	SYMBOL_CURRENT_PAGE  = "· %d ·"
	SYMBOL_NEXT_PAGE     = "%d ›"
	SYMBOL_LAST_PAGE     = "%d »"
)

func FormatCallbackPattern(CallbackPattern string, CurrentPage int) string {
	if strings.Contains(CallbackPattern, "{number}") {
		return strings.Replace(CallbackPattern, "{number}", strconv.Itoa(CurrentPage), -1)
	}
	return CallbackPattern
}

func AddButton(Text string, CallbackPattern string) telego.InlineKeyboardButton {
	return telego.InlineKeyboardButton{
		Text:         Text,
		CallbackData: CallbackPattern,
	}
}

func KeyboardPaginate(TotalPages int, CurrentPage int, CallbackPattern string) *telego.InlineKeyboardMarkup {
	var IKB [][]telego.InlineKeyboardButton
	if TotalPages <= 5 {
		IKB = append(IKB, FullPagination(TotalPages, CurrentPage, CallbackPattern)...)
	} else {
		if CurrentPage <= 3 {
			IKB = append(IKB, LeftPagination(CurrentPage, TotalPages, CallbackPattern)...)
		} else if CurrentPage > TotalPages-3 {
			IKB = append(IKB, RightPagination(CurrentPage, TotalPages, CallbackPattern)...)
		} else {
			IKB = append(IKB, MiddlePagination(CurrentPage, TotalPages, CallbackPattern)...)
		}
	}

	return telegoutil.InlineKeyboardGrid(IKB)
}

func LeftPagination(CurrentPage int, TotalPages int, callbackPattern string) [][]telego.InlineKeyboardButton {
	var row []telego.InlineKeyboardButton
	for number := 1; number <= 5; number++ {
		switch number {
		case CurrentPage:
			row = append(row, AddButton(fmt.Sprintf(SYMBOL_CURRENT_PAGE, number), FormatCallbackPattern(callbackPattern, number)))
		case 4:
			row = append(row, AddButton(fmt.Sprintf(SYMBOL_NEXT_PAGE, number), FormatCallbackPattern(callbackPattern, number)))
		case 5:
			row = append(row, AddButton(fmt.Sprintf(SYMBOL_LAST_PAGE, TotalPages), FormatCallbackPattern(callbackPattern, number)))
		default:
			row = append(row, AddButton(strconv.Itoa(number), FormatCallbackPattern(callbackPattern, number)))
		}
	}
	return [][]telego.InlineKeyboardButton{row}
}

func MiddlePagination(CurrentPage int, TotalPages int, callbackPattern string) [][]telego.InlineKeyboardButton {
	return [][]telego.InlineKeyboardButton{{
		AddButton(fmt.Sprintf(SYMBOL_FIRST_PAGE, 1), FormatCallbackPattern(callbackPattern, 1)),
		AddButton(fmt.Sprintf(SYMBOL_PREVIOUS_PAGE, CurrentPage-1), FormatCallbackPattern(callbackPattern, CurrentPage-1)),
		AddButton(fmt.Sprintf(SYMBOL_CURRENT_PAGE, CurrentPage), FormatCallbackPattern(callbackPattern, CurrentPage)),
		AddButton(fmt.Sprintf(SYMBOL_NEXT_PAGE, CurrentPage+1), FormatCallbackPattern(callbackPattern, CurrentPage+1)),
		AddButton(fmt.Sprintf(SYMBOL_LAST_PAGE, TotalPages), FormatCallbackPattern(callbackPattern, TotalPages)),
	}}
}

func RightPagination(CurrentPage int, TotalPages int, callbackPattern string) [][]telego.InlineKeyboardButton {
	var row []telego.InlineKeyboardButton
	row = append(row, AddButton(fmt.Sprintf(SYMBOL_FIRST_PAGE, 1), FormatCallbackPattern(callbackPattern, 1)))
	row = append(row, AddButton(fmt.Sprintf(SYMBOL_PREVIOUS_PAGE, TotalPages-3), FormatCallbackPattern(callbackPattern, TotalPages-3)))
	for number := TotalPages - 2; number <= TotalPages; number++ {
		if number == CurrentPage {
			row = append(row, AddButton(fmt.Sprintf(SYMBOL_CURRENT_PAGE, number), FormatCallbackPattern(callbackPattern, number)))
		} else {
			row = append(row, AddButton(strconv.Itoa(number), FormatCallbackPattern(callbackPattern, number)))
		}
	}
	return [][]telego.InlineKeyboardButton{row}
}

func FullPagination(TotalPages int, CurrentPage int, callbackPattern string) [][]telego.InlineKeyboardButton {
	var row []telego.InlineKeyboardButton
	for number := 1; number <= TotalPages; number++ {
		if number != CurrentPage {
			row = append(row, AddButton(strconv.Itoa(number), FormatCallbackPattern(callbackPattern, number)))
		} else {
			row = append(row, AddButton(fmt.Sprintf(SYMBOL_CURRENT_PAGE, number), FormatCallbackPattern(callbackPattern, number)))
		}
	}
	return [][]telego.InlineKeyboardButton{row}
}
