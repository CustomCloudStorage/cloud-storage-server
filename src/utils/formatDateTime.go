package utils

import "time"

func FormatDateTime(date time.Time) (string, error) {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return "", ErrInternal.Wrap(err, "load Moscow location")
	}

	formattedDate := date.In(location).Format("2006-01-02 15:04:05 MST")

	return formattedDate, nil
}
