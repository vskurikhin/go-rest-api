package sl

import "log/slog"

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
