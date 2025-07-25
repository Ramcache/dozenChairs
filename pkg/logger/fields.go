package logger

import "go.uber.org/zap"

func ZapAny(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

func ZapStr(key, val string) zap.Field {
	return zap.String(key, val)
}

func ZapInt(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func ZapBool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

func ZapErr(err error) zap.Field {
	return zap.Error(err)
}
