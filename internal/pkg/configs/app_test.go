package configs

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"
// )

// func Test_app_initiate(t *testing.T) {
// 	t.Run("empty fields", func(t *testing.T) {
// 		a := app{}
// 		err := a.initiate()
// 		require.NotNil(t, err)
// 	})

// 	t.Run("ok", func(t *testing.T) {
// 		a := app{Environment: "testing", AgentID: 1}
// 		err := a.initiate()
// 		require.Nil(t, err)
// 	})
// }

// func Test_app_validate(t *testing.T) {
// 	t.Run("empty fields", func(t *testing.T) {
// 		a := app{}
// 		err := a.validate()
// 		require.NotNil(t, err)
// 	})

// 	t.Run("empty agent_id", func(t *testing.T) {
// 		a := app{Environment: "testing"}
// 		err := a.validate()
// 		require.NotNil(t, err)
// 	})

// 	t.Run("ok", func(t *testing.T) {
// 		a := app{Environment: "testing", AgentID: 1}
// 		err := a.validate()
// 		require.Nil(t, err)
// 	})
// }
