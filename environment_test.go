package environment

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type envTest struct {
	number uint64
}

func TestEnvironment(t *testing.T) {
	e := New()

	st := &envTest{
		number: 1988,
	}

	e.Set("my-struct", st)
	e.Set("my-string", "value")
	e.Set("my-uint64", uint64(4576))

	structVal := e.Get("my-struct")
	require.Equal(t, uint64(1988), structVal.(*envTest).number)

	stringVal := e.Get("my-string")
	require.Equal(t, "value", stringVal)

	stringUint64 := e.Get("my-uint64")
	require.Equal(t, uint64(4576), stringUint64.(uint64))

	t.Run("Clone", func(t *testing.T) {
		clone := e.Clone()
		require.NotSame(t, e, clone) // Ensure it's a different instance
		require.Equal(t, e.m, clone.m)
	})
}

func TestEnvironmentSingleton(t *testing.T) {
	// Get singleton instance and set a value
	env1 := GetInstance()
	env1.Set("key1", "value1")

	// Get the singleton instance again
	env2 := GetInstance()

	// Verify both instances are the same
	require.Same(t, env1, env2)

	// Verify the value set in env1 is accessible in env2
	val := env2.Get("key1")
	require.Equal(t, "value1", val)

	// Set a new value in env2 and verify in env1
	env2.Set("key2", "value2")
	val2 := env1.Get("key2")
	require.Equal(t, "value2", val2)
}
