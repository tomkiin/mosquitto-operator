package v1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func (in *Mosquitto) DeepCopyInto(out *Mosquitto) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

func (in *Mosquitto) DeepCopy() *Mosquitto {
	if in == nil {
		return nil
	}
	out := new(Mosquitto)
	in.DeepCopyInto(out)
	return out
}

func (in *Mosquitto) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *MosquittoList) DeepCopyInto(out *MosquittoList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Mosquitto, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

func (in *MosquittoList) DeepCopy() *MosquittoList {
	if in == nil {
		return nil
	}
	out := new(MosquittoList)
	in.DeepCopyInto(out)
	return out
}

func (in *MosquittoList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

func (in *MosquittoSpec) DeepCopyInto(out *MosquittoSpec) {
	*out = *in
}

func (in *MosquittoSpec) DeepCopy() *MosquittoSpec {
	if in == nil {
		return nil
	}
	out := new(MosquittoSpec)
	in.DeepCopyInto(out)
	return out
}

func (in *MosquittoStatus) DeepCopyInto(out *MosquittoStatus) {
	*out = *in
}

func (in *MosquittoStatus) DeepCopy() *MosquittoStatus {
	if in == nil {
		return nil
	}
	out := new(MosquittoStatus)
	in.DeepCopyInto(out)
	return out
}
