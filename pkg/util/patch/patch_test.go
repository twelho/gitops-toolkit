package patch

import (
	"bytes"
	"testing"

	api "github.com/weaveworks/gitops-toolkit/cmd/sample-app/apis/sample"
	"github.com/weaveworks/gitops-toolkit/cmd/sample-app/apis/sample/scheme"
	"github.com/weaveworks/gitops-toolkit/pkg/runtime"
)

var (
	basebytes = []byte(`
{
	"kind": "Car",
	"apiVersion": "sample-app.weave.works/v1alpha1",
	"metadata": {
	  "name": "foo",
	  "uid": "0123456789101112"
	},
	"spec": {
	  "engine": "foo",
	  "brand": "bar"
	}
}`)
	overlaybytes = []byte(`{"status":{"speed":24.7}}`)
)

var carGVK = api.SchemeGroupVersion.WithKind("Car")
var p = NewPatcher(scheme.Serializer)

func TestCreatePatch(t *testing.T) {
	car := &api.Car{
		Spec: api.CarSpec{
			Engine: "foo",
			Brand:  "bar",
		},
	}
	car.SetGroupVersionKind(carGVK)
	b, err := p.Create(car, func(obj runtime.Object) error {
		car2 := obj.(*api.Car)
		car2.Status.Speed = 24.7
		return nil
	})
	if !bytes.Equal(b, overlaybytes) {
		t.Error(string(b), err, car.Status.Speed)
	}
}

func TestApplyPatch(t *testing.T) {
	result, err := p.Apply(basebytes, overlaybytes, carGVK)
	if err != nil {
		t.Fatal(err)
	}
	car := &api.Car{}
	if err := scheme.Serializer.DecodeInto(result, car); err != nil {
		t.Fatal(err)
	}
}
