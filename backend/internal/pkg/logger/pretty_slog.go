package logger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"

	"github.com/fatih/color"
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	opts  PrettyHandlerOptions
	l     *log.Logger
	attrs []slog.Attr
	group string
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.SlogOpts.Level.Level()
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	// Собираем все атрибуты (включая те, что были добавлены через WithAttrs)
	fields := make(map[string]interface{}, r.NumAttrs()+len(h.attrs))

	// Добавляем атрибуты из WithAttrs
	for _, attr := range h.attrs {
		fields[attr.Key] = attr.Value.Any()
	}

	// Добавляем атрибуты из текущей записи
	r.Attrs(func(a slog.Attr) bool {
		key := a.Key
		if h.group != "" {
			key = h.group + "." + key
		}
		fields[key] = a.Value.Any()
		return true
	})

	// Сериализуем в JSON с отступами
	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[2006-01-02 15:04:05.000]")
	msg := color.CyanString(r.Message)

	if len(b) > 2 { // Проверяем, что JSON не пустой ("{}")
		h.l.Println(timeStr, level, msg, color.HiBlackString(string(b)))
	} else {
		h.l.Println(timeStr, level, msg)
	}

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Создаем новый обработчик с добавленными атрибутами
	return &PrettyHandler{
		opts:  h.opts,
		l:     h.l,
		attrs: append(h.attrs[:len(h.attrs):len(h.attrs)], attrs...), // Копируем существующие атрибуты и добавляем новые
		group: h.group,
	}
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	// Создаем новый обработчик с указанной группой
	return &PrettyHandler{
		opts:  h.opts,
		l:     h.l,
		attrs: h.attrs,
		group: joinGroup(h.group, name), // Объединяем группы, если они вложенные
	}
}

// joinGroup объединяет имена групп с учетом вложенности
func joinGroup(prefix, name string) string {
	if prefix == "" {
		return name
	}
	return prefix + "." + name
}

func NewPrettyHandler(
	out io.Writer,
	opts PrettyHandlerOptions,
) *PrettyHandler {
	return &PrettyHandler{
		opts: opts,
		l:    log.New(out, "", 0),
	}
}
