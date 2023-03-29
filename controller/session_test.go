package controller

import (
	"gitea.slauson.io/slausonio/iam-ms/service"
	"github.com/gin-gonic/gin"
	"reflect"
	"testing"
)

func TestNewSessionController(t *testing.T) {
	tests := []struct {
		name string
		want *SessionController
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSessionController(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSessionController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionController_CreateEmailSession(t *testing.T) {
	type fields struct {
		s service.IamSessionService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &SessionController{
				s: tt.fields.s,
			}
			sc.CreateEmailSession(tt.args.c)
		})
	}
}

func TestSessionController_DeleteSession(t *testing.T) {
	type fields struct {
		s service.IamSessionService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &SessionController{
				s: tt.fields.s,
			}
			sc.DeleteSession(tt.args.c)
		})
	}
}
