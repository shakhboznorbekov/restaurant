package file

import "os"

func DeleteFiles(names ...string) error {
	for _, v := range names {
		if err := os.Remove(v); err != nil {
			return err
		}
	}

	return nil
}
