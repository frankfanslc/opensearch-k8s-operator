package helpers

import (
	"context"
	"fmt"
	"reflect"

	sts "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	opsterv1 "opensearch.opster.io/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OpenSearchReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

func ContainsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false

}

func (r *OpenSearchReconciler) UpdateResource(ctx context.Context, instance *sts.StatefulSet) error {
	err := r.Update(ctx, instance)
	if err != nil {
		fmt.Println(err, "Cannot update resource")
		r.Recorder.Event(instance, "Warning", "Cannot update resource", "Cannot update resource")
		return err
	}
	return nil
}

func GetField(v *sts.StatefulSetSpec, field string) interface{} {

	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field).Interface()
	return f
}

//func setField(v *sts.StatefulSetSpec, field string) reflect.Value {
//
//	//	reflect.ValueOf(v).Elem().FieldByName(field).SetString("sss")
//
//	r := reflect.ValueOf(v)
//	ty := r.Type()
//	fmt.Println(ty)
//	f := reflect.Indirect(r).FieldByName(field)
//	return f
//}

func getNamesInStruct(inter interface{}) []string {
	rv := reflect.Indirect(reflect.ValueOf(inter))

	var names []string

	for i := 0; i < rv.NumField(); i++ {
		x := rv.Type().Field(i).Name
		names = append(names, x)
	}
	return names
}

func RemoveIt(ss opsterv1.ComponentStatus, ssSlice []opsterv1.ComponentStatus) []opsterv1.ComponentStatus {
	for idx, v := range ssSlice {
		if v == ss {
			return append(ssSlice[0:idx], ssSlice[idx+1:]...)
		}
	}
	return ssSlice
}

func Remove(slice []opsterv1.ComponentStatus, componenet string) []opsterv1.ComponentStatus {
	emptyStatus := opsterv1.ComponentStatus{
		Component:   "",
		Status:      "",
		Description: "",
	}
	for i := 0; i < len(slice); i++ {
		if slice[i].Component == componenet {
			slice[i] = slice[len(slice)-1]    // Copy last element to index i.
			slice[len(slice)-1] = emptyStatus // Erase last element (write zero value).
			slice = slice[:len(slice)-1]
		}

	}
	return slice
}
