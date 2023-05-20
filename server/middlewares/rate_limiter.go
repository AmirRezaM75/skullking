package middlewares

import "context"

type (
	RateLimiterConfig struct {
		Rate int
		IdentifierExtractor Extractor
	}

	RateLimiter struct {
		// RateLimiterService
		config RateLimiterConfig
	}

	Extractor func(context context.Context) (string, error)
)
