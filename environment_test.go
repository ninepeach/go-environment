package environment

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

type envTest struct {
	number uint64
}

func TestEnvironment(t *testing.T) {
	e := New()
	st := &envTest{number: 1988}
	e.Set("my-struct", st)
	e.Set("my-string", "value")
	e.Set("my-uint64", uint64(4576))
	require.Equal(t, uint64(1988), e.Get("my-struct").(*envTest).number)
	require.Equal(t, "value", e.Get("my-string"))
	require.Equal(t, uint64(4576), e.Get("my-uint64"))

	t.Run("Lookup and Delete", func(t *testing.T) {
		_, ok := e.Lookup("missing")
		require.False(t, ok)
		require.Nil(t, e.Get("missing"))
		e.Set("nil", nil)
		value, ok := e.Lookup("nil")
		require.True(t, ok)
		require.Nil(t, value)
		e.Delete("nil")
		_, ok = e.Lookup("nil")
		require.False(t, ok)
		e.Delete("nil")
	})

	t.Run("Clone", func(t *testing.T) {
		clone := e.Clone()
		require.NotSame(t, e, clone)
		require.Equal(t, e.m, clone.m)
		clone.Set("my-string", "changed")
		require.Equal(t, "value", e.Get("my-string"))
		clone.Delete("my-uint64")
		require.Equal(t, uint64(4576), e.Get("my-uint64"))
		require.Same(t, e.Get("my-struct"), clone.Get("my-struct"))
	})
}

func TestZeroValue(t *testing.T) {
	var e Environment
	require.Nil(t, e.Get("missing"))
	e.Delete("missing")
	e.Set("key", 42)
	require.Equal(t, 42, e.Get("key"))
	require.Equal(t, 42, e.Clone().Get("key"))
}

func TestEnvironmentConcurrent(t *testing.T) {
	e := New()
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				e.Set("key", i)
				e.Get("key")
				e.Lookup("key")
				e.Clone()
				e.Delete("key")
			}
		}(i)
	}
	wg.Wait()
}

func TestGetFieldValue(t *testing.T) {
	type config struct {
		Name   string
		secret string
	}
	c := &config{Name: "ok", secret: "private"}
	pc := &c
	value, err := GetFieldValue(pc, "Name")
	require.NoError(t, err)
	require.Equal(t, "ok", value)

	for _, tc := range []struct {
		name string
		obj  interface{}
		key  string
	}{
		{"nil interface", nil, "Name"},
		{"nil pointer", (*config)(nil), "Name"},
		{"nil nested pointer", (**config)(nil), "Name"},
		{"not struct", 42, "Name"},
		{"missing field", c, "Missing"},
		{"unexported field", c, "secret"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				_, err := GetFieldValue(tc.obj, tc.key)
				require.Error(t, err)
			})
		})
	}
}

func TestEnvironmentSingleton(t *testing.T) {
	env1 := GetInstance()
	env1.Set("key1", "value1")
	env2 := GetInstance()
	require.Same(t, env1, env2)
	require.Equal(t, "value1", env2.Get("key1"))
	env2.Set("key2", "value2")
	require.Equal(t, "value2", env1.Get("key2"))
}

func TestZeroValueConcurrent(t *testing.T) {
	var e Environment
	const workers = 32
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e.Set(fmt.Sprintf("key-%d", i), i)
		}(i)
	}
	wg.Wait()
	for i := 0; i < workers; i++ {
		value, ok := e.Lookup(fmt.Sprintf("key-%d", i))
		require.True(t, ok)
		require.Equal(t, i, value)
	}
}

func TestGetFieldValueEmbedded(t *testing.T) {
	type Inner struct {
		Name string
	}
	type Embedded struct {
		*Inner
	}
	type A struct {
		Name string
	}
	type B struct {
		Name string
	}
	type Ambiguous struct {
		A
		B
	}
	value, err := GetFieldValue(Embedded{Inner: &Inner{Name: "ok"}}, "Name")
	require.NoError(t, err)
	require.Equal(t, "ok", value)

	require.NotPanics(t, func() {
		_, err := GetFieldValue(Embedded{}, "Name")
		require.Error(t, err)
	})
	require.NotPanics(t, func() {
		_, err := GetFieldValue(Ambiguous{}, "Name")
		require.Error(t, err)
	})
}
