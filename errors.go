// Copyright 2017-present Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package gramework

import (
	"errors"
)

var (
	// ErrTLSNoEmails occurs when no emails provided but user tries to use AutoTLS features
	ErrTLSNoEmails = errors.New("auto tls: no emails provided")

	// ErrArgNotFound used when no route argument is found
	ErrArgNotFound = errors.New("undefined argument")

	// ErrInvalidGQLRequest used in DecodeGQL
	ErrInvalidGQLRequest = errors.New("invalid gql request")
)
