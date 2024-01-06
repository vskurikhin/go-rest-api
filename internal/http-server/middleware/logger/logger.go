package logger

/*
 * This file was last modified at 2024.01.04 15:07 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * config.go
 * $Id$
 */

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/logger"),
		)

		log.Info("logger middleware enabled")

		handlerFunc := func(responseWriter http.ResponseWriter, request *http.Request) {
			entry := log.With(
				slog.String("method", request.Method),
				slog.String("path", request.URL.Path),
				slog.String("remote_addr", request.RemoteAddr),
				slog.String("user_agent", request.UserAgent()),
				slog.String("request_id", middleware.GetReqID(request.Context())),
			)
			wrapResponseWriter := middleware.NewWrapResponseWriter(responseWriter, request.ProtoMajor)

			tNow := time.Now()
			defer func() {
				entry.Info("request completed",
					slog.Int("status", wrapResponseWriter.Status()),
					slog.Int("bytes", wrapResponseWriter.BytesWritten()),
					slog.String("duration", time.Since(tNow).String()),
				)
			}()

			next.ServeHTTP(wrapResponseWriter, request)
		}

		return http.HandlerFunc(handlerFunc)
	}
}

/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
