package modules

func LoginStart() {

}

func LoginConfig() {
	/*
				LoginPwdCost           = GetEnvOrDefault("FINANCE_LOGIN_PWD_COST", "10")             // algorithm cost
				LoginResetDuration     = GetEnvOrDefault("FINANCE_LOGIN_RESET_DURATION", "7")        // in days
				LoginExpiresAtDuration = GetEnvOrDefault("FINANCE_LOGIN_EXPIRES_AT_DURATION", "168") // in hours (7 days)
				LoginAuthSecret        = GetEnvOrDefault("FINANCE_LOGIN_AUTH_SECRET", "")
				LoginEmailHost         = GetEnvOrDefault("FINANCE_EMAIL_HOST", "")
				LoginEmailPort         = GetEnvOrDefault("FINANCE_EMAIL_PORT", "587")
				LoginEmailPwd          = GetEnvOrDefault("FINANCE_EMAIL_PWD", "")
				LoginEmailFrom         = GetEnvOrDefault("FINANCE_EMAIL_FROM", "")
				LoginEmailResetUrl     = GetEnvOrDefault("FINANCE_EMAIL_RESET_URL", "")

				func GetPwdCost() int {
			cost, err := strconv.Atoi(LoginPwdCost)
			if err != nil {
				// TODO: unable to print to default log, might want to send error to another feedback loop
				fmt.Printf("GetPwdCost: unable to parse env var: %s", err)
				return 10
			}
			return cost
		}

		func GetResetDuration() int {
			durationInDays, err := strconv.Atoi(LoginResetDuration)
			if err != nil {
				// TODO: unable to print to default log, might want to send error to another feedback loop
				fmt.Printf("GetResetDuration: unable to parse env var: %s", err)
				return 7
			}
			return durationInDays
		}

		func GetExpiresAtDuration() int {
			durationInHours, err := strconv.Atoi(LoginExpiresAtDuration)
			if err != nil {
				// TODO: unable to print to default log, might want to send error to another feedback loop
				fmt.Printf("GetExpiredAtDuration: unable to parse env var: %s", err)
				return 7
			}
			return durationInHours
		}

		func GetEmailPort() int {
			loginEmailPort, err := strconv.Atoi(LoginEmailPort)
			if err != nil {
				// TODO: unable to print to default log, might want to send error to another feedback loop
				fmt.Printf("GetEmailPort: unable to parse env var: %s", err)
				return 7
			}
			return loginEmailPort
		}

	*/
}

func LoginErrors() {
	/*
		func PasswordValiationError(msg string) ApiError {
		return NewApiError(
			http.StatusBadRequest,
			"Password Validation Error",
			fmt.Sprintf("Invalid password, reason: %s", msg),
			false,
			nil,
		)
	}

	func EmailValidError(msg string) ApiError {
		return NewApiError(
			http.StatusBadRequest,
			"Email Validation Error",
			fmt.Sprintf("Invalid email, reason: %s", msg),
			false,
			nil,
		)
	}

	func ResetTokenInvalidError() ApiError {
		return NewApiError(
			http.StatusBadRequest,
			"Invalid Reset Token",
			"Token missing/expired, please repeat the reset password process",
			false,
			nil,
		)
	}

	func LoginActiveError() ApiError {
		return NewApiError(
			http.StatusBadRequest,
			"Login Inactive",
			"The login user inactive, please contact the site administrator",
			false,
			nil,
		)
	}

	func EmailPasswordComboError() ApiError {
		return NewApiError(
			http.StatusBadRequest,
			"Invalid Email/Password Combination",
			fmt.Sprintf("The email/password entered was not valid, please try again"),
			false,
			nil,
		)
	}

	func DuplicateEmailError(email string) ApiError {
		return NewApiError(
			http.StatusBadRequest,
			"Duplicate Email",
			fmt.Sprintf("The email: %s already exists in this system", email),
			false,
			nil,
		)
	}
	*/
}
