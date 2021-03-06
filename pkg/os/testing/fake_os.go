/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testing

import (
	"io"
	"os"
	"sync"

	"golang.org/x/net/context"

	osInterface "github.com/kubernetes-incubator/cri-containerd/pkg/os"
)

// CalledDetail is the struct contains called function name and arguments.
type CalledDetail struct {
	// Name of the function called.
	Name string
	// Arguments of the function called.
	Arguments []interface{}
}

// FakeOS mocks out certain OS calls to avoid perturbing the filesystem
// If a member of the form `*Fn` is set, that function will be called in place
// of the real call.
type FakeOS struct {
	sync.Mutex
	MkdirAllFn  func(string, os.FileMode) error
	RemoveAllFn func(string) error
	OpenFifoFn  func(context.Context, string, int, os.FileMode) (io.ReadWriteCloser, error)
	StatFn      func(string) (os.FileInfo, error)
	CopyFileFn  func(string, string, os.FileMode) error
	WriteFileFn func(string, []byte, os.FileMode) error
	MountFn     func(source string, target string, fstype string, flags uintptr, data string) error
	UnmountFn   func(target string, flags int) error
	calls       []CalledDetail
	errors      map[string]error
}

var _ osInterface.OS = &FakeOS{}

// getError get error for call
func (f *FakeOS) getError(op string) error {
	f.Lock()
	defer f.Unlock()
	err, ok := f.errors[op]
	if ok {
		delete(f.errors, op)
		return err
	}
	return nil
}

// InjectError inject error for call
func (f *FakeOS) InjectError(fn string, err error) {
	f.Lock()
	defer f.Unlock()
	f.errors[fn] = err
}

// InjectErrors inject errors for calls
func (f *FakeOS) InjectErrors(errs map[string]error) {
	f.Lock()
	defer f.Unlock()
	for fn, err := range errs {
		f.errors[fn] = err
	}
}

// ClearErrors clear errors for call
func (f *FakeOS) ClearErrors() {
	f.Lock()
	defer f.Unlock()
	f.errors = make(map[string]error)
}

func (f *FakeOS) appendCalls(name string, args ...interface{}) {
	f.Lock()
	defer f.Unlock()
	f.calls = append(f.calls, CalledDetail{Name: name, Arguments: args})
}

// GetCalls get detail of calls.
func (f *FakeOS) GetCalls() []CalledDetail {
	f.Lock()
	defer f.Unlock()
	return append([]CalledDetail{}, f.calls...)
}

// NewFakeOS creates a FakeOS.
func NewFakeOS() *FakeOS {
	return &FakeOS{
		errors: make(map[string]error),
	}
}

// MkdirAll is a fake call that invokes MkdirAllFn or just returns nil.
func (f *FakeOS) MkdirAll(path string, perm os.FileMode) error {
	f.appendCalls("MkdirAll", path, perm)
	if err := f.getError("MkdirAll"); err != nil {
		return err
	}

	if f.MkdirAllFn != nil {
		return f.MkdirAllFn(path, perm)
	}
	return nil
}

// RemoveAll is a fake call that invokes RemoveAllFn or just returns nil.
func (f *FakeOS) RemoveAll(path string) error {
	f.appendCalls("RemoveAll", path)
	if err := f.getError("RemoveAll"); err != nil {
		return err
	}

	if f.RemoveAllFn != nil {
		return f.RemoveAllFn(path)
	}
	return nil
}

// OpenFifo is a fake call that invokes OpenFifoFn or just returns nil.
func (f *FakeOS) OpenFifo(ctx context.Context, fn string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	f.appendCalls("OpenFifo", ctx, fn, flag, perm)
	if err := f.getError("OpenFifo"); err != nil {
		return nil, err
	}

	if f.OpenFifoFn != nil {
		return f.OpenFifoFn(ctx, fn, flag, perm)
	}
	return nil, nil
}

// Stat is a fake call that invokes StatFn or just return nil.
func (f *FakeOS) Stat(name string) (os.FileInfo, error) {
	f.appendCalls("Stat", name)
	if err := f.getError("Stat"); err != nil {
		return nil, err
	}

	if f.StatFn != nil {
		return f.StatFn(name)
	}
	return nil, nil
}

// CopyFile is a fake call that invokes CopyFileFn or just return nil.
func (f *FakeOS) CopyFile(src, dest string, perm os.FileMode) error {
	f.appendCalls("CopyFile", src, dest, perm)
	if err := f.getError("CopyFile"); err != nil {
		return err
	}

	if f.CopyFileFn != nil {
		return f.CopyFileFn(src, dest, perm)
	}
	return nil
}

// WriteFile is a fake call that invokes WriteFileFn or just return nil.
func (f *FakeOS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	f.appendCalls("WriteFile", filename, data, perm)
	if err := f.getError("WriteFile"); err != nil {
		return err
	}

	if f.WriteFileFn != nil {
		return f.WriteFileFn(filename, data, perm)
	}
	return nil
}

// Mount is a fake call that invokes MountFn or just return nil.
func (f *FakeOS) Mount(source string, target string, fstype string, flags uintptr, data string) error {
	f.appendCalls("Mount", source, target, fstype, flags, data)
	if err := f.getError("Mount"); err != nil {
		return err
	}

	if f.MountFn != nil {
		return f.MountFn(source, target, fstype, flags, data)
	}
	return nil
}

// Unmount is a fake call that invokes UnmountFn or just return nil.
func (f *FakeOS) Unmount(target string, flags int) error {
	f.appendCalls("Unmount", target, flags)
	if err := f.getError("Unmount"); err != nil {
		return err
	}

	if f.UnmountFn != nil {
		return f.UnmountFn(target, flags)
	}
	return nil
}
