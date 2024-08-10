package configs

import "testing"

func Test_status_validate(t *testing.T) {
	t.Run("empty fields", func(t *testing.T) {
		s := httpServer{}
		err := s.validate()
		if err == nil {
			t.Fatalf("error is expected")
		}
	})

	t.Run("ok", func(t *testing.T) {
		s := httpServer{Port: 9090}
		err := s.validate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func Test_status_initiate(t *testing.T) {
	t.Run("empty fields", func(t *testing.T) {
		s := httpServer{}
		err := s.initiate()
		if err == nil {
			t.Fatal("error is expected")
		}
	})

	t.Run("ok", func(t *testing.T) {
		s := httpServer{Port: 9090}
		err := s.initiate()
		if err != nil {
			t.Fatal("unexpected error", err)
		}
	})

}
