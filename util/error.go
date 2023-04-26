package util

func SameErrorMessage(err, target error) bool {
	if target == nil || err == nil {
		return err == target
	}
	return err.Error() == target.Error()
}
