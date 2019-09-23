package transaction

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

type bean struct {
	run      interface{}
	runArgs  []interface{}
	back     interface{}
	backArgs []interface{}
}

func NewBean() *bean {
	return &bean{}
}

func (b *bean) Run(f interface{}, args ...interface{}) error {
	if err := b.check(f, args...); err != nil {
		return err
	}

	b.run = f
	b.runArgs = args
	return nil
}

func (b *bean) Back(f interface{}, args ...interface{}) error {
	if err := b.check(f, args...); err != nil {
		return err
	}

	b.back = f
	b.backArgs = args
	return nil
}

func (b *bean) check(f interface{}, args ...interface{}) error {
	fType := reflect.TypeOf(f)
	if fType.Kind() != reflect.Func {
		return newError(errBeanFuncType, fType.Name())
	}

	funcArgsNum := fType.NumIn()
	if funcArgsNum != len(args) {
		return newError(errBeanFuncArgsNum,
			fmt.Sprintf("func args num: %d, args num: %d", funcArgsNum, len(args)))
	}

	for i := 0; i < funcArgsNum; i++ {
		iFuncArgT := fType.In(i)
		iArgT := reflect.TypeOf(args[i])
		if iArgT.Kind() != iFuncArgT.Kind() {
			if iFuncArgT.Kind() == reflect.Interface {
				continue
			}
			return newError(errBeanFuncArgsMatch,
				fmt.Sprintf("index: %d, func arg type: %s, arg type: %s", i, iFuncArgT.Kind(), iArgT.Kind()))
		}
	}

	return nil
}

type pod struct {
	beans []*bean
}

func NewPod() *pod {
	return &pod{beans: make([]*bean, 0, 2)}
}

func (p *pod) Add(bs ...*bean) {
	p.beans = append(p.beans, bs...)
}

func (p *pod) Do() error {
	for i, bean := range p.beans {
		if bean.run == nil {
			continue
		}

		runV := reflect.ValueOf(bean.run)
		var runArgsV []reflect.Value
		for _, v := range bean.runArgs {
			runArgsV = append(runArgsV, reflect.ValueOf(v))
		}

		runOutsV := runV.Call(runArgsV)
		if len(runOutsV) < 1 {
			continue
		}

		err, ok := runOutsV[len(runOutsV)-1].Interface().(error)
		if !ok {
			continue
		}

		for j := i - 1; j >= 0; j-- {
			backV := reflect.ValueOf(p.beans[j].back)
			var backArgsV []reflect.Value
			for _, v := range p.beans[j].backArgs {
				backArgsV = append(backArgsV, reflect.ValueOf(v))
			}

			backOutsV := backV.Call(backArgsV)
			if len(backOutsV) < 1 {
				continue
			}

			errBack, ok := backOutsV[len(backOutsV)-1].Interface().(error)
			if !ok {
				continue
			}

			err = errors.Wrap(err, errBack.Error())
		}

		return err
	}

	return nil
}

type BeanFunc func() error

type funcBean struct {
	run  BeanFunc
	back BeanFunc
}

func NewFuncBean() *funcBean {
	return &funcBean{}
}

func (b *funcBean) Run(f BeanFunc) {
	b.run = f
}

func (b *funcBean) Back(f BeanFunc) {
	b.back = f
}

type funcPod struct {
	beans []*funcBean
}

func NewFuncPod() *funcPod {
	return &funcPod{}
}

func (p *funcPod) Add(bs ...*funcBean) {
	p.beans = append(p.beans, bs...)
}

func (p *funcPod) Do() error {
	for i, bean := range p.beans {
		err := bean.run()
		if err == nil {
			continue
		}

		for j := i - 1; j >= 0; j-- {
			errBack := p.beans[j].back()
			if errBack != nil {
				err = errors.Wrap(err, errBack.Error())
			}
		}

		return err
	}

	return nil
}
