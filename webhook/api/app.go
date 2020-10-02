package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type App struct {
}

func (app *App) HandleMutate(w http.ResponseWriter, r *http.Request) {
	admissionReview := &admissionv1.AdmissionReview{}

	// read the AdmissionReview from the request json body
	err := readJSON(r, admissionReview)
	if err != nil {
		app.HandleError(w, r, err)
		return
	}

	// unmarshal the pod from the AdmissionRequest
	pod := &corev1.Pod{}
	if err := json.Unmarshal(admissionReview.Request.Object.Raw, pod); err != nil {
		app.HandleError(w, r, fmt.Errorf("unmarshal to pod: %v", err))
		return
	}

	// add the volume to the pod
	pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
		Name: "hello-volume",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "hello-configmap",
				},
			},
		},
	})

	// add volume mount to all containers in the pod
	for i := 0; i < len(pod.Spec.Containers); i++ {
		pod.Spec.Containers[i].VolumeMounts = append(pod.Spec.Containers[i].VolumeMounts, corev1.VolumeMount{
			Name:      "hello-volume",
			MountPath: "/etc/config",
		})
	}

	containersBytes, err := json.Marshal(&pod.Spec.Containers)
	if err != nil {
		app.HandleError(w, r, fmt.Errorf("marshall containers: %v", err))
		return
	}

	volumesBytes, err := json.Marshal(&pod.Spec.Volumes)
	if err != nil {
		app.HandleError(w, r, fmt.Errorf("marshall volumes: %v", err))
		return
	}

	// build json patch
	patch := []JSONPatchEntry{
		JSONPatchEntry{
			OP:    "add",
			Path:  "/metadata/labels/hello-added",
			Value: []byte(`"OK"`),
		},
		JSONPatchEntry{
			OP:    "replace",
			Path:  "/spec/containers",
			Value: containersBytes,
		},
		JSONPatchEntry{
			OP:    "replace",
			Path:  "/spec/volumes",
			Value: volumesBytes,
		},
	}

	patchBytes, err := json.Marshal(&patch)
	if err != nil {
		app.HandleError(w, r, fmt.Errorf("marshall jsonpatch: %v", err))
		return
	}

	patchType := admissionv1.PatchTypeJSONPatch

	// build admission response
	admissionResponse := &admissionv1.AdmissionResponse{
		UID:       admissionReview.Request.UID,
		Allowed:   true,
		Patch:     patchBytes,
		PatchType: &patchType,
	}

	respAdmissionReview := &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: admissionResponse,
	}

	jsonOk(w, &respAdmissionReview)
}

type JSONPatchEntry struct {
	OP    string          `json:"op"`
	Path  string          `json:"path"`
	Value json.RawMessage `json:"value,omitempty"`
}
