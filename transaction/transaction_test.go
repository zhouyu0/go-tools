package transaction

import (
	"fmt"
	"testing"
)

func withSuccess(p bool, args string) error {
	if p {
		fmt.Println(args)
	}

	return nil
}

func withError(p bool, args string) error {
	if p {
		fmt.Println(args)
	}

	return fmt.Errorf("%v", args)
}

func TestBean_Run(t *testing.T) {
	b := NewBean()
	err := b.Run(withSuccess, true, "test run")
	if err != nil {
		t.Error(err)
	}
}

func TestBean_Back(t *testing.T) {
	b := NewBean()
	err := b.Back(withSuccess, true, "test back")
	if err != nil {
		t.Error(err)
	}
}

func TestPod_Add(t *testing.T) {
	b1 := NewBean()
	err := b1.Run(withSuccess, true, "test run")
	if err != nil {
		t.Error(err)
	}
	err = b1.Back(withSuccess, true, "test back")
	if err != nil {
		t.Error(err)
	}

	b2 := NewBean()
	err = b2.Run(withSuccess, true, "test run")
	if err != nil {
		t.Error(err)
	}
	err = b2.Back(withSuccess, true, "test back")
	if err != nil {
		t.Error(err)
	}

	p := NewPod()
	p.Add(b1, b2)
}

func TestPod_Do(t *testing.T) {
	b1 := NewBean()
	err := b1.Run(withSuccess, true, "run1")
	if err != nil {
		t.Error(err)
	}
	err = b1.Back(withSuccess, true, "back1")
	if err != nil {
		t.Error(err)
	}

	b2 := NewBean()
	err = b2.Run(withError, true, "run2")
	if err != nil {
		t.Error(err)
	}
	err = b2.Back(withSuccess, true, "back2")
	if err != nil {
		t.Error(err)
	}

	b3 := NewBean()
	err = b3.Run(withSuccess, true, "run3")
	if err != nil {
		t.Error(err)
	}
	err = b3.Back(withSuccess, true, "back3")
	if err != nil {
		t.Error(err)
	}

	p1 := NewPod()
	p1.Add(b1, b2)
	err = p1.Do()
	if err != nil {
		if err.Error() != "run2" {
			t.Error(err)
		}
	}

	p2 := NewPod()
	p2.Add(b1, b2, b3)
	err = p2.Do()
	if err != nil {
		if err.Error() != "run2" {
			t.Error(err)
		}
	}
}

func TestFuncBean_Run(t *testing.T) {
	b := NewFuncBean()
	b.Run(func() error {
		return withSuccess(true, "test run")
	})
}

func TestFuncBean_Back(t *testing.T) {
	b := NewFuncBean()
	b.Back(func() error {
		return withSuccess(true, "test back")
	})
}

func TestFuncPod_Add(t *testing.T) {
	b1 := NewFuncBean()
	b1.Run(func() error {
		return withSuccess(true, "run1")
	})
	b1.Back(func() error {
		return withSuccess(true, "back1")
	})

	b2 := NewFuncBean()
	b2.Run(func() error {
		return withSuccess(true, "run2")
	})
	b2.Back(func() error {
		return withSuccess(true, "back2")
	})

	p := NewFuncPod()
	p.Add(b1, b2)
}

func TestFuncPod_Do(t *testing.T) {
	b1 := NewFuncBean()
	b1.Run(func() error {
		return withSuccess(true, "run1")
	})
	b1.Back(func() error {
		return withSuccess(true, "back1")
	})

	b2 := NewFuncBean()
	b2.Run(func() error {
		return withError(true, "run2")
	})
	b2.Back(func() error {
		return withSuccess(true, "back2")
	})

	p := NewFuncPod()
	p.Add(b1, b2)

	err := p.Do()
	if err != nil {
		if err.Error() != "run2" {
			t.Error(err)
		}
	}
}

func BenchmarkPod_Do(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b1 := NewBean()
		err := b1.Run(withSuccess, false, "run1")
		if err != nil {
			b.Error(err)
		}
		err = b1.Back(withSuccess, false, "back1")
		if err != nil {
			b.Error(err)
		}

		b2 := NewBean()
		err = b2.Run(withError, false, "run2")
		if err != nil {
			b.Error(err)
		}
		err = b2.Back(withSuccess, false, "back2")
		if err != nil {
			b.Error(err)
		}

		p := NewPod()
		p.Add(b1, b2)
		err = p.Do()
		if err != nil {
			if err.Error() != "run2" {
				b.Error(err)
			}
		}
	}
}

func BenchmarkFuncPod_Do(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b1 := NewFuncBean()
		b1.Run(func() error {
			return withSuccess(false, "run1")
		})
		b1.Back(func() error {
			return withSuccess(false, "back1")
		})

		b2 := NewFuncBean()
		b2.Run(func() error {
			return withError(false, "run2")
		})
		b2.Back(func() error {
			return withSuccess(false, "back2")
		})

		p := NewFuncPod()
		p.Add(b1, b2)

		err := p.Do()
		if err != nil {
			if err.Error() != "run2" {
				b.Error(err)
			}
		}
	}
}
