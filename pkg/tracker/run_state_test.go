package tracker

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1beta1"
)

func TestStateString(t *testing.T) {
	stateTests := []struct {
		s    State
		want string
	}{
		{Pending, "Pending"},
		{Failed, "Failed"},
		{Successful, "Successful"},
		{Error, "Error"},
	}

	for _, tt := range stateTests {
		t.Run(tt.want, func(t *testing.T) {
			if v := tt.s.String(); v != tt.want {
				t.Errorf("got %v, want %v", v, tt.want)
			}
		})
	}
}

func TestConditionsToState(t *testing.T) {
	condTests := []struct {
		name string
		c    duckv1.Conditions
		want State
	}{
		{"successful state", conditions(apis.ConditionSucceeded, corev1.ConditionTrue), Successful},
		{"pending state", conditions(apis.ConditionSucceeded, corev1.ConditionUnknown), Pending},
		{"failed state", conditions(apis.ConditionSucceeded, corev1.ConditionFalse), Failed},
		{"default state", conditions(apis.ConditionReady, corev1.ConditionFalse), Pending},
	}

	for _, tt := range condTests {
		t.Run(tt.name, func(t *testing.T) {
			if s := ConditionsToState(tt.c); s != tt.want {
				t.Errorf("ConditionsToState() got %v want %v", s, tt.want)
			}
		})
	}
}

func conditions(s apis.ConditionType, c corev1.ConditionStatus) duckv1.Conditions {
	return duckv1.Conditions{apis.Condition{Type: s, Status: c}}
}
