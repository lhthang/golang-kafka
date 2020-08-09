
package error

func GetErrorMessage(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}
